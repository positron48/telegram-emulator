package models

import (
	"testing"
)

func TestReplyKeyboardMarkup_IsReplyKeyboardMarkup(t *testing.T) {
	markup := &ReplyKeyboardMarkup{}
	if !markup.IsReplyKeyboardMarkup() {
		t.Error("Expected ReplyKeyboardMarkup to return true")
	}
}

func TestReplyKeyboardRemove_IsReplyKeyboardRemove(t *testing.T) {
	remove := &ReplyKeyboardRemove{}
	if !remove.IsReplyKeyboardRemove() {
		t.Error("Expected ReplyKeyboardRemove to return true")
	}
}

func TestInlineKeyboardMarkup_IsInlineKeyboardMarkup(t *testing.T) {
	markup := &InlineKeyboardMarkup{}
	if !markup.IsInlineKeyboardMarkup() {
		t.Error("Expected InlineKeyboardMarkup to return true")
	}
}

func TestKeyboardTypes(t *testing.T) {
	// Test ReplyKeyboardMarkup
	replyMarkup := &ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Button 1"},
				{Text: "Button 2"},
			},
			{
				{Text: "Button 3"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
		Selective:       true,
	}

	if !replyMarkup.IsReplyKeyboardMarkup() {
		t.Error("Expected ReplyKeyboardMarkup to be identified correctly")
	}

	// Test ReplyKeyboardRemove
	removeMarkup := &ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	if !removeMarkup.IsReplyKeyboardRemove() {
		t.Error("Expected ReplyKeyboardRemove to be identified correctly")
	}

	// Test InlineKeyboardMarkup
	inlineMarkup := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "Inline Button 1", CallbackData: "data1"},
				{Text: "Inline Button 2", CallbackData: "data2"},
			},
		},
	}

	if !inlineMarkup.IsInlineKeyboardMarkup() {
		t.Error("Expected InlineKeyboardMarkup to be identified correctly")
	}
}

func TestKeyboardButtonValidation(t *testing.T) {
	// Test regular keyboard button
	button := KeyboardButton{
		Text: "Test Button",
	}

	if button.Text != "Test Button" {
		t.Errorf("Expected button text 'Test Button', got '%s'", button.Text)
	}

	// Test inline keyboard button
	inlineButton := InlineKeyboardButton{
		Text:         "Inline Test",
		CallbackData: "test_callback",
	}

	if inlineButton.Text != "Inline Test" {
		t.Errorf("Expected inline button text 'Inline Test', got '%s'", inlineButton.Text)
	}

	if inlineButton.CallbackData != "test_callback" {
		t.Errorf("Expected callback data 'test_callback', got '%s'", inlineButton.CallbackData)
	}
}

func TestKeyboardMarkupProperties(t *testing.T) {
	// Test ReplyKeyboardMarkup properties
	replyMarkup := &ReplyKeyboardMarkup{
		Keyboard:        [][]KeyboardButton{},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       false,
	}

	if !replyMarkup.ResizeKeyboard {
		t.Error("Expected ResizeKeyboard to be true")
	}

	if !replyMarkup.OneTimeKeyboard {
		t.Error("Expected OneTimeKeyboard to be true")
	}

	if replyMarkup.Selective {
		t.Error("Expected Selective to be false")
	}

	// Test ReplyKeyboardRemove properties
	removeMarkup := &ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      true,
	}

	if !removeMarkup.RemoveKeyboard {
		t.Error("Expected RemoveKeyboard to be true")
	}

	if !removeMarkup.Selective {
		t.Error("Expected Selective to be true")
	}
}
