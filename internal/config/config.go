package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds configuration for the VGU MCP server.
type Config struct {
	MoodleURL   string `json:"moodle_url"`
	MoodleToken string `json:"moodle_token"`
	UserID      int    `json:"userid,omitempty"`
}

// Path returns the absolute path to the config file.
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config", "vgu-mcp", "config.json"), nil
}

// Load reads configuration from environment variables first,
// then falls back to ~/.config/vgu-mcp/config.json.
// Environment variables always take precedence over the file.
func Load() (*Config, error) {
	cfg := &Config{
		MoodleURL:   os.Getenv("MOODLE_URL"),
		MoodleToken: os.Getenv("MOODLE_TOKEN"),
	}
	if cfg.MoodleURL != "" && cfg.MoodleToken != "" {
		return cfg, nil
	}

	path, err := Path()
	if err != nil {
		return cfg, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	var fileCfg Config
	if err := json.NewDecoder(f).Decode(&fileCfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.MoodleURL == "" {
		cfg.MoodleURL = fileCfg.MoodleURL
	}
	if cfg.MoodleToken == "" {
		cfg.MoodleToken = fileCfg.MoodleToken
	}
	if cfg.UserID == 0 {
		cfg.UserID = fileCfg.UserID
	}
	return cfg, nil
}

// Save writes cfg to ~/.config/vgu-mcp/config.json, creating parent directories
// as needed. The file is written with 0600 permissions (owner read/write only).
func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open config for writing: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
