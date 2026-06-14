package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TagEntry struct {
	Name  string
	Count int
	Files []string
}

type TagList struct {
	entries []TagEntry
	cursor  int
	width   int
	height  int
	palette Palette
}

func NewTagList(tagIndex map[string][]string, palette Palette) TagList {
	var entries []TagEntry
	for tag, files := range tagIndex {
		entries = append(entries, TagEntry{
			Name:  tag,
			Count: len(files),
			Files: files,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Name < entries[j].Name
	})

	return TagList{
		entries: entries,
		cursor:  0,
		palette: palette,
	}
}

func (tl *TagList) MoveUp() {
	if tl.cursor > 0 {
		tl.cursor--
	}
}

func (tl *TagList) MoveDown() {
	if tl.cursor < len(tl.entries)-1 {
		tl.cursor++
	}
}

func (tl TagList) SelectedTag() string {
	if tl.cursor >= 0 && tl.cursor < len(tl.entries) {
		return tl.entries[tl.cursor].Name
	}
	return ""
}

func (tl TagList) SelectedFiles() []string {
	if tl.cursor >= 0 && tl.cursor < len(tl.entries) {
		return tl.entries[tl.cursor].Files
	}
	return nil
}

func (tl TagList) Count() int {
	return len(tl.entries)
}

func (tl TagList) View() string {
	if len(tl.entries) == 0 {
		return lipgloss.NewStyle().
			Foreground(tl.palette.TextMuted).
			Render("  No tags found")
	}

	var sb strings.Builder
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(tl.palette.Accent).
		Render(fmt.Sprintf("  Tags (%d)", len(tl.entries)))
	sb.WriteString(header)
	sb.WriteString("\n")

	for i, entry := range tl.entries {
		line := fmt.Sprintf("  #%-20s (%d)", entry.Name, entry.Count)
		if i == tl.cursor {
			line = lipgloss.NewStyle().
				Background(tl.palette.Accent).
				Foreground(tl.palette.SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(tl.palette.TextSecondary).
				Render(line)
		}
		sb.WriteString(line)
		if i < len(tl.entries)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
