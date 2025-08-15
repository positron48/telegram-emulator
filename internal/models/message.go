package models

import (
	"time"
)

// Message представляет сообщение в эмуляторе
type Message struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ChatID    string    `json:"chat_id"`
	FromID    string    `json:"from_id"`
	From      User      `json:"from" gorm:"foreignKey:FromID"`
	Text      string    `json:"text"`
	Type      string    `json:"type"` // text, file, voice, photo
	Status    string    `json:"status"` // sending, sent, delivered, read
	IsOutgoing bool     `json:"is_outgoing"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName возвращает имя таблицы для модели Message
func (Message) TableName() string {
	return "messages"
}

// MessageStatus представляет статусы сообщений
const (
	MessageStatusSending   = "sending"
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusRead      = "read"
)

// MessageType представляет типы сообщений
const (
	MessageTypeText  = "text"
	MessageTypeFile  = "file"
	MessageTypeVoice = "voice"
	MessageTypePhoto = "photo"
)

// SetStatus устанавливает статус сообщения
func (m *Message) SetStatus(status string) {
	m.Status = status
}

// IsFromBot проверяет, отправлено ли сообщение ботом
func (m *Message) IsFromBot() bool {
	return m.From.IsBot
}

// IsText проверяет, является ли сообщение текстовым
func (m *Message) IsText() bool {
	return m.Type == MessageTypeText
}

// IsFile проверяет, является ли сообщение файлом
func (m *Message) IsFile() bool {
	return m.Type == MessageTypeFile
}

// IsVoice проверяет, является ли сообщение голосовым
func (m *Message) IsVoice() bool {
	return m.Type == MessageTypeVoice
}

// IsPhoto проверяет, является ли сообщение фотографией
func (m *Message) IsPhoto() bool {
	return m.Type == MessageTypePhoto
}
