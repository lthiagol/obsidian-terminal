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
