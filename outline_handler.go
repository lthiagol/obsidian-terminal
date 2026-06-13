package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

func (m *Model) buildOutline() {
	if m.activeNote == nil {
		m.outlineItems = nil
		return
	}

	lines := markdown.ParseMarkdown(m.activeNote.RawBody)
	headings := markdown.ExtractHeadings(lines)

	m.outlineItems = make([]OutlineItem, len(headings))
	for i, h := range headings {
		m.outlineItems[i] = OutlineItem{
			Level:   h.Level,
			Text:    h.Text,
			LineIdx: h.LineIdx,
			YOffset: estimateYOffset(lines, h.LineIdx, m.viewer.Width()),
		}
	}

	m.outlineCursor = 0
}

func (m Model) renderOutline() string {
	if len(m.outlineItems) == 0 {
		return lipgloss.NewStyle().
			Foreground(m.palette.TextMuted).
			Render("  No headings in this note")
	}

	var sb strings.Builder
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.palette.Accent).
		Render(fmt.Sprintf("  Outline (%d)", len(m.outlineItems)))
	sb.WriteString(header)
	sb.WriteString("\n")

	for i, item := range m.outlineItems {
		indent := strings.Repeat("  ", item.Level-1)
		line := fmt.Sprintf("%s%s", indent, item.Text)

		if i == m.outlineCursor {
			line = lipgloss.NewStyle().
				Background(m.palette.Accent).
				Foreground(m.palette.SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(m.palette.TextSecondary).
				Render(line)
		}

		sb.WriteString(line)
		if i < len(m.outlineItems)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func estimateYOffset(lines []markdown.MarkdownLine, targetIdx, width int) int {
	yOffset := 0
	for i := 0; i < targetIdx && i < len(lines); i++ {
		line := lines[i]
		switch line.BlockType {
		case markdown.BlockEmpty:
			yOffset++
		case markdown.BlockHeading:
			yOffset++
		case markdown.BlockCodeBlock:
			codeLines := strings.Count(line.RawContent, "\n") + 1
			yOffset += codeLines + 2
		case markdown.BlockList:
			yOffset++
		case markdown.BlockBlockquote:
			yOffset++
		case markdown.BlockCallout:
			yOffset++
		case markdown.BlockHorizontalRule:
			yOffset++
		default:
			text := markdown.RenderSegmentsPlain(line.Segments)
			if width > 0 {
				runeCount := len([]rune(text))
				wrappedLines := (runeCount / width) + 1
				yOffset += wrappedLines
			} else {
				yOffset++
			}
		}
	}
	return yOffset
}

func (m Model) handleOutlineKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc || MatchRune(msg, m.keys.Outline):
		m.outlineVisible = false
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		if m.outlineCursor < len(m.outlineItems)-1 {
			m.outlineCursor++
		}
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		if m.outlineCursor > 0 {
			m.outlineCursor--
		}
		return m, nil
	case msg.Type == tea.KeyEnter:
		if m.outlineCursor < len(m.outlineItems) {
			item := m.outlineItems[m.outlineCursor]
			m.viewer.ScrollTop()
			m.viewer.ScrollDown(item.YOffset)
			m.outlineVisible = false
		}
		return m, nil
	}
	return m, nil
}
