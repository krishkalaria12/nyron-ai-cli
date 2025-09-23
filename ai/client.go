package ai

import (
	"github.com/krishkalaria12/nyron-ai-cli/ai/provider"
	"github.com/revrost/go-openrouter"
)

func OpenRouterAPI(messages []openrouter.ChatCompletionMessage, model string) (openrouter.ChatCompletionResponse, error) {
	return provider.OpenRouterAPI(messages, model)
}
