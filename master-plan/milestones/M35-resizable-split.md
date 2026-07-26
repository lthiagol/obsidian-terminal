# M35 — Resizable Tree/Viewer Split

**Status:** ⏳ pending

## Goal

The tree panel can be resized at runtime by dragging the split border with the mouse or using keyboard shortcuts. Users with wide monitors can expand the tree to see long filenames; users with tables or wide content can shrink the tree to give the viewer more room.

## Problem

The tree width is fixed at `width/4` (min 25). Users can't adjust it without resizing the terminal. On wide terminals (120+ cols), the tree wastes space at 30 cols while the viewer has 88. On narrow terminals (80 cols), the tree eats 25 cols leaving only 53 for the viewer — often too narrow for tables.

## Solution

### Mouse drag

Bubble Tea sends `tea.MouseMsg` with `X`, `Y`, and `Button` (press, release, wheel). We detect mouse press near the split boundary and track drag motion.

```
┌──────────────┬─────────────────────────────────┐
│  notes/      │  # Title                        │
│  projects/   │  Some text here...               │
│  readme.md   │  And more content...             │
│  index.md    │                                  │
│              │                                  │
│         ↑ drag this border ↑                    │
└──────────────┴─────────────────────────────────┘
   tree (adjustable)    viewer (takes remainder)
```

**Mouse states to track:**
- `MouseDown` at `X ≈ treeWidth` (±2 chars): set `dragSplit = true`, capture offset
- `MouseMotion` while `dragSplit`: `treeWidth = max(15, min(msg.X - offset, width/2))`
- `MouseUp`: set `dragSplit = false`

**Visual feedback:** While dragging, the split border renders in `Accent` instead of `TextDim` to show it's active.

### Keyboard resize

| Keybinding | Action |
|-----------|--------|
| `Ctrl+Left` / `Alt+h` | Shrink tree by 5 cols |
| `Ctrl+Right` / `Alt+l` | Grow tree by 5 cols |
| `Ctrl+0` | Reset to default (width/4) |

Works in Browse mode only (tree is visible). Clamps between 15 and `width/2`.

### Layout recalculation

Whenever `treeWidth` changes (from drag or keys), the model re-computes:

```go
func (m *Model) adjustTreeWidth(newWidth int) {
    m.treeWidth = clamp(newWidth, 15, m.width/2)
    m.fileTree.SetSize(m.treeWidth, m.height-1)
    viewerWidth := m.width - m.treeWidth - 2
    m.viewer.SetSize(viewerWidth, m.height-1)
    if m.activeNote != nil {
        m.viewer.SetContent(m.activeNote.Body, viewerWidth)
    }
}
```

This is the same logic as `WindowSizeMsg` but only changes the split ratio.

### Leftover from the `treeWidth` formula

Currently `treeWidth` is set to `max(msg.Width/4, 25)`. After resize, this formula should no longer override — the user's manual width takes precedence. When the terminal is resized, preserve the user's ratio rather than resetting to 1/4:

```go
case tea.WindowSizeMsg:
    // Preserve user's split ratio on terminal resize
    ratio := float64(m.treeWidth) / float64(m.width)
    m.width = msg.Width
    m.treeWidth = max(int(float64(m.width)*ratio), 15)
    ...
```

This ensures the user's resize persists across terminal resizes.

## Files to modify

| File | Changes |
|------|---------|
| `keys.go` | Add `ShrinkTree`, `GrowTree`, `ResetTree` keys |
| `model.go` | Add `dragSplit`, `treeRatio` fields; `adjustTreeWidth()` method; update `WindowSizeMsg` to preserve ratio |
| `mouse.go` | Add split drag detection in `handleBrowseMouse` |
| `handlers.go` | Add tree resize key handlers in `handleBrowseKey` |
| `statusbar.go` | Show resize hint when hovering near split boundary (mouse) |
| `mouse_test.go`, `model_test.go` | Add tests for drag and keyboard resize |

## Steps

### 1. Mouse drag detection
Add `dragSplit bool` and `dragOffset int` to Model. In `handleBrowseMouse`, check if `msg.Button` is press near the split. Track motion events while `dragSplit` is true. Update `treeWidth` on each motion event.

### 2. Keyboard resize keys
Add `ShrinkTree` (`Ctrl+Left`), `GrowTree` (`Ctrl+Right`), `ResetTree` (`Ctrl+0`) to `KeyMap` and `DefaultKeys()`. Handle in `handleBrowseKey`.

### 3. Layout recalculation
Create `adjustTreeWidth(newWidth)` that re-sizes both panels and re-renders the active note. Call from both mouse and keyboard paths.

### 4. Window resize preserves ratio
Update `tea.WindowSizeMsg` handler to compute `ratio = treeWidth / width` before updating width, then restore `treeWidth = int(ratio * newWidth)`. Clamp to [15, width/2].

### 5. Visual feedback
In `View()`, apply `Accent` foreground to the separator column when `dragSplit` is true (blinking or bright to show active drag).

### 6. Tests
- Mouse: simulate click at split boundary, verify `dragSplit` is set
- Mouse: simulate drag motion, verify `treeWidth` changes
- Keyboard: press `Ctrl+Left`, verify tree shrinks
- Keyboard: press `Ctrl+0`, verify tree resets to default
- Window resize: resize terminal, verify ratio is preserved (not reset to 1/4)

## Completion Criteria

- [ ] Mouse drag on split boundary resizes tree
- [ ] `Ctrl+Left`/`Ctrl+Right` shrinks/grows tree by 5 cols
- [ ] `Ctrl+0` resets tree to default width
- [ ] Tree width clamps between 15 and `width/2`
- [ ] Terminal resize preserves user's split ratio
- [ ] Visual feedback shows active drag
- [ ] All existing tests pass
- [ ] `make vet` exits 0
