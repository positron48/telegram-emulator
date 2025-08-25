package models

import (
	"testing"
)

func TestBot_TableName(t *testing.T) {
	bot := &Bot{}
	if bot.TableName() != "bots" {
		t.Errorf("Expected table name 'bots', got '%s'", bot.TableName())
	}
}

func TestBot_Activate(t *testing.T) {
	bot := &Bot{IsActive: false}
	bot.Activate()

	if !bot.IsActive {
		t.Error("Expected bot to be activated")
	}
}

func TestBot_Deactivate(t *testing.T) {
	bot := &Bot{IsActive: true}
	bot.Deactivate()

	if bot.IsActive {
		t.Error("Expected bot to be deactivated")
	}
}

func TestBot_SetWebhook(t *testing.T) {
	bot := &Bot{}
	webhookURL := "https://example.com/webhook"

	bot.SetWebhook(webhookURL)

	if bot.WebhookURL != webhookURL {
		t.Errorf("Expected webhook URL '%s', got '%s'", webhookURL, bot.WebhookURL)
	}
}

func TestBot_UpdateToken(t *testing.T) {
	bot := &Bot{Token: "old-token"}
	newToken := "new-token"

	bot.UpdateToken(newToken)

	if bot.Token != newToken {
		t.Errorf("Expected token '%s', got '%s'", newToken, bot.Token)
	}
}

func TestBot_Validation(t *testing.T) {
	// Test bot with valid data
	validBot := &Bot{
		ID:       1,
		Name:     "Test Bot",
		Username: "testbot",
		Token:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		IsActive: true,
	}

	// Test bot with minimal data
	minimalBot := &Bot{
		Username: "minimal",
	}

	// Test that methods work without panic
	validBot.Activate()
	validBot.Deactivate()
	validBot.SetWebhook("https://example.com/webhook")
	validBot.UpdateToken("new-token")

	minimalBot.Activate()
	minimalBot.Deactivate()
	minimalBot.SetWebhook("https://example.com/webhook")
	minimalBot.UpdateToken("new-token")
}
