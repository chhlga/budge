package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Search is the search view
type Search struct {
	textInput textinput.Model
	keys      KeyMap
	width     int
	height    int
}

// NewSearch creates a new search view
func NewSearch(keys KeyMap) Search {
	ti := textinput.New()
	ti.Placeholder = "Search emails..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return Search{
		textInput: ti,
		keys:      keys,
	}
}

// SetSize updates the search view dimensions
func (s *Search) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.textInput.Width = width - 4
}

// Init initializes the search view
func (s Search) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for search view
func (s Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			query := s.textInput.Value()
			if query != "" {
				return s, func() tea.Msg {
					return SearchQueryMsg{Query: query}
				}
			}
		case tea.KeyEsc:
			s.textInput.SetValue("")
			s.textInput.Blur()
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

// View renders the search view
func (s Search) View() string {
	style := lipgloss.NewStyle().
		Width(s.width).
		Height(s.height).
		Padding(2, 4)

	content := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle.Render("Search"),
		"",
		s.textInput.View(),
		"",
		StatusBarStyle.Render("Enter to search, Esc to cancel"),
	)

	return style.Render(content)
}
