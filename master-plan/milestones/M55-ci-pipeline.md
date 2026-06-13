# M55 — CI Pipeline

**Status:** ⏳ pending  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Goal

Automate `make test` and `make vet` on every push/PR; add local `make bench` smoke target promised in M17.

## Out of scope

- golangci-lint in CI (optional WP3 — only if linter config exists)
- Coverage thresholds / Codecov
- Release automation (M99)
- Race detector in default CI (too slow; document `make test-race` for local)

## Dependencies

- **After:** M50 WP4 (CI should see history tests)
- **Before:** M51/M52 large refactors (protect merges)

---

## Work packages

### WP1 — Makefile bench target (1h)

**Steps:**
1. Add to `Makefile`:

```makefile
bench:
	$(GO) test -bench=. -benchmem -run=^$$ ./...

bench-short:
	$(GO) test -bench=. -benchtime=100ms -run=^$$ ./...
```

2. Verify `make bench-short` completes < 30s on dev machine
3. Document in README under make targets (full doc sync in M53)

**Verification:**
- [ ] `make bench-short` exits 0
- [ ] No change to default `make test` runtime

---

### WP2 — GitHub Actions workflow (1h)

**Steps:**
1. Create `.github/workflows/ci.yml`:

```yaml
name: CI
on:
  push:
    branches: [main, first-version]
  pull_request:
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: make vet
      - run: make test
```

2. Push branch; confirm workflow green

**Verification:**
- [ ] Workflow runs on push
- [ ] `vet` + `test` both required steps
- [ ] Uses Go version from `go.mod`

---

### WP3 — Optional lint job (defer)

Only if `.golangci.yml` exists or AGENTS.md lint is mandatory:

```yaml
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v6
```

**Default:** skip WP3; note in milestone as optional.

---

## Acceptance criteria

- [ ] WP1 + WP2 complete
- [ ] CI green on default branch
- [ ] README mentions `make bench` (or defer explicitly to M53 with checkbox)

## Handoff notes

Do not add `-race` to CI without owner approval. Keep workflow minimal — one job is enough for this repo size.

## Estimated total

2–3 hours

## Priority

🟡 High (early in Track A)
