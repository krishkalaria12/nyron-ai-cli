package components

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
)

type streamingResponseModel struct {
	content      string
	rawContent   string // Store raw markdown for re-rendering on resize
	loading      bool
	ready        bool
	viewport     viewport.Model
	spinner      spinner.Model
	width        int
	height       int
	prompt       string
	responseChan chan ai.StreamMessage // producer -> forwarder
	forwardChan  chan ai.StreamMessage // forwarder -> tea loop
	err          error
}

// Messages
type streamChunkMsg ai.StreamMessage
type startStreamMsg string

var (
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
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
		forwardChan:  make(chan ai.StreamMessage),
	}
}

func startStreamCommand(prompt string, responseChan chan ai.StreamMessage, forwardChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		go func() {
			ai.GeminiStreamAPI(prompt, responseChan)
		}()

		go func() {
			const (
				maxBatchDelay = 50 * time.Millisecond // debounce time
				maxBatchSize  = 1024                  // flush if buffer reaches this size
			)

			var b strings.Builder
			timer := time.NewTimer(time.Hour)
			timer.Stop()
			pending := false

			flush := func(done bool) {
				if !pending && !done {
					return
				}
				content := b.String()

				forwardChan <- ai.StreamMessage{
					Content: content,
					Error:   nil,
					Done:    done,
				}

				b.Reset()
				pending = false
			}

			for {
				select {
				case ev, ok := <-responseChan:
					if !ok {
						flush(true)
						close(forwardChan)
						return
					}
					if ev.Error != nil {
						forwardChan <- ai.StreamMessage{
							Content: "",
							Error:   ev.Error,
							Done:    true,
						}
						close(forwardChan)
						return
					}

					if ev.Content != "" {
						b.WriteString(ev.Content)
						pending = true
					}

					if b.Len() >= maxBatchSize {
						flush(false)
						if !timer.Stop() {
							select {
							case <-timer.C:
							default:
							}
						}
					} else if pending {
						if !timer.Stop() {
							select {
							case <-timer.C:
							default:
							}
						}
						timer.Reset(maxBatchDelay)
					}

					if ev.Done {
						flush(true)
						close(forwardChan)
						return
					}

				case <-timer.C:
					flush(false)
				}
			}
		}()

		return startStreamMsg(prompt)
	}
}

func waitForForwardMessage(forwardChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-forwardChan
		if !ok {
			return streamChunkMsg(ai.StreamMessage{
				Content: "",
				Error:   nil,
				Done:    true,
			})
		}
		return streamChunkMsg(msg)
	}
}

// renderContentWithCurrentWidth renders the accumulated rawContent using the ai/glamour renderer
func (m *streamingResponseModel) renderContentWithCurrentWidth() string {
	// If no raw markdown accumulated yet, return whatever content we have
	if m.rawContent == "" {
		return m.content
	}

	// Calculate appropriate width for markdown rendering
	renderWidth := m.width - 4 // Leave padding for viewport borders
	if renderWidth < 20 {
		renderWidth = 20 // Minimum reasonable width
	}

	rendered, err := ai.RenderToTerminalWithWidth(m.rawContent, renderWidth)
	if err != nil {
		// Fallback: if glamour rendering fails, show raw markdown (so user still sees text)
		return m.rawContent
	}

	// update cached content with rendered ANSI
	m.content = rendered
	return m.content
}

func (m streamingResponseModel) Init() tea.Cmd {
	// Start spinner tick + start the streaming producer + forwarder
	return tea.Batch(
		m.spinner.Tick,
		startStreamCommand(m.prompt, m.responseChan, m.forwardChan),
		waitForForwardMessage(m.forwardChan), // wait for the first forwarded message
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
			m.viewport.SetContent(m.renderContentWithCurrentWidth())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			// Re-render content with new width when viewport size changes
			m.viewport.SetContent(m.renderContentWithCurrentWidth())
		}

	case streamChunkMsg:
		// Convert incoming message to ai.StreamMessage semantics
		in := ai.StreamMessage(msg)

		// Error handling from producer
		if in.Error != nil {
			m.err = in.Error
			m.loading = false
			// ensure viewport displays error
			if m.ready {
				m.viewport.SetContent(fmt.Sprintf("Error: %v\n", in.Error))
			}
			return m, nil
		}

		// If Done, mark loading false and final-render accumulated content
		if in.Done {
			m.loading = false
			// final render of accumulated content (in case any sanitization appended fences)
			if m.rawContent != "" {
				final := m.renderContentWithCurrentWidth()
				if m.ready {
					m.viewport.SetContent(final)
					m.viewport.GotoBottom()
				}
			}
			return m, nil
		}

		// Append the incoming chunk to the raw content (accumulate)
		m.rawContent += in.Content

		// Render the accumulated markdown (prevents broken fences & flicker).
		rendered := m.renderContentWithCurrentWidth()

		// Update viewport
		if m.ready {
			m.viewport.SetContent(rendered)
			// Auto-scroll to bottom for live updates
			m.viewport.GotoBottom()
		}

		// Schedule next wait for message
		return m, waitForForwardMessage(m.forwardChan)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.viewport.ScrollUp(1)
		case "down", "j":
			m.viewport.ScrollDown(1)
		case "pgup", "b":
			m.viewport.HalfPageUp()
		case "pgdown", "f":
			m.viewport.HalfPageDown()
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

// Run helper remains the same
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

// small helper
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
