package emulator

import (
	"testing"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{}, &models.Bot{}, &models.ChatMember{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserManager_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Test creating a regular user
	user, err := userManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.FirstName != "Test" {
		t.Errorf("Expected first name 'Test', got '%s'", user.FirstName)
	}

	if user.LastName != "User" {
		t.Errorf("Expected last name 'User', got '%s'", user.LastName)
	}

	if user.IsBot {
		t.Error("Expected IsBot to be false")
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be generated")
	}
}

func TestUserManager_CreateBot(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Test creating a bot
	bot, err := userManager.CreateUser("testbot", "Test", "Bot", true)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	if !bot.IsBot {
		t.Error("Expected IsBot to be true")
	}

	if bot.Username != "testbot" {
		t.Errorf("Expected username 'testbot', got '%s'", bot.Username)
	}
}

func TestUserManager_GetUser(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Create a user
	createdUser, err := userManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by ID
	retrievedUser, err := userManager.GetUser(createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %d, got %d", createdUser.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != createdUser.Username {
		t.Errorf("Expected username '%s', got '%s'", createdUser.Username, retrievedUser.Username)
	}
}

func TestUserManager_GetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Create multiple users
	_, err := userManager.CreateUser("user1", "User", "One", false)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	_, err = userManager.CreateUser("user2", "User", "Two", false)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	_, err = userManager.CreateUser("bot1", "Bot", "One", true)
	if err != nil {
		t.Fatalf("Failed to create bot1: %v", err)
	}

	// Get all users
	users, err := userManager.GetAllUsers()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}

	// Check that we have 2 regular users and 1 bot
	regularUsers := 0
	bots := 0
	for _, user := range users {
		if user.IsBot {
			bots++
		} else {
			regularUsers++
		}
	}

	if regularUsers != 2 {
		t.Errorf("Expected 2 regular users, got %d", regularUsers)
	}

	if bots != 1 {
		t.Errorf("Expected 1 bot, got %d", bots)
	}
}

func TestUserManager_GetUserByUsername(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Create a user
	createdUser, err := userManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by username
	retrievedUser, err := userManager.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %d, got %d", createdUser.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", retrievedUser.Username)
	}
}

func TestUserManager_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Create a user
	user, err := userManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.FirstName = "Updated"
	user.LastName = "Name"
	updatedUser, err := userManager.UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if updatedUser.FirstName != "Updated" {
		t.Errorf("Expected first name 'Updated', got '%s'", updatedUser.FirstName)
	}

	if updatedUser.LastName != "Name" {
		t.Errorf("Expected last name 'Name', got '%s'", updatedUser.LastName)
	}
}

func TestUserManager_DeleteUser(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Create a user
	user, err := userManager.CreateUser("testuser", "Test", "User", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Delete user
	err = userManager.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Try to get the deleted user
	_, err = userManager.GetUser(user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}

func TestUserManager_GetUserNotFound(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Try to get non-existent user
	_, err := userManager.GetUser(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestUserManager_GetUserByUsernameNotFound(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := NewUserManager(userRepo, botRepo)

	// Try to get non-existent user by username
	_, err := userManager.GetUserByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent user by username")
	}
}
