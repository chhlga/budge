package tui

import (
	"testing"

	"github.com/chhlga/budge/internal/config"
)

func TestModel_autoSelectsInboxOnMailboxesLoaded(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "SENT", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	if m.currentMailbox != "" {
		t.Fatalf("expected no selected mailbox initially")
	}

	updated, cmd := m.Update(MailboxesLoadedMsg{Mailboxes: []string{"Archive", "INBOX", "Sent"}})
	m = updated.(Model)
	if cmd == nil {
		t.Fatalf("expected auto-select to return a command")
	}

	updated, _ = m.Update(cmd())
	m = updated.(Model)
	if m.currentMailbox != "INBOX" {
		t.Fatalf("expected currentMailbox=INBOX, got %q", m.currentMailbox)
	}
	if m.state != emailListView {
		t.Fatalf("expected state=emailListView, got %v", m.state)
	}
	if m.focus != focusRight {
		t.Fatalf("expected focusRight, got %v", m.focus)
	}
}
