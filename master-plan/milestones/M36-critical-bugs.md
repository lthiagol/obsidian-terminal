# M36 — Quick Bug Fixes

**Status:** ⏳ pending

## Goal

Fix 5 critical bugs that are quick to address: accidental quit in search mode, command palette filter bug, dead SetSize code, mouse handlers missing side effects, and nil vault dereference checks.

## Issues

### C3: Accidental quit in search/command palette (`model.go:331-335`)

The global `q`/`Q` quit check fires before mode dispatch. Pressing `q` while typing a search query or in the command palette quits the app instead of treating `q` as text input.

**Fix:** Only quit on `q`/`Q` in Browse and View modes. In Search, Find, Tags, and ProfilePicker modes, let the mode handler process the key. Keep Ctrl+C as global quit.

### C4: Command palette filter shows all items for any query (`command_palette.go:238-255`)

`commandMatches()` returns `true` immediately for global commands (modes == nil), bypassing the name/description filter. Typing any search in the palette shows all global commands regardless of query.

**Fix:** Remove the early return at line 247. Only show commands whose name or description fuzzy-matches the query, regardless of their mode scope.

### C5: `FileTree.SetSize` discards width/height (`tree.go:195-198`)

`SetSize(width, height)` receives dimensions from the model but ignores them (`_ = width; _ = height`). The tree's `ft.width` field stays at the initial 25 from `NewFileTree`. The tree never reflects the actual panel width set by the model or by the user via split resize (M35).

**Fix:** Update `ft.width = width` and `ft.height = height` in `SetSize`. Recalculate any cached values that depend on width.

### C7: Mouse handlers missing side effects (`mouse.go:131-140, 143-158`)

The mouse handlers for opening notes from the tree and search results are incomplete duplicates that don't call:
- `m.buildOutline()`
- `m.backlinkPanel = NewBacklinkPanel(...)`
- `m.addRecentNote(...)`

**Impact:** Clicking a file in the tree doesn't update the outline, backlinks, or recents.

**Fix:** Replace inline code in mouse.go with calls to `transitionToNote(path)` (see M38).

### C8: Nil vault dereference in 3 locations

`ResolveWikiLink` is called without checking if `m.vault` is nil in 3 places:
1. `handlers.go:248` (embed resolver in `openNote`)
2. `handlers.go:116` (view mode link follow in `handleViewKey`)
3. `command_palette.go:75` (follow link command action)

If vault scanning failed partially, these calls panic.

**Fix:** Add nil guards before each `ResolveWikiLink` call. Skip wiki-link resolution if `m.vault == nil`.

## Files to modify

| File | Fix |
|------|-----|
| `model.go` | C3: restrict quit to browse/view modes |
| `command_palette.go` | C4: remove early return in commandMatches; C8: nil guard on vault |
| `tree.go` | C5: implement SetSize properly |
| `handlers.go` | C7: replace inline note opening with transitionToNote; C8: nil guard on vault |
| `mouse.go` | C7: replace inline note opening with transitionToNote |
| `*_test.go` | Update tests for any API changes |

## Completion Criteria

- [ ] `q` only quits in Browse and View modes
- [ ] Command palette filters correctly for all command types
- [ ] Tree width updates correctly when split is resized (M35 integration)
- [ ] Mouse click on tree/search result updates outline, backlinks, and recents
- [ ] All 3 ResolveWikiLink calls have nil vault guards
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
