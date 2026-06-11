package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

// Color palette variables for the dark theme.
var (
	Accent          = lipgloss.Color("#a78bfa")
	AccentSecondary = lipgloss.Color("#fbbf24")
	AccentTertiary  = lipgloss.Color("#2dd4bf")
	TextSecondary   = lipgloss.Color("#9ca3af")
	TextMuted       = lipgloss.Color("#6b7280")
	TextDim         = lipgloss.Color("#4b5563")
	Success         = lipgloss.Color("#34d399")
	Warning         = lipgloss.Color("#fbbf24")
	Error           = lipgloss.Color("#f87171")
	Info            = lipgloss.Color("#60a5fa")
)

// Unicode icons for the file tree.
var (
	IconFolderOpen   = "▾ "
	IconFolderClosed = "▸ "
	IconFile         = "◇ "
	IconVertical     = "│"
	IconDiamond      = "◆"
)

// Pre-defined lipgloss styles for the TUI layout.
var (
	TreeStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(Accent).
			Padding(0, 1)

	ViewerStyle = lipgloss.NewStyle().
			Padding(0, 1)

	StatusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1f2937")).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Padding(1, 2)

	SearchStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

// ModeColors maps each TUI mode to its accent color.
var ModeColors = map[Mode]lipgloss.Color{
	ModeBrowse: Accent,
	ModeView:   AccentTertiary,
	ModeSearch: AccentSecondary,
	ModeFind:   AccentSecondary,
	ModeHelp:   Info,
}

func defaultMarkdownStyle() markdown.RendererStyle {
	return markdown.RendererStyle{
		Accent:          Accent,
		AccentSecondary: AccentSecondary,
		AccentTertiary:  AccentTertiary,
		TextSecondary:   TextSecondary,
		TextDim:         TextDim,
		Success:         Success,
		CodeBackground:  lipgloss.Color("#1f2937"),
		Heading1:        lipgloss.Color("#f472b6"),
	}
}

func defaultSearchStyle() search.Style {
	return search.Style{
		Accent:        Accent,
		TextSecondary: TextSecondary,
		TextMuted:     TextMuted,
	}
}
