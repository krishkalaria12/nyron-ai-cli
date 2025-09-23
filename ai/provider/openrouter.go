package provider

import (
	"context"
	"fmt"

	"github.com/krishkalaria12/nyron-ai-cli/ai/tools"
	"github.com/krishkalaria12/nyron-ai-cli/config"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
	openrouter "github.com/revrost/go-openrouter"
)

// OpenRouterAPI generates a complete response using OpenRouter API
func OpenRouterAPI(prompt string, model string) AIResponseMessage {
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
			Tools: tools.GetAllTools(),
		},
	)

	var res AIResponseMessage

	if err != nil {
		res = AIResponseMessage{
			Thinking: "",
			Content:  "",
			Err:      fmt.Errorf("ChatCompletion error: %v\n", err),
		}
		return res
	}

	thinking := ""

	if len(resp.Choices) > 0 {
		if resp.Choices[0].Message.Reasoning != nil {
			thinking = *resp.Choices[0].Message.Reasoning
		}
	}

	res = AIResponseMessage{
		Thinking: thinking,
		Content:  resp.Choices[0].Message.Content.Text,
		Err:      nil,
	}

	return res
}
