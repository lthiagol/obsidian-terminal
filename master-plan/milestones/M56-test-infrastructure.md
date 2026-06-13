# M56 — Test Infrastructure & Targeted Coverage

**Status:** ⏳ pending  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Goal

Reduce test boilerplate; add regression tests for gaps **not** covered by M50/M46 — without chasing coverage percentage.

## Problem statement

M46 planned `testutil_test.go` — never created. Integration tests duplicate 15+ lines of setup. In-note search, profile switch, and resize lack dedicated tests.

## Out of scope

- M50 history tests (owned by M50 WP4)
- Fuzzing / property-based testing
- Coverage gate in CI
- Rewriting all tests to use helpers in one PR

## Dependencies

- **After:** M50 WP4 (history tests establish patterns)
- **Parallel with:** M52 after WP1 (helpers stabilize integration tests)

---

## Work packages

### WP1 — `testutil_test.go` helpers (2h)

**Implement:**

```go
func testVaultPath(t *testing.T) string
func newTestModel(t *testing.T, cfg *Config) Model
func sendKey(t *testing.T, m *Model, msg tea.KeyMsg) Model
func sendKeys(t *testing.T, m *Model, keys ...tea.KeyMsg) Model
func assertMode(t *testing.T, m Model, want Mode)
func assertActiveNotePath(t *testing.T, m Model, suffix string)
func navigateToFirstFile(t *testing.T, m *Model) Model // move from integration test
```

**Verification:**
- [ ] Helpers compile; used by at least one existing test

---

### WP2 — Refactor integration tests (2h)

**Steps:**
1. Replace duplicated setup in `model_integration_test.go` with helpers
2. Do not change assertions — refactor only

**Verification:**
- [ ] Same test count; all pass
- [ ] LOC in integration file reduced ~30%

---

### WP3 — In-note search tests (2h)

**New file:** `in_note_search_test.go`

| Test | Steps |
|------|-------|
| Activate | View mode → `/` → `inNoteSearchActive` |
| Type query | Runes → matches populated |
| Cycle n/N | Multiple matches → scroll position changes |
| Esc dismiss | Clears state |
| Empty note | `/` no-op |

**Verification:**
- [ ] 5 tests pass

---

### WP4 — Profile + resize regression (2h)

| Test | Purpose |
|------|---------|
| Profile switch | Theme + vault path change; mode browse |
| Resize in view | Ctrl+→; activeNote preserved; viewer renders |
| Theme + palette | After M51: assert `m.palette` not global |

**Verification:**
- [ ] 3 tests pass
- [ ] `make test && make vet` pass

---

## Optional WP5 — Makefile test-cover (30m)

```makefile
test-cover:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out | tail -1
```

Not required for milestone completion.

---

## Acceptance criteria

- [ ] WP1–WP4 complete
- [ ] `testutil_test.go` exists and is documented in DESIGN.md testing section (M53 follow-up ok)
- [ ] No test that only asserts `!= nil` without behavior check
- [ ] `make test && make vet` pass

## Handoff notes

M50 owns history. Do not duplicate T1–T6 from M50 here.

## Estimated total

1–2 days

## Priority

🟡 High (Track C, after M50)
