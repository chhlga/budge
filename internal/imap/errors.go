package imap

import (
	"errors"
	"fmt"
)

var (
	ErrNotConnected     = errors.New("not connected to IMAP server")
	ErrAlreadyConnected = errors.New("already connected to IMAP server")
	ErrAuthFailed       = errors.New("authentication failed")
	ErrConnectionLost   = errors.New("connection to IMAP server lost")
	ErrTimeout          = errors.New("operation timed out")
)

type ConnectionError struct {
	Op  string
	Err error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("imap connection error during %s: %v", e.Op, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}

type AuthenticationError struct {
	Username string
	Err      error
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("imap authentication failed for user %s: %v", e.Username, e.Err)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}
