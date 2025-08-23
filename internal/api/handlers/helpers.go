package handlers

import (
	"strconv"
)

// ParseBotID парсит ID бота из строки в int64
func ParseBotID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

// ParseChatID парсит ID чата из строки в int64
func ParseChatID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

// ParseMessageID парсит ID сообщения из строки в int64
func ParseMessageID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

// ParseUserID парсит ID пользователя из строки в int64
func ParseUserID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}
