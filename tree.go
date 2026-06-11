package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type treeItem struct {
	entry    *VaultEntry
	expanded bool
	depth    int
}

type FileTree struct {
	items  []treeItem
	cursor int
	vault  *VaultEntry
	width  int
	height int
}

func NewFileTree(vault *VaultEntry) FileTree {
	ft := FileTree{
		vault:  vault,
		width:  25,
		height: 20,
	}
	var items []treeItem
	for _, child := range vault.Children {
		items = append(items, flattenTree(child, 0, !child.IsDir)...)
	}
	ft.items = items
	return ft
}

func flattenTree(entry *VaultEntry, depth int, expanded bool) []treeItem {
	var items []treeItem

	if entry.IsDir {
		items = append(items, treeItem{entry: entry, depth: depth, expanded: expanded})
		if entry.Children != nil && expanded {
			for _, child := range entry.Children {
				items = append(items, flattenTree(child, depth+depthIncrement, !child.IsDir)...)
			}
		}
	} else {
		items = append(items, treeItem{entry: entry, depth: depth})
	}

	return items
}

func (ft *FileTree) expand() {
	if ft.cursor >= len(ft.items) {
		return
	}

	item := &ft.items[ft.cursor]
	if !item.entry.IsDir || item.expanded {
		return
	}

	item.expanded = true

	if item.entry.Children == nil {
		return
	}

	childItems := flattenChildren(item.entry.Children, item.depth+depthIncrement)

	pos := ft.cursor + 1
	ft.items = append(ft.items[:pos], append(childItems, ft.items[pos:]...)...)
}

func (ft *FileTree) collapse() {
	if ft.cursor >= len(ft.items) {
		return
	}

	item := &ft.items[ft.cursor]
	if !item.entry.IsDir || !item.expanded {
		return
	}

	item.expanded = false

	cutEnd := ft.cursor + 1
	for cutEnd < len(ft.items) && ft.items[cutEnd].depth > item.depth {
		cutEnd++
	}

	ft.items = append(ft.items[:ft.cursor+1], ft.items[cutEnd:]...)
}

func flattenChildren(children []*VaultEntry, depth int) []treeItem {
	var items []treeItem
	for _, child := range children {
		items = append(items, flattenTree(child, depth, !child.IsDir)...)
	}
	return items
}

func (ft FileTree) Update(msg tea.Msg) (FileTree, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return ft, nil
	}

	switch {
	case MatchKey(keyMsg, ft.keys().Up) || MatchRune(keyMsg, ft.keysRunes().UpRune):
		if ft.cursor > 0 {
			ft.cursor--
		}
	case MatchKey(keyMsg, ft.keys().Down) || MatchRune(keyMsg, ft.keysRunes().DownRune):
		if ft.cursor < len(ft.items)-1 {
			ft.cursor++
		}
	case MatchKey(keyMsg, ft.keys().Right) || MatchRune(keyMsg, ft.keysRunes().RightRune):
		ft.expand()
	case MatchKey(keyMsg, ft.keys().Left) || MatchRune(keyMsg, ft.keysRunes().LeftRune):
		ft.collapse()
	case MatchRune(keyMsg, ft.keysRunes().TopRune):
		ft.cursor = 0
	case MatchRune(keyMsg, ft.keysRunes().BottomRune):
		if len(ft.items) > 0 {
			ft.cursor = len(ft.items) - 1
		}
	case keyMsg.Type == tea.KeyEnter:
		ft.toggleExpand()
	}

	return ft, nil
}

func (ft FileTree) keys() KeyMap {
	return DefaultKeys()
}

func (ft FileTree) keysRunes() KeyMap {
	return DefaultKeys()
}

func (ft *FileTree) toggleExpand() {
	if ft.cursor >= len(ft.items) {
		return
	}
	item := ft.items[ft.cursor]
	if item.entry.IsDir {
		if item.expanded {
			ft.collapse()
		} else {
			ft.expand()
		}
	}
}

func (ft FileTree) SelectedEntry() *VaultEntry {
	if ft.cursor >= 0 && ft.cursor < len(ft.items) {
		return ft.items[ft.cursor].entry
	}
	return nil
}

func (ft FileTree) IsDirSelected() bool {
	entry := ft.SelectedEntry()
	return entry != nil && entry.IsDir
}

func (ft FileTree) SelectedPath() string {
	entry := ft.SelectedEntry()
	if entry != nil {
		return entry.Path
	}
	return ""
}

func (ft FileTree) SetSize(width, height int) {
	_ = width
	_ = height
}

func (ft FileTree) Cursor() int {
	return ft.cursor
}

func (ft *FileTree) MoveUp() {
	if ft.cursor > 0 {
		ft.cursor--
	}
}

func (ft *FileTree) MoveDown() {
	if ft.cursor < len(ft.items)-1 {
		ft.cursor++
	}
}

func (ft FileTree) View() string {
	if len(ft.items) == 0 {
		return lipgloss.NewStyle().
			Foreground(TextMuted).
			PaddingTop(2).
			PaddingLeft(2).
			Render("no notes found")
	}

	var sb strings.Builder
	availableWidth := ft.width - 4
	if availableWidth < 10 {
		availableWidth = 10
	}

	for i, item := range ft.items {
		isSelected := i == ft.cursor

		prefix := strings.Repeat("  ", item.depth)

		var icon string
		if item.entry.IsDir {
			if item.expanded {
				icon = IconFolderOpen
			} else {
				icon = IconFolderClosed
			}
		} else {
			icon = IconFile
		}

		name := item.entry.Name
		if item.entry.IsSymlink {
			name += " ->"
		}

		fullLine := prefix + icon + name
		if len(fullLine) > availableWidth {
			fullLine = fullLine[:availableWidth]
		}

		line := lipgloss.NewStyle()

		if item.entry.IsDir {
			line = line.Foreground(AccentSecondary)
		} else {
			line = line.Foreground(TextSecondary)
		}

		if isSelected {
			line = line.Background(Accent).Foreground(lipgloss.Color("#000000")).Bold(true)
		}

		sb.WriteString(line.Render(fullLine))
		if i < len(ft.items)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (ft FileTree) Items() []treeItem {
	return ft.items
}

func (ft FileTree) ItemCount() int {
	return len(ft.items)
}

const depthIncrement = 1
