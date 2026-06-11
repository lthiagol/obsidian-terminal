package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

// Palette holds a complete set of UI colors for a theme.
type Palette struct {
	Accent          lipgloss.Color
	AccentSecondary lipgloss.Color
	AccentTertiary  lipgloss.Color
	TextPrimary     lipgloss.Color
	TextSecondary   lipgloss.Color
	TextMuted       lipgloss.Color
	TextDim         lipgloss.Color
	Success         lipgloss.Color
	Warning         lipgloss.Color
	Error           lipgloss.Color
	Info            lipgloss.Color
	Background      lipgloss.Color
	Surface         lipgloss.Color
	Border          lipgloss.Color

	ModeBrowse lipgloss.Color
	ModeView   lipgloss.Color
	ModeSearch lipgloss.Color
	ModeFind   lipgloss.Color
	ModeHelp   lipgloss.Color

	TreeStyle   lipgloss.Style
	ViewerStyle lipgloss.Style
	StatusStyle lipgloss.Style
	HelpStyle   lipgloss.Style
	SearchStyle lipgloss.Style
}

func newDarkPalette() Palette {
	return Palette{
		Accent:          "#a78bfa",
		AccentSecondary: "#fbbf24",
		AccentTertiary:  "#2dd4bf",
		TextPrimary:     "#e5e7eb",
		TextSecondary:   "#9ca3af",
		TextMuted:       "#6b7280",
		TextDim:         "#4b5563",
		Success:         "#34d399",
		Warning:         "#fbbf24",
		Error:           "#f87171",
		Info:            "#60a5fa",
		Background:      "#111827",
		Surface:         "#1f2937",
		Border:          "#374151",

		ModeBrowse: "#a78bfa",
		ModeView:   "#2dd4bf",
		ModeSearch: "#fbbf24",
		ModeFind:   "#fbbf24",
		ModeHelp:   "#60a5fa",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#a78bfa")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#1f2937")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func newCatppuccinLatte() Palette {
	return Palette{
		Accent:          "#8839ef",
		AccentSecondary: "#7287fd",
		AccentTertiary:  "#1e66f5",
		TextPrimary:     "#4c4f69",
		TextSecondary:   "#6c6f85",
		TextMuted:       "#9ca0b0",
		TextDim:         "#8c8fa1",
		Success:         "#40a02b",
		Warning:         "#df8e1d",
		Error:           "#d20f39",
		Info:            "#1e66f5",
		Background:      "#eff1f5",
		Surface:         "#ccd0da",
		Border:          "#acb0be",

		ModeBrowse: "#8839ef",
		ModeView:   "#1e66f5",
		ModeSearch: "#7287fd",
		ModeFind:   "#7287fd",
		ModeHelp:   "#1e66f5",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#8839ef")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#ccd0da")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func newCatppuccinFrappe() Palette {
	return Palette{
		Accent:          "#ca9ee6",
		AccentSecondary: "#babbf1",
		AccentTertiary:  "#8caaee",
		TextPrimary:     "#c6d0f5",
		TextSecondary:   "#a5adce",
		TextMuted:       "#949cbb",
		TextDim:         "#838ba7",
		Success:         "#a6d189",
		Warning:         "#e5c890",
		Error:           "#e78284",
		Info:            "#8caaee",
		Background:      "#303446",
		Surface:         "#414559",
		Border:          "#626880",

		ModeBrowse: "#ca9ee6",
		ModeView:   "#8caaee",
		ModeSearch: "#babbf1",
		ModeFind:   "#babbf1",
		ModeHelp:   "#8caaee",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#ca9ee6")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#414559")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func newCatppuccinMacchiato() Palette {
	return Palette{
		Accent:          "#c6a0f6",
		AccentSecondary: "#b7bdf8",
		AccentTertiary:  "#8aadf4",
		TextPrimary:     "#cad3f5",
		TextSecondary:   "#a5adcb",
		TextMuted:       "#939ab7",
		TextDim:         "#8087a2",
		Success:         "#a6da95",
		Warning:         "#eed49f",
		Error:           "#ed8796",
		Info:            "#8aadf4",
		Background:      "#24273a",
		Surface:         "#363a4f",
		Border:          "#5b6078",

		ModeBrowse: "#c6a0f6",
		ModeView:   "#8aadf4",
		ModeSearch: "#b7bdf8",
		ModeFind:   "#b7bdf8",
		ModeHelp:   "#8aadf4",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#c6a0f6")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#363a4f")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func newCatppuccinMocha() Palette {
	return Palette{
		Accent:          "#cba6f7",
		AccentSecondary: "#b4befe",
		AccentTertiary:  "#89b4fa",
		TextPrimary:     "#cdd6f4",
		TextSecondary:   "#a6adc8",
		TextMuted:       "#9399b2",
		TextDim:         "#7f849c",
		Success:         "#a6e3a1",
		Warning:         "#f9e2af",
		Error:           "#f38ba8",
		Info:            "#89b4fa",
		Background:      "#1e1e2e",
		Surface:         "#313244",
		Border:          "#585b70",

		ModeBrowse: "#cba6f7",
		ModeView:   "#89b4fa",
		ModeSearch: "#b4befe",
		ModeFind:   "#b4befe",
		ModeHelp:   "#89b4fa",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#cba6f7")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#313244")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func newDracula() Palette {
	return Palette{
		Accent:          "#bd93f9",
		AccentSecondary: "#ff79c6",
		AccentTertiary:  "#8be9fd",
		TextPrimary:     "#f8f8f2",
		TextSecondary:   "#6272a4",
		TextMuted:       "#6272a4",
		TextDim:         "#44475a",
		Success:         "#50fa7b",
		Warning:         "#ffb86c",
		Error:           "#ff5555",
		Info:            "#8be9fd",
		Background:      "#282a36",
		Surface:         "#44475a",
		Border:          "#6272a4",

		ModeBrowse: "#bd93f9",
		ModeView:   "#8be9fd",
		ModeSearch: "#ff79c6",
		ModeFind:   "#ff79c6",
		ModeHelp:   "#8be9fd",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#bd93f9")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#44475a")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func newAlucard() Palette {
	return Palette{
		Accent:          "#b390e3",
		AccentSecondary: "#e377a9",
		AccentTertiary:  "#6cb6c5",
		TextPrimary:     "#c2c2c9",
		TextSecondary:   "#576284",
		TextMuted:       "#576284",
		TextDim:         "#45495a",
		Success:         "#43945f",
		Warning:         "#d08735",
		Error:           "#d94e5d",
		Info:            "#6cb6c5",
		Background:      "#1e2029",
		Surface:         "#323540",
		Border:          "#576284",

		ModeBrowse: "#b390e3",
		ModeView:   "#6cb6c5",
		ModeSearch: "#e377a9",
		ModeFind:   "#e377a9",
		ModeHelp:   "#6cb6c5",

		TreeStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#b390e3")).
			Padding(0, 1),

		ViewerStyle: lipgloss.NewStyle().
			Padding(0, 1),

		StatusStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#323540")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Padding(1, 2),

		SearchStyle: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

func lookupPalette(name string) (Palette, error) {
	switch name {
	case "dark":
		return newDarkPalette(), nil
	case "catppuccin-latte":
		return newCatppuccinLatte(), nil
	case "catppuccin-frappe":
		return newCatppuccinFrappe(), nil
	case "catppuccin-macchiato":
		return newCatppuccinMacchiato(), nil
	case "catppuccin-mocha":
		return newCatppuccinMocha(), nil
	case "dracula":
		return newDracula(), nil
	case "alucard":
		return newAlucard(), nil
	default:
		return Palette{}, fmt.Errorf("unknown theme %q", name)
	}
}

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

var (
	IconFolderOpen   = "▾ "
	IconFolderClosed = "▸ "
	IconFile         = "◇ "
	IconVertical     = "│"
	IconDiamond      = "◆"
)

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

var ModeColors = map[Mode]lipgloss.Color{
	ModeBrowse: Accent,
	ModeView:   AccentTertiary,
	ModeSearch: AccentSecondary,
	ModeFind:   AccentSecondary,
	ModeHelp:   Info,
}

func activatePalette(p Palette) {
	Accent = p.Accent
	AccentSecondary = p.AccentSecondary
	AccentTertiary = p.AccentTertiary
	TextSecondary = p.TextSecondary
	TextMuted = p.TextMuted
	TextDim = p.TextDim
	Success = p.Success
	Warning = p.Warning
	Error = p.Error
	Info = p.Info
	TreeStyle = p.TreeStyle
	ViewerStyle = p.ViewerStyle
	StatusStyle = p.StatusStyle
	HelpStyle = p.HelpStyle
	SearchStyle = p.SearchStyle
	ModeColors = map[Mode]lipgloss.Color{
		ModeBrowse: p.ModeBrowse,
		ModeView:   p.ModeView,
		ModeSearch: p.ModeSearch,
		ModeFind:   p.ModeFind,
		ModeHelp:   p.ModeHelp,
	}
}

func markdownStyleFrom(p Palette) markdown.RendererStyle {
	return markdown.RendererStyle{
		Accent:          p.Accent,
		AccentSecondary: p.AccentSecondary,
		AccentTertiary:  p.AccentTertiary,
		TextSecondary:   p.TextSecondary,
		TextDim:         p.TextDim,
		Success:         p.Success,
		CodeBackground:  p.Surface,
		Heading1:        p.AccentSecondary,
	}
}

func searchStyleFrom(p Palette) search.Style {
	return search.Style{
		Accent:        p.Accent,
		TextSecondary: p.TextSecondary,
		TextMuted:     p.TextMuted,
	}
}
