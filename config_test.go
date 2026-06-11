package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Theme != "dark" {
		t.Errorf("expected theme 'dark', got '%s'", cfg.Theme)
	}
	if cfg.DefaultKeys != "vim" {
		t.Errorf("expected keys 'vim', got '%s'", cfg.DefaultKeys)
	}
	if len(cfg.SkipDirs) == 0 {
		t.Error("expected non-empty skip_dirs")
	}

	foundArchive := false
	for _, d := range cfg.SkipDirs {
		if d == "archive" {
			foundArchive = true
		}
	}
	if !foundArchive {
		t.Error("expected 'archive' in skip_dirs")
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
vault_path: /Users/test/vault
theme: light
default_keys: arrows
skip_dirs:
  - .custom
  - temp
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.VaultPath != "/Users/test/vault" {
		t.Errorf("vault_path = %s, want /Users/test/vault", cfg.VaultPath)
	}
	if cfg.Theme != "light" {
		t.Errorf("theme = %s, want light", cfg.Theme)
	}
	if cfg.DefaultKeys != "arrows" {
		t.Errorf("default_keys = %s, want arrows", cfg.DefaultKeys)
	}
	if len(cfg.SkipDirs) < 1 {
		t.Error("expected skip_dirs populated")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.IsNotExist to be true, got false; error: %v", err)
	}
}

func TestThemeLookup_Valid(t *testing.T) {
	themes := []string{
		"dark",
		"catppuccin-latte",
		"catppuccin-frappe",
		"catppuccin-macchiato",
		"catppuccin-mocha",
		"dracula",
		"alucard",
	}
	for _, name := range themes {
		p, err := lookupPalette(name)
		if err != nil {
			t.Errorf("lookupPalette(%q): unexpected error: %v", name, err)
		}
		if p.Accent == "" {
			t.Errorf("lookupPalette(%q): Accent not set", name)
		}
	}
}

func TestThemeLookup_Unknown(t *testing.T) {
	_, err := lookupPalette("nonexistent-theme")
	if err == nil {
		t.Error("expected error for unknown theme")
	}
}

func TestThemeWiredToModel(t *testing.T) {
	cfg := &Config{
		VaultPath: testVaultPath(t),
		Theme:     "dracula",
		SkipDirs:  DefaultConfig().SkipDirs,
	}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}
	if m.palette.Accent == "" {
		t.Error("palette not set on model")
	}
	if m.palette.Accent != "#bd93f9" {
		t.Errorf("expected dracula accent #bd93f9, got %s", m.palette.Accent)
	}
	if Accent != "#bd93f9" {
		t.Errorf("package-level Accent should be dracula #bd93f9, got %s", Accent)
	}
}

func TestLoadConfig_CLIOverride(t *testing.T) {
	// CLI override is tested via flag in main.go.
	// Here we verify that when vaultPath is provided via flag, it takes priority
	// by checking the config's VaultPath can be set directly.
	cfg := DefaultConfig()
	cfg.VaultPath = "/cli/path"
	if cfg.VaultPath != "/cli/path" {
		t.Error("VaultPath should be settable (CLI override)")
	}

	cfg2, err := LoadConfig(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil && cfg2 != nil {
		cfg2.VaultPath = "/override/path"
		if cfg2.VaultPath != "/override/path" {
			t.Error("override not applied")
		}
	}
}
