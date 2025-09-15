package components

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	textInput textinput.Model
	err       error
	prompt    string
	submitted bool
	width     int
	height    int
}

func initialInputModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter your prompt to ask AI..."
	ti.Focus()
	ti.Width = 80 // Default width, will be updated on resize

	return inputModel{
		textInput: ti,
		width:     80,
		height:    24,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = m.width - 4
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.prompt = m.textInput.Value()
			m.submitted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m inputModel) View() string {
	return "Enter the prompt to ask ai:\n" + m.textInput.View()
}

func RunInputModel() string {
	p := tea.NewProgram(
		initialInputModel(),
		tea.WithAltScreen(),
	)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	inputModel, ok := finalModel.(inputModel)
	if !ok || !inputModel.submitted {
		return ""
	}

	return inputModel.prompt
}
