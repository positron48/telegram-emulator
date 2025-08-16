package emulator

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/pkg/logger"
	"telegram-emulator/internal/repository"

	"go.uber.org/zap"
)

// ChatManager управляет чатами в эмуляторе
type ChatManager struct {
	chatRepo    *repository.ChatRepository
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
	logger      *zap.Logger
}

// NewChatManager создает новый экземпляр ChatManager
func NewChatManager(chatRepo *repository.ChatRepository, messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *ChatManager {
	return &ChatManager{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		userRepo:    userRepo,
		logger:      logger.GetLogger(),
	}
}

// CreateChat создает новый чат
func (m *ChatManager) CreateChat(chatType, title, username, description string, userIDs []string) (*models.Chat, error) {
	// Генерируем уникальный ID
	id, err := m.generateID()
	if err != nil {
		return nil, err
	}

	chat := &models.Chat{
		ID:          id,
		Type:        chatType,
		Title:       title,
		Username:    username,
		Description: description,
		UnreadCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Создаем чат
	if err := m.chatRepo.Create(chat); err != nil {
		m.logger.Error("Ошибка создания чата", zap.Error(err))
		return nil, err
	}

	// Добавляем участников
	for _, userID := range userIDs {
		if err := m.chatRepo.AddMember(chat.ID, userID); err != nil {
			m.logger.Error("Ошибка добавления участника в чат", zap.String("chat_id", chat.ID), zap.String("user_id", userID), zap.Error(err))
			return nil, err
		}
	}

	// Загружаем участников
	chat.Members, err = m.chatRepo.GetMembers(chat.ID)
	if err != nil {
		m.logger.Error("Ошибка загрузки участников чата", zap.String("chat_id", chat.ID), zap.Error(err))
		return nil, err
	}

	m.logger.Info("Создан новый чат", 
		zap.String("id", chat.ID),
		zap.String("type", chat.Type),
		zap.String("title", chat.Title))

	return chat, nil
}

// CreatePrivateChat создает приватный чат между двумя пользователями
func (m *ChatManager) CreatePrivateChat(userID1, userID2 string) (*models.Chat, error) {
	// Проверяем, существует ли уже приватный чат
	existingChat, err := m.chatRepo.GetPrivateChat(userID1, userID2)
	if err == nil {
		return existingChat, nil
	}

	// Получаем информацию о пользователях
	user1, err := m.userRepo.GetByID(userID1)
	if err != nil {
		return nil, err
	}

	user2, err := m.userRepo.GetByID(userID2)
	if err != nil {
		return nil, err
	}

	// Создаем название чата
	title := user1.GetFullName()
	if user2.GetFullName() != "" {
		title = user2.GetFullName()
	}

	return m.CreateChat("private", title, "", "", []string{userID1, userID2})
}

// GetChat получает чат по ID
func (m *ChatManager) GetChat(id string) (*models.Chat, error) {
	chat, err := m.chatRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения чата", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return chat, nil
}

// GetUserChats получает чаты пользователя
func (m *ChatManager) GetUserChats(userID string) ([]models.Chat, error) {
	chats, err := m.chatRepo.GetByUserID(userID)
	if err != nil {
		m.logger.Error("Ошибка получения чатов пользователя", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return chats, nil
}

// GetAllChats получает все чаты
func (m *ChatManager) GetAllChats() ([]models.Chat, error) {
	chats, err := m.chatRepo.GetAll()
	if err != nil {
		m.logger.Error("Ошибка получения всех чатов", zap.Error(err))
		return nil, err
	}
	return chats, nil
}

// UpdateChat обновляет чат
func (m *ChatManager) UpdateChat(chat *models.Chat) error {
	chat.UpdatedAt = time.Now()
	if err := m.chatRepo.Update(chat); err != nil {
		m.logger.Error("Ошибка обновления чата", zap.String("id", chat.ID), zap.Error(err))
		return err
	}

	m.logger.Info("Чат обновлен", zap.String("id", chat.ID))
	return nil
}

// DeleteChat удаляет чат
func (m *ChatManager) DeleteChat(id string) error {
	if err := m.chatRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления чата", zap.String("id", id), zap.Error(err))
		return err
	}

	m.logger.Info("Чат удален", zap.String("id", id))
	return nil
}

// AddMember добавляет пользователя в чат
func (m *ChatManager) AddMember(chatID, userID string) error {
	if err := m.chatRepo.AddMember(chatID, userID); err != nil {
		m.logger.Error("Ошибка добавления участника", zap.String("chat_id", chatID), zap.String("user_id", userID), zap.Error(err))
		return err
	}

	m.logger.Info("Участник добавлен в чат", zap.String("chat_id", chatID), zap.String("user_id", userID))
	return nil
}

// RemoveMember удаляет пользователя из чата
func (m *ChatManager) RemoveMember(chatID, userID string) error {
	if err := m.chatRepo.RemoveMember(chatID, userID); err != nil {
		m.logger.Error("Ошибка удаления участника", zap.String("chat_id", chatID), zap.String("user_id", userID), zap.Error(err))
		return err
	}

	m.logger.Info("Участник удален из чата", zap.String("chat_id", chatID), zap.String("user_id", userID))
	return nil
}

// GetChatMembers получает участников чата
func (m *ChatManager) GetChatMembers(chatID string) ([]models.User, error) {
	members, err := m.chatRepo.GetMembers(chatID)
	if err != nil {
		m.logger.Error("Ошибка получения участников чата", zap.String("chat_id", chatID), zap.Error(err))
		return nil, err
	}
	return members, nil
}

// GetMembers алиас для GetChatMembers
func (m *ChatManager) GetMembers(chatID string) ([]models.User, error) {
	return m.GetChatMembers(chatID)
}

// UpdateUnreadCount обновляет количество непрочитанных сообщений
func (m *ChatManager) UpdateUnreadCount(chatID string) error {
	count, err := m.messageRepo.GetUnreadCount(chatID)
	if err != nil {
		return err
	}

	if err := m.chatRepo.UpdateUnreadCount(chatID, int(count)); err != nil {
		m.logger.Error("Ошибка обновления счетчика непрочитанных", zap.String("chat_id", chatID), zap.Error(err))
		return err
	}

	return nil
}

// MarkChatAsRead помечает все сообщения чата как прочитанные
func (m *ChatManager) MarkChatAsRead(chatID string) error {
	if err := m.messageRepo.MarkAsRead(chatID); err != nil {
		m.logger.Error("Ошибка пометки чата как прочитанного", zap.String("chat_id", chatID), zap.Error(err))
		return err
	}

	// Обновляем счетчик непрочитанных
	if err := m.chatRepo.UpdateUnreadCount(chatID, 0); err != nil {
		m.logger.Error("Ошибка обновления счетчика непрочитанных", zap.String("chat_id", chatID), zap.Error(err))
		return err
	}

	m.logger.Info("Чат помечен как прочитанный", zap.String("chat_id", chatID))
	return nil
}

// generateID генерирует уникальный ID
func (m *ChatManager) generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
