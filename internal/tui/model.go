package tui

import (
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
		// Global keys (always active)
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ViewMailboxes):
			m.state = mailboxListView
			m.statusBar.SetHelpText("enter: select | r: refresh | q: quit")
			return m, stopMonitoringCmd(m.currentMailbox)
		case key.Matches(msg, m.keys.ViewEmails):
			m.state = emailListView
			m.statusBar.SetHelpText("enter: read | s: sort | f: filter | m: mark | d: delete | /: search | q: quit")
			if m.currentMailbox != "" {
				interval := time.Duration(m.config.Behavior.PollInterval) * time.Second
				return m, startMonitoringCmd(m.imapClient, m.currentMailbox, interval)
			}
			return m, nil
		case key.Matches(msg, m.keys.ViewReader):
			m.state = emailReaderView
			m.statusBar.SetHelpText("2: back to list | q: quit")
			return m, nil
		case key.Matches(msg, m.keys.Search):
			m.state = searchView
			m.statusBar.SetHelpText("enter: search | esc: cancel")
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		// Propagate size to all sub-models
		statusBarHeight := 1
		availableHeight := m.height - statusBarHeight

		m.mailboxList.SetSize(m.width, availableHeight)
		m.emailList.SetSize(m.width, availableHeight)
		m.emailReader.SetSize(m.width, availableHeight)
		m.search.SetSize(m.width, availableHeight)
		m.statusBar.SetSize(m.width)

	case ErrorMsg:
		m.err = msg.Err
		return m, nil

	case ConnectCompleteMsg:
		return m, tea.Batch(
			func() tea.Msg { return ConnectionStateChangedMsg{State: m.imapClient.State()} },
			loadMailboxesCmd(m.imapClient),
		)

	case ConnectErrorMsg:
		m.err = msg.Err
		return m, nil

	case MailboxSelectedMsg:
		m.state = emailListView
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
		m.emailList, cmd = m.emailList.Update(msg)
	case emailReaderView:
		m.emailReader, cmd = m.emailReader.Update(msg)
	case searchView:
		m.search, cmd = m.search.Update(msg)
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

	// Render active view
	var mainView string
	switch m.state {
	case mailboxListView:
		mainView = m.mailboxList.View()
	case emailListView:
		mainView = m.emailList.View()
	case emailReaderView:
		mainView = m.emailReader.View()
	case searchView:
		mainView = m.search.View()
	default:
		mainView = "Unknown view"
	}

	// Combine main view with status bar
	return lipgloss.JoinVertical(lipgloss.Left,
		mainView,
		m.statusBar.View(),
	)
}
