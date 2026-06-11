package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWatcher_DetectsNewFile(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	// Count initial files
	initialCount := countFiles(m.vault)

	// Create a new file
	newFile := filepath.Join(testVaultPath(t), "new-test-file.md")
	err := os.WriteFile(newFile, []byte("# Test\nhello"), 0644)
	if err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	defer os.Remove(newFile)

	// Touch the vault root dir to trigger modtime change
	now := time.Now()
	os.Chtimes(testVaultPath(t), now, now)

	// Trigger rescan via Ctrl+R
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	m = model.(Model)

	newCount := countFiles(m.vault)

	if newCount <= initialCount {
		t.Errorf("expected more files after rescan, got %d (was %d)", newCount, initialCount)
	}

	// Verify new file is in search index
	if _, ok := m.searchIndex["new-test-file.md"]; !ok {
		t.Error("new file should be in search index")
	}
}

func TestWatcher_DetectsDelete(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}

	// Create temp file first
	tmpFile := filepath.Join(testVaultPath(t), "temp-delete-test.md")
	os.WriteFile(tmpFile, []byte("# temp"), 0644)

	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	initialCount := countFiles(m.vault)

	// Delete the file
	os.Remove(tmpFile)
	now := time.Now()
	os.Chtimes(testVaultPath(t), now, now)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	m = model.(Model)

	newCount := countFiles(m.vault)

	if newCount >= initialCount {
		t.Errorf("expected fewer files after delete, got %d (was %d)", newCount, initialCount)
	}
}

func TestWikiLinkResolution_ExactPath(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	vault, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault: %v", err)
	}

	resolved := ResolveWikiLink("projects/api-design", vault, testVaultPath(t))
	if resolved != "projects/api-design.md" {
		t.Errorf("exact path: %q, want 'projects/api-design.md'", resolved)
	}
}

func TestWikiLinkResolution_Basename(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	vault, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault: %v", err)
	}

	resolved := ResolveWikiLink("database", vault, testVaultPath(t))
	if resolved != "projects/database.md" {
		t.Errorf("basename: %q, want 'projects/database.md'", resolved)
	}
}

func TestWikiLinkResolution_CaseInsensitive(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	vault, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault: %v", err)
	}

	resolved := ResolveWikiLink("DATABASE", vault, testVaultPath(t))
	if resolved != "projects/database.md" {
		t.Errorf("case-insensitive: %q, want 'projects/database.md'", resolved)
	}
}

func TestWikiLinkResolution_AliasMatch(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	vault, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault: %v", err)
	}

	// index.md has aliases: "Home", "Start Here"
	resolved := ResolveWikiLink("Home", vault, testVaultPath(t))
	if resolved != "index.md" {
		t.Errorf("alias 'Home': %q, want 'index.md'", resolved)
	}

	resolved = ResolveWikiLink("Start Here", vault, testVaultPath(t))
	if resolved != "index.md" {
		t.Errorf("alias 'Start Here': %q, want 'index.md'", resolved)
	}
}

func TestToast_AddAndExpire(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)

	// No toasts initially
	if len(m.toasts) != 0 {
		t.Error("expected no toasts initially")
	}

	m.addToast("test message", ToastInfo)
	if len(m.toasts) != 1 {
		t.Errorf("expected 1 toast, got %d", len(m.toasts))
	}
	if m.toasts[0].Message != "test message" {
		t.Errorf("toast message = %q, want 'test message'", m.toasts[0].Message)
	}
	if m.toasts[0].Type != ToastInfo {
		t.Errorf("toast type = %v, want ToastInfo", m.toasts[0].Type)
	}

	// Expire toasts - should still be there (just created)
	m.expireToasts()
	if len(m.toasts) != 1 {
		t.Error("toast should not be expired immediately")
	}

	// Manually age the toast
	m.toasts[0].Created = time.Now().Add(-4 * time.Second)
	m.expireToasts()
	if len(m.toasts) != 0 {
		t.Error("toast should be expired after TTL")
	}
}

func TestRenderToast(t *testing.T) {
	toast := Toast{
		Message: "Test toast message",
		Type:    ToastWarning,
		Created: time.Now(),
		TTL:     3 * time.Second,
	}

	rendered := renderToast(toast, 80)
	if !strings.Contains(rendered, "Test toast message") {
		t.Error("rendered toast should contain the message")
	}
}
