# M30 — Table Rendering Fix

**Status:** ⏳ pending

## Goal

Tables render readably at any terminal width — cells word-wrap instead of truncating, columns distribute space fairly, and narrow terminals get a borderless fallback.

## Problem

The current table renderer has several bugs that make tables unusable on typical terminals (80-120 cols):

1. **Cells truncate instead of wrapping.** A cell with `"This is a long description"` in a 10-char column becomes `"This is a "` — the user sees truncated gibberish. Cells must word-wrap and expand rows vertically.

2. **Proportional scaling loses space to integer rounding.** When `total > width`, columns are scaled by `float64(width)/float64(total)` and truncated to int. The cumulative rounding error means tables often exceed the available width even after "scaling." For example: 3 columns scaled from [15, 12, 10] at width=30 → [9, 7, 6] = 22 + 12(borders) = 34 > 30. The leftover 4 chars overflow.

3. **Border overhead is fixed at 3 chars per column** (` | ` padding + `│` border). At narrow widths this dominates: a 4-column table needs minimum 16 chars just for borders/separators, leaving barely any room for content.

4. **Box-drawing borders compound the ANSI inflation problem.** Each `│├┤┼` styled by lipgloss adds ~10 bytes of ANSI codes per char. A single border row can be 200+ runes with ANSI for a simple table — triggering the old `softWrap` bug even after the recent fix.

5. **No graceful degradation.** When terminal width < minimum table width, we get a corrupted mess instead of a readable alternative.

## Design

### Width allocation algorithm

Replace proportional scaling with **largest remainder method** (same as parliamentary seat allocation):

```go
// 1. Compute each column's desired width (content width, min 3)
// 2. Total available = terminalWidth - (colCount * 3) - 1  // borders+padding
// 3. Give each column floor(desired / totalDesired * totalAvailable)
// 4. Distribute remaining pixels to columns with largest fraction remainders
```

This ensures no space is wasted and columns sum exactly to the available width.

### Cell word-wrapping

Replace `padCell` truncation with wrap-then-pad:

```go
func wrapCell(content string, width int) []string {
    // Word-wrap content to fit width
    // Returns multiple lines if content overflows
    // Each line is padded to exact width
}
```

Each row in the rendered table may span multiple terminal lines (if any cell wraps). The border characters are drawn only once per logical row, with blank continuation lines using `│` + spaces + `│`.

### Borderless fallback

When `terminalWidth < 40`, drop box-drawing borders entirely and render as aligned columns:

```
Name    Type    Status
foo     bar     active
baz     qux     pending
```

Uses 2-space column separation instead of box-drawing chars. This recovers ~3*n+1 chars of space.

### Minimum width check

If a table can't fit even in borderless mode (e.g., 6 columns at 80 chars with min 5 chars each = 30 + 12 separators = 42), render a single-column list view:

```
┌─ Table: Name, Type, Status ───┐
│ Name │ foo                    │
│ Type │ bar                    │
│ ...                           │
└────────────────────────────────┘
```

## Files to modify

| File | Changes |
|------|---------|
| `internal/markdown/markdown.go` | Rewrite `renderTableBlock`, add `wrapCell`, `allocateColumnWidths`, `renderBorderlessTable`, `renderSingleColumnTable`. Update `buildTableRow` to handle multi-line cells. |

## Steps

### 1. Column width allocator
Replace lines 958-989 with `allocateColumnWidths(desired []int, available int) []int` using largest remainder. Write unit tests for edge cases (equal columns, one wide column, all min-width, single column).

### 2. Cell word-wrapping
Create `wrapCell(content string, width int) []string`. Handle: empty cells, single-word cells that exceed width (hard-break), multi-word cells (word-wrap), trimming trailing whitespace. Unit tests.

### 3. Multi-line row rendering
Update `buildTableRow` to accept a row index and render continuation lines. Each logical row produces `max(lines in any cell)` terminal rows. Continuation lines show `│` + padded content + `│` without borders. The top/middle/bottom borders stay single-line.

### 4. Borderless fallback
Add `renderBorderlessTable(lines, width, style)` for width < 40. Renders as aligned columns with 2-space gaps, no box-drawing. Still supports multi-line cells.

### 5. Single-column fallback
Add `renderSingleColumnTable(lines, width, style)` for tables that can't fit even borderless. Renders as a property-list: header → value pairs.

### 6. Integration
Wire the dispatcher in `renderTableBlock`: if borderless mode fits, use it; if not, use single-column. Ensure ANSI styles are applied consistently across all rendering modes.

### 7. Visual regression tests
Add gold-string tests that render known tables at various widths (80, 60, 40, 30) and verify output matches expected. This prevents regressions in column math.

## Completion Criteria

- [ ] Tables render correctly at 80, 60, and 40 column widths
- [ ] Cells wrap instead of truncate
- [ ] Wide tables gracefully degrade to borderless mode at < 40 cols
- [ ] Very wide tables fall back to single-column list view
- [ ] No box-drawing characters are corrupted by ANSI styling
- [ ] Unit tests for column allocation, cell wrapping, and all rendering modes
- [ ] Gold-string regression tests at multiple widths
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
