package emulator

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestChatManager_CreateChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create a chat
	chat, err := chatManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	if chat.Type != "private" {
		t.Errorf("Expected type 'private', got '%s'", chat.Type)
	}

	if chat.Title != "Test Chat" {
		t.Errorf("Expected title 'Test Chat', got '%s'", chat.Title)
	}

	if chat.Username != "testchat" {
		t.Errorf("Expected username 'testchat', got '%s'", chat.Username)
	}

	if chat.ID == 0 {
		t.Error("Expected chat ID to be generated")
	}
}

func TestChatManager_GetChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create a chat
	createdChat, err := chatManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Get the chat by ID
	retrievedChat, err := chatManager.GetChat(createdChat.ID)
	if err != nil {
		t.Fatalf("Failed to get chat: %v", err)
	}

	if retrievedChat.ID != createdChat.ID {
		t.Errorf("Expected chat ID %d, got %d", createdChat.ID, retrievedChat.ID)
	}

	if retrievedChat.Title != createdChat.Title {
		t.Errorf("Expected title '%s', got '%s'", createdChat.Title, retrievedChat.Title)
	}
}

func TestChatManager_GetAllChats(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create multiple chats
	_, err := chatManager.CreateChat("private", "Chat 1", "chat1")
	if err != nil {
		t.Fatalf("Failed to create chat1: %v", err)
	}

	_, err = chatManager.CreateChat("group", "Chat 2", "chat2")
	if err != nil {
		t.Fatalf("Failed to create chat2: %v", err)
	}

	// Get all chats
	chats, err := chatManager.GetAllChats()
	if err != nil {
		t.Fatalf("Failed to get all chats: %v", err)
	}

	if len(chats) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(chats))
	}
}

func TestChatManager_UpdateChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create a chat
	chat, err := chatManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Update chat
	chat.Title = "Updated Chat"
	chat.Type = "group"
	updatedChat, err := chatManager.UpdateChat(chat)
	if err != nil {
		t.Fatalf("Failed to update chat: %v", err)
	}

	if updatedChat.Title != "Updated Chat" {
		t.Errorf("Expected title 'Updated Chat', got '%s'", updatedChat.Title)
	}

	if updatedChat.Type != "group" {
		t.Errorf("Expected type 'group', got '%s'", updatedChat.Type)
	}
}

func TestChatManager_DeleteChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create a chat
	chat, err := chatManager.CreateChat("private", "Test Chat", "testchat")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Delete chat
	err = chatManager.DeleteChat(chat.ID)
	if err != nil {
		t.Fatalf("Failed to delete chat: %v", err)
	}

	// Try to get the deleted chat
	_, err = chatManager.GetChat(chat.ID)
	if err == nil {
		t.Error("Expected error when getting deleted chat")
	}
}

func TestChatManager_GetChatNotFound(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Try to get non-existent chat
	_, err := chatManager.GetChat(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent chat")
	}
}

func TestChatManager_UpdateChatNotFound(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Try to update non-existent chat
	chat := &models.Chat{
		ID:        999999,
		Type:      "private",
		Title:     "Non-existent Chat",
		Username:  "nonexistent",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := chatManager.UpdateChat(chat)
	if err == nil {
		t.Error("Expected error when updating non-existent chat")
	}
}

func TestChatManager_DeleteChatNotFound(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Try to delete non-existent chat
	err := chatManager.DeleteChat(999999)
	if err == nil {
		t.Error("Expected error when deleting non-existent chat")
	}
}

func TestChatManager_GetChatsByType(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create chats of different types
	_, err := chatManager.CreateChat("private", "Private Chat", "private")
	if err != nil {
		t.Fatalf("Failed to create private chat: %v", err)
	}

	_, err = chatManager.CreateChat("group", "Group Chat", "group")
	if err != nil {
		t.Fatalf("Failed to create group chat: %v", err)
	}

	// Get private chats
	privateChats, err := chatManager.GetChatsByType("private")
	if err != nil {
		t.Fatalf("Failed to get private chats: %v", err)
	}

	if len(privateChats) != 1 {
		t.Errorf("Expected 1 private chat, got %d", len(privateChats))
	}

	if privateChats[0].Type != "private" {
		t.Errorf("Expected chat type 'private', got '%s'", privateChats[0].Type)
	}

	// Get group chats
	groupChats, err := chatManager.GetChatsByType("group")
	if err != nil {
		t.Fatalf("Failed to get group chats: %v", err)
	}

	if len(groupChats) != 1 {
		t.Errorf("Expected 1 group chat, got %d", len(groupChats))
	}

	if groupChats[0].Type != "group" {
		t.Errorf("Expected chat type 'group', got '%s'", groupChats[0].Type)
	}
}

func TestChatManager_AddMemberToChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create a chat
	chat, err := chatManager.CreateChat("group", "Test Group", "testgroup")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Create a user
	user, err := chatManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Add user to chat
	err = chatManager.AddMemberToChat(chat.ID, user.ID, "member")
	if err != nil {
		t.Fatalf("Failed to add member to chat: %v", err)
	}

	// Get chat members
	members, err := chatManager.GetChatMembers(chat.ID)
	if err != nil {
		t.Fatalf("Failed to get chat members: %v", err)
	}

	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}

	if members[0].UserID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, members[0].UserID)
	}

	if members[0].Role != "member" {
		t.Errorf("Expected role 'member', got '%s'", members[0].Role)
	}
}

func TestChatManager_RemoveMemberFromChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatManager := NewChatManager(chatRepo, userRepo)

	// Create a chat
	chat, err := chatManager.CreateChat("group", "Test Group", "testgroup")
	if err != nil {
		t.Fatalf("Failed to create chat: %v", err)
	}

	// Create a user
	user, err := chatManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Add user to chat
	err = chatManager.AddMemberToChat(chat.ID, user.ID, "member")
	if err != nil {
		t.Fatalf("Failed to add member to chat: %v", err)
	}

	// Remove user from chat
	err = chatManager.RemoveMemberFromChat(chat.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to remove member from chat: %v", err)
	}

	// Get chat members
	members, err := chatManager.GetChatMembers(chat.ID)
	if err != nil {
		t.Fatalf("Failed to get chat members: %v", err)
	}

	if len(members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(members))
	}
}
