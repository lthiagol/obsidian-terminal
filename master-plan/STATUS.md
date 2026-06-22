# obsidian-terminal — Build Status

**Last updated:** 2026-06-21 (M96: Graph View ✅; M99: code-complete 🟡 — PAT + first release pending)
**Language:** Go 1.26+ (see `go.mod`)
**Architecture review:** [ARCHITECTURE-REVIEW-2026-06-13.md](./ARCHITECTURE-REVIEW-2026-06-13.md)  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](./PHASE-12-EXECUTION-PLAN.md) (Phase 12), [PHASE-13-EXECUTION-PLAN.md](./PHASE-13-EXECUTION-PLAN.md) (Phase 13)  
**Templates:** [template/README.md](./template/README.md)  
**Framework:** Bubble Tea + Bubbles + Lipgloss
**Dependencies:** bubbletea, lipgloss (2 total)
**Target:** Read-only TUI markdown viewer for Obsidian vaults

## Goals

- **v1:** Read-only markdown viewer with file tree, fuzzy + content search, wiki-link navigation,
         custom Obsidian-flavored markdown parser, auto-refresh
- **v2:** Tabs, backlinks, tag browser, outline/TOC, daily-note navigation, in-note search,
         tables, checkboxes, mermaid code-blocks, LaTeX
- **Non-goals:** AI features, editing/writing operations, kanban, pomodoro, bookmarks
- **Planned read-only extras:** ASCII graph view (M49), note preview pane (M48) — see Phase 11

## Key Decisions

| Decision | Choice |
|----------|--------|
| Language | Go (single binary, minimal deps) |
| TUI framework | Bubble Tea |
| Components | Bubbles (tree, viewport, textinput) |
| Styling | Lipgloss |
| Markdown | **Custom parser** — no glamour, full Obsidian flavor |
| Config format | YAML (`~/.config/obsidian-terminal/config.yaml`) |
| Frontmatter | Custom mini YAML parser (no external dep) |
| Vault path | Required (`--vault` flag or config file, no default) |
| Keybindings | Both vim + arrow keys |
| Wiki-links | Tab cycles, Enter follows |
| Symlinks | Shown as-is in tree |
| Skip dirs | .obsidian, .git, .trash, node_modules, archive, dot-prefixed |
| Test framework | Go stdlib `testing` + Bubble Tea program tests |
| Keybinding reference | [KEYBINDINGS.md](../KEYBINDINGS.md) |

## Progress

### Phase 1: Foundation

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M0: Environment + Test Infra | ✅ done | 0 | 2026-06-10 | 2026-06-10 |
| M1: Config + Vault Scanner | ✅ done | 8 | 2026-06-10 | 2026-06-10 |
| M2: Basic TUI Shell | ✅ done | 7 | 2026-06-10 | 2026-06-10 |
| M3: File Tree Navigator | ✅ done | 6 | 2026-06-10 | 2026-06-10 |
| M4: Custom Markdown Parser | ✅ done | 10 | 2026-06-10 | 2026-06-10 |
| M5: Search (Fuzzy + Content) | ✅ done | 10 | 2026-06-10 | 2026-06-10 |
| M6: Status Bar + Help | ✅ done | 5 | 2026-06-10 | 2026-06-10 |
| M7: File Watcher + Polish | ✅ done | 9 | 2026-06-10 | 2026-06-10 |
| M8: Error Handling + Edge Cases | ✅ done | 8 | 2026-06-10 | 2026-06-10 |
| M9: Code Quality & Structure | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M10: Deduplication & DRY | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M11: Error Handling & Tests | ✅ done | 16 | 2026-06-11 | 2026-06-11 |
| M12: Performance | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M13: Theme System & Color Palettes | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M14: Code Organization & Package Structure | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M15: Polish & Complete Remaining Gaps | ✅ done | 3 | 2026-06-11 | 2026-06-11 |
| M16a: Replace Viewport Dependency | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M16b: Replace YAML Dependency | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M17: Performance (Profile-Driven) | ✅ done | 0 | 2026-06-11 | 2026-06-11 |
| M18: Mouse Support | ✅ done | 4 | 2026-06-11 | 2026-06-11 |

### Phase 2: Core Indexes

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M18.5: Vault Index System | ✅ done | 5 | 2026-06-11 | 2026-06-11 |

### Phase 3: Navigation Features

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M19: Backlinks Panel | ✅ done | 4 | 2026-06-11 | 2026-06-11 |
| M20: Tag Browsing & Filtering | ✅ done | 4 | 2026-06-11 | 2026-06-11 |
| M24: Pinned Notes | ✅ done | 6 | 2026-06-11 | 2026-06-11 |
| M25: Outline / Table of Contents | ✅ done | 6 | 2026-06-11 | 2026-06-11 |
| M26: Daily Notes + Recent Notes | ✅ done | 8 | 2026-06-11 | 2026-06-11 |

### Phase 4: Vault Management

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M21: Multiple Vault Profiles | ✅ done | 7 | 2026-06-11 | 2026-06-11 |
| M22: Custom Themes | ✅ done | 10 | 2026-06-11 | 2026-06-11 |

### Phase 5: Markdown Features

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M23: Embedded Block Embeds | ✅ done | 8 | 2026-06-11 | 2026-06-11 |
| M27: Checkboxes + Frontmatter Display | ✅ done | 8 | 2026-06-11 | 2026-06-11 |
| M28: Markdown Tables | ✅ done | 6 | 2026-06-11 | 2026-06-11 |

### Phase 6: UX Polish

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M29: Command Palette | ✅ done | 6 | 2026-06-11 | 2026-06-11 |
| M30: Table Rendering Fix | ✅ done | 4 | 2026-06-11 | 2026-06-11 |
| M31: Inline Formatting Parser Fix | ✅ done | 4 | 2026-06-11 | 2026-06-11 |
| M32: Modern Terminal Polish | ✅ done | 4 | 2026-06-11 | 2026-06-11 |
| M33: UX Refinements (Scroll, Spacing, Session) | ✅ done | 3 | 2026-06-11 | 2026-06-11 |
| M34: Horizontal Scroll for Viewer | ⏳ deferred | 0 | — | — |
| M35: Resizable Tree/Viewer Split | ✅ done | 5 | 2026-06-11 | 2026-06-11 |

### Phase 7: Critical Bug Fixes (Priority: 🔴 Immediate)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M36: Quick Bug Fixes | ✅ done | 5 | 2026-06-11 | 2026-06-11 |
| M37: Theme System Refactor | ✅ done (→ M51) | — | 2026-06-11 | 2026-06-13 |

### Phase 8: Architecture Improvements (Priority: 🟡 High)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M38: Split model.go + Consolidate Note Opening | ✅ done (via M52 + M59) | — | 2026-06-11 | 2026-06-21 |
| M39: ANSI Wrapping & Scroll Fixes | ✅ done | — | 2026-06-11 | 2026-06-11 |

### Phase 9: Code Quality (Priority: 🟢 Medium)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M40: Config & Parser Hardening | ✅ done | — | 2026-06-11 | 2026-06-11 |
| M41: Dead Code, Unused Exports & Hardcoded Colors | ✅ done | — | 2026-06-11 | 2026-06-11 |
| M42: Godoc Comments | ✅ done | 0 | 2026-06-11 | 2026-06-12 |
| M43: Performance & UX Papercuts | ✅ done | — | 2026-06-11 | 2026-06-11 |

### Phase 9b: Visual & UX Upgrade (Priority: 🟡 High)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M47: Visual Polish & Look-and-Feel | ✅ done | 0 | 2026-06-12 | 2026-06-12 |
| M48: Note Preview Pane | ✅ done | 4 | 2026-06-13 | 2026-06-13 |
| M49 → **M96** | ⏸ deferred → Phase 99 | 0 | — | — |

### Phase 10: Robustness (Priority: 🔵 Future)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M44: Config Validation | ✅ done | 19 | 2026-06-12 | 2026-06-12 |
| M45: Graceful Degradation | ✅ done | 9 | 2026-06-12 | 2026-06-12 |
| M46: Integration Test Suite | ✅ done | 7 | 2026-06-12 | 2026-06-12 |

### Phase 12: Review Remediation (Priority: 🟡 High)

From [architecture review 2026-06-13](./ARCHITECTURE-REVIEW-2026-06-13.md).

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M50: Navigation History Fix | ✅ done | 7 | 2026-06-13 | 2026-06-13 |
| M51: Theme De-globalization (finish M37) | ✅ done | 0 | 2026-06-13 | 2026-06-13 |
| M52: Decompose model.go (finish M38) | ✅ done (via M59) | 0 | 2026-06-11 | 2026-06-21 |
| M53: Documentation & Plan Sync | ✅ done (via M61) | 0 | 2026-06-13 | 2026-06-21 |
| M54: Incremental Vault Rescan | ✅ done (WP1 gate: 5k scan 74ms) | 0 | 2026-06-13 | 2026-06-13 |
| M55: CI Pipeline | ✅ done | 0 | 2026-06-13 | 2026-06-13 |
| M56: Test Infrastructure & Coverage | ✅ done | 7 | 2026-06-13 | 2026-06-13 |
| M57: Package Structure Extraction | ⏸ deferred (reactivation criteria in milestone; WPs detailed 2026-06-21) | 0 | — | — |
| M58: Fuzzy Search Optimization | ✅ done (activated out of order) | 0 | 2026-06-13 | 2026-06-13 |

**M52 done via M59:** WP1–WP5 done in M52 (extracted `vault_rescan.go`, `pin_handler.go`, `outline_handler.go`, `daily_recent_handler.go`, `render_layout.go`; `model.go` = 400 lines). WP6 (split `handlers.go` by mode) completed in M59 — `handlers.go` (624 lines) decomposed into `handlers_browse.go` (95), `handlers_view.go` (123), `handlers_search.go` (167), `handlers_note.go` (91), `in_note_search.go` (109), `history.go` (25), `profile_handler.go` (37); `handlers.go` deleted. All files < 250 lines. 298 tests pass.

**M53 done via M61:** WP1 (KEYBINDINGS), WP3 (README), WP4 (STATUS audit) done in M53. WP2 (DESIGN.md module map) completed in M61 — module map reflects all post-M59/M60 files, stale M51/M52 pending callouts removed, phantom file references eliminated, AGENTS.md styling section verified, DESIGN.md renamed to ARCHITECTURE.md. All docs in sync with code.

### Phase 13: Plan Remediation (Priority: 🟡 High)

Follow-ups to close partial milestones from Phase 12. See [PHASE-13-EXECUTION-PLAN.md](./PHASE-13-EXECUTION-PLAN.md).

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M59: Finish M52 — handlers.go Decomposition | ✅ done | 0 | 2026-06-21 | 2026-06-21 |
| M60: DRY Refactors — KeyMap + Text-Input Helpers | ✅ done | 5 | 2026-06-21 | 2026-06-21 |
| M61: Finish M53 — Doc Sync Completion (+ rename DESIGN.md → ARCHITECTURE.md) | ✅ done | 0 | 2026-06-21 | 2026-06-21 |

### Phase 99: Future (Low Priority)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M96: Graph View (ASCII, deferred from M49) | ✅ done | 20 | 2026-06-21 | 2026-06-21 |
| M97: Export to PDF/HTML | ⏳ pending (placeholder — detail when reactivated) | 0 | — | — |
| M98: Image Preview | ⏳ pending (placeholder — detail when reactivated) | 0 | — | — |
| M99: Release Automation (Homebrew Formula PR) | ✅ done | 0 | 2026-06-21 | 2026-06-22 |

**Total Tests:** 323 (refresh with `go test ./... -v -count=1 | grep -c '^--- PASS'`)

## Execution Order

... (continues from Batch 9)

### Batch 9b: Visual & UX Upgrade (🟡 High)
34. **M47** — Visual Polish & Look-and-Feel Upgrade (callout icons, heading colors, in-note search, preview pane, history, graph view)

**Rationale:** High-impact UX improvements. Independent — can be done before or after robustness work. Adds in-note search, preview pane, graph view — all read-only compatible features inspired by obsitui's design.

Milestones are organized into execution batches. Within each batch, milestones can be done sequentially or in parallel where noted.

### Batch 1: Foundation (sequential, do first)
1. **M16a** — Replace Viewport Dependency
2. **M16b** — Replace YAML Dependency
3. **M17** — Performance (Profile-Driven)
4. **M18** — Mouse Support

**Rationale:** Each builds on the previous. M16a/b reduce dependencies, M17 profiles the new code, M18 adds input handling.

### Batch 2: Index + Core Navigation (sequential → parallel)
5. **M18.5** — Vault Index System (do first)
6. **M19 + M20** — Backlinks + Tags (can parallelize after M18.5)

**Rationale:** M18.5 creates the index infrastructure that M19/M20 consume. M19 and M20 are independent and can be done in parallel.

### Batch 3: Navigation Features (all independent, can parallelize)
7. **M24** — Pinned Notes
8. **M25** — Outline / Table of Contents
9. **M26** — Daily Notes + Recent Notes

**Rationale:** All three are independent navigation features with similar complexity. Can be done in any order or in parallel.

### Batch 4: Vault Management (can parallelize)
10. **M21 + M22** — Multiple Vault Profiles + Custom Themes

**Rationale:** Both modify config structure with similar patterns. Can be done in parallel or sequentially.

### Batch 5: Markdown Features (sequential)
11. **M27** — Checkboxes + Frontmatter Display (needs M16b)
12. **M28** — Markdown Tables
13. **M23** — Embedded Block Embeds (complex, do last)

**Rationale:** M27 depends on M16b YAML parser. M28 is independent. M23 is complex and standalone.

### Batch 6: UX Polish (individual)
14. **M29** — Command Palette
15. **M30** — Table Rendering Fix
16. **M31** — Inline Formatting Parser Fix
17. **M32** — Modern Terminal Polish
18. **M33** — UX Refinements (Tree ellipsis, Line spacing, Session)
19. **M35** — Resizable Tree/Viewer Split

**Rationale:** M29-M33 are independent UX polish items completed in sequence.
M34 (Horizontal Scroll) is deferred — the milestone document is kept for
future reference but M35 (Resizable Split) is the priority. Giving the
viewer more width naturally solves the same problems (long tables, wide
content) without ANSI-aware horizontal clipping complexity.

### Batch 7: Critical Bug Fixes (Priority: 🔴 Immediate)
20. **M36** — Quick Bug Fixes (C3 quit, C4 palette, C5 SetSize, mouse side effects, nil vault checks)
21. **M37** — Theme System Refactor (C1 globals → model, C6 applyProfile, H4 data-driven themes)

**Rationale:** M36 fixes the most impactful bugs that affect daily use. M37 is a major refactor that must be done before other milestones that depend on stable theme state.

### Batch 8: Architecture Improvements (Priority: 🟡 High)
22. **M38** — Split model.go + Consolidate Note Opening (H1 868 lines, H2 6 duplicates, H3 View complexity)
23. **M39** — ANSI Wrapping & Scroll Fixes (C7 style bleed, L8 multi-width chars, H5 viewport leak)

**Rationale:** M38 reduces model.go complexity and fixes inconsistent note-opening behavior. M39 fixes visual corruption and scroll accuracy.

### Batch 9: Code Quality (Priority: 🟢 Medium)
24. **M40** — Config & Parser Hardening (C9 YAML indent, C10 duplicate parsers, M2 magic numbers)
25. **M41** — Dead Code, Unused Exports & Hardcoded Colors (M4 dead code, L1 hardcoded colors)
26. **M42** — Godoc Comments (M6 missing documentation)
27. **M43** — Performance & UX Papercuts (L3 KeyMap alloc, L4 cursor reset, L2 receiver docs)

**Rationale:** M40-M43 are quality improvements that can be done in any order. They are independent and don't block each other.

### Batch 10: Robustness (Priority: 🔵 Future)
28. **M44** — Config Validation (missing validation, unhelpful errors)
29. **M45** — Graceful Degradation (vault inaccessible, partial scan failures)
30. **M46** — Integration Test Suite (missing end-to-end tests)

**Rationale:** M44-M46 improve robustness and test coverage. They can be done in any order but are lower priority than bug fixes and architecture improvements.

### Batch 11: Review Remediation (🟡 High — see [PHASE-12-EXECUTION-PLAN.md](./PHASE-12-EXECUTION-PLAN.md))

**Track A — Correctness & safety**
1. **M50** — Navigation history fix (4 WPs, includes daily-note loader)
2. **M55** — CI pipeline (after M50 tests land)
3. **M53** — Documentation sync (after M50 behavior frozen)

**Track B — Architecture (strict order)**
4. **M51** — Theme de-globalization (finish M37)
5. **M52** — Decompose model.go (7 WPs, one per session; finish M38)

**Track C — Quality & scale**
6. **M56** — Test helpers + gap tests (after M50)
7. **M54** — WP1 benchmarks only → decision gate → WP2–4 if 5k vault p95 > 200ms

**Track D — Features**
8. **M48** — Preview pane (✅ done); **M49 → M96** — Graph view (⏸ deferred)

**Deferred (Phase 99)**
- **M57** — Package extraction (after M52; reactivation criteria in milestone)
- **M58** — Fuzzy search alloc optimization (if benchmarks prove need)

**Minimum shippable Phase 12:** Track A only (~1–2 days).

### Batch 12: Phase 13 — Plan Remediation (🟡 High — see [PHASE-13-EXECUTION-PLAN.md](./PHASE-13-EXECUTION-PLAN.md))

Closes partial milestones M52 and M53 from Phase 12, plus a targeted DRY refactor pass.

35. **M59** — Finish M52: handlers.go Decomposition (extract `in_note_search.go`, `history.go`, `profile_handler.go`; split `handlers.go` into mode files)
36. **M60** — DRY Refactors: `KeyMap.MatchDown/Up/Left/Right` helpers + shared text-input handler
37. **M61** — Finish M53: Doc Sync Completion (ARCHITECTURE.md module map, AGENTS.md styling section, optional rename to ARCHITECTURE.md)

**Rationale:** Strict sequence M59 → M60 → M61. M59 is pure file moves (low risk, gets `handlers.go` under 250 lines per file). M60 introduces helpers and replaces duplicated patterns within the now-split files. M61 documents the final state — must come last so docs reflect reality, not in-flight refactors.

**Minimum shippable Phase 13:** M59 alone (closes M52 partial, fixes the maintainability blocker). M60 and M61 can follow incrementally.

## Milestone Dependencies

```
M16a (Viewport) → M17 (Performance)
M16b (YAML) → M27 (Frontmatter Display)
M18.5 (Vault Index System)
  ├── M19 (Backlinks) — uses backlink index
  └── M20 (Tags) — uses tag index

M36 (Quick Bug Fixes)
  ↓
M37 (Theme Refactor) — required before M39 (ANSI fixes depend on stable theme)
  ↓
M38 (Split model.go + Consolidate Note Opening) — reduces duplication, makes M39 easier
  ↓
M39 (ANSI Wrapping) — can start after M37
  ↓
M40-M43 (Code Quality) — independent, can be parallel
  ↓
M44-M46 (Robustness)
  ↓
M50 (History fix)
  ↓
M51 (Theme) → M52 (model split) — strict sequence, not parallel
  ↓
M56 (Tests) + M55 (CI)
  ↓
M54 (Incremental rescan) — independent
  ↓
M48-M49 (Features) + M57 (Package extraction, optional)

Phase 13 (plan remediation):
  M52 (partial) → M59 (handlers.go split) → M60 (DRY refactors) → M61 (doc sync, finishes M53)
```

## Partial Milestones (needs follow-up)

| Milestone | Done | Remaining |
|-----------|------|-----------|
| M37 | ✅ done via M51 | No globals remain (deprecated vars kept for tests) |
| M38 | ✅ done via M52 + M59 | No remaining |
| M52 | ✅ done via M59 | No remaining — `handlers.go` split into 7 focused files, all < 250 lines |
| M53 | ✅ done via M61 | No remaining — DESIGN.md → ARCHITECTURE.md rename, module map sync, AGENTS.md verified |

## Keybinding Conflicts Resolved

All keybindings are documented in [KEYBINDINGS.md](../KEYBINDINGS.md). Key resolutions:

- `t` — Outline (View mode only)
- `T` — Tag browser (Browse mode only)
- `p` — Pin note (Browse/View modes)
- `P` — Profile switcher (Browse mode only)
- `b` — Backlinks (View mode only)

No conflicts: same key can have different actions in different modes.

## Maintenance checklist

Run after each milestone closure (see [template/README.md](./template/README.md)):

- [ ] Milestone file: all acceptance criteria checked
- [ ] Milestone status emoji matches this table
- [ ] Test count updated (`go test ./... -v -count=1 | grep -c '^--- PASS'`)
- [ ] Started / Completed dates filled
- [ ] KEYBINDINGS / DESIGN updated if behavior changed
- [ ] No ✅ without checked criteria; use 🟡 partial + follow-up milestone if needed
