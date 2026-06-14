package ansiext

import "strings"

// Undercurl wraps text with the undercurl SGR sequence (\033[4:3m).
// Terminals that don't support it (older xterm) render regular underline instead.
func Undercurl(text string) string {
	if text == "" {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\033[4:3m")
	sb.WriteString(text)
	sb.WriteString("\033[4:0m")
	return sb.String()
}

// Overline wraps text with the overline SGR sequence (\033[53m).
func Overline(text string) string {
	if text == "" {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\033[53m")
	sb.WriteString(text)
	sb.WriteString("\033[55m")
	return sb.String()
}
