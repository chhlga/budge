package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewSearch(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)

	if s.textInput.Placeholder != "Search emails..." {
		t.Errorf("Placeholder = %q, want %q", s.textInput.Placeholder, "Search emails...")
	}

	if s.textInput.CharLimit != 256 {
		t.Errorf("CharLimit = %d, want 256", s.textInput.CharLimit)
	}
}

func TestSearch_SetSize(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)

	s.SetSize(80, 24)

	if s.width != 80 {
		t.Errorf("width = %d, want 80", s.width)
	}

	if s.height != 24 {
		t.Errorf("height = %d, want 24", s.height)
	}

	if s.textInput.Width != 76 {
		t.Errorf("textInput.Width = %d, want 76", s.textInput.Width)
	}
}

func TestSearch_Init(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)

	cmd := s.Init()
	if cmd == nil {
		t.Error("Init() = nil, want non-nil blink command")
	}
}

func TestSearch_Update_EnterWithQuery(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)

	// Set value in textInput
	s.textInput.SetValue("test query")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := s.Update(msg)

	if cmd == nil {
		t.Fatal("Update() cmd = nil, want non-nil for search query")
	}

	result := cmd()
	if queryMsg, ok := result.(SearchQueryMsg); ok {
		if queryMsg.Query != "test query" {
			t.Errorf("SearchQueryMsg.Query = %q, want %q", queryMsg.Query, "test query")
		}
	} else {
		t.Errorf("cmd() returned %T, want SearchQueryMsg", result)
	}
}

func TestSearch_Update_EnterWithEmptyQuery(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := s.Update(msg)

	// Should not return a command for empty query
	if cmd != nil {
		result := cmd()
		if _, ok := result.(SearchQueryMsg); ok {
			t.Error("Update() should not return SearchQueryMsg for empty query")
		}
	}
}

func TestSearch_Update_Escape(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)

	s.textInput.SetValue("test")

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedS, cmd := s.Update(msg)

	if cmd == nil {
		t.Fatal("Update() cmd = nil, want non-nil for escape")
	}

	result := cmd()
	if _, ok := result.(SearchCancelledMsg); !ok {
		t.Errorf("cmd() returned %T, want SearchCancelledMsg", result)
	}

	// Verify textInput was reset
	if updatedS.textInput.Value() != "" {
		t.Errorf("textInput.Value() = %q, want empty after reset", updatedS.textInput.Value())
	}
}

func TestSearch_View(t *testing.T) {
	keys := NewKeyMap()
	s := NewSearch(keys)
	s.SetSize(80, 24)

	view := s.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}
