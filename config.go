package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Profile represents a vault profile with its own settings.
type Profile struct {
	Path     string   // Path is the vault directory for this profile.
	Theme    string   // Theme is the theme name for this profile.
	SkipDirs []string // SkipDirs are directory names to exclude from scanning.
}

// CustomTheme represents user-defined color overrides.
type CustomTheme struct {
	Accent          string // Accent is the primary accent color (hex).
	AccentSecondary string // AccentSecondary is a secondary accent color (hex).
	AccentTertiary  string // AccentTertiary is a tertiary accent color (hex).
	TextPrimary     string // TextPrimary is the main text color (hex).
	TextSecondary   string // TextSecondary is the secondary text color (hex).
	TextMuted       string // TextMuted is the muted text color (hex).
	TextDim         string // TextDim is the dimmed text color (hex).
	Success         string // Success is the color for success indicators (hex).
	Warning         string // Warning is the color for warning indicators (hex).
	Error           string // Error is the color for error indicators (hex).
	Info            string // Info is the color for informational text (hex).
	Background      string // Background is the terminal background color (hex).
	Surface         string // Surface is the elevated surface color (hex).
	Border          string // Border is the border/separator color (hex).
}

// Config holds user configuration loaded from YAML.
type Config struct {
	VaultPath        string            // VaultPath is the root directory of the Obsidian vault.
	Theme            string            // Theme is the active theme name.
	DefaultKeys      string            // DefaultKeys is the keybinding preset ("vim" or "arrow").
	SkipDirs         []string          // SkipDirs are directory names to exclude from scanning.
	DailyNotesDir    string            // DailyNotesDir is the directory for daily notes (default: "Journal").
	DailyNotesFormat string            // DailyNotesFormat is the date format for daily notes (default: "2006-01-02").
	LineSpacing      string            // LineSpacing controls paragraph spacing ("compact", "normal", or "relaxed").
	Profiles         map[string]Profile // Profiles maps profile names to their settings.
	CustomTheme      *CustomTheme      // CustomTheme holds user-defined color overrides.
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Theme:            "dark",
		DefaultKeys:      "vim",
		SkipDirs:         []string{".obsidian", ".git", ".trash", "node_modules", "archive"},
		DailyNotesDir:    "Journal",
		DailyNotesFormat: "2006-01-02",
		LineSpacing:      "compact",
	}
}

// LoadConfig reads and parses a YAML config file at path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	cfg := DefaultConfig()
	parseConfigYAML(data, cfg)

	return cfg, nil
}

func parseConfigYAML(data []byte, cfg *Config) {
	scanYAML(data, func(key, value string, items []string) {
		switch key {
		case "vault_path":
			if value != "" {
				cfg.VaultPath = value
			}
		case "theme":
			if value != "" {
				cfg.Theme = value
			}
		case "default_keys":
			if value != "" {
				cfg.DefaultKeys = value
			}
		case "skip_dirs":
			if len(items) > 0 {
				cfg.SkipDirs = items
			} else if value != "" {
				cfg.SkipDirs = []string{value}
			}
		case "daily_notes_dir":
			if value != "" {
				cfg.DailyNotesDir = value
			}
		case "daily_notes_format":
			if value != "" {
				cfg.DailyNotesFormat = value
			}
		case "line_spacing":
			if value != "" {
				cfg.LineSpacing = value
			}
		}
	})

	// Parse profiles (nested structure)
	profilesData := parseNestedMap(data, "profiles")
	if len(profilesData) > 0 {
		cfg.Profiles = make(map[string]Profile)
		for name, props := range profilesData {
			profile := Profile{}
			if path, ok := props["path"]; ok {
				profile.Path = path
			}
			if theme, ok := props["theme"]; ok {
				profile.Theme = theme
			}
			if skipDirs, ok := props["skip_dirs"]; ok {
				// Parse skip_dirs as inline array or single value
				if strings.HasPrefix(skipDirs, "[") {
					profile.SkipDirs = parseInlineArray(skipDirs)
				} else if skipDirs != "" {
					profile.SkipDirs = []string{skipDirs}
				}
			}
			cfg.Profiles[name] = profile
		}
	}

	// Parse custom_theme (flat structure)
	themeData := parseFlatMap(data, "custom_theme")
	if len(themeData) > 0 {
		cfg.CustomTheme = &CustomTheme{}
		if v, ok := themeData["accent"]; ok {
			cfg.CustomTheme.Accent = v
		}
		if v, ok := themeData["accent_secondary"]; ok {
			cfg.CustomTheme.AccentSecondary = v
		}
		if v, ok := themeData["accent_tertiary"]; ok {
			cfg.CustomTheme.AccentTertiary = v
		}
		if v, ok := themeData["text_primary"]; ok {
			cfg.CustomTheme.TextPrimary = v
		}
		if v, ok := themeData["text_secondary"]; ok {
			cfg.CustomTheme.TextSecondary = v
		}
		if v, ok := themeData["text_muted"]; ok {
			cfg.CustomTheme.TextMuted = v
		}
		if v, ok := themeData["text_dim"]; ok {
			cfg.CustomTheme.TextDim = v
		}
		if v, ok := themeData["success"]; ok {
			cfg.CustomTheme.Success = v
		}
		if v, ok := themeData["warning"]; ok {
			cfg.CustomTheme.Warning = v
		}
		if v, ok := themeData["error"]; ok {
			cfg.CustomTheme.Error = v
		}
		if v, ok := themeData["info"]; ok {
			cfg.CustomTheme.Info = v
		}
		if v, ok := themeData["background"]; ok {
			cfg.CustomTheme.Background = v
		}
		if v, ok := themeData["surface"]; ok {
			cfg.CustomTheme.Surface = v
		}
		if v, ok := themeData["border"]; ok {
			cfg.CustomTheme.Border = v
		}
	}
}

func configPathOrDefault(explicit string) string {
	if explicit != "" {
		return explicit
	}
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "obsidian-terminal", "config.yaml")
}

// ValidLineSpacing values for the line_spacing config field.
var ValidLineSpacing = []string{"compact", "normal", "relaxed"}

// ValidateConfig validates all configuration values, auto-fixing invalid ones
// where possible. Returns warnings for any issues found (intended for toast display).
func ValidateConfig(cfg *Config) []string {
	var warnings []string

	if cfg.Theme == "" {
		cfg.Theme = "dark"
	} else {
		valid := ValidThemeNames()
		if !stringInSlice(cfg.Theme, valid) {
			warnings = append(warnings, fmt.Sprintf(
				"Invalid theme %q — valid themes: %s",
				cfg.Theme, strings.Join(valid, ", "),
			))
			cfg.Theme = "dark"
		}
	}

	if cfg.LineSpacing == "" {
		cfg.LineSpacing = "compact"
	} else if !stringInSlice(cfg.LineSpacing, ValidLineSpacing) {
		warnings = append(warnings, fmt.Sprintf(
			"Invalid line_spacing %q — valid values: compact, normal, relaxed",
			cfg.LineSpacing,
		))
		cfg.LineSpacing = "compact"
	}

	if cfg.DailyNotesFormat != "" && !isValidDateFormat(cfg.DailyNotesFormat) {
		warnings = append(warnings, fmt.Sprintf(
			"Invalid daily_notes_format %q — using default 2006-01-02",
			cfg.DailyNotesFormat,
		))
		cfg.DailyNotesFormat = "2006-01-02"
	}

	if len(cfg.SkipDirs) == 0 {
		cfg.SkipDirs = DefaultConfig().SkipDirs
	}
	for _, d := range cfg.SkipDirs {
		if strings.Contains(d, string(filepath.Separator)) ||
			d == "" || d == "." || d == ".." {
			warnings = append(warnings, fmt.Sprintf(
				"Invalid skip_dirs entry %q — entries must be directory names, not paths",
				d,
			))
		}
	}

	if cfg.CustomTheme != nil {
		colorWarnings := validateCustomThemeColors(cfg.CustomTheme)
		warnings = append(warnings, colorWarnings...)
	}

	for name, profile := range cfg.Profiles {
		if profile.Path == "" {
			warnings = append(warnings, fmt.Sprintf(
				"Profile %q has no path — skipped", name,
			))
		}
	}

	return warnings
}

func stringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func isValidDateFormat(format string) bool {
	ref := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	formatted := ref.Format(format)
	parsed, err := time.Parse(format, formatted)
	if err != nil {
		return false
	}
	return parsed.Year() == ref.Year() &&
		parsed.Month() == ref.Month() &&
		parsed.Day() == ref.Day()
}

func validateCustomThemeColors(ct *CustomTheme) []string {
	var warnings []string

	fields := []struct {
		name  string
		value string
	}{
		{"accent", ct.Accent},
		{"accent_secondary", ct.AccentSecondary},
		{"accent_tertiary", ct.AccentTertiary},
		{"text_primary", ct.TextPrimary},
		{"text_secondary", ct.TextSecondary},
		{"text_muted", ct.TextMuted},
		{"text_dim", ct.TextDim},
		{"success", ct.Success},
		{"warning", ct.Warning},
		{"error", ct.Error},
		{"info", ct.Info},
		{"background", ct.Background},
		{"surface", ct.Surface},
		{"border", ct.Border},
	}

	for _, f := range fields {
		if f.value == "" {
			continue
		}
		if _, err := parseHexColor(f.value); err != nil {
			warnings = append(warnings, fmt.Sprintf(
				"Invalid custom_theme %s: %s", f.name, err.Error(),
			))
		}
	}

	return warnings
}
