package main

import tea "github.com/charmbracelet/bubbletea"

// KeyMap holds vim-style and arrow key bindings.
type KeyMap struct {
	Up    []tea.KeyType
	Down  []tea.KeyType
	Left  []tea.KeyType
	Right []tea.KeyType

	UpRune    rune
	DownRune  rune
	LeftRune  rune
	RightRune rune

	Enter    tea.KeyType
	Esc      tea.KeyType
	Tab      tea.KeyType
	PageUp   tea.KeyType
	PageDown tea.KeyType
	CtrlC    tea.KeyType

	QuitRune   rune
	Search     rune
	Find       rune
	Help       rune
	TopRune    rune
	BottomRune rune
	PinRune    rune

	CyclePinPrev []tea.KeyType
	CyclePinNext []tea.KeyType
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
		QuitRune:     'q',
		Search:       '/',
		Find:         's',
		Help:         '?',
		TopRune:      'g',
		BottomRune:   'G',
		PinRune:      'p',
		CyclePinPrev: []tea.KeyType{tea.KeyCtrlOpenBracket},
		CyclePinNext: []tea.KeyType{tea.KeyCtrlCloseBracket},
	}
}
