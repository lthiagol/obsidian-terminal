package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.backlinkMode {
		return m.handleBacklinkKey(msg)
	}

	switch {
	case msg.Type == m.keys.ShrinkTree:
		m.adjustTreeWidth(m.treeWidth - 5)
		if m.activeNote != nil {
			m.viewer.SetContent(m.activeNote.Body, m.width-m.treeWidth-2)
		}
		return m, nil
	case msg.Type == m.keys.GrowTree:
		m.adjustTreeWidth(m.treeWidth + 5)
		if m.activeNote != nil {
			m.viewer.SetContent(m.activeNote.Body, m.width-m.treeWidth-2)
		}
		return m, nil
	case msg.Type == m.keys.ResetTree:
		m.adjustTreeWidth(m.width / 4)
		if m.activeNote != nil {
			m.viewer.SetContent(m.activeNote.Body, m.width-m.treeWidth-2)
		}
		return m, nil
	case msg.Type == tea.KeyEsc:
		if m.inNoteSearchActive {
			m.inNoteSearchActive = false
			m.inNoteSearchQuery = ""
			m.inNoteMatches = nil
			return m, nil
		}
		m.mode = m.prevMode
		m.activeNote = nil
		return m, nil
	case MatchRune(msg, m.keys.Search):
		return m.activateInNoteSearch()
	case MatchRune(msg, 'n'):
		if m.inNoteSearchActive {
			m.cycleInNoteMatch(1)
			return m, nil
		}
	case MatchRune(msg, 'N'):
		if m.inNoteSearchActive {
			m.cycleInNoteMatch(-1)
			return m, nil
		}
	case MatchRune(msg, m.keys.Find):
		m.enterFindMode()
		return m, nil
	case MatchRune(msg, m.keys.Help):
		m.enterHelpMode()
		return m, nil
	case MatchRune(msg, 'b'):
		if m.backlinkPanel.Count() > 0 {
			m.backlinkMode = true
		}
		return m, nil
	case msg.Type == tea.KeyTab:
		m.viewer.CycleLink()
		return m, nil
	case MatchRune(msg, '['):
		m.goBackHistory()
		return m, nil
	case MatchRune(msg, ']'):
		m.goForwardHistory()
		return m, nil
	case msg.Type == tea.KeyEnter:
		if m.viewer.SelectedLinkIndex() >= 0 {
			target := m.viewer.SelectedLinkPath()
			if target != "" && m.vault != nil {
				resolved := ResolveWikiLink(target, m.vault, m.config.VaultPath)
				if resolved != "" {
					m.openNote(resolved)
				}
			}
		}
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.viewer.ScrollDown(1)
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.viewer.ScrollUp(1)
		return m, nil
	case MatchRune(msg, m.keys.TopRune):
		m.viewer.ScrollTop()
		return m, nil
	case MatchRune(msg, m.keys.BottomRune):
		m.viewer.ScrollBottom()
		return m, nil
	case msg.Type == tea.KeyPgUp:
		m.viewer.ScrollHalfPageUp()
		return m, nil
	case msg.Type == tea.KeyPgDown:
		m.viewer.ScrollHalfPageDown()
		return m, nil
	case MatchRune(msg, m.keys.PinRune):
		if m.activeNote != nil {
			m.togglePin(m.activeNote.Path)
		}
		return m, nil
	case MatchKey(msg, m.keys.CyclePinPrev):
		m.cyclePinnedPrev()
		return m, nil
	case MatchKey(msg, m.keys.CyclePinNext):
		m.cyclePinnedNext()
		return m, nil
	case MatchRune(msg, m.keys.Outline):
		if m.outlineVisible {
			m.outlineVisible = false
		} else {
			m.buildOutline()
			m.outlineVisible = true
		}
		return m, nil
	}
	return m, nil
}
