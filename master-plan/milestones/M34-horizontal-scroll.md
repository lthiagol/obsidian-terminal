# M34 — Horizontal Scroll for Viewer

**Status:** ⏳ pending

## Goal

Viewer supports horizontal scrolling for content that exceeds viewport width (long code lines, wide frontmatter values, URLs). When scrolling horizontally, soft-wrap is temporarily disabled.

## Problem

The viewport currently has `YOffset` (vertical scroll) but no `XOffset` (horizontal scroll). Content that exceeds the viewport width is soft-wrapped, which works for prose but not for code blocks, long URLs, or wide frontmatter values where wrapping is undesirable.

## Scope

- Add `XOffset int` to viewport struct
- ANSI-aware line clipping (`clipLineANSI`) that preserves escape sequences
- `ScrollLeft(n)`, `ScrollRight(n)`, `ScrollReset()` viewport methods
- Keybindings: `Shift+Left/Right` (5 cols), `0` resets to column 0
- When `XOffset > 0`, soft-wrap is disabled (lines extend full width)
- Status bar indicator: `Col: 12/80`

## Not in scope

- Horizontal scroll for the tree panel (tree width is configurable)
- Image preview (M98)
- Touch scrolling

## Completion Criteria

- [ ] Viewport has XOffset field with clamp
- [ ] `clipLineANSI` preserves ANSI sequences across clip boundaries
- [ ] `Shift+Left/Right` scrolls viewer horizontally
- [ ] `0` resets horizontal scroll position
- [ ] Status bar shows column indicator when XOffset > 0
- [ ] All existing tests pass
