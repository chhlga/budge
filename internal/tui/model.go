package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chhlga/budge/internal/cache"
	"github.com/chhlga/budge/internal/config"
	"github.com/chhlga/budge/internal/imap"
)

// viewState represents the current active view
type viewState uint

const (
	mailboxListView viewState = iota
	emailListView
	emailReaderView
	searchView
)

// Model is the root TUI model
type Model struct {
	state  viewState
	keys   KeyMap
	width  int
	height int
	err    error
	focus  paneFocus

	// Sub-models
	mailboxList MailboxList
	emailList   EmailList
	emailReader EmailReader
	search      Search
	statusBar   StatusBar

	// Services
	imapClient *imap.Client
	cache      *cache.Cache
	config     *config.Config

	currentMailbox string
	loading        bool
	loadingText    string

	inSearchResults     bool
	preSearchEmailState EmailsLoadedMsg
}

type paneFocus uint

const (
	focusLeft paneFocus = iota
	focusRight
)

// NewModel creates a new root model
func NewModel(cfg *config.Config, client *imap.Client) Model {
	keys := NewKeyMap()

	return Model{
		state:       mailboxListView,
		keys:        keys,
		mailboxList: NewMailboxList(keys),
		emailList:   NewEmailList(keys),
		emailReader: NewEmailReader(keys),
		search:      NewSearch(keys),
		statusBar:   NewStatusBar(),
		imapClient:  client,
		cache:       cache.New(100), // Cache 100 email bodies
		config:      cfg,
		focus:       focusRight,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.mailboxList.Init(),
		m.emailList.Init(),
		m.emailReader.Init(),
		m.search.Init(),
		connectCmd(m.imapClient),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Global message handling
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.focus == focusLeft {
				m.focus = focusRight
			} else {
				m.focus = focusLeft
			}
			return m, nil
		}

		switch msg.String() {
		case "h":
			m.focus = focusLeft
			return m, nil
		case "l":
			m.focus = focusRight
			return m, nil
		}

		if msg.Type == tea.KeyEsc && m.state == emailReaderView {
			m.state = emailListView
			m.focus = focusRight
			m.statusBar.SetHelpText("enter: read | s: sort | f: filter | m: mark | d: delete | /: search | q: quit")
			return m, nil
		}

		// Global keys (always active)
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ViewMailboxes):
			m.state = emailListView
			m.focus = focusLeft
			m.statusBar.SetHelpText("enter: select | /: search | q: quit")
			return m, nil
		case key.Matches(msg, m.keys.ViewEmails):
			m.state = emailListView
			m.focus = focusRight
			m.statusBar.SetHelpText("enter: read | s: sort | f: filter | m: mark | d: delete | /: search | q: quit")
			if m.currentMailbox != "" {
				interval := time.Duration(m.config.Behavior.PollInterval) * time.Second
				return m, startMonitoringCmd(m.imapClient, m.currentMailbox, interval)
			}
			return m, nil
		case key.Matches(msg, m.keys.ViewReader):
			if m.emailReader.email != nil {
				m.state = emailReaderView
				m.focus = focusRight
			}
			m.statusBar.SetHelpText("2: back to list | q: quit")
			return m, nil
		case key.Matches(msg, m.keys.Search):
			m.state = searchView
			m.focus = focusRight
			m.statusBar.SetHelpText("enter: search | esc: cancel")
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		// Propagate size to all sub-models
		statusBarHeight := 1
		headerHeight := 1
		availableHeight := m.height - statusBarHeight - headerHeight

		mailboxWidth := clampInt(m.width/4, 24, 40)
		emailListWidth := clampInt((m.width*2)/5, 40, 60)
		rightWidth := m.width - mailboxWidth

		switch m.state {
		case emailReaderView:
			m.emailList.SetSize(emailListWidth, availableHeight)
			m.emailReader.SetSize(m.width-emailListWidth, availableHeight)
			m.mailboxList.SetSize(mailboxWidth, availableHeight)
			m.search.SetSize(rightWidth, availableHeight)
		case searchView:
			m.mailboxList.SetSize(mailboxWidth, availableHeight)
			m.search.SetSize(rightWidth, availableHeight)
			m.emailList.SetSize(rightWidth, availableHeight)
			m.emailReader.SetSize(rightWidth, availableHeight)
		default:
			m.mailboxList.SetSize(mailboxWidth, availableHeight)
			m.emailList.SetSize(rightWidth, availableHeight)
			m.emailReader.SetSize(rightWidth, availableHeight)
			m.search.SetSize(rightWidth, availableHeight)
		}
		m.statusBar.SetSize(m.width)

	case ErrorMsg:
		m.err = msg.Err
		return m, nil

	case ConnectCompleteMsg:
		return m, tea.Batch(
			func() tea.Msg { return ConnectionStateChangedMsg{State: m.imapClient.State()} },
			loadMailboxesCmd(m.imapClient),
		)

	case MailboxesLoadedMsg:
		m.mailboxList.SetMailboxes(msg.Mailboxes)
		if m.currentMailbox == "" {
			defaultMailbox := pickDefaultMailbox(msg.Mailboxes, m.config.Behavior.DefaultFolder)
			if defaultMailbox != "" {
				return m, func() tea.Msg { return MailboxSelectedMsg{Mailbox: defaultMailbox} }
			}
		}
		return m, nil

	case ConnectErrorMsg:
		m.err = msg.Err
		return m, nil

	case MailboxSelectedMsg:
		m.state = emailListView
		m.focus = focusRight
		m.statusBar.SetHelpText("enter: read | s: sort | f: filter | m: mark | d: delete | /: search | q: quit")
		m.emailList.SetMailbox(msg.Mailbox)
		m.currentMailbox = msg.Mailbox

		interval := time.Duration(m.config.Behavior.PollInterval) * time.Second
		return m, tea.Batch(
			loadEmailsCmd(m.imapClient, msg.Mailbox, uint32(m.config.Behavior.PageSize)),
			startMonitoringCmd(m.imapClient, msg.Mailbox, interval),
		)
	case EmailsLoadedMsg:
		m.emailList.SetEmails(msg.Emails, msg.Total)
		m.statusBar.SetHelpText("enter: read | s: sort | f: filter | m: mark | d: delete | /: search | q: quit")
		return m, nil

	case EmailSelectedMsg:
		selectedEmail := msg.Email
		if selectedEmail.IsUnread() {
			selectedEmail.Flags = addFlag(selectedEmail.Flags, "\\Seen")
			m.emailList.markSeenLocal(selectedEmail.UID, true)
			if m.inSearchResults {
				m.preSearchEmailState.Emails = markSeenInSlice(m.preSearchEmailState.Emails, selectedEmail.UID, true)
			}
		}
		m.state = emailReaderView
		m.statusBar.SetHelpText("2: back to list | q: quit")
		m.emailReader.SetEmail(selectedEmail)
		cmds = append(cmds, loadEmailBodyCmd(m.imapClient, m.cache, selectedEmail.UID))
		if msg.Email.IsUnread() {
			cmds = append(cmds, markReadCmd(m.imapClient, selectedEmail.UID, true))
		}
		return m, tea.Batch(cmds...)

	case EmailBodyLoadedMsg:
		m.emailReader.SetBody(msg.Body)
		return m, nil

	case MarkReadRequestMsg:
		return m, markReadCmd(m.imapClient, msg.UID, msg.Read)

	case DeleteEmailRequestMsg:
		return m, deleteEmailCmd(m.imapClient, msg.UID)

	case SearchQueryMsg:
		m.inSearchResults = true
		m.preSearchEmailState = EmailsLoadedMsg{Emails: m.emailList.emails, Total: m.emailList.total}
		m.state = emailListView
		m.statusBar.SetHelpText("Searching...")
		if m.currentMailbox == "" {
			m.currentMailbox = m.config.Behavior.DefaultFolder
		}
		return m, tea.Batch(
			func() tea.Msg { return LoadingMsg{Text: "Searching..."} },
			searchEmailsCmd(m.imapClient, m.currentMailbox, msg.Query),
		)

	case SearchCancelledMsg:
		m.state = emailListView
		m.emailList.ClearFilter()
		m.inSearchResults = false
		m.statusBar, cmd = m.statusBar.Update(LoadingClearedMsg{})
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case NewEmailMsg:
		if msg.Mailbox == m.currentMailbox {
			return m, loadEmailsCmd(m.imapClient, msg.Mailbox, uint32(m.config.Behavior.PageSize))
		}
		return m, nil

	case StartIdleMonitoringMsg:
		if msg.Mailbox != "" {
			m.currentMailbox = msg.Mailbox
			interval := time.Duration(m.config.Behavior.PollInterval) * time.Second
			return m, startMonitoringCmd(m.imapClient, msg.Mailbox, interval)
		}
		return m, nil

	case StopIdleMonitoringMsg:
		if m.currentMailbox != "" {
			return m, stopMonitoringCmd(m.currentMailbox)
		}
		return m, nil
	}

	// Update status bar
	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	// Delegate to active view
	switch m.state {
	case mailboxListView:
		m.mailboxList, cmd = m.mailboxList.Update(msg)
	case emailListView:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc && m.inSearchResults {
			m.inSearchResults = false
			m.emailList.ClearFilter()
			m.emailList.SetEmails(m.preSearchEmailState.Emails, m.preSearchEmailState.Total)
			m.statusBar, cmd = m.statusBar.Update(LoadingClearedMsg{})
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
		if m.focus == focusLeft {
			m.mailboxList, cmd = m.mailboxList.Update(msg)
		} else {
			m.emailList, cmd = m.emailList.Update(msg)
		}
	case emailReaderView:
		if m.focus == focusLeft {
			m.emailList, cmd = m.emailList.Update(msg)
		} else {
			m.emailReader, cmd = m.emailReader.Update(msg)
		}
	case searchView:
		if m.focus == focusLeft {
			m.mailboxList, cmd = m.mailboxList.Update(msg)
		} else {
			m.search, cmd = m.search.Update(msg)
		}
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the current view
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Render error if present
	if m.err != nil {
		errorView := ErrorStyle.Render("Error: " + m.err.Error())
		return lipgloss.JoinVertical(lipgloss.Left,
			errorView,
			m.statusBar.View(),
		)
	}

	statusBarHeight := 1
	headerHeight := 1
	contentHeight := clampInt(m.height-statusBarHeight-headerHeight, 0, m.height)

	header := HeaderStyle.Copy().
		Width(m.width).
		MaxWidth(m.width).
		Height(headerHeight).
		MaxHeight(headerHeight).
		Render(m.breadcrumb())

	mainView := lipgloss.NewStyle().
		Height(contentHeight).
		MaxHeight(contentHeight).
		Render(m.mainView(contentHeight))

	return lipgloss.JoinVertical(lipgloss.Left, header, mainView, m.statusBar.View())
}

func (m Model) breadcrumb() string {
	if m.currentMailbox == "" {
		return "Mailboxes"
	}

	suffix := ""
	switch m.state {
	case searchView:
		suffix = "  > Search"
	case emailReaderView:
		suffix = "  > Reader"
	default:
		suffix = "  > Emails"
	}

	focus := "right"
	if m.focus == focusLeft {
		focus = "left"
	}

	return fmt.Sprintf("%s%s  (focus: %s)", m.currentMailbox, suffix, focus)
}

func (m Model) mainView(contentHeight int) string {
	mailboxWidth := clampInt(m.width/4, 24, 40)
	emailListWidth := clampInt((m.width*2)/5, 40, 60)
	rightWidth := m.width - mailboxWidth

	leftFocusStyle := SelectedItemStyle
	rightFocusStyle := SelectedItemStyle
	if m.focus != focusLeft {
		leftFocusStyle = ReadStyle
	}
	if m.focus != focusRight {
		rightFocusStyle = ReadStyle
	}

	switch m.state {
	case emailReaderView:
		left := leftFocusStyle.Copy().Width(emailListWidth).MaxWidth(emailListWidth).Height(contentHeight).MaxHeight(contentHeight).Render(m.emailList.View())
		right := rightFocusStyle.Copy().Width(m.width - emailListWidth).MaxWidth(m.width - emailListWidth).Height(contentHeight).MaxHeight(contentHeight).Render(m.emailReader.View())
		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	case searchView:
		left := leftFocusStyle.Copy().Width(mailboxWidth).MaxWidth(mailboxWidth).Height(contentHeight).MaxHeight(contentHeight).Render(m.mailboxList.View())
		right := rightFocusStyle.Copy().Width(rightWidth).MaxWidth(rightWidth).Height(contentHeight).MaxHeight(contentHeight).Render(m.search.View())
		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	default:
		left := leftFocusStyle.Copy().Width(mailboxWidth).MaxWidth(mailboxWidth).Height(contentHeight).MaxHeight(contentHeight).Render(m.mailboxList.View())
		right := rightFocusStyle.Copy().Width(rightWidth).MaxWidth(rightWidth).Height(contentHeight).MaxHeight(contentHeight).Render(m.emailList.View())
		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	}
}
