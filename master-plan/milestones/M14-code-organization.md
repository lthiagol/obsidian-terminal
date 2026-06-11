# M14 — Code Organization & Package Structure

**Status:** ⏳ pending

## Goal

Extract `internal/markdown/` and `internal/search/` packages, move wiki-link resolution out of the parser, and define clean package boundaries.

## Files to modify

- `markdown.go` → split: wiki-link resolution to new `wikilink.go`, parser + renderer to `internal/markdown/`
- `search.go` → `internal/search/search.go`
- `model.go` — add imports
- `viewer.go` — add imports
- `main_test.go` — update imports
- `markdown_test.go` → `internal/markdown/markdown_test.go`
- `search_test.go` → `internal/search/search_test.go`

## Steps

### 1. Move wiki-link resolution out of `markdown.go`

`ResolveWikiLink`, `findAlias`, `findExactPath`, `findBasename`, `extractAliasesFromFile` → new `wikilink.go`.
This breaks markdown.go's dependency on `VaultEntry` and `parseFrontmatter`, making the package extraction clean.

### 2. Create `internal/markdown/` package

Move parser + renderer. Define a `RendererStyle` struct to hold colors (injected from `main`).
Export: `ParseMarkdown`, `ExtractWikiLinks`, `RenderMarkdown`, `StripFrontmatter`, all types and block constants.
Remove dead `Styles` type.

### 3. Create `internal/search/` package

Move all search logic. Accept colors as parameters to `RenderSearchResults`.
Export: `NewSearchState`, `SearchState`, `SearchResult`, `SearchMode`, `FuzzySearch`, `ContentSearch`, `RenderSearchResults`, etc.

### 4. Update `main` package imports

Wire `model.go`, `viewer.go`, handlers to import the new packages.
Update test imports in `main_test.go`, `viewer_test.go`, `model_e2e_test.go`.

### 5. Verify

## Completion Criteria

- [ ] `wikilink.go` exists with all wiki-link resolution logic
- [ ] `internal/markdown/` package exists with parser + renderer
- [ ] `internal/search/` package exists
- [ ] `main` package imports and uses both packages correctly
- [ ] All 79 existing tests pass with updated imports
- [ ] `make vet` exits 0
