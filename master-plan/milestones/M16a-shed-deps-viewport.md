# M16a — Replace Viewport Dependency

**Status:** ✅ done

## Goal

Replace `github.com/charmbracelet/bubbles/viewport` with a minimal custom implementation. This reduces dependencies while maintaining all current viewport functionality.

## Motivation

The viewport component is simple enough that a custom implementation reduces external dependencies without significant maintenance burden. The viewport only needs:
- Vertical scrolling with offset tracking
- Content rendering (visible lines only)
- Size management (width/height)
- Scroll methods (line up/down, half-page, goto bottom)

## Implementation Plan

### 1. Create `viewport.go`

New file with minimal viewport implementation:

```go
type viewport struct {
    Width   int
    Height  int
    YOffset int
    lines   []string
}

func newViewport(width, height int) viewport
func (v *viewport) SetContent(content string)     // splits by \n
func (v viewport) View() string                   // visible slice
func (v *viewport) LineUp(n int)                  // scroll up, clamped
func (v *viewport) LineDown(n int)                // scroll down, clamped
func (v *viewport) SetYOffset(n int)              // absolute offset
func (v *viewport) GotoBottom()                   // last page
func (v *viewport) HalfViewUp()                   // half-page up
func (v *viewport) HalfViewDown()                 // half-page down
func (v viewport) TotalLineCount() int            // len(lines)
func (v *viewport) clampOffset()                  // bounds checking
```

### 2. Update `viewer.go`

Replace viewport.Model with custom viewport:

- Remove import: `github.com/charmbracelet/bubbles/viewport`
- Change field: `viewport viewport.Model` → `viewport viewport`
- Change constructor: `viewport.New(80, 20)` → `newViewport(80, 20)`
- Remove Update call: `v.viewport.Update(msg)` no longer needed
- All method calls remain the same (LineUp, LineDown, etc.)

### 3. Handle WindowSizeMsg in `model.go`

Add content re-render when window resizes:

```go
case tea.WindowSizeMsg:
    // ... existing size calculations ...
    if m.activeNote != nil {
        m.viewer.SetContent(m.activeNote.Body, viewerWidth)
    }
```

Custom viewport doesn't auto-reflow content on resize, so we need to re-render.

### 4. Remove bubbles dependency

Run `go mod tidy` to remove:
- `github.com/charmbracelet/bubbles`
- Transitive dependencies (if not used elsewhere)

## Edge Cases

- Empty content: View() returns empty string
- Content shorter than height: no scrolling needed, YOffset stays 0
- Resize to smaller height: clamp YOffset to valid range
- Resize to smaller width: content re-wraps on SetContent call
- Scroll beyond content: clampOffset() prevents invalid offsets

## Testing Strategy

- Unit tests for viewport methods (LineUp, LineDown, GotoBottom, etc.)
- Integration tests via viewer_test.go (existing tests should pass unchanged)
- Manual test: resize terminal window, verify content reflows

## Completion Criteria

- [x] Custom viewport implementation in `viewport.go`
- [x] `viewer.go` uses custom viewport instead of bubbles/viewport
- [x] Window resize triggers content re-render
- [x] `make test` passes all existing tests
- [x] `make vet` exits 0
- [x] `go mod tidy` removes bubbles dependency
- [x] Manual test: scrolling, resizing, all viewport features work
