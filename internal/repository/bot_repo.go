package repository

import (
	"telegram-emulator/internal/models"

	"gorm.io/gorm"
)

// BotRepository управляет операциями с ботами в базе данных
type BotRepository struct {
	db *gorm.DB
}

// NewBotRepository создает новый экземпляр BotRepository
func NewBotRepository(db *gorm.DB) *BotRepository {
	return &BotRepository{db: db}
}

// Create создает нового бота
func (r *BotRepository) Create(bot *models.Bot) error {
	return r.db.Create(bot).Error
}

// GetByID получает бота по ID
func (r *BotRepository) GetByID(id int64) (*models.Bot, error) {
	var bot models.Bot
	err := r.db.Where("id = ?", id).First(&bot).Error
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// GetByUsername получает бота по username
func (r *BotRepository) GetByUsername(username string) (*models.Bot, error) {
	var bot models.Bot
	err := r.db.Where("username = ?", username).First(&bot).Error
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// GetAll получает всех ботов
func (r *BotRepository) GetAll() ([]models.Bot, error) {
	var bots []models.Bot
	err := r.db.Find(&bots).Error
	return bots, err
}

// GetActive получает всех активных ботов
func (r *BotRepository) GetActive() ([]models.Bot, error) {
	var bots []models.Bot
	err := r.db.Where("is_active = ?", true).Find(&bots).Error
	return bots, err
}

// Update обновляет бота
func (r *BotRepository) Update(bot *models.Bot) error {
	return r.db.Save(bot).Error
}

// Delete удаляет бота
func (r *BotRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Bot{}).Error
}

// SetActiveStatus устанавливает статус активности бота
func (r *BotRepository) SetActiveStatus(id int64, isActive bool) error {
	return r.db.Model(&models.Bot{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// SetWebhookURL устанавливает webhook URL для бота
func (r *BotRepository) SetWebhookURL(id int64, webhookURL string) error {
	return r.db.Model(&models.Bot{}).Where("id = ?", id).Update("webhook_url", webhookURL).Error
}

// UpdateToken обновляет токен бота
func (r *BotRepository) UpdateToken(id int64, token string) error {
	return r.db.Model(&models.Bot{}).Where("id = ?", id).Update("token", token).Error
}
