package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/krishkalaria12/nyron-ai-cli/config"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var (
	OpenaiClient openai.Client
	openaiOnce   sync.Once
)

func getOpenaiClient() openai.Client {
	openaiOnce.Do(func() {
		OpenaiClient = openai.NewClient(
			option.WithAPIKey(config.Config("OPENAI_API_KEY")),
		)
	})

	return OpenaiClient
}

// OpenAIAPI generates a complete response using OpenAI API
func OpenAIAPI(prompt string) (string, error) {
	main_client := getOpenaiClient()

	chatCompletion, err := main_client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompts.FinalPrompt(prompt, "openai")),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT5Mini,
	})

	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}