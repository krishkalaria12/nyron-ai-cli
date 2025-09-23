package ai

import (
	"github.com/krishkalaria12/nyron-ai-cli/ai/provider"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
)

// StreamMessage represents a chunk of response or an error
type StreamMessage = provider.StreamMessage

// OpenRouterAPI generates a complete response using OpenRouter API
func OpenRouterAPI(userPrompt string, model string, toolchan chan<- provider.ToolCallingResponse) provider.AIResponseMessage {
	promptPair := prompts.GetPrompts(userPrompt, "openrouter")
	return provider.OpenRouterAPI(promptPair.SystemPrompt, promptPair.UserPrompt, model, toolchan)
}
