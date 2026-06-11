# M6 — Status Bar + Help Panel

**Status:** ✅ done

## Goal

Add a bottom status bar showing mode, current file path, and contextual shortcuts.
Add a `?` help overlay with full keybinding reference.

## Files to create

- `statusbar.go`

## Steps

### 1. `statusbar.go`
- `RenderStatusBar(mode Mode, filePath string, hints []string, width int) string`:
  - Full-width single line at bottom
  - Background: dark gray (`#1f2937`)
  - Left segment: mode badge (colored, bold)
    - BROWSE → violet
    - VIEW → fuchsia
    - SEARCH → amber
    - FIND → amber
    - HELP → gray
  - Middle segment: current file path (dimmed)
    - Truncated with `".../"` prefix if too long
    - Empty if no file open (browse without selection)
  - Right segment: contextual hints (dimmed)
    - Separated by ` | `
  - Layout via lipgloss: left(20) | middle(fill) | right(40)

### 2. Mode-specific hints
| Mode | Hints |
|------|-------|
| browse | `/ search \| Ctrl+F find \| Enter open \| ? help \| q quit` |
| view | `h back \| j/k scroll \| Tab link \| / search \| ? help \| q quit` |
| search | `type filter \| Enter open \| Esc cancel` |
| find | `type search \| Enter open \| Esc cancel` |
| help | `j/k scroll \| Esc back` |

### 3. Help panel
- Full-screen overlay (replaces main panels, status bar still visible)
- Title: `"obsidian-terminal — Keybindings"` in violet, bold
- Groups and their bindings:

  **Navigation**
  `j / ↑` — move down
  `k / ↓` — move up
  `h / ←` — collapse / back
  `l / →` — expand / forward
  `g / Home` — jump to top
  `G / End` — jump to bottom
  `PgUp / PgDn` — page up / down

  **File Tree**
  `Enter` — open note / toggle folder
  `← →` — collapse / expand folder

  **Viewer**
  `j / k` — scroll down / up
  `g / G` — top / bottom
  `Tab` — cycle wiki-links
  `Enter` — follow selected link
  `h / Esc` — back to browse

  **Search**
  `/` — fuzzy file name search
  `Ctrl+F` — full-text content search
  `Enter` — open selected result
  `Esc` — cancel search

  **Global**
  `?` — toggle this help
  `q` — quit

- Group headers in violet bold
- Key labels in amber
- Descriptions in gray
- `j/k/↑↓` to scroll
- `PgUp/PgDn` for page scroll
- `Esc` / `?` to close help, return to prevMode

### 4. Integration
- Status bar appended at bottom of `View()` output
- Help panel toggled via `?` in browse and view modes
- Help remembers `prevMode` to return correctly
- Scroll offset reset to 0 when help opens

## Test Spec (5 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestStatusBar_ShowsMode` | model_test.go | View() output contains "BROWSE" / "VIEW" / etc. mode badge |
| 2 | `TestStatusBar_ShowsCurrentFile` | model_test.go | File path appears in status bar after opening a note |
| 3 | `TestStatusBar_ShowsHints` | model_test.go | Mode-specific shortcut hints appear in status bar |
| 4 | `TestHelpPanel_ShowsAllSections` | model_test.go | Help View() contains Navigation, File Tree, Viewer, Search, Global |
| 5 | `TestHelpPanel_EscCloses` | model_test.go | Esc exits help mode; mode restored to prevMode |

## Completion Criteria

- [ ] Status bar with colored mode badge, truncated file path, contextual hints
- [ ] Help screen with all keybindings in 5 grouped sections
- [ ] Help scrollable via j/k/↑↓/PgUp/PgDn
- [ ] Help closes via Esc/?, returns to previous mode
- [ ] Status bar always 1 line at bottom
- [ ] All 5 tests pass
- [ ] `go vet ./...` exits 0
