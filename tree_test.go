package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTree_RootExpanded(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	if ft.ItemCount() == 0 {
		t.Error("tree should have items (root expanded, children visible)")
	}

	// First-level children should be at depth 0
	found := false
	for _, item := range ft.Items() {
		if item.depth == 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected depth-0 children in tree")
	}
}

func TestTree_ExpandCollapse(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	// Find a directory
	dirIdx := -1
	for i, item := range ft.Items() {
		if item.entry.IsDir {
			dirIdx = i
			break
		}
	}
	if dirIdx < 0 {
		t.Skip("no directories in tree")
	}

	// Navigate to directory
	for ft.cursor < dirIdx {
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	entry := ft.SelectedEntry()
	if entry == nil || !entry.IsDir {
		t.Fatal("expected a directory")
	}

	// If already expanded, collapse and test expand
	if ft.Items()[ft.cursor].expanded {
		// Collapse with Left/←
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		if ft.Items()[ft.cursor].expanded {
			t.Error("expected directory to be collapsed after Left/h")
		}

		itemsCollapsed := ft.ItemCount()

		// Expand with Right/→
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		if !ft.Items()[ft.cursor].expanded {
			t.Error("expected directory to be expanded after Right/l")
		}

		if ft.ItemCount() <= itemsCollapsed {
			t.Error("expected more items after expand")
		}
	} else {
		// Expand with Right/→
		itemsBefore := ft.ItemCount()
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		if !ft.Items()[ft.cursor].expanded {
			t.Error("expected directory to be expanded after Right/l")
		}
		if ft.ItemCount() <= itemsBefore {
			t.Error("expected more items after expand")
		}

		// Collapse with Left/←
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		if ft.Items()[ft.cursor].expanded {
			t.Error("expected directory to be collapsed after Left/h")
		}
	}
}

func TestTree_SelectionClamped(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	// Go to bottom
	maxIdx := ft.ItemCount() - 1
	for i := 0; i < ft.ItemCount(); i++ {
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	if ft.Cursor() != maxIdx {
		t.Errorf("at bottom: cursor = %d, want %d", ft.Cursor(), maxIdx)
	}

	// Down at bottom should stay at bottom
	ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyDown})
	if ft.Cursor() != maxIdx {
		t.Errorf("Down at bottom: cursor = %d, want %d", ft.Cursor(), maxIdx)
	}

	// Go to top
	for ft.Cursor() > 0 {
		ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyUp})
	}
	if ft.Cursor() != 0 {
		t.Errorf("at top: cursor = %d, want 0", ft.Cursor())
	}

	// Up at top should stay at top
	ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyUp})
	if ft.Cursor() != 0 {
		t.Errorf("Up at top: cursor = %d, want 0", ft.Cursor())
	}
}

func TestTree_ViewAfterExpand(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	// Expand all directories to reveal deep nesting
	for i := 0; i < ft.ItemCount(); i++ {
		item := ft.Items()[i]
		if item.entry.IsDir && !item.expanded {
			ft.MoveToY(i)
			ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		}
	}

	// View() must not panic with nested directories expanded
	_ = ft.View()
}

func TestTree_SymlinkShownInTree(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	foundSymlink := false
	for _, item := range ft.Items() {
		if item.entry.IsSymlink {
			foundSymlink = true
			if item.entry.Name != "readme-symlink.md" {
				t.Errorf("unexpected symlink name: %s", item.entry.Name)
			}
			break
		}
	}

	if !foundSymlink {
		t.Error("expected readme-symlink.md to appear in tree")
	}
}

func TestTree_ApplyPathFilter(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	// Filter to only show a specific file
	filter := map[string]bool{
		"readme.md": true,
	}
	ft.ApplyPathFilter(filter)

	if ft.ItemCount() == 0 {
		t.Fatal("expected at least 1 filtered item")
	}

	for _, item := range ft.Items() {
		if !item.entry.IsDir {
			if !filter[item.entry.Path] {
				t.Errorf("unfiltered file: %s", item.entry.Path)
			}
		}
	}

	if ft.Cursor() != 0 {
		t.Errorf("cursor should reset to 0 after filter, got %d", ft.Cursor())
	}
}

func TestTree_ApplyPathFilter_Empty(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)
	filter := map[string]bool{} // empty filter
	ft.ApplyPathFilter(filter)

	if ft.ItemCount() != 0 {
		t.Errorf("empty filter should result in zero items, got %d", ft.ItemCount())
	}
}

func TestTree_ResetFilter(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)
	origCount := ft.ItemCount()

	// Apply a narrow filter
	filter := map[string]bool{"readme.md": true}
	ft.ApplyPathFilter(filter)

	filteredCount := ft.ItemCount()
	if filteredCount >= origCount {
		t.Fatalf("filter should reduce items: %d >= %d", filteredCount, origCount)
	}

	// Reset back to full tree
	ft.ResetFilter(tree)

	if ft.ItemCount() != origCount {
		t.Errorf("ResetFilter should restore original count: got %d, want %d", ft.ItemCount(), origCount)
	}
	if ft.Cursor() != 0 {
		t.Errorf("cursor should reset to 0 after ResetFilter, got %d", ft.Cursor())
	}
}

func TestTree_TopBottomJumps(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)
	if ft.ItemCount() < 3 {
		t.Skip("need at least 3 items")
	}

	// g jumps to top
	ft.MoveToY(ft.ItemCount() - 1) // go to bottom
	ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if ft.Cursor() != 0 {
		t.Errorf("g should jump to top, cursor = %d", ft.Cursor())
	}

	// G jumps to bottom
	ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	if ft.Cursor() != ft.ItemCount()-1 {
		t.Errorf("G should jump to bottom, cursor = %d (want %d)", ft.Cursor(), ft.ItemCount()-1)
	}
}

func TestTree_EnterOnExpandedDir(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	ft := NewFileTree(tree)

	// Find a directory
	dirIdx := -1
	for i, item := range ft.Items() {
		if item.entry.IsDir {
			dirIdx = i
			break
		}
	}
	if dirIdx < 0 {
		t.Skip("no directories")
	}

	// Expand
	ft.MoveToY(dirIdx)
	ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	beforeCount := ft.ItemCount()

	// Enter on expanded directory should collapse
	ft, _ = ft.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if ft.Items()[dirIdx].expanded {
		t.Error("Enter on expanded dir should collapse it")
	}
	if ft.ItemCount() >= beforeCount {
		t.Error("item count should decrease after Enter collapses")
	}
}

func TestNewFileTree_EmptyVault(t *testing.T) {
	root := &VaultEntry{
		Name:     ".",
		Path:     "",
		IsDir:    true,
		Children: nil,
	}
	ft := NewFileTree(root)
	if ft.ItemCount() != 0 {
		t.Errorf("empty vault should have 0 items, got %d", ft.ItemCount())
	}
	// View should not panic
	_ = ft.View()
}
