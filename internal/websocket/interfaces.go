package websocket

import "telegram-emulator/internal/models"

// MessageManagerInterface определяет интерфейс для MessageManager
type MessageManagerInterface interface {
	SendMessage(chatID, fromUserID, text, messageType string, replyMarkup interface{}) (*models.Message, error)
}
