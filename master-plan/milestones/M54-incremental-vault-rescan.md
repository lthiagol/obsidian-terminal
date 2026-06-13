# M54 — Incremental Vault Rescan

**Status:** ⏳ pending  
**Finding:** F-4 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md)  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) §8

## Goal

Scale vault refresh for large note collections without full O(n) read on every root mtime change.

## Problem statement

`checkVaultChanges` triggers full `ScanVault` — walk entire tree, read every `.md`, rebuild indexes. Acceptable for test vault (~230µs); risky for 5k+ notes.

## Out of scope

- fsnotify new dependency (unless WP2 design proves insufficient and owner approves)
- Lazy content search index (separate future milestone)
- Changing 1s tick interval (separate UX tweak)

## Dependencies

- **WP2–WP4 gated on WP1 benchmarks**
- **Benefits from** M52 WP1 (`vault_rescan.go` extracted)

---

## Work packages

### WP1 — Benchmarks + fixtures (3h) **REQUIRED**

**Steps:**
1. Add `testdata/bench-vault/` generator script or test helper creating N markdown files
2. Add benchmarks:

```go
func BenchmarkScanVault_1k(b *testing.B)
func BenchmarkScanVault_5k(b *testing.B)
```

3. Record p50/p95 in this milestone's "Benchmark results" section
4. **Decision gate:** if ScanVault_5k p95 < 200ms → mark M54 **done at WP1** with note "full scan acceptable"

**Verification:**
- [ ] Benchmarks run via `make bench`
- [ ] Results table filled below

#### Benchmark results (fill on execution)

| Fixture | p50 | p95 | allocs/op |
|---------|-----|-----|-----------|
| 1k | | | |
| 5k | | | |

---

### WP2 — Incremental design doc (2h) **only if WP1 fails gate**

**Design questions to answer:**
1. Store per-path mtime+size in `VaultIndexes`?
2. On rescan signal: walk tree for structure, re-read only changed md files
3. How to handle deletes (remove from indexes)
4. When to fall back to full scan (profile switch, Ctrl+R, corrupt state)

**Deliverable:** "Incremental scan design" section appended to this file — no code yet.

**Verification:**
- [ ] Design reviewed; edge cases listed (delete, rename, partial read error)

---

### WP3 — Implement incremental re-index (1d)

**Steps:**
1. Add `RescanVaultIncremental` in `vault.go`
2. Unit tests: add one file, modify one, delete one — indexes correct
3. Tree structure matches full scan output

**Verification:**
- [ ] `vault_test.go` incremental tests pass
- [ ] Benchmark shows >50% improvement on 5k partial update

---

### WP4 — Wire into model (4h)

**Steps:**
1. `checkVaultChanges` uses incremental path
2. Ctrl+R and profile switch keep full scan
3. Toast on incremental failure → fallback full scan

**Verification:**
- [ ] Watcher e2e tests pass
- [ ] No duplicate entries in search index

---

## Acceptance criteria

**Minimum (WP1 only):** benchmarks documented; decision recorded  
**Full:** WP1–WP4 if gate failed

- [ ] Decision gate outcome recorded in STATUS
- [ ] `make test && make vet` pass

## Handoff notes

Do not implement WP3 before WP1 numbers exist. Small vault users may never need WP3 — that's OK.

## Estimated total

WP1: 3h | Full: 2–4 days

## Priority

🟢 Medium (Track C, gated)
