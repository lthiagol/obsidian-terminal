package main

import (
	"fmt"

	"github.com/lthiagol/obsidian-terminal/internal/search"
)

type noteNavKind int

const (
	navUser    noteNavKind = iota // user navigation: push history, clear forward
	navHistory                    // back/forward: stacks already updated
	navReload                     // rescan refresh: no history changes
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

func (m *Model) enterTagsMode() {
	m.prevMode = m.mode
	m.mode = ModeTags
	m.tagList = NewTagList(m.tagIndex, m.palette)
}
