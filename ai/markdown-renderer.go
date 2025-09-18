// ai/markdown_renderer.go
package ai

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/krishkalaria12/nyron-ai-cli/theme"
)

var GlamourOptions = []glamour.TermRendererOption{
	glamour.WithAutoStyle(),
	glamour.WithStyles(theme.MarkdownTheme()),
	glamour.WithEmoji(),
	glamour.WithChromaFormatter("terminal16m"),
	glamour.WithPreservedNewLines(),
}

func RenderToTerminalWithWidth(markdown string, width int) (string, error) {
	markdown = sanitizePartialMarkdown(markdown)
	opts := append([]glamour.TermRendererOption{}, GlamourOptions...)
	if width > 0 {
		opts = append(opts, glamour.WithWordWrap(width))
	}
	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return "", err
	}
	return r.Render(markdown)
}

func sanitizePartialMarkdown(s string) string {
	if s == "" {
		return s
	}

	backtickFenceCount := countFenceOccurrences(s, "```")
	tildeFenceCount := countFenceOccurrences(s, "~~~")

	if backtickFenceCount%2 == 1 {
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		s += "```"
	}

	if tildeFenceCount%2 == 1 {
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		s += "~~~"
	}

	if strings.HasSuffix(strings.TrimRight(s, " \t\n\r"), "<") {
		s += ">"
	}

	return s
}

func countFenceOccurrences(s, fence string) int {
	pattern := regexp.MustCompile("(?m)^\\s*" + regexp.QuoteMeta(fence))
	matches := pattern.FindAllStringIndex(s, -1)
	return len(matches)
}
