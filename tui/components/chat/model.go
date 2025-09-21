package chat

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
	"github.com/krishkalaria12/nyron-ai-cli/config"
	"github.com/krishkalaria12/nyron-ai-cli/tui/components/dialogs/models"
	editor "github.com/krishkalaria12/nyron-ai-cli/tui/components/editor"
	"github.com/krishkalaria12/nyron-ai-cli/util"
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

// Message represents a chat message
type Message struct {
	Content    string
	IsUser     bool
	Rendered   string // Cached markdown rendering
	IsRendered bool   // Whether markdown processing is complete
}

type ChatModel struct {
	messages      []Message
	loading       bool
	viewport      viewport.Model
	spinner       spinner.Model
	input         editor.InputModel
	keys          keyMap
	focused       focusState
	help          help.Model
	width         int
	height        int
	err           error
	selectedModel config.SelectedModel
	showDialog    bool
	modelDialog   *models.ModelListComponent
}

// Messages
type responseMsg struct {
	content string
	err     error
}

func NewChatModel() ChatModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	inputModel := editor.InitialInputModel()
	inputModel.InputKeys = editor.DefaultEditorKeyMap()

	vp := viewport.New(80, 20)

	return ChatModel{
		messages: []Message{},
		loading:  false,
		spinner:  s,
		input:    inputModel,
		viewport: vp,
		focused:  focusInput,
		keys:     keys,
		help:     help.New(),
		selectedModel: config.SelectedModel{
			Provider: "gemini",
			Model:    "gemini-2.5-flash",
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

func getAIResponse(prompt string, selectedModel config.SelectedModel) tea.Cmd {
	return func() tea.Msg {
		var response string
		var err error
		switch selectedModel.Provider {
		case "gemini":
			response, err = ai.GeminiAPI(prompt)
		case "openai":
			response, err = ai.OpenAIAPI(prompt)
		case "openrouter":
			response, err = ai.OpenRouterAPI(prompt, selectedModel.Model)
		default:
			response, err = ai.GeminiAPI(prompt)
		}
		return responseMsg{content: response, err: err}
	}
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		headerView := headerStyle.Render("Header")
		helpView := helpStyle.Width(m.width).Render(m.help.View(m.keys))
		inputView := focusedInputBorderStyle.Width(m.width).Render(m.input.View())

		verticalMargin := lipgloss.Height(headerView) + lipgloss.Height(inputView) + lipgloss.Height(helpView)

		m.viewport.Width = m.width
		m.viewport.Height = m.height - verticalMargin
		m.help.Width = m.width

		// Calculate input width accounting for border and padding
		inputFrameSize := focusedInputBorderStyle.GetHorizontalFrameSize()
		m.input.TextArea.SetWidth(m.width - inputFrameSize)

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
			// Only send message on plain Enter, not Shift+Enter
			if msg.String() == "enter" && !m.loading && m.input.Value() != "" {
				userMessage := m.input.Value()
				m.messages = append(m.messages, Message{Content: userMessage, IsUser: true})
				m.input.Reset()
				m.loading = true
				m.updateViewportContent()
				m.focused = focusViewport
				m.input.Blur()
				cmds = append(cmds, m.spinner.Tick, getAIResponse(userMessage, m.selectedModel))
			}
		default:
			if m.focused == focusInput {
				var updatedModel tea.Model
				updatedModel, cmd = m.input.Update(msg)
				m.input = updatedModel.(editor.InputModel)
				cmds = append(cmds, cmd)
			} else if m.focused == focusViewport {
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

	case responseMsg:
		if msg.err != nil {
			m.loading = false
			m.err = msg.err
		} else {
			messageIndex := len(m.messages)
			m.messages = append(m.messages, Message{Content: msg.content, IsUser: false, IsRendered: false})
			cmds = append(cmds, util.RenderMarkdownAsync(msg.content, m.width-4, messageIndex))
		}

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
			m.updateSpinnerContent()
		}

	case util.MarkdownRenderedMsg:
		if msg.MessageIndex < len(m.messages) {
			m.messages[msg.MessageIndex].Rendered = msg.Rendered
			m.messages[msg.MessageIndex].IsRendered = true
			m.loading = false
			m.updateViewportContent()
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
	var content string
	for _, msg := range m.messages {
		if msg.IsUser {
			userLabel := userMessageStyle.Render("You:")
			userContent := userMessageContentStyle.Width(m.width - userMessageContentStyle.GetHorizontalFrameSize() - 2).Render(msg.Content)
			content += lipgloss.JoinVertical(lipgloss.Left, userLabel, userContent) + "\n\n"
		} else if msg.IsRendered {
			content += m.renderAIMessage(msg)
		}
	}
	if m.loading {
		aiLabel := aiMessageStyle.Render("AI:")
		thinkingText := thinkingStyle.Render(m.spinner.View() + " Thinking...")
		content += lipgloss.JoinVertical(lipgloss.Left, aiLabel, thinkingText)
	}
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

func (m *ChatModel) renderAIMessage(msg Message) string {
	aiLabel := aiMessageStyle.Render("AI:")
	aiContent := aiMessageContentStyle.Width(m.width - aiMessageContentStyle.GetHorizontalFrameSize() - 2).Render(msg.Rendered)
	return lipgloss.JoinVertical(lipgloss.Left, aiLabel, aiContent) + "\n\n"
}

func (m *ChatModel) updateSpinnerContent() {
	m.updateViewportContent()
}
