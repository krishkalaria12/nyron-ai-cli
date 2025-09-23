package ai

import (
	"github.com/krishkalaria12/nyron-ai-cli/ai/provider"
)

// StreamMessage represents a chunk of response or an error
type StreamMessage = provider.StreamMessage

// OpenRouterAPI generates a complete response using OpenRouter API
func OpenRouterAPI(prompt string, model string) provider.AIResponseMessage {
	return provider.OpenRouterAPI(prompt, model)
}
