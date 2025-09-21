package util

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	tea.Model
}

type DelayedFocusMsg struct{}

type MarkdownRenderedMsg struct {
	MessageIndex int
	Rendered     string
}
