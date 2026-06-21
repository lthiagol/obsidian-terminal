# M58 — Fuzzy Search Allocation Optimization

**Status:** ✅ done (activated out of order — was ⏸ deferred, executed 2026-06-13)  
**Finding:** Review §10 — no original milestone  
**Decision:** D-10 in [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Completion summary (2026-06-13)

| WP | Status | Notes |
|----|--------|-------|
| WP1 — Profile current implementation | ✅ done | Confirmed 3.2MB alloc/op baseline on 10k paths |
| WP2 — Reuse lowered paths from vault index | ✅ done | `FuzzySearch` signature now takes `pathsLower`, `pathsRunes`, `pathsLowerRunes` (pre-computed at scan time) |
| WP3 — Allocation-free scoring path | ✅ done | `pathRune` helper reuses pre-computed `[]rune` slices; inner loop allocates nothing |

## Benchmark results (2026-06-21 verification)

```
BenchmarkFuzzySearch-4   138   863846 ns/op   2712 B/op   2 allocs/op
```

| Metric | Before (review) | After | Change |
|--------|-----------------|-------|--------|
| Latency | 1.6 ms/op | 0.86 ms/op | **-46%** |
| Allocs/op | ~3.2 MB | 2.7 KB | **-99.9%** |
| Alloc count | (not recorded) | 2/op | — |

Far exceeds the 50% alloc reduction acceptance criterion.

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

## Work packages (executed)

### WP1 — Profile current implementation ✅
- Benchmark per keystroke simulation
- Identified alloc sources in `internal/search/search.go`

### WP2 — Reuse lowered paths from vault index ✅
- `vault.go` / `vault_rescan.go` store `[]string` paths + `[]string` lowerPaths + `[][]rune` pathsRunes + `[][]rune` pathsLowerRunes at scan time
- `FuzzySearch` reads precomputed slices

### WP3 — Allocation-free scoring path ✅
- `pathRune` helper avoids copying path strings in inner loop
- `FuzzyScore` takes `[]rune` directly instead of converting per call
- Identical results via golden tests in `search_test.go`

## Acceptance criteria

- [x] ≥ 50% alloc reduction on 10k benchmark (achieved: 99.9%)
- [x] All `search_test.go` pass unchanged (23 tests pass)
- [x] No new dependencies

## Estimated total

1–2 days (actual: ~1 day)

## Priority

🔵 Low — executed early because the fix was small and the win was large

## Completion log

| Field | Value |
|-------|-------|
| Started | 2026-06-13 |
| Completed | 2026-06-13 |
| Tests added | 0 (behavior unchanged — golden tests cover) |
| Notes | Milestone was originally ⏸ deferred; activated because the optimization was small (signature change + pre-computation) and the benchmark win was large. |
