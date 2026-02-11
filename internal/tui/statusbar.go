package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusBar displays connection status and help text
type StatusBar struct {
	connectionState string
	helpText        string
	width           int
}

// NewStatusBar creates a new status bar
func NewStatusBar() StatusBar {
	return StatusBar{
		connectionState: "Disconnected",
		helpText:        "q: quit | r: refresh | ?: help",
	}
}

// SetConnectionState updates the connection status
func (s *StatusBar) SetConnectionState(state string) {
	s.connectionState = state
}

// SetHelpText updates the help text
func (s *StatusBar) SetHelpText(text string) {
	s.helpText = text
}

// SetSize updates the status bar width
func (s *StatusBar) SetSize(width int) {
	s.width = width
}

// Update handles messages for the status bar
func (s StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
	case ConnectionStateChangedMsg:
		switch msg.State {
		case 0: // StateDisconnected
			s.connectionState = "Disconnected"
		case 1: // StateConnecting
			s.connectionState = "Connecting..."
		case 2: // StateConnected
			s.connectionState = "Connected"
		case 3: // StateAuthenticated
			s.connectionState = "Authenticated"
		}
	}
	return s, nil
}

// View renders the status bar
func (s StatusBar) View() string {
	left := fmt.Sprintf("📡 %s", s.connectionState)
	right := s.helpText

	leftStyle := StatusBarStyle.Copy().Width(s.width / 2)
	rightStyle := StatusBarStyle.Copy().Width(s.width / 2).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
}
