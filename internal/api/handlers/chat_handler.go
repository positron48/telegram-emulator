package handlers

import (
	"net/http"
	"strconv"

	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"

	"github.com/gin-gonic/gin"
)



// ChatHandler обрабатывает запросы к API чатов
type ChatHandler struct {
	chatManager *emulator.ChatManager
}

// NewChatHandler создает новый экземпляр ChatHandler
func NewChatHandler(chatManager *emulator.ChatManager) *ChatHandler {
	return &ChatHandler{
		chatManager: chatManager,
	}
}

// CreateChatRequest представляет запрос на создание чата
type CreateChatRequest struct {
	Type        string   `json:"type" binding:"required"` // private, group
	Title       string   `json:"title" binding:"required"`
	Username    string   `json:"username"`
	Description string   `json:"description"`
	UserIDs     []int64  `json:"user_ids" binding:"required"`
}

// UpdateChatRequest представляет запрос на обновление чата
type UpdateChatRequest struct {
	Title       string `json:"title"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

// GetAll получает все чаты
func (h *ChatHandler) GetAll(c *gin.Context) {
	// Получаем user_id из query параметра
	userIDStr := c.Query("user_id")
	
	var chats []models.Chat
	var err error
	
	if userIDStr != "" {
		// Парсим userID
		userID, parseErr := ParseUserID(userIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
			return
		}
		// Возвращаем только чаты пользователя
		chats, err = h.chatManager.GetUserChats(userID)
	} else {
		// Возвращаем все чаты (для совместимости)
		chats, err = h.chatManager.GetAllChats()
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
		"count": len(chats),
	})
}

// Create создает новый чат
func (h *ChatHandler) Create(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := h.chatManager.CreateChat(req.Type, req.Title, req.Username, req.Description, req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"chat": chat,
	})
}

// ParseChatID конвертирует строковый ID чата в int64
func (h *ChatHandler) ParseChatID(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
		return 0, false
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID чата"})
		return 0, false
	}

	return id, true
}

// GetByID получает чат по ID
func (h *ChatHandler) GetByID(c *gin.Context) {
	id, ok := h.ParseChatID(c)
	if !ok {
		return
	}

	chat, err := h.chatManager.GetChat(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Чат не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat": chat,
	})
}

// Update обновляет чат
func (h *ChatHandler) Update(c *gin.Context) {
	id, ok := h.ParseChatID(c)
	if !ok {
		return
	}

	var req UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущий чат
	chat, err := h.chatManager.GetChat(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Чат не найден"})
		return
	}

	// Обновляем поля
	if req.Title != "" {
		chat.Title = req.Title
	}
	if req.Username != "" {
		chat.Username = req.Username
	}
	if req.Description != "" {
		chat.Description = req.Description
	}

	// Сохраняем изменения
	if err := h.chatManager.UpdateChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat": chat,
	})
}

// Delete удаляет чат
func (h *ChatHandler) Delete(c *gin.Context) {
	id, ok := h.ParseChatID(c)
	if !ok {
		return
	}

	if err := h.chatManager.DeleteChat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Чат успешно удален",
	})
}

// GetMessages получает сообщения чата
func (h *ChatHandler) GetMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
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

	// TODO: Добавить метод GetChatMessages в ChatManager
	// messages, err := h.chatManager.GetChatMessages(id, limit, offset)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"messages": []models.Message{},
		"limit": limit,
		"offset": offset,
	})
}

// AddMember добавляет участника в чат
func (h *ChatHandler) AddMember(c *gin.Context) {
	chatID, ok := h.ParseChatID(c)
	if !ok {
		return
	}

	var req struct {
		UserID int64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.chatManager.AddMember(chatID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Участник успешно добавлен",
	})
}

// GetMembers получает участников чата
func (h *ChatHandler) GetMembers(c *gin.Context) {
	chatID, ok := h.ParseChatID(c)
	if !ok {
		return
	}

	members, err := h.chatManager.GetMembers(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
		"count":   len(members),
	})
}

// RemoveMember удаляет участника из чата
func (h *ChatHandler) RemoveMember(c *gin.Context) {
	chatID, ok := h.ParseChatID(c)
	if !ok {
		return
	}
	
	userIDStr := c.Param("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	userID, err := ParseUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
		return
	}

	if err := h.chatManager.RemoveMember(chatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Участник успешно удален",
	})
}
