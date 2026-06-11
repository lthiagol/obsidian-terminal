# M25 — Outline / Table of Contents

**Status:** ⏳ pending

## Goal

Show the headings in the current note as a navigable outline sidebar. Press Enter to jump to that section.

## Steps

### 1. Extract headings from parsed markdown

When a note is loaded, extract all `BlockHeading` lines with their levels and text. Build a `[]OutlineItem{Level, Text, LineNumber}` slice.

### 2. Add outline sidebar

Show on the right side of the viewer (or toggle via `o` key). Each heading indented by level. Selected heading highlighted.

### 3. Add navigation

- `o` — toggle outline panel
- `j`/`k` — move cursor in outline
- `Enter` — jump viewer to that heading's line
- `Esc` — close outline

## Completion Criteria

- [ ] `o` toggles outline panel in view mode
- [ ] Headings shown with level-based indentation
- [ ] Enter jumps to heading in viewer
- [ ] `make test && make vet` pass
