package models

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/krishkalaria12/nyron-ai-cli/config"
	"github.com/krishkalaria12/nyron-ai-cli/tui/components/dialogs"
)

const (
	defaultWidth = 100
)

type ModelSelectedMsg struct {
	Model config.SelectedModel
}

// sent when model is chosen
type CloseModelDialog struct{}

// ModelDialog interface for the model selection dialog
type ModelDialog interface {
	dialogs.DialogModel
}

type modelDialogCmp struct {
	width   int
	wWidth  int
	wHeight int

	modelList *ModelListComponent
	keyMap    KeyMap
	help      help.Model
}

func NewModelDialogCmp() ModelDialog {
	help := help.New()
	modelList := NewModelListComponent()

	return &modelDialogCmp{
		modelList: &modelList,
		width:     defaultWidth,
		keyMap:    DefaultKeyMap(),
		help:      help,
	}
}

func (m *modelDialogCmp) Init() tea.Cmd {
	return m.modelList.Init()
}

func (m *modelDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.wWidth = msg.Width
		m.wHeight = msg.Height
	}

	var updatedModel tea.Model
	updatedModel, cmd = m.modelList.Update(msg)
	if listComponent, ok := updatedModel.(ModelListComponent); ok {
		*m.modelList = listComponent
	}

	return m, cmd
}

func (m *modelDialogCmp) View() string {
	return m.modelList.View()
}
