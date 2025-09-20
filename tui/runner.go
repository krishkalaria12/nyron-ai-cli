package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/krishkalaria12/nyron-ai-cli/tui/components/streaming"
)

// RunStreamingResponseModel starts the main streaming response TUI
func RunStreamingResponseModel() {
	p := tea.NewProgram(
		streaming.NewStreamingResponseModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}