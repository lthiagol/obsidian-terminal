package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var cachedHelpLines []string

func InvalidateHelpCache() {
	cachedHelpLines = nil
}

func buildHelpLines(p Palette) []string {
	if cachedHelpLines != nil {
		return cachedHelpLines
	}

	groups := []struct {
		title    string
		bindings []string
	}{
		{
			title: "Navigation",
			bindings: []string{
				"j / ↓  — move down",
				"k / ↑  — move up",
				"h / ←  — collapse / back",
				"l / →  — expand / forward",
				"g       — jump to top",
				"G       — jump to bottom",
				"PgUp    — page up",
				"PgDn    — page down",
			},
		},
		{
			title: "File Tree",
			bindings: []string{
			"Enter    — open note / toggle folder",
			"← →      — collapse / expand folder",
			"v        — toggle note preview",
			"T        — browse tags",
			"p        — pin/unpin note",
			"P        — switch profile",
			"Ctrl+[/] — cycle pinned notes",
			},
		},
		{
			title: "Viewer",
			bindings: []string{
				"j / k    — scroll down / up",
				"g / G    — top / bottom",
				"Tab      — cycle wiki-links",
				"Enter    — follow selected link",
				"b        — toggle backlinks",
				"t        — toggle outline",
				"p        — pin/unpin note",
				"Ctrl+[/] — cycle pinned notes",
				"h / Esc  — back to browse",
			},
		},
		{
			title: "Search",
			bindings: []string{
				"/  — fuzzy file name search",
				"s  — full-text content search",
				"Enter — open selected result",
				"Esc   — cancel search",
			},
		},
		{
			title: "Graph",
			bindings: []string{
				"Ctrl+G  — open graph (browse: global, view: local)",
				"j / k   — move between nodes",
				"Enter   — open selected note",
				"l       — toggle global ↔ local scope",
				"f       — focus local graph on selected node",
				"r       — rebuild graph from indexes",
				"Esc     — return to previous mode",
			},
		},
		{
			title: "Global",
			bindings: []string{
				"?      — toggle this help",
				"q      — quit",
				"Ctrl+D — open daily note",
				"Ctrl+G — open graph view",
				"Ctrl+O — recent notes",
				"Ctrl+K — command palette",
				"Ctrl+R — rescan vault",
			},
		},
	}

	lines := []string{
		lipgloss.NewStyle().Bold(true).Foreground(p.Accent).Render("obsidian-terminal — Keybindings"),
		"",
	}

	for _, g := range groups {
		header := lipgloss.NewStyle().Bold(true).Foreground(p.Accent).Render(g.title)
		lines = append(lines, header)
		for _, b := range g.bindings {
			parts := strings.SplitN(b, "—", 2)
			key := lipgloss.NewStyle().Foreground(p.AccentSecondary).Render(strings.TrimSpace(parts[0]))
			var desc string
			if len(parts) > 1 {
				desc = lipgloss.NewStyle().Foreground(p.TextSecondary).Render("—" + parts[1])
			}
			lines = append(lines, "  "+key+"  "+desc)
		}
		lines = append(lines, "")
	}

	cachedHelpLines = lines
	return lines
}

func (m Model) renderHelp() string {
	lines := buildHelpLines(m.palette)

	if m.helpScroll > len(lines)-1 {
		m.helpScroll = len(lines) - 1
	}

	start := m.helpScroll
	end := start + (m.height - 5)
	if end > len(lines) {
		end = len(lines)
	}
	if start >= len(lines) {
		start = 0
	}

	return strings.Join(lines[start:end], "\n")
}
