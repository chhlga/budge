package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Credentials CredentialsConfig `yaml:"credentials"`
	Behavior    BehaviorConfig    `yaml:"behavior"`
	Display     DisplayConfig     `yaml:"display"`
}

// ServerConfig contains IMAP server settings
type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	TLS      bool   `yaml:"tls"`
	STARTTLS bool   `yaml:"starttls"`
}

// CredentialsConfig contains authentication credentials
type CredentialsConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// BehaviorConfig contains application behavior settings
type BehaviorConfig struct {
	DefaultFolder string `yaml:"default_folder"`
	PageSize      int    `yaml:"page_size"`
	PollInterval  int    `yaml:"poll_interval"`
}

// DisplayConfig contains display preferences
type DisplayConfig struct {
	DateFormat string `yaml:"date_format"`
	Theme      string `yaml:"theme"`
}

// Load reads and parses a YAML configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Apply defaults
	if cfg.Behavior.DefaultFolder == "" {
		cfg.Behavior.DefaultFolder = "INBOX"
	}
	if cfg.Behavior.PageSize == 0 {
		cfg.Behavior.PageSize = 50
	}
	if cfg.Behavior.PollInterval == 0 {
		cfg.Behavior.PollInterval = 30
	}
	if cfg.Display.DateFormat == "" {
		cfg.Display.DateFormat = "Jan 02 15:04"
	}
	if cfg.Display.Theme == "" {
		cfg.Display.Theme = "auto"
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Check file permissions
	if err := CheckPermissions(path); err != nil {
		// Log warning but don't fail
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535, got %d", c.Server.Port)
	}

	if c.Server.TLS && c.Server.STARTTLS {
		return fmt.Errorf("cannot enable both TLS and STARTTLS, choose one")
	}

	if c.Credentials.Username == "" {
		return fmt.Errorf("credentials username cannot be empty")
	}

	return nil
}

// CheckPermissions verifies that the config file has secure permissions
func CheckPermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat config file: %w", err)
	}

	mode := info.Mode().Perm()

	// Check if file is readable by group or others (any bit set in 077)
	if mode&0077 != 0 {
		return fmt.Errorf("config file has insecure permissions %o, should be 0600 or similar", mode)
	}

	return nil
}
