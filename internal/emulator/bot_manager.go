package emulator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/pkg/logger"
	"telegram-emulator/internal/repository"

	"go.uber.org/zap"
)

// BotManager управляет ботами в эмуляторе
type BotManager struct {
	botRepo      *repository.BotRepository
	userRepo     *repository.UserRepository
	messageRepo  *repository.MessageRepository
	chatRepo     *repository.ChatRepository
	logger       *zap.Logger
	updateQueue  map[string][]models.Update // Очередь обновлений для каждого бота
	updateID     int64                      // Глобальный счетчик update_id
	chatIDMap    map[int64]string           // Маппинг Telegram chat_id -> внутренний chat_id
}

// NewBotManager создает новый экземпляр BotManager
func NewBotManager(botRepo *repository.BotRepository, userRepo *repository.UserRepository, messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository) *BotManager {
	return &BotManager{
		botRepo:     botRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		logger:      logger.GetLogger(),
		updateQueue: make(map[string][]models.Update),
		updateID:    1,
		chatIDMap:   make(map[int64]string),
	}
}

// CreateBot создает нового бота
func (m *BotManager) CreateBot(name, username, token, webhookURL string) (*models.Bot, error) {
	// Генерируем уникальный ID
	id, err := m.generateID()
	if err != nil {
		return nil, err
	}

	bot := &models.Bot{
		ID:         id,
		Name:       name,
		Username:   username,
		Token:      token,
		WebhookURL: webhookURL,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := m.botRepo.Create(bot); err != nil {
		m.logger.Error("Ошибка создания бота", zap.Error(err))
		return nil, err
	}

	// Создаем пользователя-бота
	botUser := &models.User{
		ID:        fmt.Sprintf("bot_%s", id),
		Username:  username,
		FirstName: name,
		IsBot:     true,
		IsOnline:  true,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if err := m.userRepo.Create(botUser); err != nil {
		m.logger.Error("Ошибка создания пользователя-бота", zap.Error(err))
		// Удаляем бота если не удалось создать пользователя
		m.botRepo.Delete(id)
		return nil, err
	}

	m.logger.Info("Создан новый бот", 
		zap.String("id", bot.ID),
		zap.String("name", bot.Name),
		zap.String("username", bot.Username))

	return bot, nil
}

// GetBot получает бота по ID
func (m *BotManager) GetBot(id string) (*models.Bot, error) {
	bot, err := m.botRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения бота", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return bot, nil
}

// GetAllBots получает всех ботов
func (m *BotManager) GetAllBots() ([]models.Bot, error) {
	bots, err := m.botRepo.GetAll()
	if err != nil {
		m.logger.Error("Ошибка получения всех ботов", zap.Error(err))
		return nil, err
	}
	return bots, nil
}

// UpdateBot обновляет бота
func (m *BotManager) UpdateBot(bot *models.Bot) error {
	bot.UpdatedAt = time.Now()
	if err := m.botRepo.Update(bot); err != nil {
		m.logger.Error("Ошибка обновления бота", zap.String("id", bot.ID), zap.Error(err))
		return err
	}

	m.logger.Info("Бот обновлен", zap.String("id", bot.ID))
	return nil
}

// DeleteBot удаляет бота
func (m *BotManager) DeleteBot(id string) error {
	if err := m.botRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления бота", zap.String("id", id), zap.Error(err))
		return err
	}

	// Удаляем пользователя-бота
	botUserID := fmt.Sprintf("bot_%s", id)
	if err := m.userRepo.Delete(botUserID); err != nil {
		m.logger.Error("Ошибка удаления пользователя-бота", zap.String("id", botUserID), zap.Error(err))
	}

	m.logger.Info("Бот удален", zap.String("id", id))
	return nil
}

// SendBotMessage отправляет сообщение через бота
func (m *BotManager) SendBotMessage(botID, chatID, text, parseMode string) (*models.Message, error) {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		return nil, err
	}

	if !bot.IsActive {
		return nil, fmt.Errorf("бот неактивен")
	}

	// Проверяем, является ли chatID Telegram chat_id (число)
	// Если да, то конвертируем в внутренний chat_id
	internalChatID := chatID
	if telegramChatID, err := strconv.ParseInt(chatID, 10, 64); err == nil {
		// Это Telegram chat_id, ищем внутренний chat_id
		if internalID, exists := m.chatIDMap[telegramChatID]; exists {
			internalChatID = internalID
			m.logger.Info("Найден маппинг chat_id", 
				zap.Int64("telegram_chat_id", telegramChatID),
				zap.String("internal_chat_id", internalID))
		} else {
			m.logger.Error("Маппинг chat_id не найден", 
				zap.Int64("telegram_chat_id", telegramChatID))
			return nil, fmt.Errorf("чат с Telegram ID %d не найден", telegramChatID)
		}
	}

	// Получаем пользователя-бота
	botUserID := fmt.Sprintf("bot_%s", botID)
	botUser, err := m.userRepo.GetByID(botUserID)
	if err != nil {
		return nil, fmt.Errorf("пользователь-бот не найден")
	}

	// Создаем сообщение
	messageID, err := m.generateID()
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		ID:        messageID,
		ChatID:    internalChatID,
		FromID:    botUser.ID,
		From:      *botUser,
		Text:      text,
		Type:      models.MessageTypeText,
		Status:    models.MessageStatusSent,
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := m.messageRepo.Create(message); err != nil {
		m.logger.Error("Ошибка создания сообщения бота", zap.Error(err))
		return nil, err
	}

	m.logger.Info("Сообщение бота отправлено", 
		zap.String("bot_id", botID),
		zap.String("chat_id", internalChatID),
		zap.String("message_id", messageID))

	return message, nil
}

// GetBotUpdates возвращает обновления для бота
func (m *BotManager) GetBotUpdates(botID string, offset, limit int) ([]models.Update, error) {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		return nil, err
	}

	if !bot.IsActive {
		return nil, fmt.Errorf("бот неактивен")
	}

	// Получаем обновления из очереди
	queue, exists := m.updateQueue[botID]
	if !exists {
		return []models.Update{}, nil
	}

	// Фильтруем по offset
	var filteredUpdates []models.Update
	for _, update := range queue {
		if update.UpdateID > int64(offset) {
			filteredUpdates = append(filteredUpdates, update)
		}
	}

	// Ограничиваем по limit
	if len(filteredUpdates) > limit {
		filteredUpdates = filteredUpdates[:limit]
	}

	m.logger.Info("Получены обновления для бота", 
		zap.String("bot_id", botID),
		zap.Int("count", len(filteredUpdates)),
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	return filteredUpdates, nil
}

// ProcessWebhook обрабатывает webhook от бота
func (m *BotManager) ProcessWebhook(botID string, update *models.Update) error {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		return err
	}

	if !bot.IsActive {
		return fmt.Errorf("бот неактивен")
	}

	// Обрабатываем сообщение из webhook
	if update.Message != nil {
		// Сохраняем сообщение
		if err := m.messageRepo.Create(update.Message); err != nil {
			m.logger.Error("Ошибка сохранения сообщения из webhook", zap.Error(err))
			return err
		}

		m.logger.Info("Webhook сообщение обработано", 
			zap.String("bot_id", botID),
			zap.String("message_id", update.Message.ID))
	}

	return nil
}

// AddUpdate добавляет обновление в очередь для бота
func (m *BotManager) AddUpdate(botID string, update *models.Update) error {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		return err
	}

	if !bot.IsActive {
		return fmt.Errorf("бот неактивен")
	}

	// Устанавливаем update_id
	update.UpdateID = m.updateID
	m.updateID++

	// Устанавливаем timestamp
	update.Timestamp = time.Now()

	// Сохраняем маппинг chat_id если есть сообщение
	if update.Message != nil {
		// Конвертируем строковый chat_id в int64 для Telegram API
		telegramChatID := int64(0)
		if len(update.Message.ChatID) > 0 {
			for i, char := range update.Message.ChatID {
				if i < 8 { // Ограничиваем длину
					telegramChatID = telegramChatID*31 + int64(char)
				}
			}
		}
		m.chatIDMap[telegramChatID] = update.Message.ChatID
	}

	// Добавляем в очередь
	if m.updateQueue[botID] == nil {
		m.updateQueue[botID] = []models.Update{}
	}
	m.updateQueue[botID] = append(m.updateQueue[botID], *update)

	// Ограничиваем размер очереди (максимум 1000 обновлений)
	if len(m.updateQueue[botID]) > 1000 {
		m.updateQueue[botID] = m.updateQueue[botID][len(m.updateQueue[botID])-1000:]
	}

	m.logger.Info("Обновление добавлено в очередь", 
		zap.String("bot_id", botID),
		zap.Int64("update_id", update.UpdateID))

	return nil
}

// ClearUpdates очищает очередь обновлений для бота
func (m *BotManager) ClearUpdates(botID string) error {
	delete(m.updateQueue, botID)
	m.logger.Info("Очередь обновлений очищена", zap.String("bot_id", botID))
	return nil
}

// generateID генерирует уникальный ID
func (m *BotManager) generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
