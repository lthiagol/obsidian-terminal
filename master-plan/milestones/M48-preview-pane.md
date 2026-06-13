# M48 — Note Preview Pane

**Status:** ⏳ pending  
**Phase:** 11 — Visual & UX  
**Priority:** 🟡 High  
**Finding:** M47 deferred item; stub field `previewVisible` in `model.go:150`  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) Track D

## Goal

In browse mode, show a rendered preview of the tree-selected markdown note in the right panel without opening it (ModeView), toggled with `v`.

## Problem statement

Users must open every note to skim content. M47 added `previewVisible bool` but never wired it. Right panel shows static `"Select a file to view"` in browse mode (`model.go:477`).

## Out of scope

- Preview in view mode (already full note)
- Tree-bottom split layout (Option 1) — rejected; see Design
- Fade/gradient effects (terminal limitation; plain truncate footer ok)
- Wiki-link Tab/Enter in preview (preview is read-only peek, not navigable)
- Embed resolution in preview (`![[embed]]` — skip or show placeholder)
- Session persistence of preview toggle
- Horizontal scroll in preview

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | **M50** (stable `loadNote`); **M51** recommended (preview uses `m.palette`) |
| **Blocks** | M53 WP1 (`v` key in KEYBINDINGS) |
| **Parallel-safe with** | M49 only after M48 WP1 design frozen (separate modes) |

## Design (approved for execution)

### Placement: Option 2 — right panel (browse mode)

When `mode == ModeBrowse` && `previewVisible && selected file is .md`:

```
┌─ tree ─┬─ preview panel ─────────────┐
│ …      │ ── preview: notes/foo.md ── │
│ ▶ file │ # Title                     │
│        │ rendered markdown (viewport)│
└────────┴──────────────────────────────┘
```

When `previewVisible == false` → keep `"Select a file to view"`.

When selected entry is **directory** → `"Preview (v): select a markdown file"`.

### Keybinding

| Key | Mode | Action |
|-----|------|--------|
| `v` | Browse | Toggle `previewVisible` |

**Conflict check:** `v` is unused in codebase (verified 2026-06-13). Add to `KeyMap.PreviewToggle`.

### State fields

```go
// model.go — extend existing stub
previewVisible bool
previewPath    string   // last loaded path (cache key)
previewCache   string   // rendered ANSI at current viewer width
previewScroll  int      // optional; WP4 — scroll within long previews
```

Use a **dedicated** `previewViewport` (or second `MarkdownViewer` named `previewPane`) — do **not** reuse main `viewer` (would conflict when switching to ModeView).

### Render pipeline (preview)

1. `LoadNote(vaultPath, relPath)` — use `note.Body` (frontmatter stripped)
2. `ParseMarkdown` → `RenderMarkdown` at width `m.width - m.treeWidth - 2`
3. **No** embed resolver, **no** wiki-link selection
4. Cap render: first **N lines** after render (e.g. 80 lines) to avoid lag on huge notes — show `… (open Enter for full note)` footer

### Performance

- Reload only when `SelectedEntry().Path` changes **and** differs from `previewPath`
- On tree `j`/`k`, skip reload if same path
- Invalidate cache on `rescanVault` (clear `previewPath`)

### Header separator

```go
header := lipgloss.NewStyle().Foreground(palette.TextMuted).
    Render("── preview: " + path + " ──")
```

---

## Work packages

### WP1 — Preview pane type + KeyMap (2h)

**Steps:**
1. Add `PreviewToggle rune` (`'v'`) to `KeyMap` + `DefaultKeys()`
2. Create `preview.go` with `PreviewPane` struct: viewport + `SetContent(body, width)` + `View()`
3. Add `preview PreviewPane` field on `Model`; init in `NewModel`
4. Wire `handleBrowseKey`: `v` toggles `previewVisible`; toast `"Preview on/off"`

**Verification:**
- [ ] `keys_test.go` or new test: `v` toggles flag
- [ ] `make test && make vet` pass

---

### WP2 — Browse View integration (2h)

**Steps:**
1. In `View()` browse branch: if `previewVisible`, call `m.updatePreview()` then render `preview.View()` instead of placeholder
2. `updatePreview()`: if entry nil or dir → set placeholder message; if `.md` → load + render
3. On `fileTree.MoveUp/Down` in browse handler, preview updates next frame via View (or call `updatePreview` after move)

**Verification:**
- [ ] Manual: browse + `v` shows rendered note
- [ ] Directory selected → hint text, no panic
- [ ] `make test && make vet` pass

---

### WP3 — Cache + line cap (1.5h)

**Steps:**
1. Skip `LoadNote` if `entry.Path == m.previewPath` && width unchanged
2. Truncate rendered output to max 80 lines (config constant `previewMaxLines`)
3. Clear preview cache in `rescanVault`

**Verification:**
- [ ] Test: two updates same path → single load (use load counter mock or path assertion)
- [ ] Large note truncates with footer message

---

### WP4 — Tests + help + status bar (1.5h)

**Steps:**
1. `preview_test.go`: toggle, md preview content contains title, dir placeholder
2. Add `v` to `help.go` browse section
3. Status bar hint when preview on: `v hide preview`

**Verification:**
- [ ] 4+ tests pass
- [ ] help.go mentions preview (M53 will sync KEYBINDINGS)

---

## Files to modify

| File | Changes |
|------|---------|
| `preview.go` | **New** — PreviewPane widget |
| `model.go` | Fields, View branch, cache, rescan invalidation |
| `handlers.go` | `v` in `handleBrowseKey` |
| `keys.go` | `PreviewToggle` |
| `help.go` | Browse help entry |
| `statusbar.go` | Optional preview indicator |
| `preview_test.go` | **New** |

## Test plan

| ID | Scenario | WP |
|----|----------|-----|
| T1 | `v` toggles `previewVisible` | WP1 |
| T2 | `.md` selected → preview contains heading text | WP2 |
| T3 | Directory selected → placeholder string | WP2 |
| T4 | Same path twice → no double load (cache) | WP3 |
| T5 | Note >80 lines → truncation footer | WP3 |
| T6 | Enter opens note → ModeView unchanged behavior | WP4 |

## Acceptance criteria

- [ ] All WPs verified
- [ ] `v` toggles preview in browse mode only
- [ ] Preview shows rendered markdown for selected `.md`
- [ ] Separator/header shows path
- [ ] Updates when tree cursor moves
- [ ] No embed/wiki-link navigation in preview
- [ ] `make test && make vet` pass
- [ ] STATUS.md updated

## Rollback / risk

| Risk | Mitigation |
|------|------------|
| Slow preview on huge notes | Line cap + cache (WP3) |
| Palette globals wrong after M51 | Prefer M51 first, or use `m.palette` in preview |
| Main viewer state corruption | Separate PreviewPane instance |

## Handoff notes

Read `model.go` View() overlay priority: command palette > recents > scan errors > broken vault **before** browse preview. Preview only applies in default browse branch.

Do not implement graph (M49) in same PR.

## Estimated total

7–8 hours (was underestimated at 1–2h)

## Completion log

| Field | Value |
|-------|-------|
| Started | — |
| Completed | — |
| Tests added | — |
