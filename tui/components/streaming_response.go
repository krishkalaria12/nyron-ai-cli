package components

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
	"github.com/muesli/reflow/wordwrap"
)

type streamingResponseModel struct {
	content      string
	loading      bool
	ready        bool
	viewport     viewport.Model
	spinner      spinner.Model
	width        int
	height       int
	prompt       string
	responseChan chan ai.StreamMessage
	err          error
}

// Messages
type streamChunkMsg ai.StreamMessage
type startStreamMsg string

var (
	loadingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C3AED")).
		Bold(true)
)

func NewStreamingResponseModel(prompt string) streamingResponseModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	return streamingResponseModel{
		prompt:       prompt,
		loading:      true,
		spinner:      s,
		responseChan: make(chan ai.StreamMessage),
	}
}

func (m *streamingResponseModel) wrapContent() string {
	if m.width <= 0 {
		return m.content
	}
	// Leave some padding for the viewport
	width := m.width - 4
	if width < 10 {
		width = 10
	}
	return wordwrap.String(m.content, width)
}

func (m streamingResponseModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		startStreamCommand(m.prompt, m.responseChan),
		waitForStreamMessage(m.responseChan),
	)
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

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.wrapContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			// Re-wrap content when viewport size changes
			m.viewport.SetContent(m.wrapContent())
		}

	case streamChunkMsg:
		if msg.Error != nil {
			m.err = msg.Error
			m.loading = false
			return m, nil
		}

		if msg.Done {
			m.loading = false
			return m, nil
		}

		// Append new content
		m.content += msg.Content
		if m.ready {
			m.viewport.SetContent(m.wrapContent())
			// Auto-scroll to bottom
			m.viewport.GotoBottom()
		}

		// Wait for next message
		return m, waitForStreamMessage(m.responseChan)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup", "b":
			m.viewport.HalfViewUp()
		case "pgdown", "f":
			m.viewport.HalfViewDown()
		case "home", "g":
			m.viewport.GotoTop()
		case "end", "G":
			m.viewport.GotoBottom()
		}

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Handle keyboard and mouse events in the viewport
	if m.ready {
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

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m streamingResponseModel) headerView() string {
	title := "AI Response"
	if m.loading {
		title = fmt.Sprintf("%s %s Generating...", m.spinner.View(), title)
	}

	titleStyled := titleStyle.Render(title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(titleStyled)))
	return lipgloss.JoinHorizontal(lipgloss.Center, titleStyled, line)
}

func (m streamingResponseModel) footerView() string {
	info := ""
	if m.loading {
		info = "⏳ Streaming..."
	} else {
		scrollInfo := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
		helpText := " • ↑/↓ j/k pgup/pgdn g/G scroll • q quit"
		info = infoStyle.Render(scrollInfo + helpText)
	}
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// Commands
func startStreamCommand(prompt string, responseChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		// Start streaming in a goroutine
		go ai.GeminiStreamAPI(prompt, responseChan)
		return startStreamMsg(prompt)
	}
}

func waitForStreamMessage(responseChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		msg := <-responseChan
		return streamChunkMsg(msg)
	}
}

func RunStreamingResponseModel(prompt string) {
	p := tea.NewProgram(
		NewStreamingResponseModel(prompt),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
