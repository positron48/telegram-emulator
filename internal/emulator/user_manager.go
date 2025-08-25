package emulator

import (
	"math/rand"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/pkg/logger"
	"telegram-emulator/internal/repository"

	"go.uber.org/zap"
)

// UserManager управляет пользователями в эмуляторе
type UserManager struct {
	userRepo *repository.UserRepository
	botRepo  *repository.BotRepository
	logger   *zap.Logger
}

// NewUserManager создает новый экземпляр UserManager
func NewUserManager(userRepo *repository.UserRepository, botRepo *repository.BotRepository) *UserManager {
	return &UserManager{
		userRepo: userRepo,
		botRepo:  botRepo,
		logger:   logger.GetLogger(),
	}
}

// CreateUser создает нового пользователя
func (m *UserManager) CreateUser(username, firstName, lastName string, isBot bool) (*models.User, error) {
	// Генерируем уникальный ID
	id, err := m.generateID()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        id,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		IsBot:     isBot,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := m.userRepo.Create(user); err != nil {
		m.logger.Error("Ошибка создания пользователя", zap.Error(err))
		return nil, err
	}

	// Если пользователь является ботом, создаем запись в таблице ботов
	if user.IsBot {
		bot := &models.Bot{
			ID:         user.ID,
			Name:       user.GetFullName(),
			Username:   user.Username,
			Token:      "", // Токен можно будет установить позже
			WebhookURL: "",
			IsActive:   true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := m.botRepo.Create(bot); err != nil {
			m.logger.Error("Ошибка создания записи бота", zap.Error(err))
			// Не удаляем пользователя, просто логируем ошибку
		} else {
			m.logger.Info("Создана запись бота",
				zap.Int64("bot_id", bot.ID),
				zap.String("username", bot.Username))
		}
	}

	m.logger.Info("Создан новый пользователь",
		zap.Int64("id", user.ID),
		zap.String("username", user.Username),
		zap.Bool("is_bot", user.IsBot))

	return user, nil
}

// GetUser получает пользователя по ID
func (m *UserManager) GetUser(id int64) (*models.User, error) {
	user, err := m.userRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetUserByUsername получает пользователя по username
func (m *UserManager) GetUserByUsername(username string) (*models.User, error) {
	user, err := m.userRepo.GetByUsername(username)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя по username", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetAllUsers получает всех пользователей
func (m *UserManager) GetAllUsers() ([]models.User, error) {
	users, err := m.userRepo.GetAll()
	if err != nil {
		m.logger.Error("Ошибка получения всех пользователей", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// UpdateUser обновляет пользователя
func (m *UserManager) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := m.userRepo.Update(user); err != nil {
		m.logger.Error("Ошибка обновления пользователя", zap.Int64("id", user.ID), zap.Error(err))
		return err
	}

	m.logger.Info("Пользователь обновлен", zap.Int64("id", user.ID))
	return nil
}

// DeleteUser удаляет пользователя
func (m *UserManager) DeleteUser(id int64) error {
	// Получаем пользователя перед удалением
	user, err := m.userRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения пользователя для удаления", zap.Int64("id", id), zap.Error(err))
		return err
	}

	// Если пользователь является ботом, удаляем запись бота
	if user.IsBot {
		if err := m.botRepo.Delete(id); err != nil {
			m.logger.Error("Ошибка удаления записи бота", zap.Int64("id", id), zap.Error(err))
			// Не прерываем удаление пользователя, просто логируем ошибку
		} else {
			m.logger.Info("Запись бота удалена", zap.Int64("id", id))
		}
	}

	if err := m.userRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления пользователя", zap.Int64("id", id), zap.Error(err))
		return err
	}

	m.logger.Info("Пользователь удален", zap.Int64("id", id))
	return nil
}

// SetUserOnline устанавливает статус онлайн для пользователя
func (m *UserManager) SetUserOnline(id int64, isOnline bool) error {
	if err := m.userRepo.SetOnlineStatus(id, isOnline); err != nil {
		m.logger.Error("Ошибка установки статуса онлайн", zap.Int64("id", id), zap.Bool("online", isOnline), zap.Error(err))
		return err
	}

	status := "онлайн"
	if !isOnline {
		status = "оффлайн"
	}
	m.logger.Info("Пользователь изменил статус", zap.Int64("id", id), zap.String("status", status))
	return nil
}

// GetOnlineUsers получает всех онлайн пользователей
func (m *UserManager) GetOnlineUsers() ([]models.User, error) {
	users, err := m.userRepo.GetOnlineUsers()
	if err != nil {
		m.logger.Error("Ошибка получения онлайн пользователей", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// GetBots получает всех ботов
func (m *UserManager) GetBots() ([]models.User, error) {
	users, err := m.userRepo.GetBots()
	if err != nil {
		m.logger.Error("Ошибка получения ботов", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// generateID генерирует уникальный ID
func (m *UserManager) generateID() (int64, error) {
	// Используем Unix timestamp в миллисекундах + случайное число для уникальности
	timestamp := time.Now().UnixMilli()
	random := rand.Int63n(1000) // случайное число от 0 до 999
	return timestamp*1000 + random, nil
}
