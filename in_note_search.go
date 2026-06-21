package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) activateInNoteSearch() (tea.Model, tea.Cmd) {
	if m.activeNote == nil || m.activeNote.Body == "" {
		return m, nil
	}
	m.inNoteSearchActive = true
	m.inNoteSearchQuery = ""
	m.inNoteSearchIdx = 0
	m.inNoteMatches = nil
	return m, nil
}

func (m *Model) updateInNoteSearch(query string) {
	m.inNoteSearchQuery = query
	m.inNoteSearchIdx = 0
	m.inNoteMatches = nil

	if query == "" {
		return
	}

	body := ""
	if m.activeNote != nil {
		body = m.activeNote.Body
	}
	lines := strings.Split(body, "\n")
	queryLower := strings.ToLower(query)
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), queryLower) {
			m.inNoteMatches = append(m.inNoteMatches, i)
		}
	}
}

func (m *Model) cycleInNoteMatch(dir int) {
	if len(m.inNoteMatches) == 0 {
		return
	}
	m.inNoteSearchIdx = (m.inNoteSearchIdx + dir) % len(m.inNoteMatches)
	if m.inNoteSearchIdx < 0 {
		m.inNoteSearchIdx += len(m.inNoteMatches)
	}
	targetLine := m.inNoteMatches[m.inNoteSearchIdx]
	m.viewer.SetScrollPosition(targetLine)
}

func (m *Model) handleInNoteSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyEsc {
		m.inNoteSearchActive = false
		m.inNoteSearchQuery = ""
		m.inNoteMatches = nil
		return m, nil
	}
	if msg.Type == tea.KeyEnter {
		m.inNoteSearchActive = false
		m.inNoteSearchQuery = ""
		m.inNoteMatches = nil
		return m, nil
	}
	if msg.Type == tea.KeyBackspace {
		if len(m.inNoteSearchQuery) > 0 {
			m.updateInNoteSearch(m.inNoteSearchQuery[:len(m.inNoteSearchQuery)-1])
		}
		return m, nil
	}
	if MatchRune(msg, 'n') {
		m.cycleInNoteMatch(1)
		return m, nil
	}
	if MatchRune(msg, 'N') {
		m.cycleInNoteMatch(-1)
		return m, nil
	}
	if len(msg.Runes) > 0 {
		m.updateInNoteSearch(m.inNoteSearchQuery + string(msg.Runes))
		return m, nil
	}
	return m, nil
}

func (m Model) renderInNoteSearch() string {
	if !m.inNoteSearchActive {
		return ""
	}
	width := m.width - m.treeWidth - 6
	if width < 20 {
		width = 20
	}

	var sb strings.Builder
	label := lipgloss.NewStyle().Bold(true).Foreground(m.palette.AccentSecondary).Render("/")
	sb.WriteString(fmt.Sprintf("%s%s_", label, m.inNoteSearchQuery))

	if len(m.inNoteMatches) > 0 {
		info := fmt.Sprintf("  (%d/%d)", m.inNoteSearchIdx+1, len(m.inNoteMatches))
		sb.WriteString(lipgloss.NewStyle().Foreground(m.palette.TextDim).Render(info))
	}

	return lipgloss.NewStyle().Width(width).Render(sb.String())
}
