# M61 — Finish M53: Doc Sync Completion

**Status:** ✅ done  
**Phase:** 13 — Plan Remediation  
**Priority:** 🟡 High  
**Finding:** F-5 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md) (M53 WP2 incomplete)  
**Execution plan:** [PHASE-13-EXECUTION-PLAN.md](../PHASE-13-EXECUTION-PLAN.md)

## Goal

Close the M53 partial by fixing the stale `ARCHITECTURE.md` module map and the `AGENTS.md` styling section, then optionally rename `ARCHITECTURE.md` → `ARCHITECTURE.md` to align with project naming conventions.

## Problem statement

M53 was marked ✅ in STATUS but WP2 (ARCHITECTURE.md module map) was never actually completed. As of 2026-06-21:

1. **`ARCHITECTURE.md` module map** (lines 63–89) still has:
   - A callout: `> **M52 pending:** Outline, recents, pinned notes, daily note functions currently live in model.go and handlers.go` — **M52 is done** (these are now in `outline_handler.go`, `daily_recent_handler.go`, `pin_handler.go`)
   - A callout: `> **M51 pending:** Theme colors currently use global activatePalette()` — **M51 is done** (colors come from `m.palette`)
   - `model.go` row says "note opening, daily notes, recent notes, pinned notes" — these are extracted
   - `handlers.go` row says "mode-specific key handlers, note-loading API, history navigation" — but `handlers.go` is split by M59
   - Missing files: `vault_rescan.go`, `pin_handler.go`, `outline_handler.go`, `daily_recent_handler.go`, `render_layout.go`, `preview.go`, `in_note_search.go` (M59), `history.go` (M59), `profile_handler.go` (M59), `handlers_browse.go`/`handlers_view.go`/`handlers_search.go`/`handlers_note.go` (M59), `textinput.go` (M60)
   - References phantom files in the review: `outline.go`, `daily.go`, `pins.go`, `recents.go` (these never existed)

2. **`AGENTS.md` styling section** says: "All colors from `theme.go` constants: `Accent`, `AccentSecondary`, `TextPrimary`, `TextSecondary`, `TextDim`, etc." and "Use the pre-defined styles: `TreeStyle`, `ViewerStyle`, `StatusStyle`" — but M51 made these **deprecated globals**. All UI code now reads from `m.palette` / `p.Accent` / `m.palette.TreeStyle`. Following the AGENTS.md guidance would reintroduce the bug M51 fixed.

3. **Filename convention drift:** The master-plan template uses `{ARCHITECTURE or DESIGN doc}` as a placeholder, AGENTS.md describes ARCHITECTURE.md as "Architecture reference", and `ARCHITECTURE-REVIEW-*.md` files exist in `master-plan/`. The filename `ARCHITECTURE.md` is the odd one out.

## Out of scope

- Rewriting the `ARCHITECTURE.md` architecture narrative (only module map + callout fixes)
- Changing the actual architecture (docs only — no code)
- User-facing README changes (M53 WP3 already did README)
- Translating docs

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | M59 (handlers.go split — module map must reflect final file structure), M60 (new `textinput.go` + `keys.go` methods to document) |
| **Blocks** | nothing (this is the last Phase 13 milestone) |
| **Parallel-safe with** | nothing — must come last so docs reflect reality |

## Design (approved for execution)

### Module map update plan

The updated module map (in `ARCHITECTURE.md` → `ARCHITECTURE.md`) will:

1. **Remove** both "M51 pending" and "M52 pending" callout boxes
2. **Add** a "Post-M52 file structure" note explaining the extraction pattern
3. **Update** the `model.go` row to reflect its trimmed responsibility
4. **Update** the `handlers.go` row → replace with rows for `handlers_browse.go`, `handlers_view.go`, `handlers_search.go`, `handlers_note.go`, `in_note_search.go`, `history.go`, `profile_handler.go`
5. **Add rows** for `vault_rescan.go`, `pin_handler.go`, `outline_handler.go`, `daily_recent_handler.go`, `render_layout.go`, `preview.go`, `textinput.go`
6. **Verify** every file in the module map exists: `for f in <listed files>; do test -f "$f" || echo "MISSING: $f"; done`

### Rename decision (recommended: yes)

**Recommendation: rename `ARCHITECTURE.md` → `ARCHITECTURE.md`.**

**Rationale:**
- Aligns with `master-plan/ARCHITECTURE-REVIEW-*.md` naming family
- Matches the master-plan template placeholder `{ARCHITECTURE or DESIGN doc}`
- Matches AGENTS.md's own description: "Architecture reference: See ARCHITECTURE.md"
- Convention in Go/Rust communities is `ARCHITECTURE.md`

**Mechanical scope:**
```bash
git mv ARCHITECTURE.md ARCHITECTURE.md
rg -l 'DESIGN\.md' --glob '*.md'   # find references to update
```
Expected references to update: `AGENTS.md`, `README.md`, `master-plan/STATUS.md`, `master-plan/PHASE-12-EXECUTION-PLAN.md`, `master-plan/milestones/*.md` (those that reference ARCHITECTURE.md).

### Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Rename `ARCHITECTURE.md` → `ARCHITECTURE.md` | **Yes (recommended)** | Naming consistency; one-time mechanical cost |
| Keep module map as a table | Yes | Current format works; only content is stale |
| Add a "File extraction history" section | No | The milestone files already record this; don't duplicate |
| Update `AGENTS.md` styling section | Yes (in WP3) | Currently guides agents toward deprecated globals — active risk |

---

## Work packages

### WP1 — Fix `ARCHITECTURE.md` module map (1.5h)

**Steps:**
1. Open `ARCHITECTURE.md`. Delete the two callout boxes:
   - Line ~63: `> **M52 pending:** Outline ...` (M52 done — delete entirely)
   - Line ~246: `> **M51 pending:** Theme colors ...` (M51 done — delete entirely)
2. Update the Root package table:
   - `model.go` row: change responsibility to "Central `Model` struct, `Init/Update` dispatch, mode constants, `tickCmd`" (drop "note opening, daily notes, recent notes, pinned notes, toast system, command palette cmd, layout")
   - `handlers.go` row: **delete** (file is gone after M59)
   - Add rows for (post-M59/M60):
     - `handlers_browse.go` | Browse mode key handler | `handleBrowseKey`
     - `handlers_view.go` | View mode key handler | `handleViewKey`
     - `handlers_search.go` | Secondary mode/overlay handlers (search, find, help, tags, backlinks, command palette, profile picker) | `handleSearchKey`, `handleTagsKey`, …
     - `handlers_note.go` | Note-loading API + mode transition helpers | `loadNote`, `openNote`, `enter*Mode`
     - `in_note_search.go` | In-note search overlay | `activateInNoteSearch`, `handleInNoteSearchKey`
     - `history.go` | Navigation history back/forward | `goBackHistory`, `goForwardHistory`
     - `profile_handler.go` | Profile switching + theme application | `switchToProfile`, `setTheme`
     - `vault_rescan.go` | Vault state machine + rescan logic | `checkVaultChanges`, `rescanVault`
     - `pin_handler.go` | Pinned notes subsystem | `togglePin`, `cyclePinned*`, `validatePins`
     - `outline_handler.go` | Outline/TOC builder + renderer | `buildOutline`, `renderOutline`
     - `daily_recent_handler.go` | Daily notes + recent notes overlay | `openDailyNote`, `toggleRecents`, `renderRecents`
     - `render_layout.go` | `View()` + panel renderers | `View`, `renderSearch*`, `renderBrokenVaultScreen`
     - `preview.go` | Note preview pane (M48) | `renderPreview`
     - `textinput.go` | Shared text-input handler (M60) | `HandleTextInput`
3. Update the "Note opening" data flow section if it still references `openNote()` living in `handlers.go` — it now lives in `handlers_note.go`.
4. Update the "Theme System" section: remove any text saying "globals updated by `activatePalette()`". Replace with "colors read from `Model.palette` (set by `setTheme` in `profile_handler.go`)".
5. **Verify every file in the module map exists.** Run:
   ```bash
   # Extract filenames from the table and check each
   for f in model.go main.go handlers_browse.go handlers_view.go handlers_search.go handlers_note.go in_note_search.go history.go profile_handler.go vault_rescan.go pin_handler.go outline_handler.go daily_recent_handler.go render_layout.go preview.go textinput.go tree.go viewer.go viewport.go vault.go session.go config.go theme.go keys.go mouse.go backlinks.go tags.go statusbar.go help.go toast.go command_palette.go wikilink.go yamlmini.go profile_picker.go; do
     test -f "$f" || echo "MISSING: $f"
   done
   ```
   Fix any MISSING entries (either the file doesn't exist yet — check M59/M60 status — or the name in the table is wrong).

**Verification:**
- [ ] No `> **M5* pending:**` callouts remain in `ARCHITECTURE.md`
- [ ] No reference to `outline.go`, `daily.go`, `pins.go`, `recents.go` (phantom files)
- [ ] Every file in the module map exists (script above prints no MISSING)
- [ ] `make test && make vet` pass (docs only — no code, but verify nothing broke)

---

### WP2 — Rename `ARCHITECTURE.md` → `ARCHITECTURE.md` (30m)

**Skip this WP if the owner decides against the rename.** Document the decision in the completion log either way.

**Steps:**
1. `git mv ARCHITECTURE.md ARCHITECTURE.md`
2. Find all references: `rg -l 'DESIGN\.md' --glob '*.md'`
3. Update each reference to `ARCHITECTURE.md`:
   - `AGENTS.md` — the "Architecture reference" line
   - `README.md` — if it links to ARCHITECTURE.md
   - `master-plan/STATUS.md` — header link
   - `master-plan/PHASE-12-EXECUTION-PLAN.md` — any references
   - `master-plan/milestones/*.md` — any milestone that links to ARCHITECTURE.md (grep to find which)
4. Verify no stale links: `rg 'DESIGN\.md'` should return 0 results (or only historical references in `ARCHITECTURE-REVIEW-*.md` which are point-in-time and OK)

**Verification:**
- [ ] `test -f ARCHITECTURE.md && ! test -f ARCHITECTURE.md` passes
- [ ] `rg 'DESIGN\.md' --glob '*.md' | grep -v ARCHITECTURE-REVIEW` returns 0 matches (review files are point-in-time, OK to leave)
- [ ] All markdown links resolve (open `AGENTS.md`, click the architecture link — should go to `ARCHITECTURE.md`)

---

### WP3 — Verify `AGENTS.md` reflects post-M59/M60 reality (30m)

> **Note:** The AGENTS.md simplification was already done in the 2026-06-21 planning session (styling section rewritten to reference `m.palette`, duplicate bench entries removed, M85–M99 range fixed, master-plan section trimmed). This WP is a **verification + touch-up** after M59/M60 land new files.

**Steps:**
1. Open `AGENTS.md`. Verify the "Patterns & Conventions" section:
   - Styling section should say "colors from `Model.palette`" (not `theme.go` globals) — if not, fix it
   - Should NOT list `Accent`, `TreeStyle`, `ViewerStyle`, `StatusStyle` as "use these pre-defined styles" — they're deprecated
2. Verify the "Navigation History" section still matches `handlers_note.go` (post-M59):
   - `loadNote(path, kind)` now lives in `handlers_note.go` — update the file reference if needed
   - `openNote(path)` still delegates to `loadNote(path, navUser)`
3. If M60 added `textinput.go` and `KeyMap.MatchDown/Up/Left/Right`, consider adding a one-line mention in the "Keybindings" section: "For navigation keys, use `m.keys.MatchDown(msg)` / `MatchUp` / `MatchLeft` / `MatchRight` instead of manual `MatchKey || MatchRune`."
4. Verify the Commands section still matches `Makefile` (no drift).

**Verification:**
- [ ] No reference to deprecated globals (`Accent`, `TreeStyle`, etc.) in AGENTS.md "Styling" section
- [ ] "Navigation History" section references `handlers_note.go` (or just `handlers*.go`)
- [ ] `make test && make vet` pass (docs only)

---

### WP4 — Final STATUS + milestone audit (30m)

**Steps:**
1. Re-count tests: `go test ./... -v -count=1 | grep -c '^--- PASS'` — update STATUS "Total Tests" if drifted
2. Verify every Phase 12 + Phase 13 milestone marked ✅ in STATUS has its milestone file also ✅ (or 🟡 partial with follow-up)
3. Verify M53 milestone file status → ✅ done (was 🟡 partial → M61)
4. Update `STATUS.md`:
   - M53 → ✅ done (was 🟡 partial → M61)
   - M61 → ✅ done with completion log
   - Update "Last updated" date
5. Update `M53-documentation-sync.md` status line: `🟡 partial → M61` → `✅ done (via M61)`
6. Spot-check 3 random ✅ milestones: do their acceptance criteria match the code?

**Verification:**
- [ ] Test count in STATUS matches `go test` output ±2
- [ ] Every ✅ in STATUS has ✅ in milestone file (or 🟡 with documented follow-up)
- [ ] M53 and M61 both ✅ in STATUS and milestone files
- [ ] `make test && make vet` pass

---

## Files to modify

| File | Changes |
|------|---------|
| `ARCHITECTURE.md` (or `ARCHITECTURE.md` after WP2) | Module map update; remove M51/M52 pending callouts; theme section fix |
| `AGENTS.md` | Verify/update styling section, navigation history file refs, keybinding helper mention |
| `STATUS.md` | M53, M61 → ✅; update Last updated; verify test count |
| `M53-documentation-sync.md` | Status → ✅ done (via M61) |
| `README.md` | Update DESIGN.md → ARCHITECTURE.md link (if WP2 done) |
| `master-plan/STATUS.md` | Update architecture review link if WP2 done |
| `master-plan/milestones/*.md` | Update any ARCHITECTURE.md references (if WP2 done) |

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | Module map file existence check (script in WP1) | build | WP1 |
| T2 | No `> **M5* pending:**` callouts in DESIGN/ARCHITECTURE.md | grep | WP1 |
| T3 | No stale `ARCHITECTURE.md` references (post-WP2) | grep | WP2 |
| T4 | `make test && make vet` pass | regression | WP3, WP4 |

## Acceptance criteria (milestone done)

- [x] WP1–WP4 complete (WP2 optional — document decision in completion log)
- [x] No `> **M5* pending:**` callouts in ARCHITECTURE.md
- [x] Every file in the module map exists (script verification)
- [x] `AGENTS.md` styling section references `m.palette`, not deprecated globals
- [x] `STATUS.md` M53 and M61 both ✅
- [x] `M53-documentation-sync.md` status updated to ✅ done (via M61)
- [x] Test count in STATUS matches reality ±2
- [x] `make test && make vet` pass

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Broken markdown link after rename | low | WP2 verification step catches; `rg` for residual references |
| Module map lists a file that doesn't exist yet (M59/M60 not done) | medium | **M61 must run after M59 and M60** — dependency enforced |
| AGENTS.md verification misses a stale reference | low | WP3 grep verification |

**Rollback:** `git revert` the WP commit. WP2 (rename) can be reverted independently of WP1 (content).

## Handoff notes

**Read first:**
- This milestone file
- The current `ARCHITECTURE.md` module map (to see what's stale)
- `M53-documentation-sync.md` for context on what was already done
- **Verify M59 and M60 are ✅ before starting** — this milestone documents their work

**Do not:**
- Rewrite the `ARCHITECTURE.md` architecture narrative — only fix the module map and remove stale callouts
- Add new sections — this is a sync, not an expansion
- Touch any `.go` files — docs only

**When stuck:**
- If a file in the module map doesn't exist, check whether M59/M60 actually completed. If not, stop M61 and do M59/M60 first.
- If the rename reveals many more references than expected, use `rg 'DESIGN\.md'` to find them all in one pass.

## Estimated total

3 hours (1.5h WP1 + 30m WP2 + 30m WP3 + 30m WP4)

## Priority

🟡 High — closes M53 partial, removes active risk of agents following stale guidance

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | 2026-06-21 |
| Completed | 2026-06-21 |
| Tests added | 0 (docs only) |
| Rename decision | **Yes** — DESIGN.md renamed to ARCHITECTURE.md; all 11 markdown references updated (ARCHITECTURE-REVIEW-2026-06-13.md excluded — point-in-time) |
| Notes | WP1: Updated module map with 32 entries covering all post-M59/M60 files; removed M51/M52 pending callouts; verified all files exist. WP2: `git mv DESIGN.md ARCHITECTURE.md`; bulk-updated references in 11 markdown files via Python script. WP3: Fixed AGENTS.md — removed stale "(pending rename)" callout; dropped phantom "handlers.go before" reference in Navigation History section. WP4: M53 → ✅; M61 → ✅; Partials table updated. 303 tests pass, `make vet` clean. |
