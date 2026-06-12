# M37 — Theme System Refactor

**Status:** ⏳ pending

## Goal

Eliminate global mutable state in the theme system and fix the broken profile switching. Convert 16 global variables to model fields, fix applyProfile to use pointer receiver, and consolidate 7 duplicate palette constructors into a data-driven approach.

## Issues

### C1: Data race on global theme variables (`theme.go:413-436`)

`activatePalette()` writes to 16 global variables (10 colors + 5 styles + 1 map) at runtime when switching profiles. Meanwhile `View()`, `renderStatusBar()`, `renderHelp()`, `FileTree.View()` read these globals. While Bubble Tea's Update/View cycle is single-threaded, profile switching triggers `rescanVault()` which may race.

**Current globals (theme.go:364-411):**
- 10 color vars: `Accent`, `AccentSecondary`, `AccentTertiary`, `TextSecondary`, `TextMuted`, `TextDim`, `Success`, `Warning`, `Error`, `Info`
- 5 style vars: `TreeStyle`, `ViewerStyle`, `StatusStyle`, `HelpStyle`, `SearchStyle`
- 1 map: `ModeColors`

**Fix:** Convert all global style variables to fields on the `Model` struct (or a `ThemeState` sub-struct). All rendering code reads from the model instead of globals.

### C6: `applyProfile` mutates a copy (`handlers.go:433-470`)

`applyProfile` has a value receiver `(m Model)` but mutates `m.config.VaultPath`, `m.config.Theme`, `m.palette`, etc. Since it returns `tea.Msg`, the mutations happen on a copy and are discarded. Profile switching via the profile picker is silently broken.

**Fix:** Change `applyProfile` to a pointer receiver `(m *Model)` and return `nil` command. Or return a `tea.Batch` that triggers rescan. Ensure the palette and styles update on the real model, not a copy.

### H4: 7 palette constructors with duplicate structure

`theme.go` has `newDarkPalette()`, `newCatppuccinLatte()`, `newCatppuccinFrappe()`, `newCatppuccinMacchiato()`, `newCatppuccinMocha()`, `newDracula()`, `newAlucard()`. Each is ~40 lines with identical structure — only the hex color values differ.

**Fix:** Replace with data-driven theme definitions using a map of color values.

## Design

### Step 1: Move globals to Model

Add a `ThemeState` sub-struct to Model:

```go
type ThemeState struct {
    Accent          lipgloss.Color
    AccentSecondary lipgloss.Color
    // ... all 10 colors
    TreeStyle       lipgloss.Style
    ViewerStyle     lipgloss.Style
    // ... all 5 styles
    ModeColors      map[Mode]lipgloss.Color
}

type Model struct {
    // ... existing fields
    theme ThemeState
}
```

Update all rendering functions to read from `m.theme` instead of globals.

### Step 2: Fix applyProfile

Change signature from `func (m Model) applyProfile() tea.Msg` to `func (m *Model) applyProfile()`. Ensure all mutations happen on the real model.

### Step 3: Data-driven themes

Replace 7 constructor functions with a single builder:

```go
type themeDef struct {
    Name   string
    Colors map[string]string
}

var themes = map[string]themeDef{
    "dark": {Name: "dark", Colors: darkColors},
    "catppuccin-latte": {Name: "catppuccin-latte", Colors: catppuccinLatteColors},
    // ...
}

func buildPalette(def themeDef) Palette {
    p := Palette{Name: def.Name}
    setColor(&p.Accent, def.Colors, "accent")
    // ...
    rebuildDerivedStyles(&p)
    return p
}
```

## Files to modify

| File | Changes |
|------|---------|
| `theme.go` | C1: remove 16 global vars; H4: replace 7 constructors with data map + builder |
| `model.go` | C1: add ThemeState field; update all rendering to use m.theme |
| `handlers.go` | C6: change applyProfile to pointer receiver |
| `statusbar.go`, `help.go`, `tree.go`, `backlinks.go`, `tags.go`, `command_palette.go`, `profile_picker.go`, `toast.go` | C1: update to read from model.theme instead of globals |
| `custom_theme_test.go` | Update palette tests to use new API |

## Completion Criteria

- [ ] No global mutable style variables — all rendering reads from model fields
- [ ] Profile theme switching works and reflects immediately
- [ ] applyProfile uses pointer receiver and works correctly
- [ ] 7 palette constructors replaced by data-driven approach
- [ ] Adding a new theme requires only a new key+color map entry
- [ ] `rebuildDerivedStyles` is the single source of truth for style construction
- [ ] All existing tests pass with no behavioral changes
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

2-3 days (major refactor touching 10+ files)
