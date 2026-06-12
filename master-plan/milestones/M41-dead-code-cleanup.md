# M41 — Dead Code, Cleanup & Papercuts

**Status:** ⏳ pending

## Goal

Remove dead code, unused exports, duplicate logic, and performance papercuts. Add missing godoc comments on exported symbols.

## Issues

### M4: Dead and duplicated code

| Location | Issue |
|----------|-------|
| `model.go:546-553` | `truncateContent()` function is defined but never called in production code |
| `markdown.go` + `vault.go` | Redundant heading parsers (also covered in M40) |
| `tree.go:158-159` | `keys()` method allocates a new `KeyMap` on every tree update |
| `theme.go:381-382` | `IconVertical` and `IconDiamond` defined but never used |

### M6: Missing godoc comments

15+ exported symbols lack godoc comments, including: `Config` struct and all its fields, `KeyMap` struct and all fields, `MatchKey`, `MatchRune`, `VaultIndexes`, `VaultEntry`, `VaultNote`, `Palette` struct, `RendererStyle`, `InlineSegment`, `MarkdownLine`, `WikiLink`, `HeadingInfo`, `Command` struct, `Result`, `State`, `Style`.

### L1: Hardcoded `lipgloss.Color("#000000")` for selection text

Used in 9 places. On terminals with dark text on dark background (or custom themes), this renders invisible text. Should be derived from the palette or computed as the inverse of the accent color.

### L3: `DefaultKeys()` allocates on every tree Update

`FileTree.keys()` → `DefaultKeys()` creates a new `KeyMap` struct with 4 slice allocations per Update call. Cache the result.

### L4: `commandPaletteSearch` resets cursor to 0 on every keystroke

Typing a query in the command palette always resets the cursor to the first result. Users can't type while maintaining their position in the filtered list.

### M7: `IconVertical` and `IconDiamond` unused

Defined in `theme.go:381-382` but never referenced. Remove.

## Files to modify

| File | Changes |
|------|---------|
| `model.go` | Remove `truncateContent` dead code |
| `theme.go` | Remove unused `IconVertical`, `IconDiamond` |
| `tree.go` | Cache `DefaultKeys()` result in FileTree |
| `command_palette.go` | Don't reset cursor to 0 on filter |
| `keys.go` | Add godoc comments |
| `vault.go` | Add godoc comments on exported types |
| `config.go` | Add godoc comments on `Config` fields |
| `internal/markdown/markdown.go` | Add godoc comments on exported types |
| `internal/search/search.go` | Add godoc comments on exported types |
| `theme.go` | Replace `#000000` with palette-derived selection text color |
| `backlinks.go`, `tags.go`, `command_palette.go`, `profile_picker.go`, `model.go`, `statusbar.go`, `search.go` | Replace `lipgloss.Color("#000000")` with computed color |

## Completion Criteria

- [ ] No dead code or unused exports
- [ ] All exported symbols have godoc comments
- [ ] Selection text color is derived from palette (visible on all themes)
- [ ] Tree KeyMap allocation happens once at construction
- [ ] Command palette cursor stays at current position during filtering
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
