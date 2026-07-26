# M32 — Modern Terminal Polish

**Status:** ⏳ pending

## Goal

The TUI looks crisp on modern terminals (Kitty, Ghostty, WezTerm, foot, iTerm2). No visual artifacts on resize, wiki-links use undercurl instead of plain underline, and redraws are flicker-free.

## Problem

Three papercuts make the app feel rough on modern terminals:

1. **Wiki-links use plain underline.** Modern terminals support undercurl (`\033[4:3m`) — a wavy underline distinct from regular hyperlink underline. This is a subtle but important visual cue that separates wiki-links from URLs.

2. **Redraw flickers on resize.** When the terminal resizes, Bubble Tea re-renders the full frame. Without synchronized output (DECSET 2026), the terminal may show partial frames, creating a visible flicker/jitter effect. Modern terminals support this protocol natively.

3. **No modern terminal detection or graceful degradation.** We emit SGR sequences blindly. If a terminal doesn't support undercurl, it should fall back to regular underline silently (the standard behavior for SGR 4:3). But we don't test this or document which terminals work.

## What's in / out of scope

| Feature | In scope | Rationale |
|---------|----------|-----------|
| Undercurl for wiki-links | ✅ | One SGR change, huge visual improvement |
| Synchronized output (DECSET 2026) | ✅ | Bubble Tea supports it natively |
| True italic verification | ✅ | Already used via lipgloss, just test+confirm |
| Overline for highlights | ✅ | `\033[53m` as alternative highlight style |
| Terminal capability detection | ❌ | Over-engineering. Terminal silently ignores unsupported SGR |
| Kitty graphics protocol | ❌ | Needs its own milestone (M98 — Image Preview) |
| CSI u keyboard | ❌ | Bubble Tea handles this already |
| Custom cursor shapes | ❌ | Not relevant to a read-only viewer |
| OSC 52 clipboard | ❌ | Read-only, no copy needed |
| OSC 8 hyperlinks | ❌ | Blocked by Bubble Tea limitation, revisit if API changes |

## Design

### Undercurl for wiki-links

Change `renderSegment` for wiki-link segments from `.Underline(true)` to using a custom SGR sequence:

```go
case seg.IsWikiLink:
    s = s.Foreground(style.AccentTertiary)
    // Use undercurl on terminals that support it (Kitty, Ghostty, WezTerm, foot, iTerm2)
    // Falls back to regular underline silently on unsupported terminals
    s = s.Underline(true)  // existing code — for fallback terminals
    // TODO: wrap with undercurl SGR if we build a raw-SGR helper
```

Since lipgloss doesn't expose undercurl natively, we need one of:
- **Option A:** Contribute a `.Undercurl()` method to lipgloss (preferred, but takes time)
- **Option B:** Wrap wiki-link text with raw `\033[4:3m...\033[4:0m` using lipgloss's `Render()` output and string manipulation
- **Option C:** Add a small `ansiext` package that provides undercurl/overline SGR wrappers

**Decision: Option C** — a tiny internal package that wraps text with modern SGR sequences. Keeps the renderer clean and avoids external dependency churn.

```go
// internal/ansiext/ansiext.go
package ansiext

// Undercurl wraps text with the undercurl SGR sequence.
// Terminals that don't support it render regular underline instead.
func Undercurl(text string) string {
    return "\033[4:3m" + text + "\033[4:0m"
}

// Overline wraps text with the overline SGR sequence.
func Overline(text string) string {
    return "\033[53m" + text + "\033[55m"
}
```

### Synchronized output

Bubble Tea supports synchronized output via `tea.WithANSICompressor()`. Check if our `main.go` already uses it. If not, add it. This batches terminal writes into a single frame buffer and flushes atomically.

Additionally, verify that `tea.WithFPS()` is set to a reasonable value (30 FPS for reading, 60 FPS for interaction).

Testing: resize the terminal rapidly during scrolling — no visible tearing or partial frames.

### True italic

We already use `.Italic(true)` in `renderSegment`. Verify that:
- The SGR sequence emitted is `\033[3m` (italic) not `\033[2m` (dim)
- Lipgloss correctly resets italic with `\033[23m`
- Test on multiple terminals to confirm visual appearance

No code changes needed — this is a verification step.

### Overline for highlighted text

Currently `==highlight==` renders as colored text (`.Foreground(style.AccentSecondary)`). Add overline for additional visual distinction:

```go
case seg.Highlight:
    s = s.Foreground(style.AccentSecondary)
    // Overline gives a distinct "highlighted" appearance on modern terminals
    // Falls back to just colored text on unsupported terminals
    return ansiext.Overline(s.Render(seg.Text))
```

## Files to modify

| File | Changes |
|------|---------|
| `internal/ansiext/ansiext.go` | **New file** — `Undercurl()`, `Overline()` helpers |
| `internal/markdown/markdown.go` | Update `renderSegment` wiki-link case to use `ansiext.Undercurl`; update highlight case to use `ansiext.Overline` |
| `main.go` | Add `tea.WithANSICompressor()` if not present; verify FPS setting |
| `internal/markdown/markdown_test.go` | Update render tests to verify undercurl/overline SGR sequences |
| `internal/ansiext/ansiext_test.go` | **New file** — test SGR sequence generation, verify fallback behavior |

## Steps

### 1. ANSI extension package
Create `internal/ansiext/` with `Undercurl(text) string` and `Overline(text) string`. Add unit tests verifying correct SGR sequences and proper reset codes.

### 2. Wiki-link undercurl
Update `renderSegment` to use `ansiext.Undercurl` for wiki-links. Keep the accent color. Test: render a note with wiki-links, verify output contains `\033[4:3m` sequences.

### 3. Highlight overline  
Update `renderSegment` to use `ansiext.Overline` for highlighted text. Keep the accent secondary color.

### 4. Synchronized output
Audit `main.go` for Bubble Tea options. Add `tea.WithANSICompressor()` if missing. Verify `tea.WithFPS()` is set. Test with rapid resize.

### 5. True italic verification
Audit that `.Italic(true)` produces `\033[3m` (not dim). Verify in rendered output. Test on Kitty and Ghostty.

### 6. Visual regression tests
Add gold-string tests: render a note containing wiki-links and highlighted text, verify the output contains the correct modern SGR sequences.

## Completion Criteria

- [ ] Wiki-links render with `\033[4:3m` undercurl in rendered output
- [ ] Highlighted text renders with `\033[53m` overline
- [ ] SGR sequences properly closed with matching reset codes
- [ ] `ansiext` package has unit tests
- [ ] Bubble Tea uses synchronized output (no flicker on resize)
- [ ] True italic confirmed working on Kitty/Ghostty
- [ ] All existing markdown render tests still pass
- [ ] Manual test: open a note with wiki-links on Kitty — links are wavy-underlined
- [ ] Manual test: rapid terminal resize — no flicker
- [ ] `make test` passes all tests
- [ ] `make vet` exits 0
