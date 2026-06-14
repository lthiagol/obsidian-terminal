package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"#fff", "#ffffff", false},
		{"#FFF", "#ffffff", false},
		{"#abc", "#aabbcc", false},
		{"#123456", "#123456", false},
		{"#ABCDEF", "#abcdef", false},
		{"#aBcDeF", "#abcdef", false},
		{"", "", true},
		{"fff", "", true},
		{"#ff", "", true},
		{"#ffff", "", true},
		{"#fffff", "", true},
		{"#fffffff", "", true},
		{"#gggggg", "", true},
		{"#xyz", "", true},
	}

	for _, tt := range tests {
		got, err := parseHexColor(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseHexColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && string(got) != tt.want {
			t.Errorf("parseHexColor(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPaletteFromCustom_Nil(t *testing.T) {
	base := newDarkPalette()
	result, err := paletteFromCustom(nil, base)
	if err != nil {
		t.Errorf("paletteFromCustom(nil) error = %v", err)
	}
	if string(result.Accent) != string(base.Accent) {
		t.Errorf("nil custom theme should return base palette")
	}
}

func TestPaletteFromCustom_Empty(t *testing.T) {
	base := newDarkPalette()
	ct := &CustomTheme{}
	result, err := paletteFromCustom(ct, base)
	if err != nil {
		t.Errorf("paletteFromCustom(empty) error = %v", err)
	}
	if string(result.Accent) != string(base.Accent) {
		t.Errorf("empty custom theme should return base palette")
	}
}

func TestPaletteFromCustom_Overrides(t *testing.T) {
	base := newDarkPalette()
	ct := &CustomTheme{
		Accent:      "#ff0000",
		TextPrimary: "#00ff00",
	}
	result, err := paletteFromCustom(ct, base)
	if err != nil {
		t.Errorf("paletteFromCustom error = %v", err)
	}
	if string(result.Accent) != "#ff0000" {
		t.Errorf("Accent = %q, want #ff0000", result.Accent)
	}
	if string(result.TextPrimary) != "#00ff00" {
		t.Errorf("TextPrimary = %q, want #00ff00", result.TextPrimary)
	}
	// Unset fields should keep base values
	if string(result.TextSecondary) != string(base.TextSecondary) {
		t.Errorf("TextSecondary should keep base value")
	}
}

func TestPaletteFromCustom_InvalidHex(t *testing.T) {
	base := newDarkPalette()
	ct := &CustomTheme{
		Accent: "invalid",
	}
	result, err := paletteFromCustom(ct, base)
	if err == nil {
		t.Error("expected error for invalid hex")
	}
	// Should still return a palette (with base values for invalid fields)
	if string(result.Accent) != string(base.Accent) {
		t.Errorf("invalid hex should keep base value")
	}
}

func TestPaletteFromCustom_3DigitHex(t *testing.T) {
	base := newDarkPalette()
	ct := &CustomTheme{
		Accent: "#f00",
	}
	result, err := paletteFromCustom(ct, base)
	if err != nil {
		t.Errorf("paletteFromCustom error = %v", err)
	}
	if string(result.Accent) != "#ff0000" {
		t.Errorf("Accent = %q, want #ff0000", result.Accent)
	}
}

func TestRebuildDerivedStyles(t *testing.T) {
	p := newDarkPalette()
	p.Accent = "#ff0000"
	p.Surface = "#00ff00"

	result := rebuildDerivedStyles(p)

	// Check that derived styles use the new colors
	if string(result.ModeBrowse) != "#ff0000" {
		t.Errorf("ModeBrowse = %q, want #ff0000", result.ModeBrowse)
	}
}

func TestLoadConfig_WithCustomTheme(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
vault_path: /some/vault
theme: dark
custom_theme:
  accent: "#ff0000"
  text_primary: "#00ff00"
  background: "#0000ff"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if cfg.CustomTheme == nil {
		t.Fatal("expected CustomTheme to be set")
	}
	if cfg.CustomTheme.Accent != "#ff0000" {
		t.Errorf("Accent = %q, want #ff0000", cfg.CustomTheme.Accent)
	}
	if cfg.CustomTheme.TextPrimary != "#00ff00" {
		t.Errorf("TextPrimary = %q, want #00ff00", cfg.CustomTheme.TextPrimary)
	}
	if cfg.CustomTheme.Background != "#0000ff" {
		t.Errorf("Background = %q, want #0000ff", cfg.CustomTheme.Background)
	}
}

func TestLoadConfig_EmptyCustomTheme(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
vault_path: /some/vault
custom_theme: {}
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	// Empty custom_theme should not set CustomTheme
	if cfg.CustomTheme != nil {
		t.Error("empty custom_theme should not set CustomTheme")
	}
}

func TestNewModel_WithCustomTheme(t *testing.T) {
	cfg := &Config{
		VaultPath: testVaultPath(t),
		SkipDirs:  DefaultConfig().SkipDirs,
		Theme:     "dark",
		CustomTheme: &CustomTheme{
			Accent: "#ff0000",
		},
	}

	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	if string(m.palette.Accent) != "#ff0000" {
		t.Errorf("palette.Accent = %q, want #ff0000", m.palette.Accent)
	}
}
