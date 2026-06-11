# M23 — Embedded Block Embeds

**Status:** ⏳ pending

## Goal

Parse and render Obsidian block embed syntax (`![[note#heading]]` and `![[note]]`) inline in the viewer, showing the referenced content directly.

## Steps

### 1. Parse `![[` embed syntax

Already handled by the markdown parser's wiki-link detection. Distinguish `![[` embeds from `[[` links and mark them as embeds.

### 2. Load and render embedded content

For `![[note.md]]`, load the full note and render it inline (first few lines or full content depending on config). For `![[note.md#heading]]`, load the note and extract only the section under the specified heading.

### 3. Render embeds in the viewer

Render embeds as indented blocks with a left border and the source note name as a header. Clickable to navigate to the source.

## Completion Criteria

- [ ] `![[note]]` renders full note content inline
- [ ] `![[note#heading]]` renders only the section under that heading
- [ ] Embeds are visually distinct from regular content
- [ ] Clickable to open the source note
- [ ] `make test && make vet` pass
