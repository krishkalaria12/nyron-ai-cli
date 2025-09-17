package ai

import (
	"github.com/charmbracelet/glamour"
)

var TerminalRenderer *glamour.TermRenderer

var GlamourOptions = []glamour.TermRendererOption{
	glamour.WithAutoStyle(),
	glamour.WithStylePath("dark"),
	glamour.WithEmoji(),
	glamour.WithChromaFormatter("terminal256"),
	glamour.WithPreservedNewLines(),
}

func init() {
	var err error
	TerminalRenderer, err = glamour.NewTermRenderer(GlamourOptions...)
	if err != nil {
		panic(err)
	}
}

func RenderToTerminal(markdown string) (string, error) {
	return TerminalRenderer.Render(markdown)
}

func RenderToTerminalWithWidth(markdown string, width int) (string, error) {
	if width <= 0 {
		return TerminalRenderer.Render(markdown)
	}

	options := append([]glamour.TermRendererOption{}, GlamourOptions...)
	options = append(options, glamour.WithWordWrap(width))

	renderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		return "", err
	}

	return renderer.Render(markdown)
}
