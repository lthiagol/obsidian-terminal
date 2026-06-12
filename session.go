package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SessionState holds the tree/viewer state to restore on next launch.
type SessionState struct {
	VaultPath  string   `json:"vault_path"`
	Expanded   []string `json:"expanded,omitempty"`
	CursorPath string   `json:"cursor_path,omitempty"`
	Version    int      `json:"version"`
}

const sessionVersion = 1

func stateFilePath() string {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		stateDir = filepath.Join(home, ".local", "state")
	}
	dir := filepath.Join(stateDir, "obsidian-terminal")
	return filepath.Join(dir, "session.json")
}

func saveSession(m Model) {
	path := stateFilePath()
	if path == "" || m.config.VaultPath == "" {
		return
	}

	s := SessionState{
		VaultPath: m.config.VaultPath,
		Version:   sessionVersion,
	}

	// Collect expanded directory paths
	for _, item := range m.fileTree.Items() {
		if item.entry.IsDir && item.expanded {
			s.Expanded = append(s.Expanded, item.entry.Path)
		}
	}

	// Save cursor position
	entry := m.fileTree.SelectedEntry()
	if entry != nil {
		s.CursorPath = entry.Path
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	os.WriteFile(path, data, 0600)
}

func restoreSession(m *Model) {
	path := stateFilePath()
	if path == "" {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return // file missing or unreadable — start fresh
	}

	var s SessionState
	if err := json.Unmarshal(data, &s); err != nil {
		return // corrupted file — start fresh
	}

	if s.VaultPath != m.config.VaultPath {
		return // different vault — don't restore
	}

	// Expand saved directories
	for _, expandedPath := range s.Expanded {
		m.fileTree.ExpandPath(expandedPath)
	}

	// Restore cursor
	if s.CursorPath != "" {
		for i, item := range m.fileTree.Items() {
			if item.entry.Path == s.CursorPath {
				m.fileTree.MoveToY(i)
				break
			}
		}
	}
}
