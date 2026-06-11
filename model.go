package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg struct{}

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

type ToastType int

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

type Toast struct {
	Message string
	Type    ToastType
	TTL     time.Duration
	Created time.Time
}

type Model struct {
	mode     Mode
	prevMode Mode

	vault       *VaultEntry
	activeNote  *VaultNote
	searchIndex map[string]string

	keys     KeyMap
	fileTree FileTree
	viewer   MarkdownViewer

	searchState SearchState
	allPaths    []string

	width     int
	height    int
	treeWidth int

	config  *Config
	ready   bool
	quitting bool

	err error

	helpScroll int

	lastRootModTime time.Time
	lastRescan      time.Time
	toasts          []Toast
}

func NewModel(cfg *Config) Model {
	keys := DefaultKeys()

	skipDirs := cfg.SkipDirs
	if len(skipDirs) == 0 {
		skipDirs = DefaultConfig().SkipDirs
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

	tree, searchIndex, _, err := ScanVault(cfg.VaultPath, skipDirs)
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
		viewer:      NewViewer(),
	}
}

func (m Model) Init() tea.Cmd {
	info, err := os.Stat(m.config.VaultPath)
	if err == nil {
		m.lastRootModTime = info.ModTime()
	}
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
		m.fileTree.SetSize(m.treeWidth, m.height-1)
		m.viewer.SetSize(m.width-m.treeWidth-2, m.height-1)
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

func (m Model) handleBrowseKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case MatchRune(msg, m.keys.Search):
		m.prevMode = m.mode
		m.mode = ModeSearch
		m.searchState = NewSearchState(SearchName, m.allPaths, m.searchIndex)
		return m, nil
	case MatchRune(msg, m.keys.Find):
		m.prevMode = m.mode
		m.mode = ModeFind
		m.searchState = NewSearchState(SearchContent, m.allPaths, m.searchIndex)
		return m, nil
	case MatchRune(msg, m.keys.Help):
		m.prevMode = m.mode
		m.mode = ModeHelp
		m.helpScroll = 0
		return m, nil
	case msg.Type == tea.KeyEnter:
		entry := m.fileTree.SelectedEntry()
		if entry != nil {
			if entry.IsDir {
				m.fileTree.toggleExpand()
			} else {
				note, err := LoadNote(m.config.VaultPath, entry.Path)
				if err == nil {
					m.activeNote = note
					m.prevMode = m.mode
					m.mode = ModeView
					m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
				}
			}
		}
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.fileTree.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.fileTree.MoveUp()
		return m, nil
	case MatchKey(msg, m.keys.Left) || MatchRune(msg, m.keys.LeftRune):
		m.fileTree.collapse()
		return m, nil
	case MatchKey(msg, m.keys.Right) || MatchRune(msg, m.keys.RightRune):
		m.fileTree.expand()
		return m, nil
	case MatchRune(msg, m.keys.TopRune):
		m.fileTree.cursor = 0
		return m, nil
	case MatchRune(msg, m.keys.BottomRune):
		if len(m.fileTree.items) > 0 {
			m.fileTree.cursor = len(m.fileTree.items) - 1
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		m.activeNote = nil
		return m, nil
	case MatchRune(msg, m.keys.Search):
		m.prevMode = m.mode
		m.mode = ModeSearch
		m.searchState = NewSearchState(SearchName, m.allPaths, m.searchIndex)
		return m, nil
	case MatchRune(msg, m.keys.Find):
		m.prevMode = m.mode
		m.mode = ModeFind
		m.searchState = NewSearchState(SearchContent, m.allPaths, m.searchIndex)
		return m, nil
	case MatchRune(msg, m.keys.Help):
		m.prevMode = m.mode
		m.mode = ModeHelp
		m.helpScroll = 0
		return m, nil
	case msg.Type == tea.KeyTab:
		m.viewer.CycleLink()
		return m, nil
	case msg.Type == tea.KeyEnter:
		if m.viewer.SelectedLinkIndex() >= 0 {
			target := m.viewer.SelectedLinkPath()
			if target != "" {
				resolved := ResolveWikiLink(target, m.vault, m.config.VaultPath)
				if resolved != "" {
					note, err := LoadNote(m.config.VaultPath, resolved)
					if err == nil {
						m.activeNote = note
						m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
					}
				}
			}
		}
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.viewer.ScrollDown(1)
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.viewer.ScrollUp(1)
		return m, nil
	case MatchRune(msg, m.keys.TopRune):
		m.viewer.ScrollTop()
		return m, nil
	case MatchRune(msg, m.keys.BottomRune):
		m.viewer.ScrollBottom()
		return m, nil
	case msg.Type == tea.KeyPgUp:
		m.viewer.ScrollHalfPageUp()
		return m, nil
	case msg.Type == tea.KeyPgDown:
		m.viewer.ScrollHalfPageDown()
		return m, nil
	}
	return m, nil
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case msg.Type == tea.KeyBackspace:
		if len(m.searchState.query) > 0 {
			m.searchState.SetQuery(m.searchState.query[:len(m.searchState.query)-1])
		}
		return m, nil
	case msg.Type == tea.KeyRunes && len(msg.Runes) > 0:
		m.searchState.SetQuery(m.searchState.query + string(msg.Runes))
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.searchState.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.searchState.MoveUp()
		return m, nil
	case msg.Type == tea.KeyEnter:
		result := m.searchState.SelectedResult()
		if result != nil {
			note, err := LoadNote(m.config.VaultPath, result.Path)
			if err == nil {
				m.activeNote = note
				m.prevMode = ModeBrowse
				m.mode = ModeView
				m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleFindKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case msg.Type == tea.KeyBackspace:
		if len(m.searchState.query) > 0 {
			m.searchState.SetQuery(m.searchState.query[:len(m.searchState.query)-1])
		}
		return m, nil
	case msg.Type == tea.KeyRunes && len(msg.Runes) > 0:
		m.searchState.SetQuery(m.searchState.query + string(msg.Runes))
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.searchState.MoveDown()
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		m.searchState.MoveUp()
		return m, nil
	case msg.Type == tea.KeyEnter:
		result := m.searchState.SelectedResult()
		if result != nil {
			note, err := LoadNote(m.config.VaultPath, result.Path)
			if err == nil {
				m.activeNote = note
				m.prevMode = ModeBrowse
				m.mode = ModeView
				m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case MatchKey(msg, m.keys.Down) || MatchRune(msg, m.keys.DownRune):
		m.helpScroll++
		return m, nil
	case MatchKey(msg, m.keys.Up) || MatchRune(msg, m.keys.UpRune):
		if m.helpScroll > 0 {
			m.helpScroll--
		}
		return m, nil
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

func (m Model) renderViewer() string {
	if m.activeNote == nil {
		if m.mode == ModeBrowse {
			return "Select a file to view"
		}
		return ""
	}

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(Accent).Render(m.activeNote.Title))
	sb.WriteString("\n\n")
	sb.WriteString(truncateContent(m.activeNote.Body, m.height-5))
	return sb.String()
}

func (m Model) renderSearch() string {
	var sb strings.Builder
	modeLabel := lipgloss.NewStyle().Bold(true).Foreground(AccentSecondary).Render("fuzzy")
	sb.WriteString(fmt.Sprintf("%s  %s_  (%d results)", modeLabel, m.searchState.query, m.searchState.ResultCount()))
	sb.WriteString("\n\n")
	sb.WriteString(RenderSearchResults(m.searchState, m.width-m.treeWidth-6))
	return sb.String()
}

func (m Model) renderFind() string {
	var sb strings.Builder
	modeLabel := lipgloss.NewStyle().Bold(true).Foreground(AccentSecondary).Render("content")
	sb.WriteString(fmt.Sprintf("%s  %s_  (%d matches)", modeLabel, m.searchState.query, m.searchState.ResultCount()))
	sb.WriteString("\n\n")
	sb.WriteString(RenderSearchResults(m.searchState, m.width-m.treeWidth-6))
	return sb.String()
}

func (m Model) renderHelp() string {
	groups := []struct {
		title  string
		bindings []string
	}{
		{
			title: "Navigation",
			bindings: []string{
				"j / ↓  — move down",
				"k / ↑  — move up",
				"h / ←  — collapse / back",
				"l / →  — expand / forward",
				"g       — jump to top",
				"G       — jump to bottom",
				"PgUp    — page up",
				"PgDn    — page down",
			},
		},
		{
			title: "File Tree",
			bindings: []string{
				"Enter — open note / toggle folder",
				"← →   — collapse / expand folder",
			},
		},
		{
			title: "Viewer",
			bindings: []string{
				"j / k  — scroll down / up",
				"g / G  — top / bottom",
				"Tab    — cycle wiki-links",
				"Enter  — follow selected link",
				"h / Esc — back to browse",
			},
		},
		{
			title: "Search",
			bindings: []string{
				"/  — fuzzy file name search",
				"s  — full-text content search",
				"Enter — open selected result",
				"Esc   — cancel search",
			},
		},
		{
			title: "Global",
			bindings: []string{
				"?  — toggle this help",
				"q  — quit",
			},
		},
	}

	lines := []string{
		lipgloss.NewStyle().Bold(true).Foreground(Accent).Render("obsidian-terminal — Keybindings"),
		"",
	}

	for _, g := range groups {
		header := lipgloss.NewStyle().Bold(true).Foreground(Accent).Render(g.title)
		lines = append(lines, header)
		for _, b := range g.bindings {
			parts := strings.SplitN(b, "—", 2)
			key := lipgloss.NewStyle().Foreground(AccentSecondary).Render(strings.TrimSpace(parts[0]))
			var desc string
			if len(parts) > 1 {
				desc = lipgloss.NewStyle().Foreground(TextSecondary).Render("—" + parts[1])
			}
			lines = append(lines, "  "+key+"  "+desc)
		}
		lines = append(lines, "")
	}

	if m.helpScroll > len(lines)-1 {
		m.helpScroll = len(lines) - 1
	}

	start := m.helpScroll
	end := start + (m.height - 5)
	if end > len(lines) {
		end = len(lines)
	}
	if start >= len(lines) {
		start = 0
	}

	return strings.Join(lines[start:end], "\n")
}

func (m Model) renderStatusBar() string {
	modeColor := ModeColors[m.mode]
	modeBadge := lipgloss.NewStyle().
		Background(modeColor).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s ", m.mode.String()))

	var info string
	switch m.mode {
	case ModeBrowse:
		info = fmt.Sprintf("%d files", countFiles(m.vault))
	case ModeView:
		if m.activeNote != nil {
			info = truncatePath(m.activeNote.Path, m.width-60)
			if m.viewer.SelectedLinkIndex() >= 0 {
				info += " → " + m.viewer.SelectedLinkPath()
			}
		}
	case ModeSearch, ModeFind:
		info = m.searchState.query
	case ModeHelp:
		info = "j/k scroll | Esc back"
	}

	midSection := lipgloss.NewStyle().Foreground(TextSecondary).Padding(0, 1).Render(info)

	hints := modeHints(m.mode)
	hintSection := lipgloss.NewStyle().Foreground(TextDim).Padding(0, 1).Render(hints)

	modeWidth := lipgloss.Width(modeBadge)
	midWidth := max(0, m.width-modeWidth-lipgloss.Width(hintSection)-4)

	fullBar := lipgloss.JoinHorizontal(lipgloss.Center,
		modeBadge,
		lipgloss.NewStyle().Width(midWidth).Render(midSection),
		hintSection,
	)

	return StatusStyle.Width(m.width).Render(fullBar)
}

func modeHints(mode Mode) string {
	switch mode {
	case ModeBrowse:
		return "/ search | Enter open | ? help | q quit"
	case ModeView:
		return "h back | j/k scroll | Tab link | / search | ? help"
	case ModeSearch:
		return "type filter | Enter open | Esc cancel"
	case ModeFind:
		return "type search | Enter open | Esc cancel"
	case ModeHelp:
		return "j/k scroll | Esc back"
	default:
		return ""
	}
}

func truncatePath(path string, maxLen int) string {
	if maxLen < 5 {
		return "..."
	}
	if len(path) <= maxLen {
		return path
	}
	return ".../" + path[len(path)-maxLen+4:]
}

func (m *Model) checkVaultChanges() {
	if time.Since(m.lastRescan) < 2*time.Second {
		return
	}

	info, err := os.Stat(m.config.VaultPath)
	if err != nil {
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

	tree, searchIndex, _, err := ScanVault(m.config.VaultPath, m.config.SkipDirs)
	if err != nil {
		return
	}

	oldActivePath := ""
	if m.activeNote != nil {
		oldActivePath = m.activeNote.Path
	}

	m.vault = tree
	m.searchIndex = searchIndex
	m.allPaths = allPaths(tree)
	m.fileTree = NewFileTree(tree)

	if oldActivePath != "" {
		_, err := LoadNote(m.config.VaultPath, oldActivePath)
		if err != nil {
			m.addToast("Note was deleted: "+oldActivePath, ToastWarning)
			m.mode = ModeBrowse
			m.activeNote = nil
		} else if m.activeNote != nil {
			note, _ := LoadNote(m.config.VaultPath, oldActivePath)
			m.activeNote = note
			m.viewer.SetContent(note.Body, m.width-m.treeWidth-2)
		}
	}
}

func (m *Model) addToast(message string, t ToastType) {
	m.toasts = append(m.toasts, Toast{
		Message: message,
		Type:    t,
		TTL:     3 * time.Second,
		Created: time.Now(),
	})
}

func (m *Model) expireToasts() {
	var active []Toast
	for _, toast := range m.toasts {
		if time.Since(toast.Created) < toast.TTL {
			active = append(active, toast)
		}
	}
	m.toasts = active
}

func (m Model) renderToasts() string {
	var lines []string
	for _, toast := range m.toasts {
		lines = append(lines, renderToast(toast, m.width))
	}
	return strings.Join(lines, "\n")
}

func renderToast(toast Toast, width int) string {
	var icon string
	var borderColor lipgloss.Color
	switch toast.Type {
	case ToastInfo:
		icon = "i"
		borderColor = Info
	case ToastSuccess:
		icon = "v"
		borderColor = Success
	case ToastWarning:
		icon = "!"
		borderColor = Warning
	case ToastError:
		icon = "x"
		borderColor = Error
	}

	iconStyle := lipgloss.NewStyle().Foreground(borderColor).Bold(true)
	msgStyle := lipgloss.NewStyle().Foreground(TextSecondary)
	borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(borderColor)

	content := iconStyle.Render(" " + icon + " ") + msgStyle.Render(toast.Message)
	return borderStyle.Width(width).Padding(0, 1).Render(content)
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
