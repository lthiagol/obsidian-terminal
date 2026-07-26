package main

import (
	"testing"

	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

func TestExtractHeadings(t *testing.T) {
	content := `# Heading 1
Some text
## Heading 2
More text
### Heading 3
Final text`

	lines := markdown.ParseMarkdown(content)
	headings := markdown.ExtractHeadings(lines)

	if len(headings) != 3 {
		t.Fatalf("expected 3 headings, got %d", len(headings))
	}

	if headings[0].Level != 1 || headings[0].Text != "Heading 1" {
		t.Errorf("heading 0: expected level 1 'Heading 1', got level %d '%s'", headings[0].Level, headings[0].Text)
	}

	if headings[1].Level != 2 || headings[1].Text != "Heading 2" {
		t.Errorf("heading 1: expected level 2 'Heading 2', got level %d '%s'", headings[1].Level, headings[1].Text)
	}

	if headings[2].Level != 3 || headings[2].Text != "Heading 3" {
		t.Errorf("heading 2: expected level 3 'Heading 3', got level %d '%s'", headings[2].Level, headings[2].Text)
	}
}

func TestExtractHeadings_Empty(t *testing.T) {
	content := `Just some text
No headings here`

	lines := markdown.ParseMarkdown(content)
	headings := markdown.ExtractHeadings(lines)

	if len(headings) != 0 {
		t.Errorf("expected 0 headings, got %d", len(headings))
	}
}

func TestBuildOutline(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	note, err := LoadNote(cfg.VaultPath, "index.md")
	if err != nil {
		t.Fatalf("LoadNote: %v", err)
	}

	model.activeNote = note
	model.width = 120
	model.height = 40
	model.treeWidth = 30
	model.viewer.SetContent(note.Body, model.width-model.treeWidth-2)

	model.buildOutline()

	if len(model.outlineItems) == 0 {
		t.Error("expected outline items to be populated")
	}

	if model.outlineCursor != 0 {
		t.Errorf("expected outlineCursor to be 0, got %d", model.outlineCursor)
	}
}

func TestBuildOutline_NoNote(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.activeNote = nil
	model.buildOutline()

	if len(model.outlineItems) != 0 {
		t.Errorf("expected 0 outline items when no active note, got %d", len(model.outlineItems))
	}
}

func TestRenderOutline_Empty(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.outlineItems = nil
	output := model.renderOutline()

	if output == "" {
		t.Error("expected non-empty output for empty outline")
	}
}

func TestRenderOutline_WithItems(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.outlineItems = []OutlineItem{
		{Level: 1, Text: "Heading 1", YOffset: 0},
		{Level: 2, Text: "Heading 2", YOffset: 5},
	}
	model.outlineCursor = 0

	output := model.renderOutline()

	if output == "" {
		t.Error("expected non-empty output for outline with items")
	}
}
