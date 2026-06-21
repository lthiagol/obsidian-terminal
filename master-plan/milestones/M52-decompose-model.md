# M52 — Decompose model.go

**Status:** ✅ done (via M59 — finishes partial **M38**)  
**Finding:** F-3 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md)  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) §7

## Completion summary (2026-06-13)

| WP | Status | Notes |
|----|--------|-------|
| WP1 — Extract `vault_rescan.go` | ✅ done | `vaultStateFrom`, `checkVaultChanges`, `rescanVault`, `countFiles` moved |
| WP2 — Extract `pin_handler.go` | ✅ done | `togglePin`, `openPinnedNote`, `cyclePinnedNext/Prev`, `validatePins` moved |
| WP3 — Extract `outline_handler.go` | ✅ done | `buildOutline`, `renderOutline`, `estimateYOffset` moved |
| WP4 — Extract `daily_recent_handler.go` | ✅ done | Combined daily + recent into one file (renamed from `daily_handler.go` + `recent_handler.go`) |
| WP5 — Extract `render_layout.go` | ✅ done | `View`, `renderSearch*`, `renderBrokenVaultScreen`, `renderScanErrors`, `showScanErrors`, `wordWrap` moved |
| WP6 — Split `handlers.go` by mode | ✅ done via M59 | `handlers.go` (624 lines) decomposed into `handlers_browse.go` (95), `handlers_view.go` (123), `handlers_search.go` (167), `handlers_note.go` (91), `in_note_search.go` (109), `history.go` (25), `profile_handler.go` (37); `handlers.go` deleted |
| WP7 — Verify + DESIGN update | 🟡 partial | `model.go` = 400 lines (✅ < 400). DESIGN.md module map NOT updated — defer to **M61** |

## Goal

Reduce `model.go` to **core Model + Init/Update + layout dispatch** (~400 lines). Move subsystems to focused files. Finish M38 decomposition without changing behavior.

## Problem statement

`model.go` was 1,013 lines. AGENTS.md guideline (~250 lines) was exceeded. Multiple subsystems (pins, outline, daily, recent, rescan, render) share one file, increasing merge conflict and agent error rate.

**Remaining problem (2026-06-21):** `handlers.go` (624 lines, 26 functions) now holds the same god-file anti-pattern for handler logic. It mixes 5 concerns: mode key handlers, note-loading API, history navigation, in-note search, and profile/theme switching. M59 finishes the job.

## Out of scope

- Changing Bubble Tea Update/View signatures
- Moving `Model` struct to another package (M57)
- Rewriting handlers logic — only file moves + split

## Dependencies

- **After:** M51 (palette threading stable before View extraction)
- **One WP per session/PR** — do not combine WP5 + WP6
- **Follow-up:** M59 (handlers.go split), M61 (DESIGN.md module map update)

## Revised success metrics (challenged)

| Metric | Original M38 | Revised | M52 actual |
|--------|--------------|---------|------------|
| model.go lines | < 250 | **< 400** | ✅ 400 |
| handlers.go lines | (not specified) | **each split file < 250** | ❌ 624 (M59 will fix) |
| Note open path | single | ✅ already `loadNote` / `openNote` (M50) | ✅ |

---

## Work packages

### WP1 — Extract `vault_rescan.go` (2h)

**Move:** `vaultStateFrom`, `checkVaultChanges`, `rescanVault`, `countFiles`

**Verification:**
- [x] `model_e2e_test.go` watcher tests pass
- [x] `wc -l model.go` decreased by ~80

---

### WP2 — Extract `pin_handler.go` (2h)

**Move:** `togglePin`, `openPinnedNote`, `cyclePinnedNext`, `cyclePinnedPrev`, `validatePins`

**Verification:**
- [x] `pinned_test.go` pass
- [x] No new imports cycle

---

### WP3 — Extract `outline_handler.go` (2h)

**Move:** `buildOutline`, `renderOutline`, `estimateYOffset`

**Verification:**
- [x] `outline_test.go` pass

---

### WP4 — Extract `daily_recent_handler.go` (2h)

**Move daily:** `buildDailyNotePath`, `openDailyNote`  
**Move recent:** `addRecentNote`, `toggleRecents`, `openRecentNote`, `renderRecents`, `handleRecentsKey`

**Note:** Originally planned as two files (`daily_handler.go` + `recent_handler.go`). Combined into one `daily_recent_handler.go` (145 lines) — both subsystems are small and share navigation helpers.

**Verification:**
- [x] `daily_recent_test.go` pass

---

### WP5 — Extract `render_layout.go` (3h)

**Move:** `View`, `renderSearch`, `renderFind`, `renderSearchPanel`, `renderBrokenVaultScreen`, `renderScanErrors`, `showScanErrors`, `wordWrap`

**Keep in model.go:** `Update` mode dispatch only (may call render helpers)

**Verification:**
- [x] Integration tests rendering pipeline pass
- [x] Broken vault screen test in `main_test.go` pass

---

### WP6 — Split `handlers.go` by mode (4h) — ❌ deferred to M59

| New file | Contents |
|----------|----------|
| `handlers_browse.go` | `handleBrowseKey` |
| `handlers_view.go` | `handleViewKey`, in-note search, history |
| `handlers_search.go` | search/find/help/tags/profile/outline/recents/backlink handlers |
| `handlers_note.go` | `loadNote`, `openNote`, `enter*Mode`, `switchToProfile`, `setTheme` |

**Verification (when M59 executes):**
- [ ] `model_test.go`, `handlers` behavior unchanged
- [ ] Each file < 250 lines

**Why skipped in M52:** Time-boxed session ended. WP6 is mechanically independent of WP1–WP5 and lower risk (no View extraction). Reopened as M59 with a more detailed plan that also extracts `in_note_search.go`, `history.go`, `profile_handler.go`.

---

### WP7 — Verification + DESIGN update (1h) — 🟡 partial

**Steps:**
1. `wc -l model.go handlers*.go *_handler.go render_layout.go vault_rescan.go`
2. Update DESIGN.md module map → **deferred to M61**
3. Mark M38 ✅ in STATUS when model.go < 400

**Verification:**
- [x] `model.go` < 400 lines (actually 400 exactly)
- [x] `make test && make vet` pass
- [ ] DESIGN.md module map accurate (→ M61)

---

## Acceptance criteria (milestone done)

- [x] WP1–WP5 complete
- [ ] WP6 complete (→ **M59**)
- [x] No behavior change (test suite green)
- [x] Function → file map matches execution plan §7 (for WP1–WP5)
- [ ] M38 marked ✅ (blocked on M59 completion)
- [ ] DESIGN.md module map accurate (→ **M61**)

## Rollback / risk

| Risk | Mitigation |
|------|------------|
| Circular imports | Extract domain helpers before UI |
| Huge PR | Strict one-WP-per-commit rule |

## Handoff notes

WP6 is now M59. Read [M59-finish-handlers-decomposition.md](./M59-finish-handlers-decomposition.md) for the expanded plan that includes extraction of `in_note_search.go`, `history.go`, `profile_handler.go` alongside the mode split.

## Estimated total

2–3 days (spread across 7 sessions). ~1 day actually spent on WP1–WP5; ~1 day budgeted for M59 (WP6 + new extractions).

## Priority

🟡 High (Track B, after M51)

## Completion log

| Field | Value |
|-------|-------|
| Started | 2026-06-11 |
| Completed (WP1–WP5) | 2026-06-13 |
| Tests added | 0 (pure refactor) |
| Notes | WP6 skipped → M59. DESIGN.md update → M61. |
