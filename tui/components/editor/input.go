package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	TextArea  textarea.Model
	err       error
	prompt    string
	submitted bool
	width     int
	height    int
	InputKeys EditorKeyMap
}

func InitialInputModel() InputModel {
	ta := textarea.New()
	ta.Placeholder = "Type your message…  Ctrl+Enter: send  Enter: new line"
	ta.Focus()
	ta.SetWidth(80) // Default width, will be updated on resize
	ta.SetHeight(3) // Multi-line height
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(true) // Allow Enter for new lines

	return InputModel{
		TextArea: ta,
		width:    80,
		height:   24,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.TextArea.SetWidth(m.width - 4)
	}

	var cmd tea.Cmd
	m.TextArea, cmd = m.TextArea.Update(msg)

	return m, cmd
}

func (m InputModel) View() string {
	return m.TextArea.View()
}

// Helper methods for the textarea
func (m *InputModel) Focus() tea.Cmd {
	return m.TextArea.Focus()
}

func (m *InputModel) Blur() {
	m.TextArea.Blur()
}

func (m *InputModel) Value() string {
	return m.TextArea.Value()
}

func (m *InputModel) Reset() {
	m.TextArea.Reset()
}
