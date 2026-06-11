# M27 — Checkboxes + Frontmatter Display

**Status:** ⏳ pending

## Goal

Render `- [ ]` and `- [x]` list items with styled checkbox icons, and display parsed frontmatter metadata at the top of the viewer.

## Steps

### 1. Checkbox rendering in lists

Add a `IsCheckbox` and `Checked` field to `MarkdownLine`. Detect `- [ ]` and `- [x]` patterns during list parsing. Render unchecked boxes as `☐` and checked boxes as `☑` with distinct styling (dimmed for unchecked, success color for checked).

### 2. Frontmatter display

At the top of the viewer, above the rendered markdown, show a metadata block:

```
Title         Tags           Aliases
My Note       tag1, tag2     alias1
```

Only shown if the note has frontmatter. Styled with dimmed colors to distinguish from note content. Toggle with `m` key.

### 3. Parser changes

- `parseListItem` — detect `[ ]` and `[x]` markers after the bullet
- `MarkdownLine` — add `IsCheckbox bool`, `Checked bool` fields
- `renderList` — use checkbox icons when `IsCheckbox` is true

## Completion Criteria

- [ ] `- [ ]` renders as unchecked checkbox icon
- [ ] `- [x]` renders as checked checkbox icon  
- [ ] Frontmatter metadata shown above rendered note
- [ ] `m` toggles frontmatter display
- [ ] `make test && make vet` pass
