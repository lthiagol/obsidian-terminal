package main

func (m *Model) goBackHistory() {
	if len(m.history) == 0 {
		return
	}
	prev := m.history[len(m.history)-1]
	m.history = m.history[:len(m.history)-1]
	if m.activeNote != nil {
		m.historyForward = append(m.historyForward, m.activeNote.Path)
	}
	m.loadNote(prev, navHistory)
}

func (m *Model) goForwardHistory() {
	if len(m.historyForward) == 0 {
		return
	}
	next := m.historyForward[len(m.historyForward)-1]
	m.historyForward = m.historyForward[:len(m.historyForward)-1]
	if m.activeNote != nil {
		m.history = append(m.history, m.activeNote.Path)
	}
	m.loadNote(next, navHistory)
}
