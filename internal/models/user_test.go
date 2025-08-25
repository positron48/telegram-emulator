package models

import (
	"testing"
)

func TestUser_TableName(t *testing.T) {
	user := &User{}
	if user.TableName() != "users" {
		t.Errorf("Expected table name 'users', got '%s'", user.TableName())
	}
}

func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		expected  string
	}{
		{"both names", "John", "Doe", "John Doe"},
		{"first name only", "John", "", "John"},
		{"last name only", "", "Doe", " Doe"},
		{"no names", "", "", ""},
		{"with spaces", "  John  ", "  Doe  ", "  John     Doe  "},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}
			
			fullName := user.GetFullName()
			if fullName != tt.expected {
				t.Errorf("Expected full name '%s', got '%s'", tt.expected, fullName)
			}
		})
	}
}

func TestUser_SetOnline(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "testuser",
		IsOnline: false,
	}
	
	// Set user online
	user.SetOnline(true)
	
	if !user.IsOnline {
		t.Error("Expected user to be online")
	}
	
	// Set user offline
	user.SetOnline(false)
	
	if user.IsOnline {
		t.Error("Expected user to be offline")
	}
	
	// Check that LastSeen is updated when going offline
	if user.LastSeen.IsZero() {
		t.Error("Expected LastSeen to be updated when going offline")
	}
}

func TestUser_Validation(t *testing.T) {
	// Test user with valid data
	validUser := &User{
		ID:        1,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsOnline:  true,
		IsBot:     false,
	}
	
	// This should not panic
	_ = validUser.GetFullName()
	validUser.SetOnline(true)
	
	// Test user with minimal data
	minimalUser := &User{
		Username: "minimal",
	}
	
	// This should not panic
	_ = minimalUser.GetFullName()
	minimalUser.SetOnline(false)
	
	// Test bot user
	botUser := &User{
		ID:        2,
		Username:  "testbot",
		FirstName: "Test",
		LastName:  "Bot",
		IsBot:     true,
	}
	
	// This should not panic
	_ = botUser.GetFullName()
	botUser.SetOnline(true)
}
