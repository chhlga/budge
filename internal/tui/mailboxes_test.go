package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMailboxItem(t *testing.T) {
	item := mailboxItem{name: "INBOX"}

	if item.Title() != "INBOX" {
		t.Errorf("Title() = %q, want %q", item.Title(), "INBOX")
	}

	if item.Description() != "" {
		t.Errorf("Description() = %q, want empty string", item.Description())
	}

	if item.FilterValue() != "" {
		t.Errorf("FilterValue() = %q, want empty string", item.FilterValue())
	}
}

func TestMailboxDelegate(t *testing.T) {
	d := mailboxDelegate{}

	if d.Height() != 1 {
		t.Errorf("Height() = %d, want 1", d.Height())
	}

	if d.Spacing() != 0 {
		t.Errorf("Spacing() = %d, want 0", d.Spacing())
	}

	cmd := d.Update(nil, nil)
	if cmd != nil {
		t.Errorf("Update() = %v, want nil", cmd)
	}
}

func TestNewMailboxList(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)

	if ml.keys.Enter.Keys()[0] != keys.Enter.Keys()[0] {
		t.Errorf("keys not set correctly")
	}

	if ml.list.Title != "Mailboxes" {
		t.Errorf("Title = %q, want %q", ml.list.Title, "Mailboxes")
	}
}

func TestMailboxList_SetSize(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)

	ml.SetSize(80, 24)
	// Size is set internally, no direct way to verify but ensures no panic
}

func TestMailboxList_SetMailboxes(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)

	mailboxes := []string{"INBOX", "Sent", "Drafts"}
	ml.SetMailboxes(mailboxes)

	if len(ml.list.Items()) != 3 {
		t.Errorf("SetMailboxes() set %d items, want 3", len(ml.list.Items()))
	}

	// Verify first item
	if item, ok := ml.list.Items()[0].(mailboxItem); ok {
		if item.name != "INBOX" {
			t.Errorf("First item name = %q, want %q", item.name, "INBOX")
		}
	} else {
		t.Error("First item is not mailboxItem")
	}
}

func TestMailboxList_Init(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)

	cmd := ml.Init()
	if cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestMailboxList_Update_MailboxSelected(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)
	ml.SetMailboxes([]string{"INBOX", "Sent"})

	// Simulate Enter key press
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedML, cmd := ml.Update(msg)

	if cmd == nil {
		t.Error("Update() cmd = nil, want non-nil for mailbox selection")
	}

	// Execute the command to get the message
	if cmd != nil {
		result := cmd()
		if selectedMsg, ok := result.(MailboxSelectedMsg); ok {
			if selectedMsg.Mailbox != "INBOX" {
				t.Errorf("MailboxSelectedMsg.Mailbox = %q, want %q", selectedMsg.Mailbox, "INBOX")
			}
		} else {
			t.Errorf("cmd() returned %T, want MailboxSelectedMsg", result)
		}
	}

	_ = updatedML
}

func TestMailboxList_Update_SkipSeparator(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)
	ml.SetMailboxes([]string{"INBOX", "---", "Sent"})

	// Move to separator
	ml.list.Select(1)

	// Simulate Enter key press on separator
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := ml.Update(msg)

	// Should not return a selection command for separator
	if cmd != nil {
		result := cmd()
		if _, ok := result.(MailboxSelectedMsg); ok {
			t.Error("Update() should not return MailboxSelectedMsg for separator")
		}
	}
}

func TestMailboxList_Update_MailboxesLoadedMsg(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)

	msg := MailboxesLoadedMsg{Mailboxes: []string{"INBOX", "Drafts", "Sent"}}
	updatedML, _ := ml.Update(msg)

	if len(updatedML.list.Items()) != 3 {
		t.Errorf("After MailboxesLoadedMsg, got %d items, want 3", len(updatedML.list.Items()))
	}
}

func TestMailboxList_View(t *testing.T) {
	keys := NewKeyMap()
	ml := NewMailboxList(keys)
	ml.SetMailboxes([]string{"INBOX"})

	view := ml.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}
