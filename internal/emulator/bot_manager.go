package emulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
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
	updateQueue  map[int64][]models.Update // Очередь обновлений для каждого бота
	updateID     int64                      // Глобальный счетчик update_id
	chatIDMap    map[int64]string           // Маппинг Telegram chat_id -> внутренний chat_id
	nextUpdateID map[int64]int64           // Персональный счетчик update_id для каждого бота
}

// NewBotManager создает новый экземпляр BotManager
func NewBotManager(botRepo *repository.BotRepository, userRepo *repository.UserRepository, messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository) *BotManager {
	// Инициализируем генератор случайных чисел (удалено rand.Seed - deprecated)
	
	return &BotManager{
		botRepo:     botRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		logger:      logger.GetLogger(),
		updateQueue: make(map[int64][]models.Update),
		updateID:    1,
		chatIDMap:   make(map[int64]string),
		nextUpdateID: make(map[int64]int64),
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
		ID:        id, // Используем тот же ID что и у бота
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
		if deleteErr := m.botRepo.Delete(id); deleteErr != nil {
			m.logger.Error("Ошибка удаления бота после неудачного создания пользователя", zap.Error(deleteErr))
		}
		return nil, err
	}

	m.logger.Info("Создан новый бот", 
		zap.Int64("id", bot.ID),
		zap.String("name", bot.Name),
		zap.String("username", bot.Username))

	return bot, nil
}

// GetBot получает бота по ID
func (m *BotManager) GetBot(id int64) (*models.Bot, error) {
	bot, err := m.botRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения бота", zap.Int64("id", id), zap.Error(err))
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
		m.logger.Error("Ошибка обновления бота", zap.Int64("id", bot.ID), zap.Error(err))
		return err
	}

	m.logger.Info("Бот обновлен", zap.Int64("id", bot.ID))
	return nil
}

// DeleteBot удаляет бота
func (m *BotManager) DeleteBot(id int64) error {
	if err := m.botRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления бота", zap.Int64("id", id), zap.Error(err))
		return err
	}

	// Удаляем пользователя-бота (используем тот же ID)
	if err := m.userRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления пользователя-бота", zap.Int64("id", id), zap.Error(err))
	}

	m.logger.Info("Бот удален", zap.Int64("id", id))
	return nil
}

// SendBotMessage отправляет сообщение через бота
func (m *BotManager) SendBotMessage(botID int64, chatID int64, text, parseMode string) (*models.Message, error) {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		return nil, err
	}

	if !bot.IsActive {
		return nil, fmt.Errorf("бот неактивен")
	}



	// Получаем пользователя-бота
	botUser, err := m.userRepo.GetByUsername(bot.Username)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя-бота по username", 
			zap.String("username", bot.Username), zap.Error(err))
		return nil, fmt.Errorf("пользователь-бот не найден")
	}

	// Создаем сообщение
	messageID, err := m.generateID()
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		ID:        messageID,
		ChatID:    chatID, // Используем chatID напрямую, так как он уже int64
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
		zap.Int64("bot_id", botID),
		zap.Int64("chat_id", chatID),
		zap.Int64("message_id", messageID))

	return message, nil
}

// GetBotUpdates возвращает обновления для бота
func (m *BotManager) GetBotUpdates(botID int64, offset, limit int) ([]models.Update, error) {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		return nil, err
	}

	if !bot.IsActive {
		return nil, fmt.Errorf("бот неактивен")
	}

	// Если offset не указан, используем сохраненный offset бота
	if offset == 0 {
		offset = int(bot.LastUpdateOffset)
	}

	// Получаем обновления из очереди
	queue, exists := m.updateQueue[botID]
	if !exists {
		m.logger.Debug("Очередь обновлений пуста для бота", zap.Int64("bot_id", botID))
		return []models.Update{}, nil
	}

	m.logger.Debug("Получена очередь обновлений", 
		zap.Int64("bot_id", botID),
		zap.Int("queue_size", len(queue)),
		zap.Int("offset", offset))

	// Проверяем, есть ли обновления с update_id больше offset
	maxUpdateID := int64(0)
	for _, update := range queue {
		if update.UpdateID > maxUpdateID {
			maxUpdateID = update.UpdateID
		}
	}

	// Если offset больше максимального update_id, возвращаем пустой список
	if int64(offset) > maxUpdateID && maxUpdateID > 0 {
		m.logger.Info("Offset больше максимального update_id, возвращаем пустой список", 
			zap.Int64("bot_id", botID),
			zap.Int("offset", offset),
			zap.Int64("max_update_id", maxUpdateID))
		return []models.Update{}, nil
	}

	// Фильтруем по offset - возвращаем обновления с update_id >= offset
	var filteredUpdates []models.Update
	for _, update := range queue {
		if update.UpdateID >= int64(offset) {
			filteredUpdates = append(filteredUpdates, update)
			// Логируем информацию о callback query
			if update.CallbackQuery != nil {
				m.logger.Debug("Найден callback query в обновлении", 
					zap.Int64("bot_id", botID),
					zap.Int64("update_id", update.UpdateID),
					zap.String("callback_data", update.CallbackQuery.Data))
			}
		}
	}

	// Ограничиваем по limit
	if len(filteredUpdates) > limit {
		filteredUpdates = filteredUpdates[:limit]
	}

	// НЕ обновляем offset бота автоматически - бот должен сам управлять своим offset
	// Это стандартное поведение Telegram Bot API

	m.logger.Info("Получены обновления для бота", 
		zap.Int64("bot_id", botID),
		zap.Int("count", len(filteredUpdates)),
		zap.Int("offset", offset),
		zap.Int("limit", limit),
		zap.Int64("last_update_offset", bot.LastUpdateOffset))

	return filteredUpdates, nil
}

// ProcessWebhook обрабатывает webhook от бота
func (m *BotManager) ProcessWebhook(botID int64, update *models.Update) error {
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
			zap.Int64("bot_id", botID),
			zap.Int64("message_id", update.Message.ID))
	}

	return nil
}

// AddUpdate добавляет обновление в очередь для бота
func (m *BotManager) AddUpdate(botID int64, update *models.Update) error {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		m.logger.Error("Ошибка получения бота в AddUpdate", 
			zap.Int64("bot_id", botID), 
			zap.Error(err))
		return err
	}

	if !bot.IsActive {
		m.logger.Error("Бот неактивен в AddUpdate", zap.Int64("bot_id", botID))
		return fmt.Errorf("бот неактивен")
	}

	// Устанавливаем update_id персонифицировано для бота
	nextID := m.nextUpdateID[botID]
	if nextID == 0 {
		nextID = 1
	}
	update.UpdateID = nextID
	m.nextUpdateID[botID] = nextID + 1

	// Устанавливаем timestamp
	update.Timestamp = time.Now()

	// Сохраняем маппинг chat_id если есть сообщение
	if update.Message != nil {
		// ChatID уже int64, сохраняем маппинг
		m.chatIDMap[update.Message.ChatID] = strconv.FormatInt(update.Message.ChatID, 10)
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
		zap.Int64("bot_id", botID),
		zap.Int64("update_id", update.UpdateID))

	// Обрабатываем команды автоматически
	if update.Message != nil && update.Message.IsCommand() {
		m.logger.Info("Найдена команда, запускаем обработку", 
			zap.Int64("bot_id", botID),
			zap.String("command", update.Message.GetCommand()),
			zap.Int64("message_id", update.Message.ID))
		go m.handleCommand(botID, update.Message)
	}

	return nil
}

// AddCallbackQuery добавляет callback query в очередь обновлений для бота
func (m *BotManager) AddCallbackQuery(botToken string, callbackQuery *models.CallbackQuery) error {
	// Находим бота по токену
	bots, err := m.GetAllBots()
	if err != nil {
		return err
	}
	
	var bot *models.Bot
	for _, b := range bots {
		if b.Token == botToken {
			bot = &b
			break
		}
	}
	
	if bot == nil {
		return fmt.Errorf("бот с токеном %s не найден", botToken)
	}
	
	if !bot.IsActive {
		return fmt.Errorf("бот неактивен")
	}
	
	// Создаем обновление с callback query с персонифицированным update_id
	nextID := m.nextUpdateID[bot.ID]
	if nextID == 0 {
		nextID = 1
	}
	update := &models.Update{
		UpdateID:      nextID,
		CallbackQuery: callbackQuery,
		Timestamp:     time.Now(),
	}
	m.nextUpdateID[bot.ID] = nextID + 1
	
	// Добавляем в очередь
	if m.updateQueue[bot.ID] == nil {
		m.updateQueue[bot.ID] = []models.Update{}
	}
	
	m.updateQueue[bot.ID] = append(m.updateQueue[bot.ID], *update)
	
	// Ограничиваем размер очереди
	if len(m.updateQueue[bot.ID]) > 1000 {
		m.updateQueue[bot.ID] = m.updateQueue[bot.ID][len(m.updateQueue[bot.ID])-1000:]
	}
	
	// Если у бота есть webhook URL, отправляем обновление в webhook
	if bot.WebhookURL != "" {
		go m.sendWebhookUpdate(bot, update)
	}
	
	m.logger.Info("Callback query добавлен в очередь", 
		zap.Int64("bot_id", bot.ID),
		zap.String("bot_token", botToken),
		zap.Int64("update_id", update.UpdateID),
		zap.String("callback_data", callbackQuery.Data))
	
	return nil
}

// ClearUpdates очищает очередь обновлений для бота
func (m *BotManager) ClearUpdates(botID int64) error {
	delete(m.updateQueue, botID)
	m.logger.Info("Очередь обновлений очищена", zap.Int64("bot_id", botID))
	return nil
}

// GetLogger возвращает логгер
func (m *BotManager) GetLogger() *zap.Logger {
	return m.logger
}

// generateID генерирует уникальный ID
func (m *BotManager) generateID() (int64, error) {
	// Используем Unix timestamp в миллисекундах + случайное число для уникальности
	timestamp := time.Now().UnixMilli()
	random := rand.Int63n(1000) // случайное число от 0 до 999
	return timestamp*1000 + random, nil
}

// sendWebhookUpdate отправляет обновление через webhook
func (m *BotManager) sendWebhookUpdate(bot *models.Bot, update *models.Update) {
	m.logger.Info("Начинаем отправку webhook", 
		zap.Int64("bot_id", bot.ID),
		zap.String("webhook_url", bot.WebhookURL),
		zap.Int64("update_id", update.UpdateID))
	
	if bot.WebhookURL == "" {
		m.logger.Error("Webhook URL не установлен для бота", zap.Int64("bot_id", bot.ID))
		return
	}

	// Конвертируем обновление в формат Telegram Bot API
	webhookUpdate := map[string]interface{}{
		"update_id": update.UpdateID,
	}

	if update.Message != nil {
		webhookUpdate["message"] = update.Message.ToTelegramMessage()
		m.logger.Debug("Добавлено сообщение в webhook", zap.Int64("bot_id", bot.ID))
	}
	if update.EditedMessage != nil {
		webhookUpdate["edited_message"] = update.EditedMessage.ToTelegramMessage()
		m.logger.Debug("Добавлено редактированное сообщение в webhook", zap.Int64("bot_id", bot.ID))
	}
	if update.CallbackQuery != nil {
		webhookUpdate["callback_query"] = update.CallbackQuery.ToTelegramCallbackQuery()
		m.logger.Debug("Добавлен callback query в webhook", 
			zap.Int64("bot_id", bot.ID),
			zap.String("callback_data", update.CallbackQuery.Data))
	}

	jsonData, err := json.Marshal(webhookUpdate)
	if err != nil {
		m.logger.Error("Ошибка маршалинга обновления для webhook", zap.Error(err))
		return
	}

	m.logger.Debug("Отправляем webhook", 
		zap.Int64("bot_id", bot.ID),
		zap.String("webhook_url", bot.WebhookURL),
		zap.String("json_data", string(jsonData)))

	resp, err := http.Post(bot.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		m.logger.Error("Ошибка отправки обновления через webhook", 
			zap.Int64("bot_id", bot.ID),
			zap.String("webhook_url", bot.WebhookURL),
			zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		m.logger.Error("Webhook вернул не OK статус", 
			zap.Int64("bot_id", bot.ID),
			zap.String("webhook_url", bot.WebhookURL),
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", string(body)))
	} else {
		m.logger.Info("Обновление успешно отправлено через webhook", 
			zap.Int64("bot_id", bot.ID),
			zap.String("webhook_url", bot.WebhookURL),
			zap.Int64("update_id", update.UpdateID))
	}
}

// handleCommand обрабатывает команды бота
func (m *BotManager) handleCommand(botID int64, message *models.Message) {
	// Получаем бота
	bot, err := m.GetBot(botID)
	if err != nil {
		m.logger.Error("Ошибка получения бота для обработки команды", 
			zap.Int64("bot_id", botID), 
			zap.Error(err))
		return
	}

	// Получаем пользователя-бота
	botUser, err := m.userRepo.GetByUsername(bot.Username)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя-бота", 
			zap.Int64("bot_id", botID),
			zap.String("bot_username", bot.Username),
			zap.Error(err))
		return
	}

	// Получаем команду
	command := message.GetCommand()
	m.logger.Info("Обрабатываем команду", 
		zap.Int64("bot_id", botID),
		zap.String("command", command),
		zap.Int64("chat_id", message.ChatID))

	// Обрабатываем команду /start
	if command == "/start" {
		responseText := fmt.Sprintf("Привет! Я бот %s. Добро пожаловать!", bot.Name)
		
		// Отправляем ответное сообщение
		responseMessage := &models.Message{
			ID:        time.Now().UnixNano(),
			ChatID:    message.ChatID,
			FromID:    botUser.ID,
			Text:      responseText,
			Type:      "text",
			Status:    "sending",
			IsOutgoing: true,
			Timestamp: time.Now(),
			CreatedAt: time.Now(),
		}

		// Сохраняем сообщение в базе данных
		if err := m.messageRepo.Create(responseMessage); err != nil {
			m.logger.Error("Ошибка сохранения ответного сообщения", 
				zap.Int64("bot_id", botID),
				zap.Error(err))
			return
		}

		m.logger.Info("Ответ на команду /start отправлен", 
			zap.Int64("bot_id", botID),
			zap.Int64("message_id", responseMessage.ID),
			zap.Int64("chat_id", message.ChatID))
	}
}
