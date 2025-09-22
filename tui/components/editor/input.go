package components

import (
	"strings"

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
	minHeight int
	maxHeight int
	InputKeys EditorKeyMap
}

func InitialInputModel() InputModel {
	ta := textarea.New()
	ta.Placeholder = "Type your message…"
	ta.Focus()
	ta.SetWidth(80) // Default width, will be updated on resize
	ta.SetHeight(2) // Multi-line height
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(true) // Allow Enter for new lines
	ta.KeyMap.InsertNewline.SetKeys("shift+enter", "ctrl+j")

	// The textarea component handles scrolling automatically
	// when content exceeds the visible area

	return InputModel{
		TextArea:  ta,
		width:     80,
		height:    24,
		minHeight: 2,
		maxHeight: 8,
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
	oldValue := m.TextArea.Value()
	m.TextArea, cmd = m.TextArea.Update(msg)

	if m.TextArea.Value() != oldValue {
		newHeight := m.CalculateHeight()
		m.TextArea.SetHeight(newHeight)
	}

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

func (m *InputModel) MinHeight() int {
	return m.minHeight
}

func (m *InputModel) CalculateHeight() int {
	content := m.TextArea.Value()
	if content == "" {
		return m.minHeight
	}

	lines := strings.Split(content, "\n")

	width := m.TextArea.Width()
	if width <= 0 {
		width = 80
	}

	totalLines := 0
	for _, line := range lines {
		if line == "" {
			totalLines++
		} else {
			wrappedLines := (len(line)-1)/width + 1
			totalLines += wrappedLines
		}
	}

	if totalLines < m.minHeight {
		return m.minHeight
	}
	if totalLines > m.maxHeight {
		return m.maxHeight
	}
	return totalLines
}
