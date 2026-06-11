# M15 — Polish & Complete Remaining Gaps

**Status:** 🚧 in progress

## Goal

Address gaps from M11, M13: wire palette to all renderers, add theme tests, fix alias error swallowing, and clean up dead code.

## Files to modify

- `model.go` — set package-level vars from palette in NewModel, theme warning toast
- `theme.go` — remove dead `defaultMarkdownStyle`/`defaultSearchStyle`
- `wikilink.go` — propagate alias extraction error
- `config_test.go` — add theme lookup tests
- `model_test.go` — add palette wiring test
- `viewer_test.go` — update NewViewer calls
- `main_test.go` — update NewViewer calls

## Steps

### 1. Wire palette to all renderers

In `NewModel`, after the palette is resolved, set the package-level color vars so every renderer picks up the active theme automatically — no signature changes needed.

On unknown theme, fall back to `"dark"` and add a warning toast.

### 2. Add theme tests

- `TestThemeLookup_Valid` — all 7 theme names resolve correctly
- `TestThemeLookup_Unknown` — unknown name falls back to dark
- `TestThemeWiredToModel` — palette from Config.Theme flows to Model.palette

### 3. Propagate alias extraction errors

`findAlias` calls `extractAliasesFromFile` — check the error and skip the file on failure instead of silently discarding.

### 4. Cleanup dead code

- Remove `defaultMarkdownStyle()` and `defaultSearchStyle()` from `theme.go`
- Update all `NewViewer(defaultMarkdownStyle())` calls in test files

## Completion Criteria

- [ ] Package-level vars set from palette in NewModel (all renderers use active theme)
- [ ] Unknown theme shows warning toast
- [ ] `TestThemeLookup_Valid` — all 7 themes resolve
- [ ] `TestThemeLookup_Unknown` — falls back to dark
- [ ] `TestThemeWiredToModel` — Config.Theme flows to palette
- [ ] Alias extraction errors properly handled in `findAlias`
- [ ] `defaultMarkdownStyle` and `defaultSearchStyle` removed
- [ ] All test `NewViewer` calls updated
- [ ] `make test` passes all 95 tests
- [ ] `make vet` exits 0
