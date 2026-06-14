# M52 — Decompose model.go

**Status:** ⏳ pending (finishes partial **M38**)  
**Finding:** F-3 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md)  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) §7

## Goal

Reduce `model.go` to **core Model + Init/Update + layout dispatch** (~400 lines). Move subsystems to focused files. Finish M38 decomposition without changing behavior.

## Problem statement

`model.go` is 1,013 lines. AGENTS.md guideline (~250 lines) was exceeded. Multiple subsystems (pins, outline, daily, recent, rescan, render) share one file, increasing merge conflict and agent error rate.

## Out of scope

- Changing Bubble Tea Update/View signatures
- Moving `Model` struct to another package (M57)
- Rewriting handlers logic — only file moves + split

## Dependencies

- **After:** M51 (palette threading stable before View extraction)
- **One WP per session/PR** — do not combine WP5 + WP6

## Revised success metrics (challenged)

| Metric | Original M38 | Revised |
|--------|--------------|---------|
| model.go lines | < 250 | **< 400** |
| handlers.go lines | (not specified) | **each split file < 250** |
| Note open path | single | ✅ already `loadNote` / `openNote` (M50) |

---

## Work packages

### WP1 — Extract `vault_rescan.go` (2h)

**Move:** `vaultStateFrom`, `checkVaultChanges`, `rescanVault`, `countFiles`

**Verification:**
- [ ] `model_e2e_test.go` watcher tests pass
- [ ] `wc -l model.go` decreased by ~80

---

### WP2 — Extract `pin_handler.go` (2h)

**Move:** `togglePin`, `openPinnedNote`, `cyclePinnedNext`, `cyclePinnedPrev`, `validatePins`

**Verification:**
- [ ] `pinned_test.go` pass
- [ ] No new imports cycle

---

### WP3 — Extract `outline_handler.go` (2h)

**Move:** `buildOutline`, `renderOutline`, `estimateYOffset`, `handleOutlineKey` (from handlers.go)

**Verification:**
- [ ] `outline_test.go` pass

---

### WP4 — Extract `daily_handler.go` + `recent_handler.go` (2h)

**Move daily:** `buildDailyNotePath`, `openDailyNote`  
**Move recent:** `addRecentNote`, `toggleRecents`, `openRecentNote`, `renderRecents`, `handleRecentsKey`

**Verification:**
- [ ] `daily_recent_test.go` pass

---

### WP5 — Extract `render_layout.go` (3h)

**Move:** `View`, `renderSearch`, `renderFind`, `renderSearchPanel`, `renderBrokenVaultScreen`, `renderScanErrors`, `showScanErrors`, `wordWrap`

**Keep in model.go:** `Update` mode dispatch only (may call render helpers)

**Verification:**
- [ ] Integration tests rendering pipeline pass
- [ ] Broken vault screen test in `main_test.go` pass

---

### WP6 — Split `handlers.go` by mode (4h)

| New file | Contents |
|----------|----------|
| `handlers_browse.go` | `handleBrowseKey` |
| `handlers_view.go` | `handleViewKey`, in-note search, history |
| `handlers_search.go` | search/find/help/tags/profile/outline/recents/backlink handlers |
| `handlers_note.go` | `loadNote`, `openNote`, `enter*Mode`, `switchToProfile`, `setTheme` |

**Verification:**
- [ ] `model_test.go`, `handlers` behavior unchanged
- [ ] Each file < 250 lines

---

### WP7 — Verification + DESIGN update (1h)

**Steps:**
1. `wc -l model.go handlers*.go *_handler.go render_layout.go vault_rescan.go`
2. Update DESIGN.md module map (or checkpoint for M53)
3. Mark M38 ✅ in STATUS when model.go < 400

**Verification:**
- [ ] `model.go` < 400 lines
- [ ] `make test && make vet` pass

---

## Acceptance criteria

- [ ] WPs 1–7 complete
- [ ] No behavior change (test suite green)
- [ ] Function → file map matches execution plan §7
- [ ] M38 marked ✅

## Rollback / risk

| Risk | Mitigation |
|------|------------|
| Circular imports | Extract domain helpers before UI |
| Huge PR | Strict one-WP-per-commit rule |

## Handoff notes

Use `git mv` mentally — pure moves first, then edits. Run tests after **every** WP.

## Estimated total

2–3 days (spread across 7 sessions)

## Priority

🟡 High (Track B, after M51)
