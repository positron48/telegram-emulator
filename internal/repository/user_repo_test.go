package repository

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"

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

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set after creation")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a user
	user := &models.User{
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by ID
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username '%s', got '%s'", user.Username, retrievedUser.Username)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a user
	user := &models.User{
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by username
	retrievedUser, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", retrievedUser.Username)
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create multiple users
	user1 := &models.User{
		Username:  "user1",
		FirstName: "User",
		LastName:  "One",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user2 := &models.User{
		Username:  "user2",
		FirstName: "User",
		LastName:  "Two",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(user1)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	err = repo.Create(user2)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// Get all users
	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a user
	user := &models.User{
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.FirstName = "Updated"
	user.LastName = "Name"
	user.IsOnline = true

	err = repo.Update(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Get the updated user
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if retrievedUser.FirstName != "Updated" {
		t.Errorf("Expected first name 'Updated', got '%s'", retrievedUser.FirstName)
	}

	if retrievedUser.LastName != "Name" {
		t.Errorf("Expected last name 'Name', got '%s'", retrievedUser.LastName)
	}

	if !retrievedUser.IsOnline {
		t.Error("Expected IsOnline to be true")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a user
	user := &models.User{
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Delete user
	err = repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Try to get the deleted user
	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}

func TestUserRepository_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Try to get non-existent user
	_, err := repo.GetByID(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestUserRepository_GetByUsernameNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Try to get non-existent user by username
	_, err := repo.GetByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent user by username")
	}
}

func TestUserRepository_UpdateNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Try to update non-existent user
	user := &models.User{
		ID:        999999,
		Username:  "nonexistent",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
		IsOnline:  false,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Update(user)
	// GORM doesn't return an error when updating non-existent records
	// It just doesn't update anything
	if err != nil {
		t.Errorf("Unexpected error when updating non-existent user: %v", err)
	}
}

func TestUserRepository_DeleteNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Try to delete non-existent user
	err := repo.Delete(999999)
	// GORM doesn't return an error when deleting non-existent records
	// It just doesn't delete anything
	if err != nil {
		t.Errorf("Unexpected error when deleting non-existent user: %v", err)
	}
}
