# M57 — Package Structure Extraction

**Status:** 🎯 design phase (adopt later)  
**Decision:** D-6 in [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Goal

Move domain logic from `package main` into `internal/*` for clearer boundaries and testability.

## Problem statement

~13k LOC in `main`. Vault, config, session, yamlmini are reusable domain code trapped in the binary package.

## Why deferred

1. M52 achieves most maintainability win with lower risk
2. Large import refactor conflicts with active Phase 12 work
3. AGENTS.md: M85–M99 complex milestones one at a time

## Reactivation criteria

Execute M57 when **any** of:
- Second contributor needs to import vault logic
- M52 complete and model.go still hard to navigate
- Preparing public library API (unlikely for this CLI)

## Proposed package map (unchanged)

| Package | Sources |
|---------|---------|
| `internal/vault` | `vault.go`, `wikilink.go` |
| `internal/config` | `config.go` |
| `internal/session` | `session.go` |
| `internal/yamlmini` | `yamlmini.go` |

## Work packages (sketch — do not execute until reactivated)

1. WP1 — Extract `internal/config` + tests
2. WP2 — Extract `internal/vault` + tests
3. WP3 — Extract `internal/session` + `yamlmini`
4. WP4 — Update DESIGN.md; verify no `main` import from internal

## Acceptance criteria (when activated)

- [ ] `go test ./...` pass
- [ ] No circular imports
- [ ] Bubble Tea model remains in `main`

## Architecture notes (design phase)

### Circular dependency: `config.go` → `parseHexColor` → `theme.go`

`config.go` calls `parseHexColor` in `ValidateConfig` (custom theme hex validation). If `config.go` moves to `internal/config`, it cannot import `main`'s `parseHexColor`.

**Resolution path:** Move `parseHexColor` + `paletteFromCustom` into `internal/config` alongside the `CustomTheme` type. These are pure validation/conversion functions — they belong with the config types, not with lipgloss style building in `theme.go`.

### Cross-package type references

| Symbol | Used by | Location after extraction |
|--------|---------|--------------------------|
| `VaultEntry` | tree.go, viewer.go, vault_rescan.go, session.go, model.go | `internal/vault` |
| `Config`, `Profile`, `CustomTheme` | model.go, handlers.go, theme.go, config_test.go | `internal/config` |
| `Session` | model.go | `internal/session` |
| `scanYAML`, `parseNestedMap` | config.go | `internal/yamlmini` |

### Extraction order (safe sequence)

1. **`internal/yamlmini`** — zero deps, can move first
2. **`internal/config`** — depends on yamlmini + needs `parseHexColor` moved too
3. **`internal/vault`** — depends on nothing internal, but heavily referenced from main
4. **`internal/session`** — depends on vault (uses VaultEntry)

### Risk: vault extraction touches ~30 files

`ScanVault`, `LoadNote`, `VaultEntry`, `VaultNote`, `VaultIndexes` are referenced across most of the codebase. This WP alone is 2-3 days of mechanical refactoring plus test updates.

## Estimated total

3–5 days

## Priority

🎯 design phase (adopt later)
