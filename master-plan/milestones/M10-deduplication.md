# M10 — Deduplication & DRY

**Status:** ✅ done

## Goal

Eliminate all copy-pasted code patterns across the codebase.

## Files to modify

- `model.go` — merge search handlers, common mode switches, render helpers
- `vault.go` + `markdown.go` — unified frontmatter parsing
- `tree.go` — remove duplicate `keys()`/`keysRunes()`
- `search.go` — extract shared style constant

## Steps

### 1. Merge `handleSearchKey` and `handleFindKey` (model.go:330-398)

These are identical apart from which `SearchMode` was set on entry.
Extract into a single method — the mode is already stored in `m.searchState.mode`.

### 2. Merge `renderSearch` and `renderFind` (model.go:489-505)

Differ only in label (`"fuzzy"` / `"content"`) and result caption.
Parameterize with `(label string, resultLabel string)`.

### 3. Extract common mode-switching dispatch (model.go:214-228, 275-289)

`handleBrowseKey` and `handleViewKey` share identical 3-branch dispatch for
Search, Find, Help. Extract into `m.switchToMode(target Mode)` helper.

### 4. Unified frontmatter parsing (vault.go:215-257 + markdown.go:461-474)

Three functions (`parseFrontmatter`, `stripFrontmatter`, `stripMarkdownFrontmatter`)
all contain the same `---\n` detection and closing-marker search.
Extract shared `findFrontmatterBounds(content string) (start, end int, ok bool)`.
Use it in all three callers.

### 5. Remove duplicate `keys()` / `keysRunes()` (tree.go:137-143)

Both return `DefaultKeys()`. Remove `keysRunes()`.

### 6. Extract selected-item style to package-level var (search.go:296-298, 315-316)

Same style appears in `renderFileList` and `formatSearchResult`.
Define `var selectedItemStyle = lipgloss.NewStyle().Background(Accent).Foreground(...)` once.

## Completion Criteria

- [ ] `handleSearchKey` + `handleFindKey` → single method
- [ ] `renderSearch` + `renderFind` → single parameterized function
- [ ] Common mode-switching helper extracted
- [ ] Single `findFrontmatterBounds` used by all 3 callers
- [ ] No duplicate `keys()` methods
- [ ] No duplicate selected-item style
- [ ] `make test` passes all 79 tests
- [ ] `make vet` exits 0
