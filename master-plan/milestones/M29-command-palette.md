# M29 — Command Palette

**Status:** ✅ done

## Goal

Fuzzy command palette (Ctrl+K) listing all available actions with keybinding hints.

## Keybinding

**Key:** `Ctrl+K`
**Mode:** Global (all modes)
**Rationale:** Standard command palette keybinding (VSCode, JetBrains)

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Implementation Plan

### 1. Command registry

```go
type Command struct {
    Name        string
    Description string
    Key         string  // keybinding hint, e.g., "/" or "Ctrl+D"
    Action      func(m *Model) (tea.Model, tea.Cmd)
    Modes       []Mode  // modes where command is available (empty = all)
}
```

Register commands by category:

**Navigation:**
- Go Back (Esc) — View mode
- Follow Link (Enter) — View mode
- Expand/Collapse (l/h) — Browse mode

**Search:**
- Fuzzy Search (/) — Browse, View modes
- Content Search (s) — Browse, View modes

**Features:**
- Toggle Help (?) — all modes
- Pin Note (p) — Browse, View modes
- Outline (t) — View mode
- Daily Note (Ctrl+D) — all modes
- Recent Notes (Ctrl+O) — all modes
- Backlinks (b) — View mode
- Tag Browser (T) — Browse mode
- Switch Profile (P) — Browse mode

**System:**
- Force Rescan (Ctrl+R) — all modes
- Quit (q) — all modes

### 2. Palette state (model.go)

Fields: `paletteVisible bool`, `paletteQuery string`, `paletteResults []Command`, `paletteCursor int`.

Methods: `openCommandPalette()`, `filterPalette(q)`, `renderPalette()`, `executeSelectedCommand()`.

### 3. Handler

`handleCommandPaletteKey()`: Esc dismiss, Backspace/Runes filter, j/k move cursor, Enter execute+dismiss.

Reuse existing `FuzzyScore` for filtering commands.

### 4. Mode filtering

Only show commands available in current mode. Check `Command.Modes` against `m.mode`.

### 5. Wiring

Global `Ctrl+K` in Update opens palette. Check `paletteVisible` before mode dispatch. View overlay: centered, half-width modal.

### Implementation order

1. Define Command struct with Modes field
2. Register all commands with appropriate modes
3. Add palette fields to Model
4. Implement open/filter/render/execute methods
5. Add handleCommandPaletteKey
6. Wire Ctrl+K + dispatch
7. Add tests

## Completion Criteria

- [ ] Command struct with Name, Description, Key, Action, Modes
- [ ] 15-20 commands registered
- [ ] Commands filtered by current mode
- [ ] `Ctrl+K` opens command palette overlay
- [ ] Fuzzy filtering works as user types
- [ ] j/k navigate results
- [ ] Enter executes command and dismisses palette
- [ ] Esc dismisses without executing
- [ ] Keybinding hints shown for each command
- [ ] Overlay is centered, half-width modal
- [ ] `make test` passes
- [ ] `make vet` exits 0
- [ ] Manual test: command palette works end-to-end
