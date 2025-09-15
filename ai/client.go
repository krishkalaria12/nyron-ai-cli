package ai

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/krishkalaria12/nyron-ai-cli/config"
	"google.golang.org/genai"
)

var (
	client *genai.Client
	once   sync.Once
)

func getClient() *genai.Client {
	once.Do(func() {
		ctx := context.Background()
		var err error

		client, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: config.Config("GEMINI_API_KEY"),
		})

		if err != nil {
			log.Fatal(err)
		}
	})

	return client
}

func GeminiAPI(prompt string) (string, error) {
	main_client := getClient()

	stream := main_client.Models.GenerateContentStream(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	var response string
	for chunk, err := range stream {
		if err != nil {
			return "", fmt.Errorf("error streaming response: %w", err)
		}
		part := chunk.Candidates[0].Content.Parts[0]
		response += part.Text
	}

	return response, nil
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

	main_client := getClient()

	stream := main_client.Models.GenerateContentStream(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompt),
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
