package main

import tea "github.com/charmbracelet/bubbletea"

// HandleTextInput processes Esc/Backspace/Runes for an in-TUI text input field.
// Returns the updated query, whether Esc was pressed (caller should exit the
// input mode), and whether the key was consumed (caller should return early).
// The caller is responsible for any after behavior (re-search, update display,
// etc.) using the returned newQuery.
func HandleTextInput(msg tea.KeyMsg, query string) (newQuery string, dismissed bool, handled bool) {
	switch msg.Type {
	case tea.KeyEsc:
		return "", true, true
	case tea.KeyBackspace:
		if len(query) > 0 {
			return query[:len(query)-1], false, true
		}
		return query, false, true
	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			return query + string(msg.Runes), false, true
		}
	}
	return query, false, false
}
