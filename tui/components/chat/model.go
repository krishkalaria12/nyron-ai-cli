package chat

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
	editor "github.com/krishkalaria12/nyron-ai-cli/tui/components/editor"
)

type focusState int

const (
	focusViewport focusState = iota
	focusInput
)

type keyMap struct {
	Tab      key.Binding
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Enter    key.Binding
	Quit     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.Tab, k.Enter, k.Quit},
	}
}

var keys = keyMap{
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
	PageUp:   key.NewBinding(key.WithKeys("pgup", "ctrl+u"), key.WithHelp("pgup", "page up")),
	PageDown: key.NewBinding(key.WithKeys("pgdown", "ctrl+d"), key.WithHelp("pgdn", "page down")),
	Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch focus")),
	Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send message")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

// Message represents a chat message
type Message struct {
	Content    string
	IsUser     bool
	Rendered   string // Cached markdown rendering
	IsRendered bool   // Whether markdown processing is complete
}

type ChatModel struct {
	messages []Message
	loading  bool
	viewport viewport.Model
	spinner  spinner.Model
	input    editor.InputModel
	keys     keyMap
	focused  focusState
	help     help.Model
	width    int
	height   int
	err      error
}

// Messages
type responseMsg struct {
	content string
	err     error
}

type delayedFocusMsg struct{}

type markdownRenderedMsg struct {
	messageIndex int
	rendered     string
}

// Helper to create delayed focus message
func delayedFocus() tea.Cmd {
	return tea.Tick(time.Millisecond*1500, func(time.Time) tea.Msg {
		return delayedFocusMsg{}
	})
}

// Background markdown rendering
func renderMarkdownAsync(content string, width int, messageIndex int) tea.Cmd {
	return func() tea.Msg {
		rendered, err := ai.RenderToTerminalWithWidth(content, width)
		if err != nil {
			// If markdown rendering fails, use plain text
			rendered = content
		}
		return markdownRenderedMsg{
			messageIndex: messageIndex,
			rendered:     rendered,
		}
	}
}

func NewChatModel() ChatModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	// Initialize input model with editor keymap
	inputModel := editor.InitialInputModel()
	inputModel.InputKeys = editor.DefaultEditorKeyMap()

	// Initialize viewport
	vp := viewport.New(80, 20)

	return ChatModel{
		messages: []Message{},
		loading:  false,
		spinner:  s,
		input:    inputModel,
		viewport: vp,
		focused:  focusInput, // Start focused on input
		keys:     keys,
		help:     help.New(),
	}
}

func (m ChatModel) Init() tea.Cmd {
	return m.input.Focus()
}

// getAIResponse sends a prompt to the AI and returns the response
func getAIResponse(prompt string) tea.Cmd {
	return func() tea.Msg {
		// For now, use Gemini API. You can add provider selection later
		response, err := ai.GeminiAPI(prompt)
		return responseMsg{
			content: response,
			err:     err,
		}
	}
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Only update if size actually changed
		if m.width != msg.Width || m.height != msg.Height {
			m.width = msg.Width
			m.height = msg.Height
			m.help.Width = msg.Width

			// ✨ CORRECTED: Calculate vertical space used by other components.
			// Header: 1 for text
			// Input: 1 for top border, 3 for text area, 1 for bottom border = 5
			// Help: 1 line
			verticalMargin := 1 + 5 + 1

			// Update viewport and input sizes
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargin

			// Set the width of the textarea component inside the input model.
			// The total horizontal space for the border and padding is 4 characters
			// (left border + left padding + right padding + right border).
			borderPadding := 4
			newTextAreaWidth := msg.Width - borderPadding
			if newTextAreaWidth > 0 {
				m.input.TextArea.SetWidth(newTextAreaWidth)
			}

			// Re-render content only if we have messages
			if len(m.messages) > 0 {
				m.updateViewportContent()
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up) || key.Matches(msg, m.keys.Down) || key.Matches(msg, m.keys.PageUp) || key.Matches(msg, m.keys.PageDown):
			// Handle scrolling first, regardless of focus
			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)
			cmds = append(cmds, vpCmd)

		case key.Matches(msg, m.keys.Tab):
			if m.focused == focusInput {
				m.focused = focusViewport
				m.input.Blur()
			} else {
				m.focused = focusInput
				cmds = append(cmds, m.input.Focus())
			}

		case key.Matches(msg, m.keys.Enter) && m.focused == focusInput:
			if !m.loading && m.input.Value() != "" {
				userMessage := m.input.Value()

				// Add user message immediately
				m.messages = append(m.messages, Message{
					Content: userMessage,
					IsUser:  true,
				})

				// Clear input
				m.input.Reset()

				// Start loading and get AI response
				m.loading = true

				// Update viewport with new message and thinking indicator
				m.updateViewportContent()

				// Switch focus to viewport so user can see the conversation
				m.focused = focusViewport
				m.input.Blur()

				cmds = append(cmds, m.spinner.Tick, getAIResponse(userMessage))
			}

		default:
			// Handle input updates when focused and not a scroll key
			if m.focused == focusInput {
				var updatedModel tea.Model
				updatedModel, cmd = m.input.Update(msg)
				m.input = updatedModel.(editor.InputModel)
				cmds = append(cmds, cmd)
			}

			// Update viewport for other keys when focused
			if m.focused == focusViewport {
				var vpCmd tea.Cmd
				m.viewport, vpCmd = m.viewport.Update(msg)
				cmds = append(cmds, vpCmd)
			}
		}

	case tea.MouseMsg:
		// Handle mouse wheel scrolling
		if msg.Type == tea.MouseWheelUp || msg.Type == tea.MouseWheelDown {
			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}

	case responseMsg:
		if msg.err != nil {
			m.loading = false
			m.err = msg.err
		} else {
			// Keep loading state - don't show message until markdown is ready
			// Add AI response but don't display it yet
			messageIndex := len(m.messages)
			m.messages = append(m.messages, Message{
				Content:    msg.content,
				IsUser:     false,
				Rendered:   "", // Will be filled by async rendering
				IsRendered: false,
			})

			// Start background markdown rendering
			cmds = append(cmds, renderMarkdownAsync(msg.content, m.width-4, messageIndex))
		}

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
			// Only update the spinner content, not the entire viewport
			m.updateSpinnerContent()
		}

	case markdownRenderedMsg:
		// Update the message with rendered markdown and stop loading
		if msg.messageIndex < len(m.messages) {
			m.messages[msg.messageIndex].Rendered = msg.rendered
			m.messages[msg.messageIndex].IsRendered = true

			// Stop loading state - markdown is ready
			m.loading = false

			// Update viewport to show rendered content
			m.updateViewportContent()

			// Now start delayed focus return to input
			cmds = append(cmds, delayedFocus())
		}

	case delayedFocusMsg:
		// Return focus to input after delay
		m.focused = focusInput
		cmds = append(cmds, m.input.Focus())
	}

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) updateViewportContent() {
	var content string

	for _, msg := range m.messages {
		if msg.IsUser {
			content += userMessageStyle.Render("You: "+msg.Content) + "\n\n"
		} else {
			// Only show AI messages that are fully rendered
			if msg.IsRendered {
				content += m.renderAIMessage(msg)
			}
		}
	}

	if m.loading {
		content += aiMessageStyle.Render("AI: ") + m.spinner.View() + " Thinking...\n"
	}

	m.viewport.SetContent(content)
	// Always scroll to bottom to show latest content
	m.viewport.GotoBottom()
}

func (m *ChatModel) renderAIMessage(msg Message) string {
	// Use cached rendered content if available, otherwise use plain text
	if msg.IsRendered && msg.Rendered != "" {
		return aiMessageStyle.Render("AI: ") + "\n" + msg.Rendered + "\n\n"
	}
	// Fallback to plain text (no expensive processing)
	return aiMessageStyle.Render("AI: "+msg.Content) + "\n\n"
}

func (m *ChatModel) updateSpinnerContent() {
	// Simply call the main update function - this is cleaner
	// The key fix is that markdown rendering only happens for complete AI messages
	m.updateViewportContent()
}

// Helper function to find max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
