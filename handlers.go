package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lthiagol/obsidian-terminal/internal/search"
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
	case msg.Type == tea.KeyEnter:
		entry := m.fileTree.SelectedEntry()
		if entry != nil {
			if entry.IsDir {
				m.fileTree.toggleExpand()
			} else {
				note, err := LoadNote(m.config.VaultPath, entry.Path)
				if err == nil {
					m.activeNote = note
					m.prevMode = m.mode
					m.mode = ModeView
					m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
				} else {
					m.addToast("Could not load note: "+err.Error(), ToastError)
				}
			}
		}
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.fileTree.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.fileTree.MoveUp()
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
	}
	return m, nil
}

func (m Model) handleViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		m.activeNote = nil
		return m, nil
	case MatchRune(msg, m.keys.Search):
		m.enterSearchMode()
		return m, nil
	case MatchRune(msg, m.keys.Find):
		m.enterFindMode()
		return m, nil
	case MatchRune(msg, m.keys.Help):
		m.enterHelpMode()
		return m, nil
	case msg.Type == tea.KeyTab:
		m.viewer.CycleLink()
		return m, nil
	case msg.Type == tea.KeyEnter:
		if m.viewer.SelectedLinkIndex() >= 0 {
			target := m.viewer.SelectedLinkPath()
			if target != "" {
				resolved := ResolveWikiLink(target, m.vault, m.config.VaultPath)
				if resolved != "" {
					note, err := LoadNote(m.config.VaultPath, resolved)
					if err == nil {
						m.activeNote = note
						m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
					} else {
						m.addToast("Could not load note: "+err.Error(), ToastError)
					}
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
	}
	return m, nil
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m.handleSearchOrFind(msg)
}

func (m Model) handleFindKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m.handleSearchOrFind(msg)
}

func (m Model) handleSearchOrFind(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case msg.Type == tea.KeyBackspace:
		if len(m.searchState.Query()) > 0 {
			m.searchState.SetQuery(m.searchState.Query()[:len(m.searchState.Query())-1])
		}
		return m, nil
	case msg.Type == tea.KeyRunes && len(msg.Runes) > 0:
		m.searchState.SetQuery(m.searchState.Query() + string(msg.Runes))
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.searchState.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.searchState.MoveUp()
		return m, nil
	case msg.Type == tea.KeyEnter:
		result := m.searchState.SelectedResult()
		if result != nil {
			note, err := LoadNote(m.config.VaultPath, result.Path)
			if err == nil {
				m.activeNote = note
				m.prevMode = ModeBrowse
				m.mode = ModeView
				m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
			} else {
				m.addToast("Could not load note: "+err.Error(), ToastError)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.helpScroll++
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		if m.helpScroll > 0 {
			m.helpScroll--
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) enterSearchMode() {
	m.prevMode = m.mode
	m.mode = ModeSearch
	m.searchState = search.NewState(search.Name, m.allPaths, m.searchIndex)
}

func (m *Model) enterFindMode() {
	m.prevMode = m.mode
	m.mode = ModeFind
	m.searchState = search.NewState(search.Content, m.allPaths, m.searchIndex)
}

func (m *Model) enterHelpMode() {
	m.prevMode = m.mode
	m.mode = ModeHelp
	m.helpScroll = 0
}
