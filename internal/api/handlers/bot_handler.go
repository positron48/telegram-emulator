package handlers

import (
	"net/http"
	"strconv"

	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"

	"github.com/gin-gonic/gin"
)

// BotHandler обрабатывает запросы к API ботов
type BotHandler struct {
	botManager *emulator.BotManager
}

// NewBotHandler создает новый экземпляр BotHandler
func NewBotHandler(botManager *emulator.BotManager) *BotHandler {
	return &BotHandler{
		botManager: botManager,
	}
}

// GetAll возвращает всех ботов
func (h *BotHandler) GetAll(c *gin.Context) {
	bots, err := h.botManager.GetAllBots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bots": bots,
	})
}

// Create создает нового бота
func (h *BotHandler) Create(c *gin.Context) {
	var botData struct {
		Name       string `json:"name" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Token      string `json:"token" binding:"required"`
		WebhookURL string `json:"webhook_url"`
	}

	if err := c.ShouldBindJSON(&botData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bot, err := h.botManager.CreateBot(botData.Name, botData.Username, botData.Token, botData.WebhookURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"bot": bot,
	})
}

// GetByID возвращает бота по ID
func (h *BotHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID бота обязателен"})
		return
	}

	id, err := ParseBotID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID бота"})
		return
	}

	bot, err := h.botManager.GetBot(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бот не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bot": bot,
	})
}

// Update обновляет бота
func (h *BotHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID бота обязателен"})
		return
	}

	id, err := ParseBotID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID бота"})
		return
	}

	var updateData struct {
		Name       string `json:"name"`
		Username   string `json:"username"`
		Token      string `json:"token"`
		WebhookURL string `json:"webhook_url"`
		IsActive   *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bot, err := h.botManager.GetBot(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бот не найден"})
		return
	}

	// Обновляем поля
	if updateData.Name != "" {
		bot.Name = updateData.Name
	}
	if updateData.Username != "" {
		bot.Username = updateData.Username
	}
	if updateData.Token != "" {
		bot.Token = updateData.Token
	}
	if updateData.WebhookURL != "" {
		bot.WebhookURL = updateData.WebhookURL
	}
	if updateData.IsActive != nil {
		bot.IsActive = *updateData.IsActive
	}

	if err := h.botManager.UpdateBot(bot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bot": bot,
	})
}

// Delete удаляет бота
func (h *BotHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID бота обязателен"})
		return
	}

	id, err := ParseBotID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID бота"})
		return
	}

	if err := h.botManager.DeleteBot(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Бот успешно удален"})
}

// SendMessage отправляет сообщение через бота
func (h *BotHandler) SendMessage(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID бота обязателен"})
		return
	}

	id, err := ParseBotID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID бота"})
		return
	}

	var messageData struct {
		ChatID    int64  `json:"chat_id" binding:"required"`
		Text      string `json:"text" binding:"required"`
		ParseMode string `json:"parse_mode"`
	}

	if err := c.ShouldBindJSON(&messageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.botManager.SendBotMessage(id, messageData.ChatID, messageData.Text, messageData.ParseMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// GetUpdates возвращает обновления для бота
func (h *BotHandler) GetUpdates(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID бота обязателен"})
		return
	}

	id, err := ParseBotID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID бота"})
		return
	}

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

	updates, err := h.botManager.GetBotUpdates(id, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"updates": updates,
	})
}

// Webhook обрабатывает webhook от бота
func (h *BotHandler) Webhook(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID бота обязателен"})
		return
	}

	id, err := ParseBotID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID бота"})
		return
	}

	var update models.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.botManager.ProcessWebhook(id, &update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
