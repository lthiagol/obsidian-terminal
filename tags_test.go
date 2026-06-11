package main

import (
	"testing"
)

func TestTagList_Empty(t *testing.T) {
	tl := NewTagList(map[string][]string{})
	if tl.Count() != 0 {
		t.Errorf("expected 0 tags, got %d", tl.Count())
	}
}

func TestTagList_SortedByCount(t *testing.T) {
	index := map[string][]string{
		"low":    {"a.md"},
		"high":   {"a.md", "b.md", "c.md"},
		"medium": {"a.md", "b.md"},
	}
	tl := NewTagList(index)
	if tl.Count() != 3 {
		t.Fatalf("expected 3 tags, got %d", tl.Count())
	}
	if tl.entries[0].Name != "high" {
		t.Errorf("first tag = %q, want 'high' (highest count)", tl.entries[0].Name)
	}
	if tl.entries[2].Name != "low" {
		t.Errorf("last tag = %q, want 'low' (lowest count)", tl.entries[2].Name)
	}
}

func TestTagList_Navigation(t *testing.T) {
	index := map[string][]string{
		"a": {"1.md"},
		"b": {"2.md"},
		"c": {"3.md"},
	}
	tl := NewTagList(index)

	tl.MoveDown()
	if tl.SelectedTag() == "" {
		t.Error("selected tag should not be empty after move down")
	}

	tl.MoveDown()
	tl.MoveDown()
	tl.MoveDown()
	prev := tl.SelectedTag()
	tl.MoveDown()
	if tl.SelectedTag() != prev {
		t.Error("should be clamped at bottom")
	}

	tl.MoveUp()
	if tl.SelectedTag() == prev {
		t.Error("should have moved up")
	}
}

func TestTagList_SelectedFiles(t *testing.T) {
	index := map[string][]string{
		"test": {"a.md", "b.md"},
	}
	tl := NewTagList(index)
	files := tl.SelectedFiles()
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}
