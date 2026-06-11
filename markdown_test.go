package main

import (
	"strings"
	"testing"
)

func TestParseMarkdown_Headings(t *testing.T) {
	input := "# Heading 1\nSome text\n## Heading 2\n### Heading 3\n#### Heading 4\n##### Heading 5\n###### Heading 6\n"
	lines := ParseMarkdown(input)

	var headingLines []MarkdownLine
	for _, l := range lines {
		if l.BlockType == BlockHeading {
			headingLines = append(headingLines, l)
		}
	}

	if len(headingLines) != 6 {
		t.Fatalf("expected 6 headings, got %d", len(headingLines))
	}

	expectedLevels := []int{1, 2, 3, 4, 5, 6}
	for i, hl := range headingLines {
		if hl.HeadingLevel != expectedLevels[i] {
			t.Errorf("heading %d: level = %d, want %d", i+1, hl.HeadingLevel, expectedLevels[i])
		}
	}

	// Check h1 text
	h1Segments := headingLines[0].Segments
	if len(h1Segments) == 0 || !strings.Contains(h1Segments[0].Text, "Heading 1") {
		t.Error("h1 should contain 'Heading 1'")
	}
}

func TestParseMarkdown_InlineFormatting(t *testing.T) {
	input := "**bold** and *italic* and ***bold italic*** and `code` and ~~strikethrough~~ and ==highlight==\n"
	lines := ParseMarkdown(input)

	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}

	segments := lines[0].Segments

	hasBold := false
	hasItalic := false
	hasBoldItalic := false
	hasCode := false
	hasStrikethrough := false
	hasHighlight := false

	for _, seg := range segments {
		if seg.Bold && !seg.Italic {
			hasBold = true
		}
		if seg.Italic && !seg.Bold {
			hasItalic = true
		}
		if seg.Bold && seg.Italic {
			hasBoldItalic = true
		}
		if seg.Code {
			hasCode = true
		}
		if seg.Strikethrough {
			hasStrikethrough = true
		}
		if seg.Highlight {
			hasHighlight = true
		}
	}

	if !hasBold {
		t.Error("should have bold segment")
	}
	if !hasItalic {
		t.Error("should have italic segment")
	}
	if !hasBoldItalic {
		t.Error("should have bold+italic segment")
	}
	if !hasCode {
		t.Error("should have code segment")
	}
	if !hasStrikethrough {
		t.Error("should have strikethrough segment")
	}
	if !hasHighlight {
		t.Error("should have highlight segment")
	}
}

func TestParseMarkdown_CodeBlocks(t *testing.T) {
	input := "```go\nfunc main() {\n    fmt.Println(\"hi\")\n}\n```\n"
	lines := ParseMarkdown(input)

	var codeBlocks []MarkdownLine
	for _, l := range lines {
		if l.BlockType == BlockCodeBlock {
			codeBlocks = append(codeBlocks, l)
		}
	}

	if len(codeBlocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(codeBlocks))
	}

	cb := codeBlocks[0]
	if cb.Language != "go" {
		t.Errorf("language = %q, want 'go'", cb.Language)
	}
	if !strings.Contains(cb.RawContent, "fmt.Println") {
		t.Error("code block should contain fmt.Println")
	}

	// Ensure no inline parsing happened inside code block
	if strings.Contains(cb.RawContent, "**") {
		t.Error("code block should contain raw '**' markers")
	}
}

func TestParseMarkdown_WikiLinks(t *testing.T) {
	input := "See [[database]] and [[api-design|API Design]] for details.\n"
	lines := ParseMarkdown(input)

	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}

	var wikiSegments []InlineSegment
	for _, seg := range lines[0].Segments {
		if seg.IsWikiLink {
			wikiSegments = append(wikiSegments, seg)
		}
	}

	if len(wikiSegments) != 2 {
		t.Fatalf("expected 2 wiki-link segments, got %d", len(wikiSegments))
	}

	if wikiSegments[0].WikiTarget != "database" {
		t.Errorf("first wiki target = %q, want 'database'", wikiSegments[0].WikiTarget)
	}
	if wikiSegments[0].WikiDisplay != "database" {
		t.Errorf("first wiki display = %q, want 'database'", wikiSegments[0].WikiDisplay)
	}

	if wikiSegments[1].WikiTarget != "api-design" {
		t.Errorf("second wiki target = %q, want 'api-design'", wikiSegments[1].WikiTarget)
	}
	if wikiSegments[1].WikiDisplay != "API Design" {
		t.Errorf("second wiki display = %q, want 'API Design'", wikiSegments[1].WikiDisplay)
	}
}

func TestParseMarkdown_Callouts(t *testing.T) {
	input := "> [!note] This is a note\n> [!warning] Watch out\n> [!tip] Useful tip\n"
	lines := ParseMarkdown(input)

	calloutCount := 0
	for _, l := range lines {
		if l.BlockType == BlockCallout {
			calloutCount++
			if l.CalloutType != "note" && l.CalloutType != "warning" && l.CalloutType != "tip" {
				t.Errorf("unexpected callout type: %s", l.CalloutType)
			}
		}
	}

	if calloutCount != 3 {
		t.Errorf("expected 3 callouts, got %d", calloutCount)
	}
}

func TestParseMarkdown_Blockquotes(t *testing.T) {
	input := "> This is a blockquote\n> Multiple lines\n"
	lines := ParseMarkdown(input)

	quoteCount := 0
	for _, l := range lines {
		if l.BlockType == BlockBlockquote {
			quoteCount++
		}
	}

	if quoteCount != 2 {
		t.Errorf("expected 2 blockquotes, got %d", quoteCount)
	}
}

func TestParseMarkdown_Lists(t *testing.T) {
	input := "- Item 1\n- Item 2\n  - Nested item\n* Star item\n+ Plus item\n1. Ordered 1\n2. Ordered 2\n"
	lines := ParseMarkdown(input)

	listCount := 0
	for _, l := range lines {
		if l.BlockType == BlockList {
			listCount++
		}
	}

	if listCount < 5 {
		t.Errorf("expected at least 5 list items, got %d", listCount)
	}
}

func TestParseMarkdown_CommentsStripped(t *testing.T) {
	input := "Hello %%hidden comment%% World\n"
	lines := ParseMarkdown(input)

	text := ""
	for _, seg := range lines[0].Segments {
		text += seg.Text
	}

	if strings.Contains(text, "hidden") {
		t.Error("comments should be stripped from output")
	}

	// "Hello World" or "Hello  World" both acceptable
	if !strings.Contains(text, "Hello") {
		t.Error("should contain 'Hello'")
	}
	if !strings.Contains(text, "World") {
		t.Error("should contain 'World'")
	}
}

func TestExtractWikiLinks(t *testing.T) {
	input := "See [[projects/api-design]] and [[database]]. Also [[notes/meeting|meeting notes]].\n"
	lines := ParseMarkdown(input)

	links := ExtractWikiLinks(lines)

	if len(links) != 3 {
		t.Fatalf("expected 3 wiki links, got %d", len(links))
	}

	if links[0].Target != "projects/api-design" {
		t.Errorf("link 0 target = %q", links[0].Target)
	}
	if links[1].Target != "database" {
		t.Errorf("link 1 target = %q", links[1].Target)
	}
	if links[2].Target != "notes/meeting" {
		t.Errorf("link 2 target = %q", links[2].Target)
	}
	if links[2].Display != "meeting notes" {
		t.Errorf("link 2 display = %q", links[2].Display)
	}
}

func TestRenderMarkdown_ANSIContent(t *testing.T) {
	input := "# Hello\n\nThis is **bold** and *italic*.\n"
	lines := ParseMarkdown(input)
	output := RenderMarkdown(lines, 60)

	if !strings.Contains(output, "Hello") {
		t.Error("rendered output should contain 'Hello'")
	}
	if !strings.Contains(output, "bold") {
		t.Error("rendered output should contain 'bold'")
	}
	if !strings.Contains(output, "italic") {
		t.Error("rendered output should contain 'italic'")
	}
}

func TestRenderMarkdown_CodeBlockStyle(t *testing.T) {
	input := "```go\nfmt.Println(\"hello\")\n```\n"
	lines := ParseMarkdown(input)
	output := RenderMarkdown(lines, 40)

	if !strings.Contains(output, "hello") {
		t.Error("should render code content")
	}
	if !strings.Contains(output, "fmt.Println") {
		t.Error("should render code content")
	}
}

func TestParseMarkdown_HorizontalRule(t *testing.T) {
	input := "Some text\n\n---\n\nMore text\n"
	lines := ParseMarkdown(input)

	foundHR := false
	for _, l := range lines {
		if l.BlockType == BlockHorizontalRule {
			foundHR = true
		}
	}
	if !foundHR {
		t.Error("should find horizontal rule")
	}
}

func TestWikiLinkResolution(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	vault, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault: %v", err)
	}

	// Exact path with .md
	resolved := ResolveWikiLink("index", vault, testVaultPath(t))
	if resolved != "index.md" {
		t.Errorf("resolved index = %q, want 'index.md'", resolved)
	}

	// Nested path
	resolved = ResolveWikiLink("projects/api-design", vault, testVaultPath(t))
	if resolved != "projects/api-design.md" {
		t.Errorf("resolved api-design = %q, want 'projects/api-design.md'", resolved)
	}

	// Basename match
	resolved = ResolveWikiLink("database", vault, testVaultPath(t))
	if resolved != "projects/database.md" {
		t.Errorf("resolved database = %q, want 'projects/database.md'", resolved)
	}

	// Non-existent
	resolved = ResolveWikiLink("nonexistent", vault, testVaultPath(t))
	if resolved != "" {
		t.Errorf("nonexistent should resolve to empty, got %q", resolved)
	}
}
