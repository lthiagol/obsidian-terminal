# M18 — Mouse Support

**Status:** ⏳ pending

## Goal

Handle `tea.MouseMsg` for click-to-select, double-click-to-open, and scroll wheel in tree, viewer, search, and help panels. Tea already enables `WithMouseCellMotion()`.

## Layout coordinate mapping

```
Terminal: width×height cells
┌──────────────────────────────────────────────┐ row 0
│ Tree [0,treeWidth) │ Right [treeWidth,width) │
│ h=height-1          │ h=height-1              │
├────────────────────┴─────────────────────────┤ row height-1
│ Status Bar (ignored)                         │
└──────────────────────────────────────────────┘
```

## Implementation Plan

### New file: `mouse.go`

```go
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)
func (m Model) handleRightPanelMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)
func (m Model) handleTreeMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)
func (m Model) handleViewerMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)  
func (m Model) handleSearchMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)
func (m Model) handleHelpMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)
func (m Model) openTreeItem() (tea.Model, tea.Cmd)        // extracted from handleBrowseKey Enter
func (m Model) openSearchResult() (tea.Model, tea.Cmd)    // extracted from handleSearchOrFind Enter
func (m *Model) isDoubleClick(x, y int) bool
func (m *Model) recordClick(x, y int)
```

### Tree panel (`handleTreeMouse`)

| Event | Action |
|-------|--------|
| Left click | Move cursor to Y index (direct mapping: `cursor = msg.Y`, clamped) |
| Left double-click | Open item (same as Enter) via `openTreeItem()` |
| Wheel up | MoveUp() × 3 |
| Wheel down | MoveDown() × 3 |

Double-click detection: `time.Since(lastClickTime) <= 500ms` AND X/Y within 1 cell tolerance.

### Viewer panel (`handleViewerMouse`)

| Event | Action |
|-------|--------|
| Wheel up/down | Forward to `m.viewer.Update(msg)` — viewport handles MouseMsg internally |

Viewport already handles wheel events, so just forward. Wiki-link click detection deferred (complex coordinate math for inline links).

### Search panel (`handleSearchMouse`)

| Event | Action |
|-------|--------|
| Left click | Select result at `index = msg.Y - 2` (header=row0, blank=row1, results=row2+) |
| Left double-click | Open result via `openSearchResult()` |
| Wheel up/down | MoveUp/MoveDown × 3 |

Need new `SetSelected(i)` and `SelectedIndex()` methods on `search.State`.

### Help panel (`handleHelpMouse`)

| Event | Action |
|-------|--------|
| Wheel up | `helpScroll = max(0, helpScroll - 3)` |
| Wheel down | `helpScroll += 3` |

### Model changes (`model.go`)

Add fields: `lastClickTime time.Time`, `lastClickX int`, `lastClickY int`

Add `case tea.MouseMsg:` in `Update()` type switch, before KeyMsg handler.

### Refactoring (`handlers.go`)

Extract shared logic:
- `openTreeItem()` from `handleBrowseKey` Enter case
- `openSearchResult()` from `handleSearchOrFind` Enter case
- Update `handleBrowseKey` and `handleSearchOrFind` to call shared helpers

### New search methods (`internal/search/search.go`)

```go
func (s *State) SetSelected(i int) { ... }  // clamp and set
func (s State) SelectedIndex() int { ... }   // return s.selected
```

### Implementation order

1. Add `SetSelected`/`SelectedIndex` to `search.State`
2. Create `mouse.go` with all handlers + shared `openTreeItem`/`openSearchResult`
3. Add `lastClickTime/X/Y` + `case tea.MouseMsg:` to model.go
4. Refactor `handlers.go` Enter cases to use shared helpers
5. Write tests in `mouse_test.go`
6. Run `make test && make vet`
