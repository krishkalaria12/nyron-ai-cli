package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/krishkalaria12/nyron-ai-cli/tui/components/chat"
)

// RunChatModel starts the main chat TUI
func RunChatModel() {
	p := tea.NewProgram(
		chat.NewChatModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}