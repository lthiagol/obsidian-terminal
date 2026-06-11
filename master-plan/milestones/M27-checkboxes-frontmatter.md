# M27 — Checkboxes + Frontmatter Display

**Status:** ⏳ pending

## Goal

Render `- [ ]` and `- [x]` with styled checkbox icons. Display parsed frontmatter metadata above note content in viewer.

## Implementation Plan

### 1. Checkbox parsing (`internal/markdown/markdown.go`)

Add `Checked bool` and `Checkable bool` to `MarkdownLine`.

In `ParseMarkdown` list item branch (~line 160): after `parseListItem`, check if text starts with `[ ]` or `[x]`/`[X]`. If so: set `Checkable=true`, `Checked=(text[1]=='x'||text[1]=='X')`, strip `[x] ` prefix from text.

### 2. Checkbox rendering (`internal/markdown/markdown.go`)

In `renderList`: if `Checkable`, render `[x]` or `[ ]` with success/dimmed color instead of bullet.

### 3. Frontmatter rendering (`viewer.go`)

New function `renderFrontmatter(rawMarkdown, width, style)`:
- Detect `---\n...\n---\n` block at start of content
- Parse key:value pairs from YAML block
- Render as formatted metadata block with border: `─── Frontmatter` header, key: value pairs, `───` footer

Update `SetContent` to prepend frontmatter block before rendered markdown.

### Edge cases

- `- [x]` followed by bold text (`- [x] **done**`) → checkbox stripped first, then bold parsed
- Empty frontmatter body → show frontmatter block only, no empty note message
- Frontmatter with complex values (arrays, objects) → show as raw YAML text
- No frontmatter → no block shown
- `- [X]` uppercase → treated as checked

### Implementation order

1. Add Checked/Checkable to MarkdownLine
2. Update ParseMarkdown list branch
3. Update renderList
4. Add renderFrontmatter to viewer.go
5. Update SetContent
6. Write tests

## Implementation Notes

**Frontmatter parsing:** After M16b (YAML replacement), use the custom YAML parser to extract frontmatter fields for display. For complex values (arrays, objects), render as raw YAML text.

## Completion Criteria

- [ ] Checkbox parsing: `- [ ]` and `- [x]`/`- [X]` detected
- [ ] Checkboxes rendered with styled icons (✓ or ☐)
- [ ] Success color for checked, dimmed for unchecked
- [ ] Checkbox + inline formatting works (`- [x] **done**`)
- [ ] Frontmatter block rendered above note content
- [ ] Frontmatter has border and header
- [ ] Key:value pairs displayed correctly
- [ ] Complex values (arrays) shown as raw YAML
- [ ] No frontmatter → no block shown
- [ ] `make test` passes
- [ ] `make vet` exits 0
- [ ] Manual test: checkboxes and frontmatter render correctly
