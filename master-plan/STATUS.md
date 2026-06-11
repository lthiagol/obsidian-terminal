# obsidian-terminal — Build Status

**Last updated:** 2026-06-11
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
| M17: Performance (Profile-Driven) | ⏳ pending | 0 | — | — |
| M18: Mouse Support | ⏳ pending | 0 | — | — |

### Phase 2: Core Indexes

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M18.5: Vault Index System | ⏳ pending | 0 | — | — |

### Phase 3: Navigation Features

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M19: Backlinks Panel | ⏳ pending | 0 | — | — |
| M20: Tag Browsing & Filtering | ⏳ pending | 0 | — | — |
| M24: Pinned Notes | ⏳ pending | 0 | — | — |
| M25: Outline / Table of Contents | ⏳ pending | 0 | — | — |
| M26: Daily Notes + Recent Notes | ⏳ pending | 0 | — | — |

### Phase 4: Vault Management

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M21: Multiple Vault Profiles | ⏳ pending | 0 | — | — |
| M22: Custom Themes | ⏳ pending | 0 | — | — |

### Phase 5: Markdown Features

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M23: Embedded Block Embeds | ⏳ pending | 0 | — | — |
| M27: Checkboxes + Frontmatter Display | ⏳ pending | 0 | — | — |
| M28: Markdown Tables | ⏳ pending | 0 | — | — |

### Phase 6: UX Polish

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M29: Command Palette | ⏳ pending | 0 | — | — |

### Phase 7: Future (Low Priority)

| Milestone | Status | Tests | Started | Completed |
|-----------|--------|-------|---------|-----------|
| M97: Export to PDF/HTML | ⏳ pending | 0 | — | — |
| M98: Image Preview | ⏳ pending | 0 | — | — |
| M99: Homebrew Distribution | ⏳ pending | 0 | — | — |

**Total Tests:** 98

## Execution Order

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

**Rationale:** References many features, better to implement after they exist.

### Batch 7: Future (individual, low priority)
- **M97** — Export to PDF/HTML
- **M98** — Image Preview
- **M99** — Homebrew Distribution

**Rationale:** Low priority, tackle individually when needed.

## Milestone Dependencies

```
M16a (Viewport) → M17 (Performance)
M16b (YAML) → M27 (Frontmatter Display)
M18.5 (Vault Index System)
  ├── M19 (Backlinks) — uses backlink index
  └── M20 (Tags) — uses tag index
```

## Keybinding Conflicts Resolved

All keybindings are documented in [KEYBINDINGS.md](../KEYBINDINGS.md). Key resolutions:

- `t` — Outline (View mode only)
- `T` — Tag browser (Browse mode only)
- `p` — Pin note (Browse/View modes)
- `P` — Profile switcher (Browse mode only)
- `b` — Backlinks (View mode only)

No conflicts: same key can have different actions in different modes.
