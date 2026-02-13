package tui

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhlga/budge/internal/cache"
	"github.com/chhlga/budge/internal/email"
	imapClient "github.com/chhlga/budge/internal/imap"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

func connectCmd(client *imapClient.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := client.Connect(ctx); err != nil {
			return ConnectErrorMsg{Err: fmt.Errorf("connection failed: %w", err)}
		}

		if err := client.Authenticate(ctx); err != nil {
			return ConnectErrorMsg{Err: fmt.Errorf("authentication failed: %w", err)}
		}

		return ConnectCompleteMsg{}
	}
}

func loadMailboxesCmd(client *imapClient.Client) tea.Cmd {
	return func() tea.Msg {
		if !client.IsConnected() {
			return ErrorMsg{Err: fmt.Errorf("not connected to IMAP server")}
		}

		imapConn := client.Client()
		if imapConn == nil {
			return ErrorMsg{Err: fmt.Errorf("IMAP client not initialized")}
		}

		listCmd := imapConn.List("", "*", nil)
		mailboxes, err := listCmd.Collect()
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to list mailboxes: %w", err)}
		}

		names := make([]string, 0, len(mailboxes))
		for _, mbox := range mailboxes {
			names = append(names, mbox.Mailbox)
		}

		names = sortMailboxes(names)

		return MailboxesLoadedMsg{Mailboxes: names}
	}
}

func loadEmailsCmd(client *imapClient.Client, mailbox string, pageSize uint32) tea.Cmd {
	return func() tea.Msg {
		if !client.IsConnected() {
			return ErrorMsg{Err: fmt.Errorf("not connected to IMAP server")}
		}

		imapConn := client.Client()
		if imapConn == nil {
			return ErrorMsg{Err: fmt.Errorf("IMAP client not initialized")}
		}

		selectData, err := imapConn.Select(mailbox, nil).Wait()
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to select mailbox %s: %w", mailbox, err)}
		}

		total := selectData.NumMessages
		if total == 0 {
			return EmailsLoadedMsg{Emails: []email.Message{}, Total: 0}
		}

		var seqSet imap.SeqSet
		if total <= pageSize {
			seqSet.AddRange(1, total)
		} else {
			start := total - pageSize + 1
			seqSet.AddRange(start, total)
		}

		fetchOptions := &imap.FetchOptions{
			UID:      true,
			Envelope: true,
			Flags:    true,
		}

		fetchCmd := imapConn.Fetch(seqSet, fetchOptions)
		messages := make([]email.Message, 0, pageSize)

		for {
			msgData := fetchCmd.Next()
			if msgData == nil {
				break
			}

			var uid imap.UID
			var envelope *imap.Envelope
			var flags []imap.Flag

			for {
				item := msgData.Next()
				if item == nil {
					break
				}

				switch item := item.(type) {
				case imapclient.FetchItemDataUID:
					uid = item.UID
				case imapclient.FetchItemDataEnvelope:
					envelope = item.Envelope
				case imapclient.FetchItemDataFlags:
					flags = item.Flags
				}
			}

			if envelope != nil {
				emailMsg := email.Message{
					UID:     uint32(uid),
					Subject: envelope.Subject,
					Date:    envelope.Date,
					Flags:   convertFlags(flags),
				}

				if len(envelope.From) > 0 {
					emailMsg.From = convertAddresses(envelope.From)
				}

				if len(envelope.To) > 0 {
					emailMsg.To = convertAddresses(envelope.To)
				}

				messages = append(messages, emailMsg)
			}
		}

		if err := fetchCmd.Close(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to fetch emails: %w", err)}
		}

		return EmailsLoadedMsg{Emails: messages, Total: total}
	}
}

func loadEmailBodyCmd(client *imapClient.Client, c *cache.Cache, uid uint32) tea.Cmd {
	return func() tea.Msg {
		cacheKey := strconv.FormatUint(uint64(uid), 10)

		if cachedBody, ok := c.Get(cacheKey); ok {
			if bodyStr, ok := cachedBody.(string); ok {
				return EmailBodyLoadedMsg{UID: uid, Body: bodyStr}
			}
		}

		if !client.IsConnected() {
			return ErrorMsg{Err: fmt.Errorf("not connected to IMAP server")}
		}

		imapConn := client.Client()
		if imapConn == nil {
			return ErrorMsg{Err: fmt.Errorf("IMAP client not initialized")}
		}

		var uidSet imap.UIDSet
		uidSet.AddNum(imap.UID(uid))

		fetchOptions := &imap.FetchOptions{
			UID:         true,
			BodySection: []*imap.FetchItemBodySection{{}},
		}

		fetchCmd := imapConn.Fetch(uidSet, fetchOptions)
		msgData := fetchCmd.Next()
		if msgData == nil {
			return ErrorMsg{Err: fmt.Errorf("email with UID %d not found", uid)}
		}

		var bodySection imapclient.FetchItemDataBodySection

		for {
			item := msgData.Next()
			if item == nil {
				break
			}

			if bs, ok := item.(imapclient.FetchItemDataBodySection); ok {
				bodySection = bs
				break
			}
		}

		if bodySection.Literal == nil {
			return ErrorMsg{Err: fmt.Errorf("email body not found for UID %d", uid)}
		}

		bodyBytes, err := io.ReadAll(bodySection.Literal)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to read email body: %w", err)}
		}

		parsedEmail, err := email.Parse(bodyBytes)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to parse email: %w", err)}
		}

		if err := fetchCmd.Close(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to close fetch command: %w", err)}
		}

		renderedBody, err := email.Render(parsedEmail.Body)
		if err != nil || renderedBody == "" {
			if parsedEmail.Body != nil {
				renderedBody = parsedEmail.Body.Text
			}
		}

		c.Set(cacheKey, renderedBody)

		return EmailBodyLoadedMsg{UID: uid, Body: renderedBody}
	}
}

func markReadCmd(client *imapClient.Client, uid uint32, read bool) tea.Cmd {
	return func() tea.Msg {
		if !client.IsConnected() {
			return ErrorMsg{Err: fmt.Errorf("not connected to IMAP server")}
		}

		imapConn := client.Client()
		if imapConn == nil {
			return ErrorMsg{Err: fmt.Errorf("IMAP client not initialized")}
		}

		var uidSet imap.UIDSet
		uidSet.AddNum(imap.UID(uid))

		storeFlags := imap.StoreFlags{
			Op:     imap.StoreFlagsSet,
			Flags:  []imap.Flag{imap.FlagSeen},
			Silent: true,
		}

		if !read {
			storeFlags.Op = imap.StoreFlagsDel
		}

		if err := imapConn.Store(uidSet, &storeFlags, nil).Close(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to mark email as read: %w", err)}
		}

		return nil
	}
}

func deleteEmailCmd(client *imapClient.Client, uid uint32) tea.Cmd {
	return func() tea.Msg {
		if !client.IsConnected() {
			return ErrorMsg{Err: fmt.Errorf("not connected to IMAP server")}
		}

		imapConn := client.Client()
		if imapConn == nil {
			return ErrorMsg{Err: fmt.Errorf("IMAP client not initialized")}
		}

		var uidSet imap.UIDSet
		uidSet.AddNum(imap.UID(uid))

		storeFlags := imap.StoreFlags{
			Op:     imap.StoreFlagsAdd,
			Flags:  []imap.Flag{imap.FlagDeleted},
			Silent: true,
		}

		if err := imapConn.Store(uidSet, &storeFlags, nil).Close(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to delete email: %w", err)}
		}

		if err := imapConn.Expunge().Close(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to expunge deleted email: %w", err)}
		}

		return nil
	}
}

func searchEmailsCmd(client *imapClient.Client, query string) tea.Cmd {
	return func() tea.Msg {
		if !client.IsConnected() {
			return ErrorMsg{Err: fmt.Errorf("not connected to IMAP server")}
		}

		imapConn := client.Client()
		if imapConn == nil {
			return ErrorMsg{Err: fmt.Errorf("IMAP client not initialized")}
		}

		criteria := &imap.SearchCriteria{
			Text: []string{query},
		}

		searchData, err := imapConn.Search(criteria, nil).Wait()
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("search failed: %w", err)}
		}

		allUIDs := searchData.AllUIDs()
		if len(allUIDs) == 0 {
			return EmailsLoadedMsg{Emails: []email.Message{}, Total: 0}
		}

		var uidSet imap.UIDSet
		for _, uid := range allUIDs {
			uidSet.AddNum(uid)
		}

		fetchOptions := &imap.FetchOptions{
			UID:      true,
			Envelope: true,
			Flags:    true,
		}

		fetchCmd := imapConn.Fetch(uidSet, fetchOptions)
		messages := make([]email.Message, 0)

		for {
			msgData := fetchCmd.Next()
			if msgData == nil {
				break
			}

			var uid imap.UID
			var envelope *imap.Envelope
			var flags []imap.Flag

			for {
				item := msgData.Next()
				if item == nil {
					break
				}

				switch item := item.(type) {
				case imapclient.FetchItemDataUID:
					uid = item.UID
				case imapclient.FetchItemDataEnvelope:
					envelope = item.Envelope
				case imapclient.FetchItemDataFlags:
					flags = item.Flags
				}
			}

			if envelope != nil {
				emailMsg := email.Message{
					UID:     uint32(uid),
					Subject: envelope.Subject,
					Date:    envelope.Date,
					Flags:   convertFlags(flags),
				}

				if len(envelope.From) > 0 {
					emailMsg.From = convertAddresses(envelope.From)
				}

				if len(envelope.To) > 0 {
					emailMsg.To = convertAddresses(envelope.To)
				}

				messages = append(messages, emailMsg)
			}
		}

		if err := fetchCmd.Close(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to fetch search results: %w", err)}
		}

		return EmailsLoadedMsg{Emails: messages, Total: uint32(len(messages))}
	}
}

func convertAddresses(imapAddrs []imap.Address) []email.Address {
	addresses := make([]email.Address, len(imapAddrs))
	for i, addr := range imapAddrs {
		addresses[i] = email.Address{
			Name:  addr.Name,
			Email: fmt.Sprintf("%s@%s", addr.Mailbox, addr.Host),
		}
	}
	return addresses
}

func convertFlags(imapFlags []imap.Flag) []string {
	flags := make([]string, len(imapFlags))
	for i, flag := range imapFlags {
		flags[i] = string(flag)
	}
	return flags
}

func sortMailboxes(mailboxes []string) []string {
	priorityOrder := []struct {
		name    string
		aliases []string
	}{
		{"INBOX", []string{"Inbox", "inbox", "INBOX"}},
		{"Sent", []string{"Sent", "sent", "SENT", "Sent Messages", "[Gmail]/Sent Mail"}},
		{"Drafts", []string{"Draft", "draft", "DRAFT", "Drafts", "DRAFTS", "[Gmail]/Drafts"}},
		{"All Mail", []string{"All Mail", "all mail", "ALL MAIL", "All mail", "[Gmail]/All Mail"}},
	}

	used := make(map[string]bool)
	var prioritized []string

	for _, entry := range priorityOrder {
		for _, name := range mailboxes {
			if used[name] {
				continue
			}

			for _, alias := range entry.aliases {
				if strings.EqualFold(name, alias) {
					prioritized = append(prioritized, name)
					used[name] = true
					break
				}
			}
		}
	}

	var others []string
	for _, name := range mailboxes {
		if !used[name] {
			others = append(others, name)
		}
	}

	sort.Slice(others, func(i, j int) bool {
		return strings.ToLower(others[i]) < strings.ToLower(others[j])
	})

	result := prioritized
	if len(others) > 0 {
		result = append(result, "---")
		result = append(result, others...)
	}

	return result
}

type emailMonitor struct {
	client     *imapClient.Client
	mailbox    string
	previous   uint32
	cancelChan chan struct{}
}

var activeMonitors = make(map[string]*emailMonitor)
var monitorsMu sync.Mutex

func startMonitoringCmd(client *imapClient.Client, mailbox string, interval time.Duration) tea.Cmd {
	monitorsMu.Lock()
	defer monitorsMu.Unlock()

	key := mailbox

	if monitor, exists := activeMonitors[key]; exists {
		if monitor.cancelChan != nil {
			close(monitor.cancelChan)
		}
	}

	monitor := &emailMonitor{
		client:     client,
		mailbox:    mailbox,
		cancelChan: make(chan struct{}),
	}
	activeMonitors[key] = monitor

	initialCount, err := client.CheckForNewMessages(context.Background(), mailbox)
	if err == nil {
		monitor.previous = initialCount
	}

	return tea.Tick(interval, func(t time.Time) tea.Msg {
		monitorsMu.Lock()
		monitor, exists := activeMonitors[key]
		if !exists {
			monitorsMu.Unlock()
			return nil
		}
		monitorsMu.Unlock()

		select {
		case <-monitor.cancelChan:
			return nil
		default:
			currentCount, err := monitor.client.CheckForNewMessages(context.Background(), monitor.mailbox)
			if err == nil && currentCount > monitor.previous {
				monitor.previous = currentCount
				return NewEmailMsg{Mailbox: monitor.mailbox, Count: currentCount}
			}
			return nil
		}
	})
}

func stopMonitoringCmd(mailbox string) tea.Cmd {
	return func() tea.Msg {
		monitorsMu.Lock()
		defer monitorsMu.Unlock()

		key := mailbox
		if monitor, exists := activeMonitors[key]; exists {
			if monitor.cancelChan != nil {
				close(monitor.cancelChan)
			}
			delete(activeMonitors, key)
		}

		return nil
	}
}
