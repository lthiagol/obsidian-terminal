package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds user configuration loaded from YAML.
type Config struct {
	VaultPath        string
	Theme            string
	DefaultKeys      string
	SkipDirs         []string
	DailyNotesDir    string
	DailyNotesFormat string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Theme:            "dark",
		DefaultKeys:      "vim",
		SkipDirs:         []string{".obsidian", ".git", ".trash", "node_modules", "archive"},
		DailyNotesDir:    "Journal",
		DailyNotesFormat: "2006-01-02",
	}
}

// LoadConfig reads and parses a YAML config file at path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	cfg := DefaultConfig()
	parseConfigYAML(data, cfg)

	return cfg, nil
}

func parseConfigYAML(data []byte, cfg *Config) {
	scanYAML(data, func(key, value string, items []string) {
		switch key {
		case "vault_path":
			if value != "" {
				cfg.VaultPath = value
			}
		case "theme":
			if value != "" {
				cfg.Theme = value
			}
		case "default_keys":
			if value != "" {
				cfg.DefaultKeys = value
			}
		case "skip_dirs":
			if len(items) > 0 {
				cfg.SkipDirs = items
			} else if value != "" {
				cfg.SkipDirs = []string{value}
			}
		case "daily_notes_dir":
			if value != "" {
				cfg.DailyNotesDir = value
			}
		case "daily_notes_format":
			if value != "" {
				cfg.DailyNotesFormat = value
			}
		}
	})
}

func configPathOrDefault(explicit string) string {
	if explicit != "" {
		return explicit
	}
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "obsidian-terminal", "config.yaml")
}
