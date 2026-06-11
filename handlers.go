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
	case MatchRune(msg, 'T'):
		m.enterTagsMode()
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
	if m.backlinkMode {
		return m.handleBacklinkKey(msg)
	}

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
	case MatchRune(msg, 'b'):
		if m.backlinkPanel.Count() > 0 {
			m.backlinkMode = true
		}
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
			m.openNote(result.Path)
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

func (m *Model) openNote(path string) {
	note, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		m.addToast("Could not load note: "+err.Error(), ToastError)
		return
	}
	m.activeNote = note
	m.prevMode = m.mode
	m.mode = ModeView
	m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
	m.backlinkPanel = NewBacklinkPanel(note.Path, m.backlinkIndex)
	m.backlinkMode = false
}

func (m Model) handleBacklinkKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc || MatchRune(msg, 'b'):
		m.backlinkMode = false
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.backlinkPanel.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.backlinkPanel.MoveUp()
		return m, nil
	case msg.Type == tea.KeyEnter:
		path := m.backlinkPanel.SelectedPath()
		if path != "" {
			m.openNote(path)
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleTagsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.tagList.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.tagList.MoveUp()
		return m, nil
	case msg.Type == tea.KeyEnter:
		tag := m.tagList.SelectedTag()
		if tag != "" {
			files := m.tagList.SelectedFiles()
			pathSet := make(map[string]bool)
			for _, f := range files {
				pathSet[f] = true
			}
			m.fileTree.ApplyPathFilter(pathSet)
			m.tagFilter = tag
			m.mode = ModeBrowse
			m.addToast("Filtered by #"+tag, ToastInfo)
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) enterTagsMode() {
	m.prevMode = m.mode
	m.mode = ModeTags
	m.tagList = NewTagList(m.tagIndex)
}
