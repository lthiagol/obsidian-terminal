package main

import (
	"regexp"
	"strings"
)

type viewport struct {
	Width   int
	Height  int
	YOffset int
	lines   []string
}

func newViewport(width, height int) viewport {
	return viewport{
		Width:  width,
		Height: height,
	}
}

func (v *viewport) SetContent(content string) {
	rawLines := strings.Split(content, "\n")
	v.lines = nil
	for _, line := range rawLines {
		v.lines = append(v.lines, softWrap(line, v.Width)...)
	}
	v.clampOffset()
}

func (v viewport) View() string {
	if len(v.lines) == 0 {
		return ""
	}

	end := v.YOffset + v.Height
	if end > len(v.lines) {
		end = len(v.lines)
	}

	if v.YOffset >= len(v.lines) {
		return ""
	}

	return strings.Join(v.lines[v.YOffset:end], "\n")
}

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleLength(s string) int {
	clean := ansiRe.ReplaceAllString(s, "")
	return len([]rune(clean))
}

func softWrap(line string, width int) []string {
	if width <= 0 {
		return []string{line}
	}

	if visibleLength(line) <= width {
		return []string{line}
	}

	// Visible content exceeds width — ANSI-aware hard wrapping.
	// Walk runes, skip ANSI sequences, split at visible-width boundaries.
	var lines []string
	runes := []rune(line)
	var current strings.Builder
	vis := 0

	for i := 0; i < len(runes); {
		if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
			current.WriteRune(runes[i])
			i++
			for i < len(runes) && runes[i] != 'm' {
				current.WriteRune(runes[i])
				i++
			}
			if i < len(runes) {
				current.WriteRune(runes[i])
				i++
			}
			continue
		}

		current.WriteRune(runes[i])
		vis++
		i++

		if vis >= width {
			lines = append(lines, current.String())
			current.Reset()
			vis = 0
		}
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	return lines
}

func (v *viewport) LineUp(n int) {
	v.YOffset -= n
	v.clampOffset()
}

func (v *viewport) LineDown(n int) {
	v.YOffset += n
	v.clampOffset()
}

func (v *viewport) SetYOffset(n int) {
	v.YOffset = n
	v.clampOffset()
}

func (v *viewport) GotoBottom() {
	v.YOffset = len(v.lines) - v.Height
	v.clampOffset()
}

func (v *viewport) HalfViewUp() {
	half := v.Height / 2
	if half < 1 {
		half = 1
	}
	v.LineUp(half)
}

func (v *viewport) HalfViewDown() {
	half := v.Height / 2
	if half < 1 {
		half = 1
	}
	v.LineDown(half)
}

func (v viewport) TotalLineCount() int {
	return len(v.lines)
}

func (v *viewport) clampOffset() {
	if v.YOffset < 0 {
		v.YOffset = 0
	}
	maxOffset := len(v.lines) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if v.YOffset > maxOffset {
		v.YOffset = maxOffset
	}
}
