# AGENTS.md

Instructions for AI agents working on this repository.

## Project

`obsidian-terminal` — a read-only terminal TUI for browsing Obsidian vaults. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

**Key constraint: read-only.** Never add editing, writing, or vault-modification features.

## Commands

```bash
make build       # compile binary
make run         # go run .
make test        # run all tests
make test-race   # run tests with race detector
make vet         # go vet
make lint        # golangci-lint (requires: brew install golangci-lint)
make fmt         # gofmt
make clean       # remove built binary
make install     # go install .
```

Always run `make test && make vet` after making changes.

## Architecture

| File | Purpose |
|------|---------|
| `main.go` | Entry point, flag parsing, config loading |
| `model.go` | TUI state machine — `Model` struct, `Init`, `Update`, `View`, all key handlers, toast management, status bar, help panel |
| `vault.go` | Vault scanning (`ScanVault`), tree building, note loading, frontmatter parsing |
| `tree.go` | File tree widget — expand/collapse, cursor navigation, rendering |
| `viewer.go` | Markdown viewer widget — viewport-based, wiki-link cycling |
| `markdown.go` | Custom Obsidian-flavored Markdown parser and renderer |
| `search.go` | Fuzzy filename search (`SearchName`) and full-text content search (`SearchContent`) |
| `keys.go` | Keymap definitions — `KeyMap` struct, `DefaultKeys()`, `MatchKey`, `MatchRune` |
| `config.go` | YAML config — `Config` struct, `LoadConfig`, `DefaultConfig` |
| `theme.go` | Lipgloss styles and ANSI color palette |
| `config_test.go` | Config loading tests |
| `keys_test.go` | Key dispatch tests |
| `model_test.go` | Model state tests (vault path errors, quit) |
| `model_e2e_test.go` | Bubble Tea program tests (mode transitions, tree interaction) |
| `vault_test.go` | Vault scanning and note loading tests |
| `tree_test.go` | File tree behavior tests |
| `viewer_test.go` | Viewer rendering and wiki-link tests |
| `markdown_test.go` | Parser and renderer tests |
| `search_test.go` | Fuzzy and content search tests |
| `testdata/test-vault/` | Test fixture vault with notes, callouts, frontmatter, symlinks |

## Patterns & Conventions

### Bubble Tea Architecture
- The TUI follows the **Elm Architecture**: `Model` → `Init()` → `Update(msg)` → `View()`
- `Model` uses **value receivers** (not pointers) — mutations return a new `Model`
- `tea.KeyMsg` dispatch is the primary input; modes route to `handleBrowseKey`, `handleViewKey`, etc.
- Modes: `ModeBrowse` → `ModeView` → `ModeSearch` → `ModeFind` → `ModeHelp`

### Styling
- **Only lipgloss.** No raw ANSI strings in render output.
- All colors from `theme.go` constants: `Accent`, `AccentSecondary`, `TextPrimary`, `TextSecondary`, `TextDim`, etc.
- Use the pre-defined styles: `TreeStyle`, `ViewerStyle`, `StatusStyle`

### Keybindings
- Vim-style: `j`/`k` for up/down, `h`/`l` for left/right, `/` for search, `?` for help
- Arrow keys must also work — every navigation binding supports both
- New keybindings must be added to the `KeyMap` struct and `DefaultKeys()`
- Use `MatchKey(msg, key)` for `tea.KeyType` and `MatchRune(msg, rune)` for runes

### Config
- YAML format, loaded from `~/.config/obsidian-terminal/config.yaml`
- All configurable behavior must have a `Config` field with a sensible default
- Don't break existing config file compatibility

### Testing
- Go stdlib `testing` package only — no test frameworks
- Bubble Tea program tests for integration: `tea.NewProgram(model).Run()` with simulated input
- Use `t.TempDir()` for vault fixtures, `os.WriteFile` for test data
- Test vault is at `testdata/test-vault/` — add fixtures there for parser/renderer tests

## Rules

- **No new dependencies** without strong justification. This is a single-binary CLI.
- **No external markdown renderers** (glamour, goldmark, etc.). The custom parser in `markdown.go` handles all Obsidian flavor.
- **Don't break existing keybindings.** Any new key must not conflict with vim or arrow navigation.
- **Don't write to the vault.** This is a viewer, not an editor. No file creation, modification, or deletion.
- **go vet and tests must pass** before considering work done.
- **Follow Go conventions** — `gofmt` formatting, exported symbols have godoc comments, errors are handled (never `_`).
- **Keep `model.go` under control.** If it grows past ~250 lines of new code, split into a new file (`handlers.go`, `toast.go`, etc.). See M9 in master-plan.
