package repository

import (
	"telegram-emulator/internal/models"

	"gorm.io/gorm"
)

// MessageRepository представляет репозиторий для работы с сообщениями
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository создает новый экземпляр MessageRepository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create создает новое сообщение
func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// GetByID получает сообщение по ID
func (r *MessageRepository) GetByID(id string) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("From").Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetByChatID получает сообщения чата
func (r *MessageRepository) GetByChatID(chatID string, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("From").
		Where("chat_id = ?", chatID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// Update обновляет сообщение
func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

// Delete удаляет сообщение
func (r *MessageRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Message{}).Error
}

// UpdateStatus обновляет статус сообщения
func (r *MessageRepository) UpdateStatus(id, status string) error {
	return r.db.Model(&models.Message{}).Where("id = ?", id).Update("status", status).Error
}

// GetLastMessage получает последнее сообщение чата
func (r *MessageRepository) GetLastMessage(chatID string) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("From").
		Where("chat_id = ?", chatID).
		Order("timestamp DESC").
		First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetUnreadCount получает количество непрочитанных сообщений в чате
func (r *MessageRepository) GetUnreadCount(chatID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("chat_id = ? AND status != ?", chatID, models.MessageStatusRead).
		Count(&count).Error
	return count, err
}

// MarkAsRead помечает сообщения как прочитанные
func (r *MessageRepository) MarkAsRead(chatID string) error {
	return r.db.Model(&models.Message{}).
		Where("chat_id = ? AND status != ?", chatID, models.MessageStatusRead).
		Update("status", models.MessageStatusRead).Error
}

// GetByType получает сообщения определенного типа
func (r *MessageRepository) GetByType(chatID, messageType string) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("From").
		Where("chat_id = ? AND type = ?", chatID, messageType).
		Order("timestamp DESC").
		Find(&messages).Error
	return messages, err
}

// SearchByText ищет сообщения по тексту
func (r *MessageRepository) SearchByText(chatID, text string) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("From").
		Where("chat_id = ? AND text LIKE ?", chatID, "%"+text+"%").
		Order("timestamp DESC").
		Find(&messages).Error
	return messages, err
}
