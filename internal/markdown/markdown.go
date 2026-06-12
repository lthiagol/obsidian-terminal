package markdown

import (
	"fmt"
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
	BlockTable
	BlockEmbed
	BlockEmbedStart
	BlockEmbedEnd
)

// TableAlignment specifies cell text alignment.
type TableAlignment int

const (
	AlignLeft TableAlignment = iota
	AlignCenter
	AlignRight
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
	Checkable    bool
	Checked      bool
	TableCells   []string
	TableAlign   []TableAlignment
	EmbedTarget  string
	EmbedHeading string
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

	for i := 0; i < len(lines); i++ {
		line := lines[i]

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

		if isTableLine(trimmed) && i+1 < len(lines) && isTableSeparator(lines[i+1]) {
			headerCells := parseTableRow(trimmed)
			aligns := parseTableAlignment(lines[i+1])
			i++

			result = append(result, MarkdownLine{
				BlockType:  BlockTable,
				TableCells: headerCells,
				TableAlign: aligns,
			})

			for i+1 < len(lines) && isTableLine(strings.TrimSpace(lines[i+1])) {
				i++
				dataCells := parseTableRow(strings.TrimSpace(lines[i]))
				result = append(result, MarkdownLine{
					BlockType:  BlockTable,
					TableCells: dataCells,
					TableAlign: aligns,
				})
			}
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
			checkable := false
			checked := false
			if strings.HasPrefix(text, "[ ] ") {
				checkable = true
				text = text[4:]
			} else if strings.HasPrefix(text, "[x] ") || strings.HasPrefix(text, "[X] ") {
				checkable = true
				checked = true
				text = text[4:]
			}
			segments := parseInline(text)
			result = append(result, MarkdownLine{
				BlockType:   BlockList,
				Segments:    segments,
				IndentLevel: indent,
				Checkable:   checkable,
				Checked:     checked,
			})
			continue
		}

		if isEmbed(line) {
			target, heading := parseEmbed(line)
			result = append(result, MarkdownLine{
				BlockType:    BlockEmbed,
				EmbedTarget:  target,
				EmbedHeading: heading,
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

func isTableLine(line string) bool {
	return strings.HasPrefix(line, "|")
}

func isTableSeparator(line string) bool {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "|") {
		return false
	}
	for _, c := range t {
		if c != '|' && c != '-' && c != ':' && c != ' ' {
			return false
		}
	}
	return strings.Contains(t, "-")
}

func parseTableRow(line string) []string {
	t := strings.TrimSpace(line)
	t = strings.ReplaceAll(t, `\|`, "\x00")
	t = strings.TrimLeft(t, "|")
	t = strings.TrimRight(t, "|")
	cells := strings.Split(t, "|")
	for i, cell := range cells {
		cells[i] = strings.ReplaceAll(strings.TrimSpace(cell), "\x00", "|")
	}
	return cells
}

func parseTableAlignment(line string) []TableAlignment {
	t := strings.TrimSpace(line)
	t = strings.TrimLeft(t, "|")
	t = strings.TrimRight(t, "|")
	parts := strings.Split(t, "|")
	var aligns []TableAlignment
	for _, p := range parts {
		p = strings.TrimSpace(p)
		left := strings.HasPrefix(p, ":")
		right := strings.HasSuffix(p, ":")
		if left && right {
			aligns = append(aligns, AlignCenter)
		} else if right {
			aligns = append(aligns, AlignRight)
		} else {
			aligns = append(aligns, AlignLeft)
		}
	}
	return aligns
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

func isEmbed(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "![[") && strings.HasSuffix(t, "]]")
}

func parseEmbed(line string) (target, heading string) {
	t := strings.TrimSpace(line)
	t = strings.TrimPrefix(t, "![[")
	t = strings.TrimSuffix(t, "]]")
	if idx := strings.Index(t, "#"); idx >= 0 {
		heading = t[idx+1:]
		target = t[:idx]
	} else {
		target = t
	}
	target = strings.SplitN(target, "|", 2)[0]
	return
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
	if idx == -1 && strings.HasSuffix(rest, "\n---") {
		return strings.TrimSuffix(rest, "\n---")
	}
	if idx == -1 && strings.HasSuffix(rest, "\n---\r") {
		return strings.TrimSuffix(rest, "\n---\r")
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

// EmbedResolver resolves an embed target to its content.
type EmbedResolver func(target, heading string) (string, error)

// ResolveEmbeds walks lines and resolves BlockEmbed lines by calling the resolver.
func ResolveEmbeds(lines []MarkdownLine, resolve EmbedResolver) []MarkdownLine {
	return resolveEmbedsRecursive(lines, resolve, make(map[string]bool), 0)
}

func resolveEmbedsRecursive(lines []MarkdownLine, resolve EmbedResolver, visited map[string]bool, depth int) []MarkdownLine {
	if depth > 2 {
		return lines
	}

	var result []MarkdownLine
	for _, line := range lines {
		if line.BlockType != BlockEmbed {
			result = append(result, line)
			continue
		}

		key := line.EmbedTarget
		if line.EmbedHeading != "" {
			key += "#" + line.EmbedHeading
		}

		if visited[key] {
			result = append(result, MarkdownLine{
				BlockType: BlockEmbedStart,
				EmbedTarget:  "(circular embed detected)",
			})
			result = append(result, MarkdownLine{BlockType: BlockEmbedEnd})
			continue
		}

		content, err := resolve(line.EmbedTarget, line.EmbedHeading)
		if err != nil || content == "" {
			result = append(result, MarkdownLine{
				BlockType: BlockEmbedStart,
				EmbedTarget:  line.EmbedTarget,
			})
			result = append(result, MarkdownLine{
				BlockType:    BlockParagraph,
				Segments:     []InlineSegment{{Text: "(embed not found: " + line.EmbedTarget + ")"}},
			})
			result = append(result, MarkdownLine{BlockType: BlockEmbedEnd})
			continue
		}

		visited[key] = true

		parsed := ParseMarkdown(content)
		resolved := resolveEmbedsRecursive(parsed, resolve, visited, depth+1)

		result = append(result, MarkdownLine{
			BlockType:    BlockEmbedStart,
			EmbedTarget:  line.EmbedTarget,
			EmbedHeading: line.EmbedHeading,
		})
		result = append(result, resolved...)
		result = append(result, MarkdownLine{BlockType: BlockEmbedEnd})
	}

	return result
}

// RenderMarkdown renders parsed markdown lines to styled terminal output.
func RenderMarkdown(lines []MarkdownLine, width int, style RendererStyle) string {
	if width < 20 {
		width = 20
	}

	var sb strings.Builder
	for i := 0; i < len(lines); i++ {
		if lines[i].BlockType == BlockTable {
			tableLines := []MarkdownLine{lines[i]}
			for i+1 < len(lines) && lines[i+1].BlockType == BlockTable {
				i++
				tableLines = append(tableLines, lines[i])
			}
			rendered := renderTableBlock(tableLines, width, style)
			if rendered != "" {
				if sb.Len() > 0 {
					sb.WriteString("\n")
				}
				sb.WriteString(rendered)
			}
			continue
		}

		if lines[i].BlockType == BlockEmbedStart {
			embedLines := []MarkdownLine{lines[i]}
			for i+1 < len(lines) && lines[i+1].BlockType != BlockEmbedEnd {
				i++
				embedLines = append(embedLines, lines[i])
			}
			if i+1 < len(lines) {
				i++
				embedLines = append(embedLines, lines[i])
			}
			rendered := renderEmbedBlock(embedLines, width, style)
			if rendered != "" {
				if sb.Len() > 0 {
					sb.WriteString("\n")
				}
				sb.WriteString(rendered)
			}
			continue
		}

		if lines[i].BlockType == BlockEmbedEnd {
			continue
		}

		rendered := renderLine(lines[i], width, style)
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
	case BlockTable:
		return "" // handled by renderTableBlock
	case BlockEmbed:
		return "" // handled by ResolveEmbeds
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

	var bullet string
	if line.Checkable {
		if line.Checked {
			bullet = lipgloss.NewStyle().Foreground(style.Success).Render("[x]")
		} else {
			bullet = lipgloss.NewStyle().Foreground(style.TextDim).Render("[ ]")
		}
	} else {
		bullet = lipgloss.NewStyle().Foreground(style.Accent).Render("•")
	}

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

func renderEmbedBlock(lines []MarkdownLine, width int, style RendererStyle) string {
	if len(lines) < 2 {
		return ""
	}

	start := lines[0]
	target := start.EmbedTarget
	if start.EmbedHeading != "" {
		target += " > " + start.EmbedHeading
	}

	var sb strings.Builder

	borderStyle := lipgloss.NewStyle().Foreground(style.Accent)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(style.AccentTertiary)
	dimStyle := lipgloss.NewStyle().Foreground(style.TextDim)

	sb.WriteString(borderStyle.Render("┌─ "))
	sb.WriteString(headerStyle.Render(target))
	sb.WriteString("\n")

	for i := 1; i < len(lines)-1; i++ {
		rendered := renderLine(lines[i], width-2, style)
		if rendered != "" {
			sb.WriteString(borderStyle.Render("│ "))
			sb.WriteString(rendered)
			sb.WriteString("\n")
		}
	}

	sb.WriteString(borderStyle.Render("└"))
	sb.WriteString(dimStyle.Render(strings.Repeat("─", width-1)))

	return sb.String()
}

func renderTableBlock(lines []MarkdownLine, width int, style RendererStyle) string {
	if len(lines) == 0 || len(lines[0].TableCells) == 0 {
		return ""
	}

	colCount := len(lines[0].TableCells)
	desired := make([]int, colCount)

	for _, line := range lines {
		for j, cell := range line.TableCells {
			if j < colCount {
				w := len([]rune(cell))
				if w > desired[j] {
					desired[j] = w
				}
			}
		}
	}

	// Box-drawing mode
	boxOverhead := colCount*3 + 1
	if boxOverhead < width {
		avail := width - boxOverhead
		colWidths := allocateColWidths(desired, avail)
		totalContent := 0
		for _, w := range colWidths {
			totalContent += w
		}
		if totalContent+boxOverhead <= width {
			return renderTableBox(lines, colWidths, style)
		}
	}

	// Borderless mode
	gapOverhead := (colCount - 1) * 2
	if gapOverhead < width {
		avail := width - gapOverhead
		colWidths := allocateColWidths(desired, avail)
		totalContent := 0
		for _, w := range colWidths {
			totalContent += w
		}
		if totalContent+gapOverhead <= width {
			return renderBorderlessTable(lines, colWidths, style)
		}
	}

	// Single-column fallback
	return renderSingleColumnTable(lines, style)
}

func allocateColWidths(desired []int, available int) []int {
	if len(desired) == 0 {
		return nil
	}
	result := make([]int, 0, len(desired))
	totalDesired := 0
	for _, w := range desired {
		totalDesired += max(w, 3)
	}
	if totalDesired <= available {
		for _, w := range desired {
			result = append(result, max(w, 3))
		}
		return result
	}

	// Largest remainder: give each floor, distribute leftovers
	given := 0
	var remainders []float64
	for _, w := range desired {
		minW := max(w, 3)
		alloc := int(float64(minW) / float64(totalDesired) * float64(available))
		if alloc < 3 {
			alloc = 3
		}
		if alloc > minW {
			alloc = minW
		}
		result = append(result, alloc)
		given += alloc
		rem := float64(minW)/float64(totalDesired)*float64(available) - float64(alloc)
		remainders = append(remainders, rem)
	}

	remaining := available - given
	for remaining > 0 {
		bestIdx := -1
		bestRem := -1.0
		for i, r := range remainders {
			if r > bestRem && result[i] < max(desired[i], 3) {
				bestRem = r
				bestIdx = i
			}
		}
		if bestIdx < 0 {
			break
		}
		result[bestIdx]++
		remainders[bestIdx] = -1
		remaining--
	}
	return result
}

func wrapCell(content string, width int) []string {
	if width < 1 {
		return []string{""}
	}
	runes := []rune(content)
	if len(runes) <= width {
		return []string{content}
	}

	var lines []string
	for len(runes) > 0 {
		if len(runes) <= width {
			lines = append(lines, string(runes))
			break
		}
		// Try to break at last space within width
		chunk := runes[:width+1]
		lastSpace := -1
		for j := len(chunk) - 1; j >= 0; j-- {
			if chunk[j] == ' ' {
				lastSpace = j
				break
			}
		}
		if lastSpace > 0 {
			lines = append(lines, string(runes[:lastSpace]))
			runes = runes[lastSpace+1:]
		} else {
			// No space found, hard-break
			lines = append(lines, string(runes[:width]))
			runes = runes[width:]
		}
	}
	return lines
}

func renderTableBox(lines []MarkdownLine, colWidths []int, style RendererStyle) string {
	var sb strings.Builder
	sepStyle := lipgloss.NewStyle().Foreground(style.TextDim)
	headerStyle := lipgloss.NewStyle().Foreground(style.Accent).Bold(true)
	cellStyle := lipgloss.NewStyle().Foreground(style.TextSecondary)

	// Wrap all cells
	type wrappedRow struct {
		lines [][]string
		maxH  int
	}
	var rows []wrappedRow
	for _, row := range lines {
		var wr wrappedRow
		for j, cell := range row.TableCells {
			w := colWidths[j]
			wrapped := wrapCell(cell, w)
			wr.lines = append(wr.lines, wrapped)
			if len(wrapped) > wr.maxH {
				wr.maxH = len(wrapped)
			}
		}
		rows = append(rows, wr)
	}

	// Ensure all cells in a row have same number of lines
	for ri := range rows {
		for ci := 0; ci < len(colWidths); ci++ {
			for len(rows[ri].lines[ci]) < rows[ri].maxH {
				rows[ri].lines[ci] = append(rows[ri].lines[ci], "")
			}
		}
	}

	aligns := lines[0].TableAlign
	border := func(left, mid, right string) {
		sb.WriteString(sepStyle.Render(left))
		for j, w := range colWidths {
			sb.WriteString(sepStyle.Render(strings.Repeat("─", w+2)))
			if j < len(colWidths)-1 {
				sb.WriteString(sepStyle.Render(mid))
			}
		}
		sb.WriteString(sepStyle.Render(right))
		sb.WriteString("\n")
	}

	rowLines := func(rowIdx int, sty lipgloss.Style) {
		wr := rows[rowIdx]
		align := aligns
		if rowIdx > 0 && len(lines[rowIdx].TableAlign) > 0 {
			align = lines[rowIdx].TableAlign
		}
		for h := 0; h < wr.maxH; h++ {
			sb.WriteString("│")
			for j, w := range colWidths {
				text := ""
				if h < len(wr.lines[j]) {
					text = wr.lines[j][h]
				}
				var padded string
				al := AlignLeft
				if j < len(align) {
					al = align[j]
				}
				switch al {
				case AlignCenter:
					padded = padCellCenter(text, w)
				case AlignRight:
					padded = padCellRight(text, w)
				default:
					padded = padCell(text, w)
				}
				sb.WriteString(" ")
				sb.WriteString(sty.Render(padded))
				sb.WriteString(" ")
				if j < len(colWidths)-1 {
					sb.WriteString("│")
				}
			}
			sb.WriteString("│\n")
		}
	}

	border("┌", "┬", "┐")
	rowLines(0, headerStyle)
	border("├", "┼", "┤")
	for ri := 1; ri < len(lines); ri++ {
		if ri > 1 {
			// No separator between data rows
		}
		rowLines(ri, cellStyle)
	}
	border("└", "┴", "┘")

	return strings.TrimRight(sb.String(), "\n")
}

func renderBorderlessTable(lines []MarkdownLine, colWidths []int, style RendererStyle) string {
	var sb strings.Builder
	headerStyle := lipgloss.NewStyle().Foreground(style.Accent).Bold(true)
	cellStyle := lipgloss.NewStyle().Foreground(style.TextSecondary)
	aligns := lines[0].TableAlign

	type wrappedRow struct {
		lines [][]string
		maxH  int
	}
	var rows []wrappedRow
	for _, row := range lines {
		var wr wrappedRow
		for j, cell := range row.TableCells {
			w := colWidths[j]
			wrapped := wrapCell(cell, w)
			wr.lines = append(wr.lines, wrapped)
			if len(wrapped) > wr.maxH {
				wr.maxH = len(wrapped)
			}
		}
		for j := 0; j < len(colWidths); j++ {
			for len(wr.lines[j]) < wr.maxH {
				wr.lines[j] = append(wr.lines[j], "")
			}
		}
		rows = append(rows, wr)
	}

	for ri, wr := range rows {
		if ri > 0 {
			// Single blank line between data rows in borderless mode
		}
		for h := 0; h < wr.maxH; h++ {
			for j, w := range colWidths {
				text := ""
				if h < len(wr.lines[j]) {
					text = wr.lines[j][h]
				}
				al := AlignLeft
				if j < len(aligns) {
					al = aligns[j]
				}
				var padded string
				switch al {
				case AlignCenter:
					padded = padCellCenter(text, w)
				case AlignRight:
					padded = padCellRight(text, w)
				default:
					padded = padCell(text, w)
				}
				sty := cellStyle
				if ri == 0 {
					sty = headerStyle
				}
				if j > 0 {
					sb.WriteString("  ")
				}
				sb.WriteString(sty.Render(padded))
			}
			sb.WriteString("\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func renderSingleColumnTable(lines []MarkdownLine, style RendererStyle) string {
	var sb strings.Builder
	dimStyle := lipgloss.NewStyle().Foreground(style.TextDim)
	headerStyle := lipgloss.NewStyle().Foreground(style.Accent).Bold(true)

	for ri, row := range lines {
		if ri > 0 {
			sb.WriteString("\n")
		}
		if ri == 0 {
			// Header row: show column names as table title
			var names []string
			for _, cell := range row.TableCells {
				names = append(names, cell)
			}
			sb.WriteString(headerStyle.Render(strings.Join(names, " / ")))
			sb.WriteString("\n")
			sb.WriteString(dimStyle.Render(strings.Repeat("─", 40)))
			continue
		}
		for j, cell := range row.TableCells {
			name := ""
			if j < len(lines[0].TableCells) {
				name = lines[0].TableCells[j]
			}
			label := lipgloss.NewStyle().Foreground(style.Accent).Bold(true).Render(name)
			value := lipgloss.NewStyle().Foreground(style.TextSecondary).Render(cell)
			sb.WriteString(fmt.Sprintf("  %s  %s", label, value))
			if j < len(row.TableCells)-1 {
				sb.WriteString("\n")
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func padCell(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-len(runes))
}

func padCellCenter(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		return string(runes[:width])
	}
	left := (width - len(runes)) / 2
	right := width - len(runes) - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func padCellRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		return string(runes[:width])
	}
	return strings.Repeat(" ", width-len(runes)) + s
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

// HeadingInfo represents a heading extracted from markdown.
type HeadingInfo struct {
	Level   int
	Text    string
	LineIdx int
}

// ExtractHeadings extracts all headings from parsed markdown lines.
func ExtractHeadings(lines []MarkdownLine) []HeadingInfo {
	var headings []HeadingInfo
	for i, line := range lines {
		if line.BlockType == BlockHeading {
			text := RenderSegmentsPlain(line.Segments)
			headings = append(headings, HeadingInfo{
				Level:   line.HeadingLevel,
				Text:    text,
				LineIdx: i,
			})
		}
	}
	return headings
}

// RenderSegmentsPlain renders segments without styling (for outline).
func RenderSegmentsPlain(segments []InlineSegment) string {
	var sb strings.Builder
	for _, seg := range segments {
		sb.WriteString(seg.Text)
	}
	return sb.String()
}
