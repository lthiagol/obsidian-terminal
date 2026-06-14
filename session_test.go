package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSessionSaveRestore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	stateFile := filepath.Join(dir, "obsidian-terminal", "session.json")
	// Override state file path by setting XDG_STATE_HOME
	t.Setenv("XDG_STATE_HOME", dir)

	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	// Save session
	saveSession(m)

	// Verify file exists
	data, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatalf("session file not created: %v", err)
	}

	var s SessionState
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("invalid session JSON: %v", err)
	}
	if s.VaultPath != testVaultPath(t) {
		t.Errorf("vault path = %q", s.VaultPath)
	}
	if s.Version != sessionVersion {
		t.Errorf("version = %d", s.Version)
	}
}

func TestSessionRestore_SkipsDifferentVault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	// Write a session for a different vault
	s := SessionState{
		VaultPath: "/different/vault",
		Version:   sessionVersion,
	}
	data, _ := json.Marshal(s)
	os.MkdirAll(filepath.Join(dir, "obsidian-terminal"), 0755)
	os.WriteFile(filepath.Join(dir, "obsidian-terminal", "session.json"), data, 0600)

	// Loading a model with a different vault should not crash
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}
	_ = m
}

func TestSessionRestore_CorruptedFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	// Write corrupted data
	os.MkdirAll(filepath.Join(dir, "obsidian-terminal"), 0755)
	os.WriteFile(filepath.Join(dir, "obsidian-terminal", "session.json"), []byte("{invalid json}"), 0600)

	// Should not crash
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}
	_ = m
}
