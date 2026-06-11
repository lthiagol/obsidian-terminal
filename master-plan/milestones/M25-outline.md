# M25 — Outline / Table of Contents

**Status:** ⏳ pending

## Goal

Show headings from current note as navigable outline overlay. Toggle with `t`, Enter to jump.

## Implementation Plan

### 1. Markdown changes (`internal/markdown/markdown.go`)

```go
type HeadingInfo struct { Level int; Text string; LineIdx int }
func ExtractHeadings(lines []MarkdownLine) []HeadingInfo
```

### 2. Model fields (`model.go`)

`outlineVisible bool`, `outlineItems []OutlineItem`, `outlineCursor int`

```go
type OutlineItem struct {
    Level   int
    Text    string
    LineIdx int
    YOffset int  // approximate Y offset in rendered viewport
}
```

### 3. New methods (`model.go`)

- `buildOutline()` — parse headings from activeNote, estimate Y offsets
- `renderOutline()` — styled list with level indentation, cursor highlight

### 4. Handler (`handlers.go`)

`handleOutlineKey()` — Esc/`t` dismiss, j/k navigate, Enter jumps to YOffset and dismisses.

In `handleViewKey`: `t` toggles outline, calls `buildOutline()`.

In `Update()`: check `outlineVisible` before mode dispatch, route to `handleOutlineKey`.

### 5. Note-load sites

Call `buildOutline()` at all places where note content is loaded:
- handleBrowseKey Enter
- handleSearchOrFind Enter
- handleViewKey link follow
- rescanVault note reload

### 6. View()

When `outlineVisible` true, show outline instead of viewer content in right panel.

### Edge cases

- Note with no headings → "No headings in this note"
- YOffset approximation → accept "good enough" for v1 (lines vary with code blocks/wrapping)
- Outline + window resize → outline stays visible, SetSize doesn't affect it
- `t` in other modes → only works in ModeView

### Implementation order

1. Add HeadingInfo + ExtractHeadings to markdown.go
2. Add outline fields to Model
3. Implement buildOutline + renderOutline
4. Add handleOutlineKey
5. Wire `t` in handleViewKey + dispatch in Update
6. Add buildOutline at note-load sites
7. Write tests
