package models

import (
	"time"
)

// Update представляет обновление от Telegram Bot API
type Update struct {
	UpdateID      int64    `json:"update_id"`
	Message       *Message `json:"message,omitempty"`
	EditedMessage *Message `json:"edited_message,omitempty"`
	ChannelPost   *Message `json:"channel_post,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// CallbackQuery представляет callback query от inline кнопок
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data,omitempty"`
	GameShortName   string   `json:"game_short_name,omitempty"`
}

// TelegramMessage представляет сообщение в формате Telegram Bot API
type TelegramMessage struct {
	MessageID int64  `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text,omitempty"`
}

// TelegramUser представляет пользователя в формате Telegram Bot API
type TelegramUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// TelegramChat представляет чат в формате Telegram Bot API
type TelegramChat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// ToTelegramMessage конвертирует внутреннее сообщение в формат Telegram Bot API
func (m *Message) ToTelegramMessage() TelegramMessage {
	// Конвертируем строковые ID в int64 (используем хеш для уникальности)
	messageID := int64(0)
	if len(m.ID) > 0 {
		for i, char := range m.ID {
			if i < 8 { // Ограничиваем длину
				messageID = messageID*31 + int64(char)
			}
		}
	}

	fromID := int64(0)
	if len(m.FromID) > 0 {
		for i, char := range m.FromID {
			if i < 8 { // Ограничиваем длину
				fromID = fromID*31 + int64(char)
			}
		}
	}

	chatID := int64(0)
	if len(m.ChatID) > 0 {
		for i, char := range m.ChatID {
			if i < 8 { // Ограничиваем длину
				chatID = chatID*31 + int64(char)
			}
		}
	}

	return TelegramMessage{
		MessageID: messageID,
		From: TelegramUser{
			ID:        fromID,
			IsBot:     m.From.IsBot,
			FirstName: m.From.FirstName,
			LastName:  m.From.LastName,
			Username:  m.From.Username,
		},
		Chat: TelegramChat{
			ID:       chatID,
			Type:     "private", // TODO: получить тип чата
			Title:    m.From.GetFullName(),
			Username: m.From.Username,
		},
		Date: m.Timestamp.Unix(),
		Text: m.Text,
	}
}

// FromTelegramMessage конвертирует сообщение из формата Telegram Bot API во внутренний формат
func FromTelegramMessage(tgMsg TelegramMessage, chatID string) *Message {
	return &Message{
		ID:        string(tgMsg.MessageID),
		ChatID:    chatID,
		FromID:    string(tgMsg.From.ID),
		Text:      tgMsg.Text,
		Type:      MessageTypeText,
		Status:    MessageStatusSent,
		Timestamp: time.Unix(tgMsg.Date, 0),
		CreatedAt: time.Now(),
	}
}
