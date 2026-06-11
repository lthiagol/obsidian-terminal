package main

import "strings"

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

func softWrap(line string, width int) []string {
	if width <= 0 {
		return []string{line}
	}

	runes := []rune(line)
	if len(runes) <= width {
		return []string{line}
	}

	var lines []string
	for len(runes) > 0 {
		if len(runes) <= width {
			lines = append(lines, string(runes))
			break
		}
		lines = append(lines, string(runes[:width]))
		runes = runes[width:]
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
