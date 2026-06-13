# M96 — Graph View (ASCII)

**Status:** ⏳ pending (deferred from M49)  
**Phase:** 99 — Future (Low Priority)  
**Priority:** 🔵 Low  
**Finding:** M47 deferred; STATUS non-goals updated — read-only ASCII graph is in charter  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) Track D

## Goal

Add `ModeGraph`: ASCII visualization of note link relationships with keyboard navigation and Enter-to-open, global and local scopes.

## Problem statement

Users cannot see vault connectivity at a glance. Backlink index exists (`vault.go` Backlinks map) but only reverse edges; no TUI to explore the graph.

## Out of scope

- Interactive force-directed layout, zoom/pan with mouse
- Graph editing (adding/removing links)
- Clustering by folder/tag
- Mermaid/layout engines; new graph dependencies
- More than **60 nodes** in global view (hard cap — see Design)
- Heading-level links (`[[note#heading]]`) as separate nodes
- Real-time graph update during mode (refresh on Enter graph only)

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | M50 (`openNote` for Enter); M51 recommended (accent colors); M52 optional (extract before `graph.go` lands) |
| **Blocks** | M53 WP1 (Ctrl+G, graph keys) |
| **Parallel-safe with** | M48 (different mode) |

## Design (approved for execution)

### New mode

```go
const ModeGraph Mode = iota + ... // append after existing modes

type GraphScope int
const (
    GraphGlobal GraphScope = iota
    GraphLocal
)

type GraphModel struct {
    scope       GraphScope
    nodes       []GraphNode   // ordered for j/k selection
    edges       []GraphEdge   // from -> to path
    selected    int
    centerPath  string        // local graph center / focus target
}

type GraphNode struct {
    Path   string
    Label  string   // basename without .md
    Angle  float64  // radians for layout
    ScreenX, ScreenY int
}
```

`Model` fields: `graph GraphModel`, `prevMode Mode` on enter.

### Edge construction

Build **directed** edges A → B when note A's content contains `[[wiki-link]]` resolving to B:

1. Iterate `m.allPaths` (or `searchIndex` keys)
2. For each path, get content from `searchIndex[path]` (already stripped frontmatter)
3. `extractWikiLinkTargets(content)` → resolve each with `ResolveWikiLink(target, vault, vaultPath)`
4. Edge `(path, resolved)` if resolved non-empty and not self-loop

**Local graph** (note open or center selected):

- Nodes: `{center}` ∪ out-neighbors ∪ in-neighbors (backlink index)
- Cap at 30 nodes; if exceeded, prioritize out-links then back-links

**Global graph:**

- Cap at **60 nodes** total
- Selection strategy: sort by `(outDegree + inDegree) desc`, take top 60
- Status line: `Graph: 60 of 412 notes (top by connectivity) — l local`

### Layout

Circular placement on canvas `width × height` of **viewer panel**:

```
cx = (viewerWidth) / 2
cy = (viewerHeight) / 2
R  = min(cx, cy) * 2 / 3
node i at angle = 2π * i / N
```

Map graph coords → rune grid (`graphWidth × graphHeight` char matrix).

**Edges:** Bresenham line on grid; use `·` or `─`/`│`/`╲`/`╱` for diagonals. Skip edges between nodes too far apart if clutter (optional WP4).

**Labels:** basename at node position; selected node wrapped with lipgloss accent `[name]`.

### Keybindings (ModeGraph only)

| Key | Action |
|-----|--------|
| `Ctrl+G` | Enter graph from browse (global) or view (local on activeNote) |
| `j` / `↓` | Next node in `nodes` slice |
| `k` / `↑` | Previous node |
| `Enter` | `openNote(selected.Path)` → ModeView |
| `l` | Toggle global ↔ local (rebuild graph; local centers on selected) |
| `f` | Focus: rebuild **local** graph centered on selected node |
| `r` | Rebuild graph (refresh edges from indexes) |
| `Esc` | `mode = prevMode` |

**Browse `l`** expands tree — no conflict (different mode).  
**View `f`** unused — no conflict.

Add `GraphToggle tea.KeyType = tea.KeyCtrlG` to KeyMap.

### Rendering

- New file `graph.go`: `BuildGraph`, `LayoutCircle`, `RenderASCII`, `handleGraphKey`
- `View()` when `mode == ModeGraph`: full viewer panel = graph + footer hints
- Mode badge: `GRAPH` with new palette color `ModeGraph` (add to Palette in M51 or use `ModeHelp` temporarily)

### Challenged decisions

| Original | Revised | Why |
|----------|---------|-----|
| All notes in global graph | Cap 60 by degree | 500+ nodes unreadable in ASCII |
| `f` reflow circle animation | `f` = local focus rebuild | Simpler; no animation infra |
| Parse vault on every open | Build from `searchIndex` | Already in memory from scan |
| 1 day estimate | 2–3 days | Layout + edges + mode + tests |

---

## Work packages

### WP1 — Graph data layer (3h)

**Steps:**
1. `graph.go`: types `GraphNode`, `GraphEdge`, `GraphModel`
2. `BuildEdges(vault, searchIndex, vaultPath) []GraphEdge`
3. `BuildLocalGraph(center string, edges, backlinks, maxNodes) GraphModel`
4. `BuildGlobalGraph(edges, allPaths, maxNodes) GraphModel`
5. Unit tests: 3-note chain A→B→C; local from B returns A,B,C

**Verification:**
- [ ] `graph_test.go`: edge extraction matches `testdata` vault links
- [ ] Self-loops excluded
- [ ] Global cap at 60 nodes

---

### WP2 — Circular layout + ASCII render (4h)

**Steps:**
1. `LayoutCircle(nodes, width, height)` sets ScreenX/Y on grid coordinates
2. `RenderGrid(edges, nodes, selected) string` — char matrix + Bresenham
3. Truncate labels to 12 chars
4. Selected node uses `m.palette.Accent` (or global until M51)

**Verification:**
- [ ] Test: 4 nodes render non-empty string containing all labels
- [ ] Selected node marker differs visually (grep ANSI or `[` bracket)

---

### WP3 — ModeGraph integration (4h)

**Steps:**
1. Add `ModeGraph` constant + `String()` + mode badge color
2. Model fields: `graph GraphModel`
3. Global Ctrl+G handler in `Update` (browse + view): set `prevMode`, build graph, `mode = ModeGraph`
4. `handleGraphKey` for j/k/Enter/l/f/r/Esc
5. Enter on node calls `openNote` then stays View (graph closes) **or** Enter → view note — **decision:** Enter opens note ModeView (exit graph)

**Verification:**
- [ ] Integration test: Ctrl+G → graph → j → Enter → ModeView + activeNote set
- [ ] Esc returns to prevMode

---

### WP4 — Polish + caps + footer (2h)

**Steps:**
1. Footer: `j/k move · Enter open · l scope · f focus · Esc back`
2. Toast when global truncates: `"Showing top 60 connected notes"`
3. Empty vault / no edges message

**Verification:**
- [ ] 0-edge vault shows friendly message, no panic

---

### WP5 — Tests + help + KEYBINDINGS prep (2h)

**Steps:**
1. `graph_test.go` additions: local/global toggle, focus rebuild
2. `help.go` graph section
3. Document keys for M53

**Verification:**
- [ ] 8+ graph tests total
- [ ] `make test && make vet` pass

---

## Files to modify

| File | Changes |
|------|---------|
| `graph.go` | **New** — build, layout, render |
| `graph_test.go` | **New** |
| `model.go` | ModeGraph, fields, View branch, Ctrl+G dispatch |
| `handlers.go` or `handlers_graph.go` | `handleGraphKey` |
| `keys.go` | `GraphToggle` |
| `theme.go` | `ModeGraph` color (or defer cosmetic to M51) |
| `help.go` | Graph key section |
| `statusbar.go` | GRAPH badge |

## Test plan

| ID | Scenario | WP |
|----|----------|-----|
| T1 | Edges from wiki-links in searchIndex | WP1 |
| T2 | Local graph includes backlink neighbor | WP1 |
| T3 | Global caps at 60 | WP1 |
| T4 | Layout places N nodes | WP2 |
| T5 | Ctrl+G browse → ModeGraph | WP3 |
| T6 | Ctrl+G view → local on activeNote | WP3 |
| T7 | j/k changes selected | WP3 |
| T8 | Enter opens note ModeView | WP3 |
| T9 | Esc restores prevMode | WP3 |
| T10 | `l` toggles scope | WP5 |

## Acceptance criteria

- [ ] All WPs verified
- [ ] Global graph renders capped node set with edges
- [ ] Local graph renders center + neighbors
- [ ] j/k selection, Enter open, l toggle, f focus, Esc back
- [ ] Ctrl+G from browse and view
- [ ] Read-only — no vault writes
- [ ] `make test && make vet` pass
- [ ] STATUS.md updated

## Rollback / risk

| Risk | Mitigation |
|------|------------|
| Unreadable dense graph | Node cap + degree filter |
| Slow BuildEdges on huge vault | Run once on mode enter; optional `r` refresh |
| Edge parse mismatch vs viewer | Reuse `extractWikiLinkTargets` + `ResolveWikiLink` |
| ModeGraph color missing | Fallback to Accent |

## Handoff notes

Implement WP1–WP2 without UI first (pure tests). WP3 wires mode last.

Alias links (`findAlias`) already handled by `ResolveWikiLink`.

Do not add fsnotify or new deps.

## Estimated total

2–3 days (15–17 hours)

## Completion log

| Field | Value |
|-------|-------|
| Started | — |
| Completed | — |
| Tests added | — |
