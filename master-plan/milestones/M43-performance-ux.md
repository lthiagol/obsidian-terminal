# M43 — Performance & UX Papercuts

**Status:** ⏳ pending

## Goal

Fix performance issues and small UX problems that degrade the user experience.

## Issues

### L3: `DefaultKeys()` allocates on every tree Update

`FileTree.keys()` → `DefaultKeys()` creates a new `KeyMap` struct with 4 slice allocations per Update call. This is called on every keystroke in the tree.

**Fix:** Cache the `KeyMap` in the `FileTree` struct. Initialize once in `NewFileTree`.

### L4: `commandPaletteSearch` resets cursor to 0 on every keystroke

Typing a query in the command palette always resets the cursor to the first result. Users can't type while maintaining their position in the filtered list.

**Location:** `command_palette.go:317-320`

**Fix:** Preserve cursor position when filtering. If the current cursor position is out of bounds after filtering, clamp to the last result.

### L2: Inconsistent pointer vs. value receiver usage

`Model.Update()` and `Model.View()` use value receivers (required by Bubble Tea's `tea.Model` interface). But `checkVaultChanges`, `rescanVault`, `addToast`, `expireToasts`, `togglePin`, `cyclePinnedNext`, `cyclePinnedPrev`, `validatePins`, `buildOutline`, `addRecentNote`, `toggleRecents` all use pointer receivers.

Since `Model` is a large struct (30 fields), passing it by value in `Update` and `View` is expensive. This is unavoidable due to Bubble Tea's interface requirements, but it should be documented.

**Fix:** Add a comment to `Model.Update` and `Model.View` explaining why value receivers are required.

### Performance: Large Model struct passed by value

The `Model` struct has 30 fields and is passed by value in `Update` and `View`. This copies ~1KB per call. While unavoidable due to Bubble Tea's interface, we should be aware of this cost.

**Future consideration:** If performance becomes an issue, consider using a pointer-based model wrapper that implements `tea.Model` and delegates to a pointer-receiver `Model`.

## Files to modify

| File | Changes |
|------|---------|
| `tree.go` | L3: cache `KeyMap` in `FileTree` struct, initialize in `NewFileTree` |
| `command_palette.go` | L4: preserve cursor position when filtering |
| `model.go` | L2: add comments explaining value receiver requirement |

## Completion Criteria

- [ ] `DefaultKeys()` is called once per `FileTree`, not on every Update
- [ ] Command palette cursor stays at current position during filtering
- [ ] Value receiver requirement is documented with comments
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

1 hour
