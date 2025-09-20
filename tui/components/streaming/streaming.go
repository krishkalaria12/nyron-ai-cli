package streaming

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/krishkalaria12/nyron-ai-cli/ai"
)

func startStreamCommand(prompt string, responseChan chan ai.StreamMessage, forwardChan chan ai.StreamMessage) tea.Cmd {
	return func() tea.Msg {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					forwardChan <- ai.StreamMessage{
						Content: "",
						Error:   fmt.Errorf("AI call panicked: %v", r),
						Done:    true,
					}
					close(forwardChan)
				}
			}()
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

			defer func() {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
			}()

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
