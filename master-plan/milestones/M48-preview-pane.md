# M48 — Note Preview Pane

**Status:** ⏳ pending

## Goal

Show a small preview of the currently highlighted note's content at the bottom of the tree panel or viewer panel, so users can peek at notes without opening them.

## Design

When browsing the file tree, a preview pane shows the body of the currently hovered `.md` file. The preview re-uses the existing markdown renderer at a narrower width. Toggle visibility with `v` key.

### Rendering

Two possible placements:
1. **Bottom of tree panel** — split the tree vertically: tree on top, preview below
2. **Right panel** — preview renders in the viewer area when no note is open

Option 1 matches obsitui's approach (preview pane with fade effect at bottom of sidebar). Option 2 is simpler to implement since we reuse the viewer panel.

Recommended: **Option 2** — when in browse mode and a `.md` file is highlighted, the right panel shows the preview instead of "Select a file to view".

### Keybinding

- `v` — toggle preview pane on/off (browse mode)

### Implementation

- `previewVisible` field already exists on Model
- Load note body for hovered item (without frontmatter)
- Render with existing markdown pipeline at viewer panel width
- Show with a thin horizontal separator above the preview

## Files to modify

| File | Changes |
|------|---------|
| `model.go` | Add `previewNote`, `renderPreview`; toggle in `handleBrowseKey` |
| `handlers.go` | Add `v` handler in `handleBrowseKey` |

## Completion Criteria

- [ ] `v` toggles preview pane in browse mode
- [ ] Preview shows rendered content of hovered markdown file
- [ ] Preview panel has a separator from the tree
- [ ] Preview updates as cursor moves
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

1-2 hours
