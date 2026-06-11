# obsidian-terminal — Build Status

**Last updated:** 2026-06-11
**Language:** Go 1.24+
**Framework:** Bubble Tea + Bubbles + Lipgloss
**Dependencies:** bubbletea, bubbles, lipgloss, gopkg.in/yaml.v3 (4 total)
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
| Frontmatter | gopkg.in/yaml.v3 |
| Vault path | Required (`--vault` flag or config file, no default) |
| Keybindings | Both vim + arrow keys |
| Wiki-links | Tab cycles, Enter follows |
| Symlinks | Shown as-is in tree |
| Skip dirs | .obsidian, .git, .trash, node_modules, archive, dot-prefixed |
| Test framework | Go stdlib `testing` + Bubble Tea program tests |

## Progress

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
| M16: Shed Dependencies | ⏳ pending | 0 | — | — |
| M17: Performance: Hot Path Optimizations | ⏳ pending | 0 | — | — |
| M18: Mouse Support | ⏳ pending | 0 | — | — |
| M19: Backlinks Panel | ⏳ pending | 0 | — | — |
| M20: Tag Browsing & Filtering | ⏳ pending | 0 | — | — |
| M21: Multiple Vault Profiles | ⏳ pending | 0 | — | — |
| M22: Custom Themes | ⏳ pending | 0 | — | — |
| M23: Embedded Block Embeds | ⏳ pending | 0 | — | — |
| M97: Export to PDF/HTML | ⏳ pending | 0 | — | — |
| M98: Image Preview | ⏳ pending | 0 | — | — |
| M99: Homebrew Distribution | ⏳ pending | 0 | — | — |
| **Total** | | **98** | | |
