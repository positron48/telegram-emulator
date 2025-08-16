package emulator

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/pkg/logger"
	"telegram-emulator/internal/repository"
	"telegram-emulator/internal/websocket"

	"go.uber.org/zap"
)

// MessageManager управляет сообщениями в эмуляторе
type MessageManager struct {
	messageRepo *repository.MessageRepository
	chatRepo    *repository.ChatRepository
	userRepo    *repository.UserRepository
	wsServer    *websocket.Server
	logger      *zap.Logger
}

// NewMessageManager создает новый экземпляр MessageManager
func NewMessageManager(messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository, userRepo *repository.UserRepository, wsServer *websocket.Server) *MessageManager {
	return &MessageManager{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		wsServer:    wsServer,
		logger:      logger.GetLogger(),
	}
}

// SendMessage отправляет сообщение в чат
func (m *MessageManager) SendMessage(chatID, fromUserID, text, messageType string) (*models.Message, error) {
	// Генерируем уникальный ID
	id, err := m.generateID()
	if err != nil {
		return nil, err
	}

	// Получаем информацию о пользователе
	fromUser, err := m.userRepo.GetByID(fromUserID)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя", zap.String("user_id", fromUserID), zap.Error(err))
		return nil, err
	}

	// Проверяем, является ли пользователь участником чата
	chat, err := m.chatRepo.GetByID(chatID)
	if err != nil {
		m.logger.Error("Ошибка получения чата", zap.String("chat_id", chatID), zap.Error(err))
		return nil, err
	}

	// Проверяем, является ли пользователь участником
	isMember := false
	for _, member := range chat.Members {
		if member.ID == fromUserID {
			isMember = true
			break
		}
	}

	// Если пользователь не участник, добавляем его (кроме приватных чатов)
	if !isMember && chat.Type != "private" {
		if err := m.chatRepo.AddMember(chatID, fromUserID); err != nil {
			m.logger.Error("Ошибка добавления пользователя в чат", zap.String("chat_id", chatID), zap.String("user_id", fromUserID), zap.Error(err))
			// Не прерываем отправку сообщения, просто логируем ошибку
		} else {
			m.logger.Info("Пользователь автоматически добавлен в чат", zap.String("chat_id", chatID), zap.String("user_id", fromUserID))
		}
	}

	// Создаем сообщение
	message := &models.Message{
		ID:        id,
		ChatID:    chatID,
		FromID:    fromUserID,
		From:      *fromUser,
		Text:      text,
		Type:      messageType,
		Status:    models.MessageStatusSending,
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	// Сохраняем сообщение
	if err := m.messageRepo.Create(message); err != nil {
		m.logger.Error("Ошибка создания сообщения", zap.Error(err))
		return nil, err
	}

	// Обновляем счетчик непрочитанных сообщений
	// Получаем количество непрочитанных сообщений
	unreadCount, err := m.messageRepo.GetUnreadCount(chatID)
	if err != nil {
		m.logger.Error("Ошибка подсчета непрочитанных сообщений", zap.String("chat_id", chatID), zap.Error(err))
	} else {
		if err := m.chatRepo.UpdateUnreadCount(chatID, int(unreadCount)); err != nil {
			m.logger.Error("Ошибка обновления счетчика непрочитанных", zap.String("chat_id", chatID), zap.Error(err))
		}
	}

	// Отправляем WebSocket уведомление
	m.broadcastMessage(message)

	// Эмулируем доставку сообщения
	go m.simulateMessageDelivery(message)

	m.logger.Info("Сообщение отправлено", 
		zap.String("message_id", message.ID),
		zap.String("chat_id", chatID),
		zap.String("from_user", fromUser.Username))

	return message, nil
}

// GetChatMessages получает сообщения чата
func (m *MessageManager) GetChatMessages(chatID string, limit, offset int) ([]models.Message, error) {
	messages, err := m.messageRepo.GetByChatID(chatID, limit, offset)
	if err != nil {
		m.logger.Error("Ошибка получения сообщений чата", zap.String("chat_id", chatID), zap.Error(err))
		return nil, err
	}
	return messages, nil
}

// GetMessage получает сообщение по ID
func (m *MessageManager) GetMessage(id string) (*models.Message, error) {
	message, err := m.messageRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения сообщения", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return message, nil
}

// UpdateMessageStatus обновляет статус сообщения
func (m *MessageManager) UpdateMessageStatus(id, status string) error {
	if err := m.messageRepo.UpdateStatus(id, status); err != nil {
		m.logger.Error("Ошибка обновления статуса сообщения", zap.String("id", id), zap.String("status", status), zap.Error(err))
		return err
	}

	// Отправляем WebSocket уведомление об обновлении статуса
	m.broadcastMessageStatusUpdate(id, status)

	m.logger.Info("Статус сообщения обновлен", zap.String("id", id), zap.String("status", status))
	return nil
}

// MarkChatAsRead помечает все сообщения чата как прочитанные
func (m *MessageManager) MarkChatAsRead(chatID, userID string) error {
	// Получаем все непрочитанные сообщения в чате
	messages, err := m.messageRepo.GetByChatID(chatID, 1000, 0)
	if err != nil {
		return err
	}

	// Помечаем сообщения как прочитанные
	for _, message := range messages {
		if message.FromID != userID && message.Status != models.MessageStatusRead {
			if err := m.messageRepo.UpdateStatus(message.ID, models.MessageStatusRead); err != nil {
				m.logger.Error("Ошибка пометки сообщения как прочитанного", zap.String("message_id", message.ID), zap.Error(err))
			}
		}
	}

	// Обновляем счетчик непрочитанных
	unreadCount, err := m.messageRepo.GetUnreadCount(chatID)
	if err != nil {
		m.logger.Error("Ошибка подсчета непрочитанных сообщений", zap.String("chat_id", chatID), zap.Error(err))
		return err
	}
	
	if err := m.chatRepo.UpdateUnreadCount(chatID, int(unreadCount)); err != nil {
		m.logger.Error("Ошибка обновления счетчика непрочитанных", zap.String("chat_id", chatID), zap.Error(err))
		return err
	}

	// Отправляем WebSocket уведомление
	m.broadcastChatRead(chatID, userID)

	m.logger.Info("Чат помечен как прочитанный", zap.String("chat_id", chatID), zap.String("user_id", userID))
	return nil
}

// DeleteMessage удаляет сообщение
func (m *MessageManager) DeleteMessage(id string) error {
	if err := m.messageRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления сообщения", zap.String("id", id), zap.Error(err))
		return err
	}

	// Отправляем WebSocket уведомление об удалении
	m.broadcastMessageDelete(id)

	m.logger.Info("Сообщение удалено", zap.String("id", id))
	return nil
}

// SearchMessages ищет сообщения по тексту
func (m *MessageManager) SearchMessages(chatID, query string) ([]models.Message, error) {
	messages, err := m.messageRepo.SearchByText(chatID, query)
	if err != nil {
		m.logger.Error("Ошибка поиска сообщений", zap.String("chat_id", chatID), zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return messages, nil
}

// simulateMessageDelivery эмулирует доставку сообщения
func (m *MessageManager) simulateMessageDelivery(message *models.Message) {
	// Эмулируем задержку сети
	time.Sleep(100 * time.Millisecond)

	// Обновляем статус на "отправлено"
	if err := m.UpdateMessageStatus(message.ID, models.MessageStatusSent); err != nil {
		m.logger.Error("Ошибка обновления статуса на 'отправлено'", zap.String("message_id", message.ID), zap.Error(err))
		return
	}

	// Эмулируем задержку доставки
	time.Sleep(200 * time.Millisecond)

	// Обновляем статус на "доставлено"
	if err := m.UpdateMessageStatus(message.ID, models.MessageStatusDelivered); err != nil {
		m.logger.Error("Ошибка обновления статуса на 'доставлено'", zap.String("message_id", message.ID), zap.Error(err))
	}
}

// broadcastMessage отправляет WebSocket уведомление о новом сообщении
func (m *MessageManager) broadcastMessage(message *models.Message) {
	if m.wsServer != nil {
		// Отправляем всем участникам чата
		chat, err := m.chatRepo.GetByID(message.ChatID)
		if err != nil {
			m.logger.Error("Ошибка получения чата для broadcast", zap.String("chat_id", message.ChatID), zap.Error(err))
			return
		}

		// Отправляем всем участникам чата, включая отправителя
		for _, member := range chat.Members {
			m.wsServer.BroadcastToUser(member.ID, "message", map[string]interface{}{
				"id":        message.ID,
				"chat_id":   message.ChatID,
				"from":      message.From,
				"text":      message.Text,
				"type":      message.Type,
				"timestamp": message.Timestamp,
				"status":    message.Status,
			})
		}
	}
}

// broadcastMessageStatusUpdate отправляет WebSocket уведомление об обновлении статуса
func (m *MessageManager) broadcastMessageStatusUpdate(messageID, status string) {
	if m.wsServer != nil {
		m.wsServer.Broadcast("message_status_update", map[string]interface{}{
			"message_id": messageID,
			"status":     status,
		})
	}
}

// broadcastMessageDelete отправляет WebSocket уведомление об удалении сообщения
func (m *MessageManager) broadcastMessageDelete(messageID string) {
	if m.wsServer != nil {
		m.wsServer.Broadcast("message_delete", map[string]interface{}{
			"message_id": messageID,
		})
	}
}

// broadcastChatRead отправляет WebSocket уведомление о прочтении чата
func (m *MessageManager) broadcastChatRead(chatID, userID string) {
	if m.wsServer != nil {
		m.wsServer.Broadcast("chat_read", map[string]interface{}{
			"chat_id": chatID,
			"user_id": userID,
		})
	}
}

// generateID генерирует уникальный ID
func (m *MessageManager) generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
