# obsidian-terminal — Architecture & Design

## Overview

A read-only TUI (terminal UI) for browsing Obsidian vaults. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea) (the Elm Architecture for terminals).

**Key constraint: read-only.** The app never writes to the vault. No editing, creating, or deleting files.

---

## Architecture

### Pattern: Elm Architecture (Model → View → Update)

The entire TUI follows Bubble Tea's three-part loop:

```
User Input → Update(msg) → new Model → View() → rendered output
                  ↑                            │
                  └────────────────────────────┘
```

| Role | File | Responsibility |
|------|------|---------------|
| **Model** | `model.go` | Central state struct (`Model`), `Init()`, `Update()`, mode dispatch. |
| **View** | `render_layout.go` | `View()` renders the full screen: tree panel + viewer panel + status bar + toasts. |
| **Update** | `model.go`, `handlers_*.go` | `Update(msg)` dispatches to mode-specific handlers. |

### Mode State Machine

```
ModeBrowse ──Enter──→ ModeView ──Esc──→ ModeBrowse
    │                     │
    ├──/──→ ModeSearch    ├──?──→ ModeHelp
    ├──T──→ ModeTags      ├──s──→ ModeFind
    ├──s──→ ModeFind      ├──b──→ backlink panel (overlay, not a separate mode)
    ├──?──→ ModeHelp      └──/──→ in-note search (overlay, not a separate mode)
    └──P──→ ModeProfilePicker

Overlays (render over main panel without mode change):
  Command palette (Ctrl+K), Recent notes (Ctrl+O), Outline (t), Backlinks (b),
  In-note search (/), Scan errors (from palette), Pins (Ctrl+[/])

Vault states:
  VaultStateOK → VaultStateBroken (inaccessible) → VaultStateOK (auto-recover)
  VaultStateOK → VaultStatePartial (some files failed to scan)
```

| Mode | Handler | Description |
|------|---------|-------------|
| `ModeBrowse` | `handleBrowseKey` | Tree navigation, open notes, search, tags |
| `ModeView` | `handleViewKey` | Scroll, cycle links, find, backlinks, outline |
| `ModeSearch` | `handleSearchKey` | Fuzzy filename search |
| `ModeFind` | `handleFindKey` | Full-text content search within the open note |
| `ModeHelp` | `handleHelpKey` | Keybinding reference panel |
| `ModeTags` | `handleTagsKey` | Tag browsing and file tree filtering |
| `ModeProfilePicker` | `handleProfilePickerKey` | Vault profile selection on startup |

---

## Module Map

### Root package (`main`)

| File | Responsibility | Key Exports |
|------|---------------|-------------|
| `main.go` | Entry point, flag parsing, config loading, Bubble Tea program start | `main()` |
| `model.go` | Central `Model` struct, `Init/Update/View` dispatch, mode constants, global key handling, layout sizing | `Model`, `Mode*` constants |
| `handlers_browse.go` | Browse mode key handler | `handleBrowseKey` |
| `handlers_view.go` | View mode key handler | `handleViewKey` |
| `handlers_search.go` | Secondary mode/overlay handlers: search, find, help, tags, backlinks, command palette, profile picker | `handleSearchKey`, `handleFindKey`, `handleHelpKey`, `handleTagsKey`, `handleBacklinkKey`, `handleCommandPaletteKey`, `handleProfilePickerKey` |
| `handlers_note.go` | Note-loading API + mode transition helpers | `loadNote`, `applyNote`, `openNote`, `enterSearchMode`, `enterFindMode`, `enterHelpMode`, `enterTagsMode`, `noteNavKind` |
| `in_note_search.go` | In-note search overlay | `activateInNoteSearch`, `handleInNoteSearchKey`, `updateInNoteSearch` |
| `history.go` | Navigation history back/forward | `goBackHistory`, `goForwardHistory` |
| `profile_handler.go` | Profile switching + theme application | `switchToProfile`, `setTheme` |
| `vault_rescan.go` | Vault state machine + rescan logic | `checkVaultChanges`, `rescanVault` |
| `pin_handler.go` | Pinned notes subsystem | `togglePin`, `cyclePinnedNext`, `cyclePinnedPrev`, `validatePins` |
| `outline_handler.go` | Outline/TOC builder + renderer | `buildOutline`, `renderOutline` |
| `daily_recent_handler.go` | Daily notes + recent notes overlay | `buildDailyNotePath`, `openDailyNote`, `addRecentNote`, `toggleRecents`, `renderRecents` |
| `render_layout.go` | `View()` + panel renderers | `View`, `renderSearch*`, `renderBrokenVaultScreen`, `renderScanErrors` |
| `preview.go` | Note preview pane | `renderPreview` |
| `textinput.go` | Shared text-input handler | `HandleTextInput` |
| `tree.go` | File tree widget: vault entry nesting, expand/collapse, cursor, filtering | `FileTree`, `NewFileTree` |
| `viewer.go` | Markdown viewer widget: wraps the markdown render pipeline, wiki-link cycling | `MarkdownViewer`, `SetContent` |
| `viewport.go` | Custom viewport: scroll, soft-wrap (ANSI-aware), X/Y offset | `viewport`, `softWrap` |
| `vault.go` | Vault scanning, tree building, note loading, frontmatter parsing | `ScanVault`, `LoadNote`, `VaultEntry` |
| `session.go` | Session state save/restore (tree expansion, cursor position) | `saveSession`, `restoreSession` |
| `config.go` | YAML config loading, defaults, validation with auto-fix, profile and theme parsing | `Config`, `LoadConfig`, `ValidateConfig` |
| `theme.go` | Color palettes, lipgloss styles, style builders | `Palette`, `lookupPalette`, `markdownStyleFrom`, `searchStyleFrom` |
| `keys.go` | Key binding definitions, vim + arrow key dispatch, navigation helpers | `KeyMap`, `DefaultKeys`, `MatchKey`, `MatchRune`, `MatchDown`, `MatchUp`, `MatchLeft`, `MatchRight` |
| `mouse.go` | Mouse event handling: tree click, split drag, scroll, double-click | `handleMouse` |
| `backlinks.go` | Backlinks panel widget (shown inside view mode) | `BacklinkPanel` |
| `tags.go` | Tag browser/filter widget | `TagList` |
| `statusbar.go` | Status bar: mode display, file name, key hints | `renderStatusBar`, `modeHints` |
| `help.go` | Help panel with all keybindings | `renderHelp`, `buildHelpSections` |
| `toast.go` | Toast notification system | `addToast`, `renderToasts`, `expireToasts` |
| `command_palette.go` | Command palette overlay: query input, result navigation, command execution | `openCommandPalette`, `executeCommand` |
| `wikilink.go` | Wiki-link resolution (basename, alias, exact path) | `ResolveWikiLink`, `findAlias`, `findBasename` |
| `yamlmini.go` | Custom mini YAML parser (no external dep) | `scanYAML`, `parseNestedMap`, `parseFlatMap` |
| `profile_picker.go` | Profile selection widget | `ProfilePicker` |

> **Note:** Files under ~250 lines are preferred. `model.go` is the intentional exception (~400 lines) because the `Model` struct, `Init()`, and `Update()` dispatcher are co-located.

### `internal/markdown/`

| File | Responsibility |
|------|---------------|
| `markdown.go` | Full Obsidian-flavored markdown parser and renderer. Parses headings, bold/italic, `[[wikilinks]]`, callouts, tables, checkboxes, code blocks, embeds, blockquotes, lists, horizontal rules, frontmatter, comments. Renders to ANSI-styled terminal output. |

Key functions: `ParseMarkdown`, `RenderMarkdown`, `RendererStyle`, `MarkdownLine`.

### `internal/search/`

| File | Responsibility |
|------|---------------|
| `search.go` | Fuzzy filename search, full-text content search, search state management, result rendering with `HighlightMatches` |

Key functions: `FuzzySearch`, `ContentSearch`, `State`, `RenderResults`.

### `internal/ansiext/`

| File | Responsibility |
|------|---------------|
| `ansiext.go` | Modern ANSI SGR helpers: undercurl (`\033[4:3m`), overline (`\033[53m`) |

---

## Data Flow

### Startup

```
main() → configPathOrDefault() → LoadConfig() → NewModel()
  → ScanVault() → buildTree() → NewFileTree()
  → restoreSession() → expand saved directories, restore cursor
  → tea.NewProgram(model, WithAltScreen, WithMouseCellMotion) → p.Run()
```

### Opening a note

```
Browse mode → Enter on a file → `openNote(path)` in `handlers_note.go` → `loadNote(vaultPath, relativePath)`
  → `LoadNote(vaultPath, relativePath)` in `vault.go`
  → parseFrontmatter() → extract title, tags, aliases
  → `applyNote(note, navUser)` in `handlers_note.go` sets activeNote, mode = ModeView
  → viewer.SetContent(note.Body, width)
    → StripFrontmatter()
    → ParseMarkdown() → []MarkdownLine
    → ResolveEmbeds() → []MarkdownLine
    → RenderMarkdown(lines, width, style) → ANSI string
    → viewport.SetContent(rendered) → softWrap each line
```

### Key dispatch

```
Model.Update(msg tea.KeyMsg)
  → check global keys (Ctrl+C, Ctrl+R, Ctrl+K, Ctrl+O, Ctrl+D, q, Q)
  → check overlays (command palette, recents, outline, scan errors)
  → check broken vault retry ('r')
  → dispatch to mode handler (`handleBrowseKey`, `handleViewKey`, etc. in `handlers_*.go`)
  → mode handler calls tree.Update, viewer methods, etc.
  → returns (Model, Cmd)
```

### Quit flow

```
User presses q → Model.Update → m.quitting = true
  → saveSession(m) → write session.json
  → return m, tea.Quit
```

---

## Rendering Pipeline

### Markdown → String

```
Raw markdown → StripFrontmatter() → stripComments()
  → ParseMarkdown() → []MarkdownLine
  → (optional) ResolveEmbeds() → []MarkdownLine
  → RenderMarkdown(lines, width, style) → ANSI string
```

### Parser phases

1. **Block parsing** — Split on `\n`. Classify each line: heading, code fence, table, list, blockquote, callout, embed, paragraph, empty.
2. **Inline parsing** — For paragraph/heading/blockquote/list lines: `parseInline()` → recursive descent through `parseSegments()` → `mergeSegments()`.

### Inline segment types

| Type | Markdown | Segment Flags |
|------|----------|---------------|
| Bold | `**text**` | `Bold: true` |
| Italic | `*text*` | `Italic: true` |
| Bold+Italic | `***text***` | `Bold: true, Italic: true` |
| Code | `` `text` `` | `Code: true` |
| Strikethrough | `~~text~~` | `Strikethrough: true` |
| Highlight | `==text==` | `Highlight: true` |
| Wiki-link | `[[target\|display]]` | `IsWikiLink: true, WikiTarget, WikiDisplay` |
| Plain text | (anything else) | `Text: "..."` |

### Renderer styles

The `RendererStyle` struct carries all colors used during rendering. It's created from the active `Palette` + `LineSpacing` setting.

---

## Config System

### Config file location

1. `--config` flag path
2. `$XDG_CONFIG_HOME/obsidian-terminal/config.yaml`
3. `~/.config/obsidian-terminal/config.yaml`

### Config options

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `vault_path` | string | — | Path to Obsidian vault (required) |
| `theme` | string | `"dark"` | Color palette name |
| `default_keys` | string | `"vim"` | Key binding style |
| `line_spacing` | string | `"compact"` | `compact`/`normal`/`relaxed` |
| `skip_dirs` | []string | `[".obsidian", ".git", ...]` | Directories to skip |
| `daily_notes_dir` | string | `"Journal"` | Daily notes subdirectory |
| `daily_notes_format` | string | `"2006-01-02"` | Go time format |
| `profiles` | map | — | Named vault profiles |
| `custom_theme` | map | — | Color overrides (hex values) |

Config validation runs at startup (`ValidateConfig`). Invalid values are auto-fixed to defaults with warning toasts:
- Unknown theme → `"dark"` with message listing valid themes
- Invalid `line_spacing` → `"compact"` with valid values
- Invalid `daily_notes_format` → `"2006-01-02"` (round-trip validated)
- Invalid `custom_theme` hex colors → warning for each bad field
- Empty profile paths → warning

### Profile system

Profiles allow switching between vaults/themes:

```yaml
profiles:
  work:
    path: /path/to/work
    theme: dark
  personal:
    path: /path/to/personal
    theme: catppuccin-mocha
```

Activate with `--profile` flag or `P` key in browse mode.

---

## Theme System

### Architecture

```
7 built-in palettes → lookupPalette(name) → Palette stored on Model.palette
                                       ↘ setTheme(model, name) updates Model.palette and derived styles
                                       ↘ markdownStyleFrom(p, lineSpacing) → RendererStyle
                                       ↘ searchStyleFrom(p) → search.Style
                                       ↘ rebuildDerivedStyles(p) → compute composite styles
```

Palette colors are read from `Model.palette` (set by `setTheme` in `profile_handler.go`). Package-level globals in `theme.go` are deprecated and kept only for test defaults.

### Available themes

| Name | Light/Dark |
|------|------------|
| `dark` | Dark (default) |
| `catppuccin-latte` | Light |
| `catppuccin-frappe` | Dark |
| `catppuccin-macchiato` | Dark |
| `catppuccin-mocha` | Dark |
| `dracula` | Dark |
| `alucard` | Dark |

### Custom theme overrides

Any palette color can be overridden via `custom_theme` in config. 15 color keys available, see `config.yaml.example`.

---

## Session System

State saved to `$XDG_STATE_HOME/obsidian-terminal/session.json` (fallback: `~/.local/state/...`).

**Saved on quit:** vault path, expanded directory paths, cursor path.  
**Restored on startup:** expand directories, set cursor. Different vault → skip restore.

---

## Testing Strategy

### Package structure

| Package | Test Type | What it covers |
|---------|-----------|----------------|
| Root (`main`) | Unit + e2e | Config loading, key dispatch, tree operations, viewer rendering, session, search state, model transitions, rescan, mouse, backlinks, tags, embeds, profile picker, pins, outline, recents, checkboxes, command palette, toast, theme |
| `internal/markdown` | Unit | Parser: headings, inline formatting, code blocks, wiki-links, callouts, tables, lists, comments, embeds. Renderer: ANSI output, code block style, horizontal rule, blockquotes, callouts, lists, table scaling, line wrapping, edge cases (nested, double-backtick, backtrack) |
| `internal/search` | Unit | Fuzzy scoring, fuzzy search, content search, rendering, state management |
| `internal/ansiext` | Unit | Undercurl/overline SGR sequences, empty input |

### Patterns

- Go stdlib `testing` — no test frameworks
- Unit tests use direct function calls with table-driven tests where appropriate
- Integration tests use `tea.NewProgram` with simulated input
- Bubble Tea program tests verify mode transitions via `Update()` calls
- Test fixtures in `testdata/test-vault/` (10+ notes with varied content)

### Key test file locations

- `model_test.go` — mode transitions, key dispatch, status bar, help, truncation, resize
- `model_e2e_test.go` — rescan watcher, wiki-link resolution, toasts
- `model_integration_test.go` — end-to-end workflows (render pipeline, search→open, tree→open, wiki-link follow, theme switch, resize, session restore)
- `tree_test.go` — expand/collapse, filtering, empty vault, cursor, ellipsis
- `viewer_test.go` — render pipelines (no broken ANSI, tables), scrolling, links
- `viewport_test.go` — ANSI-aware softWrap, visibleLength
- `vault_test.go` — scanning, note loading, frontmatter, tags, backlinks
- `session_test.go` — save/restore roundtrip, vault mismatch, corrupted file

---

## Key Design Decisions

### Why a custom markdown parser (not glamour/goldmark)?

Obsidian uses non-standard markdown extensions:
- `[[wiki-links]]` with pipe aliases and heading fragments
- `==highlight==`
- `%%comments%%`
- `> [!callout]` syntax
- Checkbox syntax `- [x]`
- Embed syntax `![[embed]]`

No existing Go parser handles all of these correctly. A custom parser gives full control.

### Why no new dependencies?

Single-binary CLI with only Bubble Tea and Lipgloss as dependencies. No external markdown renderer, YAML parser, or search library. This keeps the binary small, startup fast, and eliminates version conflicts.

### Why read-only only?

Editing a vault from a terminal is a different problem space. Read-only focus keeps the app simple, safe, and fast. Users who need editing can use the Obsidian app.

### Why a custom viewport (not Bubbles viewport)?

The Bubbles viewport doesn't support ANSI-aware wrapping, which is essential for rendering lipgloss-styled output. Our viewport also handles horizontal scroll (future), split drag, and custom scroll behavior.

---

## Performance Considerations

- Vault scanning uses `filepath.WalkDir` which is efficient for large vaults
- Search indexes (name, content, backlinks, tags) are built during initial scan
- Fuzzy search limits results to 50 matches
- Content search limits results to 100 matches
- TUI runs at default Bubble Tea FPS (no custom frame rate)
- ANSI-aware softWrap only triggers for lines that exceed visible width
- Session state is small JSON (<1KB typical)
