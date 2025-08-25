package repository

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"
)

func TestMessageRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	message := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Test message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	if message.ID == 0 {
		t.Error("Expected message ID to be set after creation")
	}
}

func TestMessageRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create a message
	message := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Test message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Get the message by ID
	retrievedMessage, err := repo.GetByID(message.ID)
	if err != nil {
		t.Fatalf("Failed to get message by ID: %v", err)
	}

	if retrievedMessage.ID != message.ID {
		t.Errorf("Expected message ID %d, got %d", message.ID, retrievedMessage.ID)
	}

	if retrievedMessage.Text != message.Text {
		t.Errorf("Expected text '%s', got '%s'", message.Text, retrievedMessage.Text)
	}
}

func TestMessageRepository_GetByChatID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create messages for different chats
	message1 := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Message in chat 1",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	message2 := &models.Message{
		ChatID:    1,
		FromID:    2,
		Text:      "Another message in chat 1",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	message3 := &models.Message{
		ChatID:    2,
		FromID:    1,
		Text:      "Message in chat 2",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	err = repo.Create(message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	err = repo.Create(message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	// Get messages for chat 1
	messages, err := repo.GetByChatID(1, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get messages by chat ID: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages for chat 1, got %d", len(messages))
	}

	// Check that all messages belong to chat 1
	for _, msg := range messages {
		if msg.ChatID != 1 {
			t.Errorf("Expected chat ID 1, got %d", msg.ChatID)
		}
	}
}

func TestMessageRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create a message
	message := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Test message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Update message
	message.Text = "Updated message"
	message.Status = "delivered"

	err = repo.Update(message)
	if err != nil {
		t.Fatalf("Failed to update message: %v", err)
	}

	// Get the updated message
	retrievedMessage, err := repo.GetByID(message.ID)
	if err != nil {
		t.Fatalf("Failed to get updated message: %v", err)
	}

	if retrievedMessage.Text != "Updated message" {
		t.Errorf("Expected text 'Updated message', got '%s'", retrievedMessage.Text)
	}

	if retrievedMessage.Status != "delivered" {
		t.Errorf("Expected status 'delivered', got '%s'", retrievedMessage.Status)
	}
}

func TestMessageRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create a message
	message := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Test message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Delete message
	err = repo.Delete(message.ID)
	if err != nil {
		t.Fatalf("Failed to delete message: %v", err)
	}

	// Try to get the deleted message
	_, err = repo.GetByID(message.ID)
	if err == nil {
		t.Error("Expected error when getting deleted message")
	}
}

func TestMessageRepository_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Try to get non-existent message
	_, err := repo.GetByID(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent message")
	}
}

func TestMessageRepository_UpdateNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

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

	err := repo.Update(message)
	// GORM doesn't return an error when updating non-existent records
	// It just doesn't update anything
	if err != nil {
		t.Errorf("Unexpected error when updating non-existent message: %v", err)
	}
}

func TestMessageRepository_DeleteNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Try to delete non-existent message
	err := repo.Delete(999999)
	// GORM doesn't return an error when deleting non-existent records
	// It just doesn't delete anything
	if err != nil {
		t.Errorf("Unexpected error when deleting non-existent message: %v", err)
	}
}

func TestMessageRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create a message
	message := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Test message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Update status
	err = repo.UpdateStatus(message.ID, "delivered")
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Get the message and check status
	retrievedMessage, err := repo.GetByID(message.ID)
	if err != nil {
		t.Fatalf("Failed to get message: %v", err)
	}

	if retrievedMessage.Status != "delivered" {
		t.Errorf("Expected status 'delivered', got '%s'", retrievedMessage.Status)
	}
}

func TestMessageRepository_GetLastMessage(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create messages for a chat
	message1 := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "First message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-time.Hour),
	}

	message2 := &models.Message{
		ChatID:    1,
		FromID:    2,
		Text:      "Last message",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	err = repo.Create(message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	// Get last message
	lastMessage, err := repo.GetLastMessage(1)
	if err != nil {
		t.Fatalf("Failed to get last message: %v", err)
	}

	if lastMessage.Text != "Last message" {
		t.Errorf("Expected text 'Last message', got '%s'", lastMessage.Text)
	}
}

func TestMessageRepository_GetUnreadCount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create messages with different statuses
	message1 := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Unread message 1",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	message2 := &models.Message{
		ChatID:    1,
		FromID:    2,
		Text:      "Unread message 2",
		Type:      "text",
		Status:    "delivered",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	message3 := &models.Message{
		ChatID:    1,
		FromID:    3,
		Text:      "Read message",
		Type:      "text",
		Status:    "read",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	err = repo.Create(message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	err = repo.Create(message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	// Get unread count
	count, err := repo.GetUnreadCount(1)
	if err != nil {
		t.Fatalf("Failed to get unread count: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 unread messages, got %d", count)
	}
}

func TestMessageRepository_MarkAsRead(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMessageRepository(db)

	// Create unread messages
	message1 := &models.Message{
		ChatID:    1,
		FromID:    1,
		Text:      "Unread message 1",
		Type:      "text",
		Status:    "sent",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	message2 := &models.Message{
		ChatID:    1,
		FromID:    2,
		Text:      "Unread message 2",
		Type:      "text",
		Status:    "delivered",
		IsOutgoing: false,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	err = repo.Create(message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	// Mark as read
	err = repo.MarkAsRead(1)
	if err != nil {
		t.Fatalf("Failed to mark as read: %v", err)
	}

	// Check unread count
	count, err := repo.GetUnreadCount(1)
	if err != nil {
		t.Fatalf("Failed to get unread count: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 unread messages, got %d", count)
	}
}
