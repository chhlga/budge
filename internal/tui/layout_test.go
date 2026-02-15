package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/chhlga/budge/internal/config"
)

func TestView_isExactlyModelHeight(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.width = 100
	m.height = 20
	m.currentMailbox = "INBOX"

	view := m.View()
	if got := lipgloss.Height(view); got != m.height {
		t.Fatalf("expected view height to be %d lines, got %d", m.height, got)
	}
}

func TestView_isExactlyModelHeight_narrowTerminal(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.width = 30
	m.height = 12
	m.currentMailbox = "INBOX"

	view := m.View()
	if got := lipgloss.Height(view); got != m.height {
		t.Fatalf("expected view height to be %d lines, got %d", m.height, got)
	}
}

func TestView_hasBreadcrumbHeader(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.width = 100
	m.height = 20
	m.currentMailbox = "INBOX"

	view := m.View()
	firstLine := strings.SplitN(strings.TrimRight(view, "\n"), "\n", 2)[0]
	if !strings.Contains(firstLine, "INBOX") {
		t.Fatalf("expected header to contain mailbox name, got: %q", firstLine)
	}
	if got := lipgloss.Width(firstLine); got == 0 {
		t.Fatalf("expected header to be non-empty")
	}
}
