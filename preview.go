package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

const previewMaxLines = 80

type PreviewPane struct {
	content string
	width   int
	height  int
}

func NewPreviewPane() PreviewPane {
	return PreviewPane{}
}

func (pp *PreviewPane) SetSize(width, height int) {
	pp.width = width
	pp.height = height
}

func (pp *PreviewPane) SetContent(body string, width int, palette Palette) {
	pp.width = width

	lines := markdown.ParseMarkdown(body)
	rendered := markdown.RenderMarkdown(lines, width, markdownStyleFrom(palette, "compact"))

	renderedLines := strings.Split(rendered, "\n")
	if len(renderedLines) > previewMaxLines {
		renderedLines = renderedLines[:previewMaxLines]
		renderedLines = append(renderedLines, lipgloss.NewStyle().Foreground(palette.TextDim).Italic(true).Render(
			fmt.Sprintf("… (%d more lines — Enter to open full note)", len(renderedLines)-previewMaxLines)))
	}

	pp.content = strings.Join(renderedLines, "\n")
}

func (m Model) renderPreview() string {
	width := m.width - m.treeWidth - 4
	if width < 20 {
		width = 20
	}
	height := m.height - 1

	entry := m.fileTree.SelectedEntry()
	if entry == nil || entry.IsDir {
		return lipgloss.NewStyle().
			Foreground(m.palette.TextMuted).
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Preview (v): select a markdown file")
	}

	note, err := LoadNote(m.config.VaultPath, entry.Path)
	if err != nil {
		return lipgloss.NewStyle().
			Foreground(m.palette.TextMuted).
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Could not load note")
	}

	var pp PreviewPane
	pp.SetContent(note.Body, width, m.palette)

	header := lipgloss.NewStyle().
		Foreground(m.palette.TextMuted).
		Render("── preview: " + entry.Path + " ──")

	content := pp.View(m.palette)

	return lipgloss.JoinVertical(lipgloss.Left, header, "", content)
}

func (pp PreviewPane) View(palette Palette) string {
	if pp.content == "" {
		return lipgloss.NewStyle().
			Foreground(palette.TextMuted).
			Width(pp.width).
			Height(pp.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Preview (v): select a markdown file")
	}

	return pp.content
}
