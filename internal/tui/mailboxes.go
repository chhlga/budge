package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// mailboxItem implements list.Item interface
type mailboxItem struct {
	name string
}

func (m mailboxItem) Title() string       { return m.name }
func (m mailboxItem) Description() string { return "" }
func (m mailboxItem) FilterValue() string { return m.name }

// mailboxDelegate handles rendering of mailbox items
type mailboxDelegate struct{}

func (d mailboxDelegate) Height() int                             { return 1 }
func (d mailboxDelegate) Spacing() int                            { return 0 }
func (d mailboxDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d mailboxDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	mailbox := item.(mailboxItem)
	str := mailbox.name

	if index == m.Index() {
		str = SelectedItemStyle.Render("▶ " + str)
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
}

// MailboxList is the mailbox list view
type MailboxList struct {
	list list.Model
	keys KeyMap
}

// NewMailboxList creates a new mailbox list view
func NewMailboxList(keys KeyMap) MailboxList {
	items := []list.Item{}
	delegate := mailboxDelegate{}

	l := list.New(items, delegate, 0, 0)
	l.Title = "Mailboxes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return MailboxList{
		list: l,
		keys: keys,
	}
}

// SetSize updates the mailbox list dimensions
func (m *MailboxList) SetSize(width, height int) {
	m.list.SetSize(width, height-3) // Account for title and status bar
}

// SetMailboxes updates the mailbox list
func (m *MailboxList) SetMailboxes(mailboxes []string) {
	items := make([]list.Item, len(mailboxes))
	for i, mb := range mailboxes {
		items[i] = mailboxItem{name: mb}
	}
	m.list.SetItems(items)
}

// Init initializes the mailbox list
func (m MailboxList) Init() tea.Cmd {
	return nil
}

// Update handles messages for mailbox list
func (m MailboxList) Update(msg tea.Msg) (MailboxList, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if selected := m.list.SelectedItem(); selected != nil {
				mailbox := selected.(mailboxItem)
				return m, func() tea.Msg {
					return MailboxSelectedMsg{Mailbox: mailbox.name}
				}
			}
		}
	case MailboxesLoadedMsg:
		m.SetMailboxes(msg.Mailboxes)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the mailbox list
func (m MailboxList) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		Render(m.list.View())
}
