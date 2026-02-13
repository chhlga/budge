package tui

import (
	"context"
	"net"
	"testing"
	"time"

	imapClient "github.com/chhlga/budge/internal/imap"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
)

func TestSearchEmailsCmd_returnsResultsWhenSubjectMatches(t *testing.T) {
	addr, cleanupServer := startIMAPMemServer(t)
	defer cleanupServer()

	client := imapClient.NewClient(&imapClient.Options{
		Host:     addr.IP.String(),
		Port:     addr.Port,
		TLS:      false,
		STARTTLS: false,
		Username: "user",
		Password: "pass",
	})
	defer func() { _ = client.Disconnect() }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect() error: %v", err)
	}
	if err := client.Authenticate(ctx); err != nil {
		t.Fatalf("Authenticate() error: %v", err)
	}

	conn := client.Client()
	if conn == nil {
		t.Fatal("expected IMAP connection")
	}

	appendMessage(t, conn, "INBOX", "Subject: hello\r\nFrom: alice@example.com\r\nTo: bob@example.com\r\n\r\nBody\r\n")

	msg := searchEmailsCmd(client, "INBOX", "hello")()
	loaded, ok := msg.(EmailsLoadedMsg)
	if !ok {
		t.Fatalf("expected EmailsLoadedMsg, got %T", msg)
	}
	if len(loaded.Emails) == 0 {
		t.Fatalf("expected at least one email result")
	}
}

func startIMAPMemServer(t *testing.T) (*net.TCPAddr, func()) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error: %v", err)
	}

	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()
		t.Fatalf("expected TCP listener addr, got %T", ln.Addr())
	}

	memServer := imapmemserver.New()
	user := imapmemserver.NewUser("user", "pass")
	if err := user.Create("INBOX", nil); err != nil {
		_ = ln.Close()
		t.Fatalf("Create(INBOX) error: %v", err)
	}
	memServer.AddUser(user)

	server := imapserver.New(&imapserver.Options{
		NewSession: func(conn *imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return memServer.NewSession(), nil, nil
		},
		Caps: imap.CapSet{
			imap.CapIMAP4rev1: {},
			imap.CapIMAP4rev2: {},
		},
		InsecureAuth: true,
	})

	done := make(chan struct{})
	go func() {
		_ = server.Serve(ln)
		close(done)
	}()

	cleanup := func() {
		_ = server.Close()
		<-done
	}

	return addr, cleanup
}

func appendMessage(t *testing.T, c *imapclient.Client, mailbox string, raw string) {
	t.Helper()

	cmd := c.Append(mailbox, int64(len(raw)), nil)
	if _, err := cmd.Write([]byte(raw)); err != nil {
		_ = cmd.Close()
		t.Fatalf("Append.Write() error: %v", err)
	}
	if err := cmd.Close(); err != nil {
		t.Fatalf("Append.Close() error: %v", err)
	}
	if _, err := cmd.Wait(); err != nil {
		t.Fatalf("Append.Wait() error: %v", err)
	}
}
