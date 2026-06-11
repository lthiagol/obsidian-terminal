# M22 — Custom Themes

**Status:** ⏳ pending

## Goal

Allow user-defined color overrides in config via `custom_theme` section. Builds on M13 palette system.

## Implementation Plan

### 1. Config changes (`config.go`)

```go
type CustomTheme struct {
    Accent, AccentSecondary, AccentTertiary string  // all hex colors
    TextPrimary, TextSecondary, TextMuted, TextDim string
    Success, Warning, Error, Info string
    Background, Surface, Border string
}
// Add CustomTheme *CustomTheme to Config
```

### 2. Theme changes (`theme.go`)

`parseHexColor(s string) (lipgloss.Color, error)` — validates `#RRGGBB` or `#RGB`, normalizes to lowercase.

`paletteFromCustom(ct *CustomTheme, base Palette) (Palette, error)` — for each non-empty field, parse hex and override base palette color. Any field unset keeps base value.

`rebuildDerivedStyles(p Palette) Palette` — rebuilds TreeStyle, ViewerStyle, StatusStyle, HelpStyle, SearchStyle, ModeColors from palette colors. Used after custom overrides applied.

### 3. Wire into NewModel (`model.go`)

After `lookupPalette(themeName)`, if `cfg.CustomTheme != nil`, call `paletteFromCustom(cfg.CustomTheme, palette)`. On error: show warning toast, fall back to base palette.

### Edge cases

- `custom_theme: {}` → empty non-nil pointer → all colors skipped → base unchanged
- 3-digit hex `#abc` → accepted
- Invalid hex → error toast, base palette used
- Missing fields → keep base palette values
- Theme + custom: `theme: dracula` + `custom_theme: {accent: "#ff0000"}` → dracula base with accent overridden

### Implementation order

1. Add CustomTheme struct to config.go
2. Add parseHexColor to theme.go
3. Add paletteFromCustom + rebuildDerivedStyles to theme.go
4. Wire into NewModel in model.go
5. Write tests
