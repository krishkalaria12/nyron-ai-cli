package streaming

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

// small helper
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}