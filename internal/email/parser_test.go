package email

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParse_PlainText(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "plain.eml"))
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Check headers
	if len(msg.From) != 1 || msg.From[0].Email != "sender@example.com" {
		t.Errorf("Expected From 'sender@example.com', got %v", msg.From)
	}

	if len(msg.To) != 1 || msg.To[0].Email != "recipient@example.com" {
		t.Errorf("Expected To 'recipient@example.com', got %v", msg.To)
	}

	if msg.Subject != "Plain Text Email" {
		t.Errorf("Expected subject 'Plain Text Email', got '%s'", msg.Subject)
	}

	if msg.MessageID != "<plain123@example.com>" {
		t.Errorf("Expected message ID '<plain123@example.com>', got '%s'", msg.MessageID)
	}

	expectedDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	if !msg.Date.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, msg.Date)
	}

	// Check body
	if msg.Body == nil {
		t.Fatal("Body is nil")
	}

	expectedText := "This is a plain text email.\nIt has multiple lines.\nAnd should be parsed correctly.\n"
	if msg.Body.Text != expectedText {
		t.Errorf("Expected text:\n%s\nGot:\n%s", expectedText, msg.Body.Text)
	}

	if msg.Body.HTML != "" {
		t.Errorf("Expected empty HTML, got '%s'", msg.Body.HTML)
	}
}

func TestParse_HTML(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "html.eml"))
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Check headers
	if msg.Subject != "HTML Email" {
		t.Errorf("Expected subject 'HTML Email', got '%s'", msg.Subject)
	}

	// Check body
	if msg.Body == nil {
		t.Fatal("Body is nil")
	}

	if msg.Body.HTML == "" {
		t.Error("Expected HTML content, got empty string")
	}

	// HTML should contain key elements
	if msg.Body.Text != "" {
		t.Errorf("Expected empty text for HTML-only email, got '%s'", msg.Body.Text)
	}
}

func TestParse_Multipart(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "multipart.eml"))
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Check headers
	if len(msg.From) != 1 || msg.From[0].Name != "Alice" || msg.From[0].Email != "alice@example.com" {
		t.Errorf("Expected From 'Alice <alice@example.com>', got %v", msg.From)
	}

	if len(msg.To) != 1 || msg.To[0].Name != "Bob" || msg.To[0].Email != "bob@example.com" {
		t.Errorf("Expected To 'Bob <bob@example.com>', got %v", msg.To)
	}

	if len(msg.Cc) != 1 || msg.Cc[0].Email != "charlie@example.com" {
		t.Errorf("Expected Cc 'charlie@example.com', got %v", msg.Cc)
	}

	if msg.Subject != "Multipart Alternative Email" {
		t.Errorf("Expected subject 'Multipart Alternative Email', got '%s'", msg.Subject)
	}

	if msg.InReplyTo != "<previous@example.com>" {
		t.Errorf("Expected InReplyTo '<previous@example.com>', got '%s'", msg.InReplyTo)
	}

	if len(msg.References) != 2 {
		t.Errorf("Expected 2 references, got %d", len(msg.References))
	}

	// Check body - should have both text and HTML
	if msg.Body == nil {
		t.Fatal("Body is nil")
	}

	if msg.Body.Text == "" {
		t.Error("Expected plain text content, got empty string")
	}

	if msg.Body.HTML == "" {
		t.Error("Expected HTML content, got empty string")
	}
}

func TestParse_InvalidEmail(t *testing.T) {
	invalidData := []byte("This is not a valid email")

	msg, err := Parse(invalidData)
	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
	if msg != nil {
		t.Error("Expected nil message for invalid email")
	}
}

func TestAddress_String(t *testing.T) {
	tests := []struct {
		addr     Address
		expected string
	}{
		{
			addr:     Address{Name: "John Doe", Email: "john@example.com"},
			expected: "John Doe <john@example.com>",
		},
		{
			addr:     Address{Email: "john@example.com"},
			expected: "john@example.com",
		},
		{
			addr:     Address{Name: "", Email: "test@test.com"},
			expected: "test@test.com",
		},
	}

	for _, tt := range tests {
		result := tt.addr.String()
		if result != tt.expected {
			t.Errorf("Address.String() = '%s', want '%s'", result, tt.expected)
		}
	}
}

func TestMessage_HasFlag(t *testing.T) {
	msg := &Message{
		Flags: []string{"\\Seen", "\\Flagged"},
	}

	if !msg.HasFlag("\\Seen") {
		t.Error("Expected message to have \\Seen flag")
	}

	if !msg.HasFlag("\\Flagged") {
		t.Error("Expected message to have \\Flagged flag")
	}

	if msg.HasFlag("\\Draft") {
		t.Error("Expected message not to have \\Draft flag")
	}
}

func TestMessage_IsUnread(t *testing.T) {
	tests := []struct {
		flags  []string
		unread bool
	}{
		{[]string{}, true},
		{[]string{"\\Flagged"}, true},
		{[]string{"\\Seen"}, false},
		{[]string{"\\Seen", "\\Flagged"}, false},
	}

	for _, tt := range tests {
		msg := &Message{Flags: tt.flags}
		if msg.IsUnread() != tt.unread {
			t.Errorf("IsUnread() with flags %v = %v, want %v", tt.flags, msg.IsUnread(), tt.unread)
		}
	}
}

func TestMessage_IsFlagged(t *testing.T) {
	tests := []struct {
		flags   []string
		flagged bool
	}{
		{[]string{}, false},
		{[]string{"\\Seen"}, false},
		{[]string{"\\Flagged"}, true},
		{[]string{"\\Seen", "\\Flagged"}, true},
	}

	for _, tt := range tests {
		msg := &Message{Flags: tt.flags}
		if msg.IsFlagged() != tt.flagged {
			t.Errorf("IsFlagged() with flags %v = %v, want %v", tt.flags, msg.IsFlagged(), tt.flagged)
		}
	}
}
