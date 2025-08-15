package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"telegram-emulator/internal/pkg/logger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Server представляет WebSocket сервер
type Server struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	logger     *zap.Logger
}

// Client представляет WebSocket клиента
type Client struct {
	server  *Server
	conn    *websocket.Conn
	send    chan []byte
	userID  string
	logger  *zap.Logger
}

// Message представляет сообщение WebSocket
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// NewServer создает новый WebSocket сервер
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger.GetLogger(),
	}
}

// Start запускает WebSocket сервер
func (s *Server) Start() {
	s.logger.Info("WebSocket сервер запущен")
	
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			s.logger.Info("Клиент подключен", zap.String("user_id", client.userID))

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mutex.Unlock()
			s.logger.Info("Клиент отключен", zap.String("user_id", client.userID))

		case message := <-s.broadcast:
			s.mutex.RLock()
			for client := range s.clients {
				select {
				case client.send <- s.serializeMessage(message):
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mutex.RUnlock()
		}
	}
}

// Broadcast отправляет сообщение всем подключенным клиентам
func (s *Server) Broadcast(messageType string, data interface{}) {
	message := &Message{
		Type: messageType,
		Data: data,
	}
	s.broadcast <- message
}

// BroadcastToUser отправляет сообщение конкретному пользователю
func (s *Server) BroadcastToUser(userID, messageType string, data interface{}) {
	message := &Message{
		Type: messageType,
		Data: data,
	}
	
	s.mutex.RLock()
	for client := range s.clients {
		if client.userID == userID {
			select {
			case client.send <- s.serializeMessage(message):
			default:
				close(client.send)
				delete(s.clients, client)
			}
		}
	}
	s.mutex.RUnlock()
}

// serializeMessage сериализует сообщение в JSON
func (s *Server) serializeMessage(message *Message) []byte {
	data, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Ошибка сериализации сообщения", zap.Error(err))
		return nil
	}
	return data
}

// HandleWebSocket обрабатывает WebSocket подключения
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из query параметров
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id обязателен", http.StatusBadRequest)
		return
	}

	// Настройка WebSocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // В продакшене нужно настроить проверку origin
		},
	}

	// Обновляем HTTP соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("Ошибка обновления соединения", zap.Error(err))
		return
	}

	// Создаем нового клиента
	client := &Client{
		server: s,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		logger: logger.GetLogger(),
	}

	// Регистрируем клиента
	s.register <- client

	// Запускаем горутины для чтения и записи
	go client.writePump()
	go client.readPump()
}

// readPump читает сообщения от клиента
func (c *Client) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("Ошибка чтения WebSocket", zap.Error(err))
			}
			break
		}

		// Обрабатываем входящие сообщения
		c.handleMessage(message)
	}
}

// writePump отправляет сообщения клиенту
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage обрабатывает входящие сообщения от клиента
func (c *Client) handleMessage(message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		c.logger.Error("Ошибка парсинга сообщения", zap.Error(err))
		return
	}

	switch msg.Type {
	case "subscribe":
		c.handleSubscribe(msg.Data)
	case "typing":
		c.handleTyping(msg.Data)
	case "ping":
		c.handlePing()
	default:
		c.logger.Warn("Неизвестный тип сообщения", zap.String("type", msg.Type))
	}
}

// handleSubscribe обрабатывает подписку на события
func (c *Client) handleSubscribe(data interface{}) {
	c.logger.Info("Клиент подписался на события", 
		zap.String("user_id", c.userID),
		zap.Any("events", data))
	
	// Отправляем подтверждение подписки
	response := &Message{
		Type: "subscribed",
		Data: map[string]interface{}{
			"user_id": c.userID,
			"events":  data,
		},
	}
	
	select {
	case c.send <- c.server.serializeMessage(response):
	default:
		c.logger.Warn("Не удалось отправить подтверждение подписки")
	}
}

// handleTyping обрабатывает событие печати
func (c *Client) handleTyping(data interface{}) {
	// Отправляем событие печати другим пользователям в чате
	c.server.Broadcast("typing", map[string]interface{}{
		"user_id": c.userID,
		"data":    data,
	})
}

// handlePing обрабатывает ping сообщения
func (c *Client) handlePing() {
	response := &Message{
		Type: "pong",
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	}
	
	select {
	case c.send <- c.server.serializeMessage(response):
	default:
		c.logger.Warn("Не удалось отправить pong")
	}
}

// GetConnectedUsers возвращает список подключенных пользователей
func (s *Server) GetConnectedUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	users := make([]string, 0, len(s.clients))
	for client := range s.clients {
		users = append(users, client.userID)
	}
	
	return users
}

// IsUserConnected проверяет, подключен ли пользователь
func (s *Server) IsUserConnected(userID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for client := range s.clients {
		if client.userID == userID {
			return true
		}
	}
	
	return false
}
