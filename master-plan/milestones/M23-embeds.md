# M23 — Embedded Block Embeds

**Status:** ⏳ pending

## Goal

Parse `![[note]]` and `![[note#heading]]` syntax, load referenced content, render inline with border.

## Implementation Plan

### 1. Markdown parser changes (`internal/markdown/markdown.go`)

Add `BlockEmbed` to BlockType constants.

Add fields to `MarkdownLine`: `EmbedTarget string`, `EmbedHeading string` (target note path + optional heading reference).

In `ParseMarkdown`, add check before paragraph detection (~line 148): if line starts with `![[` and ends with `]]`, extract target (split on `#` for heading), create `BlockEmbed` line.

### 2. Embed resolution (`internal/markdown/markdown.go`)

```go
type EmbedResolver func(target, heading string) (string, error)

func ResolveEmbeds(lines []MarkdownLine, resolve EmbedResolver) []MarkdownLine
```

Walks lines, for each `BlockEmbed`: calls resolver, parses result via `ParseMarkdown`, wraps in `BlockEmbedStart`/`BlockEmbedEnd` sentinel lines (new BlockType values).

### 3. Render embeds (`internal/markdown/markdown.go`)

Add `BlockEmbedStart`, `BlockEmbedEnd` constants.  
New render function: `renderEmbedBlock(lines []MarkdownLine, ...)` — renders with left border + source header.

### 4. Viewer integration (`viewer.go`)

In `SetContent`, after `ParseMarkdown`, call `ResolveEmbeds` with a closure that calls `LoadNote` and `extractSection`:

```go
resolveEmbed := func(target, heading string) (string, error) { ... }
lines = markdown.ResolveEmbeds(lines, resolveEmbed)
```

### 5. Section extraction (`vault.go`)

`extractSection(markdown, heading string) string` — finds heading in raw markdown, returns content under it until next heading of same or higher level, or EOF.

### Edge cases

- Embed target doesn't exist → render placeholder "(embed not found: target)"
- Circular embeds (A embeds B, B embeds A) → no recursion guard needed since we only load top-level embeds (not recursive)
- Heading not found → render full note content
- Embed width → use same viewport width

### Implementation order

1. Add BlockEmbed type + embed fields to MarkdownLine
2. Parse `![[` syntax in ParseMarkdown
3. Add ResolveEmbeds + sentinel types
4. Add renderEmbedBlock
5. Add extractSection to vault.go
6. Wire into viewer.go SetContent
7. Write tests
