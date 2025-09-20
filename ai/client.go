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

// OpenAIStreamAPI streams the response in real-time through a channel
func OpenAIStreamAPI(prompt string, responseChan chan<- StreamMessage) {
	defer close(responseChan)

	main_client := getOpenaiClient()

	ctx := context.Background()

	stream := main_client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompts.FinalPrompt(prompt, "openai")),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT5Mini,
	})

	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		// Handle finished content
		if content, ok := acc.JustFinishedContent(); ok {
			responseChan <- StreamMessage{
				Content: content,
				Error:   nil,
				Done:    false,
			}
		}

		// Handle streaming chunks
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			responseChan <- StreamMessage{
				Content: chunk.Choices[0].Delta.Content,
				Error:   nil,
				Done:    false,
			}
		}
	}

	if stream.Err() != nil {
		responseChan <- StreamMessage{
			Error: stream.Err(),
			Done:  true,
		}
		return
	}

	// Send done signal
	responseChan <- StreamMessage{
		Done: true,
	}
}
