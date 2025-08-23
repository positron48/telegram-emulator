package models

import (
	"time"
)

// Bot представляет бота в эмуляторе
type Bot struct {
	ID                int64     `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Username          string    `json:"username" gorm:"uniqueIndex"`
	Token             string    `json:"token"`
	WebhookURL        string    `json:"webhook_url"`
	IsActive          bool      `json:"is_active"`
	LastUpdateOffset  int64     `json:"last_update_offset" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName возвращает имя таблицы для модели Bot
func (Bot) TableName() string {
	return "bots"
}

// Activate активирует бота
func (b *Bot) Activate() {
	b.IsActive = true
}

// Deactivate деактивирует бота
func (b *Bot) Deactivate() {
	b.IsActive = false
}

// SetWebhook устанавливает webhook URL для бота
func (b *Bot) SetWebhook(url string) {
	b.WebhookURL = url
}

// UpdateToken обновляет токен бота
func (b *Bot) UpdateToken(token string) {
	b.Token = token
}

// BotNotFoundError представляет ошибку "бот не найден"
type BotNotFoundError struct{}

func (e *BotNotFoundError) Error() string {
	return "бот не найден"
}
