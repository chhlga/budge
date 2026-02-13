package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhlga/budge/internal/email"
)

type SortMode int

const (
	SortDateNewest SortMode = iota
	SortDateOldest
	SortSenderAZ
	SortSenderZA
	SortSubjectAZ
	SortSubjectZA
	SortUnreadFirst
	SortReadFirst
)

func (s SortMode) String() string {
	switch s {
	case SortDateNewest:
		return "Date (Newest)"
	case SortDateOldest:
		return "Date (Oldest)"
	case SortSenderAZ:
		return "Sender (A-Z)"
	case SortSenderZA:
		return "Sender (Z-A)"
	case SortSubjectAZ:
		return "Subject (A-Z)"
	case SortSubjectZA:
		return "Subject (Z-A)"
	case SortUnreadFirst:
		return "Unread First"
	case SortReadFirst:
		return "Read First"
	default:
		return "Date (Newest)"
	}
}

func (s SortMode) Next() SortMode {
	return (s + 1) % 8
}

type FilterMode int

const (
	FilterNone FilterMode = iota
	FilterUnread
	FilterRead
	FilterAttachments
)

func (f FilterMode) String() string {
	switch f {
	case FilterNone:
		return "All"
	case FilterUnread:
		return "Unread"
	case FilterRead:
		return "Read"
	case FilterAttachments:
		return "Attachments"
	default:
		return "All"
	}
}

func (f FilterMode) Next() FilterMode {
	return (f + 1) % 4
}

// emailItem implements list.Item interface
type emailItem struct {
	msg email.Message
}

func (e emailItem) Title() string {
	from := "Unknown"
	if len(e.msg.From) > 0 {
		from = e.msg.From[0].String()
	}
	return fmt.Sprintf("%s - %s", from, e.msg.Subject)
}

func (e emailItem) Description() string {
	return e.msg.Date.Format("Jan 02 15:04")
}

func (e emailItem) FilterValue() string {
	from := ""
	if len(e.msg.From) > 0 {
		from = e.msg.From[0].String()
	}
	return from + " " + e.msg.Subject
}

// emailDelegate handles rendering of email items
type emailDelegate struct{}

func (d emailDelegate) Height() int                             { return 2 }
func (d emailDelegate) Spacing() int                            { return 1 }
func (d emailDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d emailDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	email := item.(emailItem)

	// Determine style based on read/unread status
	style := ReadStyle
	if email.msg.IsUnread() {
		style = UnreadStyle
	}

	// Format the email preview
	from := "Unknown"
	if len(email.msg.From) > 0 {
		from = email.msg.From[0].String()
		if len(from) > 30 {
			from = from[:27] + "..."
		}
	}

	subject := email.msg.Subject
	if len(subject) > 50 {
		subject = subject[:47] + "..."
	}

	dateStr := email.msg.Date.Format("Jan 02 15:04")

	line1 := fmt.Sprintf("%-30s %s", from, dateStr)
	line2 := subject

	// Apply selection styling
	if index == m.Index() {
		line1 = SelectedItemStyle.Render("▶ " + line1)
		line2 = SelectedItemStyle.Render("  " + line2)
	} else {
		line1 = style.Render("  " + line1)
		line2 = style.Render("  " + line2)
	}

	fmt.Fprintf(w, "%s\n%s", line1, line2)
}

// EmailList is the email list view
type EmailList struct {
	list       list.Model
	keys       KeyMap
	mailbox    string
	total      uint32
	emails     []email.Message
	sortMode   SortMode
	filterMode FilterMode
}

func (e *EmailList) markSeenLocal(uid uint32, seen bool) {
	e.emails = markSeenInSlice(e.emails, uid, seen)
	e.applyFiltersAndSort()
}

func (e *EmailList) ClearFilter() {
	e.filterMode = FilterNone
	e.list.ResetFilter()
	e.applyFiltersAndSort()
}

// NewEmailList creates a new email list view
func NewEmailList(keys KeyMap) EmailList {
	items := []list.Item{}
	delegate := emailDelegate{}

	l := list.New(items, delegate, 0, 0)
	l.Title = "Emails"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return EmailList{
		list: l,
		keys: keys,
	}
}

// SetSize updates the email list dimensions
func (e *EmailList) SetSize(width, height int) {
	e.list.SetSize(width, height-3)
}

// SetEmails updates the email list
func (e *EmailList) SetEmails(emails []email.Message, total uint32) {
	e.emails = emails
	e.total = total
	e.applyFiltersAndSort()
}

func (e *EmailList) applyFiltersAndSort() {
	filtered := e.filterEmails(e.emails)
	sorted := e.sortEmails(filtered)

	items := make([]list.Item, len(sorted))
	for i, msg := range sorted {
		items[i] = emailItem{msg: msg}
	}
	e.list.SetItems(items)
	e.updateTitle()
}

func (e *EmailList) filterEmails(emails []email.Message) []email.Message {
	if e.filterMode == FilterNone {
		return emails
	}

	filtered := make([]email.Message, 0)
	for _, msg := range emails {
		switch e.filterMode {
		case FilterUnread:
			if msg.IsUnread() {
				filtered = append(filtered, msg)
			}
		case FilterRead:
			if !msg.IsUnread() {
				filtered = append(filtered, msg)
			}
		case FilterAttachments:
			if len(msg.Attachments) > 0 {
				filtered = append(filtered, msg)
			}
		}
	}
	return filtered
}

func (e *EmailList) sortEmails(emails []email.Message) []email.Message {
	sorted := make([]email.Message, len(emails))
	copy(sorted, emails)

	switch e.sortMode {
	case SortDateNewest:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Date.After(sorted[j].Date)
		})
	case SortDateOldest:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Date.Before(sorted[j].Date)
		})
	case SortSenderAZ:
		sort.Slice(sorted, func(i, j int) bool {
			fromI := ""
			if len(sorted[i].From) > 0 {
				fromI = sorted[i].From[0].Email
			}
			fromJ := ""
			if len(sorted[j].From) > 0 {
				fromJ = sorted[j].From[0].Email
			}
			return strings.ToLower(fromI) < strings.ToLower(fromJ)
		})
	case SortSenderZA:
		sort.Slice(sorted, func(i, j int) bool {
			fromI := ""
			if len(sorted[i].From) > 0 {
				fromI = sorted[i].From[0].Email
			}
			fromJ := ""
			if len(sorted[j].From) > 0 {
				fromJ = sorted[j].From[0].Email
			}
			return strings.ToLower(fromI) > strings.ToLower(fromJ)
		})
	case SortSubjectAZ:
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i].Subject) < strings.ToLower(sorted[j].Subject)
		})
	case SortSubjectZA:
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i].Subject) > strings.ToLower(sorted[j].Subject)
		})
	case SortUnreadFirst:
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].IsUnread() == sorted[j].IsUnread() {
				return sorted[i].Date.After(sorted[j].Date)
			}
			return sorted[i].IsUnread()
		})
	case SortReadFirst:
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].IsUnread() == sorted[j].IsUnread() {
				return sorted[i].Date.After(sorted[j].Date)
			}
			return !sorted[i].IsUnread()
		})
	}

	return sorted
}

// SetMailbox sets the current mailbox name
func (e *EmailList) SetMailbox(mailbox string) {
	e.mailbox = mailbox
	e.updateTitle()
}

func (e *EmailList) updateTitle() {
	title := e.mailbox

	displayCount := e.list.Items()
	if e.total > 0 {
		if e.filterMode != FilterNone {
			title = fmt.Sprintf("%s (%d/%d) [%s] [%s]", e.mailbox, len(displayCount), e.total, e.filterMode, e.sortMode)
		} else {
			title = fmt.Sprintf("%s (%d) [%s]", e.mailbox, e.total, e.sortMode)
		}
	}
	e.list.Title = title
}

// Init initializes the email list
func (e EmailList) Init() tea.Cmd {
	return nil
}

// Update handles messages for email list
func (e EmailList) Update(msg tea.Msg) (EmailList, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keys.Enter):
			if selected := e.list.SelectedItem(); selected != nil {
				email := selected.(emailItem).msg
				return e, func() tea.Msg {
					return EmailSelectedMsg{Email: email}
				}
			}
		case key.Matches(msg, e.keys.MarkRead):
			idx := e.list.Index()
			if idx >= 0 && idx < len(e.emails) {
				selectedEmail := e.emails[idx]
				isUnread := selectedEmail.IsUnread()
				return e, func() tea.Msg {
					return MarkReadRequestMsg{UID: selectedEmail.UID, Read: isUnread}
				}
			}
		case key.Matches(msg, e.keys.Delete):
			idx := e.list.Index()
			if idx >= 0 && idx < len(e.emails) {
				selectedEmail := e.emails[idx]
				return e, func() tea.Msg {
					return DeleteEmailRequestMsg{UID: selectedEmail.UID}
				}
			}
		case key.Matches(msg, e.keys.Sort):
			// Cycle to next sort mode
			e.sortMode = e.sortMode.Next()
			e.applyFiltersAndSort()
		case key.Matches(msg, e.keys.Filter):
			// Cycle to next filter mode
			e.filterMode = e.filterMode.Next()
			e.applyFiltersAndSort()
		}
	case EmailsLoadedMsg:
		e.SetEmails(msg.Emails, msg.Total)
	case MailboxSelectedMsg:
		e.SetMailbox(msg.Mailbox)
	}

	e.list, cmd = e.list.Update(msg)
	return e, cmd
}

// View renders the email list
func (e EmailList) View() string {
	var sb strings.Builder
	sb.WriteString(e.list.View())
	return sb.String()
}
