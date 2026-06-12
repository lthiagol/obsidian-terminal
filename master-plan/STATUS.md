# obsidian-terminal — Build Status

**Last updated:** 2026-06-12
**Language:** Go 1.24+
**Framework:** Bubble Tea + Bubbles + Lipgloss
**Dependencies:** bubbletea, lipgloss (2 total)
**Target:** Read-only TUI markdown viewer for Obsidian vaults

## Goals

- **v1:** Read-only markdown viewer with file tree, fuzzy + content search, wiki-link navigation,
         custom Obsidian-flavored markdown parser, auto-refresh
- **v2:** Tabs, backlinks, tag browser, outline/TOC, daily-note navigation, in-note search,
         tables, checkboxes, mermaid code-blocks, LaTeX
- **Non-goals:** AI features, editing/writing operations, kanban, pomodoro, graph view, bookmarks

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
| M37: Theme System Refactor | ✅ done | — | 2026-06-11 | 2026-06-11 |

### Phase 8: Architecture Improvements (Priority: 🟡 High)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M38: Split model.go + Consolidate Note Opening | ✅ done | — | 2026-06-11 | 2026-06-11 |
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
| M48: Note Preview Pane | ⏳ pending | 0 | — | — |
| M49: Graph View (ASCII) | ⏳ pending | 0 | — | — |

### Phase 10: Robustness (Priority: 🔵 Future)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M44: Config Validation | ✅ done | 19 | 2026-06-12 | 2026-06-12 |
| M45: Graceful Degradation | ✅ done | 9 | 2026-06-12 | 2026-06-12 |
| M46: Integration Test Suite | ✅ done | 7 | 2026-06-12 | 2026-06-12 |

### Phase 11: Future (Low Priority)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M97: Export to PDF/HTML | ⏳ pending | 0 | — | — |
| M98: Image Preview | ⏳ pending | 0 | — | — |
| M99: Homebrew Distribution | ⏳ pending | 0 | — | — |

**Total Tests:** 144

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
M44-M46 (Robustness) — future
```

## Keybinding Conflicts Resolved

All keybindings are documented in [KEYBINDINGS.md](../KEYBINDINGS.md). Key resolutions:

- `t` — Outline (View mode only)
- `T` — Tag browser (Browse mode only)
- `p` — Pin note (Browse/View modes)
- `P` — Profile switcher (Browse mode only)
- `b` — Backlinks (View mode only)

No conflicts: same key can have different actions in different modes.
