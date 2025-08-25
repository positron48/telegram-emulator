package emulator

import (
	"testing"

	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)



func TestBotManager_CreateBot(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Test creating a bot
	bot, err := botManager.CreateBot("TestBot", "testbot", "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11", "https://example.com/webhook")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	if bot.Name != "TestBot" {
		t.Errorf("Expected bot name 'TestBot', got '%s'", bot.Name)
	}

	if bot.Username != "testbot" {
		t.Errorf("Expected bot username 'testbot', got '%s'", bot.Username)
	}

	if bot.Token != "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" {
		t.Errorf("Expected bot token '123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11', got '%s'", bot.Token)
	}

	if !bot.IsActive {
		t.Error("Expected bot to be active")
	}

	if bot.ID == 0 {
		t.Error("Expected bot ID to be generated")
	}
}

func TestBotManager_GetBot(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Create a bot
	createdBot, err := botManager.CreateBot("TestBot", "testbot", "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11", "")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Get the bot
	retrievedBot, err := botManager.GetBot(createdBot.ID)
	if err != nil {
		t.Fatalf("Failed to get bot: %v", err)
	}

	if retrievedBot.ID != createdBot.ID {
		t.Errorf("Expected bot ID %d, got %d", createdBot.ID, retrievedBot.ID)
	}

	if retrievedBot.Name != "TestBot" {
		t.Errorf("Expected bot name 'TestBot', got '%s'", retrievedBot.Name)
	}
}

func TestBotManager_GetAllBots(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Create multiple bots
	_, err := botManager.CreateBot("Bot1", "bot1", "token1", "")
	if err != nil {
		t.Fatalf("Failed to create bot1: %v", err)
	}

	_, err = botManager.CreateBot("Bot2", "bot2", "token2", "")
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
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Create a bot
	bot, err := botManager.CreateBot("TestBot", "testbot", "token", "")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Update bot
	bot.Name = "UpdatedBot"
	bot.WebhookURL = "https://example.com/webhook"
	err = botManager.UpdateBot(bot)
	if err != nil {
		t.Fatalf("Failed to update bot: %v", err)
	}

	// Get updated bot
	updatedBot, err := botManager.GetBot(bot.ID)
	if err != nil {
		t.Fatalf("Failed to get updated bot: %v", err)
	}

	if updatedBot.Name != "UpdatedBot" {
		t.Errorf("Expected bot name 'UpdatedBot', got '%s'", updatedBot.Name)
	}

	if updatedBot.WebhookURL != "https://example.com/webhook" {
		t.Errorf("Expected webhook URL 'https://example.com/webhook', got '%s'", updatedBot.WebhookURL)
	}
}

func TestBotManager_DeleteBot(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Create a bot
	bot, err := botManager.CreateBot("TestBot", "testbot", "token", "")
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
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Try to get non-existent bot
	_, err := botManager.GetBot(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent bot")
	}
}

func TestBotManager_UpdateBotNotFound(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Try to update non-existent bot
	nonExistentBot := &models.Bot{
		ID:       999999,
		Name:     "NonExistent",
		Username: "nonexistent",
		Token:    "token",
	}

	err := botManager.UpdateBot(nonExistentBot)
	if err == nil {
		t.Error("Expected error when updating non-existent bot")
	}
}

func TestBotManager_DeleteBotNotFound(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Try to delete non-existent bot
	err := botManager.DeleteBot(999999)
	if err == nil {
		t.Error("Expected error when deleting non-existent bot")
	}
}

func TestBotManager_GetBotUpdates(t *testing.T) {
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Create a bot
	bot, err := botManager.CreateBot("TestBot", "testbot", "token", "")
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
	db := SetupTestDB(t)
	botRepo := repository.NewBotRepository(db)
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatRepo := repository.NewChatRepository(db)
	botManager := NewBotManager(botRepo, userRepo, messageRepo, chatRepo)

	// Create a bot
	bot, err := botManager.CreateBot("TestBot", "testbot", "token", "")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Create an update
	update := &models.Update{
		UpdateID: 1,
		Message: &models.Message{
			ID:      1,
			ChatID:  1,
			FromID:  1,
			Text:    "Test message",
			Type:    "text",
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
}
