# M19 — Backlinks Panel

**Status:** ⏳ pending

## Goal

Show backlinks (notes linking TO the current note) in a side panel when viewing a note. Click a backlink to navigate to it.

## Steps

### 1. Build reverse link index

When viewing a note, scan all notes in the search index for wiki-links pointing to the current note's path. Cache the reverse index in Model.

### 2. Add backlinks panel

Show backlinks in a bottom panel or right sidebar below the viewer. Each backlink shows the source note name and the context line containing the link.

### 3. Add navigation

Arrow keys + Enter to select and follow a backlink. `Tab` to switch focus between viewer and backlinks panel.

## Completion Criteria

- [ ] Backlinks displayed when viewing a note with incoming links
- [ ] "No backlinks" shown when none found
- [ ] Click/Enter to navigate to a backlink source
- [ ] `make test && make vet` pass
