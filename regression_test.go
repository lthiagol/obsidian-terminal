package main

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRegression_ProfileSwitch(t *testing.T) {
	vaultPath := testVaultPath(t)
	secondVault := t.TempDir()

	err := os.WriteFile(filepath.Join(secondVault, "test.md"), []byte("# Second Vault"), 0644)
	if err != nil {
		t.Fatalf("create second vault: %v", err)
	}

	cfg := &Config{
		VaultPath: vaultPath,
		SkipDirs:  DefaultConfig().SkipDirs,
		Profiles: map[string]Profile{
			"default": {Path: vaultPath, Theme: "dark"},
			"second":  {Path: secondVault, Theme: "alucard"},
		},
	}

	m := newTestModel(t, cfg)
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = modelFromInterface(model)

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'P'}})

	if m.mode != ModeProfilePicker {
		t.Fatalf("expected ModeProfilePicker, got %v", m.mode)
	}

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyDown})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if m.mode != ModeBrowse {
		t.Errorf("expected ModeBrowse after profile switch, got %v", m.mode)
	}
	if m.config.VaultPath != secondVault {
		t.Errorf("expected vault path %q, got %q", secondVault, m.config.VaultPath)
	}
	if m.palette.Accent == "" {
		t.Error("palette should be set after profile switch")
	}
}

func TestRegression_ThemePalette(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs, Theme: "dracula"}
	m := newTestModel(t, cfg)

	if m.palette.Accent != "#bd93f9" {
		t.Errorf("expected dracula accent #bd93f9 from m.palette, got %s", m.palette.Accent)
	}
}
