package main

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

type Command struct {
	Name        string
	Description string
	Key         string
	Modes       []Mode
	Action      func(m *Model) (tea.Model, tea.Cmd)
}

func (m *Model) registerCommands() []Command {
	cmds := []Command{
		{
			Name:        "Fuzzy Search",
			Description: "Search files by name",
			Key:         "/",
			Modes:       []Mode{ModeBrowse, ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.enterSearchMode()
				return *m, nil
			},
		},
		{
			Name:        "Content Search",
			Description: "Search across all note contents",
			Key:         "s",
			Modes:       []Mode{ModeBrowse, ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.enterFindMode()
				return *m, nil
			},
		},
		{
			Name:        "Toggle Help",
			Description: "Show or hide keybindings help",
			Key:         "?",
			Modes:       []Mode{ModeBrowse, ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.enterHelpMode()
				return *m, nil
			},
		},
		{
			Name:        "Go Back",
			Description: "Return to previous mode",
			Key:         "Esc",
			Modes:       []Mode{ModeView, ModeSearch, ModeFind, ModeHelp, ModeTags},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				if m.mode == ModeView {
					m.activeNote = nil
				}
				m.mode = m.prevMode
				return *m, nil
			},
		},
		{
			Name:        "Follow Link",
			Description: "Follow selected wiki-link",
			Key:         "Enter",
			Modes:       []Mode{ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
		if m.viewer.SelectedLinkIndex() >= 0 {
				target := m.viewer.SelectedLinkPath()
				if target != "" && m.vault != nil {
					resolved := ResolveWikiLink(target, m.vault, m.config.VaultPath)
					if resolved != "" {
							m.openNote(resolved)
						}
					}
				}
				return *m, nil
			},
		},
		{
			Name:        "Cycle Wiki-Links",
			Description: "Cycle through wiki-links in note",
			Key:         "Tab",
			Modes:       []Mode{ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.viewer.CycleLink()
				return *m, nil
			},
		},
		{
			Name:        "Toggle Backlinks",
			Description: "Show notes linking to current note",
			Key:         "b",
			Modes:       []Mode{ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				if m.backlinkPanel.Count() > 0 {
					m.backlinkMode = !m.backlinkMode
				}
				return *m, nil
			},
		},
		{
			Name:        "Toggle Outline",
			Description: "Show table of contents for current note",
			Key:         "t",
			Modes:       []Mode{ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				if m.outlineVisible {
					m.outlineVisible = false
				} else {
					m.buildOutline()
					m.outlineVisible = true
				}
				return *m, nil
			},
		},
		{
			Name:        "Pin Note",
			Description: "Pin or unpin the current note",
			Key:         "p",
			Modes:       []Mode{ModeBrowse, ModeView},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				path := ""
				if m.activeNote != nil {
					path = m.activeNote.Path
				} else {
					entry := m.fileTree.SelectedEntry()
					if entry != nil && !entry.IsDir {
						path = entry.Path
					}
				}
				m.togglePin(path)
				return *m, nil
			},
		},
		{
			Name:        "Browse Tags",
			Description: "Browse and filter notes by tags",
			Key:         "T",
			Modes:       []Mode{ModeBrowse},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.enterTagsMode()
				return *m, nil
			},
		},
		{
			Name:        "Daily Note",
			Description: "Open today's daily note",
			Key:         "Ctrl+D",
			Modes:       nil, // Global
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.openDailyNote()
				return *m, nil
			},
		},
		{
			Name:        "Recent Notes",
			Description: "Show recently opened notes",
			Key:         "Ctrl+O",
			Modes:       nil, // Global
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.toggleRecents()
				return *m, nil
			},
		},
		{
			Name:        "Switch Profile",
			Description: "Switch to a different vault profile",
			Key:         "P",
			Modes:       []Mode{ModeBrowse},
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				if len(m.config.Profiles) > 0 {
					m.prevMode = m.mode
					m.mode = ModeProfilePicker
				}
				return *m, nil
			},
		},
		{
			Name:        "Force Rescan",
			Description: "Rescan vault for changes",
			Key:         "Ctrl+R",
			Modes:       nil, // Global
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.rescanVault()
				return *m, nil
			},
		},
		{
			Name:        "Quit",
			Description: "Exit obsidian-terminal",
			Key:         "q",
			Modes:       nil, // Global
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.quitting = true
				return *m, tea.Quit
			},
		},
	}

	if len(m.scanErrors) > 0 {
		cmds = append(cmds, Command{
			Name:        "Scan Errors",
			Description: fmt.Sprintf("View %d scan errors", len(m.scanErrors)),
			Key:         "",
			Modes:       nil, // Global
			Action: func(m *Model) (tea.Model, tea.Cmd) {
				m.showScanErrors()
				return *m, nil
			},
		})
	}

	return cmds
}

func (m *Model) openCommandPalette() {
	m.commandPaletteVisible = true
	m.commandPaletteQuery = ""
	m.commandPaletteCursor = 0
	m.commandPaletteResults = m.filterCommands("")
}

func (m *Model) filterCommands(q string) []Command {
	all := m.registerCommands()
	if q == "" {
		sort.Slice(all, func(i, j int) bool {
			return all[i].Name < all[j].Name
		})
		return all
	}

	var filtered []Command
	for _, cmd := range all {
		if m.commandMatches(cmd, q) {
			filtered = append(filtered, cmd)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		scoreI := search.FuzzyScore(q, filtered[i].Name, strings.ToLower(filtered[i].Name))
		scoreJ := search.FuzzyScore(q, filtered[j].Name, strings.ToLower(filtered[j].Name))
		return scoreI > scoreJ
	})

	return filtered
}

func (m *Model) commandMatches(cmd Command, q string) bool {
	qLower := strings.ToLower(q)
	nameMatch := strings.Contains(strings.ToLower(cmd.Name), qLower)
	descMatch := strings.Contains(strings.ToLower(cmd.Description), qLower)
	if !nameMatch && !descMatch {
		return false
	}
	if len(cmd.Modes) == 0 {
		return true
	}
	for _, mode := range cmd.Modes {
		if mode == m.mode {
			return true
		}
	}
	return false
}

func (m *Model) executeCommand(index int) (tea.Model, tea.Cmd) {
	if index < 0 || index >= len(m.commandPaletteResults) {
		return *m, nil
	}
	cmd := m.commandPaletteResults[index]
	m.commandPaletteVisible = false
	return cmd.Action(m)
}

func (m Model) renderCommandPalette() string {
	var sb strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.palette.Accent).
		Render("Command Palette")
	sb.WriteString(header)

	queryStyle := lipgloss.NewStyle().
		Foreground(m.palette.AccentSecondary).
		Background(lipgloss.Color("#222222"))
	queryText := m.commandPaletteQuery
	if queryText == "" {
		queryText = "Type to filter..."
	}
	sb.WriteString("\n")
	sb.WriteString(queryStyle.Render(" " + queryText + " "))
	sb.WriteString("\n\n")

	if len(m.commandPaletteResults) == 0 {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(m.palette.TextMuted).
			Render("  No matching commands"))
		return sb.String()
	}

	for i, cmd := range m.commandPaletteResults {
		line := fmt.Sprintf("  %-25s %-10s %s", cmd.Name, cmd.Key, cmd.Description)

		if i == m.commandPaletteCursor {
			line = lipgloss.NewStyle().
				Background(m.palette.Accent).
				Foreground(m.palette.SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(m.palette.TextSecondary).
				Render(line)
		}

		sb.WriteString(line)
		if i < len(m.commandPaletteResults)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m *Model) commandPaletteSearch() {
	m.commandPaletteResults = m.filterCommands(m.commandPaletteQuery)
	if m.commandPaletteCursor >= len(m.commandPaletteResults) {
		m.commandPaletteCursor = 0
	}
}
