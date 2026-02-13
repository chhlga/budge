package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	connectionState string
	helpText        string
	loading         bool
	loadingText     string
	spinner         spinner.Model
	width           int
}

func NewStatusBar() StatusBar {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return StatusBar{
		connectionState: "Disconnected",
		helpText:        "q: quit | r: refresh | ?: help",
		spinner:         sp,
	}
}

func (s *StatusBar) SetConnectionState(state string) {
	s.connectionState = state
}

func (s *StatusBar) SetHelpText(text string) {
	s.helpText = text
}

func (s *StatusBar) SetSize(width int) {
	s.width = width
}

func (s StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
	case ConnectionStateChangedMsg:
		switch msg.State {
		case 0:
			s.connectionState = "Disconnected"
		case 1:
			s.connectionState = "Connecting..."
		case 2:
			s.connectionState = "Connected"
		case 3:
			s.connectionState = "Authenticated"
		}
	case LoadingMsg:
		s.loading = true
		s.loadingText = msg.Text
		cmd = s.spinner.Tick
	case LoadingClearedMsg:
		s.loading = false
		s.loadingText = ""
	case EmailsLoadedMsg:
		s.loading = false
		s.loadingText = ""
	case ErrorMsg:
		s.loading = false
		s.loadingText = ""
	case ConnectErrorMsg:
		s.loading = false
		s.loadingText = ""
	case spinner.TickMsg:
		if s.loading {
			s.spinner, cmd = s.spinner.Update(msg)
		}
	}

	return s, cmd
}

func (s *StatusBar) View() string {
	left := fmt.Sprintf("📡 %s", s.connectionState)

	if s.loading && s.loadingText != "" {
		left = fmt.Sprintf("%s %s", s.spinner.View(), s.loadingText)
	}

	right := s.helpText

	leftStyle := StatusBarStyle.Copy().Width(s.width / 2)
	rightStyle := StatusBarStyle.Copy().Width(s.width / 2).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
}
