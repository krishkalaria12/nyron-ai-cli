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
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		helpMsg := helpStyle.Render("Press Ctrl+C to quit")
		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Center, errorMsg, helpMsg))
	}

	// Dialog view
	if m.showDialog {
		dialog := dialogStyle.Render(m.modelDialog.View())
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dialog,
		)
	}

	// --- Main App View ---
	headerView := headerStyle.Width(m.width).Render("💬 Nyron AI Chat")
	viewportView := m.viewport.View()
	helpView := helpStyle.Width(m.width).Render(m.help.View(m.keys))

	var inputView string
	if m.focused == focusInput {
		inputView = focusedInputBorderStyle.Render(m.input.View())
	} else {
		inputView = inputBorderStyle.Render(m.input.View())
	}

	// Join all sections vertically.
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		viewportView,
		inputView,
		helpView,
	)

	return appStyle.Render(mainView)
}
