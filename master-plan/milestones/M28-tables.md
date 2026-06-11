# M28 — Markdown Tables

**Status:** ⏳ pending

## Goal

Parse and render pipe-delimited markdown tables with Unicode box-drawing borders and alignment.

## Implementation Plan

### 1. New types (`internal/markdown/markdown.go`)

```go
type TableAlignment int  // AlignLeft, AlignCenter, AlignRight
```

Add to MarkdownLine: `TableCells []string`, `TableAlign []TableAlignment`.

Add `BlockTable` to BlockType.

### 2. Table detection (`internal/markdown/markdown.go`)

New functions:
- `isTableRow(line) bool` — starts and ends with `|`
- `isTableSeparator(line) bool` — only `|`, `-`, `:`, spaces between pipes
- `parseTableRow(line) []string` — split on `|`, trim whitespace
- `parseTableAlignment(line) []TableAlignment` — detect `:---`, `:---:`, `---:`

In `ParseMarkdown` (convert loop to index-based): when line is a table row AND next line is separator: parse header + separator, then consume subsequent table rows.

### 3. Table rendering (`internal/markdown/markdown.go`)

`renderTableBlock(lines []MarkdownLine, width int, style RendererStyle) string`:
- Calculate column widths from all cells
- Clamp to available width, distribute proportionally
- Render with Unicode box-drawing: `┌─┬─┐` top border, `│ cell │` rows, `├─┼─┤` separator, `└─┴─┘` bottom
- Apply alignment (left/center/right padding)
- Header row bold + accent color, data rows secondary color

In `RenderMarkdown`: collect consecutive BlockTable lines, pass to renderTableBlock as single unit.

### 4. Update renderLine

Add `case BlockTable: ...` (single-row fallback if block rendering fails).

### Edge cases

- Single-column table → minimum padding
- Missing cells in a row → pad with empty string
- Escaped pipe `\|` → handle in v1 (treat as literal pipe in cell content)
- Table wider than viewport → proportional column scaling
- No separator row → not detected as table (treated as paragraph with pipes)
- Empty cells → render with minimum width

### Implementation Notes

**Escaped pipes:** Parse `\|` as literal `|` in cell content:

```go
func parseTableRow(line string) []string {
    // Replace \| with placeholder, split on |, restore placeholder to |
    escaped := strings.ReplaceAll(line, `\|`, "\x00")
    cells := strings.Split(escaped, "|")
    for i, cell := range cells {
        cells[i] = strings.ReplaceAll(strings.TrimSpace(cell), "\x00", "|")
    }
    // Remove first and last empty cells from leading/trailing |
    if len(cells) >= 2 && cells[0] == "" {
        cells = cells[1:]
    }
    if len(cells) >= 1 && cells[len(cells)-1] == "" {
        cells = cells[:len(cells)-1]
    }
    return cells
}
```

### Implementation order

1. Add BlockTable + new fields to MarkdownLine
2. Add table detection functions
3. Update ParseMarkdown with index-based loop + table detection
4. Add renderTableBlock
5. Update RenderMarkdown to collect consecutive table rows
6. Write tests

## Completion Criteria

- [ ] Table detection: header row + separator row + data rows
- [ ] Alignment parsing: `:---` (left), `:---:` (center), `---:` (right)
- [ ] Unicode box-drawing borders rendered correctly
- [ ] Column widths calculated and distributed proportionally
- [ ] Header row styled with bold + accent color
- [ ] Data rows styled with secondary color
- [ ] Escaped pipes (`\|`) handled as literal pipes
- [ ] Tables wider than viewport scaled proportionally
- [ ] Missing cells padded with empty strings
- [ ] `make test` passes
- [ ] `make vet` exits 0
- [ ] Manual test: tables render correctly with alignment
