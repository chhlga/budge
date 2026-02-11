package imap

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/emersion/go-imap/v2/imapclient"
)

type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateAuthenticated
)

type Client struct {
	mu     sync.RWMutex
	client *imapclient.Client
	state  ConnectionState
	opts   *Options
}

type Options struct {
	Host                 string
	Port                 int
	TLS                  bool
	STARTTLS             bool
	Username             string
	Password             string
	MaxReconnectAttempts int
	InitialBackoff       time.Duration
	MaxBackoff           time.Duration
}

func NewClient(opts *Options) *Client {
	if opts.MaxReconnectAttempts == 0 {
		opts.MaxReconnectAttempts = 5
	}
	if opts.InitialBackoff == 0 {
		opts.InitialBackoff = 1 * time.Second
	}
	if opts.MaxBackoff == 0 {
		opts.MaxBackoff = 30 * time.Second
	}

	return &Client{
		opts:  opts,
		state: StateDisconnected,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state != StateDisconnected {
		return ErrAlreadyConnected
	}

	c.state = StateConnecting

	addr := fmt.Sprintf("%s:%d", c.opts.Host, c.opts.Port)

	var client *imapclient.Client
	var err error

	dialOpts := &imapclient.Options{
		TLSConfig: &tls.Config{
			ServerName: c.opts.Host,
		},
	}

	if c.opts.TLS {
		client, err = imapclient.DialTLS(addr, dialOpts)
	} else if c.opts.STARTTLS {
		client, err = imapclient.DialStartTLS(addr, dialOpts)
	} else {
		client, err = imapclient.DialInsecure(addr, dialOpts)
	}

	if err != nil {
		c.state = StateDisconnected
		return &ConnectionError{Op: "dial", Err: err}
	}

	c.client = client
	c.state = StateConnected

	return nil
}

func (c *Client) Authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state != StateConnected {
		return ErrNotConnected
	}

	if err := c.client.Login(c.opts.Username, c.opts.Password).Wait(); err != nil {
		return &AuthenticationError{
			Username: c.opts.Username,
			Err:      err,
		}
	}

	c.state = StateAuthenticated
	return nil
}

func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateDisconnected {
		return nil
	}

	if c.client != nil {
		if err := c.client.Logout().Wait(); err != nil {
			c.client.Close()
		}
		c.client = nil
	}

	c.state = StateDisconnected
	return nil
}

func (c *Client) Reconnect(ctx context.Context) error {
	c.Disconnect()

	for attempt := 0; attempt < c.opts.MaxReconnectAttempts; attempt++ {
		if err := c.Connect(ctx); err != nil {
			backoff := c.calculateBackoff(attempt)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				continue
			}
		}

		if err := c.Authenticate(ctx); err != nil {
			c.Disconnect()
			backoff := c.calculateBackoff(attempt)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				continue
			}
		}

		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts", c.opts.MaxReconnectAttempts)
}

func (c *Client) calculateBackoff(attempt int) time.Duration {
	backoff := float64(c.opts.InitialBackoff) * math.Pow(2, float64(attempt))
	if backoff > float64(c.opts.MaxBackoff) {
		return c.opts.MaxBackoff
	}
	return time.Duration(backoff)
}

func (c *Client) State() ConnectionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state == StateConnected || c.state == StateAuthenticated
}

// Client returns the underlying IMAP client for direct operations
func (c *Client) Client() *imapclient.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client
}
