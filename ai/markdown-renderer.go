package ai

import (
	markdown "github.com/MichaelMure/go-term-markdown"
)

func RenderToTerminalWithWidth(markdownString string, width int) (string, error) {
	result := markdown.Render(string(markdownString), width, 0)

	return string(result), nil
}
