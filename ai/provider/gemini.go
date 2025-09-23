package provider

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/krishkalaria12/nyron-ai-cli/config"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
	"google.golang.org/genai"
)

var (
	GeminiClient *genai.Client
	geminiOnce   sync.Once
)

func getGeminiClient() *genai.Client {
	geminiOnce.Do(func() {
		ctx := context.Background()
		var err error

		GeminiClient, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: config.Config("GEMINI_API_KEY"),
		})

		if err != nil {
			log.Fatal(err)
		}
	})

	return GeminiClient
}

// GeminiAPI generates a complete response using Gemini API
func GeminiAPI(prompt string, model string) AIResponseMessage {
	main_client := getGeminiClient()

	result, err := main_client.Models.GenerateContent(
		context.Background(),
		model,
		genai.Text(prompts.FinalPrompt(prompt, "gemini")),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
			},
			// Tools: ,
		},
	)

	var res AIResponseMessage

	if err != nil {
		res = AIResponseMessage{
			Thinking: "",
			Content:  "",
			Err:      fmt.Errorf("error generating response: %w", err),
		}
		return res
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		res = AIResponseMessage{
			Thinking: "",
			Content:  "",
			Err:      fmt.Errorf("No Response Generated"),
		}
		return res
	}

	thinking := ""
	response := ""

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			if part.Thought {
				thinking += part.Text
			} else {
				response += part.Text
			}
		}
	}

	res = AIResponseMessage{
		Thinking: thinking,
		Content:  response,
		Err:      nil,
	}

	return res
}
