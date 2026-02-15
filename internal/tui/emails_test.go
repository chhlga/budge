package tui

import (
	"testing"
	"time"

	"github.com/chhlga/budge/internal/email"
)

func TestSortMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode SortMode
		want string
	}{
		{"date newest", SortDateNewest, "Date (Newest)"},
		{"date oldest", SortDateOldest, "Date (Oldest)"},
		{"sender a-z", SortSenderAZ, "Sender (A-Z)"},
		{"sender z-a", SortSenderZA, "Sender (Z-A)"},
		{"subject a-z", SortSubjectAZ, "Subject (A-Z)"},
		{"subject z-a", SortSubjectZA, "Subject (Z-A)"},
		{"unread first", SortUnreadFirst, "Unread First"},
		{"read first", SortReadFirst, "Read First"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("SortMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortMode_Next(t *testing.T) {
	tests := []struct {
		name    string
		current SortMode
		want    SortMode
	}{
		{"newest to oldest", SortDateNewest, SortDateOldest},
		{"oldest to sender a-z", SortDateOldest, SortSenderAZ},
		{"sender a-z to sender z-a", SortSenderAZ, SortSenderZA},
		{"sender z-a to subject a-z", SortSenderZA, SortSubjectAZ},
		{"subject a-z to subject z-a", SortSubjectAZ, SortSubjectZA},
		{"subject z-a to unread", SortSubjectZA, SortUnreadFirst},
		{"unread to read", SortUnreadFirst, SortReadFirst},
		{"read to newest (wrap)", SortReadFirst, SortDateNewest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.current.Next(); got != tt.want {
				t.Errorf("SortMode.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode FilterMode
		want string
	}{
		{"none", FilterNone, "All"},
		{"unread", FilterUnread, "Unread"},
		{"read", FilterRead, "Read"},
		{"attachments", FilterAttachments, "Attachments"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("FilterMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterMode_Next(t *testing.T) {
	tests := []struct {
		name    string
		current FilterMode
		want    FilterMode
	}{
		{"none to unread", FilterNone, FilterUnread},
		{"unread to read", FilterUnread, FilterRead},
		{"read to attachments", FilterRead, FilterAttachments},
		{"attachments to none (wrap)", FilterAttachments, FilterNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.current.Next(); got != tt.want {
				t.Errorf("FilterMode.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailItem_Title(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		from    []email.Address
		want    string
	}{
		{
			name:    "with sender name",
			subject: "Test Email",
			from:    []email.Address{{Name: "John Doe", Email: "john@example.com"}},
			want:    "John Doe <john@example.com> - Test Email",
		},
		{
			name:    "without sender name",
			subject: "Important",
			from:    []email.Address{{Email: "test@example.com"}},
			want:    "test@example.com - Important",
		},
		{
			name:    "no sender",
			subject: "Normal Email",
			from:    []email.Address{},
			want:    "Unknown - Normal Email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := emailItem{
				msg: email.Message{
					Subject: tt.subject,
					From:    tt.from,
				},
			}
			if got := item.Title(); got != tt.want {
				t.Errorf("emailItem.Title() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEmailItem_Description(t *testing.T) {
	date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "formatted date",
			date: date,
			want: "Jan 15 10:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := emailItem{
				msg: email.Message{
					Date: tt.date,
				},
			}
			if got := item.Description(); got != tt.want {
				t.Errorf("emailItem.Description() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEmailItem_FilterValue(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		from    []email.Address
		want    string
	}{
		{
			name:    "subject and sender",
			subject: "Test Subject",
			from:    []email.Address{{Name: "John Doe", Email: "john@example.com"}},
			want:    "John Doe <john@example.com> Test Subject",
		},
		{
			name:    "only subject",
			subject: "Hello World",
			from:    []email.Address{},
			want:    " Hello World",
		},
		{
			name:    "empty subject",
			subject: "",
			from:    []email.Address{{Email: "jane@example.com"}},
			want:    "jane@example.com ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := emailItem{
				msg: email.Message{
					Subject: tt.subject,
					From:    tt.from,
				},
			}
			if got := item.FilterValue(); got != tt.want {
				t.Errorf("emailItem.FilterValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEmailDelegate_Methods(t *testing.T) {
	delegate := emailDelegate{}

	if h := delegate.Height(); h != 2 {
		t.Errorf("emailDelegate.Height() = %d, want 2", h)
	}

	if s := delegate.Spacing(); s != 1 {
		t.Errorf("emailDelegate.Spacing() = %d, want 1", s)
	}

	if cmd := delegate.Update(nil, nil); cmd != nil {
		t.Error("emailDelegate.Update() should return nil")
	}
}

func TestEmailList_SetSize(t *testing.T) {
	keys := NewKeyMap()
	el := NewEmailList(keys)

	el.SetSize(200, 100)

	// Verify size was set (indirectly through list dimensions)
	// Since we can't directly check width/height, we verify it doesn't panic
	_ = el.View()
}

func TestEmailList_Init(t *testing.T) {
	keys := NewKeyMap()
	el := NewEmailList(keys)

	cmd := el.Init()

	// Init returns nil for this component
	if cmd != nil {
		t.Error("EmailList.Init() should return nil")
	}
}
