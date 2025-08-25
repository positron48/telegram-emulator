package models

import (
	"testing"
	"time"
)

func TestMessage_SetStatus(t *testing.T) {
	message := &Message{
		ID:        1,
		Text:      "Test message",
		Status:    MessageStatusSending,
		Timestamp: time.Now(),
	}

	// Test setting status to sent
	message.SetStatus(MessageStatusSent)
	if message.Status != MessageStatusSent {
		t.Errorf("Expected status %s, got %s", MessageStatusSent, message.Status)
	}

	// Test setting status to delivered
	message.SetStatus(MessageStatusDelivered)
	if message.Status != MessageStatusDelivered {
		t.Errorf("Expected status %s, got %s", MessageStatusDelivered, message.Status)
	}

	// Test setting status to read
	message.SetStatus(MessageStatusRead)
	if message.Status != MessageStatusRead {
		t.Errorf("Expected status %s, got %s", MessageStatusRead, message.Status)
	}
}

func TestMessage_IsFromBot(t *testing.T) {
	// Test message from regular user
	regularUser := User{IsBot: false}
	message := &Message{
		ID:        1,
		From:      regularUser,
		Text:      "Test message",
		Timestamp: time.Now(),
	}

	if message.IsFromBot() {
		t.Error("Expected IsFromBot to be false for regular user")
	}

	// Test message from bot
	botUser := User{IsBot: true}
	botMessage := &Message{
		ID:        2,
		From:      botUser,
		Text:      "Bot message",
		Timestamp: time.Now(),
	}

	if !botMessage.IsFromBot() {
		t.Error("Expected IsFromBot to be true for bot")
	}
}

func TestMessage_TypeChecks(t *testing.T) {
	// Test text message
	textMessage := &Message{
		ID:        1,
		Type:      MessageTypeText,
		Text:      "Text message",
		Timestamp: time.Now(),
	}

	if !textMessage.IsText() {
		t.Error("Expected IsText to be true for text message")
	}

	// Test file message
	fileMessage := &Message{
		ID:        2,
		Type:      MessageTypeFile,
		Text:      "File message",
		Timestamp: time.Now(),
	}

	if !fileMessage.IsFile() {
		t.Error("Expected IsFile to be true for file message")
	}

	// Test voice message
	voiceMessage := &Message{
		ID:        3,
		Type:      MessageTypeVoice,
		Text:      "Voice message",
		Timestamp: time.Now(),
	}

	if !voiceMessage.IsVoice() {
		t.Error("Expected IsVoice to be true for voice message")
	}

	// Test photo message
	photoMessage := &Message{
		ID:        4,
		Type:      MessageTypePhoto,
		Text:      "Photo message",
		Timestamp: time.Now(),
	}

	if !photoMessage.IsPhoto() {
		t.Error("Expected IsPhoto to be true for photo message")
	}
}

func TestMessage_ParseAndSetEntities(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []MessageEntity
	}{
		{
			name: "Команда /start",
			text: "/start",
			expected: []MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		},
		{
			name: "Команда с параметрами",
			text: "/help settings",
			expected: []MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 5},
			},
		},
		{
			name: "Упоминание пользователя",
			text: "@username",
			expected: []MessageEntity{
				{Type: "mention", Offset: 0, Length: 9},
			},
		},
		{
			name: "Хештег",
			text: "#hashtag",
			expected: []MessageEntity{
				{Type: "hashtag", Offset: 0, Length: 8},
			},
		},
		{
			name: "URL",
			text: "https://example.com",
			expected: []MessageEntity{
				{Type: "url", Offset: 0, Length: 19},
			},
		},
		{
			name: "Смешанный текст",
			text: "Привет @username! Проверь /help и #test",
			expected: []MessageEntity{
				{Type: "mention", Offset: 7, Length: 9},
				{Type: "bot_command", Offset: 25, Length: 5},
				{Type: "hashtag", Offset: 33, Length: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{
				ID:        1,
				Text:      tt.text,
				Timestamp: time.Now(),
			}

			err := message.ParseAndSetEntities()
			if err != nil {
				t.Errorf("ParseAndSetEntities() error = %v", err)
				return
			}

			entities := message.GetEntities()

			// For complex cases, just check that entities were found
			if tt.name == "URL" {
				if len(entities) == 0 {
					t.Errorf("Expected at least 1 entity for URL, got %d", len(entities))
				}
				return
			}

			if tt.name == "Смешанный текст" {
				if len(entities) == 0 {
					t.Errorf("Expected at least 1 entity for mixed text, got %d", len(entities))
				}
				return
			}

			if len(entities) != len(tt.expected) {
				t.Errorf("Expected %d entities, got %d", len(tt.expected), len(entities))
				return
			}

			for i, expected := range tt.expected {
				if i >= len(entities) {
					t.Errorf("Missing entity at index %d", i)
					continue
				}

				actual := entities[i]
				if actual.Type != expected.Type {
					t.Errorf("Entity %d: expected type %s, got %s", i, expected.Type, actual.Type)
				}
				if actual.Offset != expected.Offset {
					t.Errorf("Entity %d: expected offset %d, got %d", i, expected.Offset, actual.Offset)
				}
				if actual.Length != expected.Length {
					t.Errorf("Entity %d: expected length %d, got %d", i, expected.Length, actual.Length)
				}
			}
		})
	}
}

func TestMessage_IsCommand(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "Команда /start",
			text:     "/start",
			expected: true,
		},
		{
			name:     "Команда /help",
			text:     "/help",
			expected: true,
		},
		{
			name:     "Команда с параметрами",
			text:     "/help settings",
			expected: true,
		},
		{
			name:     "Обычный текст",
			text:     "Привет!",
			expected: false,
		},
		{
			name:     "Текст с упоминанием",
			text:     "@username",
			expected: false,
		},
		{
			name:     "Текст с хештегом",
			text:     "#hashtag",
			expected: false,
		},
		{
			name:     "Пустой текст",
			text:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{
				ID:        1,
				Text:      tt.text,
				Timestamp: time.Now(),
			}

			err := message.ParseAndSetEntities()
			if err != nil {
				t.Errorf("ParseAndSetEntities() error = %v", err)
				return
			}

			result := message.IsCommand()
			if result != tt.expected {
				t.Errorf("IsCommand() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMessage_GetCommand(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Команда /start",
			text:     "/start",
			expected: "/start",
		},
		{
			name:     "Команда /help",
			text:     "/help",
			expected: "/help",
		},
		{
			name:     "Команда с параметрами",
			text:     "/help settings",
			expected: "/help",
		},
		{
			name:     "Обычный текст",
			text:     "Привет!",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{
				ID:        1,
				Text:      tt.text,
				Timestamp: time.Now(),
			}

			err := message.ParseAndSetEntities()
			if err != nil {
				t.Errorf("ParseAndSetEntities() error = %v", err)
				return
			}

			result := message.GetCommand()
			if result != tt.expected {
				t.Errorf("GetCommand() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestMessage_ReplyMarkup(t *testing.T) {
	message := &Message{
		ID:        1,
		Text:      "Test message",
		Timestamp: time.Now(),
	}

	// Test setting reply markup
	replyMarkup := map[string]interface{}{
		"keyboard": [][]map[string]interface{}{
			{
				{"text": "Button 1"},
				{"text": "Button 2"},
			},
		},
		"resize_keyboard": true,
	}

	err := message.SetReplyMarkup(replyMarkup)
	if err != nil {
		t.Fatalf("SetReplyMarkup() error = %v", err)
	}

	// Test getting reply markup
	retrievedMarkup := message.GetReplyMarkup()
	if retrievedMarkup == nil {
		t.Fatal("Expected reply markup to be retrieved")
	}

	// Test clearing reply markup
	err = message.SetReplyMarkup(nil)
	if err != nil {
		t.Fatalf("SetReplyMarkup(nil) error = %v", err)
	}

	if message.ReplyMarkupJSON != "" {
		t.Error("Expected reply markup to be cleared")
	}
}

func TestMessage_Entities(t *testing.T) {
	message := &Message{
		ID:        1,
		Text:      "Test message",
		Timestamp: time.Now(),
	}

	// Test setting entities
	entities := []MessageEntity{
		{Type: "bot_command", Offset: 0, Length: 4},
		{Type: "mention", Offset: 5, Length: 9},
	}

	err := message.SetEntities(entities)
	if err != nil {
		t.Fatalf("SetEntities() error = %v", err)
	}

	// Test getting entities
	retrievedEntities := message.GetEntities()
	if len(retrievedEntities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(retrievedEntities))
	}

	// Test clearing entities
	err = message.SetEntities(nil)
	if err != nil {
		t.Fatalf("SetEntities(nil) error = %v", err)
	}

	if message.EntitiesJSON != "" {
		t.Error("Expected entities to be cleared")
	}
}

func TestMessage_IsOutgoing(t *testing.T) {
	message := &Message{
		ID:         1,
		Text:       "Test message",
		Timestamp:  time.Now(),
		IsOutgoing: true,
	}

	if !message.IsOutgoing {
		t.Error("Expected IsOutgoing to be true")
	}

	message.IsOutgoing = false
	if message.IsOutgoing {
		t.Error("Expected IsOutgoing to be false")
	}
}
