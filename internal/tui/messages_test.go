package tui

import (
	"testing"

	"github.com/chhlga/budge/internal/email"
)

func TestSearchQueryMsg(t *testing.T) {
	msg := SearchQueryMsg{Query: "test query"}

	if msg.Query != "test query" {
		t.Errorf("Expected Query 'test query', got '%s'", msg.Query)
	}
}

func TestLoadingMsg(t *testing.T) {
	msg := LoadingMsg{Text: "Loading..."}

	if msg.Text != "Loading..." {
		t.Errorf("Expected Text 'Loading...', got '%s'", msg.Text)
	}
}

func TestEmailsLoadedMsg(t *testing.T) {
	emails := []email.Message{
		{
			UID:     1001,
			Subject: "Test Subject",
			Flags:   []string{},
		},
		{
			UID:     1002,
			Subject: "Another Subject",
			Flags:   []string{},
		},
	}

	msg := EmailsLoadedMsg{Emails: emails, Total: 2}

	if len(msg.Emails) != 2 {
		t.Errorf("Expected 2 emails, got %d", len(msg.Emails))
	}

	if msg.Total != 2 {
		t.Errorf("Expected Total 2, got %d", msg.Total)
	}
}

func TestNewEmailMsg(t *testing.T) {
	msg := NewEmailMsg{Mailbox: "INBOX", Count: 5}

	if msg.Mailbox != "INBOX" {
		t.Errorf("Expected Mailbox 'INBOX', got '%s'", msg.Mailbox)
	}

	if msg.Count != 5 {
		t.Errorf("Expected Count 5, got %d", msg.Count)
	}
}

func TestStartIdleMonitoringMsg(t *testing.T) {
	msg := StartIdleMonitoringMsg{Mailbox: "INBOX"}

	if msg.Mailbox != "INBOX" {
		t.Errorf("Expected Mailbox 'INBOX', got '%s'", msg.Mailbox)
	}
}

func TestStopIdleMonitoringMsg(t *testing.T) {
	_ = StopIdleMonitoringMsg{}
}
