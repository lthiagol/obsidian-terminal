package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

// TickMsg is sent every second by the timer to check for vault changes.
type TickMsg struct{}

// Mode represents the current TUI mode.
type Mode int

const (
	ModeBrowse Mode = iota
	ModeView
	ModeSearch
	ModeFind
	ModeHelp
)

func (m Mode) String() string {
	switch m {
	case ModeBrowse:
		return "BROWSE"
	case ModeView:
		return "VIEW"
	case ModeSearch:
		return "SEARCH"
	case ModeFind:
		return "FIND"
	case ModeHelp:
		return "HELP"
	default:
		return "???"
	}
}

// Model is the top-level Bubble Tea model for the TUI.
type Model struct {
	mode     Mode
	prevMode Mode

	vault       *VaultEntry
	activeNote  *VaultNote
	searchIndex map[string]string

	keys     KeyMap
	fileTree FileTree
	viewer   MarkdownViewer

	searchState search.State
	searchStyle search.Style
	allPaths    []string

	width     int
	height    int
	treeWidth int

	config  *Config
	ready   bool
	quitting bool

	err        error
	scanErrors []string
	palette    Palette

	helpScroll int

	lastRootModTime time.Time
	lastRescan      time.Time
	toasts          []Toast
}

// NewModel creates a Model by scanning the vault at cfg.VaultPath.
func NewModel(cfg *Config) Model {
	keys := DefaultKeys()

	skipDirs := cfg.SkipDirs
	if len(skipDirs) == 0 {
		skipDirs = DefaultConfig().SkipDirs
	}

	palette, err := lookupPalette(cfg.Theme)
	if err != nil {
		palette = newDarkPalette()
	}

	info, err := os.Stat(cfg.VaultPath)
	if err != nil {
		return Model{
			config: cfg,
			keys:   keys,
			err:    fmt.Errorf("vault path not accessible: %w", err),
		}
	}
	if !info.IsDir() {
		return Model{
			config: cfg,
			keys:   keys,
			err:    fmt.Errorf("vault path is not a directory: %s", cfg.VaultPath),
		}
	}

	tree, searchIndex, scanErrors, err := ScanVault(cfg.VaultPath, skipDirs)
	if err != nil {
		return Model{
			config: cfg,
			keys:   keys,
			err:    fmt.Errorf("scanning vault: %w", err),
		}
	}

	paths := allPaths(tree)

	return Model{
		mode:        ModeBrowse,
		prevMode:    ModeBrowse,
		vault:       tree,
		searchIndex: searchIndex,
		allPaths:    paths,
		keys:        keys,
		config:      cfg,
		fileTree:    NewFileTree(tree),
		viewer:      NewViewer(markdownStyleFrom(palette)),
		searchStyle: searchStyleFrom(palette),
		scanErrors:  scanErrors,
		palette:     palette,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("obsidian-terminal"),
		tea.EnterAltScreen,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.treeWidth = max(msg.Width/4, 25)
		if m.treeWidth < 5 {
			m.treeWidth = 5
		}
		m.fileTree.SetSize(m.treeWidth, m.height-1)
		viewerWidth := m.width - m.treeWidth - 2
		if viewerWidth < 10 {
			viewerWidth = 10
		}
		m.viewer.SetSize(viewerWidth, m.height-1)
		m.ready = true
		return m, nil

	case TickMsg:
		m.expireToasts()
		m.checkVaultChanges()
		return m, tickCmd()

	case tea.KeyMsg:
		if m.quitting {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyCtrlR:
			m.rescanVault()
			return m, nil
		}

		if MatchRune(msg, m.keys.QuitRune) || MatchRune(msg, 'Q') {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.mode {
		case ModeBrowse:
			return m.handleBrowseKey(msg)
		case ModeView:
			return m.handleViewKey(msg)
		case ModeSearch:
			return m.handleSearchKey(msg)
		case ModeFind:
			return m.handleFindKey(msg)
		case ModeHelp:
			return m.handleHelpKey(msg)
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if !m.ready {
		return "Loading..."
	}

	if m.width < 60 || m.height < 15 {
		return lipgloss.NewStyle().
			Foreground(Warning).
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Terminal too small — please resize")
	}

	if m.quitting {
		return ""
	}

	var rightPanel string
	switch m.mode {
	case ModeSearch:
		rightPanel = m.renderSearch()
	case ModeFind:
		rightPanel = m.renderFind()
	case ModeHelp:
		rightPanel = m.renderHelp()
	case ModeView:
		rightPanel = m.viewer.View()
	default:
		rightPanel = "Select a file to view"
	}

	treePanel := m.fileTree.View()

	treeStyle := TreeStyle.Width(m.treeWidth).Height(m.height - 1)
	viewerStyle := ViewerStyle.Width(m.width - m.treeWidth - 1).Height(m.height - 1)

	leftP := treeStyle.Render(treePanel)
	rightP := viewerStyle.Render(rightPanel)

	main := lipgloss.JoinHorizontal(lipgloss.Top, leftP, rightP)

	statusBar := m.renderStatusBar()

	result := lipgloss.JoinVertical(lipgloss.Top, main, statusBar)

	if len(m.toasts) > 0 {
		toastText := m.renderToasts()
		result = lipgloss.JoinVertical(lipgloss.Bottom, result, toastText)
	}

	return result
}

func (m Model) renderSearch() string {
	return m.renderSearchPanel("fuzzy", "results")
}

func (m Model) renderFind() string {
	return m.renderSearchPanel("content", "matches")
}

func (m Model) renderSearchPanel(label, resultLabel string) string {
	var sb strings.Builder
	modeLabel := lipgloss.NewStyle().Bold(true).Foreground(AccentSecondary).Render(label)
	sb.WriteString(fmt.Sprintf("%s  %s_  (%d %s)", modeLabel, m.searchState.Query(), m.searchState.ResultCount(), resultLabel))
	sb.WriteString("\n\n")
	sb.WriteString(search.RenderResults(m.searchState, m.width-m.treeWidth-6, m.searchStyle))
	return sb.String()
}

func (m *Model) checkVaultChanges() {
	if time.Since(m.lastRescan) < 2*time.Second {
		return
	}

	info, err := os.Stat(m.config.VaultPath)
	if err != nil {
		m.addToast("Could not check vault: "+err.Error(), ToastWarning)
		return
	}

	if !info.ModTime().After(m.lastRootModTime) {
		return
	}

	m.rescanVault()
}

func (m *Model) rescanVault() {
	m.lastRescan = time.Now()

	info, err := os.Stat(m.config.VaultPath)
	if err != nil {
		return
	}
	m.lastRootModTime = info.ModTime()

	tree, searchIndex, scanErrors, err := ScanVault(m.config.VaultPath, m.config.SkipDirs)
	if err != nil {
		return
	}
	m.scanErrors = scanErrors

	oldActivePath := ""
	if m.activeNote != nil {
		oldActivePath = m.activeNote.Path
	}

	m.vault = tree
	m.searchIndex = searchIndex
	m.allPaths = allPaths(tree)
	m.fileTree = NewFileTree(tree)

	if oldActivePath != "" {
		note, err := LoadNote(m.config.VaultPath, oldActivePath)
		if err != nil {
			m.addToast("Note was deleted: "+oldActivePath, ToastWarning)
			m.mode = ModeBrowse
			m.activeNote = nil
		} else {
			m.activeNote = note
			m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
		}
	}
}

func countFiles(entry *VaultEntry) int {
	if entry == nil {
		return 0
	}
	count := 0
	for _, child := range entry.Children {
		if child.IsDir {
			count += countFiles(child)
		} else {
			count++
		}
	}
	return count
}

func truncateContent(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "...")
	}
	return strings.Join(lines, "\n")
}
