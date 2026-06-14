package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) buildDailyNotePath() string {
	now := time.Now()
	dateStr := now.Format(m.config.DailyNotesFormat)
	return filepath.Join(m.config.DailyNotesDir, dateStr+".md")
}

func (m *Model) openDailyNote() {
	path := m.buildDailyNotePath()
	note, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		dateStr := time.Now().Format(m.config.DailyNotesFormat)
		note = &VaultNote{
			Path:  path,
			Title: "Daily: " + dateStr,
			Body:  "",
		}
		if m.activeNote != nil && m.activeNote.Path != path {
			m.history = append(m.history, m.activeNote.Path)
			m.historyForward = nil
		}
		m.applyNote(note, navUser)
		return
	}
	m.openNote(note.Path)
}

func (m *Model) addRecentNote(path string) {
	if path == "" {
		return
	}

	for i, recent := range m.recentNotes {
		if recent == path {
			m.recentNotes = append(m.recentNotes[:i], m.recentNotes[i+1:]...)
			break
		}
	}

	m.recentNotes = append([]string{path}, m.recentNotes...)

	if len(m.recentNotes) > 50 {
		m.recentNotes = m.recentNotes[:50]
	}
}

func (m *Model) toggleRecents() {
	if m.recentVisible {
		m.recentVisible = false
	} else {
		m.recentVisible = true
		m.recentCursor = 0
	}
}

func (m *Model) openRecentNote(index int) {
	if index < 0 || index >= len(m.recentNotes) {
		return
	}

	path := m.recentNotes[index]
	_, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		m.addToast("Failed to load recent note: "+err.Error(), ToastError)
		m.recentNotes = append(m.recentNotes[:index], m.recentNotes[index+1:]...)
		if m.recentCursor >= len(m.recentNotes) {
			m.recentCursor = len(m.recentNotes) - 1
		}
		return
	}

	m.openNote(path)
	m.recentVisible = false
}

func (m Model) renderRecents() string {
	if len(m.recentNotes) == 0 {
		return lipgloss.NewStyle().
			Foreground(m.palette.TextMuted).
			Render("  No recent notes")
	}

	var sb strings.Builder
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.palette.Accent).
		Render(fmt.Sprintf("  Recent Notes (%d)", len(m.recentNotes)))
	sb.WriteString(header)
	sb.WriteString("\n")

	for i, path := range m.recentNotes {
		line := fmt.Sprintf("  %s", path)

		if i == m.recentCursor {
			line = lipgloss.NewStyle().
				Background(m.palette.Accent).
				Foreground(m.palette.SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(m.palette.TextSecondary).
				Render(line)
		}

		sb.WriteString(line)
		if i < len(m.recentNotes)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (m Model) handleRecentsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.recentVisible = false
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		if m.recentCursor < len(m.recentNotes)-1 {
			m.recentCursor++
		}
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		if m.recentCursor > 0 {
			m.recentCursor--
		}
		return m, nil
	case msg.Type == tea.KeyEnter:
		m.openRecentNote(m.recentCursor)
		return m, nil
	}
	return m, nil
}
