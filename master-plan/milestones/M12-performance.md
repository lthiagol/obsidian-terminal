# M12 — Performance

**Status:** ⏳ pending

## Goal

Eliminate redundant regex compilation, reduce allocations, and fix O(n²) markdown scanning.

## Files to modify

- `markdown.go` — move regexes to package-level, fix `findNextSpecial`
- `search.go` — optimize `ContentSearch` allocations

## Steps

### 1. Move regexes to package-level `var` declarations (markdown.go)

Currently these regexes are compiled on every call:
- `visibleLen` — per word in `wrapText`
- `isListItem` — 2x per line in `parseMarkdownLines`
- `extractCalloutType` — per line in callout detection
- `stripComments` — per `ParseMarkdown` call
- `stripBlockquote` — per blockquote line
- `parseListItem` — per list item line

Move all to package-level `var` using `regexp.MustCompile`.

### 2. Fix `findNextSpecial` O(n²) scanning (markdown.go:420-430)

Each call does 10 `strings.Index` scans on increasingly smaller suffixes.
For a long plain-text paragraph of length n, this becomes O(n²).
Replace with a single-pass scanner using `regexp.FindStringIndex` on
a combined regex matching all inline formatting patterns at once:
`` `|**|__|*|_|~~|[[|![[|#``

### 3. Optimize `ContentSearch` double case-lowering (search.go:226-233)

`strings.ToLower(content)` on full body + `strings.ToLower(line)` per matching line
doubles the work. Pre-compute `lowerContent := strings.ToLower(content)` once,
split on `\n`, and reuse the lowered slice.

### 4. Reduce `ContentSearch` per-line allocation (search.go:228)

`strings.Split(content, "\n")` allocates a full slice per file.
Use `strings.Index` or `bufio.Scanner` to iterate lines without full split.
Or use `strings.Cut` in a loop to avoid the allocation.

### 5. Optimize tree expand/collapse (tree.go:73-74)

`append(ft.items[:pos], append(childItems, ft.items[pos:]...)...)` double-appends.
Use `slices.Insert` (Go 1.21+) or pre-allocate a single backing array.

## Completion Criteria

- [ ] All 6 regexes are package-level `var` (no per-call compilation)
- [ ] `findNextSpecial` uses single-pass regex or scanner (no O(n²))
- [ ] `ContentSearch` does single `ToLower` for body + per-line
- [ ] `ContentSearch` avoids full `strings.Split` allocation
- [ ] Tree expand/collapse uses single allocation
- [ ] `make test -race` passes
- [ ] `make vet` exits 0
