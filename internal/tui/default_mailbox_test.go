package tui

import "testing"

func TestPickDefaultMailbox_prefersInbox(t *testing.T) {
	mailboxes := []string{"Archive", "INBOX", "Sent"}
	if got := pickDefaultMailbox(mailboxes, "Sent"); got != "INBOX" {
		t.Fatalf("expected INBOX, got %q", got)
	}
}

func TestPickDefaultMailbox_usesConfiguredDefaultWhenNoInbox(t *testing.T) {
	mailboxes := []string{"Archive", "Sent"}
	if got := pickDefaultMailbox(mailboxes, "Sent"); got != "Sent" {
		t.Fatalf("expected Sent, got %q", got)
	}
}
