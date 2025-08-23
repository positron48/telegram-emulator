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
func (m *ChatManager) CreateChat(chatType, title, username, description string, userIDs []int64) (*models.Chat, error) {
	// Генерируем уникальный ID
	id, err := m.generateID()
	if err != nil {
		return nil, err
	}

	// Конвертируем строковый ID в int64
	chatID := hashString(id)
	
	chat := &models.Chat{
		ID:          chatID,
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
			m.logger.Error("Ошибка добавления участника в чат", zap.Int64("chat_id", chat.ID), zap.Int64("user_id", userID), zap.Error(err))
			return nil, err
		}
	}

	// Загружаем участников
	chat.Members, err = m.chatRepo.GetMembers(chat.ID)
	if err != nil {
		m.logger.Error("Ошибка загрузки участников чата", zap.Int64("chat_id", chat.ID), zap.Error(err))
		return nil, err
	}

	m.logger.Info("Создан новый чат", 
		zap.Int64("id", chat.ID),
		zap.String("type", chat.Type),
		zap.String("title", chat.Title))

	return chat, nil
}

// CreatePrivateChat создает приватный чат между двумя пользователями
func (m *ChatManager) CreatePrivateChat(userID1, userID2 int64) (*models.Chat, error) {
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

	return m.CreateChat("private", title, "", "", []int64{userID1, userID2})
}

// GetChat получает чат по ID
func (m *ChatManager) GetChat(id int64) (*models.Chat, error) {
	chat, err := m.chatRepo.GetByID(id)
	if err != nil {
		m.logger.Error("Ошибка получения чата", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return chat, nil
}

// GetUserChats получает чаты пользователя
func (m *ChatManager) GetUserChats(userID int64) ([]models.Chat, error) {
	chats, err := m.chatRepo.GetByUserID(userID)
	if err != nil {
		m.logger.Error("Ошибка получения чатов пользователя", zap.Int64("user_id", userID), zap.Error(err))
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
		m.logger.Error("Ошибка обновления чата", zap.Int64("id", chat.ID), zap.Error(err))
		return err
	}

	m.logger.Info("Чат обновлен", zap.Int64("id", chat.ID))
	return nil
}

// DeleteChat удаляет чат
func (m *ChatManager) DeleteChat(id int64) error {
	if err := m.chatRepo.Delete(id); err != nil {
		m.logger.Error("Ошибка удаления чата", zap.Int64("id", id), zap.Error(err))
		return err
	}

	m.logger.Info("Чат удален", zap.Int64("id", id))
	return nil
}

// AddMember добавляет участника в чат
func (m *ChatManager) AddMember(chatID int64, userID int64) error {
	if err := m.chatRepo.AddMember(chatID, userID); err != nil {
		m.logger.Error("Ошибка добавления участника", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID), zap.Error(err))
		return err
	}
	m.logger.Info("Участник добавлен в чат", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
	return nil
}

// RemoveMember удаляет участника из чата
func (m *ChatManager) RemoveMember(chatID int64, userID int64) error {
	if err := m.chatRepo.RemoveMember(chatID, userID); err != nil {
		m.logger.Error("Ошибка удаления участника", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID), zap.Error(err))
		return err
	}
	m.logger.Info("Участник удален из чата", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
	return nil
}

// GetChatMembers получает участников чата
func (m *ChatManager) GetChatMembers(chatID int64) ([]models.User, error) {
	return m.chatRepo.GetMembers(chatID)
}

// GetMembers получает участников чата
func (m *ChatManager) GetMembers(chatID int64) ([]models.User, error) {
	return m.chatRepo.GetMembers(chatID)
}

// UpdateUnreadCount обновляет счетчик непрочитанных сообщений
func (m *ChatManager) UpdateUnreadCount(chatID int64) error {
	// Получаем количество непрочитанных сообщений
	unreadCount, err := m.messageRepo.GetUnreadCount(chatID)
	if err != nil {
		return err
	}
	
	// Обновляем счетчик в чате
	return m.chatRepo.UpdateUnreadCount(chatID, int(unreadCount))
}

// MarkChatAsRead помечает чат как прочитанный
func (m *ChatManager) MarkChatAsRead(chatID int64) error {
	// Помечаем все сообщения как прочитанные
	if err := m.messageRepo.MarkAsRead(chatID); err != nil {
		return err
	}
	
	// Обновляем счетчик непрочитанных
	return m.UpdateUnreadCount(chatID)
}

// generateID генерирует уникальный ID
func (m *ChatManager) generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashString создает простой хеш из строки для конвертации в int64
func hashString(s string) int64 {
	hash := int64(0)
	for i, char := range s {
		if i < 8 { // Ограничиваем длину для предотвращения переполнения
			hash = hash*31 + int64(char)
		}
	}
	return hash
}
