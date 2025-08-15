package models

import (
	"time"
)

// Chat представляет чат в эмуляторе
type Chat struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Type        string    `json:"type"` // private, group, channel
	Title       string    `json:"title"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Members     []User    `json:"members" gorm:"many2many:chat_members;"`
	LastMessage *Message  `json:"last_message" gorm:"foreignKey:ChatID"`
	UnreadCount int       `json:"unread_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName возвращает имя таблицы для модели Chat
func (Chat) TableName() string {
	return "chats"
}

// ChatMember представляет связь между чатом и пользователем
type ChatMember struct {
	ChatID   string    `json:"chat_id" gorm:"primaryKey"`
	UserID   string    `json:"user_id" gorm:"primaryKey"`
	JoinedAt time.Time `json:"joined_at"`
}

// TableName возвращает имя таблицы для модели ChatMember
func (ChatMember) TableName() string {
	return "chat_members"
}

// IsPrivate проверяет, является ли чат приватным
func (c *Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup проверяет, является ли чат группой
func (c *Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsChannel проверяет, является ли чат каналом
func (c *Chat) IsChannel() bool {
	return c.Type == "channel"
}

// AddMember добавляет пользователя в чат
func (c *Chat) AddMember(user User) {
	// Проверяем, не добавлен ли уже пользователь
	for _, member := range c.Members {
		if member.ID == user.ID {
			return
		}
	}
	c.Members = append(c.Members, user)
}

// RemoveMember удаляет пользователя из чата
func (c *Chat) RemoveMember(userID string) {
	for i, member := range c.Members {
		if member.ID == userID {
			c.Members = append(c.Members[:i], c.Members[i+1:]...)
			break
		}
	}
}
