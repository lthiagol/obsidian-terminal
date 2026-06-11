package markdown

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// BlockType classifies a markdown block element.
type BlockType int

const (
	BlockParagraph BlockType = iota
	BlockHeading
	BlockCodeBlock
	BlockList
	BlockBlockquote
	BlockCallout
	BlockHorizontalRule
	BlockEmpty
)

// InlineSegment represents a styled span of inline text.
type InlineSegment struct {
	Text          string
	Bold          bool
	Italic        bool
	Strikethrough bool
	Code          bool
	Highlight     bool
	IsWikiLink    bool
	WikiTarget    string
	WikiDisplay   string
}

// MarkdownLine represents a single parsed line of markdown.
type MarkdownLine struct {
	BlockType    BlockType
	HeadingLevel int
	Segments     []InlineSegment
	IndentLevel  int
	CalloutType  string
	Language     string
	RawContent   string
}

// WikiLink represents an Obsidian [[wiki-link]].
type WikiLink struct {
	Target  string
	Display string
}

// RendererStyle holds colors for the markdown renderer.
type RendererStyle struct {
	Accent          lipgloss.Color
	AccentSecondary lipgloss.Color
	AccentTertiary  lipgloss.Color
	TextSecondary   lipgloss.Color
	TextDim         lipgloss.Color
	Success         lipgloss.Color
	CodeBackground  lipgloss.Color
	Heading1        lipgloss.Color
}

var (
	calloutTypeRe     = regexp.MustCompile(`\[!(\w+)\]`)
	stripCalloutRe    = regexp.MustCompile(`\[!\w+\][+-]?\s*`)
	listItemMarkerRe  = regexp.MustCompile(`^[\-\*\+]\s`)
	listItemOrderedRe = regexp.MustCompile(`^\d+\.\s`)
	listItemParseRe   = regexp.MustCompile(`^([\-\*\+]|\d+\.)\s+`)
	commentStripRe    = regexp.MustCompile(`%%.*?%%`)
	visibleLenRe      = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	inlineSpecialRe   = regexp.MustCompile(`\x60|\*\*|__|\*|_|~~|\[\[|==`)
)

// ParseMarkdown parses markdown content into structured lines.
func ParseMarkdown(content string) []MarkdownLine {
	content = StripFrontmatter(content)
	content = stripComments(content)

	lines := strings.Split(content, "\n")
	var result []MarkdownLine

	inCodeBlock := false
	var codeLang string
	var codeLines []string

	for _, line := range lines {
		if inCodeBlock {
			if isCodeFence(line) {
				result = append(result, MarkdownLine{
					BlockType:  BlockCodeBlock,
					Language:   codeLang,
					RawContent: strings.Join(codeLines, "\n"),
				})
				inCodeBlock = false
				codeLang = ""
				codeLines = nil
				continue
			}
			codeLines = append(codeLines, line)
			continue
		}

		if isCodeFence(line) {
			inCodeBlock = true
			codeLang = codeFenceLanguage(line)
			codeLines = nil
			continue
		}

		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			result = append(result, MarkdownLine{BlockType: BlockEmpty})
			continue
		}

		if isHorizontalRule(trimmed) {
			result = append(result, MarkdownLine{BlockType: BlockHorizontalRule})
			continue
		}

		if isCalloutStart(trimmed) {
			calloutType := extractCalloutType(trimmed)
			result = append(result, MarkdownLine{
				BlockType:   BlockCallout,
				CalloutType: calloutType,
				Segments:    []InlineSegment{{Text: calloutType}},
			})
			continue
		}

		if isHeading(line) {
			level := headingLevel(line)
			text := strings.TrimSpace(line[level:])
			text = strings.TrimSuffix(text, " #")
			segments := parseInline(text)
			result = append(result, MarkdownLine{
				BlockType:    BlockHeading,
				HeadingLevel: level,
				Segments:     segments,
			})
			continue
		}

		if isBlockquote(line) {
			text := stripBlockquote(line)
			indent := blockquoteIndent(line)
			segments := parseInline(text)
			result = append(result, MarkdownLine{
				BlockType:   BlockBlockquote,
				Segments:    segments,
				IndentLevel: indent,
			})
			continue
		}

		if isListItem(line) {
			indent, _, text := parseListItem(line)
			segments := parseInline(text)
			result = append(result, MarkdownLine{
				BlockType:   BlockList,
				Segments:    segments,
				IndentLevel: indent,
			})
			continue
		}

		segments := parseInline(line)
		result = append(result, MarkdownLine{
			BlockType: BlockParagraph,
			Segments:  segments,
		})
	}

	if inCodeBlock {
		result = append(result, MarkdownLine{
			BlockType:  BlockCodeBlock,
			Language:   codeLang,
			RawContent: strings.Join(codeLines, "\n"),
		})
	}

	return result
}

func isCodeFence(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")
}

func codeFenceLanguage(line string) string {
	t := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "```"))
	t = strings.TrimPrefix(t, "~~~")
	return strings.TrimSpace(t)
}

func isHeading(line string) bool {
	if len(line) == 0 || line[0] != '#' {
		return false
	}
	level := 0
	for i, c := range line {
		if c == '#' {
			level++
		} else if c == ' ' && level > 0 {
			return level <= 6 && i+1 < len(line)
		} else {
			return false
		}
	}
	return false
}

func headingLevel(line string) int {
	level := 0
	for _, c := range line {
		if c == '#' {
			level++
		} else {
			break
		}
	}
	return level
}

func isHorizontalRule(line string) bool {
	t := strings.ReplaceAll(line, " ", "")
	return t == "---" || t == "***" || t == "___" || t == "----------"
}

func isBlockquote(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " "), ">")
}

func isCalloutStart(line string) bool {
	t := strings.TrimLeft(line, " ")
	return strings.HasPrefix(t, "> [!")
}

func extractCalloutType(line string) string {
	matches := calloutTypeRe.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return strings.ToLower(matches[1])
	}
	return "note"
}

func stripBlockquote(line string) string {
	t := strings.TrimLeft(line, " ")
	t = strings.TrimPrefix(t, ">")
	t = strings.TrimPrefix(t, " ")
	t = stripCalloutRe.ReplaceAllString(t, "")
	return strings.TrimLeft(t, " ")
}

func blockquoteIndent(line string) int {
	count := 0
	for i, c := range line {
		if c == '>' {
			count++
		} else if c == ' ' {
			break
		} else {
			if i < len(line) {
				break
			}
		}
	}
	if count == 0 {
		count = 1
	}
	return count
}

func isListItem(line string) bool {
	t := strings.TrimLeft(line, " ")
	return listItemMarkerRe.MatchString(strings.TrimLeft(t, " ")) ||
		listItemOrderedRe.MatchString(strings.TrimLeft(t, " "))
}

func parseListItem(line string) (int, string, string) {
	indent := 0
	for _, c := range line {
		if c == ' ' || c == '\t' {
			indent++
		} else {
			break
		}
	}
	t := strings.TrimLeft(line, " ")
	loc := listItemParseRe.FindStringIndex(t)
	if loc != nil {
		marker := t[loc[0]:loc[1]]
		text := t[loc[1]:]
		return indent / 2, marker, text
	}
	return 0, "", t
}

func parseInline(text string) []InlineSegment {
	var segments []InlineSegment
	parseSegments(text, &segments)
	return mergeSegments(segments)
}

func parseSegments(text string, segments *[]InlineSegment) {
	if text == "" {
		return
	}

	if strings.HasPrefix(text, "[[") {
		end := strings.Index(text, "]]")
		if end > 0 {
			inner := text[2:end]
			parts := strings.SplitN(inner, "|", 2)
			target := parts[0]
			display := target
			if len(parts) > 1 {
				display = parts[1]
			}

			target = strings.SplitN(target, "#", 2)[0]

			*segments = append(*segments, InlineSegment{
				Text:        display,
				IsWikiLink:  true,
				WikiTarget:  target,
				WikiDisplay: display,
			})
			parseSegments(text[end+2:], segments)
			return
		}
	}

	if strings.HasPrefix(text, "***") || strings.HasPrefix(text, "___") {
		marker := text[:3]
		end := strings.Index(text[3:], marker)
		if end >= 0 {
			inner := text[3 : 3+end]
			*segments = append(*segments, InlineSegment{
				Text:   inner,
				Bold:   true,
				Italic: true,
			})
			parseSegments(text[3+end+3:], segments)
			return
		}
	}

	if strings.HasPrefix(text, "**") || strings.HasPrefix(text, "__") {
		marker := text[:2]
		end := strings.Index(text[2:], marker)
		if end >= 0 {
			inner := text[2 : 2+end]
			*segments = append(*segments, InlineSegment{
				Text: inner,
				Bold: true,
			})
			parseSegments(text[2+end+2:], segments)
			return
		}
	}

	if strings.HasPrefix(text, "*") || strings.HasPrefix(text, "_") {
		marker := text[:1]
		rest := text[1:]
		end := strings.Index(rest, marker)
		if end > 0 && end < len(rest)-1 {
			inner := rest[:end]
			*segments = append(*segments, InlineSegment{
				Text:   inner,
				Italic: true,
			})
			parseSegments(text[1+end+1:], segments)
			return
		}
	}

	if strings.HasPrefix(text, "`") {
		end := strings.Index(text[1:], "`")
		if end >= 0 {
			inner := text[1 : 1+end]
			*segments = append(*segments, InlineSegment{
				Text: inner,
				Code: true,
			})
			parseSegments(text[1+end+1:], segments)
			return
		}
	}

	if strings.HasPrefix(text, "~~") {
		end := strings.Index(text[2:], "~~")
		if end >= 0 {
			inner := text[2 : 2+end]
			*segments = append(*segments, InlineSegment{
				Text:          inner,
				Strikethrough: true,
			})
			parseSegments(text[2+end+2:], segments)
			return
		}
	}

	if strings.HasPrefix(text, "==") {
		end := strings.Index(text[2:], "==")
		if end >= 0 {
			inner := text[2 : 2+end]
			*segments = append(*segments, InlineSegment{
				Text:      inner,
				Highlight: true,
			})
			parseSegments(text[2+end+2:], segments)
			return
		}
	}

	next := findNextSpecial(text)
	if next == -1 {
		*segments = append(*segments, InlineSegment{Text: text})
		return
	}
	if next > 0 {
		*segments = append(*segments, InlineSegment{Text: text[:next]})
	}
	parseSegments(text[next:], segments)
}

func findNextSpecial(text string) int {
	loc := inlineSpecialRe.FindStringIndex(text)
	if loc == nil {
		return -1
	}
	return loc[0]
}

func mergeSegments(segments []InlineSegment) []InlineSegment {
	if len(segments) <= 1 {
		return segments
	}
	var merged []InlineSegment
	current := segments[0]
	for i := 1; i < len(segments); i++ {
		next := segments[i]
		if current.Bold == next.Bold &&
			current.Italic == next.Italic &&
			current.Code == next.Code &&
			current.Strikethrough == next.Strikethrough &&
			current.Highlight == next.Highlight &&
			current.IsWikiLink == next.IsWikiLink &&
			current.WikiTarget == next.WikiTarget {
			current.Text += next.Text
		} else {
			if current.Text != "" {
				merged = append(merged, current)
			}
			current = next
		}
	}
	if current.Text != "" {
		merged = append(merged, current)
	}
	return merged
}

// StripFrontmatter removes YAML frontmatter (--- ... ---) from content.
func StripFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return content
	}
	rest := content[3:]
	idx := strings.Index(rest, "\n---\n")
	if idx == -1 {
		idx = strings.Index(rest, "\n---\r\n")
	}
	if idx == -1 {
		return content
	}
	return rest[idx+5:]
}

func stripComments(content string) string {
	return commentStripRe.ReplaceAllString(content, "")
}

// ExtractWikiLinks extracts unique wiki-links from parsed markdown.
func ExtractWikiLinks(lines []MarkdownLine) []WikiLink {
	seen := make(map[string]bool)
	var links []WikiLink
	for _, line := range lines {
		for _, seg := range line.Segments {
			if seg.IsWikiLink {
				if !seen[seg.WikiTarget] {
					seen[seg.WikiTarget] = true
					links = append(links, WikiLink{
						Target:  seg.WikiTarget,
						Display: seg.WikiDisplay,
					})
				}
			}
		}
	}
	return links
}

// RenderMarkdown renders parsed markdown lines to styled terminal output.
func RenderMarkdown(lines []MarkdownLine, width int, style RendererStyle) string {
	if width < 20 {
		width = 20
	}

	var sb strings.Builder
	for _, line := range lines {
		rendered := renderLine(line, width, style)
		if rendered != "" {
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(rendered)
		}
	}
	return sb.String()
}

func renderLine(line MarkdownLine, width int, style RendererStyle) string {
	switch line.BlockType {
	case BlockHeading:
		return renderHeading(line, width, style)
	case BlockCodeBlock:
		return renderCodeBlock(line, width, style)
	case BlockList:
		return renderList(line, width, style)
	case BlockBlockquote:
		return renderBlockquote(line, width, style)
	case BlockCallout:
		return renderCallout(line, width, style)
	case BlockHorizontalRule:
		return renderHorizontalRule(width, style)
	case BlockEmpty:
		return ""
	default:
		return renderParagraph(line, width, style)
	}
}

func renderHeading(line MarkdownLine, width int, style RendererStyle) string {
	text := renderSegments(line.Segments, style)
	var s lipgloss.Style
	switch line.HeadingLevel {
	case 1:
		s = lipgloss.NewStyle().Foreground(style.Heading1).Bold(true).Underline(true)
	case 2:
		s = lipgloss.NewStyle().Foreground(style.Accent).Bold(true)
	case 3:
		s = lipgloss.NewStyle().Foreground(style.AccentTertiary).Bold(true)
	default:
		s = lipgloss.NewStyle().Foreground(style.TextSecondary).Bold(true)
	}
	return s.Render(text)
}

func renderCodeBlock(line MarkdownLine, width int, style RendererStyle) string {
	lines := strings.Split(line.RawContent, "\n")

	header := ""
	if line.Language != "" {
		header = " " + line.Language + " "
	}

	var sb strings.Builder

	topBorder := lipgloss.NewStyle().Foreground(style.TextDim).Render("╭" + strings.Repeat("─", width-2) + "╮")
	if header != "" {
		labelStyle := lipgloss.NewStyle().Foreground(style.TextDim)
		padded := header + strings.Repeat("─", width-len(header)-2)
		topBorder = lipgloss.NewStyle().Foreground(style.TextDim).Render("╭") +
			labelStyle.Render(padded) +
			lipgloss.NewStyle().Foreground(style.TextDim).Render("╮")
	}
	sb.WriteString(topBorder)

	codeStyle := lipgloss.NewStyle().Foreground(style.Success)

	for _, l := range lines {
		sb.WriteString("\n")
		l = strings.ReplaceAll(l, "\t", "    ")
		if len(l) > width-2 {
			l = l[:width-2]
		}
		padded := l + strings.Repeat(" ", width-2-len(l))
		lineContent := lipgloss.NewStyle().Foreground(style.TextDim).Render("│") +
			codeStyle.Render(padded) +
			lipgloss.NewStyle().Foreground(style.TextDim).Render("│")
		sb.WriteString(lineContent)
	}

	sb.WriteString("\n")
	botBorder := lipgloss.NewStyle().Foreground(style.TextDim).Render("╰" + strings.Repeat("─", width-2) + "╯")
	sb.WriteString(botBorder)

	return sb.String()
}

func renderList(line MarkdownLine, width int, style RendererStyle) string {
	prefix := strings.Repeat("  ", line.IndentLevel)
	text := renderSegments(line.Segments, style)
	bullet := lipgloss.NewStyle().Foreground(style.Accent).Render("•")
	return prefix + bullet + " " + text
}

func renderBlockquote(line MarkdownLine, width int, style RendererStyle) string {
	prefix := lipgloss.NewStyle().Foreground(style.Accent).Render("│")
	text := renderSegments(line.Segments, style)
	bodyStyle := lipgloss.NewStyle().Foreground(style.TextSecondary).Italic(true)
	return prefix + " " + bodyStyle.Render(text)
}

func renderCallout(line MarkdownLine, width int, style RendererStyle) string {
	icon := "ℹ"
	switch line.CalloutType {
	case "note":
		icon = "📝"
	case "tip":
		icon = "💡"
	case "warning":
		icon = "⚠"
	case "danger":
		icon = "🚫"
	case "info":
		icon = "ℹ"
	case "todo":
		icon = "☐"
	case "question":
		icon = "❓"
	case "success":
		icon = "✅"
	case "bug":
		icon = "🐛"
	case "example":
		icon = "📋"
	}

	typeStyle := lipgloss.NewStyle().Bold(true).Foreground(style.AccentSecondary)
	bodyStyle := lipgloss.NewStyle().Foreground(style.TextSecondary)

	body := ""
	if len(line.Segments) > 0 {
		body = line.Segments[0].Text
	}
	return icon + " " + typeStyle.Render(line.CalloutType) + " " + bodyStyle.Render(body)
}

func renderHorizontalRule(width int, style RendererStyle) string {
	rule := strings.Repeat("─", width)
	return lipgloss.NewStyle().Foreground(style.TextDim).Render(rule)
}

func renderParagraph(line MarkdownLine, width int, style RendererStyle) string {
	text := renderSegments(line.Segments, style)
	return wrapText(text, width)
}

func renderSegments(segments []InlineSegment, style RendererStyle) string {
	var sb strings.Builder
	for _, seg := range segments {
		sb.WriteString(renderSegment(seg, style))
	}
	return sb.String()
}

func renderSegment(seg InlineSegment, style RendererStyle) string {
	s := lipgloss.NewStyle()

	switch {
	case seg.IsWikiLink:
		s = s.Foreground(style.AccentTertiary).Underline(true)
	case seg.Code:
		s = s.Foreground(style.Success).Background(style.CodeBackground)
	case seg.Highlight:
		s = s.Foreground(style.AccentSecondary)
	case seg.Strikethrough:
		s = s.Strikethrough(true)
	case seg.Bold && seg.Italic:
		s = s.Bold(true).Italic(true)
	case seg.Bold:
		s = s.Bold(true)
	case seg.Italic:
		s = s.Italic(true)
	}

	return s.Render(seg.Text)
}

func wrapText(text string, width int) string {
	if width <= 0 || len(text) <= width {
		return text
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	currentLine := ""
	for _, word := range words {
		wordLen := visibleLen(word)
		currentLen := visibleLen(currentLine)

		if currentLen == 0 {
			currentLine = word
		} else if currentLen+1+wordLen <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

func visibleLen(s string) int {
	clean := visibleLenRe.ReplaceAllString(s, "")
	return len([]rune(clean))
}
