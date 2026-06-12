# M36 — Fix Critical Bugs

**Status:** ⏳ pending

## Goal

Fix 6 critical bugs found during code review: theme data races, nil vault dereference, accidental quit in search mode, command palette filter bug, ANSI style bleed on wrap, and YAML parser fragility.

## Issues

### C1: Data race on global theme variables (`theme.go:413-436`)

`activatePalette()` writes to 15+ global lipgloss style variables (`Accent`, `TreeStyle`, `ViewerStyle`, etc.) at runtime when switching profiles. Meanwhile `View()`, `renderStatusBar()`, `renderHelp()`, `FileTree.View()` read these globals. Bubble Tea's Update/View cycle is single-threaded, but profile switching triggers `rescanVault()` which may race.

**Fix:** Convert global style variables to fields on the `Model` struct (or a `ThemeState` sub-struct). All rendering code reads from the model instead of globals.

### C2: Nil vault dereference in embed resolver (`handlers.go:248`)

`openNote()` sets an embed resolver that calls `ResolveWikiLink(target, m.vault, ...)` without checking if `m.vault` is nil. If vault scanning failed partially, this panics.

**Fix:** Add a nil guard before constructing the embed resolver. Skip setting the resolver if `m.vault == nil`.

### C3: Accidental quit in search/command palette (`model.go:331-335`)

The global `q`/`Q` quit check fires before mode dispatch. Pressing `q` while typing a search query or in the command palette quits the app instead of treating `q` as text input.

**Fix:** Only quit on `q`/`Q` in Browse and View modes. In Search, Find, Tags, and ProfilePicker modes, let the mode handler process the key. Keep Ctrl+C as global quit.

### C4: Command palette filter shows all items for any query (`command_palette.go:238-255`)

`commandMatches()` returns `true` immediately for global commands (modes == nil), bypassing the name/description filter. Typing any search in the palette shows all global commands regardless of query.

**Fix:** Remove the early return. Only show commands whose name or description fuzzy-matches the query, regardless of their mode scope.

### C5: `FileTree.SetSize` discards width/height (`tree.go:195-198`)

`SetSize(width, height)` receives dimensions from the model but ignores them (`_ = width; _ = height`). The tree's `ft.width` field stays at the initial 25 from `NewFileTree`. The tree never reflects the actual panel width set by the model or by the user via split resize (M35).

**Fix:** Update `ft.width = width` and `ft.height = height` in `SetSize`. Recalculate any cached values that depend on width.

### C6: `applyProfile` mutates a copy (`handlers.go:433-470`)

`applyProfile` has a value receiver `(m Model)` but mutates `m.config.VaultPath`, `m.config.Theme`, `m.palette`, etc. Since it returns `tea.Msg`, the mutations happen on a copy and are discarded. Profile switching via the profile picker is silently broken.

**Fix:** Change `applyProfile` to a pointer receiver `(m *Model)` and return `nil` command. Or return a `tea.Batch` that triggers rescan. Ensure the palette and styles update on the real model, not a copy.

## Files to modify

| File | Fix |
|------|-----|
| `theme.go` | C1: remove global mutable state, move style fields to model |
| `model.go` | C1: add theme/style fields; C3: restrict quit to browse/view; C6: pointer receiver |
| `handlers.go` | C2: nil guard on vault; C6: pointer receiver on applyProfile |
| `command_palette.go` | C4: remove early return in commandMatches |
| `tree.go` | C5: implement SetSize properly |
| `*_test.go` | Update tests for any API changes |

## Completion Criteria

- [ ] No global mutable style variables — all rendering reads from model fields
- [ ] Profile theme switching works and reflects immediately
- [ ] `q` only quits in Browse and View modes
- [ ] Command palette filters correctly for all command types
- [ ] Tree width updates correctly when split is resized (M35 integration)
- [ ] Profile switching via picker actually changes vault/theme
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
