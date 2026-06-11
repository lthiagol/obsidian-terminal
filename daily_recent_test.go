package main

import (
	"testing"
	"time"
)

func TestBuildDailyNotePath(t *testing.T) {
	cfg := &Config{
		VaultPath:        testVaultPath(t),
		SkipDirs:         DefaultConfig().SkipDirs,
		DailyNotesDir:    "Journal",
		DailyNotesFormat: "2006-01-02",
	}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	path := model.buildDailyNotePath()
	expected := "Journal/" + time.Now().Format("2006-01-02") + ".md"
	if path != expected {
		t.Errorf("expected path %q, got %q", expected, path)
	}
}

func TestAddRecentNote(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.addRecentNote("note1.md")
	model.addRecentNote("note2.md")

	if len(model.recentNotes) != 2 {
		t.Fatalf("expected 2 recent notes, got %d", len(model.recentNotes))
	}

	if model.recentNotes[0] != "note2.md" {
		t.Errorf("expected most recent to be note2.md, got %s", model.recentNotes[0])
	}
}

func TestAddRecentNote_Dedup(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.addRecentNote("note1.md")
	model.addRecentNote("note2.md")
	model.addRecentNote("note1.md")

	if len(model.recentNotes) != 2 {
		t.Errorf("expected 2 recent notes after dedup, got %d", len(model.recentNotes))
	}

	if model.recentNotes[0] != "note1.md" {
		t.Errorf("expected most recent to be note1.md, got %s", model.recentNotes[0])
	}
}

func TestAddRecentNote_Cap50(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	for i := 0; i < 60; i++ {
		model.addRecentNote("note" + string(rune('0'+i%10)) + ".md")
	}

	if len(model.recentNotes) > 50 {
		t.Errorf("expected max 50 recent notes, got %d", len(model.recentNotes))
	}
}

func TestAddRecentNote_EmptyPath(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.addRecentNote("")
	if len(model.recentNotes) != 0 {
		t.Errorf("expected 0 recent notes for empty path, got %d", len(model.recentNotes))
	}
}

func TestToggleRecents(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	if model.recentVisible {
		t.Error("expected recentVisible to be false initially")
	}

	model.toggleRecents()
	if !model.recentVisible {
		t.Error("expected recentVisible to be true after toggle")
	}

	model.toggleRecents()
	if model.recentVisible {
		t.Error("expected recentVisible to be false after second toggle")
	}
}

func TestRenderRecents_Empty(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.recentNotes = nil
	output := model.renderRecents()

	if output == "" {
		t.Error("expected non-empty output for empty recents")
	}
}

func TestRenderRecents_WithItems(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.recentNotes = []string{"note1.md", "note2.md"}
	model.recentCursor = 0

	output := model.renderRecents()

	if output == "" {
		t.Error("expected non-empty output for recents with items")
	}
}
