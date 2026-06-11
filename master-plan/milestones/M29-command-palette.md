# M29 — Command Palette

**Status:** ⏳ pending

## Goal

Fuzzy command palette (Ctrl+K) listing all available actions with keybinding hints.

## Implementation Plan

### 1. Command registry

```go
type Command struct {
    Name, Description, Key string
    Action func(m *Model)
}
```

Register 15-20 commands: Toggle Help, Fuzzy Search, Content Search, Pin Note, Outline, Daily Note, Recent Notes, Backlinks, Tag Browser, Switch Profile, Force Rescan, Quit, mode-specific (Go Back, Follow Link, Expand/Collapse).

### 2. Palette state (model.go)

Fields: `paletteVisible bool`, `paletteQuery string`, `paletteResults []Command`, `paletteCursor int`.

Methods: `openCommandPalette()`, `filterPalette(q)`, `renderPalette()`, `executeSelectedCommand()`.

### 3. Handler

`handleCommandPaletteKey()`: Esc dismiss, Backspace/Runes filter, j/k move cursor, Enter execute+dismiss.

Reuse existing `FuzzyScore` for filtering commands.

### 4. Wiring

Global `Ctrl+K` in Update opens palette. Check `paletteVisible` before mode dispatch. View overlay: centered, half-width modal.

### Implementation order

1. Define Command struct + register commands
2. Add palette fields to Model
3. Implement open/filter/render/execute methods
4. Add handleCommandPaletteKey
5. Wire Ctrl+K + dispatch
6. Add tests
