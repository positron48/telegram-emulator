package emulator

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMessageManager_SendMessage(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create a user
	user, err := messageManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a chat
	chat, err := messageManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Send a message
	message, err := messageManager.SendMessage(chat.ID, user.ID, "Hello, world!", "text")
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	if message.ChatID != chat.ID {
		t.Errorf("Expected chat ID %d, got %d", chat.ID, message.ChatID)
	}

	if message.FromID != user.ID {
		t.Errorf("Expected from ID %d, got %d", user.ID, message.FromID)
	}

	if message.Text != "Hello, world!" {
		t.Errorf("Expected text 'Hello, world!', got '%s'", message.Text)
	}

	if message.Type != "text" {
		t.Errorf("Expected type 'text', got '%s'", message.Type)
	}

	if message.Status != "sent" {
		t.Errorf("Expected status 'sent', got '%s'", message.Status)
	}

	if message.ID == 0 {
		t.Error("Expected message ID to be generated")
	}
}

func TestMessageManager_GetMessage(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create a user and chat
	user, err := messageManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	chat, err := messageManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Send a message
	sentMessage, err := messageManager.SendMessage(chat.ID, user.ID, "Test message", "text")
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Get the message
	retrievedMessage, err := messageManager.GetMessage(sentMessage.ID)
	if err != nil {
		t.Fatalf("Failed to get message: %v", err)
	}

	if retrievedMessage.ID != sentMessage.ID {
		t.Errorf("Expected message ID %d, got %d", sentMessage.ID, retrievedMessage.ID)
	}

	if retrievedMessage.Text != sentMessage.Text {
		t.Errorf("Expected text '%s', got '%s'", sentMessage.Text, retrievedMessage.Text)
	}
}

func TestMessageManager_GetChatMessages(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create users and chat
	user1, err := messageManager.CreateUser("user1", "User", "One", false)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	user2, err := messageManager.CreateUser("user2", "User", "Two", false)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	chat, err := messageManager.CreateChat("group", "Test Group", "testgroup")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Send messages
	_, err = messageManager.SendMessage(chat.ID, user1.ID, "Message 1", "text")
	if err != nil {
		t.Fatalf("Failed to send message 1: %v", err)
	}

	_, err = messageManager.SendMessage(chat.ID, user2.ID, "Message 2", "text")
	if err != nil {
		t.Fatalf("Failed to send message 2: %v", err)
	}

	// Get chat messages
	messages, err := messageManager.GetChatMessages(chat.ID)
	if err != nil {
		t.Fatalf("Failed to get chat messages: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	// Check that all messages belong to the chat
	for _, msg := range messages {
		if msg.ChatID != chat.ID {
			t.Errorf("Expected chat ID %d, got %d", chat.ID, msg.ChatID)
		}
	}
}

func TestMessageManager_UpdateMessage(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create a user and chat
	user, err := messageManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	chat, err := messageManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Send a message
	message, err := messageManager.SendMessage(chat.ID, user.ID, "Original message", "text")
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Update message
	message.Text = "Updated message"
	message.Status = "delivered"
	updatedMessage, err := messageManager.UpdateMessage(message)
	if err != nil {
		t.Fatalf("Failed to update message: %v", err)
	}

	if updatedMessage.Text != "Updated message" {
		t.Errorf("Expected text 'Updated message', got '%s'", updatedMessage.Text)
	}

	if updatedMessage.Status != "delivered" {
		t.Errorf("Expected status 'delivered', got '%s'", updatedMessage.Status)
	}
}

func TestMessageManager_DeleteMessage(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create a user and chat
	user, err := messageManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	chat, err := messageManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Send a message
	message, err := messageManager.SendMessage(chat.ID, user.ID, "Test message", "text")
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Delete message
	err = messageManager.DeleteMessage(message.ID)
	if err != nil {
		t.Fatalf("Failed to delete message: %v", err)
	}

	// Try to get the deleted message
	_, err = messageManager.GetMessage(message.ID)
	if err == nil {
		t.Error("Expected error when getting deleted message")
	}
}

func TestMessageManager_GetMessageNotFound(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Try to get non-existent message
	_, err := messageManager.GetMessage(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent message")
	}
}

func TestMessageManager_UpdateMessageNotFound(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Try to update non-existent message
	message := &models.Message{
		ID:        999999,
		ChatID:    1,
		FromID:    1,
		Text:      "Non-existent message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	_, err := messageManager.UpdateMessage(message)
	if err == nil {
		t.Error("Expected error when updating non-existent message")
	}
}

func TestMessageManager_DeleteMessageNotFound(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Try to delete non-existent message
	err := messageManager.DeleteMessage(999999)
	if err == nil {
		t.Error("Expected error when deleting non-existent message")
	}
}

func TestMessageManager_GetUserMessages(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create users and chat
	user1, err := messageManager.CreateUser("user1", "User", "One", false)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	user2, err := messageManager.CreateUser("user2", "User", "Two", false)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	chat, err := messageManager.CreateChat("group", "Test Group", "testgroup")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Send messages from different users
	_, err = messageManager.SendMessage(chat.ID, user1.ID, "Message from user 1", "text")
	if err != nil {
		t.Fatalf("Failed to send message from user1: %v", err)
	}

	_, err = messageManager.SendMessage(chat.ID, user1.ID, "Another message from user 1", "text")
	if err != nil {
		t.Fatalf("Failed to send second message from user1: %v", err)
	}

	_, err = messageManager.SendMessage(chat.ID, user2.ID, "Message from user 2", "text")
	if err != nil {
		t.Fatalf("Failed to send message from user2: %v", err)
	}

	// Get messages from user 1
	messages, err := messageManager.GetUserMessages(user1.ID)
	if err != nil {
		t.Fatalf("Failed to get user messages: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages from user 1, got %d", len(messages))
	}

	// Check that all messages are from user 1
	for _, msg := range messages {
		if msg.FromID != user1.ID {
			t.Errorf("Expected from ID %d, got %d", user1.ID, msg.FromID)
		}
	}
}

func TestMessageManager_SendMessageWithReplyMarkup(t *testing.T) {
	db := setupTestDB(t)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageManager := NewMessageManager(messageRepo, userRepo, chatRepo)

	// Create a user and chat
	user, err := messageManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	chat, err := messageManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Create reply markup
	replyMarkup := map[string]interface{}{
		"keyboard": [][]map[string]interface{}{
			{
				{"text": "Button 1"},
				{"text": "Button 2"},
			},
		},
		"resize_keyboard": true,
	}

	// Send message with reply markup
	message, err := messageManager.SendMessageWithReplyMarkup(chat.ID, user.ID, "Message with keyboard", "text", replyMarkup)
	if err != nil {
		t.Fatalf("Failed to send message with reply markup: %v", err)
	}

	if message.ReplyMarkupJSON == "" {
		t.Error("Expected reply markup to be set")
	}

	// Test getting reply markup
	retrievedMarkup := message.GetReplyMarkup()
	if retrievedMarkup == nil {
		t.Error("Expected reply markup to be retrieved")
	}
}
