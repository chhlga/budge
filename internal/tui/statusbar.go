package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
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

	half := s.width / 2
	contentWidth := half - 2
	if contentWidth < 0 {
		contentWidth = 0
	}

	left = truncateString(left, contentWidth)
	right = truncateString(right, contentWidth)

	leftStyle := StatusBarStyle.Copy().Width(half).Height(1).MaxHeight(1)
	rightStyle := StatusBarStyle.Copy().Width(half).Height(1).MaxHeight(1).Align(lipgloss.Right)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)

	return lipgloss.NewStyle().Height(1).MaxHeight(1).Render(bar)
}

func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if ansi.PrintableRuneWidth(s) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	return truncate.StringWithTail(s, uint(max), "…")
}
