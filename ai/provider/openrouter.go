package provider

import (
	"context"
	"fmt"

	"github.com/krishkalaria12/nyron-ai-cli/ai/tools"
	"github.com/krishkalaria12/nyron-ai-cli/config"
	openrouter "github.com/revrost/go-openrouter"
)

// OpenRouterAPI makes a single API call to OpenRouter with the given message history.
// It no longer loops; the conversational loop is now managed by the TUI.
func OpenRouterAPI(messages []openrouter.ChatCompletionMessage, model string) (openrouter.ChatCompletionResponse, error) {
	client := openrouter.NewClient(
		config.Config("OPENROUTER_API_KEY"),
	)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			Tools:    tools.GetAllTools(),
		},
	)

	if err != nil {
		return openrouter.ChatCompletionResponse{}, fmt.Errorf("ChatCompletion error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return openrouter.ChatCompletionResponse{}, fmt.Errorf("API returned no choices")
	}

	return resp, nil
}
