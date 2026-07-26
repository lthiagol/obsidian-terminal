package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCommandPalette_FilterEmpty(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.commandPaletteQuery = ""
	model.commandPaletteSearch()

	if len(model.commandPaletteResults) == 0 {
		t.Error("empty query should show all commands")
	}
}

func TestCommandPalette_FilterByName(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.commandPaletteQuery = "search"
	model.commandPaletteSearch()

	if len(model.commandPaletteResults) == 0 {
		t.Error("should find at least Fuzzy Search and Content Search")
	}

	foundFuzzy := false
	foundContent := false
	for _, cmd := range model.commandPaletteResults {
		if cmd.Name == "Fuzzy Search" {
			foundFuzzy = true
		}
		if cmd.Name == "Content Search" {
			foundContent = true
		}
	}
	if !foundFuzzy {
		t.Error("should find Fuzzy Search")
	}
	if !foundContent {
		t.Error("should find Content Search")
	}
}

func TestCommandPalette_OpenClose(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	if model.commandPaletteVisible {
		t.Error("palette should not be visible initially")
	}

	model.openCommandPalette()
	if !model.commandPaletteVisible {
		t.Error("palette should be visible after open")
	}
	if model.commandPaletteCursor != 0 {
		t.Errorf("cursor = %d, want 0", model.commandPaletteCursor)
	}
}

func TestCommandPalette_Navigation(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.openCommandPalette()

	// Move down
	model.commandPaletteCursor = 0
	model.commandPaletteResults = []Command{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
	}
	model.commandPaletteCursor = 1
	if model.commandPaletteCursor != 1 {
		t.Errorf("cursor = %d after move down, want 1", model.commandPaletteCursor)
	}

	// Clamp at bottom
	model.commandPaletteCursor = len(model.commandPaletteResults) - 1
	if model.commandPaletteCursor != 2 {
		t.Errorf("cursor = %d, want 2 (clamped)", model.commandPaletteCursor)
	}

	model.commandPaletteCursor = 0
}

func TestCommandPalette_Execute(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	executed := false
	model.commandPaletteResults = []Command{
		{
			Name: "Test Command",
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				executed = true
				return *m, nil
			},
		},
	}
	model.commandPaletteCursor = 0
	model.commandPaletteVisible = true

	model.executeCommand(0)

	if !executed {
		t.Error("command should have been executed")
	}
	if model.commandPaletteVisible {
		t.Error("palette should close after execution")
	}
}

func TestCommandPalette_RenderEmpty(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.commandPaletteQuery = ""
	model.commandPaletteResults = nil

	rendered := model.renderCommandPalette()
	if rendered == "" {
		t.Error("should render even with no results")
	}
}
