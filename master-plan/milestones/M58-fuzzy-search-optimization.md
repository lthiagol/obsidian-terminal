# M58 — Fuzzy Search Allocation Optimization

**Status:** ⏸ deferred (Phase 99)  
**Finding:** Review §10 — no original milestone  
**Decision:** D-10 in [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Goal

Reduce allocations in fuzzy filename search for large vaults (10k+ paths) without changing ranking behavior.

## Problem statement

Benchmark: `FuzzySearch` on 10k paths ≈ 1.6ms, **3.2 MB alloc/op**. Typing in browse search (`/`) may allocate heavily on each keystroke.

## Out of scope

- Changing fuzzy scoring algorithm
- External search library dependency
- Content search (`s` mode) — separate code path

## When to execute

Trigger if **any** of:
- M54 WP1 shows search hot in profiles
- User reports input lag in fuzzy search on large vaults
- `/` search benchmark > 5ms p95 on 10k paths

## Work packages (sketch)

### WP1 — Profile current implementation
- Benchmark per keystroke simulation
- Identify alloc sources in `internal/search/search.go`

### WP2 — Reuse lowered paths from vault index
- Store `[]string` paths + `[]string` lowerPaths at scan time
- FuzzySearch reads precomputed slice

### WP3 — Allocation-free scoring path
- Avoid copying path strings in inner loop
- Verify identical results via golden tests

## Acceptance criteria (when activated)

- [ ] ≥ 50% alloc reduction on 10k benchmark
- [ ] All `search_test.go` pass unchanged
- [ ] No new dependencies

## Estimated total

1–2 days

## Priority

🔵 Low — defer until proven need
