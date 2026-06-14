package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type BacklinkPanel struct {
	links   []string
	cursor  int
	width   int
	palette Palette
}

func NewBacklinkPanel(notePath string, backlinkIndex map[string][]string, palette Palette) BacklinkPanel {
	normalized := normalizeWikiLinkTarget(strings.TrimSuffix(notePath, filepath.Ext(notePath)))
	links := backlinkIndex[normalized]

	bp := BacklinkPanel{
		links:   links,
		cursor:  0,
		palette: palette,
	}
	return bp
}

func (bp *BacklinkPanel) MoveUp() {
	if bp.cursor > 0 {
		bp.cursor--
	}
}

func (bp *BacklinkPanel) MoveDown() {
	if bp.cursor < len(bp.links)-1 {
		bp.cursor++
	}
}

func (bp BacklinkPanel) SelectedPath() string {
	if bp.cursor >= 0 && bp.cursor < len(bp.links) {
		return bp.links[bp.cursor]
	}
	return ""
}

func (bp BacklinkPanel) Count() int {
	return len(bp.links)
}

func (bp BacklinkPanel) View() string {
	if len(bp.links) == 0 {
		return lipgloss.NewStyle().
			Foreground(bp.palette.TextMuted).
			Render("  No backlinks")
	}

	var sb strings.Builder
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(bp.palette.Accent).
		Render(fmt.Sprintf("  Backlinks (%d)", len(bp.links)))
	sb.WriteString(header)
	sb.WriteString("\n")

	for i, link := range bp.links {
		line := fmt.Sprintf("  %s", link)
		if i == bp.cursor {
			line = lipgloss.NewStyle().
				Background(bp.palette.Accent).
				Foreground(bp.palette.SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(bp.palette.TextSecondary).
				Render(line)
		}
		sb.WriteString(line)
		if i < len(bp.links)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
