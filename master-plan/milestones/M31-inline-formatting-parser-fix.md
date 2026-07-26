# M31 — Inline Formatting Parser Fix

**Status:** ⏳ pending

## Goal

No raw markup characters ( `*`, `_`, `` ` ``, `~`, `=` ) leak into rendered output. All valid inline formatting renders correctly; all invalid/ambiguous input degrades gracefully to literal text without infinite recursion.

## Problem

The `parseSegments` function in `internal/markdown/markdown.go:444` has three bugs that cause raw markup to appear in rendered output and can crash the parser:

### Bug 1: Linear fallthrough misparses `***`

When text starts with `***` but no closing `***` exists, the parser falls through to `**` which consumes the first two `*` as bold markers, leaving the orphan `*` as literal content. The `**` closing search then finds the wrong delimiter deeper in the text.

```
Input: "***text**"
  → Try *** (bold+italic) → no closing *** found
  → Fall through to ** (bold) → consumes first ** → searches for closing **
    → Finds ** after "text" → inner = "*text" → Bold("*text") ← WRONG
  → Output: \033[1m*text\033[0m   (leading * rendered as bold text)
```

Correct parse: `*` (literal) + `**text**` (bold) → `*\033[1mtext\033[0m`

The same pattern affects `___` → `__` → `_` with underscores.

### Bug 2: Infinite recursion on unmatched markers

When `findNextSpecial` returns position 0 (the current position) — meaning an unparseable special character sits at the start of text — `parseSegments` recurses infinitely:

```
text = "`"   (lone backtick, no closing)
  → ` check → end = Index(text[1:], "`") = -1 → falls through
  → ~~ check → no → falls through
  → == check → no → falls through
  → findNextSpecial("`") → loc = [0, 1] → returns 0
  → next == 0, text[:0] = "" → append empty segment
  → parseSegments(text[0:]) = parseSegments("`") → INFINITE RECURSION
```

This is catastrophic — it hangs the entire TUI. Observed during test writing (M29 tests uncovered this by accident).

### Bug 3: Double-backtick not handled

Obsidian supports `` `` ` `` `` for inline code containing backticks. Our parser only recognizes single `` ` ``:

```
Input: "``code with ` backtick``"
  → Parses as: "`" + "`code with " + code(" backtick") ← WRONG
  Correct: code("code with ` backtick")
```

## Design

### Fix 1: Backtrack on failed long markers

When a marker search fails, instead of falling through linearly, the parser retries with the next-shorter marker starting at the appropriate offset.

New flow for `***`:
```
Input: "***text**"
  → Try *** → find closing in text[3:] → not found
  → Retry from offset 1: text[1:] = "**text**"
    → Try ** → find closing in text[1+2:] = "text**" → found at pos 4 → inner="text"
    → Bold("text") ✓
  → Remaining: text[1+2+4+2:] = "" → done
  → Output: * + \033[1mtext\033[0m
```

The key insight: when `***` fails, we advance 1 char and try `**` from there. When `**` fails, we advance 1 more char and try `*` from there. This preserves the progressive marker shortening without consuming the wrong characters.

Implementation sketch:
```go
func tryMarker(text string, marker string, pos int) (consumed int, segment InlineSegment, ok bool) {
    // Try to find closing marker starting at `pos`
    // Returns (total chars consumed, segment, found?)
}

// In parseSegments:
if strings.HasPrefix(text, "***") || strings.HasPrefix(text, "___") {
    if inner, end, ok := tryMarker(text, text[:3], 3); ok {
        // bold+italic found
    } else {
        // *** failed — retry with ** from offset 1
        goto tryDoubleAsterisk
    }
}
```

### Fix 2: Infinite recursion guard

After all marker handlers fail, check if `findNextSpecial` returned position 0. If so, consume the problematic character as literal text and advance:

```go
next := findNextSpecial(text)
if next == -1 {
    *segments = append(*segments, InlineSegment{Text: text})
    return
}
if next == 0 {
    // Stuck: special char at pos 0 that we can't parse. Consume it.
    r, size := utf8.DecodeRuneInString(text)
    *segments = append(*segments, InlineSegment{Text: string(r)})
    parseSegments(text[size:], segments)
    return
}
if next > 0 {
    *segments = append(*segments, InlineSegment{Text: text[:next]})
}
parseSegments(text[next:], segments)
```

### Fix 3: Double-backtick code spans

Add a `` `` ` ` `` `` handler before the single `` ` `` handler:

```go
if strings.HasPrefix(text, "``") {
    end := strings.Index(text[2:], "``")
    if end >= 0 {
        inner := text[2 : 2+end]
        segments = append(*segments, InlineSegment{Text: inner, Code: true})
        parseSegments(text[2+end+2:], segments)
        return
    }
}
```

## Files to modify

| File | Changes |
|------|---------|
| `internal/markdown/markdown.go` | Rewrite `parseSegments` marker handling: add backtracking, recursion guard, double-backtick support |
| `internal/markdown/markdown_test.go` | Add tests for all fixed cases + regression tests for the bug scenarios |

## Steps

### 1. Infinite recursion guard (safety first)
Add the `next == 0` guard to `parseSegments`. This prevents the TUI from hanging on any future parser edge case. Write tests: lone `` ` ``, lone `*`, lone `~`, lone `=`.

### 2. Backtrack on `***` / `___` failures
Restructure the marker chain so that when `***` or `___` fails, the parser retries with the 2-char marker from offset 1. Write tests: `"***text**"`, `"___text__"`, `"***text"` (triple with no close at all), `"___"` alone.

### 3. Backtrack on `**` / `__` failures  
When `**` or `__` fails, retry with 1-char marker from offset 1. Tests: `"**text*"`, `"__text_"`.

### 4. Double-backtick code spans
Add `` `` ` ` `` `` handler before single `` ` ``. Tests: `"``code``"`, `"``code with ` backtick``"`, `"````"` (empty code span), `` "`single` and ``double``" ``.

### 5. Edge case sweep
Test all previously-crashing inputs from the infinite recursion bug. Test adjacent markers: `"**bold**`code`*italic*"`. Test unbalanced markers throughout a paragraph.

### 6. Visual diff test
Render a known markdown string through `ParseMarkdown` → `RenderMarkdown` → `viewport.SetContent` and verify no raw `*`, `_`, `` ` ``, `~`, `=` appear outside of code blocks.

## Completion Criteria

- [ ] No raw `*` `_` `` ` `` `~` `=` appear in rendered output (except in code blocks)
- [ ] `"***text**"` renders as `*` + bold(`text`) — not bold with embedded asterisk
- [ ] `"___text__"` renders as `_` + bold(`text`)
- [ ] Lone special characters (`` ` ``, `*`, `~`, `=`) don't crash the parser
- [ ] `` `` ` `` `` double-backtick code spans work correctly
- [ ] `make test` passes all tests (no infinite recursion hang)
- [ ] `make test-race` passes (ensuring no goroutine issues from the parser)
- [ ] `make vet` exits 0
