package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRenderingPipeline_FullDocument(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	output := m.viewer.View()
	if output == "" {
		t.Error("viewer should render content")
	}

	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if hasTruncatedANSI(line) {
			t.Errorf("line %d has broken ANSI escape: %q", i, line)
		}
	}

	body := m.activeNote.Body
	if !strings.Contains(body, "Welcome") || !strings.Contains(body, "#") {
		t.Error("activeNote body should contain expected content")
	}
}

func TestWorkflow_SearchAndOpen(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = model.(Model)
	if m.mode != ModeSearch {
		t.Fatalf("expected ModeSearch after /, got %v", m.mode)
	}

	for _, r := range "index" {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m = model.(Model)
	if m.searchState.Query() != "index" {
		t.Fatalf("expected query 'index', got %q", m.searchState.Query())
	}
	if m.searchState.ResultCount() == 0 {
		t.Fatal("expected at least one search result for 'index'")
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.mode != ModeView {
		t.Fatalf("expected ModeView after Enter on search result, got %v", m.mode)
	}
	if m.activeNote == nil {
		t.Fatal("expected activeNote to be set")
	}
	if !strings.HasSuffix(m.activeNote.Path, "index.md") {
		t.Errorf("expected index.md, got %s", m.activeNote.Path)
	}
}

func TestWorkflow_TreeClickAndOpen(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	if m.activeNote == nil {
		t.Fatal("activeNote should be set")
	}
	if m.activeNote.Path == "" {
		t.Error("activeNote.Path should not be empty")
	}
	if len(m.outlineItems) == 0 {
		t.Error("outlineItems should be populated for a note with headings")
	}
	if len(m.recentNotes) == 0 {
		t.Error("recentNotes should contain the opened file")
	}
	if m.recentNotes[0] != m.activeNote.Path {
		t.Errorf("first recent should be opened file %q, got %q", m.activeNote.Path, m.recentNotes[0])
	}
}

func TestWorkflow_FollowWikiLink(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	if m.viewer.LinkCount() < 1 {
		t.Fatal("expected at least one wiki-link in index.md")
	}

	model = m
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = model.(Model)
	selectedLink := m.viewer.SelectedLinkPath()
	if selectedLink == "" {
		t.Error("selected link path should not be empty after Tab")
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.mode != ModeView {
		t.Fatalf("expected ModeView after following link, got %v", m.mode)
	}
	if m.activeNote == nil {
		t.Fatal("activeNote should be set after following link")
	}
	if m.activeNote.Path == "" {
		t.Error("activeNote.Path should be the linked note path")
	}
}

func TestStatePreservation_ThemeSwitch(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	prevPath := m.activeNote.Path
	prevMode := m.mode
	prevAccent := m.palette.Accent

	m.setTheme("dracula")

	if m.mode != prevMode {
		t.Errorf("mode changed after theme switch: %v -> %v", prevMode, m.mode)
	}
	if m.activeNote == nil {
		t.Fatal("activeNote should still be set after theme switch")
	}
	if m.activeNote.Path != prevPath {
		t.Errorf("activeNote path changed: %s -> %s", prevPath, m.activeNote.Path)
	}
	if m.palette.Accent == prevAccent {
		t.Error("palette colors should change after theme switch")
	}

	output := m.viewer.View()
	if output == "" {
		t.Error("viewer should still render content after theme switch")
	}
}

func TestStatePreservation_SplitResize(t *testing.T) {
	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	m = navigateToFirstFile(t, &model)

	prevPath := m.activeNote.Path
	prevMode := m.mode
	prevWidth := m.treeWidth

	model = m
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})
	m = model.(Model)

	if m.treeWidth >= prevWidth {
		t.Errorf("treeWidth should shrink: %d -> %d", prevWidth, m.treeWidth)
	}
	if m.mode != prevMode {
		t.Errorf("mode changed after split resize: %v -> %v", prevMode, m.mode)
	}
	if m.activeNote == nil {
		t.Fatal("activeNote should still be set after resize")
	}
	if m.activeNote.Path != prevPath {
		t.Errorf("activeNote path changed: %s -> %s", prevPath, m.activeNote.Path)
	}

	output := m.viewer.View()
	if output == "" {
		t.Error("viewer should still render content after resize")
	}
}

func TestSession_SaveAndRestore(t *testing.T) {
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	m := newTestModel(t, &Config{})
	model := tea.Model(m)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = model.(Model)

	dirIdx := indexOfFirstCollapsedDir(m.fileTree)
	if dirIdx < 0 {
		t.Skip("no collapsed directory to test")
	}
	for m.fileTree.Cursor() < dirIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}
	expandedDir := m.fileTree.Items()[m.fileTree.Cursor()].entry.Path

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = model.(Model)

	for i := 0; i < 3; i++ {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}

	selected := m.fileTree.SelectedEntry()
	if selected == nil || selected.IsDir {
		firstFileAfter := indexOfFirstFileFrom(m.fileTree, m.fileTree.Cursor())
		for m.fileTree.Cursor() < firstFileAfter {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			m = model.(Model)
		}
	}

	saveSession(m)

	newM := NewModel(&Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs})
	if newM.err != nil {
		t.Fatalf("NewModel error: %v", newM.err)
	}
	newModel, _ := newM.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	restored := newModel.(Model)

	foundExpanded := false
	for _, item := range restored.fileTree.Items() {
		if item.entry.IsDir && item.entry.Path == expandedDir && item.expanded {
			foundExpanded = true
			break
		}
	}
	if !foundExpanded {
		t.Errorf("expected directory %q to be expanded after session restore", expandedDir)
	}
}

func indexOfFirstFileFrom(ft FileTree, start int) int {
	for i := start; i < len(ft.Items()); i++ {
		if !ft.Items()[i].entry.IsDir {
			return i
		}
	}
	return start
}
