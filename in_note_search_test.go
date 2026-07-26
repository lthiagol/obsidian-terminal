package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInNoteSearch_Activate(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	if !m.inNoteSearchActive {
		t.Error("expected inNoteSearchActive after / in view mode")
	}
}

func TestInNoteSearch_TypeQuery(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	for _, r := range "Welcome" {
		m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	if !m.inNoteSearchActive {
		t.Fatal("inNoteSearchActive should still be true")
	}
	if m.inNoteSearchQuery != "Welcome" {
		t.Errorf("expected query 'Welcome', got %q", m.inNoteSearchQuery)
	}
	if len(m.inNoteMatches) == 0 {
		t.Error("expected matches for 'Welcome'")
	}
}

func TestInNoteSearch_CycleMatches(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "Welcome" {
		m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	initialIdx := m.inNoteSearchIdx
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if m.inNoteSearchIdx == initialIdx && len(m.inNoteMatches) > 1 {
		t.Error("expected inNoteSearchIdx to change after n")
	}

	if len(m.inNoteMatches) > 1 {
		m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
		if m.inNoteSearchIdx != initialIdx {
			t.Error("expected N to cycle back to initial index")
		}
	}
}

func TestInNoteSearch_EscDismiss(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEsc})

	if m.inNoteSearchActive {
		t.Error("expected inNoteSearchActive to be false after Esc")
	}
	if m.inNoteSearchQuery != "" {
		t.Errorf("expected empty query after Esc, got %q", m.inNoteSearchQuery)
	}
}

func TestInNoteSearch_EmptyNote(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	path := m.activeNote.Path
	note, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		t.Fatalf("LoadNote: %v", err)
	}
	if note.Body == "" {
		t.Skip("skipped: note has empty body")
	}

	m.activeNote = &VaultNote{Path: "empty.md", Title: "Empty", Body: ""}
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	if m.inNoteSearchActive {
		t.Error("inNoteSearchActive should be false for empty note")
	}
}
