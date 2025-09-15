package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#EC4899")
	accentColor    = lipgloss.Color("#06B6D4")
	successColor   = lipgloss.Color("#10B981")
	warningColor   = lipgloss.Color("#F59E0B")
	errorColor     = lipgloss.Color("#EF4444")
	mutedColor     = lipgloss.Color("#6B7280")
	textColor      = lipgloss.Color("#F9FAFB")
	bgColor        = lipgloss.Color("#111827")
	borderColor    = lipgloss.Color("#374151")
)

var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Background(bgColor).
			Foreground(textColor)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(bgColor).
			Padding(1, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Align(lipgloss.Center)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			MarginBottom(1).
			PaddingLeft(2)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingRight(4).
			PaddingTop(1).
			PaddingBottom(1).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Foreground(textColor)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				PaddingRight(4).
				PaddingTop(1).
				PaddingBottom(1).
				MarginBottom(1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Background(primaryColor).
				Foreground(textColor).
				Bold(true)

	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(0).
			Foreground(textColor)

	SelectedListItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				PaddingRight(2).
				MarginBottom(0).
				Foreground(primaryColor).
				Bold(true)

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(2).
			Foreground(textColor).
			Background(bgColor)

	FocusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginTop(1).
				MarginBottom(2).
				Foreground(textColor).
				Background(bgColor).
				BorderStyle(lipgloss.ThickBorder())

	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1).
			PaddingLeft(4).
			Italic(true)

	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(2).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			BorderStyle(lipgloss.DoubleBorder())

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor)

	ContainerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	GradientTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7C3AED")).
				Background(lipgloss.Color("#1F2937")).
				Padding(2, 4).
				MarginBottom(2).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#7C3AED")).
				Align(lipgloss.Center)
)
