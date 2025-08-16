package main

import (
	"fmt"
	"log"
	"os"

	"telegram-emulator/internal/api"
	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"
	"telegram-emulator/internal/pkg/config"
	"telegram-emulator/internal/pkg/logger"
	"telegram-emulator/internal/repository"
	"telegram-emulator/internal/websocket"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.File); err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()

	log := logger.GetLogger()
	log.Info("Запуск Telegram эмулятора...")

	// Инициализация базы данных
	db, err := initDatabase(cfg.Database.URL)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных", zap.Error(err))
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	_ = repository.NewBotRepository(db) // Пока не используем

	// Инициализация WebSocket сервера
	wsServer := websocket.NewServer()
	
	// Инициализация менеджеров
	userManager := emulator.NewUserManager(userRepo)
	chatManager := emulator.NewChatManager(chatRepo, messageRepo, userRepo)
	messageManager := emulator.NewMessageManager(messageRepo, chatRepo, userRepo, wsServer)
	
	// Устанавливаем MessageManager в WebSocket сервер
	wsServer.SetMessageManager(messageManager)
	
	go wsServer.Start()

	// Создание тестовых данных
	if err := createTestData(userManager, chatManager); err != nil {
		log.Error("Ошибка создания тестовых данных", zap.Error(err))
	}

	log.Info("Telegram эмулятор успешно запущен!")
	log.Info(fmt.Sprintf("Сервер доступен по адресу: http://%s:%d", cfg.Emulator.Host, cfg.Emulator.Port))

	// Настройка и запуск HTTP сервера
	router := gin.Default()
	
	// Настройка CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Настройка маршрутов
	api.SetupRoutes(router, userManager, chatManager, messageManager, wsServer)

	// Запуск сервера
	addr := fmt.Sprintf("%s:%d", cfg.Emulator.Host, cfg.Emulator.Port)
	log.Info("HTTP сервер запускается", zap.String("address", addr))
	
	if err := router.Run(addr); err != nil {
		log.Fatal("Ошибка запуска HTTP сервера", zap.Error(err))
	}
}

// initDatabase инициализирует базу данных
func initDatabase(dbURL string) (*gorm.DB, error) {
	// Для SQLite используем простой путь к файлу
	dbPath := "data/emulator.db"
	
	// Создаем директорию для базы данных
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания директории для БД: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.Message{},
		&models.Bot{},
		&models.ChatMember{},
	); err != nil {
		return nil, fmt.Errorf("ошибка миграции БД: %w", err)
	}

	return db, nil
}

// createTestData создает тестовые данные
func createTestData(userManager *emulator.UserManager, chatManager *emulator.ChatManager) error {
	log := logger.GetLogger()

	// Создание тестовых пользователей
	users := []struct {
		username  string
		firstName string
		lastName  string
		isBot     bool
	}{
		{"test_user1", "Тестовый", "Пользователь 1", false},
		{"test_user2", "Тестовый", "Пользователь 2", false},
		{"test_bot", "Тестовый", "Бот", true},
	}

	for _, userData := range users {
		user, err := userManager.CreateUser(userData.username, userData.firstName, userData.lastName, userData.isBot)
		if err != nil {
			log.Error("Ошибка создания тестового пользователя", zap.Error(err))
			continue
		}
		log.Info("Создан тестовый пользователь", 
			zap.String("id", user.ID),
			zap.String("username", user.Username),
			zap.Bool("is_bot", user.IsBot))
	}

	// Получаем созданных пользователей
	allUsers, err := userManager.GetAllUsers()
	if err != nil {
		return err
	}

	if len(allUsers) >= 2 {
		// Создание приватного чата между первыми двумя пользователями
		chat, err := chatManager.CreatePrivateChat(allUsers[0].ID, allUsers[1].ID)
		if err != nil {
			log.Error("Ошибка создания приватного чата", zap.Error(err))
		} else {
			log.Info("Создан приватный чат", 
				zap.String("id", chat.ID),
				zap.String("title", chat.Title),
				zap.String("type", chat.Type))
		}
	}

	return nil
}
