package tui

import (
	"testing"
	"time"

	"github.com/chhlga/budge/internal/email"
)

func TestNewEmailReader(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)

	if r.ready {
		t.Error("NewEmailReader() ready = true, want false")
	}

	if r.email != nil {
		t.Error("NewEmailReader() email != nil, want nil")
	}
}

func TestEmailReader_SetSize(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)

	r.SetSize(80, 24)

	if !r.ready {
		t.Error("SetSize() should set ready to true")
	}

	if r.width != 80 {
		t.Errorf("width = %d, want 80", r.width)
	}

	if r.height != 24 {
		t.Errorf("height = %d, want 24", r.height)
	}

	// Test resizing after initial setup
	r.SetSize(100, 30)

	if r.width != 100 {
		t.Errorf("After resize: width = %d, want 100", r.width)
	}
}

func TestEmailReader_SetEmail(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)

	msg := email.Message{
		UID:     123,
		Subject: "Test Subject",
		From:    []email.Address{{Email: "sender@example.com"}},
	}

	r.SetEmail(msg)

	if r.email == nil {
		t.Fatal("SetEmail() email = nil, want non-nil")
	}

	if r.email.UID != 123 {
		t.Errorf("email.UID = %d, want 123", r.email.UID)
	}

	if r.body != "" {
		t.Errorf("body = %q, want empty string after SetEmail", r.body)
	}
}

func TestEmailReader_SetBody(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)
	r.SetSize(80, 24)

	body := "This is the email body content"
	r.SetBody(body)

	if r.body != body {
		t.Errorf("body = %q, want %q", r.body, body)
	}
}

func TestEmailReader_Init(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)

	cmd := r.Init()
	if cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestEmailReader_Update_EmailSelected(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)
	r.SetSize(80, 24)

	msg := EmailSelectedMsg{
		Email: email.Message{
			UID:     456,
			Subject: "Selected Email",
		},
	}

	updatedR, _ := r.Update(msg)

	if updatedR.email == nil {
		t.Fatal("Update() email = nil after EmailSelectedMsg")
	}

	if updatedR.email.UID != 456 {
		t.Errorf("email.UID = %d, want 456", updatedR.email.UID)
	}
}

func TestEmailReader_Update_EmailBodyLoaded(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)
	r.SetSize(80, 24)

	// Set an email first
	r.SetEmail(email.Message{UID: 789, Subject: "Test"})

	msg := EmailBodyLoadedMsg{
		UID:  789,
		Body: "Email body content",
	}

	updatedR, _ := r.Update(msg)

	if updatedR.body != "Email body content" {
		t.Errorf("body = %q, want %q", updatedR.body, "Email body content")
	}
}

func TestEmailReader_Update_EmailBodyLoaded_WrongUID(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)
	r.SetSize(80, 24)

	// Set an email first
	r.SetEmail(email.Message{UID: 789, Subject: "Test"})
	r.SetBody("Original body")

	// Load body for different UID
	msg := EmailBodyLoadedMsg{
		UID:  999,
		Body: "Different email body",
	}

	updatedR, _ := r.Update(msg)

	// Body should not change
	if updatedR.body != "Original body" {
		t.Errorf("body = %q, want %q (should not change for wrong UID)", updatedR.body, "Original body")
	}
}

func TestEmailReader_View_NotReady(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)

	view := r.View()
	if view != "Loading..." {
		t.Errorf("View() = %q, want %q", view, "Loading...")
	}
}

func TestEmailReader_View_NoEmail(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)
	r.SetSize(80, 24)

	view := r.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestEmailReader_View_WithEmail(t *testing.T) {
	keys := NewKeyMap()
	r := NewEmailReader(keys)
	r.SetSize(80, 24)

	msg := email.Message{
		UID:     123,
		Subject: "Test Subject",
		From:    []email.Address{{Email: "sender@example.com"}},
		To:      []email.Address{{Email: "receiver@example.com"}},
		Date:    time.Now(),
	}

	r.SetEmail(msg)
	r.SetBody("Email body")

	view := r.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}
