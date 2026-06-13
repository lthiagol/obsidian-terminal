# M57 — Package Structure Extraction

**Status:** ⏸ deferred → Phase 99  
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

## Estimated total

3–5 days

## Priority

🔵 Future (Phase 99)
