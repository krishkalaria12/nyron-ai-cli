package streaming

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
)

// renderContentWithCurrentWidth renders the accumulated rawContent using the ai/glamour renderer
func (m *StreamingResponseModel) renderContentWithCurrentWidth() string {
	// If no raw markdown accumulated yet, return whatever content we have
	if m.rawContent == "" {
		return m.content
	}

	// Calculate appropriate width for markdown rendering
	renderWidth := m.width - 4 // Leave padding for viewport borders
	if renderWidth < 20 {
		renderWidth = 20 // Minimum reasonable width
	}
	if m.width <= 0 {
		renderWidth = 76 // Default reasonable width when model width not set
	}

	rendered, err := ai.RenderToTerminalWithWidth(m.rawContent, renderWidth)
	if err != nil {
		// Fallback: if glamour rendering fails, show raw markdown (so user still sees text)
		return m.rawContent
	}

	// update cached content with rendered ANSI
	m.content = rendered
	return m.content
}

func (m StreamingResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update dimensions first
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = m.width

		// Calculate header and footer heights with updated dimensions
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		// Ensure minimum viewport height
		viewportHeight := msg.Height - verticalMarginHeight
		if viewportHeight < 3 {
			viewportHeight = 3
		}

		if !m.ready {
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.renderContentWithCurrentWidth())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
			// Re-render content with new width when viewport size changes
			m.viewport.SetContent(m.renderContentWithCurrentWidth())
		}

	case startStreamMsg:
		// Stream has started, begin loading
		m.loading = true
		m.isStreaming = true
		return m, tea.Batch(
			m.spinner.Tick,
			waitForForwardMessage(m.forwardChan),
		)

	case streamChunkMsg:
		// Convert incoming message to ai.StreamMessage semantics
		in := ai.StreamMessage(msg)

		// Error handling from producer
		if in.Error != nil {
			m.err = in.Error
			m.loading = false
			m.isStreaming = false
			// ensure viewport displays error
			if m.ready {
				m.viewport.SetContent(fmt.Sprintf("Error: %v\n", in.Error))
			}
			// Re-focus input for next query
			m.focused = focusInput
			m.input.TextArea.SetValue(m.prompt)
			cmd = m.input.Focus()
			return m, cmd
		}

		// If Done, mark loading false and final-render accumulated content
		if in.Done {
			m.loading = false
			m.isStreaming = false
			m.prompt = ""
			m.userSentMessage = false

			// Re-focus input for next query
			m.focused = focusInput
			cmd = m.input.Focus()
			return m, cmd
		}

		// Append the incoming chunk to the raw content (accumulate)
		m.rawContent += in.Content

		// Schedule next wait for message
		return m, waitForForwardMessage(m.forwardChan)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Tab):
			if m.focused == focusInput {
				m.focused = focusViewport
				m.input.Blur()
			} else {
				m.focused = focusInput
				cmd = m.input.Focus()
				cmds = append(cmds, cmd)
			}
		case key.Matches(msg, m.keys.Scroll):
			m.focused = focusViewport
			m.input.Blur()
		}

		// Handle focus-specific key events
		switch m.focused {
		case focusInput:
			// Check for send message key (Enter) press on the input.
			if key.Matches(msg, m.input.InputKeys.SendMessage) {
				inputValue := strings.TrimSpace(m.input.Value())
				if inputValue == "" {
					// Don't send empty messages
					return m, nil
				}

				// Don't start new stream if already streaming
				if m.isStreaming {
					return m, nil
				}

				// Store the prompt and clear input
				m.prompt = inputValue
				m.input.Reset()

				// Blur the input and switch focus to viewport during streaming
				m.input.Blur()
				m.focused = focusViewport

				m.isStreaming = true
				m.userSentMessage = true

				// Clear previous response content and set new content
				m.rawContent = ""
				m.content = ""

				// Create new channels for this stream
				m.responseChan = make(chan ai.StreamMessage)
				m.forwardChan = make(chan ai.StreamMessage)

				// Start the streaming command
				return m, startStreamCommand(m.prompt, m.responseChan, m.forwardChan)
			} else if key.Matches(msg, m.input.InputKeys.Newline) {
				// Handle newline explicitly
				m.input.TextArea, cmd = m.input.TextArea.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				m.input.TextArea, cmd = m.input.TextArea.Update(msg)
				cmds = append(cmds, cmd)
			}

		case focusViewport:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Handle keyboard and mouse events in the viewport (for non-key messages)
	if m.ready && m.focused == focusViewport {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}