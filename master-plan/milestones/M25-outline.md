# M25 — Outline / Table of Contents

**Status:** ✅ done

## Goal

Show headings from current note as navigable outline overlay. Toggle with `t` in View mode, Enter to jump to heading.

## Keybinding

**Key:** `t` (lowercase)
**Mode:** View mode only
**Rationale:** Mnemonic for "table of contents", only active in View mode so no conflict with `T` (tags) in Browse mode.

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Implementation Plan

### 1. Extract heading information

Add to `internal/markdown/markdown.go`:

```go
type HeadingInfo struct {
    Level   int
    Text    string
    LineIdx int  // index in parsed lines
}

func ExtractHeadings(lines []MarkdownLine) []HeadingInfo {
    var headings []HeadingInfo
    for i, line := range lines {
        if line.BlockType == BlockHeading {
            text := renderSegmentsPlain(line.Segments)
            headings = append(headings, HeadingInfo{
                Level:   line.HeadingLevel,
                Text:    text,
                LineIdx: i,
            })
        }
    }
    return headings
}

// renderSegmentsPlain renders segments without styling (for outline)
func renderSegmentsPlain(segments []InlineSegment) string {
    var sb strings.Builder
    for _, seg := range segments {
        sb.WriteString(seg.Text)
    }
    return sb.String()
}
```

### 2. Define outline types

Add to `model.go`:

```go
type OutlineItem struct {
    Level   int
    Text    string
    LineIdx int    // index in parsed markdown lines
    YOffset int    // approximate Y offset in rendered viewport
}
```

### 3. Add outline state to Model

```go
type Model struct {
    // ... existing fields ...
    outlineVisible bool
    outlineItems   []OutlineItem
    outlineCursor  int
}
```

### 4. Implement outline methods

```go
func (m *Model) buildOutline() {
    if m.activeNote == nil {
        m.outlineItems = nil
        return
    }
    
    lines := markdown.ParseMarkdown(m.activeNote.RawBody)
    headings := markdown.ExtractHeadings(lines)
    
    m.outlineItems = make([]OutlineItem, len(headings))
    for i, h := range headings {
        m.outlineItems[i] = OutlineItem{
            Level:   h.Level,
            Text:    h.Text,
            LineIdx: h.LineIdx,
            YOffset: estimateYOffset(lines, h.LineIdx, m.viewer.Width),
        }
    }
    
    m.outlineCursor = 0
}

func (m Model) renderOutline() string {
    if len(m.outlineItems) == 0 {
        return lipgloss.NewStyle().
            Foreground(TextMuted).
            Render("No headings in this note")
    }
    
    var sb strings.Builder
    for i, item := range m.outlineItems {
        indent := strings.Repeat("  ", item.Level-1)
        line := fmt.Sprintf("%s%s", indent, item.Text)
        
        if i == m.outlineCursor {
            line = lipgloss.NewStyle().
                Background(Accent).
                Foreground(lipgloss.Color("#000000")).
                Bold(true).
                Render(line)
        } else {
            line = lipgloss.NewStyle().
                Foreground(TextSecondary).
                Render(line)
        }
        
        sb.WriteString(line)
        if i < len(m.outlineItems)-1 {
            sb.WriteString("\n")
        }
    }
    return sb.String()
}

// estimateYOffset approximates the Y offset for a given line index
func estimateYOffset(lines []markdown.MarkdownLine, targetIdx, width int) int {
    yOffset := 0
    for i := 0; i < targetIdx && i < len(lines); i++ {
        line := lines[i]
        switch line.BlockType {
        case markdown.BlockEmpty:
            yOffset++
        case markdown.BlockHeading:
            yOffset++ // heading takes 1 line
        case markdown.BlockCodeBlock:
            codeLines := strings.Count(line.RawContent, "\n") + 1
            yOffset += codeLines + 2 // +2 for borders
        default:
            // Estimate wrapped lines
            text := renderSegmentsPlain(line.Segments)
            wrappedLines := (len(text) / width) + 1
            yOffset += wrappedLines
        }
    }
    return yOffset
}
```

### 5. Add keybinding

Add to `keys.go`:

```go
type KeyMap struct {
    // ... existing fields ...
    Outline rune  // 't' in View mode
}

func DefaultKeys() KeyMap {
    return KeyMap{
        // ... existing bindings ...
        Outline: 't',
    }
}
```

### 6. Add handler

Add to `handlers.go`:

```go
func (m Model) handleViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch {
    // ... existing cases ...
    
    case MatchRune(msg, m.keys.Outline):
        if m.outlineVisible {
            m.outlineVisible = false
        } else {
            m.buildOutline()
            m.outlineVisible = true
        }
        return m, nil
    }
    return m, nil
}

func (m Model) handleOutlineKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch {
    case msg.Type == tea.KeyEsc || MatchRune(msg, m.keys.Outline):
        m.outlineVisible = false
        return m, nil
    
    case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
        if m.outlineCursor < len(m.outlineItems)-1 {
            m.outlineCursor++
        }
        return m, nil
    
    case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
        if m.outlineCursor > 0 {
            m.outlineCursor--
        }
        return m, nil
    
    case msg.Type == tea.KeyEnter:
        if m.outlineCursor < len(m.outlineItems) {
            item := m.outlineItems[m.outlineCursor]
            m.viewer.ScrollTop()
            m.viewer.ScrollDown(item.YOffset)
            m.outlineVisible = false
        }
        return m, nil
    }
    return m, nil
}
```

### 7. Update Update() dispatch

In `model.go` Update():

```go
case tea.KeyMsg:
    // ... existing global handlers ...
    
    if m.outlineVisible {
        return m.handleOutlineKey(msg)
    }
    
    switch m.mode {
    // ... existing mode dispatch ...
    }
```

### 8. Update View()

In `model.go` View():

```go
var rightPanel string
switch m.mode {
// ... existing cases ...
case ModeView:
    if m.outlineVisible {
        rightPanel = m.renderOutline()
    } else {
        rightPanel = m.viewer.View()
    }
}
```

### 9. Build outline on note load

At all note-load sites, call `m.buildOutline()`:
- `handleBrowseKey` Enter case
- `handleSearchOrFind` Enter case
- `handleViewKey` link follow case
- `rescanVault` note reload

### 10. Update help text

Add to `help.go`:

```go
{
    title: "Outline",
    bindings: []string{
        "t      — toggle outline",
        "j / k  — navigate headings",
        "Enter  — jump to heading",
        "Esc    — close outline",
    },
},
```

## Edge Cases

- **Note with no headings**: Show "No headings in this note"
- **YOffset approximation**: Accept "good enough" for v1 (lines vary with wrapping)
- **Outline + window resize**: Outline stays visible, not affected by SetSize
- **`t` in other modes**: Only works in ModeView when outline not visible
- **Long heading text**: Truncate to fit panel width

## Testing Strategy

- Unit tests for ExtractHeadings
- Unit tests for estimateYOffset
- Integration tests: buildOutline populates items correctly
- Integration tests: handleOutlineKey navigates and jumps
- Manual test: outline shows all headings, navigation works, jump works

## Completion Criteria

- [x] ExtractHeadings function in markdown.go
- [x] OutlineItem type and outline state in Model
- [x] buildOutline and renderOutline methods implemented
- [x] handleOutlineKey handler implemented
- [x] `t` keybinding added to KeyMap
- [x] Outline dispatch in Update() before mode dispatch
- [x] Outline rendering in View() when visible
- [x] buildOutline called at all note-load sites
- [x] Help text updated
- [x] KEYBINDINGS.md updated
- [x] `make test` passes
- [x] `make vet` exits 0
- [x] Manual test: outline works for notes with headings
