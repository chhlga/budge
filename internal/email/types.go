package email

import (
	"time"
)

// Message represents an email message with all its metadata and content
type Message struct {
	UID         uint32
	Flags       []string
	From        []Address
	To          []Address
	Cc          []Address
	Bcc         []Address
	ReplyTo     []Address
	Subject     string
	Date        time.Time
	MessageID   string
	InReplyTo   string
	References  []string
	ContentType string
	Body        *Body
	Attachments []Attachment
}

// Address represents an email address with optional name
type Address struct {
	Name  string
	Email string
}

// Body represents the email body content
type Body struct {
	Text string // Plain text content
	HTML string // HTML content
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	ContentType string
	Size        int64
	Data        []byte
}

// String formats an address as "Name <email>" or just "email"
func (a Address) String() string {
	if a.Name != "" {
		return a.Name + " <" + a.Email + ">"
	}
	return a.Email
}

// HasFlag checks if the message has a specific flag
func (m *Message) HasFlag(flag string) bool {
	for _, f := range m.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

// IsUnread returns true if the message does not have the \Seen flag
func (m *Message) IsUnread() bool {
	return !m.HasFlag("\\Seen")
}

// IsFlagged returns true if the message has the \Flagged flag
func (m *Message) IsFlagged() bool {
	return m.HasFlag("\\Flagged")
}
