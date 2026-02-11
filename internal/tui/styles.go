package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#25A065")
	textColor      = lipgloss.Color("#FAFAFA")
	dimColor       = lipgloss.Color("#666666")
	errorColor     = lipgloss.Color("#FF0000")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	UnreadStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor)

	ReadStyle = lipgloss.NewStyle().
			Foreground(dimColor)
)
