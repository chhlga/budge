package imap

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type mockIMAPConn struct {
	connected      bool
	shouldFailConn bool
	shouldFailAuth bool
	callCount      int
}

func (m *mockIMAPConn) connect() error {
	m.callCount++
	if m.shouldFailConn {
		return errors.New("connection failed")
	}
	m.connected = true
	return nil
}

func (m *mockIMAPConn) authenticate(username, password string) error {
	if !m.connected {
		return ErrNotConnected
	}
	if m.shouldFailAuth {
		return errors.New("invalid credentials")
	}
	return nil
}

func (m *mockIMAPConn) close() error {
	m.connected = false
	return nil
}

func TestClient_Connect(t *testing.T) {
	client := &Client{
		state: StateDisconnected,
	}

	if client.IsConnected() {
		t.Error("Expected client to be disconnected initially")
	}

	if client.State() != StateDisconnected {
		t.Errorf("Expected state %v, got %v", StateDisconnected, client.State())
	}
}

func TestClient_State_Transitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  ConnectionState
		expectedState ConnectionState
	}{
		{"disconnected", StateDisconnected, StateDisconnected},
		{"connecting", StateConnecting, StateConnecting},
		{"connected", StateConnected, StateConnected},
		{"authenticated", StateAuthenticated, StateAuthenticated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{state: tt.initialState}
			if client.State() != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, client.State())
			}
		})
	}
}

func TestClient_IsConnected(t *testing.T) {
	tests := []struct {
		state    ConnectionState
		expected bool
	}{
		{StateDisconnected, false},
		{StateConnecting, false},
		{StateConnected, true},
		{StateAuthenticated, true},
	}

	for _, tt := range tests {
		client := &Client{state: tt.state}
		if client.IsConnected() != tt.expected {
			t.Errorf("For state %v, expected IsConnected=%v, got %v",
				tt.state, tt.expected, client.IsConnected())
		}
	}
}

func TestConnectionError(t *testing.T) {
	baseErr := errors.New("network timeout")
	connErr := &ConnectionError{
		Op:  "dial",
		Err: baseErr,
	}

	errMsg := connErr.Error()
	if !strings.Contains(errMsg, "dial") {
		t.Errorf("Error message should contain operation 'dial', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "network timeout") {
		t.Errorf("Error message should contain base error, got: %s", errMsg)
	}

	if unwrapped := connErr.Unwrap(); unwrapped != baseErr {
		t.Errorf("Unwrap should return base error, got: %v", unwrapped)
	}
}

func TestAuthenticationError(t *testing.T) {
	baseErr := errors.New("invalid password")
	authErr := &AuthenticationError{
		Username: "test@example.com",
		Err:      baseErr,
	}

	errMsg := authErr.Error()
	if !strings.Contains(errMsg, "test@example.com") {
		t.Errorf("Error message should contain username, got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "invalid password") {
		t.Errorf("Error message should contain base error, got: %s", errMsg)
	}

	if unwrapped := authErr.Unwrap(); unwrapped != baseErr {
		t.Errorf("Unwrap should return base error, got: %v", unwrapped)
	}
}

func TestClient_ReconnectionAttempts(t *testing.T) {
	opts := &Options{
		MaxReconnectAttempts: 3,
		InitialBackoff:       1,
		MaxBackoff:           5,
	}

	if opts.MaxReconnectAttempts != 3 {
		t.Errorf("Expected 3 reconnection attempts, got %d", opts.MaxReconnectAttempts)
	}
}

func TestClient_BackoffCalculation(t *testing.T) {
	opts := &Options{
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     10 * time.Second,
	}

	tests := []struct {
		attempt  int
		minValue time.Duration
		maxValue time.Duration
	}{
		{0, 1 * time.Second, 1 * time.Second},
		{1, 2 * time.Second, 2 * time.Second},
		{2, 4 * time.Second, 4 * time.Second},
		{3, 8 * time.Second, 8 * time.Second},
		{4, 10 * time.Second, 10 * time.Second},
		{5, 10 * time.Second, 10 * time.Second},
	}

	client := NewClient(opts)

	for _, tt := range tests {
		backoff := client.calculateBackoff(tt.attempt)
		if backoff > opts.MaxBackoff {
			t.Errorf("Backoff for attempt %d should not exceed max %v, got %v",
				tt.attempt, opts.MaxBackoff, backoff)
		}
	}
}

func TestUpdateHandler(t *testing.T) {
	handler := &UpdateHandler{
		Mailbox: "INBOX",
	}

	if handler.Mailbox != "INBOX" {
		t.Errorf("Expected mailbox 'INBOX', got '%s'", handler.Mailbox)
	}

	handler.OnNewMail = func(mailbox string, count uint32) {
		if mailbox != "INBOX" {
			t.Errorf("Expected mailbox 'INBOX', got '%s'", mailbox)
		}
		if count != 5 {
			t.Errorf("Expected count 5, got %d", count)
		}
	}

	handler.OnNewMail("INBOX", 5)
}

func TestClient_UpdateHandler(t *testing.T) {
	client := &Client{state: StateAuthenticated}

	handler := &UpdateHandler{
		Mailbox: "INBOX",
	}

	client.SetUpdateHandler(handler)

	retrieved := client.GetUpdateHandler()
	if retrieved == nil {
		t.Fatal("Expected non-nil handler")
	}
	if retrieved.Mailbox != "INBOX" {
		t.Errorf("Expected mailbox 'INBOX', got '%s'", retrieved.Mailbox)
	}
}

func TestClient_UpdateHandlerNil(t *testing.T) {
	client := &Client{state: StateAuthenticated}

	retrieved := client.GetUpdateHandler()
	if retrieved != nil {
		t.Error("Expected nil handler when not set")
	}
}
