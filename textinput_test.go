package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleTextInput(t *testing.T) {
	query := "hello"

	// Esc dismisses the input and clears the query.
	if newQuery, dismissed, handled := HandleTextInput(tea.KeyMsg{Type: tea.KeyEsc}, query); !handled || !dismissed || newQuery != "" {
		t.Errorf("Esc: got (%q, %v, %v), want (%q, true, true)", newQuery, dismissed, handled, "")
	}

	// Backspace trims the last rune.
	if newQuery, dismissed, handled := HandleTextInput(tea.KeyMsg{Type: tea.KeyBackspace}, query); !handled || dismissed || newQuery != "hell" {
		t.Errorf("Backspace: got (%q, %v, %v), want (%q, false, true)", newQuery, dismissed, handled, "hell")
	}

	// Backspace on empty query is a handled no-op.
	if newQuery, dismissed, handled := HandleTextInput(tea.KeyMsg{Type: tea.KeyBackspace}, ""); !handled || dismissed || newQuery != "" {
		t.Errorf("Backspace empty: got (%q, %v, %v), want (%q, false, true)", newQuery, dismissed, handled, "")
	}

	// Runes append to the query.
	if newQuery, dismissed, handled := HandleTextInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}}, query); !handled || dismissed || newQuery != "hello!" {
		t.Errorf("Runes: got (%q, %v, %v), want (%q, false, true)", newQuery, dismissed, handled, "hello!")
	}

	// Empty runes event is not handled.
	if newQuery, dismissed, handled := HandleTextInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}}, query); handled || dismissed || newQuery != query {
		t.Errorf("Empty runes: got (%q, %v, %v), want (%q, false, false)", newQuery, dismissed, handled, query)
	}

	// Other keys are not handled.
	if newQuery, dismissed, handled := HandleTextInput(tea.KeyMsg{Type: tea.KeyDown}, query); handled || dismissed || newQuery != query {
		t.Errorf("KeyDown: got (%q, %v, %v), want (%q, false, false)", newQuery, dismissed, handled, query)
	}
}
