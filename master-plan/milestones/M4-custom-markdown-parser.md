# M4 — Custom Markdown Parser + Viewer

**Status:** ✅ done

## Goal

Build a custom Obsidian-flavored markdown parser that converts `.md` content into
styled ANSI output. Integrate with `bubbles/viewport` for scrolling and wiki-link
Tab-cycle navigation. **No glamour — fully custom.**

## Files to create

- `markdown.go` / `markdown_test.go` — parser + renderer
- `viewer.go` / `viewer_test.go` — viewport wrapper + wiki-link navigation

## Steps

### 1. `markdown.go` — Parser

#### Types
```go
type BlockType int
const (
    BlockParagraph BlockType = iota
    BlockHeading
    BlockCodeBlock
    BlockList
    BlockBlockquote
    BlockCallout
    BlockHorizontalRule
    BlockEmpty
)

type InlineSegment struct {
    Text         string
    Bold         bool
    Italic       bool
    Strikethrough bool
    Code         bool
    Highlight    bool
    IsWikiLink   bool
    WikiTarget   string
    WikiDisplay  string
}

type MarkdownLine struct {
    BlockType    BlockType
    HeadingLevel int      // 1-6
    Segments     []InlineSegment
    IndentLevel  int      // for nested lists/quotes
    CalloutType  string   // "note", "warning", "danger", etc.
    Language     string   // fenced code block language
}
```

#### Parsing pipeline
1. **Split into blocks** — split input by `\n\n`, classify each block
2. **Parse blocks** — line-by-line state machine for each block:
   - Heading: count `#`, extract text
   - Code block: detect ` ``` ` fences, capture language, accumulate content
   - List: detect `- `, `* `, `+ `, `1. ` prefixes
   - Blockquote: detect `> ` prefix; callout if `[!type]` follows
   - Horizontal rule: detect `---`, `***`, `___`
   - Paragraph: everything else
3. **Parse inline** — within paragraphs/headings/list items:
   - `**bold**`, `***bold italic***`
   - `*italic*` (but not `*` in lists — already handled)
   - `` `code` ``
   - `~~strikethrough~~`
   - `==highlight==`
   - `[[target]]` → wiki-link
   - `[[target|display]]` → wiki-link with display text
   - Ignore `%%comments%%` — strip entirely
4. **Strip frontmatter** — if content starts with `---`, skip until closing `---`

#### Rendering
```go
func RenderMarkdown(lines []MarkdownLine, width int, styles Styles) string
```
- Per-line rendering with lipgloss:
  - h1: fuchsia, bold, underlined
  - h2: violet, bold
  - h3: teal, bold
  - h4-h6: gray, bold
  - Code block: dark background, emerald text, box-drawing frame (`╭╮╰╯`)
  - Inline code: emerald text on subtle bg
  - Bold: bold
  - Italic: italic
  - Bold+Italic: bold italic
  - Strikethrough: crossout
  - Highlight: yellow text
  - Wiki-link: teal, underlined
  - Blockquote: gray, italic, `│` prefix
  - Callout: colored icon + bold type, dimmed body (e.g., `ℹ Note: text`)
  - List: `•` bullet for unordered, `1.` for ordered
  - Horizontal rule: `───` dimmed

#### Wiki-link extraction
```go
func ExtractWikiLinks(lines []MarkdownLine) []WikiLink {
    // Collect all InlineSegment with IsWikiLink == true
    // Returns deduplicated list of {Target, Display}
}
```

#### Word wrapping
- ANSI-aware width calculation (count visible chars, ignore escape sequences)
- Wrap by `width`, splitting at word boundaries
- Preserve segment formatting across wrapped lines

### 2. `viewer.go` — Viewer
- `MarkdownViewer` struct:
  ```go
  type MarkdownViewer struct {
      viewport     viewport.Model  // from bubbles/viewport
      rawMarkdown  string
      links        []WikiLink
      selectedLink int             // -1 when no link selected
  }
  ```
- `NewViewer() MarkdownViewer` — initializes viewport with default size
- `SetContent(markdown string, width int)`:
  - Parse markdown → `[]MarkdownLine`
  - Render → ANSI string
  - Extract wiki-links
  - Set viewport content
- `SetSize(width, height int)` — resize viewport
- `ScrollUp(n)`, `ScrollDown(n)`, `ScrollTop()`, `ScrollBottom()`

### 3. View mode key handling (in `model.go`)
- `j/↓` → viewer.ScrollDown(1)
- `k/↑` → viewer.ScrollUp(1)
- `g/Home` → viewer.ScrollTop()
- `G/End` → viewer.ScrollBottom()
- `PgUp` → viewer.ScrollUp(halfPage)
- `PgDn` → viewer.ScrollDown(halfPage)
- `h/Esc/←` → mode = prevMode (browse)
- `Tab` → cycle selectedLink: `(selectedLink + 1) % len(links)`
  - Highlight active link in status bar: "→ projects/api-design.md"
  - If no links: show "No links" briefly in status bar
- `Enter` (when selectedLink >= 0) → follow wiki-link:
  - Resolve link path
  - Load note, set as activeNote
  - Reset selectedLink to -1

### 4. Wiki-link resolution (in `vault.go` or `viewer.go`)
```go
func ResolveWikiLink(target string, vault *VaultEntry) string
```
1. Strip `#Heading` and `#^block-id` fragments
2. Exact path match (append `.md` if missing)
3. Basename match: scan all files for matching filename
4. Case-insensitive basename fallback
5. Check aliases frontmatter for match
   Return resolved relative path or empty string.

## Test Spec (10 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestParseMarkdown_Headings` | markdown_test.go | h1-h6 produce correct BlockType and HeadingLevel |
| 2 | `TestParseMarkdown_InlineFormatting` | markdown_test.go | Bold, italic, code, strikethrough, highlight all parsed correctly |
| 3 | `TestParseMarkdown_CodeBlocks` | markdown_test.go | Fenced code blocks with language label, no inline parsing inside |
| 4 | `TestParseMarkdown_WikiLinks` | markdown_test.go | `[[target]]` and `[[target\|display]]` parsed correctly |
| 5 | `TestParseMarkdown_Callouts` | markdown_test.go | `> [!note]` etc. produce correct BlockCallout with type |
| 6 | `TestParseMarkdown_Blockquotes` | markdown_test.go | `>` lines produce BlockBlockquote |
| 7 | `TestParseMarkdown_Lists` | markdown_test.go | Unordered (-, *, +) and ordered (1.) lists parsed |
| 8 | `TestParseMarkdown_CommentsStripped` | markdown_test.go | `%%hidden%%` content removed from output |
| 9 | `TestExtractWikiLinks` | markdown_test.go | Returns all unique wiki-link targets from parsed document |
| 10 | `TestRenderMarkdown_ANSIContent` | markdown_test.go | Output contains expected text and ANSI escape codes |

## Completion Criteria

- [x] Parser handles: headings, inline formatting, code blocks, wiki-links, lists, blockquotes, callouts, horizontal rules, empty lines
- [x] Frontmatter stripped from rendered output
- [x] Comments (`%%`) stripped from rendered output
- [x] Wiki-links extracted with targets (stripping fragments)
- [x] Rendered output uses ANSI via lipgloss for all formatting
- [x] Word wrapping preserves inline formatting across line breaks
- [x] Viewer scrolls via j/k/↑/↓, g/G, PgUp/PgDn
- [x] Tab cycles wiki-links; Enter follows resolved link
- [x] Wiki-link resolution: exact path, basename, case-insensitive
- [x] h/Esc exits view mode back to browse
- [x] All 21 tests pass (13 markdown + 8 viewer)
- [x] `go vet ./...` exits 0

## Verification Evidence

- `go build ./...` exits 0
- `go test ./...` — 46/46 tests pass total
- `go vet ./...` exits 0
- Files created: `markdown.go`, `markdown_test.go`, `viewer.go`, `viewer_test.go`
