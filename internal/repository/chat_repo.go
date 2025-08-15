package repository

import (
	"time"

	"telegram-emulator/internal/models"

	"gorm.io/gorm"
)

// ChatRepository представляет репозиторий для работы с чатами
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository создает новый экземпляр ChatRepository
func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// Create создает новый чат
func (r *ChatRepository) Create(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

// GetByID получает чат по ID
func (r *ChatRepository) GetByID(id string) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Preload("Members").Preload("LastMessage").Where("id = ?", id).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// GetAll получает все чаты
func (r *ChatRepository) GetAll() ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.Preload("Members").Preload("LastMessage").Find(&chats).Error
	return chats, err
}

// GetByUserID получает чаты пользователя
func (r *ChatRepository) GetByUserID(userID string) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.Preload("Members").Preload("LastMessage").
		Joins("JOIN chat_members ON chats.id = chat_members.chat_id").
		Where("chat_members.user_id = ?", userID).
		Find(&chats).Error
	return chats, err
}

// Update обновляет чат
func (r *ChatRepository) Update(chat *models.Chat) error {
	return r.db.Save(chat).Error
}

// Delete удаляет чат
func (r *ChatRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Chat{}).Error
}

// AddMember добавляет пользователя в чат
func (r *ChatRepository) AddMember(chatID, userID string) error {
	member := models.ChatMember{
		ChatID:   chatID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}
	return r.db.Create(&member).Error
}

// RemoveMember удаляет пользователя из чата
func (r *ChatRepository) RemoveMember(chatID, userID string) error {
	return r.db.Where("chat_id = ? AND user_id = ?", chatID, userID).Delete(&models.ChatMember{}).Error
}

// GetMembers получает участников чата
func (r *ChatRepository) GetMembers(chatID string) ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN chat_members ON users.id = chat_members.user_id").
		Where("chat_members.chat_id = ?", chatID).
		Find(&users).Error
	return users, err
}

// UpdateUnreadCount обновляет количество непрочитанных сообщений
func (r *ChatRepository) UpdateUnreadCount(chatID string, count int) error {
	return r.db.Model(&models.Chat{}).Where("id = ?", chatID).Update("unread_count", count).Error
}

// GetPrivateChat получает приватный чат между двумя пользователями
func (r *ChatRepository) GetPrivateChat(userID1, userID2 string) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Preload("Members").
		Joins("JOIN chat_members cm1 ON chats.id = cm1.chat_id").
		Joins("JOIN chat_members cm2 ON chats.id = cm2.chat_id").
		Where("chats.type = ? AND cm1.user_id = ? AND cm2.user_id = ?", "private", userID1, userID2).
		First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}
