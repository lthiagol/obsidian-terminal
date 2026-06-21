# M60 — DRY Refactors: KeyMap + Text-Input Helpers

**Status:** ✅ done  
**Phase:** 13 — Plan Remediation  
**Priority:** 🟢 Medium  
**Finding:** Observation during 2026-06-21 plan review (not in original architecture review)  
**Execution plan:** [PHASE-13-EXECUTION-PLAN.md](../PHASE-13-EXECUTION-PLAN.md)

## Goal

Eliminate two duplicated patterns in the handler layer by adding small helpers, reducing the surface area for copy-paste bugs and making future keybinding changes single-touch.

## Problem statement

**Pattern 1 — Navigation key matching.** The expression `MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune)` appears 8+ times across handler files (browse, view, search, help, tags, backlinks, command palette, profile picker). Same for Up/Left/Right. Adding a new navigation key (e.g., `KeyPgDn` as alias for Down) requires touching every call site.

**Pattern 2 — Text-input handling.** The Esc/Backspace/Runes pattern is duplicated in 3 handlers:
- `handleSearchOrFind` (search/find mode): Esc → exit mode, Backspace → trim query, Runes → append to query
- `handleCommandPaletteKey`: Esc → close palette, Backspace → trim + re-search, Runes → append + re-search
- `handleInNoteSearchKey`: Esc → exit search, Backspace → trim + update, Runes → append + update

Each duplicate is ~10 lines of boilerplate that must be kept in sync.

## Out of scope

- Changing keybinding behavior — pure refactor, identical semantics
- Adding new keybindings
- Refactoring `MatchKey`/`MatchRune` themselves (they're the primitives)
- Touching `mouse.go` (uses a different pattern — mouse events, not key events)
- Splitting `handlers.go` (that's M59 — M60 runs **after** M59)

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | M59 (handlers.go split — refactors are cleaner on split files; avoids double-touching `handlers.go`) |
| **Blocks** | M61 (ARCHITECTURE.md should document the new helpers) |
| **Parallel-safe with** | nothing — touches the same files M59 just split |

## Design (approved for execution)

### Pattern 1 — KeyMap navigation helpers

**API (add to `keys.go`):**

```go
// MatchDown returns true if msg matches the Down key binding (KeyType or Rune).
func (k KeyMap) MatchDown(msg tea.KeyMsg) bool {
    return MatchKey(msg, k.Down) || MatchRune(msg, k.DownRune)
}

// MatchUp returns true if msg matches the Up key binding.
func (k KeyMap) MatchUp(msg tea.KeyMsg) bool {
    return MatchKey(msg, k.Up) || MatchRune(msg, k.UpRune)
}

// MatchLeft returns true if msg matches the Left key binding.
func (k KeyMap) MatchLeft(msg tea.KeyMsg) bool {
    return MatchKey(msg, k.Left) || MatchRune(msg, k.LeftRune)
}

// MatchRight returns true if msg matches the Right key binding.
func (k KeyMap) MatchRight(msg tea.KeyMsg) bool {
    return MatchKey(msg, k.Right) || MatchRune(msg, k.RightRune)
}
```

**Call site transformation:**

```go
// Before
case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):

// After
case m.keys.MatchDown(msg):
```

**Discovery command for executing agent:**
```bash
rg -n 'MatchKey\(msg, m\.keys\.(Down|Up|Left|Right)\) \|\| MatchRune\(msg, m\.keys\.(Down|Up|Left|Right)Rune\)' --glob '*.go'
```
Run this before starting WP1 to get the exact call-site list. Expected: 8+ matches across `handlers_browse.go`, `handlers_view.go`, `handlers_search.go` (post-M59 file names).

### Pattern 2 — Shared text-input handler

**API (add to a new `textinput.go` in `package main`):**

```go
// HandleTextInput processes Esc/Backspace/Runes for an in-TUI text input field.
// Returns:
//   - newQuery: the updated query string (unchanged if key not handled)
//   - dismissed: true if Esc was pressed (caller should exit the input mode)
//   - handled: true if the key was consumed (caller should return early)
//
// The caller is responsible for the "after" behavior (re-search, update display, etc.)
// using the returned newQuery.
func HandleTextInput(msg tea.KeyMsg, query string) (newQuery string, dismissed bool, handled bool) {
    switch msg.Type {
    case tea.KeyEsc:
        return "", true, true
    case tea.KeyBackspace:
        if len(query) > 0 {
            return query[:len(query)-1], false, true
        }
        return query, false, true
    case tea.KeyRunes:
        if len(msg.Runes) > 0 {
            return query + string(msg.Runes), false, true
        }
    }
    return query, false, false
}
```

**Call site transformation (example for `handleSearchOrFind`):**

```go
// Before
case msg.Type == tea.KeyEsc:
    m.mode = m.prevMode
    return m, nil
case msg.Type == tea.KeyBackspace:
    if len(m.searchState.Query()) > 0 {
        m.searchState.SetQuery(m.searchState.Query()[:len(m.searchState.Query())-1])
    }
    return m, nil
case msg.Type == tea.KeyRunes && len(msg.Runes) > 0:
    m.searchState.SetQuery(m.searchState.Query() + string(msg.Runes))
    return m, nil

// After
if newQuery, dismissed, handled := HandleTextInput(msg, m.searchState.Query()); handled {
    if dismissed {
        m.mode = m.prevMode
    } else {
        m.searchState.SetQuery(newQuery)
    }
    return m, nil
}
```

**Call sites to refactor:**
1. `handleSearchOrFind` in `handlers_search.go`
2. `handleCommandPaletteKey` in `handlers_search.go`
3. `handleInNoteSearchKey` in `in_note_search.go` (post-M59 file name)

### Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Helpers as methods on `KeyMap` | Yes | Keeps the API close to the data; `m.keys.MatchDown(msg)` reads naturally |
| `HandleTextInput` as free function in `textinput.go` | Yes | Not tied to `KeyMap` or `Model`; pure input handling |
| Helper returns `(newQuery, dismissed, handled)` instead of mutating `Model` | Yes | Caller owns the "after" behavior (re-search, update display) — abstraction without over-coupling |
| Do not touch `mouse.go` | Yes | Mouse events use `tea.MouseMsg`, not `tea.KeyMsg` — different pattern |

---

## Work packages

### WP1 — Add `KeyMap` navigation helpers + replace call sites (1.5h)

**Steps:**
1. Open `keys.go`. After the existing `MatchKey`/`MatchRune` functions, add the 4 methods: `MatchDown`, `MatchUp`, `MatchLeft`, `MatchRight` (see Design section for exact code).
2. Run `rg -n 'MatchKey\(msg, m\.keys\.(Down|Up|Left|Right)\) \|\| MatchRune' --glob '*.go'` to list all call sites.
3. For each match, replace with the corresponding `m.keys.MatchDown(msg)` / `MatchUp` / `MatchLeft` / `MatchRight` call.
4. Run `goimports -w keys.go` (no new imports needed — `tea` is already imported).
5. Verify `make build` succeeds.

**Verification:**
- [ ] `rg 'MatchKey\(msg, m\.keys\.(Down|Up|Left|Right)\) \|\| MatchRune'` returns 0 matches (old pattern gone)
- [ ] `rg 'm\.keys\.Match(Down|Up|Left|Right)\('` returns the same count as the original 8+ (new pattern present)
- [ ] `keys_test.go` passes (existing KeyMap tests)
- [ ] `make test && make vet` pass

---

### WP2 — Add `HandleTextInput` helper + refactor 3 call sites (2h)

**Steps:**
1. Create `textinput.go` with `package main` header and the `HandleTextInput` function (see Design section for exact code).
2. Add import: `tea "github.com/charmbracelet/bubbletea"`
3. Refactor `handleSearchOrFind` in `handlers_search.go`:
   - Replace the 3 `case` branches (Esc/Backspace/Runes) with a single `if` block calling `HandleTextInput`
   - Keep the Down/Up/Enter cases as-is (they're navigation, not text input)
4. Refactor `handleCommandPaletteKey` in `handlers_search.go`:
   - Replace the 3 `case` branches with `HandleTextInput` call
   - The "after" behavior: `m.commandPaletteQuery = newQuery; m.commandPaletteSearch()` (not dismissed case)
   - Keep Down/Up/Enter cases
5. Refactor `handleInNoteSearchKey` in `in_note_search.go`:
   - Replace the Esc/Backspace/Runes handling with `HandleTextInput` call
   - The "after" behavior: `m.updateInNoteSearch(newQuery)` (not dismissed case)
   - Keep the `n`/`N`/Enter cases (these are in-note-search-specific)
6. Verify each refactor: the new code must produce identical behavior. Run the relevant tests.

**Verification:**
- [ ] `in_note_search_test.go` passes (5 tests — covers Esc, Backspace, Runes, n/N)
- [ ] `command_palette_test.go` passes (6 tests)
- [ ] `model_test.go` search mode tests pass
- [ ] `make test && make vet` pass
- [ ] No `case msg.Type == tea.KeyBackspace:` remains in the 3 refactored handlers (grep to verify)

---

### WP3 — Unit tests for new helpers + final verification (1h)

**Steps:**
1. Add tests to `keys_test.go` for the 4 new methods:
   - `TestKeyMap_MatchDown` — `KeyMsg{Type: KeyDown}` and `KeyMsg{Type: KeyRunes, Runes: ['j']}` both return true; other keys return false
   - Same for `MatchUp`, `MatchLeft`, `MatchRight`
2. Create `textinput_test.go` with `TestHandleTextInput`:
   - Esc → `("", true, true)`
   - Backspace with non-empty query → `(query[:-1], false, true)`
   - Backspace with empty query → `(query, false, true)` (no-op, but handled)
   - Runes → `(query+runes, false, true)`
   - Other key (e.g., `KeyDown`) → `(query, false, false)` (not handled)
3. Run full suite.

**Verification:**
- [ ] New helper tests pass
- [ ] Total test count increased by ~5 (4 KeyMap + 1 HandleTextInput table-driven)
- [ ] `make test && make vet` pass
- [ ] Update `STATUS.md` M60 test count column

---

## Files to modify

| File | Changes |
|------|---------|
| `keys.go` | Add `MatchDown`/`MatchUp`/`MatchLeft`/`MatchRight` methods |
| `textinput.go` | **New** — `HandleTextInput` function |
| `handlers_browse.go` | Replace nav key patterns (post-M59 file) |
| `handlers_view.go` | Replace nav key patterns (post-M59 file) |
| `handlers_search.go` | Replace nav key patterns + refactor 2 text-input call sites (post-M59 file) |
| `in_note_search.go` | Refactor 1 text-input call site (post-M59 file) |
| `keys_test.go` | Add 4 tests for new methods |
| `textinput_test.go` | **New** — `TestHandleTextInput` table-driven |
| `STATUS.md` | M60 → ✅ with test count delta |

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | `keys_test.go` — 4 new method tests | unit | WP3 |
| T2 | `textinput_test.go` — table-driven Esc/Backspace/Runes/other | unit | WP3 |
| T3 | `in_note_search_test.go` passes unchanged | regression | WP2 |
| T4 | `command_palette_test.go` passes unchanged | regression | WP2 |
| T5 | `model_test.go` search mode tests pass unchanged | regression | WP2 |
| T6 | Full suite passes | regression | WP3 |

## Acceptance criteria (milestone done)

- [x] WP1–WP3 complete
- [x] `rg 'MatchKey\(msg, m\.keys\.(Down|Up|Left|Right)\) \|\| MatchRune'` returns 0 matches
- [x] No `case msg.Type == tea.KeyBackspace:` in the 3 refactored handlers
- [x] New helper tests pass (~5 new tests)
- [x] No behavior change — original 298 tests still pass
- [x] `make test && make vet` pass
- [x] `STATUS.md` updated: M60 → ✅ with dates and test count delta

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Subtle behavior change in text-input (e.g., Backspace on empty query) | low | WP3 table-driven test covers edge cases; existing tests are safety net |
| Missed a call site for nav helpers | low | WP1 grep verification (0 old patterns, N new patterns) |
| `HandleTextInput` abstraction leaks (caller needs more info) | low | If WP2 reveals this, abandon the helper for that call site and document why — do not force-fit |

**Rollback:** `git revert` the WP commit. Each WP is independent.

## Handoff notes

**Read first:**
- This milestone file (especially the Design section with exact API)
- Run the discovery `rg` command in WP1 before starting to confirm call-site count
- **M59 must be done first** — file names in this milestone assume M59's split (`handlers_search.go`, `in_note_search.go`). If M59 is not done, stop and do M59 first.

**Do not:**
- Add helpers for `MatchKey(msg, k.Top)` / `MatchBottom` — these use `MatchRune` only (no `KeyType`), so the pattern is already one-line. Only Down/Up/Left/Right have the dual KeyType+Rune pattern.
- Refactor `MatchRune(msg, m.keys.PinRune)` etc. — these are single-rune checks, no duplication.
- Touch `mouse.go` — different event type.

**When stuck:**
- If `HandleTextInput` doesn't fit one of the 3 call sites cleanly (e.g., the caller needs to know if Backspace was a no-op on empty query), document the mismatch and skip that call site. Better to leave one duplicate than force a bad abstraction.
- If a test fails after WP1, you likely changed semantics — compare the truth table of the old expression vs the new method.

## Estimated total

4–5 hours (1.5h WP1 + 2h WP2 + 1h WP3)

## Priority

🟢 Medium — quality improvement, no user-visible change

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | 2026-06-21 |
| Completed | 2026-06-21 |
| Tests added | 5 (4 KeyMap nav helper tests + 1 HandleTextInput table-driven test) |
| Notes | WP1 replaced 22 call sites of `MatchKey(msg, m.keys.*) \|\| MatchRune(msg, m.keys.*Rune)` across `handlers_browse.go`, `handlers_view.go`, `handlers_search.go`, `outline_handler.go`, `daily_recent_handler.go`. WP2 refactored 3 text-input call sites (`handleSearchOrFind`, `handleCommandPaletteKey`, `handleInNoteSearchKey`) to use new `HandleTextInput` helper. All 303 tests pass; `make vet` clean. |
