# AGENTS.md

Instructions for AI agents working on this repository.

## Project

`obsidian-terminal` — a read-only terminal TUI for browsing Obsidian vaults. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

**Key constraint: read-only.** Never add editing, writing, or vault-modification features.

For architecture, data flow, module map, state machine, and design decisions, see [DESIGN.md](./DESIGN.md) (pending rename to `ARCHITECTURE.md` in M61).

## Commands

```bash
make build       # compile binary
make run         # go run .
make test        # run all tests
make test-race   # run tests with race detector
make bench       # run all benchmarks (5s default)
make bench-short # quick benchmark smoke check (100ms benchtime)
make vet         # go vet
make lint        # golangci-lint (requires: brew install golangci-lint)
make fmt         # gofmt
make clean       # remove built binary
make install     # go install .
```

Always run `make test && make vet` after making changes.

## Patterns & Conventions

### Styling

- **Only lipgloss.** No raw ANSI strings in render output.
- **Colors come from `Model.palette`** (set by `setTheme` in `profile_handler.go`). Read palette fields directly: `m.palette.Accent`, `m.palette.TreeStyle`, etc.
- The package-level color/style vars in `theme.go` (`Accent`, `TreeStyle`, `ViewerStyle`, `StatusStyle`, `ModeColors`, …) are **deprecated** (kept for test defaults only). Do not read them in new code.
- Pass `Palette` to widgets that don't hold a `Model`: `func (t *FileTree) SetPalette(p Palette)`, or `func (t FileTree) View(p Palette) string`.

### Keybindings

- Vim-style: `j`/`k` for up/down, `h`/`l` for left/right, `/` for search, `?` for help
- Arrow keys must also work — every navigation binding supports both
- New keybindings must be added to the `KeyMap` struct and `DefaultKeys()` in `keys.go`
- Use `MatchKey(msg, key)` for `tea.KeyType` and `MatchRune(msg, rune)` for runes
- For the four navigation keys, prefer the helpers (added in M60): `m.keys.MatchDown(msg)`, `MatchUp`, `MatchLeft`, `MatchRight` — instead of manual `MatchKey || MatchRune` combinations
- Don't break existing keybindings. Any new key must not conflict with vim or arrow navigation. Check [KEYBINDINGS.md](./KEYBINDINGS.md) before allocating a key.

### Navigation History

- Use `loadNote(path, kind)` for all navigation moves (lives in `handlers_note.go` after M59; `handlers.go` before)
- `kind` is `noteNavKind`: `navUser` (explicit open), `navHistory` (back/forward), `navReload` (rescan)
- `openNote(path)` is syntactic sugar for `loadNote(path, navUser)` — pushes to history/recents
- `navHistory` — does NOT push to history or recents (prevents double-push)
- `navReload` — does NOT touch history or recents at all
- Always use the appropriate kind to avoid history corruption

### Testing

- Go stdlib `testing` package only — no test frameworks
- Bubble Tea program tests for integration: `tea.NewProgram(model).Run()` with simulated input
- Helpers live in `testutil_test.go`: `newTestModel`, `sendKey`, `sendKeys`, `assertMode`, `assertActiveNotePath`, `navigateToFirstFile` — use these instead of duplicating setup
- Use `t.TempDir()` for vault fixtures, `os.WriteFile` for test data
- Test vault is at `testdata/test-vault/` — add fixtures there for parser/renderer tests

## Rules

- **No new dependencies** without strong justification. This is a single-binary CLI.
- **No external markdown renderers** (glamour, goldmark, etc.). The custom parser in `internal/markdown/markdown.go` handles all Obsidian flavor.
- **Don't write to the vault.** This is a viewer, not an editor. No file creation, modification, or deletion.
- **go vet and tests must pass** before considering work done.
- **Follow Go conventions** — `gofmt` formatting, exported symbols have godoc comments, errors are handled (never `_`).
- **Keep files focused.** Split when a file exceeds ~250 lines or mixes more than one clear responsibility. `model.go` is the exception (target < 400 lines because the `Model` struct + `Update` dispatcher are co-located).

## Master Plan

The project uses `master-plan/` to track progress and plan larger work. See [template/README.md](./master-plan/template/README.md) for the full workflow, status legend, and templates.

**Source of truth:** [STATUS.md](./master-plan/STATUS.md) — milestone table with status, test counts, and dates.

**When to use milestones:**
- **Simple tasks** (bug fixes, typos, single-line changes) — do directly, no milestone.
- **Medium tasks** (small feature, few tests, one-file refactor) — judgment call; milestone optional.
- **Complex tasks** (multi-file refactors, new subsystems, behavior changes) — **create a milestone first** in `master-plan/milestones/` before writing code. Copy [template/MILESTONE-TEMPLATE.md](./master-plan/template/MILESTONE-TEMPLATE.md).

**Milestone lifecycle:** Create → register in STATUS → set 🚧 in progress → execute one WP per session → run `make test && make vet` after each WP → check acceptance criteria → set ✅ done (or 🟡 partial with follow-up milestone) → update STATUS dates + test count.

## Low-Priority Milestones (M96–M99)

Milestones M96 through M99 are low-priority and complex — they require new dependencies, significant architectural changes, or feature work outside the read-only core. Address them **one at a time**, never bundled into batch execution. M57 (package extraction) is separately deferred with its own reactivation criteria.

## About this file

`AGENTS.md` can be updated as the project evolves. Notify the user before changing it. Keep instructions actionable and specific. Don't duplicate information that belongs in [README.md](./README.md), [DESIGN.md](./DESIGN.md), or [master-plan/STATUS.md](./master-plan/STATUS.md).
