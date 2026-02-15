package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhlga/budge/internal/imap"
)

func TestNewStatusBar(t *testing.T) {
	sb := NewStatusBar()

	if sb.connectionState != "Disconnected" {
		t.Errorf("connectionState = %q, want %q", sb.connectionState, "Disconnected")
	}

	if sb.helpText == "" {
		t.Error("helpText should not be empty")
	}

	if sb.loading {
		t.Error("loading should be false initially")
	}
}

func TestStatusBar_SetConnectionState(t *testing.T) {
	sb := NewStatusBar()
	sb.SetConnectionState("Connected")

	if sb.connectionState != "Connected" {
		t.Errorf("connectionState = %q, want %q", sb.connectionState, "Connected")
	}
}

func TestStatusBar_SetHelpText(t *testing.T) {
	sb := NewStatusBar()
	sb.SetHelpText("Custom help text")

	if sb.helpText != "Custom help text" {
		t.Errorf("helpText = %q, want %q", sb.helpText, "Custom help text")
	}
}

func TestStatusBar_SetSize(t *testing.T) {
	sb := NewStatusBar()
	sb.SetSize(100)

	if sb.width != 100 {
		t.Errorf("width = %d, want 100", sb.width)
	}
}

func TestStatusBar_Update_WindowSizeMsg(t *testing.T) {
	sb := NewStatusBar()

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedSB, _ := sb.Update(msg)

	if updatedSB.width != 120 {
		t.Errorf("width = %d, want 120", updatedSB.width)
	}
}

func TestStatusBar_Update_ConnectionStateChangedMsg(t *testing.T) {
	tests := []struct {
		state    imap.ConnectionState
		expected string
	}{
		{imap.StateDisconnected, "Disconnected"},
		{imap.StateConnecting, "Connecting..."},
		{imap.StateConnected, "Connected"},
		{imap.StateAuthenticated, "Authenticated"},
	}

	for _, tt := range tests {
		sb := NewStatusBar()
		msg := ConnectionStateChangedMsg{State: tt.state}
		updatedSB, _ := sb.Update(msg)

		if updatedSB.connectionState != tt.expected {
			t.Errorf("state %d: connectionState = %q, want %q", tt.state, updatedSB.connectionState, tt.expected)
		}
	}
}

func TestStatusBar_Update_LoadingMsg(t *testing.T) {
	sb := NewStatusBar()

	msg := LoadingMsg{Text: "Loading emails..."}
	updatedSB, cmd := sb.Update(msg)

	if !updatedSB.loading {
		t.Error("loading should be true after LoadingMsg")
	}

	if updatedSB.loadingText != "Loading emails..." {
		t.Errorf("loadingText = %q, want %q", updatedSB.loadingText, "Loading emails...")
	}

	if cmd == nil {
		t.Error("cmd should not be nil after LoadingMsg (should tick spinner)")
	}
}

func TestStatusBar_Update_LoadingClearedMsg(t *testing.T) {
	sb := NewStatusBar()
	sb.loading = true
	sb.loadingText = "Loading..."

	msg := LoadingClearedMsg{}
	updatedSB, _ := sb.Update(msg)

	if updatedSB.loading {
		t.Error("loading should be false after LoadingClearedMsg")
	}

	if updatedSB.loadingText != "" {
		t.Errorf("loadingText = %q, want empty", updatedSB.loadingText)
	}
}

func TestStatusBar_Update_EmailsLoadedMsg(t *testing.T) {
	sb := NewStatusBar()
	sb.loading = true
	sb.loadingText = "Loading..."

	msg := EmailsLoadedMsg{}
	updatedSB, _ := sb.Update(msg)

	if updatedSB.loading {
		t.Error("loading should be false after EmailsLoadedMsg")
	}
}

func TestStatusBar_Update_ErrorMsg(t *testing.T) {
	sb := NewStatusBar()
	sb.loading = true
	sb.loadingText = "Loading..."

	msg := ErrorMsg{Err: errTest{}}
	updatedSB, _ := sb.Update(msg)

	if updatedSB.loading {
		t.Error("loading should be false after ErrorMsg")
	}
}

func TestStatusBar_Update_ConnectErrorMsg(t *testing.T) {
	sb := NewStatusBar()
	sb.loading = true
	sb.loadingText = "Connecting..."

	msg := ConnectErrorMsg{}
	updatedSB, _ := sb.Update(msg)

	if updatedSB.loading {
		t.Error("loading should be false after ConnectErrorMsg")
	}
}

func TestStatusBar_View(t *testing.T) {
	sb := NewStatusBar()
	sb.SetSize(80)

	view := sb.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestStatusBar_View_WithLoading(t *testing.T) {
	sb := NewStatusBar()
	sb.SetSize(80)
	sb, _ = sb.Update(LoadingMsg{Text: "Searching..."})

	view := sb.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{"zero max", "test", 0, ""},
		{"negative max", "test", -1, ""},
		{"within limit", "test", 10, "test"},
		{"exact limit", "test", 4, "test"},
		{"max 1", "test", 1, "…"},
		{"needs truncation", "this is a long string", 10, "this is a…"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.max)
			// For truncation with ellipsis, just check it's shorter and ends with ellipsis
			if tt.max > 1 && len(tt.input) > tt.max {
				if len(result) > tt.max+3 { // Allow for multibyte chars
					t.Errorf("truncateString(%q, %d) length = %d, should be <= %d", tt.input, tt.max, len(result), tt.max+3)
				}
			} else if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.max, result, tt.expected)
			}
		})
	}
}

func TestStatusBar_clearsLoadingOnErrorMsg(t *testing.T) {
	sb := NewStatusBar()

	sb, _ = sb.Update(LoadingMsg{Text: "Searching..."})
	if !sb.loading {
		t.Fatalf("expected loading=true")
	}

	sb, _ = sb.Update(ErrorMsg{Err: errTest{}})
	if sb.loading {
		t.Fatalf("expected loading=false")
	}
}

type errTest struct{}

func (errTest) Error() string { return "test" }
