package models

import (
	"testing"
)

func TestChat_TableName(t *testing.T) {
	chat := &Chat{}
	if chat.TableName() != "chats" {
		t.Errorf("Expected table name 'chats', got '%s'", chat.TableName())
	}
}

func TestChatMember_TableName(t *testing.T) {
	member := &ChatMember{}
	if member.TableName() != "chat_members" {
		t.Errorf("Expected table name 'chat_members', got '%s'", member.TableName())
	}
}

func TestChat_IsPrivate(t *testing.T) {
	privateChat := &Chat{Type: "private"}
	if !privateChat.IsPrivate() {
		t.Error("Expected private chat to return true")
	}

	groupChat := &Chat{Type: "group"}
	if groupChat.IsPrivate() {
		t.Error("Expected group chat to return false")
	}
}

func TestChat_IsGroup(t *testing.T) {
	groupChat := &Chat{Type: "group"}
	if !groupChat.IsGroup() {
		t.Error("Expected group chat to return true")
	}

	privateChat := &Chat{Type: "private"}
	if privateChat.IsGroup() {
		t.Error("Expected private chat to return false")
	}
}

func TestChat_GetChatIcon(t *testing.T) {
	tests := []struct {
		name     string
		chatType string
		expected string
	}{
		{"private", "private", "👤"},
		{"group", "group", "👥"},
		{"supergroup", "supergroup", "💬"},
		{"channel", "channel", "💬"},
		{"unknown", "unknown", "💬"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{Type: tt.chatType}
			icon := chat.GetChatIcon()
			if icon != tt.expected {
				t.Errorf("Expected icon '%s' for type '%s', got '%s'", tt.expected, tt.chatType, icon)
			}
		})
	}
}

func TestChat_GetChatTypeLabel(t *testing.T) {
	tests := []struct {
		name     string
		chatType string
		expected string
	}{
		{"private", "private", "Приватный чат"},
		{"group", "group", "Группа"},
		{"supergroup", "supergroup", "Чат"},
		{"channel", "channel", "Чат"},
		{"unknown", "unknown", "Чат"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{Type: tt.chatType}
			label := chat.GetChatTypeLabel()
			if label != tt.expected {
				t.Errorf("Expected label '%s' for type '%s', got '%s'", tt.expected, tt.chatType, label)
			}
		})
	}
}

func TestChat_CanUserJoin(t *testing.T) {
	// Private chats cannot be joined
	privateChat := &Chat{Type: "private"}
	if privateChat.CanUserJoin() {
		t.Error("Expected private chat to not allow joining")
	}

	// Group chats can be joined
	groupChat := &Chat{Type: "group"}
	if !groupChat.CanUserJoin() {
		t.Error("Expected group chat to allow joining")
	}

	// Supergroups cannot be joined (not implemented)
	supergroupChat := &Chat{Type: "supergroup"}
	if supergroupChat.CanUserJoin() {
		t.Error("Expected supergroup chat to not allow joining")
	}
}

func TestChat_CanUserLeave(t *testing.T) {
	// Private chats cannot be left
	privateChat := &Chat{Type: "private"}
	if privateChat.CanUserLeave() {
		t.Error("Expected private chat to not allow leaving")
	}

	// Group chats can be left
	groupChat := &Chat{Type: "group"}
	if !groupChat.CanUserLeave() {
		t.Error("Expected group chat to allow leaving")
	}

	// Supergroups cannot be left (not implemented)
	supergroupChat := &Chat{Type: "supergroup"}
	if supergroupChat.CanUserLeave() {
		t.Error("Expected supergroup chat to not allow leaving")
	}
}

func TestChat_IsUserMember(t *testing.T) {
	chat := &Chat{
		ID: 1,
		Members: []User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
		},
	}

	// Test existing member
	if !chat.IsUserMember(1) {
		t.Error("Expected user 1 to be a member")
	}

	// Test non-existing member
	if chat.IsUserMember(999) {
		t.Error("Expected user 999 to not be a member")
	}

	// Test empty members list
	emptyChat := &Chat{ID: 2, Members: []User{}}
	if emptyChat.IsUserMember(1) {
		t.Error("Expected user 1 to not be a member of empty chat")
	}
}

func TestChat_AddMember(t *testing.T) {
	chat := &Chat{
		ID:      1,
		Members: []User{},
	}

	// Add first member
	user1 := User{ID: 1, Username: "user1"}
	chat.AddMember(user1)
	if len(chat.Members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(chat.Members))
	}
	if chat.Members[0].ID != 1 {
		t.Errorf("Expected user ID 1, got %d", chat.Members[0].ID)
	}

	// Add second member
	user2 := User{ID: 2, Username: "user2"}
	chat.AddMember(user2)
	if len(chat.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(chat.Members))
	}

	// Try to add duplicate member
	chat.AddMember(user1)
	if len(chat.Members) != 2 {
		t.Errorf("Expected 2 members after duplicate add, got %d", len(chat.Members))
	}
}

func TestChat_RemoveMember(t *testing.T) {
	chat := &Chat{
		ID: 1,
		Members: []User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
			{ID: 3, Username: "user3"},
		},
	}

	// Remove existing member
	chat.RemoveMember(2)
	if len(chat.Members) != 2 {
		t.Errorf("Expected 2 members after removal, got %d", len(chat.Members))
	}

	// Check that user 2 is not in the list
	for _, member := range chat.Members {
		if member.ID == 2 {
			t.Error("Expected user 2 to be removed")
		}
	}

	// Try to remove non-existing member
	chat.RemoveMember(999)
	if len(chat.Members) != 2 {
		t.Errorf("Expected 2 members after removing non-existing user, got %d", len(chat.Members))
	}

	// Remove all members
	chat.RemoveMember(1)
	chat.RemoveMember(3)
	if len(chat.Members) != 0 {
		t.Errorf("Expected 0 members after removing all, got %d", len(chat.Members))
	}
}
