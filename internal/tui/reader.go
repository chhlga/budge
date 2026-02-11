package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chhlga/budge/internal/email"
)

// EmailReader is the email reader view with scrolling viewport
type EmailReader struct {
	viewport viewport.Model
	email    *email.Message
	body     string
	ready    bool
	width    int
	height   int
}

// NewEmailReader creates a new email reader view
func NewEmailReader(keys KeyMap) EmailReader {
	return EmailReader{
		ready: false,
	}
}

// SetSize updates the email reader dimensions
func (r *EmailReader) SetSize(width, height int) {
	r.width = width
	r.height = height

	headerHeight := 5 // Space for from, to, subject, date
	footerHeight := 1 // Space for status bar

	if !r.ready {
		r.viewport = viewport.New(width, height-headerHeight-footerHeight)
		r.viewport.YPosition = headerHeight
		r.viewport.MouseWheelEnabled = true
		r.ready = true
	} else {
		r.viewport.Width = width
		r.viewport.Height = height - headerHeight - footerHeight
	}
}

// SetEmail sets the current email
func (r *EmailReader) SetEmail(msg email.Message) {
	r.email = &msg
	r.body = "" // Reset body, will be loaded separately
}

// SetBody sets the rendered email body
func (r *EmailReader) SetBody(body string) {
	r.body = body
	r.viewport.SetContent(body)
	r.viewport.GotoTop()
}

// Init initializes the email reader
func (r EmailReader) Init() tea.Cmd {
	return nil
}

// Update handles messages for email reader
func (r EmailReader) Update(msg tea.Msg) (EmailReader, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case EmailSelectedMsg:
		r.SetEmail(msg.Email)
	case EmailBodyLoadedMsg:
		if r.email != nil && r.email.UID == msg.UID {
			r.SetBody(msg.Body)
		}
	}

	r.viewport, cmd = r.viewport.Update(msg)
	return r, cmd
}

// View renders the email reader
func (r EmailReader) View() string {
	if !r.ready {
		return "Loading..."
	}

	if r.email == nil {
		return lipgloss.NewStyle().
			Width(r.width).
			Height(r.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No email selected")
	}

	// Render header
	from := "Unknown"
	if len(r.email.From) > 0 {
		from = r.email.From[0].String()
	}

	to := "Unknown"
	if len(r.email.To) > 0 {
		to = r.email.To[0].String()
	}

	headerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Padding(0, 1)

	header := headerStyle.Render(fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\nDate: %s",
		from,
		to,
		r.email.Subject,
		r.email.Date.Format("Mon, Jan 02, 2006 at 15:04"),
	))

	// Render body viewport
	var body string
	if r.body == "" {
		body = lipgloss.NewStyle().
			Foreground(dimColor).
			Render("Loading email body...")
	} else {
		body = r.viewport.View()
	}

	// Footer with scroll position
	footer := lipgloss.NewStyle().
		Foreground(dimColor).
		Render(fmt.Sprintf("%3.f%%", r.viewport.ScrollPercent()*100))

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		footer,
	)
}
