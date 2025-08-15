package handlers

import (
	"net/http"
	"strconv"

	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"

	"github.com/gin-gonic/gin"
)

// UserHandler обрабатывает запросы к API пользователей
type UserHandler struct {
	userManager *emulator.UserManager
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userManager *emulator.UserManager) *UserHandler {
	return &UserHandler{
		userManager: userManager,
	}
}

// CreateUserRequest представляет запрос на создание пользователя
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	IsBot     bool   `json:"is_bot"`
}

// UpdateUserRequest представляет запрос на обновление пользователя
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsOnline  *bool  `json:"is_online"`
}

// GetAll получает всех пользователей
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userManager.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// Create создает нового пользователя
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userManager.CreateUser(req.Username, req.FirstName, req.LastName, req.IsBot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

// GetByID получает пользователя по ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	user, err := h.userManager.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Update обновляет пользователя
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущего пользователя
	user, err := h.userManager.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Обновляем поля
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.IsOnline != nil {
		user.SetOnline(*req.IsOnline)
	}

	// Сохраняем изменения
	if err := h.userManager.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Delete удаляет пользователя
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	if err := h.userManager.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно удален",
	})
}

// GetChats получает чаты пользователя
func (h *UserHandler) GetChats(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
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

	// TODO: Добавить метод GetUserChats в ChatManager
	// chats, err := h.chatManager.GetUserChats(id, limit, offset)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"chats": []models.Chat{},
		"limit": limit,
		"offset": offset,
	})
}
