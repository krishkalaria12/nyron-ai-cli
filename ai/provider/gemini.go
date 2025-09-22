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
func GeminiAPI(prompt string) (string, error) {
	main_client := getGeminiClient()

	result, err := main_client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompts.FinalPrompt(prompt, "gemini")),
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}