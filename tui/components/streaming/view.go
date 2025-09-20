package streaming

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
)

func (m StreamingResponseModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	if !m.isStreaming && !m.loading {
		// final render of accumulated content (in case any sanitization appended fences)
		if m.rawContent != "" {
			final := m.renderContentWithCurrentWidth()
			if m.ready {
				m.viewport.SetContent(final)
				m.viewport.GotoBottom()
			}
		}
	}

	if m.userSentMessage {
		// Add the user's message to the viewport
		userMessage := fmt.Sprintf("\n**You:** %s\n\n**AI:** ", m.prompt)
		currentContent := ""
		if m.ready {
			currentContent = m.viewport.View()
		}
		newContent := currentContent + userMessage

		userMarkMsg := ""
		userMarkdownMsg, err := ai.RenderToTerminalWithWidth(string(newContent), m.width)

		if err != nil {
			userMarkMsg = newContent
		} else {
			userMarkMsg = userMarkdownMsg
		}

		if m.ready {
			m.viewport.SetContent(userMarkMsg)
			m.viewport.GotoBottom()
		}
	}

	// Render the accumulated markdown (prevents broken fences & flicker).
	rendered := m.renderContentWithCurrentWidth()
	// Update viewport
	if m.ready {
		m.viewport.SetContent(rendered)
		// Auto-scroll to bottom for live updates
		m.viewport.GotoBottom()
	}

	var viewportStyle lipgloss.Style
	if m.focused != focusViewport {
		viewportStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	}
	m.viewport.Style = viewportStyle

	// Get current viewport content
	viewportContent := m.viewport.View()

	// Add loading indicator if streaming
	if m.loading && m.isStreaming {
		viewportContent += "\n" + m.spinner.View() + " Thinking..."
	}

	// Create a styled viewport container
	styledViewport := viewportStyle.Render(viewportContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.headerView(),
		styledViewport,
		m.footerView(),
	)
}

func (m StreamingResponseModel) headerView() string {
	title := "Nyron AI - Your AI Based Terminal"
	titleStyled := titleStyle.Render(title)

	// Use model width instead of viewport width for reliable measurement
	width := m.width
	if width <= 0 {
		width = 80 // Default fallback width
	}

	line := strings.Repeat("─", max(0, width-lipgloss.Width(titleStyled)))
	return lipgloss.JoinHorizontal(lipgloss.Center, titleStyled, line)
}

func (m StreamingResponseModel) footerView() string {
	helpView := m.help.View(m.keys)

	var inputStyle lipgloss.Style
	if m.focused == focusInput {
		inputStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("63"))
	} else {
		inputStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false, false, false, true)
	}

	return lipgloss.JoinVertical(lipgloss.Left, inputStyle.Render(m.input.View()), helpView)
}