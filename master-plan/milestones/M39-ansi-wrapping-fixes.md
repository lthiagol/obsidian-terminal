# M39 — ANSI Wrapping & Scroll Estimation Fixes

**Status:** ⏳ pending

## Goal

Fix ANSI style bleed when soft-wrapping styled lines. Fix the scroll estimation drift that causes outline navigation to jump to wrong positions.

## Issues

### C7: ANSI style bleed on soft-wrap boundaries (`viewport.go:55-102`)

When `softWrap` splits a styled line at a wrap boundary, it does not close active ANSI style sequences on the first line or re-open them on the continuation line. A bold span starting on line 1 of a paragraph and continuing to line 2 loses its styling on line 2.

**Fix:** At each wrap boundary, scan backwards through the current line to collect any open SGR sequences (e.g., `\033[1m` without a matching `\033[0m`). Close them before the line break. Then re-open them at the start of the continuation line.

```go
// After wrapping at position `split`:
activeStyles := extractOpenStyles(runes[:split]) // find unclosed SGR starts
lines = append(lines, currentLine + closeStyles(activeStyles))
current.Reset()
current.WriteString(reopenStyles(activeStyles)) // re-open on next line
```

### L8: Scroll estimation drifts for multi-width characters (`model.go:717-748`)

`estimateYOffset` approximates rendered heading positions by dividing raw text length by viewport width. This uses `len(text)` which counts bytes/runes, not display width. For text containing multi-width characters (CJK characters, emoji, box-drawing characters), the estimate is inaccurate because these characters occupy 2 terminal columns but count as 1 rune.

**Note:** The original diagnosis mentioned ANSI escape sequences, but `RenderSegmentsPlain` does not produce ANSI codes. The actual issue is multi-width character handling.

**Fix:** Use a proper display width calculation that accounts for multi-width characters. For paragraphs, use the actual rendered line count from the viewport rather than estimating from raw text length.

### H5: Viewport is exposed through MarkdownViewer

`model.go:671` reads `m.viewer.viewport.Width` directly, reaching through the viewer into its internal viewport. This leak couples outline rendering to viewport internals.

**Fix:** Add `(v MarkdownViewer) Width() int` method and use it instead of reaching into `v.viewport`.

## Files to modify

| File | Changes |
|------|---------|
| `viewport.go` | C7: add `extractOpenStyles()`, `closeStyles()`, `reopenStyles()`; use in softWrap |
| `model.go` | L8: fix `estimateYOffset` to use display width; H5: use `m.viewer.Width()` |
| `viewer.go` | H5: add `Width()` method |
| `viewport_test.go` | C7: test style preservation across wrap boundaries |
| `outline_test.go` | L8: test scroll estimation accuracy with multi-width content |

## Completion Criteria

- [ ] Styled text that spans a wrap boundary retains its styling on continuation lines
- [ ] No visual style bleed (e.g., bold continuing where it shouldn't)
- [ ] Outline navigation scrolls to the correct heading position, within ±2 lines
- [ ] `estimateYOffset` uses display width (accounting for multi-width chars) for paragraph width estimation
- [ ] Viewport is not accessed directly from outside `viewer.go`
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

1 day
