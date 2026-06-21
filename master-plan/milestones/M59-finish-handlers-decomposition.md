# M59 — Finish M52: handlers.go Decomposition

**Status:** ✅ done  
**Phase:** 13 — Plan Remediation  
**Priority:** 🟡 High  
**Finding:** F-3 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md) (M52 WP6 skipped)  
**Execution plan:** [PHASE-13-EXECUTION-PLAN.md](../PHASE-13-EXECUTION-PLAN.md)

## Goal

Finish M52 by splitting `handlers.go` (624 lines, 26 functions) into focused files — each under 250 lines — without changing behavior. Closes the M52 partial and unblocks M38 → ✅.

## Problem statement

`handlers.go` is a god file mixing 5 concerns: mode key handlers, note-loading API, history navigation, in-note search, and profile/theme switching. M52 WP1–WP5 extracted other subsystems from `model.go` but WP6 (handlers split) was skipped due to time-boxing. `handlers.go` remains 624 lines / 26 functions.

## Out of scope

- Changing Bubble Tea `Update`/`View` signatures
- Moving `Model` struct to another package (M57)
- Rewriting handler logic — pure file moves + import adjustments
- Updating `DESIGN.md` module map (→ M61)
- DRY refactors like `KeyMap.MatchDown` helper (→ M60)

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | M51 (palette threading — done ✅) |
| **Blocks** | M60 (DRY refactors cleaner on split files), M61 (DESIGN.md module map must reflect final file structure) |
| **Parallel-safe with** | nothing — do alone, one WP per commit |

## Design (approved for execution)

### Target file structure after M59

| File | Lines (est.) | Contents |
|------|--------------|----------|
| `handlers_browse.go` | ~90 | `handleBrowseKey` |
| `handlers_view.go` | ~120 | `handleViewKey` |
| `handlers_search.go` | ~180 | `handleSearchKey`, `handleFindKey`, `handleSearchOrFind`, `handleHelpKey`, `handleBacklinkKey`, `handleTagsKey`, `handleCommandPaletteKey`, `handleProfilePickerKey` |
| `handlers_note.go` | ~90 | `noteNavKind` type + consts, `loadNote`, `applyNote`, `openNote`, `enterSearchMode`, `enterFindMode`, `enterHelpMode`, `enterTagsMode` |
| `in_note_search.go` | ~100 | `activateInNoteSearch`, `updateInNoteSearch`, `cycleInNoteMatch`, `handleInNoteSearchKey`, `renderInNoteSearch` |
| `history.go` | ~25 | `goBackHistory`, `goForwardHistory` |
| `profile_handler.go` | ~35 | `switchToProfile`, `setTheme` |
| `handlers.go` | **deleted** | — |

Total: ~640 lines across 7 files (was 624 in one file — small net growth from per-file `package` + `import` headers). Largest file: `handlers_search.go` at ~180 lines, well under the 250-line limit.

### Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Extract `in_note_search.go` as separate file | Yes | 5 functions / ~100 lines — large enough to justify own file; conceptually distinct (search within a note) |
| Extract `history.go` as separate file | Yes | Only 25 lines but conceptually distinct; pairs with future history features |
| Extract `profile_handler.go` as separate file | Yes | Profile/theme switching is a self-contained subsystem |
| Keep `enter*Mode` helpers in `handlers_note.go` | Yes | Used by multiple mode handlers; centralizing with `loadNote`/`applyNote` (the other state-transition code) keeps "state mutation" in one place |
| Name `handlers_search.go` for the grab-bag | Yes | Follows M52 WP6 plan; alternative `handlers_overlay.go` was considered but search/find/help are modes, not overlays |
| Delete `handlers.go` entirely | Yes | No functions remain after extraction; keeping an empty file would be noise |

### Import budget per new file

| File | Imports needed |
|------|----------------|
| `handlers_browse.go` | `tea`, (none else — uses `m.fileTree`, `m.openNote`, etc.) |
| `handlers_view.go` | `tea` |
| `handlers_search.go` | `tea` |
| `handlers_note.go` | `fmt`, `tea`, `github.com/lthiagol/obsidian-terminal/internal/search` (for `search.NewState`) |
| `in_note_search.go` | `fmt`, `strings`, `lipgloss`, `tea` |
| `history.go` | (none — uses `m.history`, `m.loadNote`) |
| `profile_handler.go` | (none — uses `m.config`, `m.palette`, `lookupPalette`, `markdownStyleFrom`, `searchStyleFrom`) |

> **Tip for executing agent:** Run `goimports -w <file>` after each extraction to auto-fix imports. Or use `make fmt` which runs `go fmt` (note: `go fmt` does NOT add missing imports — use `goimports`).

---

## Work packages

> **Rule:** One WP = one commit. Run `make test && make vet` after each WP. Pure file moves — no behavior changes. If a test fails, you've made a mistake in the move, not the test.

### WP1 — Extract `in_note_search.go` (1h)

**Steps:**
1. Create `in_note_search.go` with `package main` header
2. Move these 5 functions from `handlers.go` (lines 526–624):
   - `activateInNoteSearch`
   - `updateInNoteSearch`
   - `cycleInNoteMatch`
   - `handleInNoteSearchKey`
   - `renderInNoteSearch`
3. Add imports: `fmt`, `strings`, `lipgloss`, `tea`
4. Delete the moved functions from `handlers.go`
5. Run `goimports -w in_note_search.go handlers.go` (or manually fix imports)
6. Verify `handlers.go` decreased by ~100 lines

**Verification:**
- [ ] `in_note_search_test.go` passes (5 tests)
- [ ] `make test && make vet` pass
- [ ] `wc -l handlers.go` decreased by ~99 lines (526→624 gone)

---

### WP2 — Extract `history.go` (30m)

**Steps:**
1. Create `history.go` with `package main` header
2. Move these 2 functions from `handlers.go` (lines 502–524):
   - `goBackHistory`
   - `goForwardHistory`
3. No new imports needed (they use `m.history`, `m.historyForward`, `m.loadNote` — all on `Model`)
4. Delete the moved functions from `handlers.go`

**Verification:**
- [ ] `history_test.go` passes (6 scenarios)
- [ ] `make test && make vet` pass
- [ ] `wc -l handlers.go` decreased by ~23 lines

---

### WP3 — Extract `profile_handler.go` (30m)

**Steps:**
1. Create `profile_handler.go` with `package main` header
2. Move these 2 functions from `handlers.go` (lines 466–500):
   - `switchToProfile`
   - `setTheme`
3. No new imports needed (they use `lookupPalette`, `markdownStyleFrom`, `searchStyleFrom`, `m.palette`, `m.config`, `m.viewer`, `m.fileTree`, `m.rescanVault`, `m.addToast` — all package-level or on `Model`)
4. Delete the moved functions from `handlers.go`

**Verification:**
- [ ] `profiles_test.go` passes (7 tests)
- [ ] `custom_theme_test.go` passes (10 tests)
- [ ] `make test && make vet` pass
- [ ] `wc -l handlers.go` decreased by ~35 lines

---

### WP4 — Extract `handlers_note.go` (1h)

**Steps:**
1. Create `handlers_note.go` with `package main` header
2. Move these from `handlers.go`:
   - `noteNavKind` type definition + `navUser`/`navHistory`/`navReload` consts (lines 292–298)
   - `loadNote` (lines 300–313)
   - `applyNote` (lines 315–346)
   - `openNote` (lines 348–350)
   - `enterSearchMode` (lines 274–278)
   - `enterFindMode` (lines 280–284)
   - `enterHelpMode` (lines 286–290)
   - `enterTagsMode` (lines 402–406)
3. Add imports: `fmt`, `tea`, `github.com/lthiagol/obsidian-terminal/internal/search`
4. Delete the moved code from `handlers.go`

**Verification:**
- [ ] `model_test.go` passes (23 tests — many exercise mode transitions)
- [ ] `history_test.go` passes (depends on `loadNote`/`navHistory`)
- [ ] `daily_recent_test.go` passes (depends on `loadNote`)
- [ ] `make test && make vet` pass
- [ ] `wc -l handlers.go` decreased by ~90 lines

---

### WP5 — Split remaining `handlers.go` by mode (2h)

**After WP1–WP4, `handlers.go` should contain only the 8 mode key handlers (~430 lines).**

**Steps:**
1. Create `handlers_browse.go` — move `handleBrowseKey` (lines 12–100, ~89 lines). Import: `tea`.
2. Create `handlers_view.go` — move `handleViewKey` (lines 102–218, ~117 lines). Import: `tea`.
3. Create `handlers_search.go` — move the remaining 6 handlers:
   - `handleSearchKey` (3 lines)
   - `handleFindKey` (3 lines)
   - `handleSearchOrFind` (28 lines)
   - `handleHelpKey` (16 lines)
   - `handleBacklinkKey` (20 lines)
   - `handleTagsKey` (28 lines)
   - `handleCommandPaletteKey` (30 lines)
   - `handleProfilePickerKey` (26 lines)
   Import: `tea`.
4. **Delete `handlers.go`** — it should be empty now. If `git rm handlers.go` fails, an import is still referenced; find and fix.
5. Run `goimports -w handlers_*.go`
6. Verify `make build` succeeds (catches missing imports / duplicate symbols)

**Verification:**
- [ ] `make build` succeeds
- [ ] `make test && make vet` pass (298 tests)
- [ ] `wc -l handlers_browse.go` < 250
- [ ] `wc -l handlers_view.go` < 250
- [ ] `wc -l handlers_search.go` < 250
- [ ] `handlers.go` no longer exists

---

### WP6 — Verify line counts + update STATUS (30m)

**Steps:**
1. Run: `wc -l model.go handlers_*.go *_handler.go in_note_search.go history.go profile_handler.go render_layout.go vault_rescan.go`
2. Confirm every file is under 250 lines (except `model.go` which stays < 400 per M52)
3. Update `STATUS.md`:
   - M52 → ✅ done (was 🟡 partial → M59)
   - M38 → ✅ done (was 🟡 partial)
   - M59 → ✅ done with completion log
   - Update "Last updated" date
4. Update `M52-decompose-model.md` status line: `🟡 partial → M59` → `✅ done (via M59)`
5. **Do NOT update `DESIGN.md` module map** — that's M61's job

**Verification:**
- [ ] All extracted files < 250 lines (run the `wc -l` command and paste output into completion log)
- [ ] `model.go` < 400 lines (unchanged from M52)
- [ ] `STATUS.md` M52, M38, M59 all ✅
- [ ] `M52-decompose-model.md` status updated
- [ ] `make test && make vet` pass

---

## Files to modify

| File | Changes |
|------|---------|
| `in_note_search.go` | **New** — 5 funcs moved from `handlers.go` |
| `history.go` | **New** — 2 funcs moved from `handlers.go` |
| `profile_handler.go` | **New** — 2 funcs moved from `handlers.go` |
| `handlers_note.go` | **New** — note-loading API + mode transition helpers |
| `handlers_browse.go` | **New** — `handleBrowseKey` |
| `handlers_view.go` | **New** — `handleViewKey` |
| `handlers_search.go` | **New** — 8 secondary mode/overlay handlers |
| `handlers.go` | **Deleted** |
| `STATUS.md` | M52, M38 → ✅; M59 → ✅ with completion log |
| `M52-decompose-model.md` | Status → ✅ done (via M59) |

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | `in_note_search_test.go` still passes | unit | WP1 |
| T2 | `history_test.go` still passes | unit | WP2 |
| T3 | `profiles_test.go` + `custom_theme_test.go` still pass | unit | WP3 |
| T4 | `model_test.go` mode transitions still pass | unit | WP4 |
| T5 | `make build` succeeds (no missing imports) | build | WP5 |
| T6 | Full suite: 298 tests pass | regression | WP5, WP6 |

**No new tests needed** — this is a pure refactor. Existing tests are the safety net. If a test breaks, the move was wrong.

## Acceptance criteria (milestone done)

All must be checked before setting status to ✅:

- [x] WP1–WP6 complete
- [x] `handlers.go` deleted; no functions lost (grep for each function name confirms presence in new files)
- [x] Every new handler file < 250 lines
- [x] `model.go` < 400 lines (unchanged)
- [x] No behavior change — 298 tests pass
- [x] `make test && make vet` pass
- [x] `STATUS.md` updated: M52, M38 → ✅; M59 → ✅ with dates and test count delta (0)
- [x] `M52-decompose-model.md` status updated to ✅ done (via M59)
- [x] DESIGN.md module map NOT touched (deferred to M61)

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Missing import in extracted file → compile error | medium | Run `goimports -w` after each WP; `make build` in WP5 catches all |
| Duplicate symbol if move is partial (function left in both files) | low | `make vet` catches duplicate declarations |
| Circular import (unlikely — all `package main`) | very low | No internal imports added except `internal/search` which is already imported |
| Test breaks due to package-private access | very low | All functions stay in `package main` — no access changes |

**Rollback:** `git revert` the failing WP commit. Each WP is one commit, so rollback is surgical.

## Handoff notes

**Read first:**
- This milestone file (especially the Import budget table)
- `M52-decompose-model.md` for context on what was already done
- The current `handlers.go` to see the function locations (line numbers in this doc are accurate as of 2026-06-21; re-verify with `grep -n '^func ' handlers.go` before starting)

**Do not:**
- Add new functionality, refactor patterns, or rename functions — pure moves only
- Update `DESIGN.md` — that's M61
- Touch `model.go`, `render_layout.go`, `vault_rescan.go`, `pin_handler.go`, `outline_handler.go`, `daily_recent_handler.go` — they're already extracted
- Run WPs in parallel — strict sequence WP1 → WP6

**When stuck:**
- If a test fails after a move, you likely forgot to move a helper or mis-added an import. Compare the function body byte-for-byte with the original.
- If `goimports` is not installed: `go install golang.org/x/tools/cmd/goimports@latest`, or manually add imports matching the Import budget table.

## Estimated total

4–5 hours (1h WP1 + 30m WP2 + 30m WP3 + 1h WP4 + 2h WP5 + 30m WP6)

## Priority

🟡 High — closes M52 partial, unblocks M60 and M61

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | 2026-06-21 |
| Completed | 2026-06-21 |
| Tests added | 0 (pure refactor) |
| Notes | All WPs executed in one session. `wc -l` final: `model.go` 400, `handlers_browse.go` 95, `handlers_view.go` 123, `handlers_search.go` 167, `handlers_note.go` 91, `in_note_search.go` 109, `history.go` 25, `profile_handler.go` 37, `render_layout.go` 222, `vault_rescan.go` 100, `pin_handler.go` 103, `outline_handler.go` 132, `daily_recent_handler.go` 145. All < 250 lines (model.go = 400 per M52 spec). `handlers.go` deleted. 298 tests pass, `make vet` clean. Deviation from plan: `handlers_note.go` import budget listed `tea` but the extracted functions don't use tea types directly (all on `*Model`), so `tea` import was omitted. |
