package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines key bindings for the application
type KeyMap struct {
	// Global keys
	Quit    key.Binding
	Refresh key.Binding
	Help    key.Binding

	// Navigation keys
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Search key.Binding

	// View switching keys
	ViewMailboxes key.Binding
	ViewEmails    key.Binding
	ViewReader    key.Binding

	// Email actions
	MarkRead key.Binding
	Delete   key.Binding
	Sort     key.Binding
	Filter   key.Binding
}

// NewKeyMap creates a new KeyMap with default bindings
func NewKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "ctrl+r"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		ViewMailboxes: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "mailboxes"),
		),
		ViewEmails: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "emails"),
		),
		ViewReader: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "reader"),
		),
		MarkRead: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "mark read/unread"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Sort: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "cycle sort"),
		),
		Filter: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "toggle filter"),
		),
	}
}
