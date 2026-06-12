# M49 — Graph View (ASCII)

**Status:** ⏳ pending

## Goal

Render a graph view of note connections using ASCII art — nodes positioned in a circle with Bresenham-style lines for edges. Supports global graph (all notes) and local graph (current note's connections). Navigate with keyboard, open notes by selecting nodes.

## Design

### Layout algorithm

- Place N nodes evenly around a circle of radius R
- Draw edges between connected nodes using line-drawing algorithm
- Nodes display abbreviated note names

### Modes

- **Global graph**: all notes and their wiki-link connections from backlink index
- **Local graph**: only the current note and its direct connections

### Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+G` | Open graph view (global or local if a note is open) |
| `j`/`k` | Select node down/up |
| `Enter` | Open selected note |
| `l` | Toggle local / global |
| `f` | Focus on selected node |
| `Esc` | Close, return to previous mode |

### Drawing

- Canvas size: viewer panel dimensions
- Node positions: circular layout with radius = min(width, height) / 3
- Nodes: labeled circles rendered as `(name)` or styled label
- Edges: dashed or dotted ASCII lines between connected nodes
- Selected node: highlighted with accent color

## Files to modify

| File | Changes |
|------|---------|
| New: `graph.go` | Graph state, layout, rendering |
| `model.go` | Add `ModeGraph`, graph state fields |
| `handlers.go` | Add `handleGraphKey`, dispatch from `handleBrowseKey`/`handleViewKey` |
| `keys.go` | Add `GraphToggle` key (`Ctrl+G`) |

## Completion Criteria

- [ ] Global graph renders with all note nodes and edges
- [ ] Local graph renders with current note + direct connections
- [ ] Node selection works with j/k
- [ ] Enter opens selected note
- [ ] `l` toggles local/global
- [ ] `Esc` returns to previous mode
- [ ] Ctrl+G opens graph from browse or view mode
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

1 day
