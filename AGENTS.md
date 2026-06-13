# AGENTS.md

Instructions for AI agents working on this repository.

## Project

`obsidian-terminal` — a read-only terminal TUI for browsing Obsidian vaults. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

**Key constraint: read-only.** Never add editing, writing, or vault-modification features.

> **Architecture reference:** See [DESIGN.md](./DESIGN.md) for the full architecture, data flow, module map, state machine, and design decisions. This file covers commands, conventions, and rules only.

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

## Patterns & Conventions

### Styling
- **Only lipgloss.** No raw ANSI strings in render output.
- All colors from `theme.go` constants: `Accent`, `AccentSecondary`, `TextPrimary`, `TextSecondary`, `TextDim`, etc.
- Use the pre-defined styles: `TreeStyle`, `ViewerStyle`, `StatusStyle`

### Keybindings
- Vim-style: `j`/`k` for up/down, `h`/`l` for left/right, `/` for search, `?` for help
- Arrow keys must also work — every navigation binding supports both
- New keybindings must be added to the `KeyMap` struct and `DefaultKeys()`
- Use `MatchKey(msg, key)` for `tea.KeyType` and `MatchRune(msg, rune)` for runes

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

## Master Plan

The project uses a `master-plan/` folder to track progress and plan larger work.

### Structure

```
master-plan/
├── STATUS.md                              # Progress overview (see template/STATUS-TEMPLATE.md)
├── REVIEW-TEMPLATE.md                     # Project-agnostic review template
├── PHASE-12-EXECUTION-PLAN.md             # Work packages, challenged decisions
├── ARCHITECTURE-REVIEW-{date}.md          # Point-in-time review findings
├── template/                              # Templates for new milestones & STATUS
│   ├── README.md
│   ├── MILESTONE-TEMPLATE.md
│   └── STATUS-TEMPLATE.md
└── milestones/                            # Individual milestone documents
    ├── M0-environment.md
    └── ...
```

`STATUS.md` is the source of truth for overall progress — it contains the milestone table with status, test counts, and dates.

**Creating milestones:** copy [template/MILESTONE-TEMPLATE.md](./master-plan/template/MILESTONE-TEMPLATE.md) → `milestones/M{N}-slug.md`. Register in STATUS.md.

**Architecture reviews:** copy [REVIEW-TEMPLATE.md](./master-plan/REVIEW-TEMPLATE.md). Complex work uses work packages (WPs) per [PHASE-12-EXECUTION-PLAN.md](./master-plan/PHASE-12-EXECUTION-PLAN.md).

### When to use milestones

- **Simple tasks** (bug fixes, typos, single-line changes) — do them directly, no milestone needed.
- **Medium tasks** (a small feature, adding a few tests, refactoring one file) — use judgment; a milestone is optional.
- **Complex tasks** (multi-file refactors, new subsystems, significant behavior changes) — **create a new milestone first** in `master-plan/milestones/` before writing code. This ensures the plan is discussed and approved before execution.

### Milestone workflow

1. **Create** — Copy [template/MILESTONE-TEMPLATE.md](./master-plan/template/MILESTONE-TEMPLATE.md) to `master-plan/milestones/M<N>-<slug>.md`. Fill all sections (goal, out of scope, dependencies, WPs, acceptance criteria). Register in `STATUS.md`.
2. **Start** — Change status to `🚧 in progress`, update Started date in `STATUS.md`.
3. **Work** — Execute one **work package (WP)** per session; run `make test && make vet` after each WP.
4. **Complete** — Check all acceptance criteria in the milestone file; set status to `✅ done` (or `🟡 partial → M<N>` if follow-up remains); update Completed date and test count in `STATUS.md`.

See [template/README.md](./master-plan/template/README.md) for status legend and rules.

### Milestone document template

Use [template/MILESTONE-TEMPLATE.md](./master-plan/template/MILESTONE-TEMPLATE.md) — do not use the shortened inline version below.

Legacy minimal shape (deprecated):

```markdown
# M<N> — <Title>
**Status:** ⏳ pending
## Goal
...
```

## Low-Priority Milestones (M85-M99)

Milestones numbered M85 through M99 are low-priority and complex — they require better planning, new dependencies, or significant architectural changes. These must **always be addressed individually** (one at a time) and never bundled into batch execution loops.

## About this file

`AGENTS.md` can be updated as the project evolves. When editing it:
- **Notify the user** before making changes — describe what you intend to change and why.
- Keep instructions actionable and specific.
- Don't duplicate information that belongs in `README.md`, `DESIGN.md`, or `master-plan/STATUS.md`.
