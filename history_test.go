package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHistory_Navigation(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel: %v", m.err)
	}

	noteA := "notes/meeting.md"
	noteB := "notes/frontmatter-test.md"
	noteC := "projects/api-design.md"
	noteD := "notes/callouts.md"

	tests := []struct {
		name string
		step func(m *Model, t *testing.T)
	}{
		{
			name: "T1: Open A, B, C then [ -> activeNote=B, history=[A], forward=[C]",
			step: func(m *Model, t *testing.T) {
				m.openNote(noteA)
				m.openNote(noteB)
				m.openNote(noteC)

				if m.mode != ModeView {
					t.Fatalf("expected ModeView, got %v", m.mode)
				}
				if m.activeNote == nil || m.activeNote.Path != noteC {
					t.Fatalf("expected activeNote=%s, got %v", noteC, pathOrNil(m.activeNote))
				}

				model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'['}})
				updated := model.(Model)
				*m = updated

				if updated.activeNote == nil || updated.activeNote.Path != noteB {
					t.Fatalf("after [, expected activeNote=%s, got %v", noteB, pathOrNil(updated.activeNote))
				}
				if len(updated.history) != 1 || updated.history[0] != noteA {
					t.Fatalf("after [, expected history=[%s], got %v", noteA, updated.history)
				}
				if len(updated.historyForward) != 1 || updated.historyForward[0] != noteC {
					t.Fatalf("after [, expected forward=[%s], got %v", noteC, updated.historyForward)
				}
			},
		},
		{
			name: "T2: After T1, [ again -> activeNote=A, forward=[C,B] (stack order)",
			step: func(m *Model, t *testing.T) {
				model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'['}})
				updated := model.(Model)
				*m = updated

				if updated.activeNote == nil || updated.activeNote.Path != noteA {
					t.Fatalf("after second [, expected activeNote=%s, got %v", noteA, pathOrNil(updated.activeNote))
				}
				if len(updated.history) != 0 {
					t.Fatalf("after second [, expected empty history, got %v", updated.history)
				}
				if len(updated.historyForward) != 2 {
					t.Fatalf("after second [, expected 2-item forward stack, got %v", updated.historyForward)
				}
				// B was added after C, so B is at the end (popped first on forward)
				if updated.historyForward[0] != noteC || updated.historyForward[1] != noteB {
					t.Fatalf("after second [, expected forward=[%s,%s], got %v", noteC, noteB, updated.historyForward)
				}
			},
		},
		{
			name: "T3: After T2, ] -> activeNote=B",
			step: func(m *Model, t *testing.T) {
				model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{']'}})
				updated := model.(Model)
				*m = updated

				if updated.activeNote == nil || updated.activeNote.Path != noteB {
					t.Fatalf("after ], expected activeNote=%s, got %v", noteB, pathOrNil(updated.activeNote))
				}
				if len(updated.historyForward) != 1 || updated.historyForward[0] != noteC {
					t.Fatalf("after ], expected forward=[%s], got %v", noteC, updated.historyForward)
				}
				if len(updated.history) != 1 || updated.history[0] != noteA {
					t.Fatalf("after ], expected history=[%s], got %v", noteA, updated.history)
				}
			},
		},
		{
			name: "T4: Open D from tree after back -> forward cleared",
			step: func(m *Model, t *testing.T) {
				model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'['}})
				updated := model.(Model)
				*m = updated

				if len(updated.historyForward) == 0 {
					t.Fatal("expected non-empty forward stack before opening new note")
				}

				m.openNote(noteD)

				if m.activeNote == nil || m.activeNote.Path != noteD {
					t.Fatalf("expected activeNote=%s, got %v", noteD, pathOrNil(m.activeNote))
				}
				if len(m.historyForward) != 0 {
					t.Fatalf("expected forward stack cleared, got %v", m.historyForward)
				}
			},
		},
		{
			name: "T5: Ctrl+O in view mode -> same as [",
			step: func(m *Model, t *testing.T) {
				m.openNote(noteA)
				m.openNote(noteC)

				model, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
				updated := model.(Model)
				*m = updated

				if updated.activeNote == nil || updated.activeNote.Path != noteA {
					t.Fatalf("after Ctrl+O, expected activeNote=%s, got %v", noteA, pathOrNil(updated.activeNote))
				}
			},
		},
		{
			name: "T6: Rescan with active note -> history unchanged",
			step: func(m *Model, t *testing.T) {
				m.openNote(noteA)
				m.openNote(noteB)
				m.openNote(noteC)

				originalHistoryLen := len(m.history)
				originalForwardLen := len(m.historyForward)

				model, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
				updated := model.(Model)
				*m = updated

				if len(updated.history) != originalHistoryLen {
					t.Fatalf("expected history length %d after rescan, got %d", originalHistoryLen, len(updated.history))
				}
				if len(updated.historyForward) != originalForwardLen {
					t.Fatalf("expected forward length %d after rescan, got %d", originalForwardLen, len(updated.historyForward))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.step(&m, t)
		})
	}
}

func pathOrNil(n *VaultNote) string {
	if n == nil {
		return "<nil>"
	}
	return n.Path
}
