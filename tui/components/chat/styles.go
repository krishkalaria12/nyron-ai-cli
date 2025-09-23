package chat

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for consistent theming
var (
	// Primary colors
	primaryColor   = lipgloss.Color("#6366f1") // Indigo
	secondaryColor = lipgloss.Color("#8b5cf6") // Violet
	accentColor    = lipgloss.Color("#06b6d4") // Cyan
	successColor   = lipgloss.Color("#10b981") // Emerald
	errorColor     = lipgloss.Color("#ef4444") // Red

	// Neutral colors
	textPrimary   = lipgloss.Color("#111827") // Gray-900
	textMuted     = lipgloss.Color("#9ca3af") // Gray-400
	borderColor   = lipgloss.Color("#d1d5db") // Gray-300
	borderFocused = primaryColor
)

var (
	// A master style for the entire application, providing a container.
	// FIX: Removed Margin(1, 0) to eliminate extra padding and horizontal shift.
	appStyle = lipgloss.NewStyle()

	// Loading spinner style
	loadingStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Input border styles with focus states
	inputBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Padding(0, 1)

	focusedInputBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(borderFocused).
				Padding(0, 1)

	// Message styles with better visual hierarchy
	userMessageStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	userMessageContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))

	aiMessageStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	aiMessageContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))

	// Header style - a slim, single-line bar
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D84797")).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1).
			Height(3)

	// Error and status styles
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor)

	thinkingStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true).
			PaddingLeft(2)

	// Thinking header style with greyer color
	thinkingHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6b7280")). // Gray-500 - greyer than textMuted
				Bold(true).
				PaddingLeft(2)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(textMuted)

	// Dialog styles for modal appearance
	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Tool calling styles
	toolCallStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			PaddingLeft(2)

	toolCallHeaderStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				PaddingLeft(2)

	toolCallContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				PaddingLeft(4).
				Border(lipgloss.Border{Left: "│"}).
				BorderForeground(accentColor)
)
