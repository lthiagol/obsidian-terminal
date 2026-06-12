# M38 — Split model.go and Eliminate Duplication

**Status:** ⏳ pending

## Goal

Reduce `model.go` from 868 lines to under 250 by extracting subsystems into dedicated files. Eliminate the 5 duplicated "open note" patterns by creating a single `transitionToNote` method.

## Issues

### H1: model.go is 868 lines (3.5x over the 250-line guideline)

The file contains: Model struct, Init/Update/View, mode-change helpers, vault management, pin management, outline management, daily notes, recent notes, command palette, profile picker. Most of these should be standalone files.

### H2: "Open note" logic duplicated 5 times

The same pattern (load note → set activeNote → set mode → set viewer content → build backlinks → build outline → add recent) appears in `handlers.go:236`, `model.go:597`, `mouse.go:123`, `mouse.go:143`, `model.go:767`. Each variant has different subsets of side effects.

### H3: View() method is 83 lines with 6 levels of nesting

The View method branches on mode, then sub-mode (outline/backlink/recent/profile picker), building inline styles and recalculating dimensions. The layout logic should be extracted.

### M10: applyProfile broken (value receiver)

`applyProfile` mutates config/palette/style fields on a value-receiver copy. They're discarded. Fix by using pointer receiver.

## Design

### New file structure

| File | Extracted from model.go |
|------|------------------------|
| `pin_handler.go` | Pin toggle, cycle, validate |
| `outline_render.go` | buildOutline, renderOutline |
| `daily_handler.go` | buildDailyNotePath, openDailyNote |
| `recent_handler.go` | addRecentNote, toggleRecents, openRecentNote, renderRecents |
| `profile_handler.go` | applyProfile (fixed), profile switch handling |

### Single note transition

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

Replace all 5 call sites with `m.transitionToNote(path)`.

## Files to modify

| File | Changes |
|------|---------|
| `model.go` | Remove extracted code, keep only Model struct + Init/Update/View + search/index fields |
| `pin_handler.go` | **New** — pin toggle/cycle/validate from model.go |
| `outline_render.go` | **New** — outline rendering from model.go |
| `daily_handler.go` | **New** — daily note logic from model.go |
| `recent_handler.go` | **New** — recent notes logic from model.go |
| `profile_handler.go` | **New** — applyProfile + profile switch from handlers.go |
| `handlers.go` | Replace openNote with transitionToNote call; remove applyProfile |
| `mouse.go` | Replace openTreeItem/openSearchResult with transitionToNote call |

## Completion Criteria

- [ ] model.go is under 250 lines
- [ ] All 5 open-note code paths use the single `transitionToNote` method
- [ ] applyProfile uses pointer receiver and works correctly
- [ ] Extracted files have clear, single responsibilities
- [ ] All existing tests pass (adjust for new file locations)
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
