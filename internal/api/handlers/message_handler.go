package handlers

import (
	"net/http"

	"telegram-emulator/internal/emulator"

	"github.com/gin-gonic/gin"
)

// MessageHandler обрабатывает запросы к API сообщений
type MessageHandler struct {
	chatManager *emulator.ChatManager
}

// NewMessageHandler создает новый экземпляр MessageHandler
func NewMessageHandler(chatManager *emulator.ChatManager) *MessageHandler {
	return &MessageHandler{
		chatManager: chatManager,
	}
}

// UpdateMessageStatusRequest представляет запрос на обновление статуса сообщения
type UpdateMessageStatusRequest struct {
	Status string `json:"status" binding:"required"` // sending, sent, delivered, read
}

// GetByID получает сообщение по ID
func (h *MessageHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID сообщения обязателен"})
		return
	}

	// TODO: Добавить метод GetMessage в ChatManager
	// message, err := h.chatManager.GetMessage(id)
	// if err != nil {
	//     c.JSON(http.StatusNotFound, gin.H{"error": "Сообщение не найдено"})
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Метод пока не реализован",
	})
}

// UpdateStatus обновляет статус сообщения
func (h *MessageHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID сообщения обязателен"})
		return
	}

	var req UpdateMessageStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Добавить метод UpdateMessageStatus в ChatManager
	// if err := h.chatManager.UpdateMessageStatus(id, req.Status); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Статус сообщения обновлен",
	})
}
