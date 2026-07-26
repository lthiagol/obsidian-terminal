# M96 — Graph View (ASCII)

**Status:** ✅ done  
**Phase:** 99 — Future (Low Priority)  
**Priority:** 🔵 Low  
**Finding:** M47 deferred; STATUS non-goals updated — read-only ASCII graph is in charter  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) Track D

## Goal

Add `ModeGraph`: ASCII visualization of note link relationships with keyboard navigation and Enter-to-open, global and local scopes.

## Problem statement

Users cannot see vault connectivity at a glance. Backlink index exists (`vault.go` `Backlinks` map) but only reverse edges; no TUI to explore the graph.

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
| **Blocked by** | M50 (`openNote` for Enter — ✅ done), M59 (handlers split — graph handler goes in `handlers_graph.go`), M60 (use `MatchDown`/`MatchUp` helpers in graph key handler) |
| **Blocks** | M61 WP1 (KEYBINDINGS must document Ctrl+G + graph keys in module map) |
| **Parallel-safe with** | M48 (different mode — ✅ done) |

## Design (approved for execution — detailed 2026-06-21)

### New mode constant

Add to `model.go` after `ModeProfilePicker`:

```go
const (
    ModeBrowse Mode = iota
    ModeView
    ModeSearch
    ModeFind
    ModeHelp
    ModeTags
    ModeProfilePicker
    ModeGraph  // NEW — M96
)
```

Add to `Mode.String()`:
```go
case ModeGraph:
    return "GRAPH"
```

### New KeyMap field

Add to `KeyMap` struct in `keys.go`:

```go
GraphToggle tea.KeyType // GraphToggle opens the graph view (Ctrl+G).
```

Add to `DefaultKeys()`:
```go
GraphToggle: tea.KeyCtrlG,
```

### New types (in `graph.go`)

```go
// GraphScope determines whether the graph shows all notes or a local neighborhood.
type GraphScope int

const (
    GraphGlobal GraphScope = iota // GraphGlobal shows top-connected notes across the vault.
    GraphLocal                    // GraphLocal shows a center note and its neighbors.
)

// GraphNode represents a single note in the graph visualization.
type GraphNode struct {
    Path     string  // Path is the vault-relative note path.
    Label    string  // Label is the basename without .md extension.
    OutDeg   int     // OutDeg is the number of outgoing edges.
    InDeg    int     // InDeg is the number of incoming edges.
    ScreenX  int     // ScreenX is the column position on the rune grid.
    ScreenY  int     // ScreenY is the row position on the rune grid.
}

// GraphEdge represents a directed link from one note to another.
type GraphEdge struct {
    From string // From is the source note path.
    To   string // To is the target note path.
}

// GraphModel holds the state for the graph view mode.
type GraphModel struct {
    scope      GraphScope
    nodes      []GraphNode
    edges      []GraphEdge
    selected   int
    centerPath string // centerPath is the focus note for local scope.
}
```

### New `Model` fields

Add to `Model` struct in `model.go`:

```go
graph GraphModel  // graph holds state for ModeGraph (M96).
```

### Edge construction

Build **directed** edges A → B when note A's content contains `[[wiki-link]]` resolving to B:

```go
// BuildEdges constructs directed edges from wiki-links found in note contents.
// Uses searchIndex (path → content) and ResolveWikiLink for target resolution.
func BuildEdges(allPaths []string, searchIndex map[string]string, vault *VaultEntry, vaultPath string) []GraphEdge {
    var edges []GraphEdge
    for _, path := range allPaths {
        content, ok := searchIndex[path]
        if !ok {
            continue
        }
        for _, target := range extractWikiLinkTargets(content) {
            resolved := ResolveWikiLink(target, vault, vaultPath)
            if resolved != "" && resolved != path {
                edges = append(edges, GraphEdge{From: path, To: resolved})
            }
        }
    }
    return edges
}
```

**Note:** `extractWikiLinkTargets` and `ResolveWikiLink` are currently in `vault.go` / `wikilink.go` (root package). After M57, they become `vault.ExtractWikiLinkTargets` and `vault.ResolveWikiLink`. Use the appropriate path based on whether M57 has run.

### Local graph

```go
// BuildLocalGraph builds a graph centered on centerPath with its neighbors.
// Nodes: {center} ∪ out-neighbors ∪ in-neighbors (from backlinkIndex).
// Cap at 30 nodes; prioritize out-links then back-links.
func BuildLocalGraph(center string, edges []GraphEdge, backlinks map[string][]string, maxNodes int) GraphModel
```

### Global graph

```go
// BuildGlobalGraph builds a graph of the top-connected notes.
// Sort by (outDegree + inDegree) desc, take top maxNodes.
func BuildGlobalGraph(edges []GraphEdge, allPaths []string, maxNodes int) GraphModel
```

### Circular layout

```go
// LayoutCircle positions nodes in a circle on a width×height grid.
// Sets ScreenX, ScreenY on each node.
func LayoutCircle(nodes []GraphNode, width, height int) {
    if len(nodes) == 0 {
        return
    }
    cx := width / 2
    cy := height / 2
    r := min(cx, cy) * 2 / 3
    if r < 1 {
        r = 1
    }
    for i := range nodes {
        angle := 2 * math.Pi * float64(i) / float64(len(nodes))
        nodes[i].ScreenX = cx + int(float64(r)*math.Cos(angle))
        nodes[i].ScreenY = cy + int(float64(r)*math.Sin(angle))
    }
}
```

### ASCII rendering

```go
// RenderGrid renders the graph as an ASCII string on a width×height grid.
// Uses Bresenham line algorithm for edges. Selected node wrapped with [brackets].
func RenderGrid(nodes []GraphNode, edges []GraphEdge, selected int, width, height int, palette Palette) string
```

**Edge characters:** `─` (horizontal), `│` (vertical), `╲` (diag down-right), `╱` (diag up-right), `·` (fallback).

**Labels:** basename at node position, truncated to 12 chars. Selected node: `[label]` with `palette.Accent` color.

### Keybindings (ModeGraph only)

| Key | Action | Implementation |
|-----|--------|----------------|
| `Ctrl+G` | Enter graph from browse (global) or view (local on activeNote) | Global dispatch in `Update` |
| `j` / `↓` | Next node in `nodes` slice | `m.keys.MatchDown(msg)` (M60 helper) |
| `k` / `↑` | Previous node | `m.keys.MatchUp(msg)` (M60 helper) |
| `Enter` | `openNote(selected.Path)` → ModeView | `msg.Type == tea.KeyEnter` |
| `l` | Toggle global ↔ local (rebuild graph) | `MatchRune(msg, m.keys.RightRune)` — no conflict (graph mode only) |
| `f` | Focus: rebuild local graph centered on selected | `MatchRune(msg, 'f')` — unused in graph mode |
| `r` | Rebuild graph (refresh edges) | `MatchRune(msg, 'r')` — no conflict (graph mode only) |
| `Esc` | `mode = prevMode` | `msg.Type == tea.KeyEsc` |

**Conflict check (verified 2026-06-21):**
- `Ctrl+G` — not allocated anywhere (check `KEYBINDINGS.md` + `keys.go`)
- `f` — unused in any mode (check `grep -rn "MatchRune.*'f'" --include='*.go'`)
- `r` — used in broken vault retry (global dispatch, but graph mode is not broken-vault state — no conflict)
- `l` — used for tree expand in browse/view, but graph mode has no tree — no conflict

### Rendering integration

In `render_layout.go` `View()`, add a case for `ModeGraph`:

```go
case ModeGraph:
    rightPanel = m.renderGraph()
```

New method in `graph.go`:

```go
func (m Model) renderGraph() string {
    width := m.width - m.treeWidth - 2
    height := m.height - 2 // leave room for footer
    
    grid := RenderGrid(m.graph.nodes, m.graph.edges, m.graph.selected, width, height, m.palette)
    
    footer := lipgloss.NewStyle().Foreground(m.palette.TextDim).Render(
        "j/k move · Enter open · l scope · f focus · r refresh · Esc back",
    )
    
    scopeLabel := "global"
    if m.graph.scope == GraphLocal {
        scopeLabel = "local: " + truncate(m.graph.centerPath, 30)
    }
    status := lipgloss.NewStyle().Foreground(m.palette.Accent).Render(
        fmt.Sprintf("Graph (%s) — %d nodes", scopeLabel, len(m.graph.nodes)),
    )
    
    return lipgloss.JoinVertical(lipgloss.Left, status, "", grid, "", footer)
}
```

### Challenged decisions

| Original | Revised | Why |
|----------|---------|-----|
| All notes in global graph | Cap 60 by degree | 500+ nodes unreadable in ASCII |
| `f` reflow circle animation | `f` = local focus rebuild | Simpler; no animation infra |
| Parse vault on every open | Build from `searchIndex` | Already in memory from scan |
| 1 day estimate | 2–3 days | Layout + edges + mode + tests |
| `handlers.go` for graph handler | `handlers_graph.go` (post-M59) | M59 splits handlers by mode |

---

## Work packages

### WP1 — Graph data layer (3h)

**Steps:**
1. Create `graph.go` with `package main` header
2. Add types: `GraphScope`, `GraphNode`, `GraphEdge`, `GraphModel` (see Design section for exact code)
3. Implement `BuildEdges(allPaths, searchIndex, vault, vaultPath) []GraphEdge`:
   - Iterate `allPaths`; for each, get content from `searchIndex[path]`
   - Call `extractWikiLinkTargets(content)` → resolve each with `ResolveWikiLink(target, vault, vaultPath)`
   - Append `GraphEdge{From: path, To: resolved}` if resolved non-empty and not self-loop
4. Implement `BuildLocalGraph(center, edges, backlinks, maxNodes) GraphModel`:
   - Find out-neighbors: edges where `From == center` → collect `To` paths
   - Find in-neighbors: `backlinks[center]` (already available as `m.backlinkIndex`)
   - Nodes: `{center}` ∪ out-neighbors ∪ in-neighbors
   - Cap at `maxNodes` (30); if exceeded, prioritize out-links then back-links
   - Set `centerPath = center`, `scope = GraphLocal`
5. Implement `BuildGlobalGraph(edges, allPaths, maxNodes) GraphModel`:
   - Count out-degree and in-degree for each path
   - Sort paths by `(outDeg + inDeg)` descending
   - Take top `maxNodes` (60)
   - Filter edges to only those between selected nodes
   - Set `scope = GraphGlobal`
6. Create `graph_test.go` with unit tests (see Test plan below)

**Import budget for `graph.go`:**
- `math` (for `math.Pi`, `math.Cos`, `math.Sin` in layout)
- `sort` (for sorting nodes by degree)
- `strings` (for label truncation)
- `github.com/charmbracelet/lipgloss` (for rendering)

**Discovery commands for executing agent:**
```bash
# Verify extractWikiLinkTargets and ResolveWikiLink exist and their signatures:
grep -n 'func extractWikiLinkTargets' vault.go
grep -n 'func ResolveWikiLink' wikilink.go

# Check what backlinkIndex looks like on Model:
grep -n 'backlinkIndex' model.go
```

**Verification:**
- [ ] `graph_test.go` passes:
  - T1: 3-note chain A→B→C: `BuildEdges` returns 2 edges (A→B, B→C)
  - T2: Self-loops excluded (note with `[[self]]` link to itself → no edge)
  - T3: `BuildLocalGraph("B", edges, backlinks, 30)` returns nodes {A, B, C}
  - T4: `BuildGlobalGraph(edges, allPaths, 60)` caps at 60 nodes
  - T5: `BuildGlobalGraph` sorts by degree (highest-degree node is first)
- [ ] `make test && make vet` pass

---

### WP2 — Circular layout + ASCII render (4h)

**Steps:**
1. Implement `LayoutCircle(nodes, width, height)` (see Design section for exact code)
2. Implement `RenderGrid(nodes, edges, selected, width, height, palette) string`:
   - Create a 2D rune grid (`[][]rune`) of size `height × width`
   - Fill with spaces
   - For each edge: draw Bresenham line between `from.ScreenX/Y` and `to.ScreenX/Y`
   - Edge chars: `─`, `│`, `╲`, `╱`, `·` (based on line direction)
   - For each node: place truncated label (12 chars max) at `ScreenX/Y`
   - Selected node: wrap with `[` `]` and color with `palette.Accent`
   - Convert grid to string with newlines
3. Implement `truncate(label, maxLen) string` helper (truncate to maxLen-1 + `…` if too long)
4. Add tests to `graph_test.go`:
   - T6: 4 nodes → `LayoutCircle` sets non-zero ScreenX/Y on all nodes
   - T7: `RenderGrid` returns non-empty string containing all 4 labels
   - T8: Selected node rendered with `[` bracket (grep for `[` in output)
   - T9: 0 nodes → `RenderGrid` returns empty string, no panic
   - T10: 0 edges → `RenderGrid` returns node labels only, no edge chars

**Verification:**
- [ ] `graph_test.go` passes (T6–T10)
- [ ] `make test && make vet` pass

---

### WP3 — ModeGraph integration (4h)

**Steps:**
1. Add `ModeGraph` to mode constants in `model.go` (after `ModeProfilePicker`)
2. Add `case ModeGraph: return "GRAPH"` to `Mode.String()`
3. Add `graph GraphModel` field to `Model` struct
4. Add `GraphToggle tea.KeyType` to `KeyMap` struct in `keys.go`
5. Add `GraphToggle: tea.KeyCtrlG` to `DefaultKeys()` in `keys.go`
6. Add global `Ctrl+G` handler in `Update` (in `model.go`, before mode dispatch):
   ```go
   if msg.Type == m.keys.GraphToggle {
       m.enterGraphMode()
       return m, nil
   }
   ```
7. Implement `enterGraphMode()` in `graph.go`:
   ```go
   func (m *Model) enterGraphMode() {
       m.prevMode = m.mode
       edges := BuildEdges(m.allPaths, m.searchIndex, m.vault, m.config.VaultPath)
       if m.mode == ModeView && m.activeNote != nil {
           m.graph = BuildLocalGraph(m.activeNote.Path, edges, m.backlinkIndex, 30)
       } else {
           m.graph = BuildGlobalGraph(edges, m.allPaths, 60)
       }
       m.mode = ModeGraph
   }
   ```
8. Implement `handleGraphKey` in `handlers_graph.go` (new file, post-M59 pattern):
   ```go
   func (m Model) handleGraphKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
       switch {
       case msg.Type == tea.KeyEsc:
           m.mode = m.prevMode
           return m, nil
       case msg.Type == tea.KeyEnter:
           if m.graph.selected < len(m.graph.nodes) {
               m.openNote(m.graph.nodes[m.graph.selected].Path)
           }
           return m, nil
       case m.keys.MatchDown(msg):
           if m.graph.selected < len(m.graph.nodes)-1 {
               m.graph.selected++
           }
           return m, nil
       case m.keys.MatchUp(msg):
           if m.graph.selected > 0 {
               m.graph.selected--
           }
           return m, nil
       case MatchRune(msg, 'l'):
           m.toggleGraphScope()
           return m, nil
       case MatchRune(msg, 'f'):
           if m.graph.selected < len(m.graph.nodes) {
               m.focusGraphOn(m.graph.nodes[m.graph.selected].Path)
           }
           return m, nil
       case MatchRune(msg, 'r'):
           m.enterGraphMode() // rebuild
           return m, nil
       }
       return m, nil
   }
   ```
9. Implement `toggleGraphScope()` and `focusGraphOn(path)` in `graph.go`
10. Add `case ModeGraph: return m.handleGraphKey(msg)` to `Update` mode dispatch
11. Add `case ModeGraph: rightPanel = m.renderGraph()` to `View` in `render_layout.go`
12. Add `GRAPH` badge color in `statusbar.go` (use `m.palette.Accent` or add a dedicated color)
13. Add integration tests to `graph_test.go`:
    - T11: Ctrl+G from browse → ModeGraph, scope=global
    - T12: Ctrl+G from view → ModeGraph, scope=local, center=activeNote.Path
    - T13: `j` → `graph.selected` increments
    - T14: `Enter` → ModeView, `activeNote.Path` = selected node path
    - T15: `Esc` → `mode == prevMode`

**Import budget for `handlers_graph.go`:**
- `tea "github.com/charmbracelet/bubbletea"` (for `tea.KeyMsg`)

**Discovery commands:**
```bash
# Verify Ctrl+G is not already allocated:
grep -rn 'KeyCtrlG\|CtrlG\|GraphToggle' --include='*.go'
# Verify 'f' is not allocated in any mode:
grep -rn "MatchRune.*'f'" --include='*.go'
# Verify 'r' allocation (broken vault retry — check it's global, not mode-specific):
grep -rn "MatchRune.*'r'" --include='*.go'
```

**Verification:**
- [ ] `graph_test.go` passes (T11–T15)
- [ ] `model_test.go` passes (23 tests — mode transitions unchanged)
- [ ] `make test && make vet` pass
- [ ] `grep 'KeyCtrlG' keys.go` returns 1 match (GraphToggle in DefaultKeys)

---

### WP4 — Polish + caps + footer (2h)

**Steps:**
1. Implement `renderGraph()` in `graph.go` (see Design section for exact code)
2. Footer: `j/k move · Enter open · l scope · f focus · r refresh · Esc back`
3. Status line: `Graph (global) — 60 nodes` or `Graph (local: path/to/note.md) — 12 nodes`
4. Toast when global truncates: `m.addToast("Showing top 60 connected notes", ToastInfo)`
5. Empty vault / no edges: friendly message `"No wiki-link connections found in this vault"`
6. Test:
    - T16: 0-edge vault → graph mode shows friendly message, no panic
    - T17: Global graph with >60 notes → toast shown on enter

**Verification:**
- [ ] `graph_test.go` passes (T16–T17)
- [ ] `make test && make vet` pass

---

### WP5 — Tests + help + KEYBINDINGS prep (2h)

**Steps:**
1. Add to `graph_test.go`:
    - T18: `l` toggles scope (global → local centered on selected; local → global)
    - T19: `f` focuses on selected node (rebuild local graph)
    - T20: `r` rebuilds graph (edges refreshed)
2. Add graph section to `help.go` `buildHelpSections`:
   ```
   Graph Mode:
     Ctrl+G  Open graph view (browse: global, view: local)
     j/k     Move between nodes
     Enter   Open selected note
     l       Toggle global ↔ local scope
     f       Focus local graph on selected node
     r       Rebuild graph from indexes
     Esc     Return to previous mode
   ```
3. Update `KEYBINDINGS.md` with the graph mode section (or defer to M61 if M61 hasn't run)
4. Run full suite

**Verification:**
- [ ] `graph_test.go` passes (20 tests total: T1–T20)
- [ ] `help.go` includes graph section
- [ ] `KEYBINDINGS.md` includes graph keys (or note: deferred to M61)
- [ ] `make test && make vet` pass
- [ ] Total test count: 298 + 20 = ~318

---

## Files to modify

| File | Changes |
|------|---------|
| `graph.go` | **New** — types, BuildEdges, BuildLocalGraph, BuildGlobalGraph, LayoutCircle, RenderGrid, renderGraph, enterGraphMode, toggleGraphScope, focusGraphOn, truncate |
| `graph_test.go` | **New** — 20 tests (T1–T20) |
| `handlers_graph.go` | **New** — `handleGraphKey` (post-M59 pattern) |
| `model.go` | Add `ModeGraph` constant, `String()` case, `graph GraphModel` field, Ctrl+G dispatch in `Update` |
| `keys.go` | Add `GraphToggle` field + `tea.KeyCtrlG` in `DefaultKeys()` |
| `render_layout.go` | Add `case ModeGraph: rightPanel = m.renderGraph()` in `View` |
| `statusbar.go` | Add GRAPH badge color (use `m.palette.Accent` or dedicated) |
| `help.go` | Add graph key section to `buildHelpSections` |
| `KEYBINDINGS.md` | Document graph keys (or defer to M61) |
| `STATUS.md` | M96 → ✅ with test count delta (+20) |

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | Edges from wiki-links in searchIndex | unit | WP1 |
| T2 | Self-loops excluded | unit | WP1 |
| T3 | Local graph includes backlink neighbor | unit | WP1 |
| T4 | Global caps at 60 nodes | unit | WP1 |
| T5 | Global sorts by degree (highest first) | unit | WP1 |
| T6 | LayoutCircle sets non-zero ScreenX/Y | unit | WP2 |
| T7 | RenderGrid returns non-empty string with all labels | unit | WP2 |
| T8 | Selected node rendered with `[` bracket | unit | WP2 |
| T9 | 0 nodes → RenderGrid returns empty, no panic | unit | WP2 |
| T10 | 0 edges → RenderGrid shows labels only | unit | WP2 |
| T11 | Ctrl+G browse → ModeGraph, scope=global | integration | WP3 |
| T12 | Ctrl+G view → ModeGraph, scope=local | integration | WP3 |
| T13 | `j` increments selected | integration | WP3 |
| T14 | Enter opens note → ModeView | integration | WP3 |
| T15 | Esc restores prevMode | integration | WP3 |
| T16 | 0-edge vault shows friendly message | integration | WP4 |
| T17 | Global >60 notes → toast shown | integration | WP4 |
| T18 | `l` toggles scope | integration | WP5 |
| T19 | `f` focuses on selected | integration | WP5 |
| T20 | `r` rebuilds graph | integration | WP5 |

## Acceptance criteria (milestone done)

- [ ] WP1–WP5 complete
- [ ] Global graph renders capped node set with edges
- [ ] Local graph renders center + neighbors
- [ ] j/k selection, Enter open, l toggle, f focus, r refresh, Esc back
- [ ] Ctrl+G from browse (global) and view (local)
- [ ] Read-only — no vault writes
- [ ] 20 graph tests pass
- [ ] `make test && make vet` pass
- [ ] `STATUS.md` updated: M96 → ✅ with test count delta
- [ ] `KEYBINDINGS.md` updated (or deferred to M61)

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Unreadable dense graph | high | Node cap (60 global, 30 local) + degree filter |
| Slow BuildEdges on huge vault | medium | Run once on mode enter; `r` refresh is user-initiated |
| Edge parse mismatch vs viewer | low | Reuse `extractWikiLinkTargets` + `ResolveWikiLink` (same as viewer) |
| `Ctrl+G` conflict with existing binding | low | WP3 discovery command verifies; `keys.go` has no `KeyCtrlG` today |
| `f` conflict with existing binding | low | WP3 discovery command; `f` is unused in all modes |
| Bresenham line rendering looks bad | medium | WP4 polish; can switch to `·` dots if diagonal chars look noisy |
| ModeGraph color missing | low | Fallback to `m.palette.Accent` |

**Rollback:** `git revert` the WP commit. WP1–WP2 are pure additions (new files, no integration); WP3 is the integration point (revert WP3 + WP4 + WP5 together if integration is broken).

## Handoff notes

**Read first:**
- This milestone file (especially the Design section with exact function signatures)
- `vault.go` lines with `extractWikiLinkTargets` and `ResolveWikiLink` — the graph reuses these
- `model.go` Mode constants and `Update` dispatch — to see where ModeGraph fits
- **Check if M59 and M60 are done** — this milestone assumes `handlers_graph.go` pattern (M59) and `m.keys.MatchDown/MatchUp` helpers (M60). If not done, use `handlers.go` and manual `MatchKey || MatchRune` instead.

**Do not:**
- Add new dependencies (no graph libraries, no layout engines)
- Write to the vault (read-only constraint)
- Implement mouse-based zoom/pan (out of scope)
- Add heading-level links as separate nodes (out of scope)
- Update the graph in real-time (refresh on mode enter only)

**When stuck:**
- If `BuildEdges` returns 0 edges on a vault with links: check that `searchIndex` contains content (not just filenames). The content is stripped of frontmatter — `extractWikiLinkTargets` looks for `[[...]]` in body text.
- If `RenderGrid` looks garbled: start with just node labels (no edges), then add Bresenham lines one at a time. Use `·` for all edge chars first, then switch to directional chars.
- If `Ctrl+G` doesn't enter graph mode: check that the handler is in `Update` **before** the mode dispatch (global keys), not inside a mode handler.
- If `f` or `r` conflicts: run the discovery commands in WP3 to find the conflict. If `r` conflicts with broken-vault retry, move broken-vault retry to only fire when `vaultState == VaultStateBroken`.

## Estimated total

2–3 days (3h WP1 + 4h WP2 + 4h WP3 + 2h WP4 + 2h WP5 = 15h ≈ 2 days focused, 3 days with buffer)

## Priority

🔵 Low — Phase 99, execute when prioritized

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | 2026-06-21 |
| Completed | 2026-06-21 |
| Tests added | 20 |
| Notes | All WPs executed in one session. graph.go (458 lines), graph_test.go (495 lines). model.go → 411 lines (+11 for ModeGraph). `make build && make test && make vet` all pass. 323 total tests. Deviations: T14 uses `os.WriteFile` + `ScanVault` for real vault fixture; T17 similarly creates real notes; KEYBINDINGS.md + help.go + statusbar.go updated with Graph section. |
