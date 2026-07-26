# Phase 12 — Execution Plan (Ready for Adoption)

**Status:** Planning complete — no code executed yet  
**Source review:** [ARCHITECTURE-REVIEW-2026-06-13.md](./ARCHITECTURE-REVIEW-2026-06-13.md)  
**Review template:** [REVIEW-TEMPLATE.md](./REVIEW-TEMPLATE.md)

This document **challenges** the initial Phase 12 backlog, **refines execution order**, and defines **work packages (WPs)** small enough for one focused agent session each (~2–6 hours).

---

## 1. Challenged decisions (decision log)

| ID | Original decision | Challenge | Revised decision |
|----|-------------------|-----------|------------------|
| D-1 | M50 before everything | Doc sync (M53) first avoids rework | **Keep M50 first** — behavior must be fixed before documenting Ctrl+O/history |
| D-2 | M55 CI after M56 tests | CI without tests still helps; refactors need CI sooner | **Move M55 earlier** — after M50 (M50 includes its own tests); CI protects M51/M52 |
| D-3 | M51 ∥ M52 in parallel | Both touch render paths and theme reads | **Strict sequence: M51 → M52** |
| D-4 | M56 as big “test infra” bucket | Overlaps M50 tests; becomes a dumping ground | **Split:** M50 ships history tests; M56 = testutil + gap tests only |
| D-5 | M54 incremental rescan now | May be premature for typical personal vaults | **Gate M54:** WP1 benchmarks on 1k/5k fixtures; WP2+ only if p95 scan > 200ms |
| D-6 | M57 package extraction in Phase 12 | High merge conflict risk; optional value | **Defer M57 to Phase 99** unless M52 blocked by package boundaries |
| D-7 | Single `openNoteOpts` flag | Too coarse for back vs forward | **Three behaviors:** `recordHistory`, `restoreForward`, used via `loadNote(path, navMode)` enum |
| D-8 | M52 target model.go < 250 lines | Unrealistic in one pass; Update/View alone are ~200+ | **Target < 400 lines** for model.go; handlers split is separate success metric |
| D-9 | Fold daily-note fix into M52 | Small correctness fix | **WP0 in M50** — daily missing-file uses same loader as openNote (30 min) |
| D-10 | Fuzzy search perf (review §10) | No milestone | **M58 deferred** — only if M54 benchmarks show search also hot |

---

## 2. Refined execution order

```
Track A — Correctness & safety (do first)
  M50  WP1–WP4  History + daily-note loader
  M55  WP1–WP2  CI pipeline
  M53  WP1–WP4  Documentation sync (post M50)

Track B — Architecture (sequential)
  M51  WP1–WP5  Theme de-globalization
  M52  WP1–WP7  model.go decomposition (one WP per PR/session)

Track C — Quality & scale (parallel after Track B starts)
  M56  WP1–WP4  testutil + regression gaps
  M54  WP1 only → decision gate → WP2–WP4 if needed

Track D — Features (Phase 11)
  M48  WP1–WP4  Preview pane (~8h; after M50, prefer M51)
  M49  WP1–WP5  ASCII graph (~2–3 days; 60-node global cap)

Deferred
  M57  Package extraction (Phase 99)
  M58  Fuzzy search alloc optimization (Phase 99)
```

**Minimum viable Phase 12:** Track A only (M50 + M55 + M53) — shippable improvement in ~1–2 days.

**Recommended Phase 12:** Track A + Track B — fixes bug, CI, docs, theme, file structure (~1–2 weeks).

---

## 3. Cross-milestone rules (executing agents)

1. **One WP = one commit** (or one PR if using git workflow). Run `make test && make vet` after each WP.
2. **No scope expansion** — new findings → new milestone or WP appended to STATUS, not bundled silently.
3. **Update docs in M53**, not ad-hoc during M51/M52 (except godoc on touched symbols).
4. **Read-only constraint** — no milestone may write to vault.
5. **Partial milestones** — M37/M38 stay partial until M51/M52 acceptance criteria fully checked.

---

## 4. Work package index (all Phase 12 milestones)

| Milestone | WP | Title | Est. | Deliverable |
|-----------|-----|-------|------|-------------|
| **M50** | WP1 | Design `loadNote` API + call-site audit | 1h | Comment table in milestone |
| | WP2 | Implement loader + fix history handlers | 2h | Passing unit tests |
| | WP3 | Fix `openDailyNote` missing-file path | 30m | Test for empty daily note |
| | WP4 | Integration test + help/status hints | 1h | 4+ history scenarios |
| **M55** | WP1 | Makefile bench target + workflow skeleton | 1h | `make bench` works locally |
| | WP2 | GitHub Actions CI (test + vet) | 1h | Green CI on push |
| **M53** | WP1 | KEYBINDINGS.md full pass | 2h | Matches keys.go/help.go |
| | WP2 | ARCHITECTURE.md module map | 1h | No phantom files |
| | WP3 | README + AGENTS alignment | 1h | Theme/history rules accurate |
| | WP4 | STATUS + milestone status audit | 30m | All ✅ verified |
| **M51** | WP1 | Global read audit (grep + checklist) | 1h | List in milestone |
| | WP2 | Thread `Palette` through tree/status/help | 3h | No global reads in widgets |
| | WP3 | Thread through handlers/toast/command palette | 2h | Compile + tests pass |
| | WP4 | Remove runtime `activatePalette` mutation | 1h | Theme switch test |
| | WP5 | Deprecate or shrink global vars | 1h | Document in DESIGN |
| **M52** | WP1 | Extract `vault_rescan.go` | 2h | model.go −80 lines |
| | WP2 | Extract `pin_handler.go` | 2h | No behavior change |
| | WP3 | Extract `outline_handler.go` | 2h | outline tests pass |
| | WP4 | Extract `daily_handler.go` + `recent_handler.go` | 2h | daily/recent tests pass |
| | WP5 | Extract `render_layout.go` (View + panels) | 3h | View tests pass |
| | WP6 | Split `handlers.go` by mode | 4h | handlers under 200 lines each |
| | WP7 | Line-count verification + DESIGN update | 1h | model.go < 400 lines |
| **M56** | WP1 | Create `testutil_test.go` | 2h | Helpers used by 1 test |
| | WP2 | Refactor integration tests to helpers | 2h | Less duplication |
| | WP3 | In-note search test suite | 2h | 5+ cases |
| | WP4 | Profile switch + resize regression tests | 2h | 3+ cases |
| **M54** | WP1 | Large vault benchmarks + fixtures | 3h | Benchmark doc in milestone |
| | WP2 | Incremental index design | 2h | Design section only |
| | WP3 | Implement incremental re-read | 1d | Tests for add/mod/del |
| | WP4 | Wire into `checkVaultChanges` | 4h | No full scan on single edit |
| **M48** | WP1 | PreviewPane + KeyMap `v` | 2h | preview.go |
| | WP2 | Browse View integration | 2h | Right panel preview |
| | WP3 | Cache + 80-line cap | 1.5h | Perf guard |
| | WP4 | Tests + help | 1.5h | 4+ tests |
| **M49** | WP1 | Graph edge builder from searchIndex | 3h | graph_test.go |
| | WP2 | Circle layout + ASCII render | 4h | Bresenham grid |
| | WP3 | ModeGraph + Ctrl+G | 4h | Integration tests |
| | WP4 | Caps + footer polish | 2h | 60 node global max |
| | WP5 | Tests + help | 2h | 8+ tests total |

---

## 5. Call-site contract: note navigation (M50)

All paths must use the loader; **`openNote(path)`** remains the public API for user-initiated navigation.

| Call site | File | `recordHistory` | Notes |
|-----------|------|-----------------|-------|
| Tree Enter | handlers.go | yes | clear forward |
| Wiki-link Enter | handlers.go | yes | |
| Search result Enter | handlers.go | yes | |
| Backlink Enter | handlers.go | yes | |
| Command palette open | command_palette.go | yes | |
| Mouse tree open | mouse.go | yes | |
| Mouse search open | mouse.go | yes | |
| Pinned open | model.go | yes | via openNote |
| Recent open | model.go | yes | via openNote |
| Daily (exists) | model.go | yes | via openNote |
| Daily (missing) | model.go | yes | **after WP3:** use loader |
| History back | handlers.go | **no** | manual stack already updated |
| History forward | handlers.go | **no** | manual stack already updated |
| Rescan reload | model.go | **no** | same path reload; no new navigation |

**Proposed API** (for executing agent — finalize in M50 WP1):

```go
type noteNavKind int

const (
    navUser noteNavKind = iota // push history, clear forward
    navHistory                 // do not push; stacks managed by caller
    navReload                  // rescan/refresh; no history side effects
)

func (m *Model) loadNote(path string, kind noteNavKind) error
func (m *Model) openNote(path string) { _ = m.loadNote(path, navUser) }
```

---

## 6. M51 global read audit (starting checklist)

Executing agent runs `rg '\b(Accent|TreeStyle|ViewerStyle|StatusStyle|HelpStyle|SearchStyle|TextMuted|ModeColors)\b' --glob '*.go'` and checks off each file.

**Must convert before M51 done:**

- `tree.go`, `statusbar.go`, `help.go`, `toast.go`, `tags.go`, `backlinks.go`, `profile_picker.go`, `command_palette.go`, `model.go` (render*), `handlers.go` (renderInNoteSearch)

**May keep using `RendererStyle` from palette (already threaded):**

- `viewer.go`, `internal/markdown`

**AGENTS.md update (in M53):** change “colors from theme.go constants” → “colors from `Model.palette` / passed Palette”.

---

## 7. M52 extraction map (functions → files)

| Target file | Functions moved from model.go |
|-------------|--------------------------------|
| `vault_rescan.go` | `checkVaultChanges`, `rescanVault`, `countFiles`, `vaultStateFrom` |
| `pin_handler.go` | `togglePin`, `openPinnedNote`, `cyclePinnedNext/Prev`, `validatePins` |
| `outline_handler.go` | `buildOutline`, `renderOutline`, `estimateYOffset` |
| `daily_handler.go` | `buildDailyNotePath`, `openDailyNote` |
| `recent_handler.go` | `addRecentNote`, `toggleRecents`, `openRecentNote`, `renderRecents` |
| `render_layout.go` | `View`, `renderSearch*`, `renderBrokenVaultScreen`, `renderScanErrors`, `showScanErrors`, `wordWrap` |
| `handlers_browse.go` etc. | Split handlers.go by `handleBrowseKey`, `handleViewKey`, … |

**Stay in model.go:** `Model` struct, `NewModel`, `Init`, `Update`, `adjustTreeWidth`, `tickCmd`, mode constants.

---

## 8. M54 decision gate (benchmarks)

Generate synthetic vaults in test/bench fixture:

| Fixture | Files | Purpose |
|---------|-------|---------|
| `vault-1k` | 1,000 .md | Baseline |
| `vault-5k` | 5,000 .md | Large personal vault |
| `vault-10k` | 10,000 .md | Stress (optional) |

Measure `ScanVault` p95 over 10 runs. **Proceed to M54 WP2–WP4 only if** p95 > 200ms on vault-5k **or** user reports rescan jank.

Otherwise: document “full scan acceptable” in M54 milestone and close at WP1.

---

## 9. Verification matrix (Phase 12 exit)

| Check | Command / criterion |
|-------|---------------------|
| Tests | `make test && make vet` |
| Test count | Update STATUS (~285 + new) |
| CI | GitHub Actions green |
| History | M50 acceptance criteria all checked |
| Theme | No runtime global mutation; grep `activatePalette` only sets test defaults |
| model.go | `wc -l model.go` < 400 |
| Docs | KEYBINDINGS matches help.go; DESIGN module map accurate |
| Plan | No milestone ✅ with unchecked completion criteria |

---

## 10. Related milestone documents

Each milestone file contains **expanded WPs** with step-level verification:

- [M50-navigation-history-fix.md](./milestones/M50-navigation-history-fix.md)
- [M55-ci-pipeline.md](./milestones/M55-ci-pipeline.md)
- [M53-documentation-sync.md](./milestones/M53-documentation-sync.md)
- [M51-theme-deglobalization.md](./milestones/M51-theme-deglobalization.md)
- [M52-decompose-model.md](./milestones/M52-decompose-model.md)
- [M56-test-infrastructure.md](./milestones/M56-test-infrastructure.md)
- [M54-incremental-vault-rescan.md](./milestones/M54-incremental-vault-rescan.md)
- [M48-preview-pane.md](./milestones/M48-preview-pane.md)
- [M49-graph-view.md](./milestones/M49-graph-view.md)
- [M58-fuzzy-search-optimization.md](./milestones/M58-fuzzy-search-optimization.md) (deferred)

**Templates for new work:** [template/MILESTONE-TEMPLATE.md](./template/MILESTONE-TEMPLATE.md)

---

## 11. Handoff for next executing agent

**Start here:** M50 WP1 — read call-site table (Section 5), write failing history tests (WP4 can be red-first), then WP2 implementation.

**Do not start:** M51 before M50 merged; M52 before M51 done; M54 WP2 before WP1 benchmarks.

**When stuck:** Add finding to ARCHITECTURE-REVIEW or spawn new milestone — do not expand current WP.
