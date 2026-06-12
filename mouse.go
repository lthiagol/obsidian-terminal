package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type clickState struct {
	lastTime time.Time
	lastX    int
	lastY    int
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	splitZone := m.treeWidth
	// Detect drag on split boundary
	if m.dragSplit {
		if msg.Action == tea.MouseActionRelease {
			m.dragSplit = false
			return m, nil
		}
		if msg.Action == tea.MouseActionMotion {
			m.adjustTreeWidth(msg.X)
			return m, nil
		}
	}
	if msg.Action == tea.MouseActionPress && abs(msg.X-splitZone) <= 2 {
		m.dragSplit = true
		return m, nil
	}
	if msg.X < m.treeWidth {
		return m.handleTreeMouse(msg)
	}
	return m.handleRightPanelMouse(msg)
}

func (m Model) handleRightPanelMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeView:
		return m.handleViewerMouse(msg)
	case ModeSearch, ModeFind:
		return m.handleSearchMouse(msg)
	case ModeHelp:
		return m.handleHelpMouse(msg)
	}
	return m, nil
}

func (m Model) handleTreeMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			if m.isDoubleClick(msg.X, msg.Y) {
				return m.openTreeItem()
			}
			m.recordClick(msg.X, msg.Y)
			m.fileTree.MoveToY(msg.Y)
		}
	case tea.MouseButtonWheelUp:
		for i := 0; i < 3; i++ {
			m.fileTree.MoveUp()
		}
	case tea.MouseButtonWheelDown:
		for i := 0; i < 3; i++ {
			m.fileTree.MoveDown()
		}
	}
	return m, nil
}

func (m Model) handleViewerMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		m.viewer.ScrollUp(3)
	case tea.MouseButtonWheelDown:
		m.viewer.ScrollDown(3)
	}
	return m, nil
}

func (m Model) handleSearchMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			idx := msg.Y - 2
			if idx >= 0 {
				m.searchState.SetSelected(idx)
			}
			if m.isDoubleClick(msg.X, msg.Y) {
				return m.openSearchResult()
			}
			m.recordClick(msg.X, msg.Y)
		}
	case tea.MouseButtonWheelUp:
		for i := 0; i < 3; i++ {
			m.searchState.MoveUp()
		}
	case tea.MouseButtonWheelDown:
		for i := 0; i < 3; i++ {
			m.searchState.MoveDown()
		}
	}
	return m, nil
}

func (m Model) handleHelpMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.helpScroll > 0 {
			m.helpScroll -= 3
			if m.helpScroll < 0 {
				m.helpScroll = 0
			}
		}
	case tea.MouseButtonWheelDown:
		m.helpScroll += 3
	}
	return m, nil
}

func (m Model) openTreeItem() (tea.Model, tea.Cmd) {
	entry := m.fileTree.SelectedEntry()
	if entry == nil {
		return m, nil
	}
	if entry.IsDir {
		m.fileTree.toggleExpand()
		return m, nil
	}
	note, err := LoadNote(m.config.VaultPath, entry.Path)
	if err != nil {
		m.addToast("Could not load note: "+err.Error(), ToastError)
		return m, nil
	}
	m.activeNote = note
	m.prevMode = m.mode
	m.mode = ModeView
	m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
	return m, nil
}

func (m Model) openSearchResult() (tea.Model, tea.Cmd) {
	result := m.searchState.SelectedResult()
	if result == nil {
		return m, nil
	}
	note, err := LoadNote(m.config.VaultPath, result.Path)
	if err != nil {
		m.addToast("Could not load note: "+err.Error(), ToastError)
		return m, nil
	}
	m.activeNote = note
	m.prevMode = ModeBrowse
	m.mode = ModeView
	m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
	return m, nil
}

func (m *Model) isDoubleClick(x, y int) bool {
	if m.click == nil {
		return false
	}
	elapsed := time.Since(m.click.lastTime)
	dx := abs(x - m.click.lastX)
	dy := abs(y - m.click.lastY)
	return elapsed <= 500*time.Millisecond && dx <= 1 && dy <= 1
}

func (m *Model) recordClick(x, y int) {
	if m.click == nil {
		m.click = &clickState{}
	}
	m.click.lastTime = time.Now()
	m.click.lastX = x
	m.click.lastY = y
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
