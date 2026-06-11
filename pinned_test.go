package main

import (
	"testing"
)

func TestTogglePin_PinAndUnpin(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	path := "index.md"
	model.togglePin(path)

	if len(model.pinnedNotes) != 1 {
		t.Errorf("expected 1 pinned note, got %d", len(model.pinnedNotes))
	}
	if model.pinnedNotes[0].Path != path {
		t.Errorf("expected pinned path %q, got %q", path, model.pinnedNotes[0].Path)
	}

	model.togglePin(path)
	if len(model.pinnedNotes) != 0 {
		t.Errorf("expected 0 pinned notes after unpin, got %d", len(model.pinnedNotes))
	}
}

func TestTogglePin_EmptyPath(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.togglePin("")
	if len(model.pinnedNotes) != 0 {
		t.Errorf("expected 0 pinned notes for empty path, got %d", len(model.pinnedNotes))
	}
}

func TestCyclePinnedNext_NoPins(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.cyclePinnedNext()
	if len(model.toasts) == 0 {
		t.Error("expected warning toast when cycling with no pins")
	}
}

func TestCyclePinnedNext_WrapsAround(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.togglePin("index.md")
	model.togglePin("notes/no-frontmatter.md")

	model.cyclePinnedNext()
	if model.activePinnedIdx != 0 {
		t.Errorf("expected activePinnedIdx 0, got %d", model.activePinnedIdx)
	}

	model.cyclePinnedNext()
	if model.activePinnedIdx != 1 {
		t.Errorf("expected activePinnedIdx 1, got %d", model.activePinnedIdx)
	}

	model.cyclePinnedNext()
	if model.activePinnedIdx != 0 {
		t.Errorf("expected activePinnedIdx to wrap to 0, got %d", model.activePinnedIdx)
	}
}

func TestCyclePinnedPrev_WrapsAround(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.togglePin("index.md")
	model.togglePin("notes/no-frontmatter.md")

	model.cyclePinnedPrev()
	if model.activePinnedIdx != 1 {
		t.Errorf("expected activePinnedIdx to wrap to 1, got %d", model.activePinnedIdx)
	}

	model.cyclePinnedPrev()
	if model.activePinnedIdx != 0 {
		t.Errorf("expected activePinnedIdx 0, got %d", model.activePinnedIdx)
	}
}

func TestValidatePins_RemovesDeleted(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.togglePin("index.md")
	model.togglePin("nonexistent.md")

	model.validatePins()

	if len(model.pinnedNotes) != 1 {
		t.Errorf("expected 1 valid pin after validation, got %d", len(model.pinnedNotes))
	}
	if model.pinnedNotes[0].Path != "index.md" {
		t.Errorf("expected valid pin to be index.md, got %q", model.pinnedNotes[0].Path)
	}
}
