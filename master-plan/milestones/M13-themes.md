# M13 — Theme System & Color Palettes

**Status:** ✅ done

## Goal

Implement a theme system that supports multiple color palettes, starting with Catppuccin (4 flavors) and Dracula/Alucard, making `Config.Theme` actually functional.

## Files to modify

- `theme.go` — restructure for multiple themes, add palette constants
- `config.go` — validate theme string, wire up theme selection
- `model.go` — pass theme through to all renderers

## Steps

### 1. Define a `Theme` struct

Instead of bare `lipgloss.Color` package-level vars, define:

```go
type Palette struct {
    Accent          lipgloss.Color
    AccentSecondary lipgloss.Color
    TextPrimary     lipgloss.Color
    TextSecondary   lipgloss.Color
    TextDim         lipgloss.Color
    Success         lipgloss.Color
    Warning         lipgloss.Color
    Error           lipgloss.Color
    Info            lipgloss.Color
    Background      lipgloss.Color
    Surface         lipgloss.Color
    Border          lipgloss.Color

    // Mode-specific badge colors
    ModeBrowse lipgloss.Color
    ModeView   lipgloss.Color
    ModeSearch lipgloss.Color
    ModeFind   lipgloss.Color
    ModeHelp   lipgloss.Color

    // Tree styles
    TreeStyle      lipgloss.Style
    ViewerStyle    lipgloss.Style
    StatusStyle    lipgloss.Style
    SelectedStyle  lipgloss.Style
}
```

### 2. Add Catppuccin flavors

**Latte** (light):
| Semantic | Catppuccin Name | Hex |
|----------|-----------------|-----|
| Background | Base | `#eff1f5` |
| Surface | Surface0 | `#ccd0da` |
| TextPrimary | Text | `#4c4f69` |
| TextSecondary | Subtext0 | `#6c6f85` |
| TextDim | Overlay1 | `#8c8fa1` |
| Accent | Mauve | `#8839ef` |
| AccentSecondary | Lavender | `#7287fd` |
| Success | Green | `#40a02b` |
| Warning | Yellow | `#df8e1d` |
| Error | Red | `#d20f39` |
| Info | Blue | `#1e66f5` |
| Border | Surface2 | `#acb0be` |

**Frappé** (dark):
| Semantic | Catppuccin Name | Hex |
|----------|-----------------|-----|
| Background | Base | `#303446` |
| Surface | Surface0 | `#414559` |
| TextPrimary | Text | `#c6d0f5` |
| TextSecondary | Subtext0 | `#a5adce` |
| TextDim | Overlay1 | `#838ba7` |
| Accent | Mauve | `#ca9ee6` |
| AccentSecondary | Lavender | `#babbf1` |
| Success | Green | `#a6d189` |
| Warning | Yellow | `#e5c890` |
| Error | Red | `#e78284` |
| Info | Blue | `#8caaee` |
| Border | Surface2 | `#626880` |

**Macchiato** (dark):
| Semantic | Catppuccin Name | Hex |
|----------|-----------------|-----|
| Background | Base | `#24273a` |
| Surface | Surface0 | `#363a4f` |
| TextPrimary | Text | `#cad3f5` |
| TextSecondary | Subtext0 | `#a5adcb` |
| TextDim | Overlay1 | `#8087a2` |
| Accent | Mauve | `#c6a0f6` |
| AccentSecondary | Lavender | `#b7bdf8` |
| Success | Green | `#a6da95` |
| Warning | Yellow | `#eed49f` |
| Error | Red | `#ed8796` |
| Info | Blue | `#8aadf4` |
| Border | Surface2 | `#5b6078` |

**Mocha** (dark):
| Semantic | Catppuccin Name | Hex |
|----------|-----------------|-----|
| Background | Base | `#1e1e2e` |
| Surface | Surface0 | `#313244` |
| TextPrimary | Text | `#cdd6f4` |
| TextSecondary | Subtext0 | `#a6adc8` |
| TextDim | Overlay1 | `#7f849c` |
| Accent | Mauve | `#cba6f7` |
| AccentSecondary | Lavender | `#b4befe` |
| Success | Green | `#a6e3a1` |
| Warning | Yellow | `#f9e2af` |
| Error | Red | `#f38ba8` |
| Info | Blue | `#89b4fa` |
| Border | Surface2 | `#585b70` |

### 3. Add Dracula

**Dracula** (dark):
| Semantic | Dracula Name | Hex |
|----------|-------------|-----|
| Background | Background | `#282a36` |
| Surface | CurrentLine | `#44475a` |
| TextPrimary | Foreground | `#f8f8f2` |
| TextSecondary | Comment | `#6272a4` (dimmed) |
| TextDim | Selection | `#44475a` (dimmed) |
| Accent | Purple | `#bd93f9` |
| AccentSecondary | Pink | `#ff79c6` |
| Success | Green | `#50fa7b` |
| Warning | Orange | `#ffb86c` |
| Error | Red | `#ff5555` |
| Info | Cyan | `#8be9fd` |
| Border | Comment | `#6272a4` |

**Alucard** (dark, vampire-themed variant):
Alucard is Dracula's darker cousin with deeper backgrounds.
| Semantic | Alucard Name | Hex |
|----------|-------------|-----|
| Background | Background | `#1e2029` |
| Surface | CurrentLine | `#323540` |
| TextPrimary | Foreground | `#c2c2c9` |
| TextSecondary | Comment | `#576284` |
| TextDim | — | `#45495a` |
| Accent | Purple | `#b390e3` |
| AccentSecondary | Pink | `#e377a9` |
| Success | Green | `#43945f` |
| Warning | Orange | `#d08735` |
| Error | Red | `#d94e5d` |
| Info | Cyan | `#6cb6c5` |
| Border | Comment | `#576284` |

### 4. Store active theme in Model

- Add `palette Palette` field to `Model`
- Initialize from `Config.Theme` in `NewModel`
- Pass palette to `FileTree`, `MarkdownViewer`, and all render methods

### 5. Make renderers use palette

- All current hardcoded `Accent`, `TextPrimary`, etc. package-level vars → use `m.palette.Accent`, etc.
- `View()`, `renderSearch()`, `renderFind()`, `renderHelp()`, `renderStatusBar()`, `renderToasts()` — all need the palette
- Tree rendering (`tree.go`) needs the palette for folder/file colors
- Viewer (`viewer.go`) needs the palette, or passes it through to markdown renderer

### 6. Config validation

- `Config.Theme` accepts: `"catppuccin-latte"`, `"catppuccin-frappe"`, `"catppuccin-macchiato"`, `"catppuccin-mocha"`, `"dracula"`, `"alucard"`, `"dark"` (existing default)
- Unknown theme → fall back to `"dark"` with a warning toast

### 7. Keep backward compatibility

- Existing `Config.Theme = "dark"` must still work as before (map to current hardcoded colors)
- Package-level color vars remain for backward compat (can be served by `DefaultPalette()`)

## Completion Criteria

- [ ] `Theme` struct and `Palette` struct defined in `theme.go`
- [ ] Catppuccin Latte, Frappé, Macchiato, Mocha palettes implemented
- [ ] Dracula palette implemented
- [ ] Alucard palette implemented
- [ ] All renderers use palette from model instead of hardcoded color vars
- [ ] Config validates theme and falls back on unknown values
- [ ] `"dark"` theme still works (backward compat)
- [ ] New tests for theme lookup and fallback
- [ ] `make test` passes all existing tests
- [ ] `make vet` exits 0
