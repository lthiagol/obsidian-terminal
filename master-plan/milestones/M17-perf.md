# M17 — Performance: Hot Path Optimizations

**Status:** ⏳ pending

## Goal

Reduce per-frame and per-keystroke allocations in tree rendering, fuzzy search, and help text.

## Implementation Plan

### 1. Pre-computed tree styles + prefix cache (`tree.go`)

**Problem:** `View()` allocates 3-5 lipgloss.Style objects per item per frame (300k allocs/sec at 60fps with 1000 files).

**Fix:** Add 3 cached styles + prefix cache to `FileTree` struct:
```go
type FileTree struct {
    items, cursor, vault, width, height  // existing
    fileStyle, dirStyle, selectedStyle lipgloss.Style  // NEW
    prefixCache []string  // NEW: index=depth, value=strings.Repeat("  ", depth)
}
```

Pre-compute in `NewFileTree`: compute max depth, build prefixCache, create 3 styles once.

Rewrite `View()` body: use `ft.prefixCache[item.depth]`, use pre-computed styles instead of `lipgloss.NewStyle()` per item.

### 2. Pre-computed lowercase paths in search (`internal/search/search.go`)

**Problem:** `FuzzyScore` calls `strings.ToLower(target)` once per file per keystroke (10k calls per typed character with 10k files).

**Fix:** Add `allPathsLower []string` to `State` struct. Pre-compute in `NewState`. Thread through as extra parameter to `FuzzySearch` and `FuzzyScore`:
- `FuzzySearch(query, paths, pathsLower []string)` 
- `FuzzyScore(query, target, targetLower string)`
- Update all internal callers and test call sites

**Signature changes** (only affect search package + tests):
- `FuzzyScore(query, target, targetLower string) float64`
- `FuzzySearch(query string, paths, pathsLower []string) []Result`

**Test impact:** All `FuzzyScore`/`FuzzySearch` test calls need 3rd argument. Add `lowerPaths()` helper to test file. No external callers affected.

### 3. Cache help text (`help.go`)

**Problem:** `renderHelp()` rebuilds entire static help text every frame.

**Fix:** Extract into `buildHelpLines()`, store in package-level `var cachedHelpLines []string`. Compute once on first render. Add `InvalidateHelpCache()` called from `activatePalette` in `theme.go` when palette changes.

### Files changed

| File | Changes |
|------|---------|
| `tree.go` | 3 cached styles + prefixCache in FileTree; pre-compute in NewFileTree; rewrite View() body |
| `internal/search/search.go` | allPathsLower in State; pre-compute in NewState; extra params to FuzzySearch/FuzzyScore |
| `internal/search/search_test.go` | lowerPaths() helper; update all call sites |
| `help.go` | cachedHelpLines var + buildHelpLines() + InvalidateHelpCache(); simplify renderHelp() |
| `theme.go` | Call InvalidateHelpCache() in activatePalette() |

### Optimization impact summary

| Optimization | Saved per frame | At scale (1000 files) |
|-------------|----------------|----------------------|
| Tree styles | 3-5 allocs/item | ~4,000 allocs/frame |
| Tree prefixes | 1 alloc/item | ~1,000 allocs/frame |
| Lowercase paths | 1 alloc/file/keystroke | ~10,000 allocs/keystroke |
| Help cache | 1 alloc/frame | 60 allocs/sec |

### Implementation order
1. Tree: add fields, pre-compute in NewFileTree, rewrite View()
2. Search: add allPathsLower, update signatures, update callers/tests
3. Help: add cachedHelpLines, update renderHelp(), wire InvalidateHelpCache()
4. Run `make test && make vet`
