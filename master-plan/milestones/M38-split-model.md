# M38 — Split model.go and Consolidate Note Opening

**Status:** 🟡 partial — see **M52** for remaining work

## Goal

Reduce `model.go` from 868 lines to under 250 by extracting subsystems into dedicated files. Create a single `transitionToNote` method that handles all side effects when opening a note, replacing 6 duplicated code paths.

## Issues

### H1: model.go is 868 lines (3.5x over the 250-line guideline)

The file contains: Model struct, Init/Update/View, mode-change helpers, vault management, pin management, outline management, daily notes, recent notes, command palette, profile picker. Most of these should be standalone files.

### H2: "Open note" logic duplicated 6 times

The same pattern (load note → set activeNote → set mode → set viewer content) appears in 6 locations, but with inconsistent side effects:

| Location | Function | Sets embed resolver? | buildOutline? | addRecentNote? | backlinkPanel? |
|----------|----------|---------------------|---------------|----------------|----------------|
| `handlers.go:236` | `openNote()` | ✅ | ✅ | ✅ | ✅ |
| `model.go:580` | `openPinnedNote()` | ❌ | ❌ | ❌ | ❌ |
| `model.go:756` | `openDailyNote()` | ❌ | ✅ | ✅ | ❌ |
| `model.go:805` | `openRecentNote()` | ❌ | ❌ | ❌ | ❌ |
| `mouse.go:131` | tree click (inline) | ❌ | ❌ | ❌ | ❌ |
| `mouse.go:143` | `openSearchResult()` | ❌ | ❌ | ❌ | ❌ |

**Impact:** Clicking a file in the tree doesn't update the outline, backlinks, or recents. Opening a pinned note doesn't add it to recents. The behavior is inconsistent and buggy.

### H3: View() method is 83 lines with 6 levels of nesting

The View method branches on mode, then sub-mode (outline/backlink/recent/profile picker), building inline styles and recalculating dimensions. The layout logic should be extracted.

## Design

### Part 1: Consolidate note opening

Create one method that opens a note and handles ALL side effects:

```go
func (m *Model) transitionToNote(path string) {
    note, err := LoadNote(m.config.VaultPath, path)
    if err != nil {
        m.addToast("Could not load note: " + err.Error(), ToastError)
        return
    }
    m.activeNote = note
    m.prevMode = m.mode
    m.mode = ModeView
    m.viewer.SetContent(note.Body, m.width - m.treeWidth - 2)
    m.addRecentNote(path)
    m.buildOutline()
    m.backlinkPanel = NewBacklinkPanel(m.backlinkIndex, note.Path)
    if m.vault != nil {
        m.viewer.SetEmbedResolver(func(target, heading string) (string, error) {
            resolved := ResolveWikiLink(target, m.vault, m.config.VaultPath)
            if resolved == "" { return "", nil }
            return LoadSection(m.config.VaultPath, resolved, heading)
        })
    }
}
```

Replace all 6 call sites with `m.transitionToNote(path)`.

### Part 2: Extract subsystems

| File | Extracted from model.go |
|------|------------------------|
| `pin_handler.go` | Pin toggle, cycle, validate |
| `outline_render.go` | buildOutline, renderOutline |
| `daily_handler.go` | buildDailyNotePath, openDailyNote |
| `recent_handler.go` | addRecentNote, toggleRecents, openRecentNote, renderRecents |

Note: `applyProfile` and profile switch handling are already fixed in M37 and can remain in `handlers.go`.

## Files to modify

| File | Changes |
|------|---------|
| `handlers.go` | Add `transitionToNote` method; replace `openNote` body with call to it |
| `model.go` | Replace `openPinnedNote`, `openDailyNote`, `openRecentNote` bodies with calls to `transitionToNote`; extract pin/outline/daily/recent code to new files |
| `mouse.go` | Replace inline note opening and `openSearchResult` with calls to `transitionToNote` |
| `pin_handler.go` | **New** — pin toggle/cycle/validate from model.go |
| `outline_render.go` | **New** — outline rendering from model.go |
| `daily_handler.go` | **New** — daily note logic from model.go |
| `recent_handler.go` | **New** — recent notes logic from model.go |

## Completion Criteria

- [ ] All 6 open-note code paths use the single `transitionToNote` method
- [ ] Clicking a file in the tree updates outline, backlinks, and recents
- [ ] Opening a pinned note adds it to recents
- [ ] All note-opening paths set the embed resolver consistently
- [ ] model.go is under 250 lines
- [ ] Extracted files have clear, single responsibilities
- [ ] All existing tests pass (adjust for new file locations)
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

2-3 days
