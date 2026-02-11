package tui

import (
	"github.com/chhlga/budge/internal/email"
	"github.com/chhlga/budge/internal/imap"
)

// Custom message types for inter-component communication

// ConnectCompleteMsg is sent when IMAP connection is established
type ConnectCompleteMsg struct{}

// ConnectErrorMsg is sent when IMAP connection fails
type ConnectErrorMsg struct {
	Err error
}

// MailboxesLoadedMsg is sent when mailbox list is fetched
type MailboxesLoadedMsg struct {
	Mailboxes []string
}

// MailboxSelectedMsg is sent when user selects a mailbox
type MailboxSelectedMsg struct {
	Mailbox string
}

// EmailsLoadedMsg is sent when email list is fetched
type EmailsLoadedMsg struct {
	Emails []email.Message
	Total  uint32
}

// EmailSelectedMsg is sent when user selects an email
type EmailSelectedMsg struct {
	Email email.Message
}

// EmailBodyLoadedMsg is sent when email body is fetched and rendered
type EmailBodyLoadedMsg struct {
	UID  uint32
	Body string
}

// SearchQueryMsg is sent when user submits search query
type SearchQueryMsg struct {
	Query string
}

// ConnectionStateChangedMsg is sent when connection state changes
type ConnectionStateChangedMsg struct {
	State imap.ConnectionState
}

// ErrorMsg is a generic error message
type ErrorMsg struct {
	Err error
}

// MarkReadRequestMsg requests marking an email as read/unread
type MarkReadRequestMsg struct {
	UID  uint32
	Read bool
}

// DeleteEmailRequestMsg requests deleting an email
type DeleteEmailRequestMsg struct {
	UID uint32
}
