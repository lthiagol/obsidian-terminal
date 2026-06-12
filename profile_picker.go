package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProfilePicker displays a list of vault profiles for selection.
type ProfilePicker struct {
	profiles []string
	cursor   int
	width    int
	height   int
}

// NewProfilePicker creates a ProfilePicker from a map of profiles.
func NewProfilePicker(profiles map[string]Profile) ProfilePicker {
	var names []string
	for name := range profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	return ProfilePicker{
		profiles: names,
		cursor:   0,
	}
}

// MoveUp moves the cursor up.
func (pp *ProfilePicker) MoveUp() {
	if pp.cursor > 0 {
		pp.cursor--
	}
}

// MoveDown moves the cursor down.
func (pp *ProfilePicker) MoveDown() {
	if pp.cursor < len(pp.profiles)-1 {
		pp.cursor++
	}
}

// Selected returns the currently selected profile name.
func (pp ProfilePicker) Selected() string {
	if pp.cursor >= 0 && pp.cursor < len(pp.profiles) {
		return pp.profiles[pp.cursor]
	}
	return ""
}

// Count returns the number of profiles.
func (pp ProfilePicker) Count() int {
	return len(pp.profiles)
}

// SetSize sets the picker dimensions.
func (pp *ProfilePicker) SetSize(width, height int) {
	pp.width = width
	pp.height = height
}

// View renders the profile picker.
func (pp ProfilePicker) View() string {
	if len(pp.profiles) == 0 {
		return lipgloss.NewStyle().
			Foreground(TextMuted).
			Width(pp.width).
			Height(pp.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No profiles defined")
	}

	var sb strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(Accent).
		Render("Select Vault Profile")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	// Profile list
	for i, name := range pp.profiles {
		var line string
		if i == pp.cursor {
			line = lipgloss.NewStyle().
				Background(Accent).
				Foreground(SelectionText).
				Bold(true).
				Render(fmt.Sprintf("  %s  ", name))
		} else {
			line = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Render(fmt.Sprintf("  %s", name))
		}
		sb.WriteString(line)
		if i < len(pp.profiles)-1 {
			sb.WriteString("\n")
		}
	}

	// Footer
	sb.WriteString("\n\n")
	footer := lipgloss.NewStyle().
		Foreground(TextDim).
		Render("↑/↓ navigate • Enter select • Esc quit")
	sb.WriteString(footer)

	content := sb.String()
	return lipgloss.NewStyle().
		Width(pp.width).
		Height(pp.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}
