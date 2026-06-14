package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if !m.ready {
		return "Loading..."
	}

	if m.width < 60 || m.height < 15 {
		return lipgloss.NewStyle().
			Foreground(m.palette.Warning).
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Terminal too small — please resize")
	}

	if m.quitting {
		return ""
	}

	var rightPanel string
	if m.commandPaletteVisible {
		rightPanel = m.renderCommandPalette()
	} else if m.recentVisible {
		rightPanel = m.renderRecents()
	} else if m.scanErrorsVisible {
		rightPanel = m.renderScanErrors()
	} else if m.vaultState == VaultStateBroken && m.mode != ModeHelp && m.mode != ModeProfilePicker {
		rightPanel = m.renderBrokenVaultScreen()
	} else {
		switch m.mode {
		case ModeSearch:
			rightPanel = m.renderSearch()
		case ModeFind:
			rightPanel = m.renderFind()
		case ModeHelp:
			rightPanel = m.renderHelp()
		case ModeTags:
			rightPanel = m.tagList.View()
		case ModeProfilePicker:
			rightPanel = m.profilePicker.View()
		case ModeView:
			if m.outlineVisible {
				rightPanel = m.renderOutline()
			} else if m.backlinkMode {
				viewerHeight := (m.height - 1) * 7 / 10
				backlinkHeight := m.height - 1 - viewerHeight - 1
				viewerStyle := m.palette.ViewerStyle.Width(m.width - m.treeWidth - 1).Height(viewerHeight)
				backlinkStyle := lipgloss.NewStyle().
					Border(lipgloss.NormalBorder(), true, false, false, false).
					BorderForeground(m.palette.Accent).
					Width(m.width - m.treeWidth - 1).
					Height(backlinkHeight)
				rightPanel = lipgloss.JoinVertical(lipgloss.Left,
					viewerStyle.Render(m.viewer.View()),
					backlinkStyle.Render(m.backlinkPanel.View()),
				)
			} else {
				viewerOutput := m.viewer.View()
				if m.inNoteSearchActive {
					searchBar := m.renderInNoteSearch()
					viewerOutput = lipgloss.JoinVertical(lipgloss.Left, searchBar, "", viewerOutput)
				}
				rightPanel = viewerOutput
			}
		default:
			if m.mode == ModeBrowse && m.previewVisible {
				rightPanel = m.renderPreview()
			} else {
				rightPanel = "Select a file to view"
			}
		}
	}

	treePanel := m.fileTree.View()

	treeStyle := m.palette.TreeStyle.Width(m.treeWidth).Height(m.height - 1)
	viewerStyle := m.palette.ViewerStyle.Width(m.width - m.treeWidth - 1).Height(m.height - 1)

	leftP := treeStyle.Render(treePanel)
	rightP := viewerStyle.Render(rightPanel)

	main := lipgloss.JoinHorizontal(lipgloss.Top, leftP, rightP)

	statusBar := m.renderStatusBar()

	result := lipgloss.JoinVertical(lipgloss.Top, main, statusBar)

	if len(m.toasts) > 0 {
		toastText := m.renderToasts()
		result = lipgloss.JoinVertical(lipgloss.Bottom, result, toastText)
	}

	return result
}

func (m Model) renderSearch() string {
	return m.renderSearchPanel("fuzzy", "results")
}

func (m Model) renderFind() string {
	return m.renderSearchPanel("content", "matches")
}

func (m Model) renderSearchPanel(label, resultLabel string) string {
	var sb strings.Builder
	modeLabel := lipgloss.NewStyle().Bold(true).Foreground(m.palette.AccentSecondary).Render(label)
	sb.WriteString(fmt.Sprintf("%s  %s_  (%d %s)", modeLabel, m.searchState.Query(), m.searchState.ResultCount(), resultLabel))
	sb.WriteString("\n\n")
	sb.WriteString(search.RenderResults(m.searchState, m.width-m.treeWidth-6, m.searchStyle))
	return sb.String()
}

func (m Model) renderBrokenVaultScreen() string {
	width := m.width - m.treeWidth - 4
	if width < 30 {
		width = 30
	}
	padH := 4
	textWidth := width - padH*2
	if textWidth < 1 {
		textWidth = 1
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(m.palette.Error).Render("Vault is inaccessible")
	msg := "The vault directory could not be read. It may have been moved, deleted, or permissions may have changed."
	wrapped := wordWrap(msg, textWidth)

	recovery := []string{
		"r  Retry rescan",
		"P  Switch profile",
		"q  Quit",
	}
	recoveryText := lipgloss.NewStyle().Foreground(m.palette.TextSecondary).Render(strings.Join(recovery, "  │  "))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		wrapped,
		"",
		recoveryText,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.palette.Error).
		Padding(padH/2, padH).
		Width(width).
		Render(content)

	return lipgloss.NewStyle().
		Height(m.height - 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(box)
}

func (m Model) renderScanErrors() string {
	width := m.width - m.treeWidth - 6
	if width < 20 {
		width = 20
	}

	var sb strings.Builder
	title := lipgloss.NewStyle().Bold(true).Foreground(m.palette.Warning).Render(
		fmt.Sprintf("Scan Errors (%d)", len(m.scanErrors)))
	sb.WriteString(title)
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(m.palette.TextDim).Render(strings.Repeat("─", width)))
	sb.WriteString("\n\n")

	for _, err := range m.scanErrors {
		sb.WriteString(lipgloss.NewStyle().Foreground(m.palette.TextSecondary).Render(" • " + err))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	hint := lipgloss.NewStyle().Foreground(m.palette.TextDim).Render("Ctrl+R to rescan  •  Esc to close")
	sb.WriteString(hint)

	return sb.String()
}

func (m *Model) showScanErrors() {
	m.scanErrorsVisible = true
}

func wordWrap(text string, width int) string {
	if width < 1 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	var lines []string
	current := ""
	for _, word := range words {
		if current == "" {
			current = word
		} else if len(current)+1+len(word) <= width {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}
