# M53 — Documentation & Plan Sync

**Status:** 🟡 partial → **M61**  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Completion summary (2026-06-13)

| WP | Status | Notes |
|----|--------|-------|
| WP1 — KEYBINDINGS.md | ✅ done | Matches `keys.go`, `help.go`, `model.go` global dispatch |
| WP2 — DESIGN.md module map | ❌ **incomplete** | Module map still references phantom files (`outline.go`, `daily.go`, `pins.go`, `recents.go`); still has "M51 pending" and "M52 pending" callouts; missing 5 files extracted in M52 (`vault_rescan.go`, `pin_handler.go`, `outline_handler.go`, `daily_recent_handler.go`, `render_layout.go`). Finish in **M61**. |
| WP3 — README + AGENTS.md | 🟡 partial | README updated. AGENTS.md styling section still references deprecated globals (`Accent`, `TreeStyle`, etc.) — M51 made these deprecated. Finish in **M61**. |
| WP4 — STATUS + milestone audit | ✅ done | Test count matches reality ±5 |

## Goal

Eliminate plan/code drift so agents and contributors do not reintroduce fixed bugs or wrong keybindings.

## Out of scope

- User-facing marketing copy
- Rewriting DESIGN.md architecture narrative (only module map + testing section updates)
- Translating docs

## Dependencies

- **After:** M50 (history behavior final), M55 (document CI in README)
- **Before:** M48/M49 feature work (agents need accurate KEYBINDINGS)
- **Follow-up:** M61 (finishes WP2 + AGENTS.md styling section; optional rename DESIGN.md → ARCHITECTURE.md)

---

## Work packages

### WP1 — KEYBINDINGS.md (2h)

**Checklist — every row must match `keys.go`, `help.go`, `model.go` global dispatch:**

Browse mode:
- [ ] `/` fuzzy search, `s` content search
- [ ] `T` tags, `P` profiles, `p` pin, `Ctrl+K` palette, `Ctrl+O` recents overlay
- [ ] `Ctrl+D` daily, `Ctrl+←/→` resize split

View mode:
- [ ] `/` **in-note search** (not fuzzy — M47)
- [ ] `n` / `N` cycle in-note matches when search active
- [ ] `[` back, `]` forward, **Ctrl+O back** (not recents)
- [ ] `b` backlinks, `t` outline, `Tab` wiki-links

Global:
- [ ] `Ctrl+R` rescan, `q`/`Q` quit browse/view only

**Steps:**
1. Move all implemented keys from "Planned" → "Current"
2. Delete or mark "backup" keys that were never implemented (`o` for outline)
3. Add "Mode-specific" column where same key differs by mode

**Verification:**
- [ ] Side-by-side diff against `help.go` `buildHelpSections` — no conflicts

---

### WP2 — DESIGN.md module map (1h)

**Replace phantom entries:**

| Remove / wrong | Actual location (until M52) |
|----------------|----------------------------|
| `outline.go` | `model.go` → M52 moves to `outline_handler.go` |
| `recents.go` | `model.go` |
| `pins.go`, `pinned_display.go` | `model.go` |
| `daily.go` | `model.go` |

**Add:**
- `handlers.go`, `command_palette.go`, `mouse.go`, `model_integration_test.go`
- Note M51 pending: theme globals vs `m.palette`
- Note `loadNote` / `noteNavKind` after M50

**Verification:**
- [ ] Every file in module map exists (`test -f` or glob)

---

### WP3 — README + AGENTS.md (1h)

README:
- [ ] Features list: in-note search, navigation history
- [ ] Quick reference: `/` in view = in-note search; Ctrl+O context-dependent
- [ ] CI badge (optional after M55)
- [ ] `make bench` in make targets table

AGENTS.md:
- [ ] Styling section: "read colors from `Model.palette`" (post-M51) or note "transition: globals until M51"
- [ ] Master plan section: link `REVIEW-TEMPLATE.md`, `PHASE-12-EXECUTION-PLAN.md`
- [ ] Testing: mention `testutil_test.go` when M56 lands

**Verification:**
- [ ] No README claim without code backing

---

### WP4 — STATUS + milestone audit (30m)

**Steps:**
1. Re-count tests: `go test ./... -v -count=1 | grep -c '^--- PASS'`
2. For each Phase 12 milestone marked ✅ in future: ensure completion criteria checked
3. Verify M37/M38 remain partial until M51/M52 done
4. Add link to PHASE-12-EXECUTION-PLAN in STATUS header

**Verification:**
- [ ] Test count matches reality ±5
- [ ] Spot-check 5 random ✅ milestones: criteria match code

---

## Acceptance criteria (milestone done)

- [x] WP1 complete
- [ ] WP2 complete (→ **M61**)
- [ ] WP3 AGENTS.md styling section (→ **M61**)
- [x] WP4 complete
- [ ] Agent can implement M48 using only KEYBINDINGS + DESIGN without reading all handlers (blocked on WP2)
- [x] `make test && make vet` pass (docs only — no code required)

## Handoff notes

Do not change keybindings in this milestone — docs only. If doc reveals a code bug, file new finding / M50 follow-up.

WP2 + AGENTS.md styling section are now owned by **M61**. M61 also considers renaming `DESIGN.md` → `ARCHITECTURE.md` to align with `ARCHITECTURE-REVIEW-*.md` naming and the master-plan template's `{ARCHITECTURE or DESIGN doc}` placeholder.

## Estimated total

4–5 hours (WP1, WP3-README, WP4 done in ~2h; WP2 + AGENTS.md ~2h in M61)

## Priority

🟡 High (Track A, after M50)

## Completion log

| Field | Value |
|-------|-------|
| Started | 2026-06-13 |
| Completed (WP1, WP3-README, WP4) | 2026-06-13 |
| Tests added | 0 (docs only) |
| Notes | WP2 + AGENTS.md styling section → M61. Optional DESIGN.md → ARCHITECTURE.md rename also in M61. |
