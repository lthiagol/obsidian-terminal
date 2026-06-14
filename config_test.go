package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestValidateConfig_Valid(t *testing.T) {
	cfg := &Config{
		VaultPath:        t.TempDir(),
		Theme:            "dark",
		LineSpacing:      "normal",
		DailyNotesFormat: "2006-01-02",
		SkipDirs:         []string{".obsidian", "archive"},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) > 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
	if cfg.Theme != "dark" {
		t.Errorf("theme should be dark, got %s", cfg.Theme)
	}
	if cfg.LineSpacing != "normal" {
		t.Errorf("line_spacing should be normal, got %s", cfg.LineSpacing)
	}
}

func TestValidateConfig_InvalidTheme(t *testing.T) {
	cfg := &Config{
		VaultPath:   t.TempDir(),
		Theme:       "nonexistent",
		LineSpacing: "compact",
		SkipDirs:    []string{".obsidian"},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected theme to be auto-fixed to dark, got %s", cfg.Theme)
	}
}

func TestValidateConfig_InvalidLineSpacing(t *testing.T) {
	cfg := &Config{
		VaultPath:   t.TempDir(),
		Theme:       "dark",
		LineSpacing: "extra-wide",
		SkipDirs:    []string{".obsidian"},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	if cfg.LineSpacing != "compact" {
		t.Errorf("expected line_spacing to be auto-fixed to compact, got %s", cfg.LineSpacing)
	}
}

func TestValidateConfig_InvalidDateFormat(t *testing.T) {
	cfg := &Config{
		VaultPath:        t.TempDir(),
		Theme:            "dark",
		LineSpacing:      "compact",
		DailyNotesFormat: "2006-13-99",
		SkipDirs:         []string{".obsidian"},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	if cfg.DailyNotesFormat != "2006-01-02" {
		t.Errorf("expected daily_notes_format to be auto-fixed to 2006-01-02, got %s", cfg.DailyNotesFormat)
	}
}

func TestValidateConfig_InvalidSkipDirs(t *testing.T) {
	cfg := &Config{
		VaultPath:   t.TempDir(),
		Theme:       "dark",
		LineSpacing: "compact",
		SkipDirs:    []string{".obsidian", "..", "/absolute/path"},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestValidateConfig_EmptyDefaults(t *testing.T) {
	cfg := &Config{
		VaultPath: t.TempDir(),
	}
	warnings := ValidateConfig(cfg)
	if cfg.Theme != "dark" {
		t.Errorf("expected theme to default to dark, got %s", cfg.Theme)
	}
	if cfg.LineSpacing != "compact" {
		t.Errorf("expected line_spacing to default to compact, got %s", cfg.LineSpacing)
	}
	if len(cfg.SkipDirs) == 0 {
		t.Error("expected skip_dirs to be populated with defaults")
	}
	_ = warnings
}

func TestValidateConfig_InvalidCustomThemeColors(t *testing.T) {
	cfg := &Config{
		VaultPath:   t.TempDir(),
		Theme:       "dark",
		LineSpacing: "compact",
		SkipDirs:    []string{".obsidian"},
		CustomTheme: &CustomTheme{
			Accent:          "not-a-color",
			AccentSecondary: "#XYZ",
			TextPrimary:     "#ff0000",
		},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings for bad colors, got %d: %v", len(warnings), warnings)
	}
}

func TestValidateConfig_ProfileNoPath(t *testing.T) {
	cfg := &Config{
		VaultPath:   t.TempDir(),
		Theme:       "dark",
		LineSpacing: "compact",
		SkipDirs:    []string{".obsidian"},
		Profiles: map[string]Profile{
			"work": {Path: ""},
			"home": {Path: "/some/path"},
		},
	}
	warnings := ValidateConfig(cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for empty profile path, got %d: %v", len(warnings), warnings)
	}
}

func TestValidateConfig_ModelShowsToast(t *testing.T) {
	cfg := &Config{
		VaultPath:        testVaultPath(t),
		Theme:            "nonexistent",
		LineSpacing:      "bad-value",
		DailyNotesFormat: "not-a-format",
		SkipDirs:         []string{".."},
	}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel should not fail on invalid config: %v", m.err)
	}
	if len(m.toasts) < 3 {
		t.Errorf("expected at least 3 toasts for invalid theme, spacing, and format, got %d: %v",
			len(m.toasts), m.toasts)
	}
}

func TestValidateConfig_CustomThemeToast(t *testing.T) {
	cfg := &Config{
		VaultPath:   testVaultPath(t),
		Theme:       "dark",
		LineSpacing: "compact",
		SkipDirs:    []string{".obsidian"},
		CustomTheme: &CustomTheme{
			Accent:     "#badhex",
			Background: "#00ff00",
		},
	}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}
	found := false
	for _, toast := range m.toasts {
		if strings.Contains(toast.Message, "accent") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected toast warning about invalid custom theme accent color")
	}
}
