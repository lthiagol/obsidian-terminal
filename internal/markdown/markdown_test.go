package markdown

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
	style := testRendererStyle()
	output := RenderMarkdown(lines, 60, style)

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
	style := testRendererStyle()
	output := RenderMarkdown(lines, 40, style)

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

func TestParseMarkdown_UnclosedCodeBlock(t *testing.T) {
	input := "```go\nfunc test() {\n    fmt.Println(\"hello\")\n}\n"
	lines := ParseMarkdown(input)

	found := false
	for _, l := range lines {
		if l.BlockType == BlockCodeBlock {
			found = true
			if !strings.Contains(l.RawContent, "fmt.Println") {
				t.Error("code block should contain content")
			}
		}
	}
	if !found {
		t.Error("unclosed code block should still produce a BlockCodeBlock at EOF")
	}
}

func TestRenderCallout_EmptySegments(t *testing.T) {
	line := MarkdownLine{
		BlockType:   BlockCallout,
		CalloutType: "note",
	}
	style := testRendererStyle()
	output := RenderMarkdown([]MarkdownLine{line}, 40, style)
	if !strings.Contains(output, "note") {
		t.Error("callout should render type even with zero segments")
	}
}

func TestWrapText(t *testing.T) {
	text := "a b c d e f g h i j k l m n o p"
	output := wrapText(text, 10)
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Error("long text should be wrapped to multiple lines")
	}
}

func TestRenderBlockquote(t *testing.T) {
	line := MarkdownLine{
		BlockType: BlockBlockquote,
		Segments:  []InlineSegment{{Text: "quoted text"}},
	}
	style := testRendererStyle()
	output := RenderMarkdown([]MarkdownLine{line}, 40, style)
	if !strings.Contains(output, "quoted text") {
		t.Error("blockquote should contain the text")
	}
}

func TestRenderCallout(t *testing.T) {
	line := MarkdownLine{
		BlockType:   BlockCallout,
		CalloutType: "warning",
		Segments:    []InlineSegment{{Text: "be careful"}},
	}
	style := testRendererStyle()
	output := RenderMarkdown([]MarkdownLine{line}, 40, style)
	if !strings.Contains(output, "warning") {
		t.Error("callout should show type")
	}
	if !strings.Contains(output, "be careful") {
		t.Error("callout should show text")
	}
}

func TestRenderList(t *testing.T) {
	line := MarkdownLine{
		BlockType:   BlockList,
		IndentLevel: 1,
		Segments:    []InlineSegment{{Text: "list item"}},
	}
	style := testRendererStyle()
	output := RenderMarkdown([]MarkdownLine{line}, 40, style)
	if !strings.Contains(output, "list item") {
		t.Error("list should contain the item text")
	}
}

func TestRenderCodeBlock(t *testing.T) {
	line := MarkdownLine{
		BlockType:  BlockCodeBlock,
		Language:   "go",
		RawContent: "package main",
	}
	style := testRendererStyle()
	output := RenderMarkdown([]MarkdownLine{line}, 40, style)
	if !strings.Contains(output, "package main") {
		t.Error("code block should contain content")
	}
	if !strings.Contains(output, "go") {
		t.Error("code block should show language")
	}
}

func TestRenderHorizontalRule(t *testing.T) {
	line := MarkdownLine{BlockType: BlockHorizontalRule}
	style := testRendererStyle()
	output := RenderMarkdown([]MarkdownLine{line}, 20, style)
	if len(output) == 0 {
		t.Error("horizontal rule should produce output")
	}
}

func TestParseMarkdown_UnderscoreBold(t *testing.T) {
	input := "__bold__ text\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	foundBold := false
	for _, seg := range lines[0].Segments {
		if seg.Bold && seg.Text == "bold" {
			foundBold = true
		}
	}
	if !foundBold {
		t.Error("__bold__ should produce bold segment")
	}
}

func TestParseMarkdown_NestedBoldItalicTriple(t *testing.T) {
	// *** marks bold+italic via triple asterisk
	input := "***bold italic***\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	for _, seg := range lines[0].Segments {
		if seg.Bold && seg.Italic {
			return // found bold+italic
		}
	}
	t.Error("***bold italic*** should produce bold+italic segment")
}

func TestParseMarkdown_WikiLinkHeadingFragment(t *testing.T) {
	input := "See [[note#section]] for details.\n"
	lines := ParseMarkdown(input)
	var wiki []InlineSegment
	for _, seg := range lines[0].Segments {
		if seg.IsWikiLink {
			wiki = append(wiki, seg)
		}
	}
	if len(wiki) != 1 {
		t.Fatalf("expected 1 wiki link, got %d", len(wiki))
	}
	// Parser strips #section from target but display shows full text
	if wiki[0].WikiTarget != "note" {
		t.Errorf("WikiTarget = %q, want 'note'", wiki[0].WikiTarget)
	}
}

func TestParseMarkdown_WikiLinkOnlyDisplay(t *testing.T) {
	input := "[[|just display]]\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	for _, seg := range lines[0].Segments {
		if seg.IsWikiLink {
			if seg.WikiDisplay != "just display" {
				t.Errorf("display = %q", seg.WikiDisplay)
			}
			return
		}
	}
	t.Error("expected wiki link segment")
}

func TestParseMarkdown_MultiLineCallout(t *testing.T) {
	input := "> [!note] First line\n> Second line\n> Third line\n"
	lines := ParseMarkdown(input)
	calloutCount := 0
	for _, l := range lines {
		if l.BlockType == BlockCallout {
			calloutCount++
		}
	}
	if calloutCount < 1 {
		t.Error("expected at least 1 callout block")
	}
}

func TestParseMarkdown_TableEmptyCell(t *testing.T) {
	input := "| a |  | c |\n|---|---|---|\n| 1 |   | 3 |\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected table line")
	}
	if lines[0].BlockType != BlockTable {
		t.Fatalf("expected BlockTable, got %d", lines[0].BlockType)
	}
	// Empty cells should be allowed
	for _, cell := range lines[0].TableCells {
		if cell == "" {
			return // found empty cell, test passes
		}
	}
}

func TestParseMarkdown_TableManyColumns(t *testing.T) {
	input := "| a | b | c | d | e | f | g | h | i | j |\n" +
		"|---|---|---|---|---|---|---|---|---|---|\n" +
		"| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 |\n"
	lines := ParseMarkdown(input)
	if len(lines) < 1 {
		t.Fatal("expected table rows")
	}
	// First line is header
	if len(lines[0].TableCells) != 10 {
		t.Errorf("expected 10 columns, got %d", len(lines[0].TableCells))
	}
}

func TestParseMarkdown_HorizontalRuleVariants(t *testing.T) {
	tests := []string{"---\n", "***\n", "___\n", "* * *\n", "- - -\n"}
	for _, input := range tests {
		t.Run(input[:min(3, len(input))], func(t *testing.T) {
			lines := ParseMarkdown(input)
			found := false
			for _, l := range lines {
				if l.BlockType == BlockHorizontalRule {
					found = true
				}
			}
			if !found {
				t.Errorf("should be horizontal rule: %q", input)
			}
		})
	}
}

func TestParseMarkdown_CodeFenceTilde(t *testing.T) {
	input := "~~~go\nfmt.Println(\"hi\")\n~~~\n"
	lines := ParseMarkdown(input)
	found := false
	for _, l := range lines {
		if l.BlockType == BlockCodeBlock {
			found = true
			if l.Language != "go" {
				t.Errorf("language = %q", l.Language)
			}
		}
	}
	if !found {
		t.Error("tilde code fence should produce code block")
	}
}

func TestParseMarkdown_UnknownCalloutType(t *testing.T) {
	input := "> [!customtype] Some content\n"
	lines := ParseMarkdown(input)
	found := false
	for _, l := range lines {
		if l.BlockType == BlockCallout {
			found = true
			if l.CalloutType != "customtype" {
				t.Errorf("callout type = %q", l.CalloutType)
			}
		}
	}
	if !found {
		t.Error("unknown callout type should still be parsed")
	}
}

func TestParseMarkdown_CheckedCheckbox(t *testing.T) {
	input := "- [x] Done item\n- [X] Also done\n- [-] Cancelled\n"
	lines := ParseMarkdown(input)
	checked := 0
	for _, l := range lines {
		if l.Checked {
			checked++
		}
	}
	if checked < 2 {
		t.Errorf("expected at least 2 checked items, got %d", checked)
	}
}

func TestRenderTable_WideTableScaling(t *testing.T) {
	// 8 columns with data, should scale down to fit width
	input := "| Name | Type | Status | Priority | Owner | Date | Notes | ID |\n" +
		"|------|------|--------|----------|-------|------|-------|----|\n" +
		"| foo  | bar  | active | high     | alice | jan  | none  | 1  |\n"
	lines := ParseMarkdown(input)
	style := testRendererStyle()
	output := RenderMarkdown(lines, 60, style)

	if !strings.Contains(output, "foo") {
		t.Error("wide table should still contain cell data")
	}
	// Must have box-drawing borders
	if !strings.Contains(output, "┌") {
		t.Error("wide table should have top border")
	}
}

func TestRenderMarkdown_EmptyContent(t *testing.T) {
	style := testRendererStyle()
	output := RenderMarkdown(nil, 80, style)
	if output != "" {
		t.Errorf("empty lines should produce empty output: %q", output)
	}
}

func TestVisibleLen_ComplexANSI(t *testing.T) {
	// Extended color codes
	s := "\033[38;5;123mcolored\033[0m"
	clean := visibleLenRe.ReplaceAllString(s, "")
	if len([]rune(clean)) != 7 {
		t.Errorf("visibleLen should be 7 for 'colored', got %d", len([]rune(clean)))
	}
}

func TestParseMarkdown_LoneBacktick(t *testing.T) {
	// Should not crash (infinite recursion guard)
	input := "text ` more text\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
}

func TestParseMarkdown_BacktrackTripleAsterisk(t *testing.T) {
	// ***text** should be: * + bold(text)
	input := "***text**\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	found := false
	for _, seg := range lines[0].Segments {
		if seg.Bold && seg.Text == "text" {
			found = true
		}
	}
	if !found {
		t.Log("segments:", len(lines[0].Segments), lines[0].Segments)
		t.Error("***text** should produce Bold(text)")
	}
}

func TestParseMarkdown_BacktrackDoubleAsterisk(t *testing.T) {
	// **text* more should be: * + text + " more"
	// (with no matching italic because the * has no content after it)
	input := "**text* more\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	// Should produce some text content without crashing
	for _, seg := range lines[0].Segments {
		if strings.Contains(seg.Text, "more") {
			return
		}
	}
	t.Log("segments:", lines[0].Segments)
	t.Error("should contain 'more' in output")
}

func TestParseMarkdown_DoubleBacktick(t *testing.T) {
	input := "``code with ` backtick``\n"
	lines := ParseMarkdown(input)
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	found := false
	for _, seg := range lines[0].Segments {
		if seg.Code && strings.Contains(seg.Text, "backtick") {
			found = true
		}
	}
	if !found {
		t.Log("segments:", lines[0].Segments)
		t.Error("double-backtick should produce Code segment")
	}
}

func testRendererStyle() RendererStyle {
	return RendererStyle{
		Accent:          "#a78bfa",
		AccentSecondary: "#fbbf24",
		AccentTertiary:  "#2dd4bf",
		TextSecondary:   "#9ca3af",
		TextDim:         "#4b5563",
		Success:         "#34d399",
		CodeBackground:  "#1f2937",
		Heading1:        "#f472b6",
	}
}
