package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhlga/budge/internal/config"
	"github.com/chhlga/budge/internal/email"
)

func TestEscInSearchView_exitsSearchViewAndClearsFilter(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.state = emailListView
	m.emailList.SetMailbox("INBOX")
	m.emailList.SetEmails([]email.Message{{UID: 1, Subject: "hello"}}, 1)
	m.emailList.filterMode = FilterUnread

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = updated.(Model)
	if m.state != searchView {
		t.Fatalf("expected state=searchView, got %v", m.state)
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if cmd == nil {
		t.Fatalf("expected Esc to return a command")
	}
	updated, _ = m.Update(cmd())
	m = updated.(Model)
	if m.state != emailListView {
		t.Fatalf("expected state=emailListView after Esc, got %v", m.state)
	}
	if m.emailList.filterMode != FilterNone {
		t.Fatalf("expected filter cleared to FilterNone, got %v", m.emailList.filterMode)
	}
	if m.emailList.list.IsFiltered() {
		t.Fatalf("expected bubbles list filter to be reset")
	}
	if got := m.search.textInput.Value(); got != "" {
		t.Fatalf("expected search input cleared, got %q", got)
	}
}

func TestEscInEmailListView_exitsSearchResultsRestoresPreviousListAndClearsFilter(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.state = emailListView
	m.currentMailbox = "INBOX"
	m.emailList.SetMailbox("INBOX")

	before := []email.Message{{UID: 1, Subject: "a"}, {UID: 2, Subject: "b"}}
	m.emailList.SetEmails(before, uint32(len(before)))
	m.emailList.filterMode = FilterAttachments

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = updated.(Model)
	if m.state != searchView {
		t.Fatalf("expected state=searchView, got %v", m.state)
	}

	updated, _ = m.Update(SearchQueryMsg{Query: "hello"})
	m = updated.(Model)
	if m.state != emailListView {
		t.Fatalf("expected state=emailListView after submitting search query, got %v", m.state)
	}

	searchResults := []email.Message{{UID: 10, Subject: "hello"}}
	updated, _ = m.Update(EmailsLoadedMsg{Emails: searchResults, Total: uint32(len(searchResults))})
	m = updated.(Model)
	if len(m.emailList.emails) != 1 || m.emailList.emails[0].UID != 10 {
		t.Fatalf("expected search results to be loaded")
	}

	m.emailList.filterMode = FilterUnread

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if len(m.emailList.emails) != len(before) {
		t.Fatalf("expected previous emails restored, got %d", len(m.emailList.emails))
	}
	if m.emailList.emails[0].UID != before[0].UID || m.emailList.emails[1].UID != before[1].UID {
		t.Fatalf("expected restored emails to match snapshot")
	}
	if m.emailList.filterMode != FilterNone {
		t.Fatalf("expected filter cleared to FilterNone, got %v", m.emailList.filterMode)
	}
	if m.emailList.list.IsFiltered() {
		t.Fatalf("expected bubbles list filter to be reset")
	}
}

func TestClearFilter_resetsBubblesListFilterQuery(t *testing.T) {
	el := NewEmailList(NewKeyMap())
	el.SetMailbox("INBOX")
	el.SetEmails([]email.Message{{UID: 1, Subject: "hello"}}, 1)

	el.list.SetFilterText("hel")
	if !el.list.IsFiltered() {
		t.Fatalf("expected list to be filtered")
	}

	el.filterMode = FilterUnread
	el.ClearFilter()
	if el.filterMode != FilterNone {
		t.Fatalf("expected FilterNone")
	}
	if el.list.IsFiltered() {
		t.Fatalf("expected bubbles list filter to be reset")
	}
}

func TestEscInSearchView_clearsStatusBarLoading(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.state = searchView

	updated, _ := m.Update(LoadingMsg{Text: "Searching..."})
	m = updated.(Model)
	if !m.statusBar.loading {
		t.Fatalf("expected status bar loading=true")
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if cmd == nil {
		t.Fatalf("expected Esc to return a command")
	}
	updated, _ = m.Update(cmd())
	m = updated.(Model)

	if m.statusBar.loading {
		t.Fatalf("expected status bar loading=false after cancel")
	}
}

func TestEscInEmailListView_exitsSearchResultsAndClearsStatusBarLoading(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.state = emailListView
	m.currentMailbox = "INBOX"
	m.emailList.SetMailbox("INBOX")

	before := []email.Message{{UID: 1, Subject: "a"}, {UID: 2, Subject: "b"}}
	m.emailList.SetEmails(before, uint32(len(before)))

	updated, _ := m.Update(LoadingMsg{Text: "Searching..."})
	m = updated.(Model)
	if !m.statusBar.loading {
		t.Fatalf("expected status bar loading=true")
	}

	updated, _ = m.Update(SearchQueryMsg{Query: "hello"})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if m.statusBar.loading {
		t.Fatalf("expected status bar loading=false after exiting search results")
	}
}

func TestOpenUnreadEmail_marksItReadLocallyImmediately(t *testing.T) {
	cfg := &config.Config{Behavior: config.BehaviorConfig{DefaultFolder: "INBOX", PageSize: 50, PollInterval: 30}}

	m := NewModel(cfg, nil)
	m.state = emailListView
	m.currentMailbox = "INBOX"
	m.emailList.SetMailbox("INBOX")

	unread := email.Message{UID: 1, Subject: "hello", Flags: []string{}}
	m.emailList.SetEmails([]email.Message{unread}, 1)

	updated, _ := m.Update(EmailSelectedMsg{Email: unread})
	m = updated.(Model)

	if m.emailReader.email.IsUnread() {
		t.Fatalf("expected reader email to be marked read locally")
	}
	if m.emailList.emails[0].IsUnread() {
		t.Fatalf("expected list email to be marked read locally")
	}
}
