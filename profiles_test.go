package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseNestedMap(t *testing.T) {
	yaml := `
profiles:
  work:
    path: /path/to/work
    theme: dracula
  personal:
    path: /path/to/personal
    theme: catppuccin-mocha
`
	result := parseNestedMap([]byte(yaml), "profiles")

	if len(result) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(result))
	}

	work, ok := result["work"]
	if !ok {
		t.Fatal("expected 'work' profile")
	}
	if work["path"] != "/path/to/work" {
		t.Errorf("work path = %q, want /path/to/work", work["path"])
	}
	if work["theme"] != "dracula" {
		t.Errorf("work theme = %q, want dracula", work["theme"])
	}

	personal, ok := result["personal"]
	if !ok {
		t.Fatal("expected 'personal' profile")
	}
	if personal["path"] != "/path/to/personal" {
		t.Errorf("personal path = %q, want /path/to/personal", personal["path"])
	}
}

func TestParseNestedMap_NoProfiles(t *testing.T) {
	yaml := `
vault_path: /some/path
theme: dark
`
	result := parseNestedMap([]byte(yaml), "profiles")
	if len(result) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(result))
	}
}

func TestLoadConfig_WithProfiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
vault_path: /default/vault
theme: dark
profiles:
  work:
    path: /work/vault
    theme: dracula
  personal:
    path: /personal/vault
    theme: catppuccin-mocha
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if cfg.VaultPath != "/default/vault" {
		t.Errorf("VaultPath = %q, want /default/vault", cfg.VaultPath)
	}

	if len(cfg.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(cfg.Profiles))
	}

	work, ok := cfg.Profiles["work"]
	if !ok {
		t.Fatal("expected 'work' profile")
	}
	if work.Path != "/work/vault" {
		t.Errorf("work.Path = %q, want /work/vault", work.Path)
	}
	if work.Theme != "dracula" {
		t.Errorf("work.Theme = %q, want dracula", work.Theme)
	}

	personal, ok := cfg.Profiles["personal"]
	if !ok {
		t.Fatal("expected 'personal' profile")
	}
	if personal.Path != "/personal/vault" {
		t.Errorf("personal.Path = %q, want /personal/vault", personal.Path)
	}
}

func TestProfilePicker_New(t *testing.T) {
	profiles := map[string]Profile{
		"work":     {Path: "/work", Theme: "dracula"},
		"personal": {Path: "/personal", Theme: "dark"},
		"archive":  {Path: "/archive", Theme: "light"},
	}

	pp := NewProfilePicker(profiles)

	if pp.Count() != 3 {
		t.Errorf("Count() = %d, want 3", pp.Count())
	}

	// Profiles should be sorted
	if pp.profiles[0] != "archive" {
		t.Errorf("profiles[0] = %q, want archive", pp.profiles[0])
	}
	if pp.profiles[1] != "personal" {
		t.Errorf("profiles[1] = %q, want personal", pp.profiles[1])
	}
	if pp.profiles[2] != "work" {
		t.Errorf("profiles[2] = %q, want work", pp.profiles[2])
	}
}

func TestProfilePicker_Navigation(t *testing.T) {
	profiles := map[string]Profile{
		"a": {Path: "/a"},
		"b": {Path: "/b"},
		"c": {Path: "/c"},
	}

	pp := NewProfilePicker(profiles)

	if pp.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", pp.cursor)
	}

	pp.MoveDown()
	if pp.cursor != 1 {
		t.Errorf("after MoveDown cursor = %d, want 1", pp.cursor)
	}

	pp.MoveDown()
	pp.MoveDown()
	if pp.cursor != 2 {
		t.Errorf("after 2x MoveDown cursor = %d, want 2", pp.cursor)
	}

	// Should clamp at bottom
	pp.MoveDown()
	if pp.cursor != 2 {
		t.Errorf("after extra MoveDown cursor = %d, want 2 (clamped)", pp.cursor)
	}

	pp.MoveUp()
	if pp.cursor != 1 {
		t.Errorf("after MoveUp cursor = %d, want 1", pp.cursor)
	}

	// Go to top and try to go above
	pp.MoveUp()
	pp.MoveUp()
	if pp.cursor != 0 {
		t.Errorf("after 2x MoveUp cursor = %d, want 0", pp.cursor)
	}
	pp.MoveUp()
	if pp.cursor != 0 {
		t.Errorf("after extra MoveUp cursor = %d, want 0 (clamped)", pp.cursor)
	}
}

func TestProfilePicker_Selected(t *testing.T) {
	profiles := map[string]Profile{
		"work":     {Path: "/work"},
		"personal": {Path: "/personal"},
	}

	pp := NewProfilePicker(profiles)

	// Should return first profile (sorted)
	selected := pp.Selected()
	if selected != "personal" {
		t.Errorf("Selected() = %q, want personal", selected)
	}

	pp.MoveDown()
	selected = pp.Selected()
	if selected != "work" {
		t.Errorf("after MoveDown Selected() = %q, want work", selected)
	}
}

func TestNewModel_ProfilePickerMode(t *testing.T) {
	cfg := &Config{
		VaultPath: "", // No vault path
		Profiles: map[string]Profile{
			"work":     {Path: testVaultPath(t), Theme: "dark"},
			"personal": {Path: testVaultPath(t), Theme: "dracula"},
		},
	}

	m := NewModel(cfg)

	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	if m.mode != ModeProfilePicker {
		t.Errorf("mode = %v, want ModeProfilePicker", m.mode)
	}

	if m.profilePicker.Count() != 2 {
		t.Errorf("profilePicker.Count() = %d, want 2", m.profilePicker.Count())
	}
}
