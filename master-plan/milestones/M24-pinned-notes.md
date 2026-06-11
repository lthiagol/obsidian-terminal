# M24 — Pinned Notes

**Status:** ⏳ pending

## Goal

Pin notes to a working set, cycle through them with `Ctrl+]`/`Ctrl+[`, and remember scroll position per pin. Better fit for TUI than full tab support — no wasted screen rows.

## Steps

### 1. Add pinned notes to Model

```go
type PinnedNote struct {
    Path    string
    Title   string
    Body    string
    YOffset int // remembered scroll position
}

// Model fields
pinnedNotes []PinnedNote
pinnedIdx   int
```

### 2. Add keybindings

- `p` — pin/unpin current note (toggle)
- `Ctrl+]` — next pinned note
- `Ctrl+[` — previous pinned note
- Status bar shows `📌 note-name (2/3)` when cycling, `📌 3` when not cycling

### 3. Pin management

- Max 10 pinned notes
- Pinning a note already pinned → unpin it
- When a pinned note is deleted (rescan detects), remove it from pins gracefully
- Each pin remembers its scroll Y-offset when unpinning/cycling away

### 4. UI

- Status bar shows pin count: `📌 3`
- Cycling shows `📌 {title} (2/3)` in the status bar info section
- No extra panel or sidebar — minimal footprint

## Completion Criteria

- [ ] `p` toggles pin on current note
- [ ] `Ctrl+]`/`Ctrl+[` cycle through pinned notes
- [ ] Scroll position remembered per pin
- [ ] Deleted notes removed from pins on rescan
- [ ] Status bar shows pin count
- [ ] `make test && make vet` pass
