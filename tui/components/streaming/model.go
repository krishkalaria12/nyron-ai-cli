package streaming

import (
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

type StreamingResponseModel struct {
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

func NewStreamingResponseModel() StreamingResponseModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	// Initialize input model with editor keymap
	inputModel := editor.InitialInputModel()
	inputModel.InputKeys = editor.DefaultEditorKeyMap()

	return StreamingResponseModel{
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

func (m StreamingResponseModel) Init() tea.Cmd {
	// Initialize input focus
	return m.input.Focus()
}