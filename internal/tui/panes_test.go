package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhlga/budge/internal/config"
)

func TestTabTogglesFocus(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	if m.focus != focusRight {
		t.Fatalf("expected initial focusRight")
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.focus != focusLeft {
		t.Fatalf("expected focusLeft")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.focus != focusRight {
		t.Fatalf("expected focusRight")
	}
}

func TestHLSetFocus(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(Model)
	if m.focus != focusLeft {
		t.Fatalf("expected focusLeft")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = updated.(Model)
	if m.focus != focusRight {
		t.Fatalf("expected focusRight")
	}
}
