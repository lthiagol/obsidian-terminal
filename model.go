package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lthiagol/obsidian-terminal/internal/markdown"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

// TickMsg is sent every second by the timer to check for vault changes.
type TickMsg struct{}

// Mode represents the current TUI mode.
type Mode int

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
	scanErrors []string
	palette    Palette

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

	commandPaletteVisible bool
	commandPaletteQuery   string
	commandPaletteCursor  int
	commandPaletteResults []Command
}

// NewModel creates a Model by scanning the vault at cfg.VaultPath.
func NewModel(cfg *Config) Model {
	keys := DefaultKeys()

	skipDirs := cfg.SkipDirs
	if len(skipDirs) == 0 {
		skipDirs = DefaultConfig().SkipDirs
	}

	themeWarning := ""
	themeName := cfg.Theme
	if themeName == "" {
		themeName = "dark"
	}
	palette, err := lookupPalette(themeName)
	if err != nil {
		palette = newDarkPalette()
		themeWarning = "Unknown theme " + themeName + " — using dark"
	}

	// Apply custom theme overrides if present
	if cfg.CustomTheme != nil {
		customPalette, customErr := paletteFromCustom(cfg.CustomTheme, palette)
		if customErr != nil {
			if themeWarning != "" {
				themeWarning += "; "
			}
			themeWarning += customErr.Error()
		}
		palette = customPalette
	}

	activatePalette(palette)

	// If no vault path but profiles exist, enter picker mode
	if cfg.VaultPath == "" && len(cfg.Profiles) > 0 {
		m := Model{
			mode:          ModeProfilePicker,
			prevMode:      ModeProfilePicker,
			config:        cfg,
			keys:          keys,
			palette:       palette,
			profilePicker: NewProfilePicker(cfg.Profiles),
		}
		if themeWarning != "" {
			m.addToast(themeWarning, ToastWarning)
		}
		return m
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
		fileTree:        NewFileTree(tree),
		viewer:          NewViewer(markdownStyleFrom(palette, cfg.LineSpacing)),
		searchStyle:     searchStyleFrom(palette),
		scanErrors:      scanErrors,
		palette:         palette,
		activePinnedIdx: -1,
		profilePicker:   NewProfilePicker(cfg.Profiles),
	}
	if themeWarning != "" {
		m.addToast(themeWarning, ToastWarning)
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
			m.toggleRecents()
			return m, nil
		case tea.KeyCtrlK:
			m.openCommandPalette()
			return m, nil
		}

		if (m.mode == ModeBrowse || m.mode == ModeView) && !m.commandPaletteVisible && !m.recentVisible && !m.outlineVisible && (MatchRune(msg, m.keys.QuitRune) || MatchRune(msg, 'Q')) {
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
	if m.commandPaletteVisible {
		rightPanel = m.renderCommandPalette()
	} else if m.recentVisible {
		rightPanel = m.renderRecents()
	} else {
		switch m.mode {
		case ModeSearch:
			rightPanel = m.renderSearch()
		case ModeFind:
			rightPanel = m.renderFind()
		case ModeHelp:
			rightPanel = m.renderHelp()
		case ModeTags:
			rightPanel = m.tagList.View()
		case ModeProfilePicker:
			rightPanel = m.profilePicker.View()
		case ModeView:
			if m.outlineVisible {
				rightPanel = m.renderOutline()
			} else if m.backlinkMode {
				viewerHeight := (m.height - 1) * 7 / 10
				backlinkHeight := m.height - 1 - viewerHeight - 1
				viewerStyle := ViewerStyle.Width(m.width - m.treeWidth - 1).Height(viewerHeight)
				backlinkStyle := lipgloss.NewStyle().
					Border(lipgloss.NormalBorder(), true, false, false, false).
					BorderForeground(Accent).
					Width(m.width - m.treeWidth - 1).
					Height(backlinkHeight)
				rightPanel = lipgloss.JoinVertical(lipgloss.Left,
					viewerStyle.Render(m.viewer.View()),
					backlinkStyle.Render(m.backlinkPanel.View()),
				)
			} else {
				rightPanel = m.viewer.View()
			}
		default:
			rightPanel = "Select a file to view"
		}
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

	tree, indexes, scanErrors, err := ScanVault(m.config.VaultPath, m.config.SkipDirs)
	if err != nil {
		return
	}
	m.scanErrors = scanErrors

	oldActivePath := ""
	if m.activeNote != nil {
		oldActivePath = m.activeNote.Path
	}

	m.vault = tree
	m.searchIndex = indexes.Search
	m.backlinkIndex = indexes.Backlinks
	m.tagIndex = indexes.Tags
	m.allPaths = allPaths(tree)
	m.fileTree = NewFileTree(tree)
	m.validatePins()

	if oldActivePath != "" {
		note, err := LoadNote(m.config.VaultPath, oldActivePath)
		if err != nil {
			m.addToast("Note was deleted: "+oldActivePath, ToastWarning)
			m.mode = ModeBrowse
			m.activeNote = nil
		} else {
			m.openNote(note.Path)
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

func (m *Model) togglePin(path string) {
	if path == "" {
		return
	}

	for i, pin := range m.pinnedNotes {
		if pin.Path == path {
			m.pinnedNotes = append(m.pinnedNotes[:i], m.pinnedNotes[i+1:]...)
			if m.activePinnedIdx >= len(m.pinnedNotes) {
				m.activePinnedIdx = len(m.pinnedNotes) - 1
			}
			m.addToast("Unpinned note", ToastInfo)
			return
		}
	}

	scrollY := 0
	if m.activeNote != nil && m.activeNote.Path == path {
		scrollY = m.viewer.GetScrollPosition()
	}

	m.pinnedNotes = append(m.pinnedNotes, PinnedNote{Path: path, ScrollY: scrollY})
	m.addToast("Pinned note", ToastInfo)
}

func (m *Model) openPinnedNote(index int) {
	if index < 0 || index >= len(m.pinnedNotes) {
		return
	}

	pin := m.pinnedNotes[index]

	// Validate pin still exists before opening
	_, err := LoadNote(m.config.VaultPath, pin.Path)
	if err != nil {
		m.addToast("Pinned note deleted: "+err.Error(), ToastError)
		m.pinnedNotes = append(m.pinnedNotes[:index], m.pinnedNotes[index+1:]...)
		if m.activePinnedIdx >= len(m.pinnedNotes) {
			m.activePinnedIdx = len(m.pinnedNotes) - 1
		}
		return
	}

	m.openNote(pin.Path)
	m.viewer.SetScrollPosition(pin.ScrollY)
	m.activePinnedIdx = index
}

func (m *Model) cyclePinnedNext() {
	if len(m.pinnedNotes) == 0 {
		m.addToast("No pinned notes", ToastWarning)
		return
	}

	if m.activePinnedIdx >= 0 && m.activePinnedIdx < len(m.pinnedNotes) && m.activeNote != nil {
		m.pinnedNotes[m.activePinnedIdx].ScrollY = m.viewer.GetScrollPosition()
	}

	m.activePinnedIdx++
	if m.activePinnedIdx >= len(m.pinnedNotes) {
		m.activePinnedIdx = 0
	}

	m.openPinnedNote(m.activePinnedIdx)
}

func (m *Model) cyclePinnedPrev() {
	if len(m.pinnedNotes) == 0 {
		m.addToast("No pinned notes", ToastWarning)
		return
	}

	if m.activePinnedIdx >= 0 && m.activePinnedIdx < len(m.pinnedNotes) && m.activeNote != nil {
		m.pinnedNotes[m.activePinnedIdx].ScrollY = m.viewer.GetScrollPosition()
	}

	m.activePinnedIdx--
	if m.activePinnedIdx < 0 {
		m.activePinnedIdx = len(m.pinnedNotes) - 1
	}

	m.openPinnedNote(m.activePinnedIdx)
}

func (m *Model) validatePins() {
	var valid []PinnedNote
	for _, pin := range m.pinnedNotes {
		path := filepath.Join(m.config.VaultPath, pin.Path)
		if _, err := os.Stat(path); err == nil {
			valid = append(valid, pin)
		}
	}
	m.pinnedNotes = valid
	if m.activePinnedIdx >= len(m.pinnedNotes) {
		m.activePinnedIdx = len(m.pinnedNotes) - 1
	}
}

func (m *Model) buildOutline() {
	if m.activeNote == nil {
		m.outlineItems = nil
		return
	}

	lines := markdown.ParseMarkdown(m.activeNote.RawBody)
	headings := markdown.ExtractHeadings(lines)

	m.outlineItems = make([]OutlineItem, len(headings))
	for i, h := range headings {
		m.outlineItems[i] = OutlineItem{
			Level:   h.Level,
			Text:    h.Text,
			LineIdx: h.LineIdx,
			YOffset: estimateYOffset(lines, h.LineIdx, m.viewer.viewport.Width),
		}
	}

	m.outlineCursor = 0
}

func (m Model) renderOutline() string {
	if len(m.outlineItems) == 0 {
		return lipgloss.NewStyle().
			Foreground(TextMuted).
			Render("  No headings in this note")
	}

	var sb strings.Builder
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(Accent).
		Render(fmt.Sprintf("  Outline (%d)", len(m.outlineItems)))
	sb.WriteString(header)
	sb.WriteString("\n")

	for i, item := range m.outlineItems {
		indent := strings.Repeat("  ", item.Level-1)
		line := fmt.Sprintf("%s%s", indent, item.Text)

		if i == m.outlineCursor {
			line = lipgloss.NewStyle().
				Background(Accent).
				Foreground(SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Render(line)
		}

		sb.WriteString(line)
		if i < len(m.outlineItems)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func estimateYOffset(lines []markdown.MarkdownLine, targetIdx, width int) int {
	yOffset := 0
	for i := 0; i < targetIdx && i < len(lines); i++ {
		line := lines[i]
		switch line.BlockType {
		case markdown.BlockEmpty:
			yOffset++
		case markdown.BlockHeading:
			yOffset++
		case markdown.BlockCodeBlock:
			codeLines := strings.Count(line.RawContent, "\n") + 1
			yOffset += codeLines + 2
		case markdown.BlockList:
			yOffset++
		case markdown.BlockBlockquote:
			yOffset++
		case markdown.BlockCallout:
			yOffset++
		case markdown.BlockHorizontalRule:
			yOffset++
		default:
			text := markdown.RenderSegmentsPlain(line.Segments)
			if width > 0 {
				wrappedLines := (len(text) / width) + 1
				yOffset += wrappedLines
			} else {
				yOffset++
			}
		}
	}
	return yOffset
}

func (m *Model) buildDailyNotePath() string {
	now := time.Now()
	dateStr := now.Format(m.config.DailyNotesFormat)
	return filepath.Join(m.config.DailyNotesDir, dateStr+".md")
}

func (m *Model) openDailyNote() {
	path := m.buildDailyNotePath()
	note, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		dateStr := time.Now().Format(m.config.DailyNotesFormat)
		m.activeNote = &VaultNote{
			Path:  path,
			Title: "Daily: " + dateStr,
			Body:  "",
		}
		m.prevMode = m.mode
		m.mode = ModeView
		m.viewer.SetContent(m.activeNote.Body, m.width-m.treeWidth-2)
		m.buildOutline()
		m.addRecentNote(path)
		return
	}
	m.openNote(note.Path)
}

func (m *Model) addRecentNote(path string) {
	if path == "" {
		return
	}

	for i, recent := range m.recentNotes {
		if recent == path {
			m.recentNotes = append(m.recentNotes[:i], m.recentNotes[i+1:]...)
			break
		}
	}

	m.recentNotes = append([]string{path}, m.recentNotes...)

	if len(m.recentNotes) > 50 {
		m.recentNotes = m.recentNotes[:50]
	}
}

func (m *Model) toggleRecents() {
	if m.recentVisible {
		m.recentVisible = false
	} else {
		m.recentVisible = true
		m.recentCursor = 0
	}
}

func (m *Model) openRecentNote(index int) {
	if index < 0 || index >= len(m.recentNotes) {
		return
	}

	path := m.recentNotes[index]
	// Validate note still exists
	_, err := LoadNote(m.config.VaultPath, path)
	if err != nil {
		m.addToast("Failed to load recent note: "+err.Error(), ToastError)
		m.recentNotes = append(m.recentNotes[:index], m.recentNotes[index+1:]...)
		if m.recentCursor >= len(m.recentNotes) {
			m.recentCursor = len(m.recentNotes) - 1
		}
		return
	}

	m.openNote(path)
	m.recentVisible = false
}

func (m Model) renderRecents() string {
	if len(m.recentNotes) == 0 {
		return lipgloss.NewStyle().
			Foreground(TextMuted).
			Render("  No recent notes")
	}

	var sb strings.Builder
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(Accent).
		Render(fmt.Sprintf("  Recent Notes (%d)", len(m.recentNotes)))
	sb.WriteString(header)
	sb.WriteString("\n")

	for i, path := range m.recentNotes {
		line := fmt.Sprintf("  %s", path)

		if i == m.recentCursor {
			line = lipgloss.NewStyle().
				Background(Accent).
				Foreground(SelectionText).
				Bold(true).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Render(line)
		}

		sb.WriteString(line)
		if i < len(m.recentNotes)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
