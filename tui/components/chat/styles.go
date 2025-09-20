package chat

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Loading spinner style
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	// Input styles
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true)

	inputBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(0, 1)

	// Message styles
	userMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Bold(true)

	aiMessageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Bold(true)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("211")).
			Bold(true).
			Padding(0, 1).
			Width(0)
)
