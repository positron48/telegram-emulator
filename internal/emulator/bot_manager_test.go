package emulator

import (
	"testing"
	"time"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBotManager_CreateBot(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	if bot.Name != "Test Bot" {
		t.Errorf("Expected name 'Test Bot', got '%s'", bot.Name)
	}

	if bot.Username != "testbot" {
		t.Errorf("Expected username 'testbot', got '%s'", bot.Username)
	}

	if bot.Token != "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz" {
		t.Errorf("Expected token '1234567890:ABCdefGHIjklMNOpqrsTUVwxyz', got '%s'", bot.Token)
	}

	if !bot.IsActive {
		t.Error("Expected IsActive to be true")
	}

	if bot.ID == 0 {
		t.Error("Expected bot ID to be generated")
	}
}

func TestBotManager_GetBot(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	createdBot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Get the bot by ID
	retrievedBot, err := botManager.GetBot(createdBot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot: %v", err)
	}

	if retrievedBot.ID != createdBot.ID {
		t.Errorf("Expected bot ID %d, got %d", createdBot.ID, retrievedBot.ID)
	}

	if retrievedBot.Username != createdBot.Username {
		t.Errorf("Expected username '%s', got '%s'", createdBot.Username, retrievedBot.Username)
	}
}

func TestBotManager_GetBotByToken(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	createdBot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Get the bot by token
	retrievedBot, err := botManager.GetBotByToken("1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to get bot by token: %v", err)
	}

	if retrievedBot.ID != createdBot.ID {
		t.Errorf("Expected bot ID %d, got %d", createdBot.ID, retrievedBot.ID)
	}

	if retrievedBot.Token != "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz" {
		t.Errorf("Expected token '1234567890:ABCdefGHIjklMNOpqrsTUVwxyz', got '%s'", retrievedBot.Token)
	}
}

func TestBotManager_GetAllBots(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create multiple bots
	_, err := botManager.CreateBot("Bot 1", "bot1", "1111111111:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot1: %v", err)
	}

	_, err = botManager.CreateBot("Bot 2", "bot2", "2222222222:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot2: %v", err)
	}

	// Get all bots
	bots, err := botManager.GetAllBots()
	if err != nil {
		t.Fatalf("Failed to get all bots: %v", err)
	}

	if len(bots) != 2 {
		t.Errorf("Expected 2 bots, got %d", len(bots))
	}
}

func TestBotManager_UpdateBot(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Update bot
	bot.Name = "Updated Bot"
	bot.WebhookURL = "https://updated.com/webhook"
	bot.IsActive = false

	updatedBot, err := botManager.UpdateBot(bot)
	if err != nil {
		t.Fatalf("Failed to update bot: %v", err)
	}

	if updatedBot.Name != "Updated Bot" {
		t.Errorf("Expected name 'Updated Bot', got '%s'", updatedBot.Name)
	}

	if updatedBot.WebhookURL != "https://updated.com/webhook" {
		t.Errorf("Expected webhook URL 'https://updated.com/webhook', got '%s'", updatedBot.WebhookURL)
	}

	if updatedBot.IsActive {
		t.Error("Expected IsActive to be false")
	}
}

func TestBotManager_DeleteBot(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Delete bot
	err = botManager.DeleteBot(bot.ID)
	if err != nil {
		t.Fatalf("Failed to delete bot: %v", err)
	}

	// Try to get the deleted bot
	_, err = botManager.GetBot(bot.ID)
	if err == nil {
		t.Error("Expected error when getting deleted bot")
	}
}

func TestBotManager_GetBotNotFound(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Try to get non-existent bot
	_, err := botManager.GetBot(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent bot")
	}
}

func TestBotManager_GetBotByTokenNotFound(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Try to get non-existent bot by token
	_, err := botManager.GetBotByToken("nonexistent:token")
	if err == nil {
		t.Error("Expected error when getting non-existent bot by token")
	}
}

func TestBotManager_UpdateBotNotFound(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

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

	_, err := botManager.UpdateBot(bot)
	if err == nil {
		t.Error("Expected error when updating non-existent bot")
	}
}

func TestBotManager_DeleteBotNotFound(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Try to delete non-existent bot
	err := botManager.DeleteBot(999999)
	if err == nil {
		t.Error("Expected error when deleting non-existent bot")
	}
}

func TestBotManager_GetActiveBots(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create bots with different active states
	_, err := botManager.CreateBot("Active Bot", "activebot", "1111111111:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create active bot: %v", err)
	}

	// Create an inactive bot
	inactiveBot, err := botManager.CreateBot("Inactive Bot", "inactivebot", "2222222222:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create inactive bot: %v", err)
	}

	// Deactivate the second bot
	inactiveBot.IsActive = false
	_, err = botManager.UpdateBot(inactiveBot)
	if err != nil {
		t.Fatalf("Failed to update bot: %v", err)
	}

	// Get active bots
	activeBots, err := botManager.GetActiveBots()
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

func TestBotManager_SetWebhook(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Set webhook
	webhookURL := "https://example.com/webhook"
	err = botManager.SetWebhook(bot.ID, webhookURL)
	if err != nil {
		t.Fatalf("Failed to set webhook: %v", err)
	}

	// Get the bot and check webhook
	retrievedBot, err := botManager.GetBot(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot: %v", err)
	}

	if retrievedBot.WebhookURL != webhookURL {
		t.Errorf("Expected webhook URL '%s', got '%s'", webhookURL, retrievedBot.WebhookURL)
	}
}

func TestBotManager_RemoveWebhook(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot with webhook
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Set webhook first
	err = botManager.SetWebhook(bot.ID, "https://example.com/webhook")
	if err != nil {
		t.Fatalf("Failed to set webhook: %v", err)
	}

	// Remove webhook
	err = botManager.RemoveWebhook(bot.ID)
	if err != nil {
		t.Fatalf("Failed to remove webhook: %v", err)
	}

	// Get the bot and check webhook is removed
	retrievedBot, err := botManager.GetBot(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot: %v", err)
	}

	if retrievedBot.WebhookURL != "" {
		t.Errorf("Expected empty webhook URL, got '%s'", retrievedBot.WebhookURL)
	}
}

func TestBotManager_GetBotUpdates(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Get updates (should be empty initially)
	updates, err := botManager.GetBotUpdates(bot.ID, 0, 10)
	if err != nil {
		t.Fatalf("Failed to get bot updates: %v", err)
	}

	if len(updates) != 0 {
		t.Errorf("Expected 0 updates, got %d", len(updates))
	}
}

func TestBotManager_AddUpdate(t *testing.T) {
	db := setupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	botManager := NewBotManager(botRepo, userRepo)

	// Create a bot
	bot, err := botManager.CreateBot("Test Bot", "testbot", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Create an update
	update := &models.Update{
		UpdateID: 1,
		Message: &models.Message{
			Text:      "Test message",
			Type:      "text",
			Status:    "sent",
			Timestamp: time.Now(),
		},
	}

	// Add update
	err = botManager.AddUpdate(bot.ID, update)
	if err != nil {
		t.Fatalf("Failed to add update: %v", err)
	}

	// Get updates
	updates, err := botManager.GetBotUpdates(bot.ID, 0, 10)
	if err != nil {
		t.Fatalf("Failed to get bot updates: %v", err)
	}

	if len(updates) != 1 {
		t.Errorf("Expected 1 update, got %d", len(updates))
	}

	if updates[0].UpdateID != 1 {
		t.Errorf("Expected update ID 1, got %d", updates[0].UpdateID)
	}
}
