# M50 — Navigation History Fix

**Status:** ⏳ pending  
**Finding:** F-1 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md)  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) §5

## Goal

Fix back/forward navigation so history stacks are not corrupted, and unify note-loading side effects across all entry points.

## Problem statement

`goBackHistory` / `goForwardHistory` update stacks manually, then call `openNote()`, which pushes the current path onto `history` whenever `activeNote.Path != path`. This re-pushes the note the user is leaving.

## Out of scope

- Changing keybindings (Ctrl+O, `[`, `]`) — behavior fixed, keys unchanged
- Browser-style history across app restarts (session.json)
- Preview pane (M48)
- Refactoring `model.go` file layout (M52)

## Dependencies

- **Blocks:** M53 WP1 (KEYBINDINGS history section should describe fixed behavior)
- **Blocked by:** nothing
- **Parallel-safe with:** nothing (do first)

---

## Design (approved for execution)

### API

```go
type noteNavKind int

const (
    navUser noteNavKind = iota // user navigation: push history, clear forward
    navHistory                 // back/forward: stacks already updated
    navReload                  // rescan refresh: no history changes
)

// loadNote loads path into viewer with full side effects (embed, outline, backlinks, recents).
func (m *Model) loadNote(path string, kind noteNavKind) error

// openNote is the public entry for user-initiated navigation.
func (m *Model) openNote(path string) { _ = m.loadNote(path, navUser) }
```

### History rules (`navUser` only)

1. If `activeNote != nil` and `activeNote.Path != path`: append current path to `history`
2. Set `historyForward = nil`
3. Then load note (mode, viewer, outline, backlinks, recents, embed resolver)

### `navHistory`

- No push to `history`
- No clear of `historyForward` (caller already adjusted stacks)
- Still updates outline, backlinks, recents, viewer

### `navReload`

- Used by `rescanVault` when reloading active note
- No history or recent mutation beyond what `loadNote` always does — **decision:** recents should still update on explicit user open only; for reload, skip `addRecentNote` (add `skipRecents bool` or sub-kind if needed)

**WP1 decision:** `navReload` must **not** call `addRecentNote`.

---

## Work packages

### WP1 — Call-site audit + API skeleton (1h)

**Steps:**
1. List all 11 `openNote(` call sites (see execution plan §5)
2. Classify each as `navUser`, `navHistory`, or `navReload`
3. Add `loadNote` stub and `noteNavKind` in `handlers.go`
4. Change `openNote` body to delegate to `loadNote(path, navUser)`

**Verification:**
- [ ] Table in this file matches grep output
- [ ] `make test && make vet` pass (no behavior change yet if loadNote wraps existing body)

---

### WP2 — Implement loader + fix history handlers (2h)

**Steps:**
1. Move body of current `openNote` into `loadNote`
2. Gate history push on `kind == navUser`
3. Gate `addRecentNote` on `kind == navUser || kind == navHistory` — **not** on `navReload`
4. Update `goBackHistory` / `goForwardHistory` to call `loadNote(..., navHistory)`
5. Update `rescanVault` reload to `loadNote(..., navReload)`

**Verification:**
- [ ] Manual: A → B → C → `[` → `[` lands on A; `]` → `]` lands on C
- [ ] `make test && make vet` pass

---

### WP3 — Daily note missing-file path (30m)

**Steps:**
1. In `openDailyNote`, when `LoadNote` fails, call `loadNote` with synthetic empty note **or** extract shared empty-note setup
2. Ensure embed resolver + backlink panel initialized (match `loadNote` behavior)

**Verification:**
- [ ] New test: daily path missing → ModeView, empty body, no panic, backlinks panel set
- [ ] `make test && make vet` pass

---

### WP4 — History test suite (1h)

**Steps:**
1. Create `history_test.go` (or `handlers_test.go`)
2. Table-driven tests using `NewModel` + `Update(tea.KeyMsg)` where possible

**Required scenarios:**

| # | Scenario | Assert |
|---|----------|--------|
| T1 | Open A, B, C then `[` | activeNote=B, history=[A], forward=[C] |
| T2 | After T1, `[` again | activeNote=A, forward=[B,C] |
| T3 | After T2, `]` | activeNote=B |
| T4 | Open D from tree after back | forward cleared |
| T5 | Ctrl+O in view mode | same as `[` |
| T6 | Rescan with active note | history length unchanged |

**Verification:**
- [ ] All 6 scenarios pass
- [ ] `make test && make vet` pass

---

## Files to modify

| File | Changes |
|------|---------|
| `handlers.go` | `loadNote`, `noteNavKind`, history handlers |
| `model.go` | `rescanVault` → `navReload`; daily note WP3 |
| `history_test.go` | **New** — WP4 |

## Acceptance criteria (milestone done)

- [ ] All WPs verified
- [ ] Call-site table complete; no direct duplicate load logic
- [ ] 6 history tests pass
- [ ] Daily missing-file test passes
- [ ] `make test && make vet` pass
- [ ] STATUS.md: M50 → ✅ with test count delta noted

## Rollback / risk

| Risk | Mitigation |
|------|------------|
| Recents no longer update on rescan | Intended; verify T6 |
| Breaking pinned scroll restore | Run `pinned_test.go` after WP2 |

## Handoff notes

Read execution plan §5 before coding. Write **failing tests in WP4 first** if doing TDD. Do not touch theme globals (M51).

## Estimated total

4–5 hours

## Priority

🔴 Immediate
