package util

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
)

// Background markdown rendering
func RenderMarkdownAsync(content string, width int, messageIndex int) tea.Cmd {
	return func() tea.Msg {
		rendered, err := ai.RenderToTerminalWithWidth(content, width)
		if err != nil {
			// If markdown rendering fails, use plain text
			rendered = content
		}
		return MarkdownRenderedMsg{
			MessageIndex: messageIndex,
			Rendered:     rendered,
		}
	}
}
