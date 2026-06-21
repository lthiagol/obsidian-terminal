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

func TestKeyMap_MatchDown(t *testing.T) {
	keys := DefaultKeys()
	if !keys.MatchDown(tea.KeyMsg{Type: tea.KeyDown}) {
		t.Error("KeyDown should match MatchDown")
	}
	if !keys.MatchDown(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) {
		t.Error("'j' should match MatchDown")
	}
	if keys.MatchDown(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}) {
		t.Error("'k' should NOT match MatchDown")
	}
	if keys.MatchDown(tea.KeyMsg{Type: tea.KeyEnter}) {
		t.Error("Enter should NOT match MatchDown")
	}
}

func TestKeyMap_MatchUp(t *testing.T) {
	keys := DefaultKeys()
	if !keys.MatchUp(tea.KeyMsg{Type: tea.KeyUp}) {
		t.Error("KeyUp should match MatchUp")
	}
	if !keys.MatchUp(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}) {
		t.Error("'k' should match MatchUp")
	}
	if keys.MatchUp(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) {
		t.Error("'j' should NOT match MatchUp")
	}
}

func TestKeyMap_MatchLeft(t *testing.T) {
	keys := DefaultKeys()
	if !keys.MatchLeft(tea.KeyMsg{Type: tea.KeyLeft}) {
		t.Error("KeyLeft should match MatchLeft")
	}
	if !keys.MatchLeft(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}) {
		t.Error("'h' should match MatchLeft")
	}
	if keys.MatchLeft(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}) {
		t.Error("'l' should NOT match MatchLeft")
	}
}

func TestKeyMap_MatchRight(t *testing.T) {
	keys := DefaultKeys()
	if !keys.MatchRight(tea.KeyMsg{Type: tea.KeyRight}) {
		t.Error("KeyRight should match MatchRight")
	}
	if !keys.MatchRight(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}) {
		t.Error("'l' should match MatchRight")
	}
	if keys.MatchRight(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}) {
		t.Error("'h' should NOT match MatchRight")
	}
}

