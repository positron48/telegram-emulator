package websocket

import "telegram-emulator/internal/models"

// MessageManagerInterface определяет интерфейс для MessageManager
type MessageManagerInterface interface {
	SendMessage(chatID int64, fromUserID int64, text, messageType string, replyMarkup interface{}) (*models.Message, error)
}
