package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds user configuration loaded from YAML.
type Config struct {
	VaultPath   string   `yaml:"vault_path"`
	Theme       string   `yaml:"theme"`
	DefaultKeys string   `yaml:"default_keys"`
	SkipDirs    []string `yaml:"skip_dirs"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Theme:       "dark",
		DefaultKeys: "vim",
		SkipDirs:    []string{".obsidian", ".git", ".trash", "node_modules", "archive"},
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
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func configPathOrDefault(explicit string) string {
	if explicit != "" {
		return explicit
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "obsidian-terminal", "config.yaml")
}
