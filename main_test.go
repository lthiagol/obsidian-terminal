package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("", 80)
	view := v.View()

	// Should not panic
	if view == "" {
		t.Error("viewer should render something for empty content")
	}
}

func TestViewer_SingleLongLine(t *testing.T) {
	longLine := strings.Repeat("x", 5000)
	v := NewViewer(defaultMarkdownStyle())
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
