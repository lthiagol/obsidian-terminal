package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

type noteNavKind int

const (
	navUser noteNavKind = iota // user navigation: push history, clear forward
	navHistory                 // back/forward: stacks already updated
	navReload                  // rescan refresh: no history changes
)

func (m *Model) loadNote(path string, kind noteNavKind) {
	if kind == navUser {
		if m.activeNote != nil && m.activeNote.Path != path {
			m.history = append(m.history, m.activeNote.Path)
			m.historyForward = nil
		}
	}
	note, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		m.addToast("Could not load note: "+err.Error(), ToastError)
		return
	}
	m.applyNote(note, kind)
}

func (m *Model) applyNote(note *VaultNote, kind noteNavKind) {
	m.activeNote = note
	m.prevMode = m.mode
	m.mode = ModeView

	m.viewer.SetEmbedResolver(func(target, heading string) (string, error) {
		if m.vault == nil {
			return "", fmt.Errorf("vault not available")
		}
		resolved := ResolveWikiLink(target, m.vault, m.config.VaultPath)
		if resolved == "" {
			return "", fmt.Errorf("not found: %s", target)
		}
		note, err := LoadNote(m.config.VaultPath, resolved)
		if err != nil {
			return "", err
		}
		if heading != "" {
			return extractSection(note.RawBody, heading), nil
		}
		return note.Body, nil
	})

	m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
	m.backlinkPanel = NewBacklinkPanel(note.Path, m.backlinkIndex, m.palette)
	m.backlinkMode = false
	m.buildOutline()

	if kind != navReload {
		m.addRecentNote(note.Path)
	}
}

func (m *Model) openNote(path string) {
	m.loadNote(path, navUser)
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
	m.tagList = NewTagList(m.tagIndex, m.palette)
}

func (m Model) handleCommandPaletteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.commandPaletteVisible = false
		return m, nil
	case msg.Type == tea.KeyBackspace:
		if len(m.commandPaletteQuery) > 0 {
			m.commandPaletteQuery = m.commandPaletteQuery[:len(m.commandPaletteQuery)-1]
			m.commandPaletteSearch()
		}
		return m, nil
	case msg.Type == tea.KeyRunes && len(msg.Runes) > 0:
		m.commandPaletteQuery += string(msg.Runes)
		m.commandPaletteSearch()
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		if m.commandPaletteCursor < len(m.commandPaletteResults)-1 {
			m.commandPaletteCursor++
		}
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
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
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.profilePicker.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
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

func (m *Model) switchToProfile(profileName string) {
	profile, ok := m.config.Profiles[profileName]
	if !ok {
		m.addToast("Profile not found: "+profileName, ToastError)
		return
	}

	// Apply profile settings
	if profile.Path != "" {
		m.config.VaultPath = profile.Path
	}
	if profile.Theme != "" {
		m.config.Theme = profile.Theme
		m.setTheme(profile.Theme)
	}
	if len(profile.SkipDirs) > 0 {
		m.config.SkipDirs = profile.SkipDirs
	}

	// Rescan vault with new settings
	m.rescanVault()
	m.mode = ModeBrowse
	m.addToast("Switched to profile: "+profileName, ToastInfo)
}

func (m *Model) setTheme(themeName string) {
	palette, err := lookupPalette(themeName)
	if err != nil {
		return
	}
	m.palette = palette
	m.viewer.renderStyle = markdownStyleFrom(palette, m.config.LineSpacing)
	m.searchStyle = searchStyleFrom(palette)
	m.fileTree.SetPalette(palette)
}

func (m *Model) goBackHistory() {
	if len(m.history) == 0 {
		return
	}
	prev := m.history[len(m.history)-1]
	m.history = m.history[:len(m.history)-1]
	if m.activeNote != nil {
		m.historyForward = append(m.historyForward, m.activeNote.Path)
	}
	m.loadNote(prev, navHistory)
}

func (m *Model) goForwardHistory() {
	if len(m.historyForward) == 0 {
		return
	}
	next := m.historyForward[len(m.historyForward)-1]
	m.historyForward = m.historyForward[:len(m.historyForward)-1]
	if m.activeNote != nil {
		m.history = append(m.history, m.activeNote.Path)
	}
	m.loadNote(next, navHistory)
}

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
