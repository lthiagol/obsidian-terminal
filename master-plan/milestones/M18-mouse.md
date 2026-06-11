# M18 — Mouse Support

**Status:** ⏳ pending

## Goal

Add mouse support for the file tree, viewer, search results, and help panel. Tea already enables `WithMouseCellMotion()` — we just need to handle `tea.MouseMsg`.

## Files to modify

- `model.go` — add `tea.MouseMsg` dispatch
- `tree.go` — click-to-select, double-click-to-open, scroll wheel
- `viewer.go` — scroll wheel, click wiki-links
- `search.go` (internal) — no view changes; search results rendered by model
- `handlers.go` — mouse handler stubs per mode

## Steps

### 1. Add MouseMsg dispatch in model.go

In `Update()`, add a `tea.MouseMsg` case that dispatches to mode-specific mouse handlers:

```go
case tea.MouseMsg:
    return m.handleMouse(msg)
```

`handleMouse` routes to `handleMouseBrowse`, `handleMouseView`, `handleMouseSearch`, `handleMouseFind`, `handleMouseHelp`.

### 2. Tree click support

Map the mouse Y-coordinate to a tree item index. The tree starts at (0, 0) on the left panel. Each item is one line. Click → move cursor to that item. Double-click → toggle folder (if dir) or open note (if file). Scroll wheel → move cursor up/down.

### 3. Viewer scroll + link clicking

- Scroll wheel → `viewport.ScrollUp/Down(3)` (3 lines per tick feels natural)
- Click on a wiki-link line → cycle to that link or follow it. The viewer knows which lines contain links via the `links` slice positions.

### 4. Search result clicking

Map Y-coordinate to search result index. Click → select result. Double-click → open note.

### 5. Help panel scrolling

Scroll wheel → scroll help text.

## Completion Criteria

- [ ] Click tree items to select, double-click to toggle/open
- [ ] Scroll wheel works in tree, viewer, search, and help
- [ ] Click wiki-links in viewer to cycle/follow
- [ ] Click search results to open notes
- [ ] All 98 tests pass
- [ ] `make build && make vet` exit 0
