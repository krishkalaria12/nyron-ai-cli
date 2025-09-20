package components

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
	editor "github.com/krishkalaria12/nyron-ai-cli/tui/components/editor"
)

type focusState int

const (
	focusViewport focusState = iota
	focusInput
)

type keyMap struct {
	Tab    key.Binding
	Up     key.Binding
	Down   key.Binding
	Scroll key.Binding
	Enter  key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Tab, k.Enter, k.Quit},
	}
}

var keys = keyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
	Scroll: key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("scroll up", "scroll down")),
	Tab:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch focus")),
	Enter:  key.NewBinding(key.WithKeys("ctrl+enter"), key.WithHelp("ctrl+enter", "send message")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

type streamingResponseModel struct {
	content         string
	userSentMessage bool
	rawContent      string // Store raw markdown for re-rendering on resize
	loading         bool
	ready           bool
	viewport        viewport.Model
	spinner         spinner.Model
	input           editor.InputModel
	keys            keyMap
	focused         focusState
	help            help.Model
	width           int
	height          int
	prompt          string
	responseChan    chan ai.StreamMessage // producer -> forwarder
	forwardChan     chan ai.StreamMessage // forwarder -> tea loop
	err             error
	isStreaming     bool // Track if we're currently streaming
}

// Messages
type streamChunkMsg ai.StreamMessage
type startStreamMsg string

var (
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).MaxHeight(10)
	}()
)

func NewStreamingResponseModel() streamingResponseModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	// Initialize input model with editor keymap
	inputModel := editor.InitialInputModel()
	inputModel.InputKeys = editor.DefaultEditorKeyMap()

	return streamingResponseModel{
		loading:         false, // Don't start loading immediately
		spinner:         s,
		input:           inputModel, // Assign the input model
		focused:         focusInput, // Start focused on input
		keys:            keys,
		help:            help.New(),
		responseChan:    make(chan ai.StreamMessage),
		forwardChan:     make(chan ai.StreamMessage),
		isStreaming:     false,
		userSentMessage: false,
	}
}

func startStreamCommand(prompt string, responseChan chan ai.StreamMessage, forwardChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		go func() {
			ai.GeminiStreamAPI(prompt, responseChan)
		}()

		go func() {
			const (
				maxBatchDelay = 50 * time.Millisecond // debounce time
				maxBatchSize  = 1024                  // flush if buffer reaches this size
			)

			var b strings.Builder
			timer := time.NewTimer(time.Hour)
			timer.Stop()
			pending := false

			flush := func(done bool) {
				if !pending && !done {
					return
				}
				content := b.String()

				forwardChan <- ai.StreamMessage{
					Content: content,
					Error:   nil,
					Done:    done,
				}

				b.Reset()
				pending = false
			}

			for {
				select {
				case ev, ok := <-responseChan:
					if !ok {
						flush(true)
						close(forwardChan)
						return
					}
					if ev.Error != nil {
						forwardChan <- ai.StreamMessage{
							Content: "",
							Error:   ev.Error,
							Done:    true,
						}
						close(forwardChan)
						return
					}

					if ev.Content != "" {
						b.WriteString(ev.Content)
						pending = true
					}

					if b.Len() >= maxBatchSize {
						flush(false)
						if !timer.Stop() {
							select {
							case <-timer.C:
							default:
							}
						}
					} else if pending {
						if !timer.Stop() {
							select {
							case <-timer.C:
							default:
							}
						}
						timer.Reset(maxBatchDelay)
					}

					if ev.Done {
						flush(true)
						close(forwardChan)
						return
					}

				case <-timer.C:
					flush(false)
				}
			}
		}()

		return startStreamMsg(prompt)
	}
}

func waitForForwardMessage(forwardChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-forwardChan
		if !ok {
			return streamChunkMsg(ai.StreamMessage{
				Content: "",
				Error:   nil,
				Done:    true,
			})
		}
		return streamChunkMsg(msg)
	}
}

// renderContentWithCurrentWidth renders the accumulated rawContent using the ai/glamour renderer
func (m *streamingResponseModel) renderContentWithCurrentWidth() string {
	// If no raw markdown accumulated yet, return whatever content we have
	if m.rawContent == "" {
		return m.content
	}

	// Calculate appropriate width for markdown rendering
	renderWidth := m.width - 4 // Leave padding for viewport borders
	if renderWidth < 20 {
		renderWidth = 20 // Minimum reasonable width
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

func (m streamingResponseModel) Init() tea.Cmd {
	// Initialize input focus
	return m.input.Focus()
}

func (m streamingResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = m.width

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.renderContentWithCurrentWidth())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
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

func (m streamingResponseModel) View() string {
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

func (m streamingResponseModel) headerView() string {
	title := "Nyron AI - Your AI Based Terminal"
	titleStyled := titleStyle.Render(title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(titleStyled)))
	return lipgloss.JoinHorizontal(lipgloss.Center, titleStyled, line)
}

func (m streamingResponseModel) footerView() string {
	helpView := m.help.View(m.keys)

	var inputStyle lipgloss.Style
	if m.focused == focusInput {
		inputStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("63"))
	} else {
		inputStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false, false, false, true)
	}

	return lipgloss.JoinVertical(lipgloss.Left, inputStyle.Render(m.input.View()), helpView)
}

// Run helper remains the same
func RunStreamingResponseModel() {
	p := tea.NewProgram(
		NewStreamingResponseModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// small helper
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
