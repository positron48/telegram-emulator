package models

import (
	"encoding/json"
	"regexp"
	"time"
)

// Message представляет сообщение в эмуляторе
type Message struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	ChatID          int64     `json:"chat_id"`
	FromID          int64     `json:"from_id"`
	From            User      `json:"from" gorm:"foreignKey:FromID"`
	Text            string    `json:"text"`
	Type            string    `json:"type"`   // text, file, voice, photo
	Status          string    `json:"status"` // sending, sent, delivered, read
	IsOutgoing      bool      `json:"is_outgoing"`
	Timestamp       time.Time `json:"timestamp"`
	CreatedAt       time.Time `json:"created_at"`
	ReplyMarkupJSON string    `json:"reply_markup,omitempty" gorm:"column:reply_markup"` // Клавиатура в JSON формате
	EntitiesJSON    string    `json:"entities,omitempty" gorm:"column:entities"`         // Сущности в JSON формате
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

// SetReplyMarkup устанавливает клавиатуру и сериализует её в JSON
func (m *Message) SetReplyMarkup(replyMarkup interface{}) error {
	if replyMarkup == nil {
		m.ReplyMarkupJSON = ""
		return nil
	}

	// Сериализуем клавиатуру в JSON
	jsonData, err := json.Marshal(replyMarkup)
	if err != nil {
		return err
	}

	m.ReplyMarkupJSON = string(jsonData)
	return nil
}

// GetReplyMarkup десериализует клавиатуру из JSON
func (m *Message) GetReplyMarkup() interface{} {
	if m.ReplyMarkupJSON == "" {
		return nil
	}

	var replyMarkup interface{}
	if err := json.Unmarshal([]byte(m.ReplyMarkupJSON), &replyMarkup); err != nil {
		return nil
	}

	return replyMarkup
}

// SetEntities устанавливает сущности и сериализует их в JSON
func (m *Message) SetEntities(entities []MessageEntity) error {
	if len(entities) == 0 {
		m.EntitiesJSON = ""
		return nil
	}

	// Сериализуем сущности в JSON
	jsonData, err := json.Marshal(entities)
	if err != nil {
		return err
	}

	m.EntitiesJSON = string(jsonData)
	return nil
}

// GetEntities десериализует сущности из JSON
func (m *Message) GetEntities() []MessageEntity {
	if m.EntitiesJSON == "" {
		return nil
	}

	var entities []MessageEntity
	if err := json.Unmarshal([]byte(m.EntitiesJSON), &entities); err != nil {
		return nil
	}

	return entities
}

// ParseAndSetEntities парсит текст сообщения и устанавливает сущности
func (m *Message) ParseAndSetEntities() error {
	if m.Text == "" {
		return nil
	}

	var entities []MessageEntity

	// Парсим команды
	commandEntities := parseCommands(m.Text)
	entities = append(entities, commandEntities...)

	// Парсим упоминания
	mentionEntities := parseMentions(m.Text)
	entities = append(entities, mentionEntities...)

	// Парсим URL (до хештегов, чтобы избежать конфликтов)
	urlEntities := parseURLs(m.Text)
	entities = append(entities, urlEntities...)

	// Парсим хештеги (после URL)
	hashtagEntities := parseHashtags(m.Text)
	entities = append(entities, hashtagEntities...)

	return m.SetEntities(entities)
}

// IsCommand проверяет, является ли сообщение командой
func (m *Message) IsCommand() bool {
	entities := m.GetEntities()
	for _, entity := range entities {
		if entity.Type == "bot_command" {
			return true
		}
	}
	return false
}

// GetCommand возвращает команду из сообщения
func (m *Message) GetCommand() string {
	entities := m.GetEntities()
	for _, entity := range entities {
		if entity.Type == "bot_command" {
			if entity.Offset < len(m.Text) && entity.Offset+entity.Length <= len(m.Text) {
				return m.Text[entity.Offset : entity.Offset+entity.Length]
			}
		}
	}
	return ""
}

// parseCommands парсит команды в тексте
func parseCommands(text string) []MessageEntity {
	var entities []MessageEntity
	commandRegex := regexp.MustCompile(`/([a-zA-Z0-9_]+)`)

	matches := commandRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		entities = append(entities, MessageEntity{
			Type:   "bot_command",
			Offset: match[0],
			Length: match[1] - match[0],
		})
	}

	return entities
}

// parseMentions парсит упоминания в тексте
func parseMentions(text string) []MessageEntity {
	var entities []MessageEntity
	mentionRegex := regexp.MustCompile(`@([a-zA-Z0-9_]{5,32})`)

	matches := mentionRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		entities = append(entities, MessageEntity{
			Type:   "mention",
			Offset: match[0],
			Length: match[1] - match[0],
		})
	}

	return entities
}

// parseHashtags парсит хештеги в тексте
func parseHashtags(text string) []MessageEntity {
	var entities []MessageEntity
	hashtagRegex := regexp.MustCompile(`#([a-zA-Z0-9_]+)`)

	matches := hashtagRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		entities = append(entities, MessageEntity{
			Type:   "hashtag",
			Offset: match[0],
			Length: match[1] - match[0],
		})
	}

	return entities
}

// parseURLs парсит URL в тексте
func parseURLs(text string) []MessageEntity {
	var entities []MessageEntity
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)

	matches := urlRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		entities = append(entities, MessageEntity{
			Type:   "url",
			Offset: match[0],
			Length: match[1] - match[0],
		})
	}

	return entities
}
