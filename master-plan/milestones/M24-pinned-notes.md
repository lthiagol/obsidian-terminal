# M24 — Pinned Notes

**Status:** ⏳ pending

## Goal

Pin notes to a working set. Cycle with `Ctrl+]`/`Ctrl+[`, remember scroll position, show indicator in status bar. Better fit for TUI than full tabs.

## Keybindings

**Keys:**
- `p` — pin/unpin current note (Browse and View modes)
- `Ctrl+]` — cycle to next pinned note
- `Ctrl+[` — cycle to previous pinned note

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Implementation Plan

### 1. New type (`model.go`)

```go
type PinnedNote struct { Path string; ScrollY int }
```

### 2. Model fields (`model.go`)

`pinnedNotes []PinnedNote`, `activePinnedIdx int` (-1 when not cycling).

### 3. New methods (`model.go`)

- `togglePin(path string)` — add/remove from pinnedNotes, show toast
- `openPinnedNote(index int)` — loads note, restores ScrollY, sets activePinnedIdx. On load error: removes invalid pin
- `cyclePinnedNext()` — wraps to 0 at end
- `cyclePinnedPrev()` — wraps to last at start

### 4. Viewer changes (`viewer.go`)

Add `GetScrollPosition() int` and `SetScrollPosition(y int)` (delegate to viewport).

### 5. Keybindings (`keys.go`)

Add `PinRune rune` (`p`), `CyclePinPrev []tea.KeyType` (`Ctrl+[`), `CyclePinNext []tea.KeyType` (`Ctrl+]`).

### 6. Handlers (`handlers.go`)

In `handleBrowseKey`: `p` calls `togglePin(selectedPath)`, `Ctrl+[`/`Ctrl+]` cycle.

In `handleViewKey`: Esc saves scroll if pinned, `p` toggles, cycle saves scroll first.

### 7. Statusbar

When viewing a pinned note: append `📌` to info line.

### 8. Rescan

After rescan: validate pinned paths still exist, remove invalid. Reset activePinnedIdx if needed.

### Edge cases

- Pin directory → guard: only pin `.md` files
- Pin same note twice → togglePin is idempotent (2nd call unpins)
- Empty pins + cycle → toast "No pinned notes"
- Scroll beyond content → SetScrollPosition clamps in viewport
- Resize after pin → scroll position may be invalid, clamp on restore

### Implementation order

1. Add PinnedNote + fields to Model
2. Add GetScrollPosition/SetScrollPosition to viewer
3. Add keys to KeyMap
4. Implement togglePin/openPinnedNote/cycle methods
5. Wire p + Ctrl+[/Ctrl+] in handleBrowseKey + handleViewKey
6. Add 📌 to statusbar
7. Validate pins in rescanVault
8. Add help section
9. Write tests

## Completion Criteria

- [ ] PinnedNote type with Path and ScrollY fields
- [ ] `p` keybinding pins/unpins current note
- [ ] `Ctrl+]` cycles to next pinned note
- [ ] `Ctrl+[` cycles to previous pinned note
- [ ] Scroll position saved when switching pinned notes
- [ ] Scroll position restored when opening pinned note
- [ ] 📌 indicator shown in status bar for pinned notes
- [ ] Only `.md` files can be pinned (not directories)
- [ ] Invalid pins removed on rescan
- [ ] Help text updated
- [ ] KEYBINDINGS.md updated
- [ ] `make test` passes
- [ ] `make vet` exits 0
- [ ] Manual test: pin cycling and scroll memory work
