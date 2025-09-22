package provider

import (
	"context"
	"fmt"

	"github.com/krishkalaria12/nyron-ai-cli/config"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
	openrouter "github.com/revrost/go-openrouter"
)

// OpenRouterAPI generates a complete response using OpenRouter API
func OpenRouterAPI(prompt string, model string) (string, error) {
	client := openrouter.NewClient(
		config.Config("OPENROUTER_API_KEY"),
	)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: model,
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(prompts.FinalPrompt(prompt, "openrouter")),
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Message.Content.Text, nil
}