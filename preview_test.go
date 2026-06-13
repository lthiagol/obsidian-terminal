package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPreview_Toggle(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = modelFromInterface(model)

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if !m.previewVisible {
		t.Error("expected previewVisible after v")
	}
	if len(m.toasts) == 0 || !strings.Contains(m.toasts[0].Message, "Preview on") {
		t.Error("expected 'Preview on' toast")
	}

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if m.previewVisible {
		t.Error("expected previewVisible false after second v")
	}
}

func TestPreview_ShowsNoteContent(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = modelFromInterface(model)

	// Navigate to a file, then go back to browse
	m = navigateToFirstFile(t, &model)
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEsc})

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if !m.previewVisible {
		t.Fatal("expected preview to be on")
	}

	output := m.renderPreview()
	if strings.Contains(output, "select a markdown file") || strings.Contains(output, "select a markdown") {
		t.Error("expected note content, got directory placeholder")
	}
}

func TestPreview_DirPlaceholder(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = modelFromInterface(model)

	entry := m.fileTree.SelectedEntry()
	output := m.renderPreview()
	if entry != nil && entry.IsDir {
		if !strings.Contains(output, "select a markdown file") {
			t.Error("expected dir placeholder for directory entry")
		}
	} else if entry != nil {
		if strings.Contains(output, "select a markdown file") {
			t.Error("unexpected placeholder for a file entry")
		}
	}
}

func TestPreview_EnterStillOpensNote(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = modelFromInterface(model)

	m = navigateToFirstFile(t, &model)
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if m.mode != ModeView {
		t.Errorf("expected ModeView after Enter with preview on, got %v", m.mode)
	}
	if m.activeNote == nil {
		t.Error("expected activeNote after Enter")
	}
}
