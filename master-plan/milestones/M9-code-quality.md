# M9 — Code Quality & Structure

**Status:** ✅ done

## Goal

Consolidate model.go into focused files, fix receiver bugs, and add godoc comments on all exported symbols.

## Files to modify

- `model.go` → split: `modal.go`, `help.go`, `statusbar.go`, `toast.go`
- `tree.go` — fix `FileTree.Update` value receiver
- All `.go` files — add godoc comments

## Steps

### 1. Split model.go (795 lines → multiple files)

- `model.go` — `Model` struct, `NewModel`, `Init`, `Update`, key dispatch (~150 lines)
- `toast.go` — `ToastType`, `Toast`, `addToast`, `expireToasts`, `renderToast`, `renderToasts` (~100 lines)
- `statusbar.go` — `renderStatusBar`, `modeHints`, `truncatePath` (~60 lines)
- `help.go` — `renderHelp` (~70 lines)
- Keep `handleBrowseKey`, `handleViewKey`, `handleSearchKey`, `handleFindKey`, `handleHelpKey` in `model.go` (or a `handlers.go`)

### 2. Fix Init() value receiver dead code (model.go:138-142)

`Init()` uses value receiver, so `m.lastRootModTime` modification is lost.
Move the stat check to `NewModel()` or use pointer receiver (not possible for `Init()`).
Best fix: remove the stat from `Init()`; `checkVaultChanges` / `rescanVault` already handle it on first tick.

### 3. Fix FileTree.Update value receiver mutation (tree.go:105-135)

`func (ft FileTree) Update(...)` modifies cursor, expanded state on a copy.
Change to pointer receiver: `func (ft *FileTree) Update(...)`.
Requires updating all callers that currently do `ft, _ = ft.Update(msg)`.

### 4. Pointer vs value receiver consistency

- `Model.Update`, `Model.handleBrowseKey`, etc. — all use value receiver despite mutating state.
  Bubble Tea convention is value receiver + return, so this is acceptable idiomatically,
  but document the pattern clearly.

### 5. Add godoc comments on all exported symbols

Every exported type, function, method, constant, and variable needs a doc comment.
Key items:
- `Model`, `NewModel`, `Mode`, `ToastType`, `Toast`, `TickMsg`
- `Config`, `DefaultConfig`, `LoadConfig`
- `VaultEntry`, `VaultNote`, `ScanVault`, `LoadNote`
- `FileTree`, `NewFileTree`
- `MarkdownViewer`, `NewViewer`
- `KeyMap`, `DefaultKeys`, `MatchKey`, `MatchRune`
- `SearchState`, `NewSearchState`, `SearchResult`, `FuzzyScore`, `FuzzySearch`, `ContentSearch`
- `ParseMarkdown`, `ExtractWikiLinks`, `ResolveWikiLink`, `RenderMarkdown`
- All `BlockType`, `InlineSegment`, `MarkdownLine`, `WikiLink`
- All theme color/style constants

## Completion Criteria

- [ ] `model.go` split: `toast.go`, `statusbar.go`, `help.go` exist and `model.go` < 250 lines
- [ ] `Init()` no longer dead-modifies `lastRootModTime`
- [ ] `FileTree.Update` uses pointer receiver; all callers updated
- [ ] All exported symbols have godoc comments
- [ ] `make test` passes all 79 tests
- [ ] `make vet` exits 0
