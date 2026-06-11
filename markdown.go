package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

type MarkdownLine struct {
	BlockType    BlockType
	HeadingLevel int
	Segments     []InlineSegment
	IndentLevel  int
	CalloutType  string
	Language     string
	RawContent   string
}

type WikiLink struct {
	Target  string
	Display string
}

type Styles struct {
	Width int
}

func ParseMarkdown(content string) []MarkdownLine {
	content = stripMarkdownFrontmatter(content)
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
					BlockType: BlockCodeBlock,
					Language:  codeLang,
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
	re := regexp.MustCompile(`\[!(\w+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return strings.ToLower(matches[1])
	}
	return "note"
}

func stripBlockquote(line string) string {
	t := strings.TrimLeft(line, " ")
	t = strings.TrimPrefix(t, ">")
	t = strings.TrimPrefix(t, " ")
	re := regexp.MustCompile(`\[!\w+\][+-]?\s*`)
	t = re.ReplaceAllString(t, "")
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
	return regexp.MustCompile(`^[\-\*\+]\s`).MatchString(strings.TrimLeft(t, " ")) ||
		regexp.MustCompile(`^\d+\.\s`).MatchString(strings.TrimLeft(t, " "))
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
	re := regexp.MustCompile(`^([\-\*\+]|\d+\.)\s+`)
	loc := re.FindStringIndex(t)
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

	// Wiki-links [[target]] or [[target|display]]
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

	// Bold+Italic *** or bold **
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

	// Italic * or _
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

	// Inline code `
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

	// Strikethrough ~~
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

	// Highlight ==
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

	// Plain text up to next special character
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
	delims := []string{"[[", "***", "___", "**", "__", "*", "_", "`", "~~", "=="}
	earliest := -1
	for _, d := range delims {
		i := strings.Index(text, d)
		if i >= 0 && (earliest == -1 || i < earliest) {
			earliest = i
		}
	}
	return earliest
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

func stripMarkdownFrontmatter(content string) string {
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
	re := regexp.MustCompile(`%%.*?%%`)
	return re.ReplaceAllString(content, "")
}

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

func ResolveWikiLink(target string, vault *VaultEntry, vaultRoot string) string {
	if target == "" {
		return ""
	}

	target = strings.SplitN(target, "#", 2)[0]

	if !strings.HasSuffix(target, ".md") {
		exact := findExactPath(vault, "", target+".md")
		if exact != "" {
			return exact
		}
		exact = findExactPath(vault, "", target+".markdown")
		if exact != "" {
			return exact
		}
	} else {
		exact := findExactPath(vault, "", target)
		if exact != "" {
			return exact
		}
	}

	basename := strings.ToLower(target)
	if !strings.HasSuffix(basename, ".md") {
		basename += ".md"
	}
	result := findBasename(vault, "", basename)
	if result != "" {
		return result
	}

	result = findAlias(vault, target, vaultRoot)
	if result != "" {
		return result
	}

	return ""
}

func findAlias(vault *VaultEntry, alias string, vaultRoot string) string {
	aliasLower := strings.ToLower(alias)
	for _, child := range vault.Children {
		if child.IsDir {
			found := findAlias(child, alias, vaultRoot)
			if found != "" {
				return found
			}
			continue
		}
		aliasEntries, _ := extractAliasesFromFile(vaultRoot, child.Path)
		for _, a := range aliasEntries {
			if strings.ToLower(a) == aliasLower {
				return child.Path
			}
		}
	}
	return ""
}

func extractAliasesFromFile(vaultRoot, relativePath string) ([]string, error) {
	fullPath := relativePath
	if vaultRoot != "" {
		fullPath = filepath.Join(vaultRoot, relativePath)
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	fm, _ := parseFrontmatter(content)
	return fm.Aliases, nil
}

func findExactPath(vault *VaultEntry, prefix, target string) string {
	for _, child := range vault.Children {
		childPath := child.Path
		if prefix != "" {
			childPath = prefix + "/" + child.Name
		}
		if childPath == target && !child.IsDir {
			return childPath
		}
		if child.IsDir {
			found := findExactPath(child, childPath, target)
			if found != "" {
				return found
			}
		}
	}
	return ""
}

func findBasename(vault *VaultEntry, prefix, target string) string {
	targetLower := strings.ToLower(target)
	for _, child := range vault.Children {
		childPath := child.Path
		if prefix != "" {
			childPath = prefix + "/" + child.Name
		}
		nameLower := strings.ToLower(child.Name)
		if nameLower == targetLower && !child.IsDir {
			return childPath
		}
		if child.IsDir {
			found := findBasename(child, childPath, target)
			if found != "" {
				return found
			}
		}
	}
	return ""
}

func RenderMarkdown(lines []MarkdownLine, width int) string {
	if width < 20 {
		width = 20
	}

	var sb strings.Builder
	for _, line := range lines {
		rendered := renderLine(line, width)
		if rendered != "" {
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(rendered)
		}
	}
	return sb.String()
}

func renderLine(line MarkdownLine, width int) string {
	switch line.BlockType {
	case BlockHeading:
		return renderHeading(line, width)
	case BlockCodeBlock:
		return renderCodeBlock(line, width)
	case BlockList:
		return renderList(line, width)
	case BlockBlockquote:
		return renderBlockquote(line, width)
	case BlockCallout:
		return renderCallout(line, width)
	case BlockHorizontalRule:
		return renderHorizontalRule(width)
	case BlockEmpty:
		return ""
	default:
		return renderParagraph(line, width)
	}
}

func renderHeading(line MarkdownLine, width int) string {
	text := renderSegments(line.Segments)
	var style lipgloss.Style
	switch line.HeadingLevel {
	case 1:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#f472b6")).Bold(true).Underline(true)
	case 2:
		style = lipgloss.NewStyle().Foreground(Accent).Bold(true)
	case 3:
		style = lipgloss.NewStyle().Foreground(AccentTertiary).Bold(true)
	default:
		style = lipgloss.NewStyle().Foreground(TextSecondary).Bold(true)
	}
	return style.Render(text)
}

func renderCodeBlock(line MarkdownLine, width int) string {
	lines := strings.Split(line.RawContent, "\n")

	header := ""
	if line.Language != "" {
		header = " " + line.Language + " "
	}

	var sb strings.Builder

	topBorder := lipgloss.NewStyle().Foreground(TextDim).Render("╭" + strings.Repeat("─", width-2) + "╮")
	if header != "" {
		labelStyle := lipgloss.NewStyle().Foreground(TextDim)
		padded := header + strings.Repeat("─", width-len(header)-2)
		topBorder = lipgloss.NewStyle().Foreground(TextDim).Render("╭") +
			labelStyle.Render(padded) +
			lipgloss.NewStyle().Foreground(TextDim).Render("╮")
	}
	sb.WriteString(topBorder)

	codeStyle := lipgloss.NewStyle().Foreground(Success)

	for _, l := range lines {
		sb.WriteString("\n")
		l = strings.ReplaceAll(l, "\t", "    ")
		if len(l) > width-2 {
			l = l[:width-2]
		}
		padded := l + strings.Repeat(" ", width-2-len(l))
		lineContent := lipgloss.NewStyle().Foreground(TextDim).Render("│") +
			codeStyle.Render(padded) +
			lipgloss.NewStyle().Foreground(TextDim).Render("│")
		sb.WriteString(lineContent)
	}

	sb.WriteString("\n")
	botBorder := lipgloss.NewStyle().Foreground(TextDim).Render("╰" + strings.Repeat("─", width-2) + "╯")
	sb.WriteString(botBorder)

	return sb.String()
}

func renderList(line MarkdownLine, width int) string {
	prefix := strings.Repeat("  ", line.IndentLevel)
	text := renderSegments(line.Segments)
	bullet := lipgloss.NewStyle().Foreground(Accent).Render("•")
	return prefix + bullet + " " + text
}

func renderBlockquote(line MarkdownLine, width int) string {
	prefix := lipgloss.NewStyle().Foreground(Accent).Render("│")
	text := renderSegments(line.Segments)
	bodyStyle := lipgloss.NewStyle().Foreground(TextSecondary).Italic(true)
	return prefix + " " + bodyStyle.Render(text)
}

func renderCallout(line MarkdownLine, width int) string {
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

	typeStyle := lipgloss.NewStyle().Bold(true).Foreground(AccentSecondary)
	bodyStyle := lipgloss.NewStyle().Foreground(TextSecondary)

	return icon + " " + typeStyle.Render(line.CalloutType) + " " + bodyStyle.Render(line.Segments[0].Text)
}

func renderHorizontalRule(width int) string {
	rule := strings.Repeat("─", width)
	return lipgloss.NewStyle().Foreground(TextDim).Render(rule)
}

func renderParagraph(line MarkdownLine, width int) string {
	text := renderSegments(line.Segments)
	return wrapText(text, width)
}

func renderSegments(segments []InlineSegment) string {
	var sb strings.Builder
	for _, seg := range segments {
		sb.WriteString(renderSegment(seg))
	}
	return sb.String()
}

func renderSegment(seg InlineSegment) string {
	style := lipgloss.NewStyle()

	switch {
	case seg.IsWikiLink:
		style = style.Foreground(AccentTertiary).Underline(true)
	case seg.Code:
		style = style.Foreground(Success).Background(lipgloss.Color("#1f2937"))
	case seg.Highlight:
		style = style.Foreground(AccentSecondary)
	case seg.Strikethrough:
		style = style.Strikethrough(true)
	case seg.Bold && seg.Italic:
		style = style.Bold(true).Italic(true)
	case seg.Bold:
		style = style.Bold(true)
	case seg.Italic:
		style = style.Italic(true)
	}

	return style.Render(seg.Text)
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
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	clean := re.ReplaceAllString(s, "")
	return len([]rune(clean))
}
