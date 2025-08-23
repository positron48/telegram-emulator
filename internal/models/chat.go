package models

import (
	"time"
)

// Chat представляет чат в эмуляторе
type Chat struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Type        string    `json:"type"` // private, group
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
	ChatID   int64     `json:"chat_id" gorm:"primaryKey"`
	UserID   int64     `json:"user_id" gorm:"primaryKey"`
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



// GetChatIcon возвращает иконку для типа чата
func (c *Chat) GetChatIcon() string {
	switch c.Type {
	case "private":
		return "👤"
	case "group":
		return "👥"
	default:
		return "💬"
	}
}

// GetChatTypeLabel возвращает человекочитаемое название типа чата
func (c *Chat) GetChatTypeLabel() string {
	switch c.Type {
	case "private":
		return "Приватный чат"
	case "group":
		return "Группа"
	default:
		return "Чат"
	}
}

// CanUserJoin проверяет, может ли пользователь присоединиться к чату
func (c *Chat) CanUserJoin() bool {
	return c.Type == "group"
}

// CanUserLeave проверяет, может ли пользователь покинуть чат
func (c *Chat) CanUserLeave() bool {
	return c.Type == "group"
}

// IsUserMember проверяет, является ли пользователь участником чата
func (c *Chat) IsUserMember(userID int64) bool {
	for _, member := range c.Members {
		if member.ID == userID {
			return true
		}
	}
	return false
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
func (c *Chat) RemoveMember(userID int64) {
	for i, member := range c.Members {
		if member.ID == userID {
			c.Members = append(c.Members[:i], c.Members[i+1:]...)
			break
		}
	}
}
