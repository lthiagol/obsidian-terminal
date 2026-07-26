package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

// MarkdownViewer renders and navigates markdown content.
type MarkdownViewer struct {
	viewport      viewport
	rawMarkdown   string
	links         []markdown.WikiLink
	selectedLink  int
	renderStyle   markdown.RendererStyle
	embedResolver markdown.EmbedResolver
}

// SetEmbedResolver sets the resolver for ![[embed]] directives.
func (v *MarkdownViewer) SetEmbedResolver(r markdown.EmbedResolver) {
	v.embedResolver = r
}

// NewViewer creates a MarkdownViewer with the given renderer style.
func NewViewer(style markdown.RendererStyle) MarkdownViewer {
	vp := newViewport(80, 20)
	return MarkdownViewer{
		viewport:     vp,
		selectedLink: -1,
		renderStyle:  style,
	}
}

func (v *MarkdownViewer) SetContent(md string, width int) {
	v.rawMarkdown = md
	v.viewport.Width = width - 2

	if strings.TrimSpace(md) == "" {
		v.viewport.SetContent("(empty note)")
		v.links = nil
		v.selectedLink = -1
		return
	}

	fmBlock := renderFrontmatter(md, v.viewport.Width-2, v.renderStyle)

	if strings.HasPrefix(md, "---\n") {
		afterFM := markdown.StripFrontmatter(md)
		if strings.TrimSpace(afterFM) == "" {
			if fmBlock != "" {
				v.viewport.SetContent(fmBlock)
			} else {
				v.viewport.SetContent("(empty note)")
			}
			v.links = nil
			v.selectedLink = -1
			return
		}
	}

	lines := markdown.ParseMarkdown(md)
	if v.embedResolver != nil {
		lines = markdown.ResolveEmbeds(lines, v.embedResolver)
	}
	rendered := markdown.RenderMarkdown(lines, v.viewport.Width-2, v.renderStyle)

	if fmBlock != "" {
		rendered = fmBlock + "\n" + rendered
	}

	v.viewport.SetContent(rendered)
	v.links = markdown.ExtractWikiLinks(lines)
	v.selectedLink = -1
}

func (v MarkdownViewer) View() string {
	return v.viewport.View()
}

func (v *MarkdownViewer) ScrollUp(n int) {
	v.viewport.LineUp(n)
}

func (v *MarkdownViewer) ScrollDown(n int) {
	v.viewport.LineDown(n)
}

func (v *MarkdownViewer) ScrollTop() {
	v.viewport.SetYOffset(0)
}

func (v *MarkdownViewer) ScrollBottom() {
	v.viewport.GotoBottom()
}

func (v *MarkdownViewer) GetScrollPosition() int {
	return v.viewport.YOffset
}

func (v *MarkdownViewer) SetScrollPosition(y int) {
	v.viewport.SetYOffset(y)
}

func (v *MarkdownViewer) ScrollHalfPageUp() {
	v.viewport.HalfViewUp()
}

func (v *MarkdownViewer) ScrollHalfPageDown() {
	v.viewport.HalfViewDown()
}

func (v *MarkdownViewer) SetSize(width, height int) {
	v.viewport.Width = max(width-2, 10)
	v.viewport.Height = max(height-2, 5)
}

func (v *MarkdownViewer) CycleLink() {
	if len(v.links) == 0 {
		v.selectedLink = -1
		return
	}
	v.selectedLink = (v.selectedLink + 1) % len(v.links)
}

func (v MarkdownViewer) SelectedLinkPath() string {
	if v.selectedLink < 0 || v.selectedLink >= len(v.links) {
		return ""
	}
	return v.links[v.selectedLink].Target
}

func (v MarkdownViewer) SelectedLinkDisplay() string {
	if v.selectedLink < 0 || v.selectedLink >= len(v.links) {
		return ""
	}
	return v.links[v.selectedLink].Display
}

// Width returns the current viewport content width.
func (v MarkdownViewer) Width() int {
	return v.viewport.Width
}

func (v MarkdownViewer) LinkCount() int {
	return len(v.links)
}

func (v MarkdownViewer) SelectedLinkIndex() int {
	return v.selectedLink
}

func renderFrontmatter(content string, width int, style markdown.RendererStyle) string {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return ""
	}

	rest := content[3:]
	endIdx := strings.Index(rest, "\n---\n")
	if endIdx < 0 {
		endIdx = strings.Index(rest, "\n---\r\n")
	}
	if endIdx < 0 && strings.HasSuffix(rest, "\n---") {
		endIdx = len(rest) - 4
	}
	if endIdx < 0 {
		return ""
	}

	yamlBlock := rest[:endIdx]
	if strings.TrimSpace(yamlBlock) == "" {
		return ""
	}

	var sb strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(style.AccentTertiary).
		Render("─── Frontmatter")
	sb.WriteString(header)
	sb.WriteString("\n")

	lines := strings.Split(yamlBlock, "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		colonIdx := -1
		for i, c := range line {
			if c == ':' {
				colonIdx = i
				break
			}
		}

		if colonIdx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		keyStyled := lipgloss.NewStyle().Foreground(style.Accent).Render(key)
		valueStyled := lipgloss.NewStyle().Foreground(style.TextSecondary).Render(value)

		sb.WriteString(fmt.Sprintf("  %s: %s", keyStyled, valueStyled))
		sb.WriteString("\n")
	}

	footer := lipgloss.NewStyle().
		Foreground(style.TextDim).
		Render(strings.Repeat("─", width))
	sb.WriteString(footer)

	return sb.String()
}
