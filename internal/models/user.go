package models

import (
	"time"
)

// User представляет пользователя в эмуляторе
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsBot     bool      `json:"is_bot"`
	IsOnline  bool      `json:"is_online"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName возвращает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}

// GetFullName возвращает полное имя пользователя
func (u *User) GetFullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// SetOnline устанавливает статус онлайн
func (u *User) SetOnline(online bool) {
	u.IsOnline = online
	if !online {
		u.LastSeen = time.Now()
	}
}
