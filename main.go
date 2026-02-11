package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhlga/budge/internal/config"
	"github.com/chhlga/budge/internal/imap"
	"github.com/chhlga/budge/internal/tui"
)

func main() {
	// Determine config path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".config", "budge", "config.yaml")

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config from %s: %v", configPath, err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// Check config file permissions
	if err := config.CheckPermissions(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run: chmod 600 %s\n\n", configPath)
	}

	// Create IMAP client
	imapOpts := &imap.Options{
		Host:     cfg.Server.Host,
		Port:     cfg.Server.Port,
		TLS:      cfg.Server.TLS,
		STARTTLS: cfg.Server.STARTTLS,
		Username: cfg.Credentials.Username,
		Password: cfg.Credentials.Password,
	}

	client := imap.NewClient(imapOpts)

	// Note: Actual IMAP connection and operations will be done via tea.Cmd
	// to avoid blocking the TUI event loop. For this initial version,
	// we're just creating the client wrapper.

	// Create TUI model
	model := tui.NewModel(cfg, client)

	// Run the TUI
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatalf("error running TUI: %v", err)
	}
}
