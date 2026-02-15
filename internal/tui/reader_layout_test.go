package tui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/chhlga/budge/internal/email"
)

func TestEmailReaderView_isExactlyAssignedHeight(t *testing.T) {
	r := NewEmailReader(KeyMap{})
	r.SetSize(60, 18)

	r.SetEmail(email.Message{Subject: "hello"})
	r.SetBody("line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\n")

	view := r.View()
	if got := lipgloss.Height(view); got != 18 {
		t.Fatalf("expected reader view height to be %d, got %d", 18, got)
	}
}
