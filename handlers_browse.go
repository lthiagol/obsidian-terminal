package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleBrowseKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case MatchRune(msg, m.keys.Search):
		m.enterSearchMode()
		return m, nil
	case MatchRune(msg, m.keys.Find):
		m.enterFindMode()
		return m, nil
	case MatchRune(msg, m.keys.Help):
		m.enterHelpMode()
		return m, nil
	case MatchRune(msg, 'T'):
		m.enterTagsMode()
		return m, nil
	case MatchRune(msg, m.keys.ProfileSwitch):
		if len(m.config.Profiles) > 0 {
			m.prevMode = m.mode
			m.mode = ModeProfilePicker
		}
		return m, nil
	case msg.Type == m.keys.ShrinkTree:
		m.adjustTreeWidth(m.treeWidth - 5)
		return m, nil
	case msg.Type == m.keys.GrowTree:
		m.adjustTreeWidth(m.treeWidth + 5)
		return m, nil
	case msg.Type == m.keys.ResetTree:
		m.adjustTreeWidth(m.width / 4)
		return m, nil
	case msg.Type == tea.KeyEnter:
		entry := m.fileTree.SelectedEntry()
		if entry != nil {
			if entry.IsDir {
				m.fileTree.toggleExpand()
			} else {
				m.openNote(entry.Path)
			}
		}
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.fileTree.MoveDown()
		if m.previewVisible {
			m.previewPath = ""
		}
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.fileTree.MoveUp()
		if m.previewVisible {
			m.previewPath = ""
		}
		return m, nil
	case MatchKey(msg, m.keys.Left) || MatchRune(msg, m.keys.LeftRune):
		m.fileTree.collapse()
		return m, nil
	case MatchKey(msg, m.keys.Right) || MatchRune(msg, m.keys.RightRune):
		m.fileTree.expand()
		return m, nil
	case MatchRune(msg, m.keys.TopRune):
		m.fileTree.cursor = 0
		return m, nil
	case MatchRune(msg, m.keys.BottomRune):
		if len(m.fileTree.items) > 0 {
			m.fileTree.cursor = len(m.fileTree.items) - 1
		}
		return m, nil
	case MatchRune(msg, m.keys.PreviewToggle):
		m.previewVisible = !m.previewVisible
		if m.previewVisible {
			m.previewPath = ""
			m.addToast("Preview on", ToastInfo)
		} else {
			m.addToast("Preview off", ToastInfo)
		}
		return m, nil
	case MatchRune(msg, m.keys.PinRune):
		entry := m.fileTree.SelectedEntry()
		if entry != nil && !entry.IsDir {
			m.togglePin(entry.Path)
		}
		return m, nil
	case MatchKey(msg, m.keys.CyclePinPrev):
		m.cyclePinnedPrev()
		return m, nil
	case MatchKey(msg, m.keys.CyclePinNext):
		m.cyclePinnedNext()
		return m, nil
	}
	return m, nil
}
