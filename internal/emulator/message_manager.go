package emulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
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
	botManager  *BotManager
	wsServer    *websocket.Server
	logger      *zap.Logger
}

// NewMessageManager создает новый экземпляр MessageManager
func NewMessageManager(messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository, userRepo *repository.UserRepository, botManager *BotManager, wsServer *websocket.Server) *MessageManager {
	return &MessageManager{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		botManager:  botManager,
		wsServer:    wsServer,
		logger:      logger.GetLogger(),
	}
}

// SendMessage отправляет сообщение в чат
func (m *MessageManager) SendMessage(chatID int64, fromUserID int64, text, messageType string, replyMarkup interface{}) (*models.Message, error) {
	// Генерируем уникальный ID
	id, err := m.generateID()
	if err != nil {
		return nil, err
	}

	// Получаем информацию о пользователе
	fromUser, err := m.userRepo.GetByID(fromUserID)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя", zap.Int64("user_id", fromUserID), zap.Error(err))
		return nil, err
	}

	// Проверяем, является ли пользователь участником чата
	chat, err := m.chatRepo.GetByID(chatID)
	if err != nil {
		m.logger.Warn("Чат не найден, создаем чат в базе данных", zap.Int64("chat_id", chatID))
		// Создаем чат в базе данных
		newChat := &models.Chat{
			ID:        chatID,
			Type:      "private",
			Title:     "Budget Chat",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := m.chatRepo.Create(newChat); err != nil {
			m.logger.Error("Ошибка создания чата в базе данных", zap.Error(err))
			// Если не удалось создать, используем виртуальный чат
			chat = &models.Chat{
				ID:    chatID,
				Type:  "private",
				Title: "Virtual Chat",
			}
		} else {
			// Добавляем бота в чат
			if err := m.chatRepo.AddMember(chatID, fromUserID); err != nil {
				m.logger.Error("Ошибка добавления бота в чат", zap.Error(err))
			}
			chat = newChat
			m.logger.Info("Чат успешно создан в базе данных", zap.Int64("chat_id", chatID))
		}
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
			m.logger.Error("Ошибка добавления пользователя в чат", zap.Int64("chat_id", chatID), zap.Int64("user_id", fromUserID), zap.Error(err))
			// Не прерываем отправку сообщения, просто логируем ошибку
		} else {
			m.logger.Info("Пользователь автоматически добавлен в чат", zap.Int64("chat_id", chatID), zap.Int64("user_id", fromUserID))
		}
	}

	// Создаем сообщение
	message := &models.Message{
		ID:         id,
		ChatID:     chatID,
		FromID:     fromUserID,
		From:       *fromUser,
		Text:       text,
		Type:       messageType,
		Status:     models.MessageStatusSending,
		IsOutgoing: false,
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}

	// Устанавливаем клавиатуру, если она есть
	if replyMarkup != nil {
		if err := message.SetReplyMarkup(replyMarkup); err != nil {
			m.logger.Error("Ошибка установки клавиатуры", zap.Error(err))
			// Не прерываем отправку сообщения, просто логируем ошибку
		}
	}

	// Парсим и устанавливаем сущности (команды, упоминания, хештеги, URL)
	if err := message.ParseAndSetEntities(); err != nil {
		m.logger.Error("Ошибка парсинга сущностей", zap.Error(err))
		// Не прерываем отправку сообщения, просто логируем ошибку
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
		m.logger.Error("Ошибка подсчета непрочитанных сообщений", zap.Int64("chat_id", chatID), zap.Error(err))
	} else {
		if err := m.chatRepo.UpdateUnreadCount(chatID, int(unreadCount)); err != nil {
			m.logger.Error("Ошибка обновления счетчика непрочитанных", zap.Int64("chat_id", chatID), zap.Error(err))
		}
	}

	// Отправляем WebSocket уведомление
	m.broadcastMessage(message)

	// Эмулируем доставку сообщения
	go m.simulateMessageDelivery(message)

	// Уведомляем ботов о новом сообщении
	m.notifyBots(message)

	m.logger.Info("Сообщение отправлено",
		zap.Int64("message_id", message.ID),
		zap.Int64("chat_id", chatID),
		zap.String("from_user", fromUser.Username))

	return message, nil
}

// GetChatMessages получает сообщения чата
func (m *MessageManager) GetChatMessages(chatID int64, limit, offset int) ([]models.Message, error) {
	messages, err := m.messageRepo.GetByChatID(chatID, limit, offset)
	if err != nil {
		m.logger.Error("Ошибка получения сообщений чата", zap.Int64("chat_id", chatID), zap.Error(err))
		return nil, err
	}

	// Клавиатуры автоматически десериализуются при обращении к GetReplyMarkup()

	return messages, nil
}

// GetMessage получает сообщение по ID
func (m *MessageManager) GetMessage(id int64) (*models.Message, error) {
	message, err := m.messageRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения сообщения", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// Клавиатура автоматически десериализуется при обращении к GetReplyMarkup()

	return message, nil
}

// UpdateMessageStatus обновляет статус сообщения
func (m *MessageManager) UpdateMessageStatus(id int64, status string) error {
	if err := m.messageRepo.UpdateStatus(id, status); err != nil {
		m.logger.Error("Ошибка обновления статуса сообщения", zap.Int64("id", id), zap.String("status", status), zap.Error(err))
		return err
	}

	// Отправляем WebSocket уведомление об обновлении статуса
	m.broadcastMessageStatusUpdate(id, status)

	m.logger.Info("Статус сообщения обновлен", zap.Int64("id", id), zap.String("status", status))
	return nil
}

// MarkChatAsRead помечает все сообщения чата как прочитанные
func (m *MessageManager) MarkChatAsRead(chatID int64, userID int64) error {
	// Получаем все непрочитанные сообщения в чате
	messages, err := m.messageRepo.GetByChatID(chatID, 1000, 0)
	if err != nil {
		return err
	}

	// Помечаем сообщения как прочитанные
	for _, message := range messages {
		if message.FromID != userID && message.Status != models.MessageStatusRead {
			if err := m.messageRepo.UpdateStatus(message.ID, models.MessageStatusRead); err != nil {
				m.logger.Error("Ошибка пометки сообщения как прочитанного", zap.Int64("message_id", message.ID), zap.Error(err))
			}
		}
	}

	// Обновляем счетчик непрочитанных
	unreadCount, err := m.messageRepo.GetUnreadCount(chatID)
	if err != nil {
		m.logger.Error("Ошибка подсчета непрочитанных сообщений", zap.Int64("chat_id", chatID), zap.Error(err))
		return err
	}

	if err := m.chatRepo.UpdateUnreadCount(chatID, int(unreadCount)); err != nil {
		m.logger.Error("Ошибка обновления счетчика непрочитанных", zap.Int64("chat_id", chatID), zap.Error(err))
		return err
	}

	// Отправляем WebSocket уведомление
	m.broadcastChatRead(chatID, userID)

	m.logger.Info("Чат помечен как прочитанный", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
	return nil
}

// DeleteMessage удаляет сообщение
func (m *MessageManager) DeleteMessage(id int64) error {
	if err := m.messageRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления сообщения", zap.Int64("id", id), zap.Error(err))
		return err
	}

	// Отправляем WebSocket уведомление об удалении
	m.broadcastMessageDelete(id)

	m.logger.Info("Сообщение удалено", zap.Int64("id", id))
	return nil
}

// SearchMessages ищет сообщения по тексту
func (m *MessageManager) SearchMessages(chatID int64, query string) ([]models.Message, error) {
	messages, err := m.messageRepo.SearchByText(chatID, query)
	if err != nil {
		m.logger.Error("Ошибка поиска сообщений", zap.Int64("chat_id", chatID), zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return messages, nil
}

// HandleCallbackQuery обрабатывает callback query от inline кнопки
func (m *MessageManager) HandleCallbackQuery(userID int64, messageID int64, callbackData string) (*models.CallbackQuery, error) {
	// Получаем сообщение
	message, err := m.messageRepo.GetByID(messageID)
	if err != nil {
		m.logger.Error("Ошибка получения сообщения для callback query", zap.Int64("message_id", messageID), zap.Error(err))
		return nil, err
	}

	// Получаем пользователя
	user, err := m.userRepo.GetByID(userID)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя для callback query", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Генерируем уникальный ID для callback query как строку
	callbackID, err := m.generateID()
	if err != nil {
		m.logger.Error("Ошибка генерации ID для callback query", zap.Error(err))
		return nil, err
	}

	// Создаем callback query
	callbackQuery := &models.CallbackQuery{
		ID:      fmt.Sprintf("cq_%d", callbackID),
		From:    *user,
		Message: message,
		Data:    callbackData,
	}

	// Уведомляем ботов о callback query
	m.notifyBotsCallbackQuery(callbackQuery)

	m.logger.Info("Callback query обработан",
		zap.String("callback_id", callbackQuery.ID),
		zap.Int64("message_id", messageID),
		zap.Int64("user_id", userID),
		zap.String("callback_data", callbackData))

	return callbackQuery, nil
}

// notifyBotsCallbackQuery уведомляет ботов о callback query
func (m *MessageManager) notifyBotsCallbackQuery(callbackQuery *models.CallbackQuery) {
	if m.botManager == nil {
		m.logger.Error("botManager равен nil - уведомления ботов отключены")
		return
	}

	// Получаем всех активных ботов
	bots, err := m.botManager.GetAllBots()
	if err != nil {
		m.logger.Error("Ошибка получения ботов для уведомления callback query", zap.Error(err))
		return
	}

	// Создаем обновление для каждого бота
	for _, bot := range bots {
		if !bot.IsActive {
			continue
		}

		// Проверяем, является ли сообщение от этого бота
		botUser, err := m.userRepo.GetByUsername(bot.Username)
		if err != nil {
			m.logger.Error("Ошибка получения пользователя-бота",
				zap.Int64("bot_id", bot.ID),
				zap.String("bot_username", bot.Username),
				zap.Error(err))
			continue
		}

		// Если сообщение от этого бота, уведомляем его о callback query
		if callbackQuery.Message.FromID == botUser.ID {
			// Создаем обновление
			update := &models.Update{
				CallbackQuery: callbackQuery,
			}

			// Добавляем в очередь обновлений бота
			if err := m.botManager.AddUpdate(bot.ID, update); err != nil {
				m.logger.Error("Ошибка добавления callback query для бота",
					zap.Int64("bot_id", bot.ID),
					zap.Error(err))
			}

			// Если у бота есть webhook URL, отправляем обновление в webhook
			if bot.WebhookURL != "" {
				go m.sendWebhookUpdate(&bot, update)
			}
		}
	}

	m.logger.Info("Боты уведомлены о callback query",
		zap.String("callback_id", callbackQuery.ID),
		zap.Int("bots_count", len(bots)))
}

// simulateMessageDelivery эмулирует доставку сообщения
func (m *MessageManager) simulateMessageDelivery(message *models.Message) {
	// Эмулируем задержку сети
	time.Sleep(100 * time.Millisecond)

	// Обновляем статус на "отправлено"
	if err := m.UpdateMessageStatus(message.ID, models.MessageStatusSent); err != nil {
		m.logger.Error("Ошибка обновления статуса на 'отправлено'", zap.Int64("message_id", message.ID), zap.Error(err))
		return
	}

	// Эмулируем задержку доставки
	time.Sleep(200 * time.Millisecond)

	// Обновляем статус на "доставлено"
	if err := m.UpdateMessageStatus(message.ID, models.MessageStatusDelivered); err != nil {
		m.logger.Error("Ошибка обновления статуса на 'доставлено'", zap.Int64("message_id", message.ID), zap.Error(err))
	}
}

// broadcastMessage отправляет WebSocket уведомление о новом сообщении
func (m *MessageManager) broadcastMessage(message *models.Message) {
	if m.wsServer != nil {
		m.logger.Info("Начинаем broadcast сообщения",
			zap.Int64("message_id", message.ID),
			zap.Int64("chat_id", message.ChatID),
			zap.String("text", message.Text))

		// Отправляем всем участникам чата
		chat, err := m.chatRepo.GetByID(message.ChatID)
		if err != nil {
			m.logger.Error("Ошибка получения чата для broadcast", zap.Int64("chat_id", message.ChatID), zap.Error(err))
			return
		}

		m.logger.Info("Найдены участники чата",
			zap.Int64("chat_id", message.ChatID),
			zap.Int("members_count", len(chat.Members)))

		// Отправляем всем участникам чата, включая отправителя
		for _, member := range chat.Members {
			m.logger.Info("Отправляем сообщение участнику",
				zap.Int64("message_id", message.ID),
				zap.Int64("member_id", member.ID),
				zap.String("member_username", member.Username))

			messageData := map[string]interface{}{
				"id":        message.ID,
				"chat_id":   message.ChatID,
				"from":      message.From,
				"text":      message.Text,
				"type":      message.Type,
				"timestamp": message.Timestamp,
				"status":    message.Status,
			}

			// Добавляем клавиатуру, если она есть
			if replyMarkup := message.GetReplyMarkup(); replyMarkup != nil {
				messageData["reply_markup"] = replyMarkup
			}

			m.wsServer.BroadcastToUser(member.ID, "message", messageData)
		}
	} else {
		m.logger.Error("wsServer равен nil - broadcast отключен")
	}
}

// broadcastMessageStatusUpdate отправляет WebSocket уведомление об обновлении статуса
func (m *MessageManager) broadcastMessageStatusUpdate(messageID int64, status string) {
	if m.wsServer != nil {
		m.wsServer.Broadcast("message_status_update", map[string]interface{}{
			"message_id": messageID,
			"status":     status,
		})
	}
}

// broadcastMessageDelete отправляет WebSocket уведомление об удалении сообщения
func (m *MessageManager) broadcastMessageDelete(messageID int64) {
	if m.wsServer != nil {
		m.wsServer.Broadcast("message_delete", map[string]interface{}{
			"message_id": messageID,
		})
	}
}

// broadcastChatRead отправляет WebSocket уведомление о прочтении чата
func (m *MessageManager) broadcastChatRead(chatID int64, userID int64) {
	if m.wsServer != nil {
		m.wsServer.Broadcast("chat_read", map[string]interface{}{
			"chat_id": chatID,
			"user_id": userID,
		})
	}
}

// generateID генерирует уникальный ID
func (m *MessageManager) generateID() (int64, error) {
	// Используем Unix timestamp в миллисекундах + случайное число для уникальности
	timestamp := time.Now().UnixMilli()
	random := rand.Int63n(1000) // случайное число от 0 до 999
	return timestamp*1000 + random, nil
}

// notifyBots уведомляет всех активных ботов о новом сообщении
func (m *MessageManager) notifyBots(message *models.Message) {
	if m.botManager == nil {
		m.logger.Error("botManager равен nil - уведомления ботов отключены")
		return
	}

	// Получаем всех активных ботов
	bots, err := m.botManager.GetAllBots()
	if err != nil {
		m.logger.Error("Ошибка получения ботов для уведомления", zap.Error(err))
		return
	}

	// Создаем обновление для каждого бота только если бот является участником чата
	for _, bot := range bots {
		if !bot.IsActive {
			continue
		}

		// Проверяем, не является ли отправитель сообщения этим ботом
		// Если да, то не уведомляем бота о его собственном сообщении
		botUser, err := m.userRepo.GetByUsername(bot.Username)
		if err != nil {
			m.logger.Error("Ошибка получения пользователя-бота",
				zap.Int64("bot_id", bot.ID),
				zap.String("bot_username", bot.Username),
				zap.Error(err))
			continue
		}

		// Если сообщение от этого бота, пропускаем
		if message.FromID == botUser.ID {
			m.logger.Debug("Пропускаем уведомление бота о его собственном сообщении",
				zap.Int64("bot_id", bot.ID),
				zap.Int64("message_id", message.ID))
			continue
		}

		// Проверяем, состоит ли бот в этом чате
		chat, err := m.chatRepo.GetByID(message.ChatID)
		if err != nil {
			m.logger.Warn("Чат не найден для проверки членства бота, считаем что бот является участником", zap.Int64("chat_id", message.ChatID))
			// Если чат не найден, считаем что бот является участником (для виртуальных чатов)
			chat = &models.Chat{
				ID:    message.ChatID,
				Type:  "private",
				Title: "Virtual Chat",
			}
		}
		isBotMember := false
		// Для виртуальных чатов считаем что бот является участником
		if chat.Title == "Virtual Chat" {
			isBotMember = true
			m.logger.Debug("Виртуальный чат, бот считается участником",
				zap.Int64("bot_id", bot.ID),
				zap.Int64("chat_id", message.ChatID))
		} else {
			for _, member := range chat.Members {
				if member.Username == bot.Username || member.ID == botUser.ID {
					isBotMember = true
					m.logger.Debug("Бот найден в участниках чата",
						zap.Int64("bot_id", bot.ID),
						zap.Int64("chat_id", message.ChatID),
						zap.String("member_username", member.Username))
					break
				}
			}
		}
		if !isBotMember {
			// Бот не состоит в этом чате — пропускаем уведомление
			m.logger.Debug("Бот не является участником чата, уведомление пропущено",
				zap.Int64("bot_id", bot.ID),
				zap.Int64("chat_id", message.ChatID))
			continue
		}

		// Создаем обновление
		update := &models.Update{
			Message: message,
		}

		// Добавляем в очередь обновлений бота
		if err := m.botManager.AddUpdate(bot.ID, update); err != nil {
			m.logger.Error("Ошибка добавления обновления для бота",
				zap.Int64("bot_id", bot.ID),
				zap.Error(err))
		}

		// Если у бота есть webhook URL, отправляем обновление в webhook
		if bot.WebhookURL != "" {
			go m.sendWebhookUpdate(&bot, update)
		}
	}

	m.logger.Info("Боты уведомлены о новом сообщении",
		zap.Int64("message_id", message.ID),
		zap.Int("bots_count", len(bots)))
}

// sendWebhookUpdate отправляет обновление в webhook бота
func (m *MessageManager) sendWebhookUpdate(bot *models.Bot, update *models.Update) {
	// Конвертируем обновление в формат Telegram Bot API
	webhookUpdate := map[string]interface{}{
		"update_id": update.UpdateID,
	}

	if update.Message != nil {
		webhookUpdate["message"] = update.Message.ToTelegramMessage()
	}
	if update.EditedMessage != nil {
		webhookUpdate["edited_message"] = update.EditedMessage.ToTelegramMessage()
	}
	if update.CallbackQuery != nil {
		webhookUpdate["callback_query"] = update.CallbackQuery.ToTelegramCallbackQuery()
	}

	// Отправляем POST запрос в webhook
	jsonData, err := json.Marshal(webhookUpdate)
	if err != nil {
		m.logger.Error("Ошибка сериализации webhook обновления",
			zap.Int64("bot_id", bot.ID),
			zap.Error(err))
		return
	}

	resp, err := http.Post(bot.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		m.logger.Error("Ошибка отправки webhook",
			zap.Int64("bot_id", bot.ID),
			zap.String("webhook_url", bot.WebhookURL),
			zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		m.logger.Error("Webhook вернул ошибку",
			zap.Int64("bot_id", bot.ID),
			zap.String("webhook_url", bot.WebhookURL),
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)))
		return
	}

	m.logger.Info("Webhook обновление отправлено",
		zap.Int64("bot_id", bot.ID),
		zap.String("webhook_url", bot.WebhookURL),
		zap.Int64("update_id", update.UpdateID))
}
