package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TelegramBotAPI представляет API совместимый с Telegram Bot API
type TelegramBotAPI struct {
	botManager     *emulator.BotManager
	userManager    *emulator.UserManager
	chatManager    *emulator.ChatManager
	messageManager *emulator.MessageManager
	logger         *zap.Logger
}

// NewTelegramBotAPI создает новый экземпляр TelegramBotAPI
func NewTelegramBotAPI(botManager *emulator.BotManager, userManager *emulator.UserManager, chatManager *emulator.ChatManager, messageManager *emulator.MessageManager) *TelegramBotAPI {
	return &TelegramBotAPI{
		botManager:     botManager,
		userManager:    userManager,
		chatManager:    chatManager,
		messageManager: messageManager,
		logger:         botManager.GetLogger(),
	}
}

// SetupTelegramBotRoutes настраивает маршруты Telegram Bot API
func (api *TelegramBotAPI) SetupTelegramBotRoutes(router *gin.Engine) {
	// Telegram Bot API маршруты с правильными заголовками
	botAPI := router.Group("/bot:token")
	botAPI.Use(func(c *gin.Context) {
		// Устанавливаем правильные заголовки для Telegram Bot API
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Server", "Telegram-Emulator/1.0")
		c.Next()
	})
	
	{
		// Основные методы - поддерживаем и GET и POST для совместимости
		botAPI.GET("/getMe", api.GetMe)
		botAPI.POST("/getMe", api.GetMe)
		botAPI.GET("/getUpdates", api.GetUpdates)
		botAPI.POST("/getUpdates", api.GetUpdates)
		botAPI.POST("/sendMessage", api.SendMessage)
		botAPI.POST("/setWebhook", api.SetWebhook)
		botAPI.GET("/deleteWebhook", api.DeleteWebhook)
		botAPI.POST("/deleteWebhook", api.DeleteWebhook)
		botAPI.GET("/getWebhookInfo", api.GetWebhookInfo)
		botAPI.POST("/getWebhookInfo", api.GetWebhookInfo)
		
		// Callback query методы
		botAPI.POST("/answerCallbackQuery", api.AnswerCallbackQuery)
		botAPI.POST("/editMessageText", api.EditMessageText)
		botAPI.POST("/editMessageReplyMarkup", api.EditMessageReplyMarkup)
	}
}

// GetMe возвращает информацию о боте
func (api *TelegramBotAPI) GetMe(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Нормализуем токен
	token = api.normalizeToken(token)

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	// Конвертируем строковый ID в числовой (как в Telegram Bot API)
	botID := api.convertStringIDToInt64(bot.ID)

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"result": gin.H{
			"id":         botID,
			"is_bot":     true,
			"first_name": bot.Name,
			"username":   bot.Username,
			"can_join_groups": true,
			"can_read_all_group_messages": false,
			"supports_inline_queries": false,
		},
	})
}

// GetUpdates возвращает обновления для бота
func (api *TelegramBotAPI) GetUpdates(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Нормализуем токен
	token = api.normalizeToken(token)

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	// Получаем параметры запроса из query, form (POST) или JSON
	var params struct {
		Offset  int `form:"offset" json:"offset"`
		Limit   int `form:"limit" json:"limit"`
		Timeout int `form:"timeout" json:"timeout"`
	}
	// Значения по умолчанию
	params.Limit = 100
	params.Timeout = 0

	// Унифицированный биндинг: поддерживает GET query, POST form, POST JSON
	_ = c.ShouldBind(&params)

	// Валидация и границы
	offset := params.Offset
	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	timeout := params.Timeout
	if timeout < 0 {
		timeout = 0
	}
	if timeout > 50 {
		timeout = 50
	}

	// Получаем обновления
	updates, err := api.botManager.GetBotUpdates(bot.ID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Internal Server Error"})
		return
	}

	// Если нет обновлений и указан timeout, ждем новые обновления
	if len(updates) == 0 && timeout > 0 {
		api.logger.Info("Long polling: ожидаем новые обновления", 
			zap.String("bot_id", bot.ID),
			zap.Int("timeout", timeout))
		
		// Ждем новые обновления в течение timeout секунд
		startTime := time.Now()
		for time.Since(startTime) < time.Duration(timeout)*time.Second {
			// Проверяем новые обновления каждые 1 секунду (увеличено с 100ms для снижения нагрузки на БД)
			time.Sleep(1 * time.Second)
			
			updates, err = api.botManager.GetBotUpdates(bot.ID, offset, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Internal Server Error"})
				return
			}
			
			if len(updates) > 0 {
				api.logger.Info("Long polling: получены новые обновления", 
					zap.String("bot_id", bot.ID),
					zap.Int("count", len(updates)),
					zap.Duration("wait_time", time.Since(startTime)))
				break
			}
		}
	}

	// Конвертируем в формат Telegram Bot API
	telegramUpdates := make([]gin.H, 0)
	for _, update := range updates {
		telegramUpdate := gin.H{
			"update_id": update.UpdateID,
		}

		if update.Message != nil {
			telegramUpdate["message"] = update.Message.ToTelegramMessage()
		}
		if update.EditedMessage != nil {
			telegramUpdate["edited_message"] = update.EditedMessage.ToTelegramMessage()
		}
		if update.CallbackQuery != nil {
			telegramUpdate["callback_query"] = update.CallbackQuery
		}

		telegramUpdates = append(telegramUpdates, telegramUpdate)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": telegramUpdates,
	})
}

// SendMessage отправляет сообщение
func (api *TelegramBotAPI) SendMessage(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Нормализуем токен
	token = api.normalizeToken(token)

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	var request struct {
		ChatID                string      `json:"chat_id" binding:"required"`
		Text                  string      `json:"text" binding:"required"`
		ParseMode             string      `json:"parse_mode"`
		DisableWebPagePreview bool        `json:"disable_web_page_preview"`
		DisableNotification   bool        `json:"disable_notification"`
		ProtectContent        bool        `json:"protect_content"`
		ReplyToMessageID      int64       `json:"reply_to_message_id"`
		AllowSendingWithoutReply bool     `json:"allow_sending_without_reply"`
		ReplyMarkup           interface{} `json:"reply_markup"`
	}

	// Улучшенная обработка JSON - поддерживаем как JSON, так и form data
	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&request); err != nil {
			rawData, _ := c.GetRawData()
			api.logger.Error("Ошибка парсинга JSON", zap.Error(err), zap.String("body", string(rawData)))
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: invalid JSON format"})
			return
		}
	} else {
		// Поддержка form data для совместимости
		if err := c.ShouldBind(&request); err != nil {
			api.logger.Error("Ошибка парсинга form data", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: invalid form data"})
			return
		}
	}

	// Валидация обязательных полей
	if request.ChatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: chat_id is required"})
		return
	}
	if request.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: text is required"})
		return
	}

	// Получаем пользователя-бота
	botUser, err := api.userManager.GetUserByUsername(bot.Username)
	if err != nil {
		api.logger.Error("Ошибка получения пользователя-бота", zap.String("bot_username", bot.Username), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Bot user not found"})
		return
	}

	// Конвертируем Telegram chat_id в внутренний chat_id
	internalChatID := request.ChatID
	if telegramChatID, err := strconv.ParseInt(request.ChatID, 10, 64); err == nil {
		// Это Telegram chat_id, нужно найти внутренний chat_id
		chats, err := api.chatManager.GetAllChats()
		if err != nil {
			api.logger.Error("Ошибка получения чатов", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Failed to get chats"})
			return
		}
		
		for _, chat := range chats {
			// Конвертируем внутренний chat_id в Telegram chat_id
			chatTelegramID := int64(0)
			if len(chat.ID) > 0 {
				for i, char := range chat.ID {
					if i < 8 { // Ограничиваем длину
						chatTelegramID = chatTelegramID*31 + int64(char)
					}
				}
			}
			
			if chatTelegramID == telegramChatID {
				internalChatID = chat.ID
				break
			}
		}
	}

	// Валидация reply_markup если он передан
	if request.ReplyMarkup != nil {
		if err := api.validateReplyMarkup(request.ReplyMarkup); err != nil {
			api.logger.Error("Ошибка валидации reply_markup", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: invalid reply_markup format"})
			return
		}
	}

	// Отправляем сообщение через обычный API с клавиатурой
	message, err := api.messageManager.SendMessage(internalChatID, botUser.ID, request.Text, "text", request.ReplyMarkup)
	if err != nil {
		api.logger.Error("Ошибка отправки сообщения", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Failed to send message"})
		return
	}

	// Конвертируем в формат Telegram Bot API
	telegramMessage := message.ToTelegramMessage()

	api.logger.Info("Сообщение успешно отправлено", 
		zap.String("bot_id", bot.ID),
		zap.String("chat_id", request.ChatID),
		zap.String("message_id", message.ID),
		zap.Bool("is_command", message.IsCommand()))

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": telegramMessage,
	})
}

// SetWebhook устанавливает webhook для бота
func (api *TelegramBotAPI) SetWebhook(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	var request struct {
		URL                string   `json:"url" binding:"required"`
		Certificate        string   `json:"certificate"`
		IPAddress          string   `json:"ip_address"`
		MaxConnections     int      `json:"max_connections"`
		AllowedUpdates     []string `json:"allowed_updates"`
		DropPendingUpdates bool     `json:"drop_pending_updates"`
		SecretToken        string   `json:"secret_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: " + err.Error()})
		return
	}

	// Обновляем webhook URL
	bot.WebhookURL = request.URL
	if err := api.botManager.UpdateBot(bot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": true,
	})
}

// DeleteWebhook удаляет webhook
func (api *TelegramBotAPI) DeleteWebhook(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	// Удаляем webhook URL
	bot.WebhookURL = ""
	if err := api.botManager.UpdateBot(bot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": true,
	})
}

// GetWebhookInfo возвращает информацию о webhook
func (api *TelegramBotAPI) GetWebhookInfo(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	webhookInfo := gin.H{
		"url": bot.WebhookURL,
	}

	if bot.WebhookURL == "" {
		webhookInfo["url"] = ""
		webhookInfo["has_custom_certificate"] = false
		webhookInfo["pending_update_count"] = 0
		webhookInfo["max_connections"] = 40
		webhookInfo["last_error_date"] = nil
		webhookInfo["last_error_message"] = nil
		webhookInfo["last_synchronization_error_date"] = nil
		webhookInfo["allowed_updates"] = []string{}
	} else {
		webhookInfo["has_custom_certificate"] = false
		webhookInfo["pending_update_count"] = 0
		webhookInfo["max_connections"] = 40
		webhookInfo["last_error_date"] = nil
		webhookInfo["last_error_message"] = nil
		webhookInfo["last_synchronization_error_date"] = nil
		webhookInfo["allowed_updates"] = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": webhookInfo,
	})
}

// normalizeToken убирает двоеточие из начала токена
func (api *TelegramBotAPI) normalizeToken(token string) string {
	if len(token) > 0 && token[0] == ':' {
		return token[1:]
	}
	return token
}

// findBotByToken находит бота по токену
func (api *TelegramBotAPI) findBotByToken(token string) (*models.Bot, error) {
	// Получаем всех ботов и ищем по токену
	bots, err := api.botManager.GetAllBots()
	if err != nil {
		api.logger.Error("Ошибка получения ботов", zap.Error(err))
		return nil, err
	}

	api.logger.Info("Поиск бота по токену", 
		zap.String("token", token),
		zap.Int("total_bots", len(bots)))

	for _, bot := range bots {
		api.logger.Debug("Проверяем бота", 
			zap.String("bot_id", bot.ID),
			zap.String("bot_token", bot.Token),
			zap.String("search_token", token))
		
		if bot.Token == token {
			api.logger.Info("Бот найден", zap.String("bot_id", bot.ID))
			return &bot, nil
		}
	}

	api.logger.Warn("Бот не найден", zap.String("token", token))
	return nil, &models.BotNotFoundError{}
}

// AnswerCallbackQuery отвечает на callback query
func (api *TelegramBotAPI) AnswerCallbackQuery(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	var request struct {
		CallbackQueryID string `json:"callback_query_id" binding:"required"`
		Text            string `json:"text,omitempty"`
		ShowAlert       bool   `json:"show_alert,omitempty"`
		URL             string `json:"url,omitempty"`
		CacheTime       int    `json:"cache_time,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: " + err.Error()})
		return
	}

	// Логируем ответ на callback query
	api.logger.Info("Ответ на callback query",
		zap.String("bot_id", bot.ID),
		zap.String("callback_query_id", request.CallbackQueryID),
		zap.String("text", request.Text),
		zap.Bool("show_alert", request.ShowAlert),
		zap.String("url", request.URL))

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": true,
	})
}

// EditMessageReplyMarkup редактирует клавиатуру сообщения
func (api *TelegramBotAPI) EditMessageReplyMarkup(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	var request struct {
		ChatID      string      `json:"chat_id" binding:"required"`
		MessageID   string      `json:"message_id" binding:"required"`
		ReplyMarkup interface{} `json:"reply_markup,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: " + err.Error()})
		return
	}

	// Валидация reply_markup если он передан
	if request.ReplyMarkup != nil {
		if err := api.validateReplyMarkup(request.ReplyMarkup); err != nil {
			api.logger.Error("Ошибка валидации reply_markup", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: invalid reply_markup format"})
			return
		}
	}

	// Логируем редактирование клавиатуры
	api.logger.Info("Редактирование клавиатуры сообщения",
		zap.String("bot_id", bot.ID),
		zap.String("chat_id", request.ChatID),
		zap.String("message_id", request.MessageID))

	// Конвертируем строковый ID в числовой
	botID := api.convertStringIDToInt64(bot.ID)

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": gin.H{
			"message_id": request.MessageID,
			"from": gin.H{
				"id":         botID,
				"is_bot":     true,
				"first_name": bot.Name,
				"username":   bot.Username,
			},
			"chat": gin.H{
				"id":    request.ChatID,
				"type":  "private",
				"title": bot.Name,
			},
			"date": time.Now().Unix(),
			"reply_markup": request.ReplyMarkup,
		},
	})
}

// EditMessageText редактирует текст сообщения
func (api *TelegramBotAPI) EditMessageText(c *gin.Context) {
	token := api.normalizeToken(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: token is empty"})
		return
	}

	// Находим бота по токену
	bot, err := api.findBotByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error_code": 401, "description": "Unauthorized"})
		return
	}

	var request struct {
		ChatID      string      `json:"chat_id" binding:"required"`
		MessageID   string      `json:"message_id" binding:"required"`
		Text        string      `json:"text" binding:"required"`
		ReplyMarkup interface{} `json:"reply_markup,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: " + err.Error()})
		return
	}

	// Логируем редактирование сообщения
	api.logger.Info("Редактирование сообщения",
		zap.String("chat_id", request.ChatID),
		zap.String("message_id", request.MessageID),
		zap.String("text", request.Text))

	// Конвертируем строковый ID в числовой
	botID := api.convertStringIDToInt64(bot.ID)

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": gin.H{
			"message_id": request.MessageID,
			"from": gin.H{
				"id":         botID,
				"is_bot":     true,
				"first_name": bot.Name,
				"username":   bot.Username,
			},
			"chat": gin.H{
				"id":    request.ChatID,
				"type":  "private",
				"title": bot.Name,
			},
			"date": time.Now().Unix(),
			"text": request.Text,
		},
	})
}

// validateReplyMarkup валидирует формат reply_markup
func (api *TelegramBotAPI) validateReplyMarkup(replyMarkup interface{}) error {
	// Базовая валидация - проверяем, что это map
	markupMap, ok := replyMarkup.(map[string]interface{})
	if !ok {
		return fmt.Errorf("reply_markup must be an object")
	}

	// Проверяем наличие одного из типов клавиатур
	hasInlineKeyboard := false
	hasKeyboard := false
	hasRemoveKeyboard := false
	hasForceReply := false

	if _, exists := markupMap["inline_keyboard"]; exists {
		hasInlineKeyboard = true
	}
	if _, exists := markupMap["keyboard"]; exists {
		hasKeyboard = true
	}
	if _, exists := markupMap["remove_keyboard"]; exists {
		hasRemoveKeyboard = true
	}
	if _, exists := markupMap["force_reply"]; exists {
		hasForceReply = true
	}

	// Должен быть только один тип клавиатуры
	keyboardTypes := 0
	if hasInlineKeyboard { keyboardTypes++ }
	if hasKeyboard { keyboardTypes++ }
	if hasRemoveKeyboard { keyboardTypes++ }
	if hasForceReply { keyboardTypes++ }

	if keyboardTypes == 0 {
		return fmt.Errorf("reply_markup must contain one of: inline_keyboard, keyboard, remove_keyboard, force_reply")
	}
	if keyboardTypes > 1 {
		return fmt.Errorf("reply_markup can contain only one keyboard type")
	}

	// Валидация inline_keyboard
	if hasInlineKeyboard {
		if err := api.validateInlineKeyboard(markupMap["inline_keyboard"]); err != nil {
			return fmt.Errorf("invalid inline_keyboard: %w", err)
		}
	}

	// Валидация keyboard
	if hasKeyboard {
		if err := api.validateKeyboard(markupMap["keyboard"]); err != nil {
			return fmt.Errorf("invalid keyboard: %w", err)
		}
	}

	return nil
}

// validateInlineKeyboard валидирует inline клавиатуру
func (api *TelegramBotAPI) validateInlineKeyboard(keyboard interface{}) error {
	keyboardArray, ok := keyboard.([]interface{})
	if !ok {
		return fmt.Errorf("inline_keyboard must be an array")
	}

	for i, row := range keyboardArray {
		rowArray, ok := row.([]interface{})
		if !ok {
			return fmt.Errorf("inline_keyboard row %d must be an array", i)
		}

		for j, button := range rowArray {
			buttonMap, ok := button.(map[string]interface{})
			if !ok {
				return fmt.Errorf("inline_keyboard button %d in row %d must be an object", j, i)
			}

			// Проверяем обязательные поля
			text, exists := buttonMap["text"].(string)
			if !exists || text == "" {
				return fmt.Errorf("inline_keyboard button %d in row %d must have non-empty text", j, i)
			}

			// Проверяем наличие хотя бы одного callback_data или url
			hasCallbackData := false
			hasURL := false
			if _, exists := buttonMap["callback_data"]; exists {
				hasCallbackData = true
			}
			if _, exists := buttonMap["url"]; exists {
				hasURL = true
			}

			if !hasCallbackData && !hasURL {
				return fmt.Errorf("inline_keyboard button %d in row %d must have either callback_data or url", j, i)
			}
		}
	}

	return nil
}

// validateKeyboard валидирует обычную клавиатуру
func (api *TelegramBotAPI) validateKeyboard(keyboard interface{}) error {
	keyboardArray, ok := keyboard.([]interface{})
	if !ok {
		return fmt.Errorf("keyboard must be an array")
	}

	for i, row := range keyboardArray {
		rowArray, ok := row.([]interface{})
		if !ok {
			return fmt.Errorf("keyboard row %d must be an array", i)
		}

		for j, button := range rowArray {
			buttonMap, ok := button.(map[string]interface{})
			if !ok {
				return fmt.Errorf("keyboard button %d in row %d must be an object", j, i)
			}

			// Проверяем обязательные поля
			text, exists := buttonMap["text"].(string)
			if !exists || text == "" {
				return fmt.Errorf("keyboard button %d in row %d must have non-empty text", j, i)
			}
		}
	}

	return nil
}

// convertStringIDToInt64 конвертирует строковый ID в числовой (как в Telegram Bot API)
func (api *TelegramBotAPI) convertStringIDToInt64(id string) int64 {
	result := int64(0)
	if len(id) > 0 {
		for i, char := range id {
			if i < 8 { // Ограничиваем длину для предотвращения переполнения
				result = result*31 + int64(char)
			}
		}
	}
	return result
}
