package repository

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"
)

func TestChatRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	chat := &models.Chat{
		Type:      "private",
		Title:     "Test Chat",
		Username:  "testchat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(chat)
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	if chat.ID == 0 {
		t.Error("Expected chat ID to be set after creation")
	}
}

func TestChatRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Create a chat
	chat := &models.Chat{
		Type:      "private",
		Title:     "Test Chat",
		Username:  "testchat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(chat)
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Get the chat by ID
	retrievedChat, err := repo.GetByID(chat.ID)
	if err != nil {
		t.Fatalf("Failed to get chat by ID: %v", err)
	}

	if retrievedChat.ID != chat.ID {
		t.Errorf("Expected chat ID %d, got %d", chat.ID, retrievedChat.ID)
	}

	if retrievedChat.Title != chat.Title {
		t.Errorf("Expected title '%s', got '%s'", chat.Title, retrievedChat.Title)
	}
}

func TestChatRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Create multiple chats
	chat1 := &models.Chat{
		Type:      "private",
		Title:     "Chat 1",
		Username:  "chat1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	chat2 := &models.Chat{
		Type:      "group",
		Title:     "Chat 2",
		Username:  "chat2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(chat1)
	if err != nil {
		t.Fatalf("Failed to create chat1: %v", err)
	}

	err = repo.Create(chat2)
	if err != nil {
		t.Fatalf("Failed to create chat2: %v", err)
	}

	// Get all chats
	chats, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all chats: %v", err)
	}

	if len(chats) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(chats))
	}
}

func TestChatRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Create a chat
	chat := &models.Chat{
		Type:      "private",
		Title:     "Test Chat",
		Username:  "testchat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(chat)
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Update chat
	chat.Title = "Updated Chat"
	chat.Type = "group"

	err = repo.Update(chat)
	if err != nil {
		t.Fatalf("Failed to update chat: %v", err)
	}

	// Get the updated chat
	retrievedChat, err := repo.GetByID(chat.ID)
	if err != nil {
		t.Fatalf("Failed to get updated chat: %v", err)
	}

	if retrievedChat.Title != "Updated Chat" {
		t.Errorf("Expected title 'Updated Chat', got '%s'", retrievedChat.Title)
	}

	if retrievedChat.Type != "group" {
		t.Errorf("Expected type 'group', got '%s'", retrievedChat.Type)
	}
}

func TestChatRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Create a chat
	chat := &models.Chat{
		Type:      "private",
		Title:     "Test Chat",
		Username:  "testchat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(chat)
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Delete chat
	err = repo.Delete(chat.ID)
	if err != nil {
		t.Fatalf("Failed to delete chat: %v", err)
	}

	// Try to get the deleted chat
	_, err = repo.GetByID(chat.ID)
	if err == nil {
		t.Error("Expected error when getting deleted chat")
	}
}

func TestChatRepository_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Try to get non-existent chat
	_, err := repo.GetByID(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent chat")
	}
}

func TestChatRepository_UpdateNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Try to update non-existent chat
	chat := &models.Chat{
		ID:        999999,
		Type:      "private",
		Title:     "Non-existent Chat",
		Username:  "nonexistent",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Update(chat)
	// GORM doesn't return an error when updating non-existent records
	// It just doesn't update anything
	if err != nil {
		t.Errorf("Unexpected error when updating non-existent chat: %v", err)
	}
}

func TestChatRepository_DeleteNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)

	// Try to delete non-existent chat
	err := repo.Delete(999999)
	// GORM doesn't return an error when deleting non-existent records
	// It just doesn't delete anything
	if err != nil {
		t.Errorf("Unexpected error when deleting non-existent chat: %v", err)
	}
}
