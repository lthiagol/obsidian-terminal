package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel(t *testing.T, cfg *Config) Model {
	t.Helper()
	if cfg.VaultPath == "" {
		cfg.VaultPath = testVaultPath(t)
	}
	if cfg.SkipDirs == nil {
		cfg.SkipDirs = DefaultConfig().SkipDirs
	}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}
	return m
}

func sendKey(t *testing.T, m Model, msg tea.KeyMsg) Model {
	t.Helper()
	model, _ := m.Update(msg)
	return modelFromInterface(model)
}

func sendKeys(t *testing.T, m Model, msgs ...tea.KeyMsg) Model {
	t.Helper()
	for _, msg := range msgs {
		m = sendKey(t, m, msg)
	}
	return m
}

func assertMode(t *testing.T, m Model, want Mode) {
	t.Helper()
	if m.mode != want {
		t.Errorf("expected mode %v, got %v", want, m.mode)
	}
}

func assertActiveNotePath(t *testing.T, m Model, suffix string) {
	t.Helper()
	if m.activeNote == nil {
		t.Fatal("activeNote is nil")
	}
	if !strings.HasSuffix(m.activeNote.Path, suffix) {
		t.Errorf("expected activeNote path to end with %q, got %q", suffix, m.activeNote.Path)
	}
}

func indexOfFirstCollapsedDir(ft FileTree) int {
	for i, item := range ft.Items() {
		if item.entry.IsDir && !item.expanded {
			return i
		}
	}
	return -1
}

func indexOfFirstFile(ft FileTree) int {
	for i, item := range ft.Items() {
		if !item.entry.IsDir {
			return i
		}
	}
	return 0
}

func modelFromInterface(v tea.Model) Model {
	switch m := v.(type) {
	case Model:
		return m
	case *Model:
		return *m
	default:
		panic("unexpected tea.Model type")
	}
}

func navigateToFirstFile(t *testing.T, model *tea.Model) Model {
	t.Helper()
	m := modelFromInterface(*model)
	firstFileIdx := indexOfFirstFile(m.fileTree)
	for m.fileTree.Cursor() < firstFileIdx {
		*model, _ = (*model).Update(tea.KeyMsg{Type: tea.KeyDown})
		m = modelFromInterface(*model)
	}
	*model, _ = (*model).Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = modelFromInterface(*model)
	if m.mode != ModeView {
		t.Fatalf("navigateToFirstFile: expected ModeView, got %v", m.mode)
	}
	return m
}
