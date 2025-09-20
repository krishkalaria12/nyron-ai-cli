package chat

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m ChatModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Error display
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress Ctrl+C to quit", m.err)
	}

	// Header
	headerText := headerStyle.Render("💬 Nyron AI Chat")

	// Main content area
	content := m.viewport.View()

	// Input area
	// The width of the m.input.TextArea component is already being managed in the
	// model's Update function. We just let the border style wrap it.
	input := m.input.View()
	inputSection := inputBorderStyle.Render(input)

	// Help
	help := m.help.View(m.keys)

	// Combine all sections vertically.
	// The viewport's height is set in the model, ensuring it fills the
	// available space without pushing other components off-screen.
	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerText,
		content,
		inputSection,
		help,
	)
}
