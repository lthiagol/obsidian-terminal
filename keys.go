package main

import tea "github.com/charmbracelet/bubbletea"

// KeyMap holds vim-style and arrow key bindings.
type KeyMap struct {
	Up    []tea.KeyType // Up moves the cursor or scrolls up.
	Down  []tea.KeyType // Down moves the cursor or scrolls down.
	Left  []tea.KeyType // Left collapses a tree directory or navigates back.
	Right []tea.KeyType // Right expands a tree directory or navigates forward.

	UpRune    rune // UpRune is the vim-style key for moving up (k).
	DownRune  rune // DownRune is the vim-style key for moving down (j).
	LeftRune  rune // LeftRune is the vim-style key for moving left (h).
	RightRune rune // RightRune is the vim-style key for moving right (l).

	Enter    tea.KeyType // Enter confirms selection or opens a note.
	Esc      tea.KeyType // Esc cancels or returns to browse mode.
	Tab      tea.KeyType // Tab cycles between wiki-links.
	PageUp   tea.KeyType // PageUp scrolls up one page.
	PageDown tea.KeyType // PageDown scrolls down one page.
	CtrlC    tea.KeyType // CtrlC quits the application.

	QuitRune      rune // QuitRune quits the application (q).
	Search        rune // Search activates the fuzzy search bar (/).
	Find          rune // Find activates content search (s).
	Help          rune // Help shows the keybinding reference (?).
	TopRune       rune // TopRune jumps to the top (g).
	BottomRune    rune // BottomRune jumps to the bottom (G).
	PinRune       rune // PinRune pins the current note (p).
	Outline       rune // Outline shows the table of contents (t).
	ProfileSwitch rune // ProfileSwitch opens the profile picker (P).
	PreviewToggle rune // PreviewToggle toggles the note preview pane (v).

	CyclePinPrev []tea.KeyType // CyclePinPrev cycles to the previous pinned note.
	CyclePinNext []tea.KeyType // CyclePinNext cycles to the next pinned note.

	ShrinkTree tea.KeyType // ShrinkTree decreases the tree panel width.
	GrowTree   tea.KeyType // GrowTree increases the tree panel width.
	ResetTree  tea.KeyType // ResetTree restores the default tree panel width.
}

// MatchKey reports whether msg matches any of the given key types.
func MatchKey(msg tea.KeyMsg, keys []tea.KeyType) bool {
	for _, k := range keys {
		if msg.Type == k {
			return true
		}
	}
	return false
}

// MatchRune reports whether msg contains the given rune.
func MatchRune(msg tea.KeyMsg, r rune) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	for _, msgRune := range msg.Runes {
		if msgRune == r {
			return true
		}
	}
	return false
}

// DefaultKeys returns the default vim+arrow key bindings.
func DefaultKeys() KeyMap {
	return KeyMap{
		Up:         []tea.KeyType{tea.KeyUp},
		Down:       []tea.KeyType{tea.KeyDown},
		Left:       []tea.KeyType{tea.KeyLeft},
		Right:      []tea.KeyType{tea.KeyRight},
		UpRune:     'k',
		DownRune:   'j',
		LeftRune:   'h',
		RightRune:  'l',
		Enter:      tea.KeyEnter,
		Esc:        tea.KeyEsc,
		Tab:        tea.KeyTab,
		PageUp:     tea.KeyPgUp,
		PageDown:   tea.KeyPgDown,
		CtrlC:      tea.KeyCtrlC,
		QuitRune:      'q',
		Search:        '/',
		Find:          's',
		Help:          '?',
		TopRune:       'g',
		BottomRune:    'G',
		PinRune:       'p',
		Outline:       't',
		ProfileSwitch: 'P',
		PreviewToggle: 'v',
		CyclePinPrev:  []tea.KeyType{tea.KeyCtrlOpenBracket},
		CyclePinNext:  []tea.KeyType{tea.KeyCtrlCloseBracket},
		ShrinkTree:    tea.KeyCtrlLeft,
		GrowTree:      tea.KeyCtrlRight,
		ResetTree:     tea.KeyCtrlBackslash,
	}
}
