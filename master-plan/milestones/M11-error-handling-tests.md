# M11 — Error Handling & Tests

**Status:** ✅ done

## Goal

Fix silent error swallowing across the codebase and add missing test coverage for error paths, renderers, and utility functions.

## Files to modify

- `model.go` — add toasts for silent LoadNote errors, rescanVault fixes
- `vault.go` — surface scanErrors, fix empty-filename panic, propagate parse errors
- `markdown.go` — fix Segments[0] bounds, propagate alias errors
- `viewer.go` — clamp negative dimensions
- Test files — add missing test coverage

## Steps

### 1. Silent LoadNote errors → toasts (model.go:236, 299, 352, 387)

When Enter fails to load a note in browse, view, search, or find mode,
add `m.addToast("Could not load note: "+err.Error(), ToastError)` before returning.

### 2. Surface scanErrors (vault.go:96-99, model.go:691)

`ScanVault` returns `scanErrors []string` but caller discards it.
Store in `Model.scanErrors`, display count in status bar when > 0.

### 3. Fix double LoadNote in rescanVault (model.go:707-714)

Load once, check for nil, reuse result — avoid loading the file twice.

### 4. checkVaultChanges silent stat error (model.go:672)

If `os.Stat` fails on vault root, add a toast warning.

### 5. Empty filename panic (vault.go:202)

`strings.ToUpper(name[:1])` panics if filename is `.md`.
Check `len(name) > 0` before slicing.

### 6. renderCallout Segments[0] bounds check (markdown.go:754)

Guard `line.Segments[0].Text` with a len check.

### 7. Negative dimension clamping (model.go:166-167, viewer.go:88-89)

Clamp `treeWidth`, viewport width/height to minimums to prevent panic on narrow terminals.

### 8. Silent YAML parse failure in frontmatter (vault.go:235-237)

`parseFrontmatter` returns empty data on YAML error. Instead, return the error;
caller can fall back to filename-based title.

### 9. Silenced alias extraction errors (markdown.go:550)

`findAlias` calls `extractAliasesFromFile` and discards the error.
Return the error and handle at callsite.

### 10. Propagate markdown parse errors

`ParseMarkdown` currently ignores parse errors with `_`. Log them at debug level
or return an `[]error` alongside the lines.

### 11. Add missing tests

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestCheckVaultChanges_StatError` | model_test.go | os.Stat failure on vault root is handled |
| 2 | `TestRescanVault_ScanError` | model_test.go | ScanVault error produces toast |
| 3 | `TestLoadNote_EmptyFilename` | vault_test.go | `.md` file doesn't panic |
| 4 | `TestRenderCallout_EmptySegments` | markdown_test.go | Callout with no segments doesn't panic |
| 5 | `TestSetSize_NegativeDimensions` | viewer_test.go | Width 0 doesn't panic |
| 6 | `TestParseFrontmatter_InvalidYAML` | vault_test.go | Bad YAML returns error, not silently empty |
| 7 | `TestFindAlias_FileReadError` | markdown_test.go | Alias extraction on unreadable file |
| 8 | `TestWrapText` | markdown_test.go | Word wrapping at various widths |
| 9 | `TestTruncatePath` | model_test.go | Path truncation with long/noisy paths |
| 10 | `TestTruncateContent` | model_test.go | Content truncation at line limit |
| 11 | `TestRenderBlockquote` | markdown_test.go | Blockquote rendering |
| 12 | `TestRenderCallout` | markdown_test.go | Callout with all types |
| 13 | `TestRenderList` | markdown_test.go | Ordered/unordered list rendering |
| 14 | `TestRenderCodeBlock` | markdown_test.go | Fenced code block rendering |
| 15 | `TestRenderHorizontalRule` | markdown_test.go | HR rendering |
| 16 | `TestFuzzyScore_EmptyInput` | search_test.go | Empty query and target edge cases |

## Completion Criteria

- [ ] All silent error paths produce toasts or status bar warnings
- [ ] scanErrors displayed in status bar when > 0
- [ ] No more double-loading in rescanVault
- [ ] Empty filename no longer panics
- [ ] Negative dimensions clamped, no panics
- [ ] YAML parse errors propagated from frontmatter
- [ ] 16 new tests added and passing
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
