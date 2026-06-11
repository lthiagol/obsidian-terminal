# M4 — Markdown Viewer

**Status:** ✅ done

## Goal

Render markdown notes on the right panel using `glamour` for syntax highlighting, and `bubbles/viewport` for scrollable navigation. Support `[[wiki-link]]` following.

## Files to create

- `viewer.go` — Markdown renderer + viewport wrapper

## Steps

### 1. `viewer.go`
- `MarkdownViewer` struct:
  - `viewport viewport.Model` — from `bubbles/viewport`
  - `glamour *glamour.TermRenderer` — from `glamour`
  - `links []string` — list of [[wiki-link]] targets extracted from current note
  - `selectedLink int` — index for Tab-based link navigation
- `newViewer()` — initializes glamour renderer with default theme:
  ```go
  renderer, _ := glamour.NewTermRenderer(
      glamour.WithStandardStyle("dark"),
  )
  ```
- `SetContent(markdown string)` — renders via glamour → sets viewport content → extracts wiki-links
- `Update(msg tea.Msg)` — delegates to viewport
- `View()` — returns rendered viewport
- `ScrollUp(n)`, `ScrollDown(n)`, `ScrollTop()`, `ScrollBottom()`
- `ExtractWikiLinks(markdown string) []string` — regex `\[\[([^\]|#]+)` to find link targets
- `SelectedLinkPath()` — returns the path for the currently selected link

### 2. Wiki-link resolution
- `resolveWikiLink(target string, vault *VaultEntry) string`:
  1. Exact relative path match (e.g., `projects/api-design`)
  2. Basename match with `.md` extension (e.g., `api-design` → `projects/api-design.md`)
  3. Case-insensitive basename fallback
- Returns empty string if not found

### 3. View mode key handling (in `model.go`)
- `j/↓` → `viewer.ScrollDown(1)`
- `k/↑` → `viewer.ScrollUp(1)`
- `g/Home` → `viewer.ScrollTop()`
- `G/End` → `viewer.ScrollBottom()`
- `PgUp` → `viewer.ScrollUp(halfPage)`
- `PgDn` → `viewer.ScrollDown(halfPage)`
- `h/Esc/←` → switch back to "browse"
- `Tab` → cycle through/wiki-links
  - If links present: highlight next link in status bar, `Enter` follows it
  - If no links: no-op
- `Enter` (when link selected) → follow wiki-link: load target note, set as activeNote

### 4. Edge cases
- Empty note → show "(empty note)" placeholder
- Note with no frontmatter → render body directly
- Very long note (1000+ lines, 50KB+) → viewport handles with lazy rendering
- Note with no wiki-links → Tab does nothing
- Wiki-link to non-existent note → show error toast

## Test Spec (~8 tests)

| Test | File | Description |
|------|------|-------------|
| `TestViewer_RendersMarkdown` | viewer_test.go | glamour output contains ANSI escape codes and original text |
| `TestViewer_ScrollDown` | viewer_test.go | j/↓ increments viewport Y offset |
| `TestViewer_ScrollUp` | viewer_test.go | k/↑ decrements viewport Y offset |
| `TestViewer_ScrollToTop` | viewer_test.go | g/Home sets Y offset to 0 |
| `TestViewer_ScrollToBottom` | viewer_test.go | G/End sets Y offset to total content height |
| `TestViewer_BackToBrowse` | model_test.go | h/Esc/← exits view mode, mode returns to "browse" |
| `TestViewer_WikiLinkExtraction` | viewer_test.go | `ExtractWikiLinks` finds all `[[target]]` without `#` fragments |
| `TestViewer_LongNoteScroll` | viewer_test.go | Viewport handles content taller than terminal (multiple pages) |

## Completion Criteria

- [ ] Glamour renders markdown with syntax-highlighted code blocks, tables, callouts, etc.
- [ ] Viewport scrolls via j/k/↑↓, g/G/Home/End, PgUp/PgDn
- [ ] h/Esc/← returns to browse mode
- [ ] Tab cycles through [[wiki-links]]; Enter follows
- [ ] Wiki-link resolution finds notes by exact path, basename, case-insensitive
- [ ] Empty notes and missing wiki-link targets handled gracefully
- [ ] All 8 tests pass
- [ ] `go vet ./...` exits 0
