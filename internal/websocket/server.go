package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"telegram-emulator/internal/models"
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
	messageManager MessageManagerInterface // MessageManager для обработки сообщений
	botManager interface{} // BotManager для обработки callback query
}

// Client представляет WebSocket клиента
type Client struct {
	server  *Server
	conn    *websocket.Conn
	send    chan []byte
	userID  int64
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

// SetMessageManager устанавливает MessageManager для обработки сообщений
func (s *Server) SetMessageManager(messageManager MessageManagerInterface) {
	s.messageManager = messageManager
}

// SetBotManager устанавливает BotManager для обработки callback query
func (s *Server) SetBotManager(botManager interface{}) {
	s.botManager = botManager
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
			s.logger.Info("Клиент подключен", zap.Int64("user_id", client.userID))

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mutex.Unlock()
			s.logger.Info("Клиент отключен", zap.Int64("user_id", client.userID))

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
func (s *Server) BroadcastToUser(userID int64, messageType string, data interface{}) {
	s.logger.Info("BroadcastToUser вызван", 
		zap.Int64("user_id", userID),
		zap.String("message_type", messageType))
	
	message := &Message{
		Type: messageType,
		Data: data,
	}
	
	s.mutex.RLock()
	clientCount := 0
	for client := range s.clients {
		if client.userID == userID {
			clientCount++
			s.logger.Info("Найден клиент для отправки", 
				zap.Int64("user_id", userID),
				zap.String("message_type", messageType))
			select {
			case client.send <- s.serializeMessage(message):
				s.logger.Info("Сообщение отправлено клиенту", 
					zap.Int64("user_id", userID),
					zap.String("message_type", messageType))
			default:
				s.logger.Warn("Канал клиента переполнен, закрываем соединение", zap.Int64("user_id", userID))
				close(client.send)
				delete(s.clients, client)
			}
		}
	}
	s.mutex.RUnlock()
	
	s.logger.Info("BroadcastToUser завершен", 
		zap.Int64("user_id", userID),
		zap.String("message_type", messageType),
		zap.Int("clients_found", clientCount))
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
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id обязателен", http.StatusBadRequest)
		return
	}
	
	// Конвертируем userID в int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат user_id", http.StatusBadRequest)
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
				c.logger.Error("Ошибка получения writer", zap.Error(err))
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				c.logger.Error("Ошибка закрытия writer", zap.Error(err))
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
	case "send_message":
		c.handleSendMessage(msg.Data)
	case "callback_query":
		c.handleCallbackQuery(msg.Data)
	default:
		c.logger.Warn("Неизвестный тип сообщения", zap.String("type", msg.Type))
	}
}

// handleSubscribe обрабатывает подписку на события
func (c *Client) handleSubscribe(data interface{}) {
	c.logger.Info("Клиент подписался на события", 
		zap.Int64("user_id", c.userID),
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

// handleSendMessage обрабатывает отправку сообщений
func (c *Client) handleSendMessage(data interface{}) {
	// Преобразуем data в map для извлечения параметров
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		c.logger.Error("Неверный формат данных для send_message")
		return
	}

	// Обрабатываем chat_id как float64 или строку
	var chatID int64
	if chatIDFloat, ok := dataMap["chat_id"].(float64); ok {
		chatID = int64(chatIDFloat)
	} else if chatIDStr, ok := dataMap["chat_id"].(string); ok {
		chatIDParsed, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			c.logger.Error("Неверный формат chat_id", zap.Error(err))
			return
		}
		chatID = chatIDParsed
	} else {
		c.logger.Error("Отсутствует или неверный формат chat_id в send_message")
		return
	}

	text, ok := dataMap["text"].(string)
	if !ok {
		c.logger.Error("Отсутствует text в send_message")
		return
	}

	// Получаем from_user_id, если не передан - используем текущего пользователя
	var fromUserID int64
	if fromUserIDStr, ok := dataMap["from_user_id"].(string); ok {
		// Конвертируем fromUserID в int64
		parsedUserID, err := strconv.ParseInt(fromUserIDStr, 10, 64)
		if err != nil {
			c.logger.Error("Неверный формат from_user_id", zap.Error(err))
			return
		}
		fromUserID = parsedUserID
		
		// Проверяем, что отправитель совпадает с текущим пользователем
		if fromUserID != c.userID {
			c.logger.Warn("Попытка отправить сообщение от имени другого пользователя", 
				zap.Int64("from_user_id", fromUserID), 
				zap.Int64("current_user_id", c.userID))
			return
		}
	} else {
		// Если from_user_id не передан, используем текущего пользователя
		fromUserID = c.userID
		c.logger.Debug("from_user_id не передан, используем текущего пользователя", 
			zap.Int64("user_id", c.userID))
	}

	// Используем MessageManager для отправки сообщения
	if c.server.messageManager != nil {
		c.logger.Info("Отправляем сообщение через MessageManager", 
			zap.Int64("chat_id", chatID),
			zap.Int64("from_user_id", fromUserID),
			zap.String("text", text))
		
		// Вызываем метод SendMessage напрямую через интерфейс
		message, err := c.server.messageManager.SendMessage(chatID, fromUserID, text, "text", nil)
		
		if err != nil {
			c.logger.Error("Ошибка отправки сообщения", zap.Error(err))
		} else if message != nil {
			c.logger.Info("Сообщение успешно отправлено", 
				zap.Int64("message_id", message.ID),
				zap.Int64("chat_id", chatID),
				zap.Int64("from_user_id", fromUserID),
				zap.String("text", text))
		}
	} else {
		c.logger.Error("MessageManager не установлен")
	}
}

// handleCallbackQuery обрабатывает callback query от inline кнопок
func (c *Client) handleCallbackQuery(data interface{}) {
	// Преобразуем data в map для извлечения параметров
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		c.logger.Error("Неверный формат данных для callback_query")
		return
	}

	buttonData, ok := dataMap["button"].(map[string]interface{})
	if !ok {
		c.logger.Error("Отсутствует button в callback_query")
		return
	}

	// Получаем chat_id из данных
	chatIDFloat, ok := dataMap["chat_id"].(float64)
	if !ok {
		c.logger.Error("Отсутствует chat_id в callback_query")
		return
	}
	chatID := int64(chatIDFloat)

	// Генерируем уникальный ID для callback query как строку
	callbackQueryID := fmt.Sprintf("cq_%d", time.Now().UnixNano())

	c.logger.Info("Получен callback query",
		zap.Int64("user_id", c.userID),
		zap.Int64("chat_id", chatID),
		zap.String("callback_query_id", callbackQueryID),
		zap.Any("button", buttonData))

	// Создаем CallbackQuery для BotManager
	callbackData, ok := buttonData["callback_data"].(string)
	if !ok {
		c.logger.Error("Не удалось получить callback_data из button", zap.Any("button", buttonData))
		return
	}
	
	if callbackData == "" {
		c.logger.Error("callback_data пустой", zap.Any("button", buttonData))
		return
	}
	
	// Создаем сообщение для callback query
	// Используем переданный chat_id
	message := &models.Message{
		ID:     time.Now().UnixNano(), // Генерируем int64 ID
		ChatID: chatID, // Используем переданный chat_id
		From: models.User{
			ID: c.userID,
		},
		Text:      "Inline keyboard message",
		Type:      "text",
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	callbackQuery := &models.CallbackQuery{
		ID:   callbackQueryID,
		Data: callbackData,
		From: models.User{
			ID: c.userID,
		},
		Message:      message,
		ChatInstance: "chat_instance",
	}

	// Добавляем callback query в BotManager для обработки ботом
	if c.server.botManager != nil {
		// Используем reflection для вызова метода AddCallbackQuery
		botManagerValue := reflect.ValueOf(c.server.botManager)
		addCallbackQueryMethod := botManagerValue.MethodByName("AddCallbackQuery")
		
		if addCallbackQueryMethod.IsValid() {
			// Находим токен бота - пробуем несколько способов
			var botToken string
			
			// Способ 1: Пытаемся извлечь из callback_data (если это формат cat:message_id)
			if strings.HasPrefix(callbackData, "cat:") {
				// Здесь можно попробовать найти сообщение по ID и определить бота
				// Пока используем дефолтный токен
				botToken = "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
			} else {
				// Способ 2: Используем дефолтный токен для всех callback query
				botToken = "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
			}
			
			// Способ 3: Получаем всех ботов и используем первого активного
			getAllBotsMethod := botManagerValue.MethodByName("GetAllBots")
			if getAllBotsMethod.IsValid() {
				results := getAllBotsMethod.Call(nil)
				if len(results) > 0 && !results[0].IsNil() {
					bots := results[0].Interface().([]models.Bot)
					for _, bot := range bots {
						if bot.IsActive {
							botToken = bot.Token
							break
						}
					}
				}
			}
			
			args := []reflect.Value{
				reflect.ValueOf(botToken),
				reflect.ValueOf(callbackQuery),
			}
			
			results := addCallbackQueryMethod.Call(args)
			
			if len(results) > 0 && !results[0].IsNil() {
				err := results[0].Interface().(error)
				c.logger.Error("Ошибка добавления callback query в BotManager", zap.Error(err))
			} else {
							c.logger.Info("Callback query добавлен в BotManager", 
				zap.String("callback_query_id", callbackQueryID),
				zap.String("callback_data", callbackData),
					zap.String("bot_token", botToken))
			}
		}
	}

	// Отправляем callback query всем участникам чата
	c.server.Broadcast("callback_query", map[string]interface{}{
		"id":       callbackQueryID,
		"user_id":  c.userID,
		"button":   buttonData,
		"data":     buttonData["callback_data"],
		"message": map[string]interface{}{
			"message_id": "msg_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			"chat": map[string]interface{}{
				"id": chatID, // Используем переданный chat_id
			},
		},
	})
}

// GetConnectedUsers возвращает список подключенных пользователей
func (s *Server) GetConnectedUsers() []int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	users := make([]int64, 0, len(s.clients))
	for client := range s.clients {
		users = append(users, client.userID)
	}
	
	return users
}

// IsUserConnected проверяет, подключен ли пользователь
func (s *Server) IsUserConnected(userID int64) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for client := range s.clients {
		if client.userID == userID {
			return true
		}
	}
	
	return false
}
