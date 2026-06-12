# M37 — Refactor Theme System

**Status:** ⏳ pending

## Goal

Replace 7 near-identical palette constructors (~400 lines) with data-driven theme definitions. Eliminate redundant color definitions. Consolidate duplicate inline style construction (the 7 identical `TreeStyle`, `ViewerStyle` blocks).

## Issues

### H4: 7 palette constructors with duplicate structure

`theme.go` has `newDarkPalette()`, `newCatppuccinLatte()`, `newCatppuccinFrappe()`, `newCatppuccinMacchiato()`, `newCatppuccinMocha()`, `newDracula()`, `newAlucard()`. Each is ~40 lines with identical structure — only the hex color values differ.

### M5: Palette struct has 25 fields

The `Palette` struct carries 25 color fields + 7 lipgloss style fields. Many fields are only used in one place (`ModeFind`, `ModeProfilePicker`, `ModeBacklinks`, etc.). The struct is a catch-all.

### H1 (partial): Inline style construction duplicated per theme

Each palette constructor builds `TreeStyle`, `ViewerStyle`, `StatusStyle`, `HelpStyle`, `SearchStyle` inline with the same pattern but different colors. When a new style is needed (e.g., `TreeStyle` border changes), it must be added to 7 places.

## Design

### Data-driven themes

Replace the 7 constructor functions with a single `themeData` struct and a map:

```go
type themeDef struct {
    Name   string            // e.g., "dark", "dracula"
    Colors map[string]string // e.g., "accent": "#a78bfa", ...
}

var themes = map[string]themeDef{
    "dark":               {Name: "dark", Colors: darkColors},
    "catppuccin-latte":   {Name: "catppuccin-latte", Colors: catppuccinLatteColors},
    // ...
}
```

```go
func buildPalette(def themeDef) Palette {
    p := Palette{Name: def.Name}
    setColor(&p.Accent, def.Colors, "accent")
    setColor(&p.AccentSecondary, def.Colors, "accent_secondary")
    // ...
    rebuildDerivedStyles(p)  // existing function
    return p
}
```

This reduces ~400 lines to ~150 lines of data + one builder function.

### Unified style construction

Move `TreeStyle`, `ViewerStyle`, `StatusStyle`, `HelpStyle`, `SearchStyle` construction from each palette constructor into a single `rebuildDerivedStyles` function that takes a `Palette` and populates the style fields. This already exists at `theme.go:612` — ensure all palette constructors use it.

## Files to modify

| File | Changes |
|------|---------|
| `theme.go` | Replace 7 constructors with data map + `buildPalette`; consolidate style construction |
| `custom_theme_test.go` | Update palette tests to use new API |
| `model.go` | Update palette initialization calls |
| `DESIGN.md` | Update theme system architecture section |

## Completion Criteria

- [ ] 7 constructor functions replaced by data-driven approach
- [ ] Adding a new theme requires only a new key+color map entry
- [ ] All existing tests pass with no behavioral changes
- [ ] `rebuildDerivedStyles` is the single source of truth for style construction
- [ ] Palette struct has no unused fields
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
