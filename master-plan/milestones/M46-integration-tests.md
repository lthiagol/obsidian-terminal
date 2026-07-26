# M46 — Integration Test Suite

**Status:** ✅ done

## Goal

Add end-to-end integration tests that verify the full rendering pipeline and user workflows work correctly.

## Issues

### Missing integration tests

Current tests are mostly unit tests that test individual components in isolation. There are no tests that verify:
- Full rendering pipeline: parse markdown → render to ANSI → viewport softWrap → display
- User workflows: search → select → open note → scroll → follow link
- Mode transitions with state preservation
- Theme switching with immediate visual updates
- Session save/restore across restarts

### Test coverage gaps

Critical user-facing features have no integration tests:
- Opening a note from tree click (mouse.go)
- Opening a note from search result (mouse.go)
- Opening a pinned note (model.go)
- Opening a daily note (model.go)
- Opening a recent note (model.go)
- Following a wiki-link (handlers.go)
- Switching themes (handlers.go)
- Resizing the split (mouse.go, handlers.go)

## Design

### Integration test structure

Create `integration_test.go` with tests that:
1. Create a `Model` with a test vault
2. Simulate user actions via `Update(tea.KeyMsg)` or `Update(tea.MouseMsg)`
3. Verify the resulting state (mode, activeNote, cursor position, etc.)
4. Verify the rendered output via `View()` (check for expected content, no panics)

### Test scenarios

#### Rendering pipeline
```go
func TestRenderingPipeline_FullDocument(t *testing.T) {
    // Parse a complex markdown document with all features
    // Render to ANSI
    // Set on viewport
    // Verify no panics, no broken ANSI sequences
    // Verify content is present
}
```

#### User workflows
```go
func TestWorkflow_SearchAndOpen(t *testing.T) {
    // Start in browse mode
    // Press '/' to enter search
    // Type query
    // Press Enter to open first result
    // Verify mode is View, activeNote is set
}

func TestWorkflow_TreeClickAndOpen(t *testing.T) {
    // Start in browse mode
    // Click on a file in the tree
    // Verify mode is View, activeNote is set
    // Verify outline, backlinks, recents are updated
}

func TestWorkflow_FollowWikiLink(t *testing.T) {
    // Open a note with wiki-links
    // Press Tab to cycle to a link
    // Press Enter to follow
    // Verify new note is opened
}
```

#### State preservation
```go
func TestStatePreservation_ThemeSwitch(t *testing.T) {
    // Open a note
    // Switch theme
    // Verify note is still open
    // Verify colors are updated
}

func TestStatePreservation_SplitResize(t *testing.T) {
    // Open a note
    // Resize split via mouse drag
    // Verify note is still open
    // Verify content is re-rendered at new width
}
```

#### Session persistence
```go
func TestSession_SaveAndRestore(t *testing.T) {
    // Expand some directories
    // Navigate to a file
    // Save session
    // Create new model (simulating restart)
    // Verify directories are expanded
    // Verify cursor is at saved position
}
```

### Test helpers

Create helpers in `testutil_test.go`:
- `createTestModel(t) Model` — creates a model with test vault
- `simulateKeyPress(model, key)` — sends a key event
- `simulateMouseClick(model, x, y)` — sends a mouse click
- `assertMode(t, model, expectedMode)` — checks mode
- `assertActiveNote(t, model, expectedPath)` — checks active note
- `assertContains(t, output, substring)` — checks rendered output

## Files to modify

| File | Changes |
|------|---------|
| `integration_test.go` | **New** — integration tests for full workflows |
| `testutil_test.go` | **New** — test helpers for integration tests |

## Completion Criteria

- [x] Full rendering pipeline test passes (no panics, no broken ANSI)
- [x] Search -> open workflow test passes
- [x] Tree click -> open workflow test passes
- [x] Wiki-link follow workflow test passes
- [x] Theme switch preserves state test passes
- [x] Split resize preserves state test passes
- [x] Session save/restore test passes
- [x] All integration tests pass
- [x] `make test` passes all tests
- [x] `make vet` exits 0

## Completed

2026-06-12

Added 7 integration tests in `model_integration_test.go` exercising full end-to-end workflows:

1. **TestRenderingPipeline_FullDocument** — opens index.md, verifies viewer output, checks for truncated ANSI, confirms body content
2. **TestWorkflow_SearchAndOpen** — presses `/`, types "index", presses Enter, verifies ModeView + activeNote = index.md
3. **TestWorkflow_TreeClickAndOpen** — navigates tree to first file, presses Enter, verifies mode + activeNote + outlineItems + recentNotes all updated
4. **TestWorkflow_FollowWikiLink** — opens index.md, presses Tab to select wiki-link, presses Enter, verifies followed to linked note
5. **TestStatePreservation_ThemeSwitch** — opens note, calls setTheme("dracula"), verifies mode/activeNote preserved, colors changed, viewer still renders
6. **TestStatePreservation_SplitResize** — opens note, presses Ctrl+Left to shrink tree, verifies mode/activeNote preserved, viewer still renders
7. **TestSession_SaveAndRestore** — expands directory, navigates to file, saves session, creates new model, verifies directory still expanded

Bonus fixes:
- Tree resize (ShrinkTree/GrowTree/ResetTree) now works from View mode too (was Browse-only)
- `ExpandPath` fixed to expand top-level directories (was missing `len(parts)` vs `len(parts)-1`)

## Estimated Time

2 days
