package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestVaultPath_NotExist(t *testing.T) {
	cfg := &Config{VaultPath: "/nonexistent/path/vault", SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)

	if m.err == nil {
		t.Error("expected error for nonexistent vault path")
	}
}

func TestVaultPath_IsFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "not-a-dir.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	cfg := &Config{VaultPath: tmpFile, SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)

	if m.err == nil {
		t.Error("expected error when vault path is a file")
	}
}

func TestLoadNote_MalformedFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "---\ntitle: [bad yaml\n---\n\n# Content\nhello world\n"
	os.WriteFile(filepath.Join(dir, "note.md"), []byte(content), 0644)

	note, err := LoadNote(dir, "note.md")
	if err != nil {
		t.Fatalf("LoadNote should not error on malformed frontmatter: %v", err)
	}

	if note.Title != "Note" {
		t.Errorf("expected title fallback 'Note', got %q", note.Title)
	}
	if len(note.Tags) != 0 {
		t.Errorf("expected no tags from malformed frontmatter, got %v", note.Tags)
	}
}

func TestViewer_EmptyNote(t *testing.T) {
	v := NewViewer(markdownStyleFrom(newDarkPalette(), "compact"))
	v.SetContent("", 80)
	view := v.View()

	// Should not panic
	if view == "" {
		t.Error("viewer should render something for empty content")
	}
}

func TestViewer_SingleLongLine(t *testing.T) {
	longLine := strings.Repeat("x", 5000)
	v := NewViewer(markdownStyleFrom(newDarkPalette(), "compact"))
	v.SetContent(longLine, 60)
	view := v.View()

	// Should not panic
	if len(view) == 0 {
		t.Error("viewer should render long line")
	}

	// Should wrap to multiple lines
	lineCount := strings.Count(view, "\n")
	if lineCount < 10 {
		t.Errorf("long line should wrap to many lines, got %d", lineCount+1)
	}
}

func TestModel_Quit(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	// q key quits
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = model.(Model)
	if !m.quitting {
		t.Error("'q' should set quitting=true")
	}
}

func TestModel_CtrlCQuit(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = model.(Model)
	if !m.quitting {
		t.Error("Ctrl+C should set quitting=true")
	}
}

func TestVaultState_InitialOK(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}
	if m.vaultState != VaultStateOK {
		t.Errorf("expected VaultStateOK, got %d", m.vaultState)
	}
}

func TestVaultState_CheckVaultChanges_Broken(t *testing.T) {
	vaultDir := t.TempDir()
	cfg := &Config{VaultPath: vaultDir, SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	os.RemoveAll(vaultDir)
	m.vaultState = VaultStateOK
	m.checkVaultChanges()
	if m.vaultState != VaultStateBroken {
		t.Error("expected VaultStateBroken after checkVaultChanges on deleted path")
	}
	if len(m.toasts) < 1 {
		t.Error("expected toast warning for broken vault")
	}
}

func TestVaultState_RescanVault_Broken(t *testing.T) {
	vaultDir := t.TempDir()
	cfg := &Config{VaultPath: vaultDir, SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	os.RemoveAll(vaultDir)
	m.rescanVault()
	if m.vaultState != VaultStateBroken {
		t.Error("expected VaultStateBroken after rescanVault on deleted path")
	}
}

func TestBrokenVault_RetryKey(t *testing.T) {
	vaultDir := t.TempDir()
	cfg := &Config{VaultPath: vaultDir, SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	os.RemoveAll(vaultDir)
	m.vaultState = VaultStateBroken
	m.ready = true
	m.width = 80
	m.height = 24

	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := model.(Model)

	if updated.vaultState != VaultStateBroken {
		t.Error("'r' should trigger rescan, vault still broken")
	}
}

func TestBrokenVault_ErrorScreen(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	m.vaultState = VaultStateBroken
	m.ready = true
	m.width = 100
	m.height = 24

	view := m.View()
	if view == "" {
		t.Error("view should not be empty for broken vault")
	}
	if !strings.Contains(view, "inaccessible") && !strings.Contains(view, "Vault") {
		t.Error("view should show broken vault error screen")
	}
}

func TestScanErrors_Display(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	m.scanErrors = []string{"read: bad.md: permission denied", "walk: dir: not a directory"}
	m.scanErrorsVisible = true
	m.ready = true
	m.width = 100
	m.height = 24

	view := m.View()
	if !strings.Contains(view, "Scan Errors") {
		t.Error("view should show scan errors display")
	}
	if !strings.Contains(view, "bad.md") {
		t.Error("view should show scan error details")
	}
}

func TestScanErrors_Dismiss(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}
	m.scanErrorsVisible = true
	m.ready = true
	m.width = 100
	m.height = 24

	// Esc should dismiss scan errors
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := model.(Model)
	if updated.scanErrorsVisible {
		t.Error("Esc should close scan errors display")
	}

	// 'q' should also dismiss
	m.scanErrorsVisible = true
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	updated = model.(Model)
	if updated.scanErrorsVisible {
		t.Error("'q' should close scan errors display when visible")
	}
}

func TestVaultState_Recovery(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	// Simulate broken state then recovery
	m.vaultState = VaultStateBroken
	m.lastRescan = m.lastRescan.Add(-5 * time.Second)

	// checkVaultChanges should detect vault is accessible and recover
	m.checkVaultChanges()
	if m.vaultState != VaultStateOK {
		t.Errorf("expected VaultStateOK after recovery, got %d", m.vaultState)
	}
}

func TestVaultState_CommandPalette_IncludesScanErrors(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	// No scan errors — command not present
	commands := m.registerCommands()
	for _, cmd := range commands {
		if cmd.Name == "Scan Errors" {
			t.Error("Scan Errors command should not appear when there are no scan errors")
		}
	}

	// With scan errors — command present
	m.scanErrors = []string{"read: bad.md: permission denied"}
	commands = m.registerCommands()
	found := false
	for _, cmd := range commands {
		if cmd.Name == "Scan Errors" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Scan Errors command should appear when scan errors exist")
	}
}
