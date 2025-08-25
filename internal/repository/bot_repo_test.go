package repository

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"
)

func TestBotRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	if bot.ID == 0 {
		t.Error("Expected bot ID to be set after creation")
	}
}

func TestBotRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create a bot
	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Get the bot by ID
	retrievedBot, err := repo.GetByID(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot by ID: %v", err)
	}

	if retrievedBot.ID != bot.ID {
		t.Errorf("Expected bot ID %d, got %d", bot.ID, retrievedBot.ID)
	}

	if retrievedBot.Username != bot.Username {
		t.Errorf("Expected username '%s', got '%s'", bot.Username, retrievedBot.Username)
	}
}

func TestBotRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create a bot
	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Get the bot by username
	retrievedBot, err := repo.GetByUsername("testbot")
	if err != nil {
		t.Fatalf("Failed to get bot by username: %v", err)
	}

	if retrievedBot.ID != bot.ID {
		t.Errorf("Expected bot ID %d, got %d", bot.ID, retrievedBot.ID)
	}

	if retrievedBot.Username != "testbot" {
		t.Errorf("Expected username 'testbot', got '%s'", retrievedBot.Username)
	}
}

func TestBotRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create multiple bots
	bot1 := &models.Bot{
		Name:       "Bot 1",
		Username:   "bot1",
		Token:      "1111111111:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example1.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	bot2 := &models.Bot{
		Name:       "Bot 2",
		Username:   "bot2",
		Token:      "2222222222:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example2.com/webhook",
		IsActive:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot1)
	if err != nil {
		t.Fatalf("Failed to create bot1: %v", err)
	}

	err = repo.Create(bot2)
	if err != nil {
		t.Fatalf("Failed to create bot2: %v", err)
	}

	// Get all bots
	bots, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all bots: %v", err)
	}

	if len(bots) != 2 {
		t.Errorf("Expected 2 bots, got %d", len(bots))
	}
}

func TestBotRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create a bot
	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Update bot
	bot.Name = "Updated Bot"
	bot.WebhookURL = "https://updated.com/webhook"
	bot.IsActive = false

	err = repo.Update(bot)
	if err != nil {
		t.Fatalf("Failed to update bot: %v", err)
	}

	// Get the updated bot
	retrievedBot, err := repo.GetByID(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get updated bot: %v", err)
	}

	if retrievedBot.Name != "Updated Bot" {
		t.Errorf("Expected name 'Updated Bot', got '%s'", retrievedBot.Name)
	}

	if retrievedBot.WebhookURL != "https://updated.com/webhook" {
		t.Errorf("Expected webhook URL 'https://updated.com/webhook', got '%s'", retrievedBot.WebhookURL)
	}

	if retrievedBot.IsActive {
		t.Error("Expected IsActive to be false")
	}
}

func TestBotRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create a bot
	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Delete bot
	err = repo.Delete(bot.ID)
	if err != nil {
		t.Fatalf("Failed to delete bot: %v", err)
	}

	// Try to get the deleted bot
	_, err = repo.GetByID(bot.ID)
	if err == nil {
		t.Error("Expected error when getting deleted bot")
	}
}

func TestBotRepository_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Try to get non-existent bot
	_, err := repo.GetByID(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent bot")
	}
}

func TestBotRepository_GetByUsernameNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Try to get non-existent bot by username
	_, err := repo.GetByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent bot by username")
	}
}

func TestBotRepository_UpdateNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Try to update non-existent bot
	bot := &models.Bot{
		ID:         999999,
		Name:       "Non-existent Bot",
		Username:   "nonexistent",
		Token:      "9999999999:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://nonexistent.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Update(bot)
	// GORM doesn't return an error when updating non-existent records
	// It just doesn't update anything
	if err != nil {
		t.Errorf("Unexpected error when updating non-existent bot: %v", err)
	}
}

func TestBotRepository_DeleteNonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Try to delete non-existent bot
	err := repo.Delete(999999)
	// GORM doesn't return an error when deleting non-existent records
	// It just doesn't delete anything
	if err != nil {
		t.Errorf("Unexpected error when deleting non-existent bot: %v", err)
	}
}

func TestBotRepository_GetActive(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create bots with different active states
	activeBot := &models.Bot{
		Name:       "Active Bot",
		Username:   "activebot",
		Token:      "1111111111:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://active.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	inactiveBot := &models.Bot{
		Name:       "Inactive Bot",
		Username:   "inactivebot",
		Token:      "2222222222:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://inactive.com/webhook",
		IsActive:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(activeBot)
	if err != nil {
		t.Fatalf("Failed to create active bot: %v", err)
	}

	err = repo.Create(inactiveBot)
	if err != nil {
		t.Fatalf("Failed to create inactive bot: %v", err)
	}

	// Get active bots
	activeBots, err := repo.GetActive()
	if err != nil {
		t.Fatalf("Failed to get active bots: %v", err)
	}

	if len(activeBots) != 1 {
		t.Errorf("Expected 1 active bot, got %d", len(activeBots))
	}

	if !activeBots[0].IsActive {
		t.Error("Expected bot to be active")
	}

	if activeBots[0].Username != "activebot" {
		t.Errorf("Expected username 'activebot', got '%s'", activeBots[0].Username)
	}
}

func TestBotRepository_SetWebhookURL(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create a bot
	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Set webhook URL
	webhookURL := "https://example.com/webhook"
	err = repo.SetWebhookURL(bot.ID, webhookURL)
	if err != nil {
		t.Fatalf("Failed to set webhook URL: %v", err)
	}

	// Get the bot and check webhook URL
	retrievedBot, err := repo.GetByID(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot: %v", err)
	}

	if retrievedBot.WebhookURL != webhookURL {
		t.Errorf("Expected webhook URL '%s', got '%s'", webhookURL, retrievedBot.WebhookURL)
	}
}

func TestBotRepository_SetActiveStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBotRepository(db)

	// Create a bot
	bot := &models.Bot{
		Name:       "Test Bot",
		Username:   "testbot",
		Token:      "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
		WebhookURL: "https://example.com/webhook",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(bot)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Set inactive status
	err = repo.SetActiveStatus(bot.ID, false)
	if err != nil {
		t.Fatalf("Failed to set active status: %v", err)
	}

	// Get the bot and check status
	retrievedBot, err := repo.GetByID(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot: %v", err)
	}

	if retrievedBot.IsActive {
		t.Error("Expected IsActive to be false")
	}
}
