package chat

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
	"github.com/krishkalaria12/nyron-ai-cli/ai/tools"
	"github.com/krishkalaria12/nyron-ai-cli/config"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
	"github.com/krishkalaria12/nyron-ai-cli/tui/components/dialogs/models"
	editor "github.com/krishkalaria12/nyron-ai-cli/tui/components/editor"
	"github.com/krishkalaria12/nyron-ai-cli/util"
	openrouter "github.com/revrost/go-openrouter"
)

type focusState int

const (
	focusViewport focusState = iota
	focusInput
)

type keyMap struct {
	Tab        key.Binding
	Up         key.Binding
	Down       key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	Enter      key.Binding
	Quit       key.Binding
	OpenDialog key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.Tab, k.Enter, k.OpenDialog, k.Quit},
	}
}

var keys = keyMap{
	Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
	Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
	PageUp:     key.NewBinding(key.WithKeys("pgup", "ctrl+u"), key.WithHelp("pgup", "page up")),
	PageDown:   key.NewBinding(key.WithKeys("pgdown", "ctrl+d"), key.WithHelp("pgdn", "page down")),
	Tab:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch focus")),
	Enter:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send message")),
	Quit:       key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	OpenDialog: key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "choose model")),
}

// Message represents a chat message for UI rendering
type Message struct {
	Content    string
	IsUser     bool
	Rendered   string     // Cached markdown rendering
	IsRendered bool       // Whether markdown processing is complete
	Thinking   string     // AI thinking process (if available)
	ToolCalls  []ToolCall // Tool calls made during this message
}

// ToolCall represents a single tool call for UI rendering
type ToolCall struct {
	Step    string // The tool call step description
	Content string // The content/result of the tool call
}

type ChatModel struct {
	messages            []Message                          // For UI rendering
	conversationHistory []openrouter.ChatCompletionMessage // For API calls
	loading             bool
	viewport            viewport.Model
	spinner             spinner.Model
	input               editor.InputModel
	keys                keyMap
	focused             focusState
	help                help.Model
	width               int
	height              int
	err                 error
	selectedModel       config.SelectedModel
	showDialog          bool
	modelDialog         *models.ModelListComponent
}

// --- New Message Types for the event loop ---
type responseMsg struct {
	response openrouter.ChatCompletionResponse
	err      error
}

type toolResultsMsg struct {
	results []openrouter.ChatCompletionMessage
}

func NewChatModel() ChatModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	inputModel := editor.InitialInputModel()
	inputModel.InputKeys = editor.DefaultEditorKeyMap()

	vp := viewport.New(80, 20)

	return ChatModel{
		messages:            []Message{},
		conversationHistory: []openrouter.ChatCompletionMessage{},
		loading:             false,
		spinner:             s,
		input:               inputModel,
		viewport:            vp,
		focused:             focusInput,
		keys:                keys,
		help:                help.New(),
		selectedModel: config.SelectedModel{
			Provider: "openrouter",
			Model:    "google/gemini-2.5-flash",
		},
		showDialog: false,
		modelDialog: func() *models.ModelListComponent {
			component := models.NewModelListComponent()
			return &component
		}(),
	}
}

func (m ChatModel) Init() tea.Cmd {
	return m.input.Focus()
}

func getAIResponse(history []openrouter.ChatCompletionMessage, selectedModel config.SelectedModel) tea.Cmd {
	return func() tea.Msg {
		modelID := selectedModel.Model
		if selectedModel.Provider != "openrouter" {
			modelID = "google/gemini-2.5-flash"
		}
		resp, err := ai.OpenRouterAPI(history, modelID)
		return responseMsg{response: resp, err: err}
	}
}

// executeToolsCmd processes the tool calls requested by the AI.
func executeToolsCmd(calls []openrouter.ToolCall) tea.Cmd {
	return func() tea.Msg {
		var results []openrouter.ChatCompletionMessage
		for _, call := range calls {
			toolResult := tools.ExecuteTool(call.Function.Name, call.Function.Arguments)
			results = append(results, openrouter.ChatCompletionMessage{
				Role:       openrouter.ChatMessageRoleTool,
				Content:    openrouter.Content{Text: toolResult},
				ToolCallID: call.ID,
			})
		}
		return toolResultsMsg{results: results}
	}
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.viewport.Width = m.width
		m.help.Width = m.width

		// Calculate input width accounting for border and padding
		inputFrameSize := focusedInputBorderStyle.GetHorizontalFrameSize()
		m.input.TextArea.SetWidth(m.width - inputFrameSize)

		// Update viewport height based on current input height
		m.updateViewportHeight()

		if len(m.messages) > 0 {
			m.updateViewportContent()
		}

	case tea.KeyMsg:
		if m.showDialog {
			var cmd tea.Cmd
			updatedModel, cmd := m.modelDialog.Update(msg)
			if listComponent, ok := updatedModel.(models.ModelListComponent); ok {
				*m.modelDialog = listComponent
			}
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.OpenDialog):
			m.showDialog = true
			return m, m.modelDialog.Init()
		case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.PageUp), key.Matches(msg, m.keys.PageDown):
			switch m.focused {
			case focusViewport:
				var vpCmd tea.Cmd
				m.viewport, vpCmd = m.viewport.Update(msg)
				cmds = append(cmds, vpCmd)
			case focusInput:
				// Pass navigation keys to the input for cursor movement and scrolling
				oldInputHeight := m.input.TextArea.Height()
				var updatedModel tea.Model
				updatedModel, cmd = m.input.Update(msg)
				m.input = updatedModel.(editor.InputModel)

				// Update viewport height if input height changed
				if m.input.TextArea.Height() != oldInputHeight {
					m.updateViewportHeight()
				}

				cmds = append(cmds, cmd)
			}
		case key.Matches(msg, m.keys.Tab):
			if m.focused == focusInput {
				m.focused = focusViewport
				m.input.Blur()
			} else {
				m.focused = focusInput
				cmds = append(cmds, m.input.Focus())
			}
		case key.Matches(msg, m.keys.Enter) && m.focused == focusInput:
			// Only send message on plain Enter, not Shift+Enter
			if msg.String() == "enter" && !m.loading && m.input.Value() != "" {
				userMessageContent := m.input.Value()
				m.messages = append(m.messages, Message{Content: userMessageContent, IsUser: true})
				m.input.Reset()

				// If this is a new conversation, add the system prompt first.
				if len(m.conversationHistory) == 0 {
					promptPair := prompts.GetPrompts(userMessageContent, "openrouter")
					m.conversationHistory = append(m.conversationHistory, openrouter.ChatCompletionMessage{
						Role:    openrouter.ChatMessageRoleSystem,
						Content: openrouter.Content{Text: promptPair.SystemPrompt},
					})
				}

				// Append the user message to the API history
				m.conversationHistory = append(m.conversationHistory, openrouter.ChatCompletionMessage{
					Role:    openrouter.ChatMessageRoleUser,
					Content: openrouter.Content{Text: userMessageContent},
				})

				// Reset input height to minimum after clearing
				m.input.TextArea.SetHeight(m.input.MinHeight())
				m.updateViewportHeight()

				m.loading = true
				m.updateViewportContentWithScroll(true)
				m.focused = focusViewport
				m.input.Blur()
				cmds = append(cmds, m.spinner.Tick, getAIResponse(m.conversationHistory, m.selectedModel))
			}
		default:
			switch m.focused {
			case focusInput:
				oldInputHeight := m.input.TextArea.Height()
				var updatedModel tea.Model
				updatedModel, cmd = m.input.Update(msg)
				m.input = updatedModel.(editor.InputModel)

				// Update viewport height if input height changed
				if m.input.TextArea.Height() != oldInputHeight {
					m.updateViewportHeight()
				}

				cmds = append(cmds, cmd)
			case focusViewport:
				var vpCmd tea.Cmd
				m.viewport, vpCmd = m.viewport.Update(msg)
				cmds = append(cmds, vpCmd)
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && (msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown) {
			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}

	// Handling the conversation cycle

	case responseMsg:
		if msg.err != nil {
			m.loading = false
			m.err = msg.err
			// Clear history on error to start fresh
			m.conversationHistory = []openrouter.ChatCompletionMessage{}
			return m, nil
		}

		assistantMessage := msg.response.Choices[0].Message
		m.conversationHistory = append(m.conversationHistory, assistantMessage)

		if len(assistantMessage.ToolCalls) > 0 {
			// AI wants to use tools
			var uiToolCalls []ToolCall
			for _, call := range assistantMessage.ToolCalls {
				uiToolCalls = append(uiToolCalls, ToolCall{
					Step:    call.Function.Name,
					Content: call.Function.Arguments,
				})
			}
			// Add a new message to the UI to show what the AI is doing
			m.messages = append(m.messages, Message{
				IsUser:     false,
				IsRendered: true, // Mark as rendered to show tool call info
				ToolCalls:  uiToolCalls,
			})
			m.updateViewportContentWithScroll(true)
			// Dispatch a command to execute the tools
			cmds = append(cmds, executeToolsCmd(assistantMessage.ToolCalls))
		} else {
			// This is the final text response
			messageIndex := len(m.messages)
			finalContent := assistantMessage.Content.Text
			thinking := ""
			if assistantMessage.Reasoning != nil {
				thinking = *assistantMessage.Reasoning
			}
			m.messages = append(m.messages, Message{
				Content:    finalContent,
				IsUser:     false,
				IsRendered: false,
				Thinking:   thinking,
			})
			// Render the final markdown response
			cmds = append(cmds, util.RenderMarkdownAsync(finalContent, m.width-4, messageIndex))
			// Conversation is over, clear history for the next prompt
			m.conversationHistory = []openrouter.ChatCompletionMessage{}
		}

	case toolResultsMsg:
		// Append tool results to history
		m.conversationHistory = append(m.conversationHistory, msg.results...)

		// add a UI message here to show the tool's raw output.
		// For now, we immediately call the AI again with the new context.
		cmds = append(cmds, getAIResponse(m.conversationHistory, m.selectedModel))

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
			m.updateSpinnerContent()
		}

	case util.MarkdownRenderedMsg:
		// This handles the final response rendering
		if msg.MessageIndex < len(m.messages) {
			m.messages[msg.MessageIndex].Rendered = msg.Rendered
			m.messages[msg.MessageIndex].IsRendered = true
			m.loading = false // Stop loading only after final render
			m.updateViewportContentWithScroll(true)
			cmds = append(cmds, util.DelayedFocus())
		}

	case util.DelayedFocusMsg:
		m.focused = focusInput
		cmds = append(cmds, m.input.Focus())

	case models.ModelSelectedMsg:
		m.selectedModel = msg.Model
		m.showDialog = false
		cmds = append(cmds, m.input.Focus())

	case models.CloseModelDialog:
		m.showDialog = false
		cmds = append(cmds, m.input.Focus())
	}

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) updateViewportContent() {
	m.updateViewportContentWithScroll(false)
}

func (m *ChatModel) updateViewportContentWithScroll(autoScroll bool) {
	var content string
	for _, msg := range m.messages {
		if msg.IsUser {
			userLabel := userMessageStyle.Render("You:")
			userContent := userMessageContentStyle.Width(m.width - userMessageContentStyle.GetHorizontalFrameSize()).Render(msg.Content)
			content += userLabel + " " + userContent + "\n\n"
		} else if msg.IsRendered {
			content += m.renderAIMessage(msg)
		}
	}
	if m.loading {
		aiLabel := aiMessageStyle.Render("AI:")
		thinkingText := thinkingStyle.Render(m.spinner.View() + " Thinking...")
		content += aiLabel + " " + thinkingText
	}
	m.viewport.SetContent(content)
	if autoScroll {
		m.viewport.GotoBottom()
	}
}

func (m *ChatModel) renderAIMessage(msg Message) string {
	aiLabel := aiMessageStyle.Render("AI:")

	var content string

	// Add thinking section if available
	if msg.Thinking != "" {
		thinkingHeader := thinkingHeaderStyle.Render("🤔 Thinking:")
		thinkingContent := thinkingStyle.Width(m.width - thinkingStyle.GetHorizontalFrameSize()).Render(msg.Thinking)
		content = lipgloss.JoinVertical(lipgloss.Left, thinkingHeader, thinkingContent)
		content += "\n\n"
	}

	// Add tool calls if available
	if len(msg.ToolCalls) > 0 {
		toolHeader := toolCallHeaderStyle.Render("🔧 Tool Calls:")
		content += toolHeader + "\n"

		for _, toolCall := range msg.ToolCalls {
			stepText := toolCallStyle.Render("Step: " + toolCall.Step)
			contentText := toolCallContentStyle.Width(m.width - toolCallContentStyle.GetHorizontalFrameSize()).Render(toolCall.Content)

			content += stepText + "\n"
			content += contentText + "\n\n"
		}
	}

	// Add main response content if available
	if msg.Rendered != "" {
		aiContent := aiMessageContentStyle.Width(m.width - aiMessageContentStyle.GetHorizontalFrameSize()).Render(msg.Rendered)
		content += lipgloss.JoinVertical(lipgloss.Left, aiLabel, aiContent)
	} else if len(msg.ToolCalls) > 0 {
		// If we only have tool calls, still show the AI label
		content = lipgloss.JoinVertical(lipgloss.Left, aiLabel, content)
	}

	return content + "\n\n"
}

func (m *ChatModel) updateSpinnerContent() {
	m.updateViewportContent()
}

func (m *ChatModel) updateViewportHeight() {
	headerView := headerStyle.Render("Header")
	helpView := helpStyle.Width(m.width).Render(m.help.View(m.keys))
	inputView := focusedInputBorderStyle.Width(m.width).Render(m.input.View())

	verticalMargin := lipgloss.Height(headerView) + lipgloss.Height(inputView) + lipgloss.Height(helpView)
	m.viewport.Height = m.height - verticalMargin
}
