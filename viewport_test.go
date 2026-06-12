package main

import (
	"strings"
	"testing"
)

func TestVisibleLength_PlainText(t *testing.T) {
	if n := visibleLength("hello"); n != 5 {
		t.Errorf("visibleLength(hello) = %d, want 5", n)
	}
	if n := visibleLength(""); n != 0 {
		t.Errorf("visibleLength('') = %d, want 0", n)
	}
}

func TestVisibleLength_WithANSI(t *testing.T) {
	s := "\033[1mbold\033[0m text"
	n := visibleLength(s)
	if n != 9 { // "bold text" = 9 runes
		t.Errorf("visibleLength with ANSI = %d, want 9", n)
	}
}

func TestSoftWrap_PlainText_Fits(t *testing.T) {
	lines := softWrap("hello world", 20)
	if len(lines) != 1 || lines[0] != "hello world" {
		t.Errorf("plain text should fit: got %v", lines)
	}
}

func TestSoftWrap_PlainText_Wraps(t *testing.T) {
	lines := softWrap("0123456789abc", 10)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "0123456789" {
		t.Errorf("line[0] = %q, want '0123456789'", lines[0])
	}
	if lines[1] != "abc" {
		t.Errorf("line[1] = %q, want 'abc'", lines[1])
	}
}

func TestSoftWrap_ANSIContent_Fits(t *testing.T) {
	// Simulate a lipgloss-rendered line: visible text fits within width,
	// but ANSI codes inflate raw rune count.
	styled := "\033[38;2;107;114;128m│\033[0m      hello\033[38;2;107;114;128m│\033[0m"
	visLen := visibleLength(styled)
	rawLen := len([]rune(styled))
	if rawLen <= visLen {
		t.Skip("test requires ANSI codes to inflate rune count")
	}

	lines := softWrap(styled, visLen)
	if len(lines) != 1 {
		t.Fatalf("ANSI line fitting visible width should stay as single line, got %d lines", len(lines))
	}
	// Verify no truncated escape sequences
	for _, line := range lines {
		if hasTruncatedANSI(line) {
			t.Errorf("broken ANSI in line: %q", line)
		}
	}
}

func TestSoftWrap_ANSIContent_ExceedsWidth(t *testing.T) {
	// Content that truly exceeds width in visible characters
	styled := "\033[1m0123456789\033[0mabcdef"
	lines := softWrap(styled, 10)
	if len(lines) < 2 {
		t.Fatalf("expected >= 2 lines, got %d", len(lines))
	}
	// No truncated escape sequences in any line
	for _, line := range lines {
		if hasTruncatedANSI(line) {
			t.Errorf("truncated ANSI escape in line: %q", line)
		}
	}
}

func hasTruncatedANSI(s string) bool {
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
			j := i + 2
			for j < len(runes) && runes[j] != 'm' {
				j++
			}
			if j >= len(runes) {
				return true // unterminated escape
			}
			i = j + 1
			continue
		}
		i++
	}
	return false
}

func TestSoftWrap_WidthZero(t *testing.T) {
	lines := softWrap("hello", 0)
	if len(lines) != 1 || lines[0] != "hello" {
		t.Errorf("width=0 should return original: got %v", lines)
	}
}

func TestViewport_SetContent_VisualRegression(t *testing.T) {
	// Full pipeline: styled content → viewport → view should not corrupt ANSI
	vp := newViewport(40, 10)

	// Build a styled line similar to what lipgloss produces
	styled := "\033[1mBold Title\033[0m and some text here that is long enough"
	vp.SetContent(styled)

	output := vp.View()
	if output == "" {
		t.Fatal("viewport view returned empty")
	}

	// Output must not contain truncated escape sequences
	for i, line := range strings.Split(output, "\n") {
		if hasTruncatedANSI(line) {
			t.Errorf("line %d has truncated ANSI escape: %q", i, line)
		}
	}
}
