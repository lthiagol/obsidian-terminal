# M40 — Config & Parser Hardening

**Status:** ⏳ pending

## Goal

Harden the YAML mini-parser against non-standard formatting. Extract hardcoded magic numbers into named constants. Fix the redundant heading-parser duplication between vault.go and markdown.go.

## Issues

### C9: YAML parser indent assumptions (`yamlmini.go:51`)

`scanYAML` checks `itemLine[0]` for indentation on array items but only checks the first byte, not the trimmed line. Tab-indented YAML or non-standard indentation silently breaks array parsing. `parseNestedMap` (line 281) hardcodes `rootIndent + 2` for profile indent depth.

**Fix:** Measure indent consistently using string prefixes rather than byte offsets. Accept tabs as valid indentation. Allow variable indent depth in `parseNestedMap`.

### C10: Redundant heading parsers

`vault.go:400-415` (`isMarkdownHeading`) and `markdown.go:320-334` (`isHeading`) are near-identical. `vault.go:417-426` (`countHeadingLevel`) and `markdown.go:337-346` (`headingLevel`) are identical. The vault package duplicates markdown parsing logic unnecessarily.

**Subtle behavioral difference:** `isMarkdownHeading` (vault.go) returns true for `"# "` (heading with trailing space but no content), while `isHeading` (markdown.go) requires `i+1 < len(line)` — i.e., there must be content after the space.

**Fix:** Remove the vault.go duplicates. Export `IsHeading` and `HeadingLevel` from the markdown package and call them from vault.go.

### M2: Hardcoded magic numbers

| File | Value | Suggested Constant |
|------|-------|---------------------|
| `search.go:215,232` | `50` | `maxSearchResults` |
| `search.go:293` | `100` | `maxContentResults` |
| `search.go:285` | `80` | `contentResultContextLen` |
| `mouse.go:167` | `500*time.Millisecond` | `doubleClickWindow` |
| `mouse.go:61` | `3` (in 4 places) | `mouseScrollStep` |
| `toast.go:32` | `3*time.Second` | `toastTTL` |
| `model.go:50` | `50` | `maxRecentNotes` |
| `markdown.go:689` | `2` | `maxEmbedDepth` |
| `markdown.go:746` | `20` | `minRenderWidth` |
| `model.go:572` | `15` | `minTreeWidth` |

## Files to modify

| File | Changes |
|------|---------|
| `yamlmini.go` | C9: consistent indent handling, variable indent depth |
| `vault.go` | C10: remove duplicate heading parsers; import from markdown package |
| `internal/search/search.go` | M2: extract constants |
| `mouse.go` | M2: extract `doubleClickWindow`, `mouseScrollStep` |
| `toast.go` | M2: extract `toastTTL` |
| `model.go` | M2: extract `maxRecentNotes`, `minTreeWidth` |
| `internal/markdown/markdown.go` | M2: extract `maxEmbedDepth`, `minRenderWidth`; export `IsHeading`/`HeadingLevel` |
| `*_test.go` | Update references to extracted constants |

## Completion Criteria

- [ ] YAML parser handles tab indentation
- [ ] Profile parsing works with variable indent depth
- [ ] No duplicate heading-parser code between vault and markdown packages
- [ ] All hardcoded magic numbers replaced with named constants
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

1 day
