package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	configData := `
server:
  host: imap.example.com
  port: 993
  tls: true
credentials:
  username: user@example.com
  password: "secret123"
behavior:
  default_folder: INBOX
  page_size: 50
display:
  date_format: "Jan 02 15:04"
  theme: auto
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(configData), 0600); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Validate server config
	if cfg.Server.Host != "imap.example.com" {
		t.Errorf("Expected host 'imap.example.com', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 993 {
		t.Errorf("Expected port 993, got %d", cfg.Server.Port)
	}
	if !cfg.Server.TLS {
		t.Error("Expected TLS to be true")
	}

	// Validate credentials
	if cfg.Credentials.Username != "user@example.com" {
		t.Errorf("Expected username 'user@example.com', got '%s'", cfg.Credentials.Username)
	}
	if cfg.Credentials.Password != "secret123" {
		t.Errorf("Expected password 'secret123', got '%s'", cfg.Credentials.Password)
	}

	// Validate behavior
	if cfg.Behavior.DefaultFolder != "INBOX" {
		t.Errorf("Expected default folder 'INBOX', got '%s'", cfg.Behavior.DefaultFolder)
	}
	if cfg.Behavior.PageSize != 50 {
		t.Errorf("Expected page size 50, got %d", cfg.Behavior.PageSize)
	}

	// Validate display
	if cfg.Display.DateFormat != "Jan 02 15:04" {
		t.Errorf("Expected date format 'Jan 02 15:04', got '%s'", cfg.Display.DateFormat)
	}
	if cfg.Display.Theme != "auto" {
		t.Errorf("Expected theme 'auto', got '%s'", cfg.Display.Theme)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	invalidYAML := `
server:
  host: imap.example.com
  invalid yaml here: [unclosed
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0600); err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "",
			Port: 993,
			TLS:  true,
		},
		Credentials: CredentialsConfig{
			Username: "user@example.com",
			Password: "secret",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for empty host, got nil")
	}
}

func TestValidate_EmptyUsername(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "imap.example.com",
			Port: 993,
			TLS:  true,
		},
		Credentials: CredentialsConfig{
			Username: "",
			Password: "secret",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for empty username, got nil")
	}
}

func TestValidate_InvalidPortZero(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "imap.example.com",
			Port: 0,
			TLS:  true,
		},
		Credentials: CredentialsConfig{
			Username: "user@example.com",
			Password: "secret",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for port 0, got nil")
	}
}

func TestValidate_InvalidPortTooHigh(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "imap.example.com",
			Port: 70000,
			TLS:  true,
		},
		Credentials: CredentialsConfig{
			Username: "user@example.com",
			Password: "secret",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for port > 65535, got nil")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "imap.example.com",
			Port: 993,
			TLS:  true,
		},
		Credentials: CredentialsConfig{
			Username: "user@example.com",
			Password: "secret",
		},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

func TestCheckPermissions_TooOpen(t *testing.T) {
	configData := `
server:
  host: imap.example.com
  port: 993
  tls: true
credentials:
  username: user@example.com
  password: "secret"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write with overly permissive mode
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	err := CheckPermissions(configPath)
	if err == nil {
		t.Error("Expected permissions warning for mode 0644, got nil")
	}
}

func TestCheckPermissions_Secure(t *testing.T) {
	configData := `
server:
  host: imap.example.com
  port: 993
  tls: true
credentials:
  username: user@example.com
  password: "secret"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write with secure mode
	if err := os.WriteFile(configPath, []byte(configData), 0600); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	err := CheckPermissions(configPath)
	if err != nil {
		t.Errorf("Expected no permissions error for mode 0600, got: %v", err)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	// Config with minimal required fields
	configData := `
server:
  host: imap.example.com
  port: 993
  tls: true
credentials:
  username: user@example.com
  password: "secret"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(configData), 0600); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check defaults are applied
	if cfg.Behavior.DefaultFolder != "INBOX" {
		t.Errorf("Expected default folder 'INBOX', got '%s'", cfg.Behavior.DefaultFolder)
	}
	if cfg.Behavior.PageSize != 50 {
		t.Errorf("Expected default page size 50, got %d", cfg.Behavior.PageSize)
	}
	if cfg.Display.DateFormat != "Jan 02 15:04" {
		t.Errorf("Expected default date format 'Jan 02 15:04', got '%s'", cfg.Display.DateFormat)
	}
	if cfg.Display.Theme != "auto" {
		t.Errorf("Expected default theme 'auto', got '%s'", cfg.Display.Theme)
	}
}
