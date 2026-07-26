package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMouse_TreeClick(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.width = 120
	model.height = 40
	model.treeWidth = 30
	model.ready = true

	initialCursor := model.fileTree.Cursor()

	mouseMsg := tea.MouseMsg{
		X:      5,
		Y:      3,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	result, _ := model.Update(mouseMsg)
	m := result.(Model)

	if m.fileTree.Cursor() == initialCursor && len(m.fileTree.Items()) > 3 {
		t.Error("tree click should move cursor")
	}
}

func TestMouse_WheelUp(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.width = 120
	model.height = 40
	model.treeWidth = 30
	model.ready = true

	model.fileTree.MoveDown()
	model.fileTree.MoveDown()
	model.fileTree.MoveDown()
	cursorBefore := model.fileTree.Cursor()

	mouseMsg := tea.MouseMsg{
		X:      5,
		Y:      3,
		Button: tea.MouseButtonWheelUp,
		Action: tea.MouseActionPress,
	}
	result, _ := model.Update(mouseMsg)
	m := result.(Model)

	if m.fileTree.Cursor() >= cursorBefore {
		t.Error("wheel up should decrease cursor")
	}
}

func TestDoubleClick_Detection(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.recordClick(5, 5)
	if !model.isDoubleClick(5, 5) {
		t.Error("immediate second click should be detected as double-click")
	}
	if model.isDoubleClick(20, 20) {
		t.Error("distant click should not be double-click")
	}
}

func TestSearchState_SetSelected(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.enterSearchMode()

	model.searchState.SetSelected(0)
	if model.searchState.SelectedIndex() != 0 {
		t.Errorf("selected = %d, want 0", model.searchState.SelectedIndex())
	}

	model.searchState.SetSelected(-5)
	if model.searchState.SelectedIndex() != 0 {
		t.Errorf("clamped low: selected = %d, want 0", model.searchState.SelectedIndex())
	}

	model.searchState.SetSelected(99999)
	if model.searchState.SelectedIndex() != model.searchState.ResultCount()-1 {
		t.Errorf("clamped high: selected = %d, want %d",
			model.searchState.SelectedIndex(), model.searchState.ResultCount()-1)
	}
}

func TestMouse_SplitDrag(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.width = 120
	model.height = 40
	model.treeWidth = 30
	model.ready = true

	// Press near the split boundary
	mouseMsg := tea.MouseMsg{
		X:      30,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	result, _ := model.Update(mouseMsg)
	m := result.(Model)

	if !m.dragSplit {
		t.Fatal("press near split should start drag")
	}

	// Drag to the right
	mouseMsg = tea.MouseMsg{
		X:      45,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	result, _ = m.Update(mouseMsg)
	m = result.(Model)

	if m.treeWidth <= 30 {
		t.Errorf("drag right should widen tree: got %d", m.treeWidth)
	}

	// Release
	mouseMsg = tea.MouseMsg{
		X:      45,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionRelease,
	}
	result, _ = m.Update(mouseMsg)
	m = result.(Model)

	if m.dragSplit {
		t.Error("release should stop drag")
	}
}

func TestMouse_SplitDrag_ClampsMinWidth(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	model := NewModel(cfg)
	if model.err != nil {
		t.Fatalf("NewModel: %v", model.err)
	}

	model.width = 120
	model.height = 40
	model.treeWidth = 30
	model.ready = true

	// Start drag
	mouseMsg := tea.MouseMsg{
		X:      31,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	result, _ := model.Update(mouseMsg)
	m := result.(Model)

	// Drag way left (past minimum)
	mouseMsg = tea.MouseMsg{
		X:      1,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	result, _ = m.Update(mouseMsg)
	m = result.(Model)

	if m.treeWidth < 15 {
		t.Errorf("tree width should not go below 15: got %d", m.treeWidth)
	}
}
