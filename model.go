package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

// TickMsg is sent every second by the timer to check for vault changes.
type TickMsg struct{}

// Mode represents the current TUI mode.
type Mode int

// VaultState tracks the health of the vault connection.
type VaultState int

const (
	VaultStateOK      VaultState = iota // VaultStateOK indicates the vault is fully accessible.
	VaultStatePartial                   // VaultStatePartial indicates some files/dirs failed to scan.
	VaultStateBroken                    // VaultStateBroken indicates the vault is inaccessible.
)


// PinnedNote represents a pinned note with saved scroll position.
type PinnedNote struct {
	Path    string
	ScrollY int
}

// OutlineItem represents a heading in the outline.
type OutlineItem struct {
	Level   int
	Text    string
	LineIdx int
	YOffset int
}

const (
	ModeBrowse Mode = iota
	ModeView
	ModeSearch
	ModeFind
	ModeHelp
	ModeTags
	ModeProfilePicker
	ModeGraph
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
	case ModeTags:
		return "TAGS"
	case ModeProfilePicker:
		return "PROFILES"
	case ModeGraph:
		return "GRAPH"
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
	treeRatio float64

	dragSplit bool

	config  *Config
	ready   bool
	quitting bool

	err        error
	scanErrors        []string
	scanErrorsVisible bool
	vaultState        VaultState
	palette           Palette

	helpScroll int

	lastRootModTime time.Time
	lastRescan      time.Time
	toasts          []Toast
	click           *clickState

	backlinkIndex map[string][]string
	backlinkPanel BacklinkPanel
	backlinkMode  bool

	tagIndex  map[string][]string
	tagList   TagList
	tagFilter string

	pinnedNotes     []PinnedNote
	activePinnedIdx int

	outlineVisible bool
	outlineItems   []OutlineItem
	outlineCursor  int

	recentNotes   []string
	recentVisible bool
	recentCursor  int

	profilePicker      ProfilePicker
	inNoteSearchActive bool
	inNoteSearchQuery  string
	inNoteSearchIdx    int
	inNoteMatches      []int
	history            []string
	historyForward     []string
	previewVisible     bool
	previewPath        string
	previewPane        PreviewPane

	commandPaletteVisible bool
	commandPaletteQuery   string
	commandPaletteCursor  int
	commandPaletteResults []Command

	graph GraphModel
}

// NewModel creates a Model by scanning the vault at cfg.VaultPath.
func NewModel(cfg *Config) Model {
	keys := DefaultKeys()

	validationWarnings := ValidateConfig(cfg)

	skipDirs := cfg.SkipDirs

	palette, _ := lookupPalette(cfg.Theme)
	if cfg.CustomTheme != nil {
		palette, _ = paletteFromCustom(cfg.CustomTheme, palette)
	}
	// If no vault path but profiles exist, enter picker mode
	if cfg.VaultPath == "" && len(cfg.Profiles) > 0 {
		m := Model{
			mode:          ModeProfilePicker,
			prevMode:      ModeProfilePicker,
			config:        cfg,
			keys:          keys,
			palette:       palette,
			profilePicker: NewProfilePicker(cfg.Profiles, palette),
		}
		if len(validationWarnings) > 0 {
			for _, w := range validationWarnings {
				m.addToast(w, ToastWarning)
			}
		}
		return m
	}

	info, err := os.Stat(cfg.VaultPath)
	if err != nil {
		suggestion := ""
		if os.IsNotExist(err) {
			suggestion = " — directory does not exist, create it first"
		} else {
			suggestion = " — check file permissions"
		}
		return Model{
			config: cfg,
			keys:   keys,
			err:    fmt.Errorf("vault path %q is not accessible%s: %w", cfg.VaultPath, suggestion, err),
		}
	}
	if !info.IsDir() {
		return Model{
			config: cfg,
			keys:   keys,
			err:    fmt.Errorf("vault path %q is not a directory", cfg.VaultPath),
		}
	}

	tree, indexes, scanErrors, err := ScanVault(cfg.VaultPath, skipDirs)
	if err != nil {
		return Model{
			config: cfg,
			keys:   keys,
			err:    fmt.Errorf("scanning vault: %w", err),
		}
	}

	paths := allPaths(tree)

	m := Model{
		mode:            ModeBrowse,
		prevMode:        ModeBrowse,
		vault:           tree,
		searchIndex:     indexes.Search,
		backlinkIndex:   indexes.Backlinks,
		tagIndex:        indexes.Tags,
		allPaths:        paths,
		keys:            keys,
		config:          cfg,
		fileTree:        NewFileTree(tree, palette),
		viewer:          NewViewer(markdownStyleFrom(palette, cfg.LineSpacing)),
		searchStyle:     searchStyleFrom(palette),
		scanErrors:      scanErrors,
		vaultState:      vaultStateFrom(len(scanErrors)),
		palette:         palette,
		activePinnedIdx: -1,
		profilePicker:   NewProfilePicker(cfg.Profiles, palette),
		previewPane:     NewPreviewPane(),
	}
	if len(validationWarnings) > 0 {
		for _, w := range validationWarnings {
			m.addToast(w, ToastWarning)
		}
	}
	restoreSession(&m)
	return m
}

func (m *Model) adjustTreeWidth(newWidth int) {
	if newWidth < 15 {
		newWidth = 15
	}
	if newWidth > m.width/2 {
		newWidth = m.width / 2
	}
	m.treeWidth = newWidth
	if m.width > 0 {
		m.treeRatio = float64(m.treeWidth) / float64(m.width)
	}
	m.fileTree.SetSize(m.treeWidth, m.height-1)
	viewerWidth := m.width - m.treeWidth - 2
	if viewerWidth < 10 {
		viewerWidth = 10
	}
	m.viewer.SetSize(viewerWidth, m.height-1)
	if m.activeNote != nil {
		m.viewer.SetContent(m.activeNote.Body, viewerWidth)
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

// Update implements tea.Model. Value receiver required by Bubble Tea interface.
// Model is large (~30 fields), but Bubble Tea handles the value copy efficiently
// because Updates return a new Model that replaces the old one via the event loop.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.treeRatio == 0 {
			m.treeRatio = 0.25
		}
		m.treeWidth = max(int(float64(m.width)*m.treeRatio), 25)
		m.fileTree.SetSize(m.treeWidth, m.height-1)
		viewerWidth := m.width - m.treeWidth - 2
		if viewerWidth < 10 {
			viewerWidth = 10
		}
		m.viewer.SetSize(viewerWidth, m.height-1)
		m.profilePicker.SetSize(m.width, m.height)
		if m.activeNote != nil {
			m.viewer.SetContent(m.activeNote.Body, viewerWidth)
		}
		m.ready = true
		return m, nil

	case TickMsg:
		m.expireToasts()
		m.checkVaultChanges()
		return m, tickCmd()

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		if m.quitting {
			saveSession(m)
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			saveSession(m)
			return m, tea.Quit
		case tea.KeyCtrlR:
			m.rescanVault()
			return m, nil
		case tea.KeyCtrlD:
			m.openDailyNote()
			return m, nil
		case tea.KeyCtrlO:
			if m.mode == ModeView {
				m.goBackHistory()
			} else {
				m.toggleRecents()
			}
			return m, nil
		case tea.KeyCtrlK:
			m.openCommandPalette()
			return m, nil
		}
		if msg.Type == m.keys.GraphToggle {
			m.enterGraphMode()
			return m, nil
		}

		if (m.mode == ModeBrowse || m.mode == ModeView) && !m.commandPaletteVisible && !m.recentVisible && !m.outlineVisible && !m.scanErrorsVisible && (MatchRune(msg, m.keys.QuitRune) || MatchRune(msg, 'Q')) {
			m.quitting = true
			saveSession(m)
			return m, tea.Quit
		}

		if m.commandPaletteVisible {
			return m.handleCommandPaletteKey(msg)
		}

		if m.recentVisible {
			return m.handleRecentsKey(msg)
		}

		if m.outlineVisible {
			return m.handleOutlineKey(msg)
		}

		if m.scanErrorsVisible {
			if msg.Type == tea.KeyEsc || MatchRune(msg, 'q') || MatchRune(msg, m.keys.Help) {
				m.scanErrorsVisible = false
				return m, nil
			}
			return m, nil
		}

		if m.inNoteSearchActive {
			return m.handleInNoteSearchKey(msg)
		}

		// Retry rescan when vault is broken
		if m.vaultState == VaultStateBroken && MatchRune(msg, 'r') {
			m.rescanVault()
			return m, nil
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
		case ModeTags:
			return m.handleTagsKey(msg)
		case ModeProfilePicker:
			return m.handleProfilePickerKey(msg)
		case ModeGraph:
			return m.handleGraphKey(msg)
		}
	}

	return m, nil
}

