# M16 — Shed Dependencies

**Status:** ⏳ pending

## Goal

Remove `gopkg.in/yaml.v3` and `github.com/charmbracelet/bubbles` by implementing their functionality ourselves. Reduce from 22 → 19 deps (all remaining are bubbletea/lipgloss transitives).

## Files to modify

- `config.go` — replace yaml.Unmarshal with own parser
- `vault.go` — replace yaml.Unmarshal with own parser
- `viewer.go` — replace bubbles/viewport with own viewport
- `viewer_test.go` — update for new viewport API
- `go.mod` — remove deps

## Steps

### 1. Replace yaml.v3 with own parser

We parse two YAML shapes, both trivial:

```yaml
# Config (config.go)
vault_path: /path/to/vault
theme: dark
default_keys: vim
skip_dirs:
  - .obsidian
  - .git

# Frontmatter (vault.go)  
---
title: My Note
tags:
  - tag1
  - tag2
aliases:
  - alias1
---

Add `parseYAMLConfig(data []byte) map[string][]string` to handle both. Key-value pairs for scalars, indented `- item` for arrays. No nesting, no quotes, no complex types needed.

### 2. Replace bubbles/viewport with own viewport

The viewport is: content string → split into lines → track Y-offset → render visible window.

Add a `Viewport` struct in `viewer.go`:

```go
type Viewport struct {
    content string
    lines   []string
    yOffset int
    Width   int
    Height  int
}
```

Methods: `SetContent(s)`, `View() string`, `LineUp(n)`, `LineDown(n)`, `SetYOffset(n)`, `GotoBottom()`, `HalfViewUp()`, `HalfViewDown()`, `TotalLineCount() int`

Update `MarkdownViewer` to use `Viewport` instead of `viewport.Model`. Remove `Update(tea.Msg)` from the viewer (we handle key events directly in handlers.go).

### 3. Cleanup

- Remove `gopkg.in/yaml.v3` and `github.com/charmbracelet/bubbles` from go.mod
- Run `go mod tidy`
- Update README dependency table

## Completion Criteria

- [ ] Own YAML parser in `config.go`/`vault.go` (no yaml.v3 import)
- [ ] Own viewport in `viewer.go` (no bubbles/viewport import)
- [ ] `go mod` shows 19 deps (down from 22)
- [ ] All 98 tests pass
- [ ] `make build && make vet` exit 0
