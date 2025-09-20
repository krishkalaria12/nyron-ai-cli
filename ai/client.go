package ai

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/krishkalaria12/nyron-ai-cli/config"
	prompts "github.com/krishkalaria12/nyron-ai-cli/config/prompts"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"google.golang.org/genai"
)

var (
	GeminiClient *genai.Client
	OpenaiClient openai.Client
	once         sync.Once
)

func getGeminiClient() *genai.Client {
	once.Do(func() {
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

func getOpenaiClient() openai.Client {
	once.Do(func() {
		OpenaiClient = openai.NewClient(
			option.WithAPIKey(config.Config("OPENAI_API_KEY")),
		)
	})

	return OpenaiClient
}

// StreamMessage represents a chunk of response or an error
type StreamMessage struct {
	Content string
	Error   error
	Done    bool
}

// GeminiStreamAPI streams the response in real-time through a channel
func GeminiStreamAPI(prompt string, responseChan chan<- StreamMessage) {
	defer close(responseChan)

	main_client := getGeminiClient()

	stream := main_client.Models.GenerateContentStream(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompts.FinalPrompt(prompt, "gemini")),
		nil,
	)

	for chunk, err := range stream {
		if err != nil {
			responseChan <- StreamMessage{
				Error: fmt.Errorf("error streaming response: %w", err),
				Done:  true,
			}
			return
		}

		part := chunk.Candidates[0].Content.Parts[0]
		responseChan <- StreamMessage{
			Content: part.Text,
			Error:   nil,
			Done:    false,
		}
	}

	// Send done signal
	responseChan <- StreamMessage{
		Done: true,
	}
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
