# M28 — Markdown Tables

**Status:** ⏳ pending

## Goal

Parse and render markdown tables with box-drawing borders in the viewer.

## Steps

### 1. Parse table syntax

Add `BlockTable` type. Detect table rows:
```
| Header 1 | Header 2 |
|----------|----------|
| cell 1   | cell 2   |
```

Parse alignment from separator row (`:---`, `:---:`, `---:`).

### 2. Render with box-drawing characters

Use Unicode box-drawing characters (`│`, `─`, `┼`, etc.) for borders. Auto-size columns based on content width. Handle cells that exceed the viewport width (truncate with `…`).

### 3. Handle edge cases

- Single-column tables
- Missing cells (fewer cells in a row than header)
- Escaped pipe characters in cells (`\|`)
- Empty cells

## Completion Criteria

- [ ] Tables parsed and rendered with box-drawing borders
- [ ] Column widths auto-sized
- [ ] Alignment respected (left/center/right)
- [ ] Wide cells truncated gracefully
- [ ] `make test && make vet` pass
