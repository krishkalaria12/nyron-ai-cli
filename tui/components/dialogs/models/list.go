package models

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/config"
)

var (
	primaryColor   = lipgloss.Color("#6366f1")
	secondaryColor = lipgloss.Color("#8b5cf6")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#FFFFFF"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(secondaryColor).Bold(true)
	providerHeadStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).PaddingLeft(2)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type ListItem interface {
	list.Item
	IsHeader() bool
}

type ProviderHeader struct {
	name string
}

func (h ProviderHeader) FilterValue() string { return "" }
func (h ProviderHeader) IsHeader() bool     { return true }

type ModelItem struct {
	provider config.Provider
	model    config.Model
}

func (i ModelItem) FilterValue() string { return i.model.Name }
func (i ModelItem) Title() string       { return i.model.Name }
func (i ModelItem) Description() string {
	return fmt.Sprintf("%s - %s", i.provider.Name, i.model.Description)
}
func (i ModelItem) IsHeader() bool { return false }

type ModelListComponent struct {
	list     list.Model
	choice   *ModelItem
	quitting bool
}

func NewModelListComponent() ModelListComponent {
	var items []list.Item
	providers := config.GetAllProviders()
	for _, provider := range providers {
		models := config.GetModelsByProvider(provider.ID)
		if len(models) > 0 {
			// Add provider header
			items = append(items, ProviderHeader{name: provider.Name})
			// Add models for this provider
			for _, model := range models {
				items = append(items, ModelItem{provider: provider, model: model})
			}
		}
	}

	// Calculate initial dimensions (will be updated on window resize)
	initialWidth := 80  // Make it wider initially
	initialHeight := 20
	l := list.New(items, itemDelegate{}, initialWidth, initialHeight)
	l.Title = "Choose an AI Model"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false) // Disable filtering to preserve grouping
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Styles.Title = titleStyle

	// Set initial selection to first model (skip headers)
	for i, item := range items {
		if modelItem, ok := item.(ModelItem); ok {
			l.Select(i)
			_ = modelItem // Use the variable to avoid unused warning
			break
		}
	}

	return ModelListComponent{list: l}
}

func (m ModelListComponent) Init() tea.Cmd {
	return nil
}

func (m ModelListComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Calculate dynamic width based on terminal size
		// Use 80% of terminal width, with min 60 and max 150
		dynamicWidth := int(float64(msg.Width) * 0.8)
		if dynamicWidth < 60 {
			dynamicWidth = 60
		} else if dynamicWidth > 150 {
			dynamicWidth = 150
		}

		// Calculate height based on terminal height (60% with min 15)
		dynamicHeight := int(float64(msg.Height) * 0.6)
		if dynamicHeight < 15 {
			dynamicHeight = 15
		}

		m.list.SetSize(dynamicWidth, dynamicHeight)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(ModelItem); ok {
				m.choice = &i
				return m, func() tea.Msg {
					return ModelSelectedMsg{Model: config.SelectedModel{Provider: i.provider.ID, Model: i.model.ID}}
				}
			}
			// If it's a header, do nothing
			return m, nil
		case "esc":
			return m, func() tea.Msg { return CloseModelDialog{} }
		case "up", "k":
			// Move up and skip headers
			currentIndex := m.list.Index()
			for i := currentIndex - 1; i >= 0; i-- {
				if item, ok := m.list.Items()[i].(ListItem); ok && !item.IsHeader() {
					m.list.Select(i)
					break
				}
			}
			return m, nil
		case "down", "j":
			// Move down and skip headers
			currentIndex := m.list.Index()
			items := m.list.Items()
			for i := currentIndex + 1; i < len(items); i++ {
				if item, ok := items[i].(ListItem); ok && !item.IsHeader() {
					m.list.Select(i)
					break
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ModelListComponent) View() string {
	if m.quitting {
		return quitTextStyle.Render("Goodbye!")
	}
	return m.list.View()
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	// Handle provider headers
	if header, ok := listItem.(ProviderHeader); ok {
		headerText := providerHeadStyle.Render(header.name)
		fmt.Fprint(w, headerText)
		return
	}

	// Handle model items
	if modelItem, ok := listItem.(ModelItem); ok {
		modelName := modelItem.Title()
		line := fmt.Sprintf("  • %s", modelName)

		fn := itemStyle.Render
		if index == m.Index() {
			fn = func(s ...string) string {
				return selectedItemStyle.Render("> " + s[0])
			}
		}
		fmt.Fprint(w, fn(line))
	}
}
