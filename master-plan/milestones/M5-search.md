# M5 — Search (Fuzzy Name + Full-Text Content)

**Status:** ✅ done

## Verification Evidence

- `go build ./...` exits 0
- `go test ./...` — 58/58 tests pass
- `go vet ./...` exits 0
- Files created: `search.go`, `search_test.go`

## Goal

Implement fuzzy file name search (`/`) and full-text content search (`Ctrl+F`)
with real-time filtering, cached content index, and 200ms type debounce.

## Files to create

- `search.go` / `search_test.go`

## Steps

### 1. `search.go`
- `SearchMode` enum:
  ```go
  type SearchMode int
  const (
      SearchName SearchMode = iota
      SearchContent
  )
  ```
- `SearchState` struct:
  ```go
  type SearchState struct {
      mode        SearchMode
      query       string
      results     []SearchResult
      selected    int
      textInput   textinput.Model
      allPaths    []string            // all .md file paths (for fuzzy)
      searchIndex map[string]string   // filename → plain text (built in M1)
  }
  ```
- `SearchResult` struct:
  ```go
  type SearchResult struct {
      Path    string  // relative file path
      Score   float64 // for fuzzy ranking (0 for content)
      Context string  // matching line with line number (content search)
  }
  ```
- `NewSearchState(mode SearchMode, allPaths []string, index map[string]string) SearchState`

### 2. Fuzzy name search
- `FuzzyScore(query, target string) float64`:
  - +12 bonus: query matches start of target
  - +10 bonus: matching character is at word boundary (`/`, `-`, `_`, ` `, `.`)
  - +8 bonus: consecutive matching character
  - +2 bonus: exact case match
  - –1 penalty: gap between matches (cumulative)
  - Returns 0 if no characters match
  - All comparisons case-insensitive
- `FuzzySearch(query string, paths []string) []SearchResult`:
  - Score each path, filter score > 0
  - Sort descending by score
  - Limit to 50 results
- `HighlightMatches(query, path string) string`:
  - Returns path with matching characters wrapped in ANSI bold+underline

### 3. Full-text content search
- `ContentSearch(query string, index map[string]string) []SearchResult`:
  - Iterate index map entries
  - Case-insensitive `strings.Contains(body, query)`
  - For matches: find line containing match, show line number + text (truncated to 80 chars)
  - Sort by path alphabetically (not scored)
  - Limit to 100 results

### 4. Search mode key handling (in `model.go`)
- `/` → open search in SearchName mode
- `Ctrl+F` → open search in SearchContent mode
- Character keys → append to query; update results (debounced at 200ms for content search)
- `Backspace` → remove last char; update results
- `↑/k` → selected-- (clamped to 0)
- `↓/j` → selected++ (clamped to len(results)-1)
- `Enter` → open selected result:
  - Load note via resolved path
  - Set activeNote
  - Switch mode to "view"
  - Clear search state
- `Esc` → cancel search; restore prevMode; clear query

### 5. Search overlay rendering
- Top bar (full width):
  - Name search: `"🔍 fuzzy  query_here  (N results)"`
  - Content search: `"🔍 content  query_here  (N matches in M files)"`
- Results list (scrollable):
  - Name search: highlighted filename + score indicator
  - Content search: `path:line  matching line text...`
- Selected item: violet highlight
- Unselected: gray text
- "No results" when results empty and query non-empty

### 6. Search index building (in M1 `ScanVault`)
- During vault scan, for each `.md` file:
  - Read and store full plain text in `searchIndex[relativePath]`
  - Strip frontmatter before indexing
- Index rebuilt on vault refresh (M7)

## Test Spec (10 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestFuzzyScore_ExactMatch` | search_test.go | Identical strings score > 100 |
| 2 | `TestFuzzyScore_Substring` | search_test.go | Consecutive matching chars score higher than scattered |
| 3 | `TestFuzzyScore_BoundaryBonus` | search_test.go | Word boundary match gets +10 bonus |
| 4 | `TestFuzzyScore_ExactCaseBonus` | search_test.go | Exact case match gets +2 bonus |
| 5 | `TestFuzzyScore_NoMatch` | search_test.go | Returns 0 when no characters match |
| 6 | `TestFuzzySearch_ResultsSorted` | search_test.go | Results ranked by score descending |
| 7 | `TestFuzzySearch_NoMatchingFiles` | search_test.go | Returns empty slice when no file matches query |
| 8 | `TestContentSearch_FindsInBody` | search_test.go | Index-based search finds content in note bodies |
| 9 | `TestContentSearch_ReturnsLineContext` | search_test.go | Results include file path, line number, matching text |
| 10 | `TestContentSearch_CaseInsensitive` | search_test.go | "TODO" matches "todo" in content |

## Completion Criteria

- [ ] `/` opens fuzzy name search; `Ctrl+F` opens content search
- [ ] Fuzzy search ranks by score, shows top 50
- [ ] Content search uses cached index (no disk reads per keystroke)
- [ ] Content search debounced at 200ms
- [ ] Real-time filtering as user types
- [ ] `↑↓` navigate results; `Enter` opens note; `Esc` cancels
- [ ] Search overlay renders result count and highlighted matches
- [ ] Empty query shows all files (fuzzy) or no results (content)
- [ ] All 10 tests pass
- [ ] `go vet ./...` exits 0
