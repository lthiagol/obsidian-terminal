# M16 — Shed Dependencies

**Status:** ⏳ pending

## Goal

Replace `gopkg.in/yaml.v3` and `github.com/charmbracelet/bubbles/viewport` with own implementations. Go from 22 → 19 deps.

## Implementation Plan

### Part A: Replace yaml.v3

**New file: `yamlmini.go`** — shared mini YAML parser:
```go
type yamlKeyValue struct { key, value string; items []string }
func scanYAML(data []byte, fn func(key, value string))  // handles scalar, quoted, inline array, block array
func stripQuotes(s string) string
func parseInlineArray(s string) []string
```

scanYAML handles: `key: value`, `key: "value"`, `key: [a, b]`, block arrays with `- item`, comments, blank lines.

**config.go changes:**
- Remove `gopkg.in/yaml.v3` import
- Remove yaml struct tags from `Config`
- New `parseConfigYAML(data []byte, cfg *Config)` using scanYAML
- skip_dirs replaces defaults (not merges) — matches yaml.Unmarshal behavior

**vault.go changes:**
- Remove yaml struct tags from `frontmatterData`
- Rewrite `parseFrontmatter` using scanYAML
- scanYAML([]byte(yamlBlock), func(key, value string) { switch key { case "title","tags","aliases": ... } })

**Edge cases:** Quoted strings, Windows \r\n, comments, colons in values, invalid YAML (skipped silently), block arrays followed by scalars.

### Part B: Replace bubbles/viewport

**New file: `viewport.go`** — minimal vertical scroll window:
```go
type viewport struct { Width, Height, YOffset int; lines []string }
func newViewport(w, h int) viewport
func (v *viewport) SetContent(s string)     // splits content by \n
func (v viewport) View() string             // visible slice of lines
func (v *viewport) LineUp/Down(n int)       // scroll, clamped
func (v *viewport) SetYOffset(n int)        // absolute offset, clamped
func (v *viewport) GotoBottom()             // last page
func (v *viewport) HalfViewUp/Down()        // half-page scroll
func (v viewport) TotalLineCount() int      // len(lines)
func (v *viewport) clampOffset()            // bounds checking
```

**viewer.go changes:**
- Remove `bubbles/viewport` import
- Change `viewport viewport.Model` → `viewport viewport`
- Change `viewport.New(80,20)` → `newViewport(80,20)`
- Simplify `Update(msg)` → `return *v, nil` (no message processing needed)
- All other methods unchanged — new viewport has same field/method names

**model.go change:** Add content re-render in `WindowSizeMsg` handler so resized viewport gets re-wrapped content (new viewport doesn't auto-wrap).

**Test impact:** None — `viewer_test.go` accesses `v.viewport.YOffset`, `.TotalLineCount()`, `.Width`, `.Height` which all match exactly.

### Part C: Go module cleanup

Run `go mod tidy` — removes yaml.v3, bubbles, and their transitives (cellbuf, displaywidth, stringish, uax29, go-colorful).

### Implementation order
1. Create `yamlmini.go`
2. Rewrite `config.go` + remove yaml tags
3. Rewrite `vault.go` + remove yaml tags
4. Create `viewport.go`
5. Update `viewer.go`
6. Add resize re-render in `model.go`
7. Run tests → `go mod tidy` → tests again
