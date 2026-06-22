package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderStatusBar() string {
	modeColor := m.palette.ModeColor(m.mode)
	modeBadge := lipgloss.NewStyle().
		Background(modeColor).
		Foreground(m.palette.SelectionText).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s ", m.mode.String()))

	var info string
	switch m.mode {
	case ModeBrowse:
		info = fmt.Sprintf("%d files", countFiles(m.vault))
		if m.vaultState == VaultStateBroken {
			info += " BROKEN"
		} else if m.vaultState == VaultStatePartial {
			info += fmt.Sprintf(" (%d scan errors)", len(m.scanErrors))
		}
	case ModeView:
		if m.activeNote != nil {
			info = truncatePath(m.activeNote.Path, m.width-60)
			if m.viewer.SelectedLinkIndex() >= 0 {
				info += " → " + m.viewer.SelectedLinkPath()
			}
			if m.activePinnedIdx >= 0 && m.activePinnedIdx < len(m.pinnedNotes) &&
				m.pinnedNotes[m.activePinnedIdx].Path == m.activeNote.Path {
				info += " 📌"
			}
		}
		if m.vaultState == VaultStateBroken {
			info += " BROKEN"
		}
	case ModeSearch, ModeFind:
		info = m.searchState.Query()
	case ModeHelp:
		info = "j/k scroll | Esc back"
	}

	midSection := lipgloss.NewStyle().Foreground(m.palette.TextSecondary).Padding(0, 1).Render(info)

	hints := modeHints(m.mode)
	hintSection := lipgloss.NewStyle().Foreground(m.palette.TextDim).Padding(0, 1).Render(hints)

	modeWidth := lipgloss.Width(modeBadge)
	midWidth := max(0, m.width-modeWidth-lipgloss.Width(hintSection)-4)

	fullBar := lipgloss.JoinHorizontal(lipgloss.Center,
		modeBadge,
		lipgloss.NewStyle().Width(midWidth).Render(midSection),
		hintSection,
	)

	return m.palette.StatusStyle.Width(m.width).Render(fullBar)
}

func modeHints(mode Mode) string {
	switch mode {
	case ModeBrowse:
		return "/ search | v preview | Enter open | ? help | q quit"
	case ModeView:
		return "h back | j/k scroll | Tab link | / search | ? help"
	case ModeSearch:
		return "type filter | Enter open | Esc cancel"
	case ModeFind:
		return "type search | Enter open | Esc cancel"
	case ModeHelp:
		return "j/k scroll | Esc back"
	case ModeGraph:
		return "j/k move | Enter open | l scope | f focus | r refresh | Esc back"
	default:
		return ""
	}
}

func truncatePath(path string, maxLen int) string {
	if maxLen < 5 {
		return "..."
	}
	if len(path) <= maxLen {
		return path
	}
	return ".../" + path[len(path)-maxLen+4:]
}
