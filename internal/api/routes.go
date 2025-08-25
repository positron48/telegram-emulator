package api

import (
	"telegram-emulator/internal/api/handlers"
	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты API
func SetupRoutes(router *gin.Engine, userManager *emulator.UserManager, chatManager *emulator.ChatManager, messageManager *emulator.MessageManager, botManager *emulator.BotManager, wsServer *websocket.Server) {
	// Telegram Bot API
	telegramAPI := NewTelegramBotAPI(botManager, userManager, chatManager, messageManager)
	telegramAPI.SetupTelegramBotRoutes(router)
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
			chats.GET("/:id/members", chatHandler.GetMembers)
			chats.POST("/:id/members", chatHandler.AddMember)
			chats.DELETE("/:id/members/:userID", chatHandler.RemoveMember)
		}

		// Сообщения
		messages := api.Group("/messages")
		{
			messageHandler := handlers.NewMessageHandler(messageManager)
			messages.GET("/:id", messageHandler.GetByID)
			messages.PUT("/:id/status", messageHandler.UpdateStatus)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
			messages.POST("/:id/callback", messageHandler.HandleCallbackQuery)
		}

		// Сообщения чатов
		messageHandler := handlers.NewMessageHandler(messageManager)
		chats.GET("/:id/messages", messageHandler.GetChatMessages)
		chats.POST("/:id/messages", messageHandler.SendMessage)
		chats.PUT("/:id/read", messageHandler.MarkChatAsRead)
		chats.GET("/:id/search", messageHandler.SearchMessages)
	}

	// Боты
	bots := api.Group("/bots")
	{
		botHandler := handlers.NewBotHandler(botManager)
		bots.GET("", botHandler.GetAll)
		bots.POST("", botHandler.Create)
		bots.GET("/:id", botHandler.GetByID)
		bots.PUT("/:id", botHandler.Update)
		bots.DELETE("/:id", botHandler.Delete)
		bots.POST("/:id/sendMessage", botHandler.SendMessage)
		bots.GET("/:id/updates", botHandler.GetUpdates)
		bots.POST("/:id/webhook", botHandler.Webhook)
	}

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		wsServer.HandleWebSocket(c.Writer, c.Request)
	})

	// Статические файлы для веб-интерфейса
	router.Static("/web", "./web/dist")
	router.StaticFile("/", "./web/dist/index.html")
	router.StaticFile("/index.html", "./web/dist/index.html")

	// Главная страница (fallback)
	router.GET("/api", func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSON(200, gin.H{
			"message": "Telegram Emulator API",
			"version": "1.0.0",
			"docs":    "/api",
		})
	})
}
