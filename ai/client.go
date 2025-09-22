package ai

import (
	"github.com/krishkalaria12/nyron-ai-cli/ai/provider"
)

// StreamMessage represents a chunk of response or an error
type StreamMessage = provider.StreamMessage

// GeminiAPI generates a complete response using Gemini API
func GeminiAPI(prompt string) (string, error) {
	return provider.GeminiAPI(prompt)
}

// OpenAIAPI generates a complete response using OpenAI API
func OpenAIAPI(prompt string) (string, error) {
	return provider.OpenAIAPI(prompt)
}

// OpenRouterAPI generates a complete response using OpenRouter API
func OpenRouterAPI(prompt string, model string) (string, error) {
	return provider.OpenRouterAPI(prompt, model)
}
