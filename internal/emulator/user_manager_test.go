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
	userManager := NewUserManager(userRepo)

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

	if user.ID == "" {
		t.Error("Expected user ID to be generated")
	}
}

func TestUserManager_CreateBot(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userManager := NewUserManager(userRepo)

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
	userManager := NewUserManager(userRepo)

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
		t.Errorf("Expected user ID '%s', got '%s'", createdUser.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != createdUser.Username {
		t.Errorf("Expected username '%s', got '%s'", createdUser.Username, retrievedUser.Username)
	}
}

func TestUserManager_GetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userManager := NewUserManager(userRepo)

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
