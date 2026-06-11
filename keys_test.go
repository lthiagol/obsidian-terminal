package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestKeyDispatch_BothVimAndArrows(t *testing.T) {
	keys := DefaultKeys()

	// j and ↓ both match Down
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, keys.DownRune) {
		t.Error("'j' should match DownRune")
	}
	if !MatchKey(tea.KeyMsg{Type: tea.KeyDown}, keys.Down) {
		t.Error("KeyDown should match Down")
	}

	// k and ↑ both match Up
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}, keys.UpRune) {
		t.Error("'k' should match UpRune")
	}
	if !MatchKey(tea.KeyMsg{Type: tea.KeyUp}, keys.Up) {
		t.Error("KeyUp should match Up")
	}

	// h and ← both match Left
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}, keys.LeftRune) {
		t.Error("'h' should match LeftRune")
	}
	if !MatchKey(tea.KeyMsg{Type: tea.KeyLeft}, keys.Left) {
		t.Error("KeyLeft should match Left")
	}

	// l and → both match Right
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, keys.RightRune) {
		t.Error("'l' should match RightRune")
	}
	if !MatchKey(tea.KeyMsg{Type: tea.KeyRight}, keys.Right) {
		t.Error("KeyRight should match Right")
	}

	// q matches QuitRune
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, keys.QuitRune) {
		t.Error("'q' should match QuitRune")
	}

	// / matches Search
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}, keys.Search) {
		t.Error("'/' should match Search")
	}

	// s matches Find
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}, keys.Find) {
		t.Error("'s' should match Find")
	}

	// ? matches Help
	if !MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, keys.Help) {
		t.Error("'?' should match Help")
	}

	// Non-matching rune
	if MatchRune(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}, keys.QuitRune) {
		t.Error("'x' should NOT match QuitRune")
	}

	// Non-rune key should not match rune
	if MatchRune(tea.KeyMsg{Type: tea.KeyEnter}, keys.QuitRune) {
		t.Error("Enter should NOT match QuitRune")
	}
}
