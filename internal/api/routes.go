package api

import (
	"telegram-emulator/internal/api/handlers"
	"telegram-emulator/internal/emulator"

	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты API
func SetupRoutes(router *gin.Engine, userManager *emulator.UserManager, chatManager *emulator.ChatManager) {
	// API группа
	api := router.Group("/api")
	{
		// Пользователи
		users := api.Group("/users")
		{
			userHandler := handlers.NewUserHandler(userManager)
			users.GET("", userHandler.GetAll)
			users.POST("", userHandler.Create)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.GET("/:id/chats", userHandler.GetChats)
		}

		// Чаты
		chats := api.Group("/chats")
		{
			chatHandler := handlers.NewChatHandler(chatManager)
			chats.GET("", chatHandler.GetAll)
			chats.POST("", chatHandler.Create)
			chats.GET("/:id", chatHandler.GetByID)
			chats.PUT("/:id", chatHandler.Update)
			chats.DELETE("/:id", chatHandler.Delete)
			chats.GET("/:id/messages", chatHandler.GetMessages)
			chats.POST("/:id/members", chatHandler.AddMember)
			chats.DELETE("/:id/members/:userID", chatHandler.RemoveMember)
		}

		// Сообщения
		messages := api.Group("/messages")
		{
			messageHandler := handlers.NewMessageHandler(chatManager)
			messages.GET("/:id", messageHandler.GetByID)
			messages.PUT("/:id/status", messageHandler.UpdateStatus)
		}
	}

	// Статические файлы для веб-интерфейса
	router.Static("/static", "./web/public")
	router.LoadHTMLGlob("web/templates/*")
	
	// Главная страница
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Telegram Emulator",
		})
	})
}
