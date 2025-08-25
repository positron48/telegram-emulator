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
func (r *MessageRepository) GetByID(id int64) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("From").Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetByChatID получает сообщения чата
func (r *MessageRepository) GetByChatID(chatID int64, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	query := r.db.Where("chat_id = ?", chatID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	// Загружаем связанные данные
	for i := range messages {
		if err := r.db.Model(&messages[i]).Association("From").Find(&messages[i].From); err != nil {
			return nil, err
		}
	}

	return messages, nil
}

// Update обновляет сообщение
func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

// Delete удаляет сообщение
func (r *MessageRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Message{}).Error
}

// UpdateStatus обновляет статус сообщения
func (r *MessageRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&models.Message{}).Where("id = ?", id).Update("status", status).Error
}

// GetLastMessage получает последнее сообщение чата
func (r *MessageRepository) GetLastMessage(chatID int64) (*models.Message, error) {
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
func (r *MessageRepository) GetUnreadCount(chatID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("chat_id = ? AND status != ?", chatID, models.MessageStatusRead).
		Count(&count).Error
	return count, err
}

// MarkAsRead помечает сообщения как прочитанные
func (r *MessageRepository) MarkAsRead(chatID int64) error {
	return r.db.Model(&models.Message{}).
		Where("chat_id = ? AND status != ?", chatID, models.MessageStatusRead).
		Update("status", models.MessageStatusRead).Error
}

// GetByType получает сообщения определенного типа
func (r *MessageRepository) GetByType(chatID int64, messageType string) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("From").
		Where("chat_id = ? AND type = ?", chatID, messageType).
		Order("timestamp DESC").
		Find(&messages).Error
	return messages, err
}

// SearchByText ищет сообщения по тексту
func (r *MessageRepository) SearchByText(chatID int64, text string) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("From").
		Where("chat_id = ? AND text LIKE ?", chatID, "%"+text+"%").
		Order("timestamp DESC").
		Find(&messages).Error
	return messages, err
}
