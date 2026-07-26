package main

import (
	"os"
	"path/filepath"
)

func (m *Model) togglePin(path string) {
	if path == "" {
		return
	}

	for i, pin := range m.pinnedNotes {
		if pin.Path == path {
			m.pinnedNotes = append(m.pinnedNotes[:i], m.pinnedNotes[i+1:]...)
			if m.activePinnedIdx >= len(m.pinnedNotes) {
				m.activePinnedIdx = len(m.pinnedNotes) - 1
			}
			m.addToast("Unpinned note", ToastInfo)
			return
		}
	}

	scrollY := 0
	if m.activeNote != nil && m.activeNote.Path == path {
		scrollY = m.viewer.GetScrollPosition()
	}

	m.pinnedNotes = append(m.pinnedNotes, PinnedNote{Path: path, ScrollY: scrollY})
	m.addToast("Pinned note", ToastInfo)
}

func (m *Model) openPinnedNote(index int) {
	if index < 0 || index >= len(m.pinnedNotes) {
		return
	}

	pin := m.pinnedNotes[index]

	_, err := LoadNote(m.config.VaultPath, pin.Path)
	if err != nil {
		m.addToast("Pinned note deleted: "+err.Error(), ToastError)
		m.pinnedNotes = append(m.pinnedNotes[:index], m.pinnedNotes[index+1:]...)
		if m.activePinnedIdx >= len(m.pinnedNotes) {
			m.activePinnedIdx = len(m.pinnedNotes) - 1
		}
		return
	}

	m.openNote(pin.Path)
	m.viewer.SetScrollPosition(pin.ScrollY)
	m.activePinnedIdx = index
}

func (m *Model) cyclePinnedNext() {
	if len(m.pinnedNotes) == 0 {
		m.addToast("No pinned notes", ToastWarning)
		return
	}

	if m.activePinnedIdx >= 0 && m.activePinnedIdx < len(m.pinnedNotes) && m.activeNote != nil {
		m.pinnedNotes[m.activePinnedIdx].ScrollY = m.viewer.GetScrollPosition()
	}

	m.activePinnedIdx++
	if m.activePinnedIdx >= len(m.pinnedNotes) {
		m.activePinnedIdx = 0
	}

	m.openPinnedNote(m.activePinnedIdx)
}

func (m *Model) cyclePinnedPrev() {
	if len(m.pinnedNotes) == 0 {
		m.addToast("No pinned notes", ToastWarning)
		return
	}

	if m.activePinnedIdx >= 0 && m.activePinnedIdx < len(m.pinnedNotes) && m.activeNote != nil {
		m.pinnedNotes[m.activePinnedIdx].ScrollY = m.viewer.GetScrollPosition()
	}

	m.activePinnedIdx--
	if m.activePinnedIdx < 0 {
		m.activePinnedIdx = len(m.pinnedNotes) - 1
	}

	m.openPinnedNote(m.activePinnedIdx)
}

func (m *Model) validatePins() {
	var valid []PinnedNote
	for _, pin := range m.pinnedNotes {
		path := filepath.Join(m.config.VaultPath, pin.Path)
		if _, err := os.Stat(path); err == nil {
			valid = append(valid, pin)
		}
	}
	m.pinnedNotes = valid
	if m.activePinnedIdx >= len(m.pinnedNotes) {
		m.activePinnedIdx = len(m.pinnedNotes) - 1
	}
}
