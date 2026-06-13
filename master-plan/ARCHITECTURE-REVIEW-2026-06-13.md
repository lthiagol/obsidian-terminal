# Architecture Review — 2026-06-13

Deep-dive review of codebase, design docs, and master-plan. No code changes — findings only.

## Executive Summary

**obsidian-terminal** is in strong shape for a v1 read-only TUI: custom Obsidian markdown parser, solid test suite (~285 tests, ~65% main-package coverage), clear read-only constraint, and thoughtful milestone history. The biggest gaps are **plan/code drift**, **incomplete refactors marked done**, **scalability of vault rescan**, and **missing CI + test helpers**.

| Area | Grade | Notes |
|------|-------|-------|
| Architecture (Elm/Bubble Tea) | B+ | Clear mode machine; `model.go` still a god file |
| Code organization | C+ | `internal/` for markdown/search is good; ~13k LOC still in `package main` |
| Test quality | B | Good unit + integration coverage; gaps on history, perf regressions, CI |
| Performance | B- | Acceptable for small/medium vaults; full rescan + fuzzy alloc heavy at scale |
| Documentation | C | DESIGN/AGENTS strong; KEYBINDINGS and milestone files stale |
| Master-plan hygiene | C+ | STATUS.md useful but out of sync with milestone files and reality |

---

## What Is Done (Verified)

### Core product (Phases 1–6)
- Config, vault scan, file tree, custom markdown parser, fuzzy + content search, status bar, help, watcher/rescan, error handling
- Vault indexes (search, backlinks, tags), backlinks panel, tag browser, pins, outline, daily/recent notes, profiles, custom themes
- Embeds, checkboxes/frontmatter display, tables, command palette, resizable split, mouse support
- Visual polish (M47): Unicode icons, heading hierarchy, in-note search, navigation history keys

### Robustness (Phase 10)
- Config validation with auto-fix toasts (M44)
- Graceful degradation for broken/partial vaults (M45)
- Integration tests in `model_integration_test.go` (M46)

### Note-opening consolidation (partial M38)
`openNote()` in `handlers.go` is now the canonical path for most flows (tree Enter, mouse, pinned, recent, backlinks, wiki-links). Mouse side-effect bug from M36-M41 review appears fixed.

---

## Critical Findings

### 1. Navigation history corrupts the stack (Correctness — High)

`goBackHistory` / `goForwardHistory` manually manipulate stacks, then call `openNote()`, which **always pushes the current note onto `history`** when paths differ.

Example: history `[A, B]`, viewing `C`, user presses `[`:
1. Pop `B`, push `C` to forward — correct so far
2. `openNote(B)` sees `C != B` → pushes `C` back onto history → history becomes `[A, C]` instead of `[A]`

**Impact:** Back/forward navigation (M47) is unreliable. No tests cover this.

**Fix:** Add `openNoteAt(path string, opts NoteOpenOpts)` with `RecordHistory bool`, or internal `navigateToNote` used by history handlers.

**Proposed milestone:** M50

---

### 2. M37 marked done but theme globals remain (Architecture — High)

`Model` has `palette Palette`, but rendering still reads package-level globals (`Accent`, `TreeStyle`, `ModeColors`, etc.) updated by `activatePalette()`. Every `setTheme()` and `NewModel()` writes globals; `View()`, `renderRecents()`, `FileTree.View()`, status bar, etc. read them.

**Impact:**
- Dual source of truth (`m.palette` vs globals)
- Harder to test multiple themes in one process
- M37 completion criteria not met

**Fix:** Move all style reads to `m.palette` (or `ThemeState` on Model). Remove `activatePalette` global mutation; keep globals only as package defaults for tests if needed.

**Proposed milestone:** M51

---

### 3. M38 marked done but decomposition not done (Maintainability — High)

| M38 criterion | Status |
|---------------|--------|
| Single note-open path | ✅ Mostly (`openNote`) |
| `model.go` under 250 lines | ❌ **1013 lines** |
| Extract pin/outline/daily/recent files | ❌ Still in `model.go` |
| Milestone doc status | ❌ Still says `⏳ pending` |

`handlers.go` (639 lines) is also oversized. AGENTS.md rule: split when model grows past ~250 lines of new code.

**Proposed milestone:** M52

---

### 4. Full vault rescan on any change (Performance — High at scale)

`checkVaultChanges()` compares **vault root directory mtime** only. Any change triggers `rescanVault()` → `ScanVault()` → `filepath.WalkDir` + **read every `.md` file** + rebuild all indexes.

Benchmarks (M2 Max, current code):
- `ScanVault` test vault: ~230µs (tiny fixture)
- `FuzzySearch` 10k paths: ~1.6ms, **3.2MB alloc/op**

For vaults with 5k–20k notes, periodic full rescans will cause noticeable UI stalls every second (tick interval).

**Fix options (incremental):**
1. Debounce + coalesce rescans (partially exists: 2s cooldown)
2. Incremental index update for changed paths only
3. Optional `fsnotify` for targeted invalidation
4. Lazy content index: index filenames at scan, load bodies on demand for content search

**Proposed milestone:** M54

---

## Medium Findings

### 5. Documentation drift (Maintainability)

| Document | Issue |
|----------|-------|
| `KEYBINDINGS.md` | Lists M19–M29 keybindings as "⏳ pending"; `/` in view mode still documented as fuzzy search (M47 changed to in-note search); Ctrl+O behavior changed |
| `DESIGN.md` module map | References `outline.go`, `daily.go`, `pins.go`, `recents.go` — **files don't exist**; logic lives in `model.go` |
| `STATUS.md` | Says **144 tests**; actual count is **~285** (`go test -v` PASS lines) |
| Milestone files M37, M38 | Still `⏳ pending` while STATUS marks ✅ done |
| `STATUS.md` non-goals | Lists "graph view" but M49 plans ASCII graph (M47 deferred it) |

**Proposed milestone:** M53

---

### 6. Dead / stub state (Correctness — Medium)

- `previewVisible bool` on Model — field exists, never read or toggled (M48 not started)
- `openDailyNote()` when note file missing: creates empty note inline **without** embed resolver or backlink panel (inconsistent with `openNote`)

---

### 7. Test infrastructure gaps (Quality — Medium)

M46 planned `testutil_test.go` with helpers — **not created**. Integration tests duplicate setup boilerplate.

Missing test coverage for high-risk behavior:
- Navigation history back/forward
- Ctrl+O mode-specific behavior (history vs recents)
- In-note search (`/`, `n`, `N`)
- Daily note when file doesn't exist
- Profile switch end-to-end
- Render-on-resize doesn't corrupt scroll position

No CI (`.github/workflows` absent). No coverage gate. M17 promised `make bench` targets — **not in Makefile**.

**Proposed milestones:** M55 (CI), M56 (test infra)

---

### 8. Package structure (Architecture — Medium, long-term)

Almost all application code is `package main` (~13k LOC). Only `internal/markdown`, `internal/search`, `internal/ansiext` are extracted.

**Consequences:**
- Hard to reuse vault/tree/config without importing main
- Large compile unit; harder for agents to navigate
- Bubble Tea model must stay in main, but vault, tree, config, session could move to `internal/vault`, `internal/tui/tree`, etc.

**Proposed milestone:** M57 (optional, larger refactor)

---

### 9. Rendering pipeline re-parses on resize (Performance — Medium)

`adjustTreeWidth` and window resize call `viewer.SetContent()` → full `ParseMarkdown` + `RenderMarkdown` + `softWrap`. Acceptable for now; cache parsed lines keyed by `(path, body hash, width)` if profiling shows resize jank.

---

### 10. Fuzzy search allocation profile (Performance — Medium)

`FuzzySearch` on 10k paths allocates ~3MB per query. For large vaults with `/` search on every keystroke, consider:
- Reusing lowered-path slice from vault index
- Early termination (already limits to 50 results)
- Subsequence scoring without full path copy

---

## Low Findings

- **v2 goals** in STATUS (mermaid, LaTeX) have no milestones — either add M85+ entries or move to explicit non-goals
- **M34** horizontal scroll deferred correctly; document relationship to wide tables
- **Graph view vs non-goals:** reconcile — ASCII graph is read-only compatible; update non-goals text
- **Go version:** `go.mod` says 1.26.4; STATUS says 1.24+ — align docs
- **Receiver mix:** value receivers on `Update`/`View` with pointer receivers on helpers — idiomatic for Bubble Tea but worth documenting in DESIGN.md

---

## Master-Plan Reorganization Proposal

### Problems today
1. Milestone file status ≠ STATUS.md status (M37, M38)
2. Phase numbering inconsistent (9, 9b, 10, 11)
3. M36-M41-review.md proposed renumbering that wasn't applied — confusing for agents
4. Test counts not updated after M46/M47 additions

### Proposed structure

```
Phase 1–6   Foundation → UX Polish        (unchanged, all done)
Phase 7–10  Bug fixes → Robustness         (unchanged, all done)
Phase 11    Visual UX (M47 done, M48–M49 pending)
Phase 12    Review remediation (M50–M57)   ← NEW
Phase 99    Distribution & long-tail (M97–M99)
```

### Execution order (Phase 12)

1. **M50** — History fix (small, user-visible bug)
2. **M53** — Doc sync (unblocks agents, low risk)
3. **M51** — Theme de-globalization (enables safer UI work)
4. **M52** — model.go split (maintainability)
5. **M56** — Test helpers + coverage for above
6. **M55** — CI pipeline
7. **M54** — Incremental rescan (when vault size becomes an issue)
8. **M48, M49** — Feature work (after foundations)
9. **M57** — Package extraction (optional, when team wants larger refactor)

---

## Positive Patterns to Preserve

1. **Read-only constraint** — consistently enforced; keep in AGENTS.md
2. **Custom markdown parser** — right call for Obsidian flavor; good test coverage in `internal/markdown`
3. **No dependency creep** — only bubbletea + lipgloss direct deps
4. **Integration tests** — `model_integration_test.go` catches real workflow bugs
5. **Config validation with toasts** — good UX for misconfiguration
6. **Milestone template + STATUS.md** — good agent workflow once synced
7. **Benchmarks exist** — extend with Makefile targets and CI smoke bench

---

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| History bug confuses users | High | Medium | M50 |
| Large vault rescan lag | Medium | High | M54 |
| Agent follows stale KEYBINDINGS | High | Medium | M53 |
| Theme globals cause subtle bugs | Low | Medium | M51 |
| model.go grows unbounded | High | Medium | M52 |
| Regressions without CI | Medium | High | M55 |

---

## Handoff for execution

Planning artifacts for Phase 12:

- [PHASE-12-EXECUTION-PLAN.md](./PHASE-12-EXECUTION-PLAN.md) — start here when implementing M50+
- [REVIEW-TEMPLATE.md](./REVIEW-TEMPLATE.md) — use for future reviews

---

## References

- Review performed against branch `first-version`, 2026-06-13
- Tests: `make test && make vet` — all pass
- Coverage: main 65%, markdown 68%, search 85%, ansiext 100%
