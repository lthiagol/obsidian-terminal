package main

import (
	"strings"
	"testing"

	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

func TestCheckbox_Unchecked(t *testing.T) {
	content := "- [ ] Unchecked item"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	line := lines[0]
	if !line.Checkable {
		t.Error("expected Checkable to be true")
	}
	if line.Checked {
		t.Error("expected Checked to be false")
	}

	rendered := markdown.RenderMarkdown(lines, 80, defaultRendererStyle())
	if !strings.Contains(rendered, "[ ]") {
		t.Error("rendered output should contain [ ]")
	}
	if strings.Contains(rendered, "•") {
		t.Error("rendered output should not contain bullet for checkbox")
	}
}

func TestCheckbox_Checked(t *testing.T) {
	content := "- [x] Checked item"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	line := lines[0]
	if !line.Checkable {
		t.Error("expected Checkable to be true")
	}
	if !line.Checked {
		t.Error("expected Checked to be true")
	}
}

func TestCheckbox_CheckedUppercase(t *testing.T) {
	content := "- [X] Checked uppercase"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	line := lines[0]
	if !line.Checkable {
		t.Error("expected Checkable to be true for [X]")
	}
	if !line.Checked {
		t.Error("expected Checked to be true for [X]")
	}
}

func TestCheckbox_WithFormatting(t *testing.T) {
	content := "- [x] **bold** done"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	line := lines[0]
	if !line.Checkable {
		t.Error("expected Checkable to be true")
	}

	rendered := markdown.RenderMarkdown(lines, 80, defaultRendererStyle())
	if !strings.Contains(rendered, "[x]") {
		t.Error("rendered output should contain [x]")
	}
}

func TestCheckbox_NormalListItem(t *testing.T) {
	content := "- Normal list item"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	line := lines[0]
	if line.Checkable {
		t.Error("normal list item should not be checkable")
	}

	rendered := markdown.RenderMarkdown(lines, 80, defaultRendererStyle())
	if !strings.Contains(rendered, "•") {
		t.Error("normal list item should have bullet")
	}
}

func TestFrontmatter_Displayed(t *testing.T) {
	content := "---\ntitle: Test\ntags: [a, b]\n---\n\n# Heading\n\nBody text"
	v := NewViewer(markdownStyleFrom(newDarkPalette(), "compact"))
	v.SetContent(content, 80)

	view := v.View()
	if !strings.Contains(view, "Frontmatter") {
		t.Error("view should show frontmatter block")
	}
	if !strings.Contains(view, "title") {
		t.Error("frontmatter should show 'title' key")
	}
	if !strings.Contains(view, "Test") {
		t.Error("frontmatter should show title value")
	}
}

func TestFrontmatter_NoFrontmatter(t *testing.T) {
	content := "# Just a heading\n\nNo frontmatter here"
	v := NewViewer(markdownStyleFrom(newDarkPalette(), "compact"))
	v.SetContent(content, 80)

	view := v.View()
	if strings.Contains(view, "Frontmatter") {
		t.Error("view should not show frontmatter block when none exists")
	}
}

func TestFrontmatter_EmptyBody(t *testing.T) {
	content := "---\ntitle: Metadata Only\n---"
	v := NewViewer(markdownStyleFrom(newDarkPalette(), "compact"))
	v.SetContent(content, 80)

	view := v.View()
	if !strings.Contains(view, "Frontmatter") {
		t.Error("should show frontmatter even with empty body")
	}
	if !strings.Contains(view, "Metadata Only") {
		t.Error("frontmatter should show title value")
	}
}

func defaultRendererStyle() markdown.RendererStyle {
	return markdown.RendererStyle{
		Accent:          "#a78bfa",
		AccentSecondary: "#fbbf24",
		AccentTertiary:  "#2dd4bf",
		TextSecondary:   "#9ca3af",
		TextDim:         "#4b5563",
		Success:         "#34d399",
		CodeBackground:  "#1f2937",
		Heading1:        "#fbbf24",
	}
}
