package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_InitialMode_IsBrowse(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)

	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}
	if m.mode != ModeBrowse {
		t.Errorf("initial mode = %v, want ModeBrowse", m.mode)
	}
}

func TestKeyDispatch_Browse_JK(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	if m.fileTree.Cursor() != 0 {
		t.Errorf("initial cursor = %d, want 0", m.fileTree.Cursor())
	}

	// Down with 'j'
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = model.(Model)
	if m.fileTree.Cursor() != 1 {
		t.Errorf("after j: cursor = %d, want 1", m.fileTree.Cursor())
	}

	// Down with Down arrow
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = model.(Model)
	if m.fileTree.Cursor() != 2 {
		t.Errorf("after Down: cursor = %d, want 2", m.fileTree.Cursor())
	}

	// Up with 'k'
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = model.(Model)
	if m.fileTree.Cursor() != 1 {
		t.Errorf("after k: cursor = %d, want 1", m.fileTree.Cursor())
	}

	// Up with Up arrow
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = model.(Model)
	if m.fileTree.Cursor() != 0 {
		t.Errorf("after Up: cursor = %d, want 0", m.fileTree.Cursor())
	}

	// Up at 0 should stay at 0
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = model.(Model)
	if m.fileTree.Cursor() != 0 {
		t.Errorf("Up at 0: cursor = %d, want 0", m.fileTree.Cursor())
	}
}

func TestKeyDispatch_View_Esc(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	firstFileIdx := indexOfFirstFile(m.fileTree)
	for m.fileTree.Cursor() < firstFileIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.mode != ModeView {
		t.Fatalf("after Enter: mode = %v, want ModeView", m.mode)
	}
	if m.activeNote == nil {
		t.Fatal("expected activeNote to be set after Enter")
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("after Esc: mode = %v, want ModeBrowse", m.mode)
	}
	if m.activeNote != nil {
		t.Error("expected activeNote to be nil after Esc")
	}
}

func TestKeyDispatch_Search_OpenClose(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m := model.(Model)
	if m.mode != ModeSearch {
		t.Errorf("after /: mode = %v, want ModeSearch", m.mode)
	}
	if m.prevMode != ModeBrowse {
		t.Errorf("prevMode = %v, want ModeBrowse", m.prevMode)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("after Esc: mode = %v, want ModeBrowse", m.mode)
	}
}

func TestKeyDispatch_Help_Toggle(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m := model.(Model)
	if m.mode != ModeHelp {
		t.Errorf("after ?: mode = %v, want ModeHelp", m.mode)
	}
	if m.prevMode != ModeBrowse {
		t.Errorf("prevMode = %v, want ModeBrowse", m.prevMode)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("after Esc: mode = %v, want ModeBrowse", m.mode)
	}
}

func TestModeTransitions(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	firstFileIdx := indexOfFirstFile(m.fileTree)
	for m.fileTree.Cursor() < firstFileIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}

	// browse → view (Enter on a file)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.mode != ModeView {
		t.Errorf("browse→view: mode = %v", m.mode)
	}

	// view → browse (Esc)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("view→browse: mode = %v", m.mode)
	}

	// browse → search (/)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = model.(Model)
	if m.mode != ModeSearch {
		t.Errorf("browse→search: mode = %v", m.mode)
	}

	// search → browse (Esc)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("search→browse: mode = %v", m.mode)
	}

	// browse → help (?)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = model.(Model)
	if m.mode != ModeHelp {
		t.Errorf("browse→help: mode = %v", m.mode)
	}

	// help → browse (Esc)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("help→browse: mode = %v", m.mode)
	}

	// browse → view → help
	for m.fileTree.Cursor() < firstFileIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.mode != ModeView {
		t.Errorf("2nd browse→view: mode = %v", m.mode)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = model.(Model)
	if m.mode != ModeHelp {
		t.Errorf("view→help: mode = %v", m.mode)
	}

	// help → view (Esc)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeView {
		t.Errorf("help→view: mode = %v, prevMode=%v", m.mode, m.prevMode)
	}
}

func TestTree_OpenNoteSwitchesToView(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	firstFileIdx := indexOfFirstFile(m.fileTree)
	for m.fileTree.Cursor() < firstFileIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}

	// Enter on .md file should switch to view mode
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.mode != ModeView {
		t.Fatalf("expected ModeView after Enter on md file, got %v", m.mode)
	}
	if m.activeNote == nil {
		t.Fatal("expected activeNote to be set")
	}
	if m.activeNote.Path == "" {
		t.Error("activeNote.Path should not be empty")
	}
}

func TestTree_EnterOnFolderExpands(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	dirIdx := indexOfFirstCollapsedDir(m.fileTree)
	if dirIdx < 0 {
		t.Skip("no collapsed directory to test")
	}

	for m.fileTree.Cursor() < dirIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}

	entry := m.fileTree.SelectedEntry()
	if entry == nil || !entry.IsDir {
		t.Fatalf("expected a directory at cursor %d", m.fileTree.Cursor())
	}

	itemCountBefore := m.fileTree.ItemCount()

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	itemCountAfter := m.fileTree.ItemCount()

	if itemCountAfter <= itemCountBefore {
		t.Errorf("expected children to appear after Enter on dir, items before=%d after=%d", itemCountBefore, itemCountAfter)
	}
}

func TestStatusBar_ShowsMode(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	// Simulate window resize to make model ready
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	view := model.View()

	if !strings.Contains(view, "BROWSE") {
		t.Error("status bar should show mode BROWSE")
	}
}

func TestStatusBar_ShowsCurrentFile(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Open first file
	firstFileIdx := indexOfFirstFile(m.fileTree)
	m = model.(Model)
	for m.fileTree.Cursor() < firstFileIdx {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	view := model.View()

	// Status bar should show VIEW mode
	if !strings.Contains(view, "VIEW") {
		t.Error("status bar should show VIEW mode")
	}
}

func TestStatusBar_ShowsHints(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Browse mode hints
	view := model.View()
	if !strings.Contains(view, "search") {
		t.Error("browse status bar should show hints with 'search'")
	}
	if !strings.Contains(view, "help") {
		t.Error("browse status bar should show hints with 'help'")
	}

	// Open help - should show help hints
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	view = model.View()
	if !strings.Contains(view, "HELP") {
		t.Error("status bar should show HELP mode")
	}
}

func TestHelpPanel_ShowsAllSections(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	if model.(Model).err != nil {
		t.Fatalf("NewModel error: %v", model.(Model).err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	view := model.View()

	sections := []string{"Navigation", "File Tree", "Viewer", "Search", "Global"}
	for _, s := range sections {
		if !strings.Contains(view, s) {
			t.Errorf("help panel should contain section: %s", s)
		}
	}
}

func TestHelpPanel_EscCloses(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	if model.(Model).err != nil {
		t.Fatalf("NewModel error: %v", model.(Model).err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

	if model.(Model).mode != ModeHelp {
		t.Fatal("expected ModeHelp")
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if model.(Model).mode != ModeBrowse {
		t.Errorf("after Esc: mode = %v, want ModeBrowse", model.(Model).mode)
	}
}
func TestCheckVaultChanges_StatError(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	originalPath := m.config.VaultPath
	m.config.VaultPath = "/nonexistent/path/for/stat/test"
	m.checkVaultChanges()
	if len(m.toasts) == 0 {
		t.Error("expected toast on stat error")
	}
	m.config.VaultPath = originalPath
}

func TestRescanVault_ScanError(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	originalPath := m.config.VaultPath
	m.config.VaultPath = "/nonexistent/path/for/rescan/test"
	m.rescanVault()
	m.config.VaultPath = originalPath
}

func TestTruncatePath(t *testing.T) {
	if p := truncatePath("short", 20); p != "short" {
		t.Errorf("short path: got %q, want 'short'", p)
	}
	long := "a/very/long/path/that/should/be/truncated.md"
	p := truncatePath(long, 20)
	if len(p) > 20 {
		t.Errorf("truncated path should be <= maxLen: %d > 20", len(p))
	}
	if !strings.HasPrefix(p, ".../") {
		t.Error("truncated path should start with .../")
	}
}

func TestKeyDispatch_TagsEnterExit(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// T enters tags mode
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})
	m = model.(Model)
	if m.mode != ModeTags {
		t.Errorf("after T: mode = %v, want ModeTags", m.mode)
	}
	if m.prevMode != ModeBrowse {
		t.Errorf("prevMode = %v, want ModeBrowse", m.prevMode)
	}

	// Esc exits back to browse
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeBrowse {
		t.Errorf("after Esc from tags: mode = %v, want ModeBrowse", m.mode)
	}
}

func TestKeyDispatch_FindEnterExit(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)

	// First open a file to get into view mode
	m := navigateToFirstFile(t, &model)

	// 's' enters find mode from view
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = model.(Model)
	if m.mode != ModeFind {
		t.Errorf("after s: mode = %v, want ModeFind", m.mode)
	}

	// Esc exits back to view
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.mode != ModeView {
		t.Errorf("after Esc from find: mode = %v, want ModeView", m.mode)
	}
}

func TestKeyDispatch_BacklinkPanel(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Navigate to a file that has backlinks (index.md has alias "Home" and is referenced)
	// First file in tree should work
	m = navigateToFirstFile(t, &model)
	if m.mode != ModeView {
		t.Fatalf("expected ModeView, got %v", m.mode)
	}

	// Press b to try activating backlink panel
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = model.(Model)
	// backlinkMode depends on whether the note has backlinks
	if m.backlinkMode {
		// Successfully entered backlink mode, now Esc should exit
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		m = model.(Model)
		if m.backlinkMode {
			t.Error("Esc should exit backlink mode")
		}
	}
	// backlinkMode may be false if the note has no backlinks — that's fine
}

func TestModel_ViewReturnsNonEmpty(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	view := model.View()
	if view == "" {
		t.Error("View() should not return empty string after WindowSizeMsg")
	}
}

func navigateToFirstFile(t *testing.T, model *tea.Model) Model {
	t.Helper()
	m := (*model).(Model)
	firstFileIdx := indexOfFirstFile(m.fileTree)
	for m.fileTree.Cursor() < firstFileIdx {
		*model, _ = (*model).Update(tea.KeyMsg{Type: tea.KeyDown})
		m = (*model).(Model)
	}
	*model, _ = (*model).Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = (*model).(Model)
	if m.mode != ModeView {
		t.Fatalf("navigateToFirstFile: expected ModeView, got %v", m.mode)
	}
	return m
}

func indexOfFirstFile(ft FileTree) int {
	for i, item := range ft.Items() {
		if !item.entry.IsDir {
			return i
		}
	}
	return 0
}

func indexOfFirstCollapsedDir(ft FileTree) int {
	for i, item := range ft.Items() {
		if item.entry.IsDir && !item.expanded {
			return i
		}
	}
	return -1
}

func TestKeyDispatch_ShrinkTree(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	initialWidth := m.treeWidth

	// Ctrl+Left shrinks tree
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})
	m = model.(Model)
	if m.treeWidth >= initialWidth {
		t.Errorf("Ctrl+Left should shrink tree: was %d, now %d", initialWidth, m.treeWidth)
	}
}

func TestKeyDispatch_GrowTree(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)

	// Shrink first so we can grow
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})
	m = model.(Model)
	shrunkWidth := m.treeWidth

	// Ctrl+Right grows tree
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlRight})
	m = model.(Model)
	if m.treeWidth <= shrunkWidth {
		t.Errorf("Ctrl+Right should grow tree: was %d, now %d", shrunkWidth, m.treeWidth)
	}
}

func TestKeyDispatch_ResetTree(t *testing.T) {
	cfg := &Config{VaultPath: testVaultPath(t), SkipDirs: DefaultConfig().SkipDirs}
	var model tea.Model = NewModel(cfg)
	m := model.(Model)
	if m.err != nil {
		t.Fatalf("NewModel error: %v", m.err)
	}

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)

	// Shrink then reset
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlBackslash})
	m = model.(Model)
	// Should be close to default (width/4 for width=100 = 25)
	if m.treeWidth < 22 || m.treeWidth > 28 {
		t.Errorf("Reset should restore ~25: got %d", m.treeWidth)
	}
}
