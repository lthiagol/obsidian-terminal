package main

import (
	"testing"
)

func TestBacklinkPanel_Empty(t *testing.T) {
	bp := NewBacklinkPanel("nonexistent.md", map[string][]string{})
	if bp.Count() != 0 {
		t.Errorf("expected 0 backlinks, got %d", bp.Count())
	}
}

func TestBacklinkPanel_WithLinks(t *testing.T) {
	index := map[string][]string{
		"note.md": {"a.md", "b.md", "c.md"},
	}
	bp := NewBacklinkPanel("note.md", index)
	if bp.Count() != 3 {
		t.Errorf("expected 3 backlinks, got %d", bp.Count())
	}
}

func TestBacklinkPanel_Navigation(t *testing.T) {
	index := map[string][]string{
		"note.md": {"a.md", "b.md", "c.md"},
	}
	bp := NewBacklinkPanel("note.md", index)

	if bp.SelectedPath() != "a.md" {
		t.Errorf("initial selection = %q, want 'a.md'", bp.SelectedPath())
	}

	bp.MoveDown()
	if bp.SelectedPath() != "b.md" {
		t.Errorf("after down = %q, want 'b.md'", bp.SelectedPath())
	}

	bp.MoveDown()
	bp.MoveDown()
	if bp.SelectedPath() != "c.md" {
		t.Errorf("clamped at bottom = %q, want 'c.md'", bp.SelectedPath())
	}

	bp.MoveUp()
	if bp.SelectedPath() != "b.md" {
		t.Errorf("after up = %q, want 'b.md'", bp.SelectedPath())
	}
}

func TestBacklinkPanel_Normalization(t *testing.T) {
	index := map[string][]string{
		"note.md": {"source.md"},
	}
	bp := NewBacklinkPanel("Note.md", index)
	if bp.Count() != 1 {
		t.Errorf("expected 1 backlink with normalized path, got %d", bp.Count())
	}
}
