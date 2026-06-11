package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

// MarkdownViewer renders and navigates markdown content.
type MarkdownViewer struct {
	viewport     viewport.Model
	rawMarkdown  string
	links        []markdown.WikiLink
	selectedLink int
	renderStyle  markdown.RendererStyle
}

// NewViewer creates a MarkdownViewer with the given renderer style.
func NewViewer(style markdown.RendererStyle) MarkdownViewer {
	vp := viewport.New(80, 20)
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

	if strings.HasPrefix(md, "---\n") {
		afterFM := markdown.StripFrontmatter(md)
		if strings.TrimSpace(afterFM) == "" {
			v.viewport.SetContent("(empty note)")
			v.links = nil
			v.selectedLink = -1
			return
		}
	}

	lines := markdown.ParseMarkdown(md)
	rendered := markdown.RenderMarkdown(lines, v.viewport.Width-2, v.renderStyle)
	v.viewport.SetContent(rendered)
	v.links = markdown.ExtractWikiLinks(lines)
	v.selectedLink = -1
}

func (v *MarkdownViewer) Update(msg tea.Msg) (MarkdownViewer, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)
	return *v, cmd
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

func (v MarkdownViewer) LinkCount() int {
	return len(v.links)
}

func (v MarkdownViewer) SelectedLinkIndex() int {
	return v.selectedLink
}
