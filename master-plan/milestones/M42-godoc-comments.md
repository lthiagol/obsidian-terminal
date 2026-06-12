# M42 — Godoc Comments

**Status:** ⏳ pending

## Goal

Add missing godoc comments on all exported symbols to improve code documentation and IDE support.

## Issues

### M6: Missing godoc comments

15+ exported symbols lack godoc comments, including:

**Root package:**
- `Config` struct and all its fields
- `Profile` struct
- `CustomTheme` struct
- `KeyMap` struct and all its fields
- `MatchKey`, `MatchRune` functions
- `DefaultKeys` function
- `VaultIndexes` struct
- `VaultEntry` struct
- `VaultNote` struct
- `frontmatterData` struct (if exported)

**internal/markdown package:**
- `RendererStyle` struct
- `InlineSegment` struct
- `MarkdownLine` struct
- `WikiLink` struct
- `HeadingInfo` struct
- `EmbedResolver` type
- `BlockType` constants

**internal/search package:**
- `Result` struct
- `State` struct
- `Style` struct
- `Mode` type

## Files to modify

| File | Changes |
|------|---------|
| `config.go` | Add godoc comments on `Config`, `Profile`, `CustomTheme` |
| `keys.go` | Add godoc comments on `KeyMap`, `MatchKey`, `MatchRune`, `DefaultKeys` |
| `vault.go` | Add godoc comments on `VaultIndexes`, `VaultEntry`, `VaultNote` |
| `internal/markdown/markdown.go` | Add godoc comments on all exported types |
| `internal/search/search.go` | Add godoc comments on all exported types |

## Completion Criteria

- [ ] All exported symbols have godoc comments
- [ ] Comments follow Go conventions (start with symbol name, describe purpose)
- [ ] `go doc` shows useful documentation for all public APIs
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0

## Estimated Time

2 hours
