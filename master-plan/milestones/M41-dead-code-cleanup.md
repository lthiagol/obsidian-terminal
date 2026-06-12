# M41 — Dead Code, Unused Exports & Hardcoded Colors

**Status:** ⏳ pending

## Goal

Remove dead code, unused exports, and fix hardcoded colors that cause visibility issues on certain themes.

## Issues

### M4: Dead and duplicated code

| Location | Issue |
|----------|-------|
| `model.go:546-553` | `truncateContent()` function is defined but never called in production code |
| `theme.go:381-382` | `IconVertical` and `IconDiamond` defined but never used |

### L1: Hardcoded `lipgloss.Color("#000000")` for selection text

Used in 9 places. On terminals with dark text on dark background (or custom themes), this renders invisible text. Should be derived from the palette or computed as the inverse of the accent color.

**Locations:**
- `tree.go:47`
- `backlinks.go:71`
- `tags.go:97`
- `command_palette.go:298`
- `profile_picker.go:93`
- `model.go:700`
- `model.go:853`
- `statusbar.go:13`
- `search.go:373`

**Fix:** Add a `SelectionText` color to the Palette struct (computed as inverse of accent or a fixed contrasting color). Use `palette.SelectionText` instead of hardcoded `#000000`.

## Files to modify

| File | Changes |
|------|---------|
| `model.go` | Remove `truncateContent` dead code |
| `theme.go` | Remove unused `IconVertical`, `IconDiamond`; add `SelectionText` to Palette |
| `backlinks.go`, `tags.go`, `command_palette.go`, `profile_picker.go`, `model.go`, `statusbar.go`, `search.go`, `tree.go` | Replace `lipgloss.Color("#000000")` with `palette.SelectionText` |

## Completion Criteria

- [ ] No dead code or unused exports
- [ ] Selection text color is derived from palette (visible on all themes)
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

2 hours
