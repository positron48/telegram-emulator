package handlers

import (
	"net/http"
	"strconv"

	"telegram-emulator/internal/emulator"

	"github.com/gin-gonic/gin"
)



// MessageHandler обрабатывает запросы к API сообщений
type MessageHandler struct {
	messageManager *emulator.MessageManager
}

// NewMessageHandler создает новый экземпляр MessageHandler
func NewMessageHandler(messageManager *emulator.MessageManager) *MessageHandler {
	return &MessageHandler{
		messageManager: messageManager,
	}
}



// SendMessageRequest представляет запрос на отправку сообщения
type SendMessageRequest struct {
	FromUserID int64  `json:"from_user_id" binding:"required"`
	Text       string `json:"text" binding:"required"`
	Type       string `json:"type"` // text, file, voice, photo
}

// UpdateMessageStatusRequest представляет запрос на обновление статуса сообщения
type UpdateMessageStatusRequest struct {
	Status string `json:"status" binding:"required"` // sending, sent, delivered, read
}

// GetByID получает сообщение по ID
func (h *MessageHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID сообщения обязателен"})
		return
	}

	id, err := ParseMessageID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID сообщения"})
		return
	}

	message, err := h.messageManager.GetMessage(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Сообщение не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// SendMessage отправляет сообщение в чат
func (h *MessageHandler) SendMessage(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
		return
	}
	
	chatID, err := ParseChatID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID чата"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем тип сообщения по умолчанию
	if req.Type == "" {
		req.Type = "text"
	}

	message, err := h.messageManager.SendMessage(chatID, req.FromUserID, req.Text, req.Type, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

// GetChatMessages получает сообщения чата
func (h *MessageHandler) GetChatMessages(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
		return
	}
	
	chatID, err := ParseChatID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID чата"})
		return
	}

	// Получаем параметры пагинации
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	messages, err := h.messageManager.GetChatMessages(chatID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateStatus обновляет статус сообщения
func (h *MessageHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID сообщения обязателен"})
		return
	}

	id, err := ParseMessageID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID сообщения"})
		return
	}

	var req UpdateMessageStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.messageManager.UpdateMessageStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Статус сообщения обновлен",
	})
}

// MarkChatAsRead помечает чат как прочитанный
func (h *MessageHandler) MarkChatAsRead(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
		return
	}
	
	chatID, err := ParseChatID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID чата"})
		return
	}
	
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	userID, err := ParseUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
		return
	}

	if err := h.messageManager.MarkChatAsRead(chatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Чат помечен как прочитанный",
	})
}

// DeleteMessage удаляет сообщение
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID сообщения обязателен"})
		return
	}

	id, err := ParseMessageID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID сообщения"})
		return
	}

	if err := h.messageManager.DeleteMessage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Сообщение удалено",
	})
}

// SearchMessages ищет сообщения по тексту
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
		return
	}
	
	chatID, err := ParseChatID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID чата"})
		return
	}
	
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поисковый запрос обязателен"})
		return
	}

	messages, err := h.messageManager.SearchMessages(chatID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"query":    query,
	})
}

// HandleCallbackQueryRequest представляет запрос на обработку callback query
type HandleCallbackQueryRequest struct {
	UserID       int64  `json:"user_id" binding:"required"`
	CallbackData string `json:"callback_data" binding:"required"`
}

// HandleCallbackQuery обрабатывает callback query от inline кнопки
func (h *MessageHandler) HandleCallbackQuery(c *gin.Context) {
	messageIDStr := c.Param("id")
	if messageIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID сообщения обязателен"})
		return
	}

	messageID, err := ParseMessageID(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID сообщения"})
		return
	}

	var req HandleCallbackQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	callbackQuery, err := h.messageManager.HandleCallbackQuery(req.UserID, messageID, req.CallbackData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"callback_query": callbackQuery,
	})
}
