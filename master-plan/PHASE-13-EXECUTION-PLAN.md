# Phase 13 — Execution Plan (Plan Remediation)

**Status:** Planning complete — no code executed yet
**Source review:** [ARCHITECTURE-REVIEW-2026-06-13.md](./ARCHITECTURE-REVIEW-2026-06-13.md) + 2026-06-21 plan audit
**Parent plan:** [PHASE-12-EXECUTION-PLAN.md](./PHASE-12-EXECUTION-PLAN.md) (Phase 12)

This document closes the partial milestones left over from Phase 12 (M52, M53) with three focused milestones. Each milestone has expanded work packages in its own file — this plan coordinates ordering and cross-milestone rules.

---

## 1. Why Phase 13 exists

Phase 12 shipped M50–M58 but two milestones were marked ✅ in `STATUS.md` while their milestone files and acceptance criteria were incomplete:

| Milestone | What was marked | What was actually done | What remains |
|-----------|-----------------|------------------------|--------------|
| M52 | ✅ done | WP1–WP5 done (model.go = 400 lines) | WP6: split `handlers.go` by mode — `handlers.go` is still 624 lines / 26 functions |
| M53 | ✅ done | WP1 (KEYBINDINGS), WP3 (README), WP4 (STATUS audit) | WP2: ARCHITECTURE.md module map still has phantom files + M51/M52 pending callouts; AGENTS.md styling section still references deprecated globals |

Phase 13 closes these partials and adds a targeted DRY refactor pass (M60) that was identified during the 2026-06-21 plan audit but not in the original architecture review.

---

## 2. Challenged decisions (decision log)

| ID | Original decision | Challenge | Revised decision |
|----|-------------------|-----------|------------------|
| D-1 | Reopen M52/M53 as 🟡 partial | Loses momentum; just create new milestones | **Reopen as 🟡 partial with explicit follow-up** — preserves the link to original acceptance criteria; new milestones M59/M61 own the remaining work |
| D-2 | M52 WP6 as a single big PR | 624-line file split in one commit is hard to review | **Split into 6 WPs in M59** — extract subsystems first (in_note_search, history, profile_handler), then split mode handlers, then verify |
| D-3 | Do DRY refactors inside M52 WP6 | Mixes file moves with behavior refactors | **Separate milestone M60** — pure file moves (M59) land first, then DRY refactors on the split files (M60). Each has independent rollback. |
| D-4 | Update ARCHITECTURE.md as part of M53 | Module map depends on M59/M60 final file structure | **Defer ARCHITECTURE.md update to M61** — must run after M59/M60 so docs reflect reality, not in-flight refactors |
| D-5 | Keep ARCHITECTURE.md filename | Inconsistent with ARCHITECTURE-REVIEW-*.md naming and template placeholder | **Recommend rename to ARCHITECTURE.md in M61 WP2** — optional, owner's call |
| D-6 | Detail Phase 99 milestones (M97–M99) now | They're low-priority and not near activation | **Keep as placeholders** — detail when reactivated (per 2026-06-21 user decision) |
| D-7 | Activate M57 (package extraction) now | Reactivation criteria not met (no second contributor; model.go navigable post-M52) | **Keep M57 deferred** — flesh out WPs when criteria met |

---

## 3. Refined execution order

```
Phase 13 (strict sequence — do not parallelize)

  M59  WP1–WP6  Finish M52: handlers.go decomposition
       ↓ (blocks M60 — refactors cleaner on split files)
  M60  WP1–WP3  DRY refactors: KeyMap helpers + text-input helper
       ↓ (blocks M61 — docs must reflect final file structure + helpers)
  M61  WP1–WP4  Finish M53: ARCHITECTURE.md module map + AGENTS.md verify + optional rename
```

**Minimum shippable Phase 13:** M59 alone (closes M52 partial, fixes the maintainability blocker). M60 and M61 can follow incrementally.

**Full Phase 13:** M59 + M60 + M61 — closes both partials, removes duplicated patterns, syncs docs to reality. ~8–10 hours total.

---

## 4. Cross-milestone rules (executing agents)

1. **One WP = one commit.** Run `make test && make vet` after each WP.
2. **No scope expansion** — new findings → new milestone or WP appended to STATUS, not bundled silently.
3. **Read-only constraint** — no milestone may write to vault.
4. **Strict sequence** — M59 → M60 → M61. Do not start M60 until M59 is ✅. Do not start M61 until M60 is ✅.
5. **Update STATUS after each milestone closes** — dates, test count delta, status emoji.
6. **M61 must not update ARCHITECTURE.md module map until M59 and M60 are ✅** — otherwise the map will list files that don't exist yet.

---

## 5. Work package index (all Phase 13 milestones)

| Milestone | WP | Title | Est. | Deliverable |
|-----------|-----|-------|------|-------------|
| **M59** | WP1 | Extract `in_note_search.go` | 1h | 5 funcs moved; `in_note_search_test.go` passes |
| | WP2 | Extract `history.go` | 30m | 2 funcs moved; `history_test.go` passes |
| | WP3 | Extract `profile_handler.go` | 30m | 2 funcs moved; profile/theme tests pass |
| | WP4 | Extract `handlers_note.go` | 1h | noteNavKind + loadNote + enter*Mode moved; model_test passes |
| | WP5 | Split remaining `handlers.go` by mode | 2h | 3 new files; `handlers.go` deleted; each < 250 lines |
| | WP6 | Verify + update STATUS | 30m | M52, M38 → ✅; M59 → ✅ |
| **M60** | WP1 | `KeyMap.MatchDown/Up/Left/Right` + replace call sites | 1.5h | 8+ call sites updated; 0 old patterns remain |
| | WP2 | `HandleTextInput` helper + refactor 3 call sites | 2h | search/command-palette/in-note-search refactored |
| | WP3 | Unit tests for new helpers + verify | 1h | ~5 new tests; full suite passes |
| **M61** | WP1 | Fix ARCHITECTURE.md module map | 1.5h | No pending callouts; every listed file exists |
| | WP2 | Rename DESIGN.md → ARCHITECTURE.md (optional) | 30m | `git mv` + reference updates |
| | WP3 | Verify AGENTS.md reflects post-M59/M60 reality | 30m | Styling section correct; file refs updated |
| | WP4 | Final STATUS + milestone audit | 30m | M53, M61 → ✅; test count verified |

**Total estimated effort:** ~11 hours across 13 WPs (3 focused sessions).

---

## 6. Dependency contract

```
M52 (partial) ──M59──→ ✅ (closes M38 too)
                        │
                        ↓
                       M60 ──→ ✅
                                │
                                ↓
                               M61 ──→ ✅ (closes M53)
```

| Milestone | Blocked by | Blocks | Parallel-safe with |
|-----------|------------|--------|---------------------|
| M59 | M51 (✅ done) | M60, M61 | nothing |
| M60 | M59 | M61 | nothing |
| M61 | M59, M60 | nothing (last in Phase 13) | nothing |

---

## 7. Verification matrix (Phase 13 exit)

| Check | Command / criterion |
|-------|---------------------|
| Tests | `make test && make vet` |
| Test count | Update STATUS (expect 298 + M60's ~5 new = ~303) |
| `handlers.go` | Does not exist (deleted in M59 WP5) |
| Handler file sizes | `wc -l handlers_*.go *_handler.go in_note_search.go history.go` — each < 250 |
| `model.go` | `wc -l model.go` < 400 (unchanged from M52) |
| DRY: nav helpers | `rg 'MatchKey\(msg, m\.keys\.(Down\|Up\|Left\|Right)\) \|\| MatchRune'` returns 0 |
| DRY: text input | `rg 'case msg.Type == tea.KeyBackspace' handlers_search.go in_note_search.go` returns 0 |
| ARCHITECTURE.md | No `> **M5* pending:**` callouts; every listed file exists |
| AGENTS.md | Styling section references `m.palette`, not deprecated globals |
| Plan hygiene | Every ✅ in STATUS has ✅ in milestone file (or 🟡 with documented follow-up) |

---

## 8. Related milestone documents

Each milestone file contains **expanded WPs** with step-level verification:

- [M59-finish-handlers-decomposition.md](./milestones/M59-finish-handlers-decomposition.md)
- [M60-dry-refactors.md](./milestones/M60-dry-refactors.md)
- [M61-doc-sync-completion.md](./milestones/M61-doc-sync-completion.md)

**Templates for new work:** [template/MILESTONE-TEMPLATE.md](./template/MILESTONE-TEMPLATE.md)

---

## 9. Handoff for next executing agent

**Start here:** M59 WP1 — read the Import budget table in [M59](./milestones/M59-finish-handlers-decomposition.md), then `grep -n '^func ' handlers.go` to confirm the function locations (line numbers in M59 are accurate as of 2026-06-21 but re-verify before moving).

**Do not start:**
- M60 before M59 is ✅
- M61 before M60 is ✅
- Any WP without first reading the milestone's "Handoff notes" section

**When stuck:**
- If a test fails after a pure file move, you've made a mistake in the move, not the test. Compare function bodies byte-for-byte with the original.
- If `goimports` is not installed: `go install golang.org/x/tools/cmd/goimports@latest`
- If a file in M61's module map doesn't exist, M59 or M60 didn't complete — stop and fix that first.
- Add new findings to a future review or spawn a new milestone — do not expand the current WP.

---

## 10. After Phase 13

Once M59, M60, M61 are all ✅:

1. Run the verification matrix in §7 — all checks pass
2. Update STATUS.md "Last updated" date
3. Consider an architecture review refresh: copy [REVIEW-TEMPLATE.md](./REVIEW-TEMPLATE.md) → `ARCHITECTURE-REVIEW-{date}.md` and re-score the dimensions from the 2026-06-13 review. Expect improvements in Code organization (C+ → B), Maintainability, and Plan hygiene.
4. The next active milestones are Phase 11 features (M48 done; M49 → M96 deferred) and Phase 99 (M96–M99 low-priority). M57 (package extraction) remains deferred until reactivation criteria met.
