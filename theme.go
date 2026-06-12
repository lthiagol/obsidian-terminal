package main

import (
	"fmt"
	"strings"

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
	SelectionText   lipgloss.Color

	ModeBrowse     lipgloss.Color
	ModeView       lipgloss.Color
	ModeSearch     lipgloss.Color
	ModeFind       lipgloss.Color
	ModeHelp       lipgloss.Color
	ModeTags       lipgloss.Color
	ModeProfile    lipgloss.Color
	Heading1       lipgloss.Color

	TreeStyle   lipgloss.Style
	ViewerStyle lipgloss.Style
	StatusStyle lipgloss.Style
	HelpStyle   lipgloss.Style
	SearchStyle lipgloss.Style
}

type themeDef struct {
	Colors map[string]string
}

var themeData = map[string]themeDef{
	"dark": {
		Colors: map[string]string{
			"accent":           "#a78bfa",
			"accent_secondary": "#fbbf24",
			"accent_tertiary":  "#2dd4bf",
			"text_primary":     "#e5e7eb",
			"text_secondary":   "#9ca3af",
			"text_muted":       "#6b7280",
			"text_dim":         "#4b5563",
			"success":          "#34d399",
			"warning":          "#fbbf24",
			"error":            "#f87171",
			"info":             "#60a5fa",
			"background":       "#111827",
			"surface":          "#1f2937",
			"border":           "#374151",
			"selection_text":   "#000000",
			"mode_browse":      "#a78bfa",
			"mode_view":        "#e879f9",
			"mode_search":      "#fbbf24",
			"mode_find":        "#fbbf24",
			"mode_help":        "#60a5fa",
			"mode_tags":        "#fb923c",
			"mode_profile":     "#a78bfa",
			"heading1":         "#e879f9",
		},
	},
	"catppuccin-latte": {
		Colors: map[string]string{
			"accent":           "#8839ef",
			"accent_secondary": "#7287fd",
			"accent_tertiary":  "#1e66f5",
			"text_primary":     "#4c4f69",
			"text_secondary":   "#6c6f85",
			"text_muted":       "#9ca0b0",
			"text_dim":         "#8c8fa1",
			"success":          "#40a02b",
			"warning":          "#df8e1d",
			"error":            "#d20f39",
			"info":             "#1e66f5",
			"background":       "#eff1f5",
			"surface":          "#ccd0da",
			"border":           "#acb0be",
			"mode_browse":      "#8839ef",
			"mode_view":        "#ea76cb",
			"mode_search":      "#7287fd",
			"mode_find":        "#7287fd",
			"mode_help":        "#1e66f5",
			"mode_tags":        "#fe640b",
			"mode_profile":     "#8839ef",
			"heading1":         "#ea76cb",
		},
	},
	"catppuccin-frappe": {
		Colors: map[string]string{
			"accent":           "#ca9ee6",
			"accent_secondary": "#babbf1",
			"accent_tertiary":  "#8caaee",
			"text_primary":     "#c6d0f5",
			"text_secondary":   "#a5adce",
			"text_muted":       "#949cbb",
			"text_dim":         "#838ba7",
			"success":          "#a6d189",
			"warning":          "#e5c890",
			"error":            "#e78284",
			"info":             "#8caaee",
			"background":       "#303446",
			"surface":          "#414559",
			"border":           "#626880",
			"mode_browse":      "#ca9ee6",
			"mode_view":        "#f4b8e4",
			"mode_search":      "#babbf1",
			"mode_find":        "#babbf1",
			"mode_help":        "#8caaee",
			"mode_tags":        "#ef9f76",
			"mode_profile":     "#ca9ee6",
			"heading1":         "#f4b8e4",
		},
	},
	"catppuccin-macchiato": {
		Colors: map[string]string{
			"accent":           "#c6a0f6",
			"accent_secondary": "#b7bdf8",
			"accent_tertiary":  "#8aadf4",
			"text_primary":     "#cad3f5",
			"text_secondary":   "#a5adcb",
			"text_muted":       "#939ab7",
			"text_dim":         "#8087a2",
			"success":          "#a6da95",
			"warning":          "#eed49f",
			"error":            "#ed8796",
			"info":             "#8aadf4",
			"background":       "#24273a",
			"surface":          "#363a4f",
			"border":           "#5b6078",
			"mode_browse":      "#c6a0f6",
			"mode_view":        "#f5bde6",
			"mode_search":      "#b7bdf8",
			"mode_find":        "#b7bdf8",
			"mode_help":        "#8aadf4",
			"mode_tags":        "#f5a97f",
			"mode_profile":     "#c6a0f6",
			"heading1":         "#f5bde6",
		},
	},
	"catppuccin-mocha": {
		Colors: map[string]string{
			"accent":           "#cba6f7",
			"accent_secondary": "#b4befe",
			"accent_tertiary":  "#89b4fa",
			"text_primary":     "#cdd6f4",
			"text_secondary":   "#a6adc8",
			"text_muted":       "#9399b2",
			"text_dim":         "#7f849c",
			"success":          "#a6e3a1",
			"warning":          "#f9e2af",
			"error":            "#f38ba8",
			"info":             "#89b4fa",
			"background":       "#1e1e2e",
			"surface":          "#313244",
			"border":           "#585b70",
			"mode_browse":      "#cba6f7",
			"mode_view":        "#f5c2e7",
			"mode_search":      "#b4befe",
			"mode_find":        "#b4befe",
			"mode_help":        "#89b4fa",
			"mode_tags":        "#fab387",
			"mode_profile":     "#cba6f7",
			"heading1":         "#f5c2e7",
		},
	},
	"dracula": {
		Colors: map[string]string{
			"accent":           "#bd93f9",
			"accent_secondary": "#ff79c6",
			"accent_tertiary":  "#8be9fd",
			"text_primary":     "#f8f8f2",
			"text_secondary":   "#6272a4",
			"text_muted":       "#6272a4",
			"text_dim":         "#44475a",
			"success":          "#50fa7b",
			"warning":          "#ffb86c",
			"error":            "#ff5555",
			"info":             "#8be9fd",
			"background":       "#282a36",
			"surface":          "#44475a",
			"border":           "#6272a4",
			"mode_browse":      "#bd93f9",
			"mode_view":        "#ff79c6",
			"mode_search":      "#f1fa8c",
			"mode_find":        "#f1fa8c",
			"mode_help":        "#8be9fd",
			"mode_tags":        "#ffb86c",
			"mode_profile":     "#bd93f9",
			"heading1":         "#ff79c6",
		},
	},
	"alucard": {
		Colors: map[string]string{
			"accent":           "#b390e3",
			"accent_secondary": "#e377a9",
			"accent_tertiary":  "#6cb6c5",
			"text_primary":     "#c2c2c9",
			"text_secondary":   "#576284",
			"text_muted":       "#576284",
			"text_dim":         "#45495a",
			"success":          "#43945f",
			"warning":          "#d08735",
			"error":            "#d94e5d",
			"info":             "#6cb6c5",
			"background":       "#1e2029",
			"surface":          "#323540",
			"border":           "#576284",
			"mode_browse":      "#b390e3",
			"mode_view":        "#e377a9",
			"mode_search":      "#f0c062",
			"mode_find":        "#f0c062",
			"mode_help":        "#6cb6c5",
			"mode_tags":        "#d08735",
			"mode_profile":     "#b390e3",
			"heading1":         "#e377a9",
		},
	},
}

func buildPalette(name string) (Palette, error) {
	def, ok := themeData[name]
	if !ok {
		return Palette{}, fmt.Errorf("unknown theme %q", name)
	}
	p := Palette{
		Accent:          parseHex(def.Colors["accent"]),
		AccentSecondary: parseHex(def.Colors["accent_secondary"]),
		AccentTertiary:  parseHex(def.Colors["accent_tertiary"]),
		TextPrimary:     parseHex(def.Colors["text_primary"]),
		TextSecondary:   parseHex(def.Colors["text_secondary"]),
		TextMuted:       parseHex(def.Colors["text_muted"]),
		TextDim:         parseHex(def.Colors["text_dim"]),
		Success:         parseHex(def.Colors["success"]),
		Warning:         parseHex(def.Colors["warning"]),
		Error:           parseHex(def.Colors["error"]),
		Info:            parseHex(def.Colors["info"]),
		Background:      parseHex(def.Colors["background"]),
		Surface:         parseHex(def.Colors["surface"]),
		Border:          parseHex(def.Colors["border"]),
		SelectionText:   parseHexOrDefault(def.Colors["selection_text"], "#000000"),
		Heading1:        parseHexOrDefault(def.Colors["heading1"], "#e879f9"),
	}
	return rebuildDerivedStyles(p), nil
}

func parseHex(s string) lipgloss.Color {
	return lipgloss.Color(s)
}

func parseHexOrDefault(s, def string) lipgloss.Color {
	if s == "" {
		return lipgloss.Color(def)
	}
	return lipgloss.Color(s)
}

// ValidThemeNames returns a list of all available theme names.
func ValidThemeNames() []string {
	var names []string
	for name := range themeData {
		names = append(names, name)
	}
	return names
}

func lookupPalette(name string) (Palette, error) {
	return buildPalette(name)
}

// newDarkPalette is a test helper that returns the dark theme palette.
func newDarkPalette() Palette {
	p, err := buildPalette("dark")
	if err != nil {
		p = rebuildDerivedStyles(Palette{
			Accent: "#a78bfa", AccentSecondary: "#fbbf24", AccentTertiary: "#2dd4bf",
			TextPrimary: "#e5e7eb", TextSecondary: "#9ca3af", TextMuted: "#6b7280",
			TextDim: "#4b5563", Success: "#34d399", Warning: "#fbbf24",
			Error: "#f87171", Info: "#60a5fa", Background: "#111827",
			Surface: "#1f2937", Border: "#374151", Heading1: "#e879f9",
		})
	}
	return p
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
)

var SelectionText = lipgloss.Color("#000000")

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
	ModeBrowse:        Accent,
	ModeView:          AccentSecondary,
	ModeSearch:        AccentSecondary,
	ModeFind:          AccentSecondary,
	ModeHelp:          Info,
	ModeTags:          AccentSecondary,
	ModeProfilePicker: Accent,
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
	SelectionText = p.SelectionText
	TreeStyle = p.TreeStyle
	ViewerStyle = p.ViewerStyle
	StatusStyle = p.StatusStyle
	HelpStyle = p.HelpStyle
	SearchStyle = p.SearchStyle
	ModeColors = map[Mode]lipgloss.Color{
		ModeBrowse:        p.ModeBrowse,
		ModeView:          p.ModeView,
		ModeSearch:        p.ModeSearch,
		ModeFind:          p.ModeFind,
		ModeHelp:          p.ModeHelp,
		ModeTags:          p.ModeTags,
		ModeProfilePicker: p.ModeProfile,
	}
}

func markdownStyleFrom(p Palette, lineSpacing string) markdown.RendererStyle {
	return markdown.RendererStyle{
		Accent:          p.Accent,
		AccentSecondary: p.AccentSecondary,
		AccentTertiary:  p.AccentTertiary,
		TextSecondary:   p.TextSecondary,
		TextDim:         p.TextDim,
		Success:         p.Success,
		CodeBackground:  p.Surface,
		Heading1:        p.Heading1,
		LineSpacing:     lineSpacing,
	}
}

func searchStyleFrom(p Palette) search.Style {
	return search.Style{
		Accent:        p.Accent,
		TextSecondary: p.TextSecondary,
		TextMuted:     p.TextMuted,
		SelectionText: p.SelectionText,
	}
}

// parseHexColor validates and parses a hex color string.
// Accepts #RGB or #RRGGBB format, returns normalized lowercase.
func parseHexColor(s string) (lipgloss.Color, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("empty color string")
	}
	if s[0] != '#' {
		return "", fmt.Errorf("color must start with #")
	}

	hex := s[1:]
	switch len(hex) {
	case 3:
		// #RGB -> #RRGGBB
		r, g, b := hex[0], hex[1], hex[2]
		hex = string([]byte{r, r, g, g, b, b})
	case 6:
		// #RRGGBB - already correct
	default:
		return "", fmt.Errorf("invalid hex color length: %d (expected 3 or 6)", len(hex))
	}

	// Validate hex characters
	for _, c := range hex {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return "", fmt.Errorf("invalid hex character: %c", c)
		}
	}

	// Normalize to lowercase
	return lipgloss.Color("#" + strings.ToLower(hex)), nil
}

// paletteFromCustom applies custom theme overrides to a base palette.
// Returns the modified palette and any error encountered.
func paletteFromCustom(ct *CustomTheme, base Palette) (Palette, error) {
	if ct == nil {
		return base, nil
	}

	p := base
	var errors []string

	if ct.Accent != "" {
		if c, err := parseHexColor(ct.Accent); err == nil {
			p.Accent = c
		} else {
			errors = append(errors, "accent: "+err.Error())
		}
	}
	if ct.AccentSecondary != "" {
		if c, err := parseHexColor(ct.AccentSecondary); err == nil {
			p.AccentSecondary = c
		} else {
			errors = append(errors, "accent_secondary: "+err.Error())
		}
	}
	if ct.AccentTertiary != "" {
		if c, err := parseHexColor(ct.AccentTertiary); err == nil {
			p.AccentTertiary = c
		} else {
			errors = append(errors, "accent_tertiary: "+err.Error())
		}
	}
	if ct.TextPrimary != "" {
		if c, err := parseHexColor(ct.TextPrimary); err == nil {
			p.TextPrimary = c
		} else {
			errors = append(errors, "text_primary: "+err.Error())
		}
	}
	if ct.TextSecondary != "" {
		if c, err := parseHexColor(ct.TextSecondary); err == nil {
			p.TextSecondary = c
		} else {
			errors = append(errors, "text_secondary: "+err.Error())
		}
	}
	if ct.TextMuted != "" {
		if c, err := parseHexColor(ct.TextMuted); err == nil {
			p.TextMuted = c
		} else {
			errors = append(errors, "text_muted: "+err.Error())
		}
	}
	if ct.TextDim != "" {
		if c, err := parseHexColor(ct.TextDim); err == nil {
			p.TextDim = c
		} else {
			errors = append(errors, "text_dim: "+err.Error())
		}
	}
	if ct.Success != "" {
		if c, err := parseHexColor(ct.Success); err == nil {
			p.Success = c
		} else {
			errors = append(errors, "success: "+err.Error())
		}
	}
	if ct.Warning != "" {
		if c, err := parseHexColor(ct.Warning); err == nil {
			p.Warning = c
		} else {
			errors = append(errors, "warning: "+err.Error())
		}
	}
	if ct.Error != "" {
		if c, err := parseHexColor(ct.Error); err == nil {
			p.Error = c
		} else {
			errors = append(errors, "error: "+err.Error())
		}
	}
	if ct.Info != "" {
		if c, err := parseHexColor(ct.Info); err == nil {
			p.Info = c
		} else {
			errors = append(errors, "info: "+err.Error())
		}
	}
	if ct.Background != "" {
		if c, err := parseHexColor(ct.Background); err == nil {
			p.Background = c
		} else {
			errors = append(errors, "background: "+err.Error())
		}
	}
	if ct.Surface != "" {
		if c, err := parseHexColor(ct.Surface); err == nil {
			p.Surface = c
		} else {
			errors = append(errors, "surface: "+err.Error())
		}
	}
	if ct.Border != "" {
		if c, err := parseHexColor(ct.Border); err == nil {
			p.Border = c
		} else {
			errors = append(errors, "border: "+err.Error())
		}
	}

	// Rebuild derived styles after applying custom colors
	p = rebuildDerivedStyles(p)

	if len(errors) > 0 {
		return p, fmt.Errorf("custom theme errors: %s", strings.Join(errors, "; "))
	}

	return p, nil
}

// rebuildDerivedStyles rebuilds TreeStyle, ViewerStyle, StatusStyle, HelpStyle,
// SearchStyle, and ModeColors from the palette colors.
func rebuildDerivedStyles(p Palette) Palette {
	p.TreeStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(p.Accent).
		Padding(0, 1)

	p.ViewerStyle = lipgloss.NewStyle().
		Padding(0, 1)

	p.StatusStyle = lipgloss.NewStyle().
		Background(p.Surface).
		Padding(0, 1)

	p.HelpStyle = lipgloss.NewStyle().
		Padding(1, 2)

	p.SearchStyle = lipgloss.NewStyle().
		Padding(1, 2)

	p.ModeBrowse = p.Accent
	p.ModeView = p.Heading1
	p.ModeSearch = p.AccentSecondary
	p.ModeFind = p.AccentSecondary
	p.ModeHelp = p.Info
	p.ModeTags = p.AccentSecondary
	p.ModeProfile = p.Accent

	return p
}
