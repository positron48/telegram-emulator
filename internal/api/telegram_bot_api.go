package api

import (
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
	// Telegram Bot API маршруты
	botAPI := router.Group("/bot:token")
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
	}
}

// GetMe возвращает информацию о боте
func (api *TelegramBotAPI) GetMe(c *gin.Context) {
	token := c.Param("token")
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
	token := c.Param("token")
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

	// Получаем параметры запроса
	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsed
		}
	}

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	timeout := 0
	if timeoutStr := c.Query("timeout"); timeoutStr != "" {
		if parsed, err := strconv.Atoi(timeoutStr); err == nil && parsed >= 0 && parsed <= 50 {
			timeout = parsed
		}
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
	var telegramUpdates []gin.H
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
	token := c.Param("token")
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
		ChatID                string `json:"chat_id" binding:"required"`
		Text                  string `json:"text" binding:"required"`
		ParseMode             string `json:"parse_mode"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
		DisableNotification   bool   `json:"disable_notification"`
		ProtectContent        bool   `json:"protect_content"`
		ReplyToMessageID      int64  `json:"reply_to_message_id"`
		AllowSendingWithoutReply bool `json:"allow_sending_without_reply"`
		ReplyMarkup           interface{} `json:"reply_markup"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error_code": 400, "description": "Bad Request: " + err.Error()})
		return
	}

	// Получаем пользователя-бота
	botUser, err := api.userManager.GetUserByUsername(bot.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Bot user not found: " + err.Error()})
		return
	}

	// Конвертируем Telegram chat_id в внутренний chat_id
	internalChatID := request.ChatID
	if telegramChatID, err := strconv.ParseInt(request.ChatID, 10, 64); err == nil {
		// Это Telegram chat_id, нужно найти внутренний chat_id
		chats, err := api.chatManager.GetAllChats()
		if err != nil {
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

	// Отправляем сообщение через обычный API с клавиатурой
	message, err := api.messageManager.SendMessage(internalChatID, botUser.ID, request.Text, "text", request.ReplyMarkup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error_code": 500, "description": "Failed to send message: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": message.ToTelegramMessage(),
	})
}

// SetWebhook устанавливает webhook для бота
func (api *TelegramBotAPI) SetWebhook(c *gin.Context) {
	token := c.Param("token")
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
	token := c.Param("token")
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
	token := c.Param("token")
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

// findBotByToken находит бота по токену
func (api *TelegramBotAPI) findBotByToken(token string) (*models.Bot, error) {
	// Получаем всех ботов и ищем по токену
	bots, err := api.botManager.GetAllBots()
	if err != nil {
		return nil, err
	}

	for _, bot := range bots {
		if bot.Token == token {
			return &bot, nil
		}
	}

	return nil, &models.BotNotFoundError{}
}

// AnswerCallbackQuery отвечает на callback query
func (api *TelegramBotAPI) AnswerCallbackQuery(c *gin.Context) {
	token := c.Param("token")
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
		zap.Bool("show_alert", request.ShowAlert))

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": true,
	})
}

// EditMessageText редактирует текст сообщения
func (api *TelegramBotAPI) EditMessageText(c *gin.Context) {
	token := c.Param("token")
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
