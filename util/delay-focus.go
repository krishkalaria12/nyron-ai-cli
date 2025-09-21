package util

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Helper to create delayed focus message
func DelayedFocus() tea.Cmd {
	return tea.Tick(time.Millisecond*1500, func(time.Time) tea.Msg {
		return DelayedFocusMsg{}
	})
}
