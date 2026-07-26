package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m.handleSearchOrFind(msg)
}

func (m Model) handleFindKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m.handleSearchOrFind(msg)
}

func (m Model) handleSearchOrFind(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if newQuery, dismissed, handled := HandleTextInput(msg, m.searchState.Query()); handled {
		if dismissed {
			m.mode = m.prevMode
		} else {
			m.searchState.SetQuery(newQuery)
		}
		return m, nil
	}
	switch {
	case m.keys.MatchDown(msg):
		m.searchState.MoveDown()
		return m, nil
	case m.keys.MatchUp(msg):
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
	case m.keys.MatchDown(msg):
		m.helpScroll++
		return m, nil
	case m.keys.MatchUp(msg):
		if m.helpScroll > 0 {
			m.helpScroll--
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleBacklinkKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc || MatchRune(msg, 'b'):
		m.backlinkMode = false
		return m, nil
	case m.keys.MatchDown(msg):
		m.backlinkPanel.MoveDown()
		return m, nil
	case m.keys.MatchUp(msg):
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
	case m.keys.MatchDown(msg):
		m.tagList.MoveDown()
		return m, nil
	case m.keys.MatchUp(msg):
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

func (m Model) handleCommandPaletteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if newQuery, dismissed, handled := HandleTextInput(msg, m.commandPaletteQuery); handled {
		if dismissed {
			m.commandPaletteVisible = false
		} else {
			m.commandPaletteQuery = newQuery
			m.commandPaletteSearch()
		}
		return m, nil
	}
	switch {
	case m.keys.MatchDown(msg):
		if m.commandPaletteCursor < len(m.commandPaletteResults)-1 {
			m.commandPaletteCursor++
		}
		return m, nil
	case m.keys.MatchUp(msg):
		if m.commandPaletteCursor > 0 {
			m.commandPaletteCursor--
		}
		return m, nil
	case msg.Type == tea.KeyEnter:
		return m.executeCommand(m.commandPaletteCursor)
	}
	return m, nil
}

func (m Model) handleProfilePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		// If we have no vault loaded, Esc quits the app
		if m.vault == nil {
			m.quitting = true
			return m, tea.Quit
		}
		// Otherwise return to previous mode
		m.mode = m.prevMode
		return m, nil
	case m.keys.MatchDown(msg):
		m.profilePicker.MoveDown()
		return m, nil
	case m.keys.MatchUp(msg):
		m.profilePicker.MoveUp()
		return m, nil
	case msg.Type == tea.KeyEnter:
		profileName := m.profilePicker.Selected()
		if profileName != "" {
			m.switchToProfile(profileName)
		}
		return m, nil
	}
	return m, nil
}
