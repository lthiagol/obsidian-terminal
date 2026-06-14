package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ToastType classifies the severity of a toast notification.
type ToastType int

const toastTTL = 3 * time.Second

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

// Toast represents a temporary status bar notification.
type Toast struct {
	Message string
	Type    ToastType
	TTL     time.Duration
	Created time.Time
}

func (m *Model) addToast(message string, t ToastType) {
	m.toasts = append(m.toasts, Toast{
		Message: message,
		Type:    t,
		TTL:     toastTTL,
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
		lines = append(lines, renderToast(toast, m.width, m.palette))
	}
	return strings.Join(lines, "\n")
}

func renderToast(toast Toast, width int, p Palette) string {
	var icon string
	var borderColor lipgloss.Color
	switch toast.Type {
	case ToastInfo:
		icon = "\u2139" // ℹ
		borderColor = p.Info
	case ToastSuccess:
		icon = "\u2714" // ✔
		borderColor = p.Success
	case ToastWarning:
		icon = "\u26A0" // ⚠
		borderColor = p.Warning
	case ToastError:
		icon = "\u2716" // ✖
		borderColor = p.Error
	}

	iconStyle := lipgloss.NewStyle().Foreground(borderColor).Bold(true)
	msgStyle := lipgloss.NewStyle().Foreground(p.TextSecondary)
	borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(borderColor)

	content := iconStyle.Render(" " + icon + " ") + msgStyle.Render(toast.Message)
	return borderStyle.Width(width).Padding(0, 1).Render(content)
}
