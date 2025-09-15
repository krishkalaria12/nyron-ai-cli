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

func GeminiAPI() {
	main_client := getClient()

	stream := main_client.Models.GenerateContentStream(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text("Explain AI in simple terms"),
		nil,
	)

	for chunk, _ := range stream {
		part := chunk.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)
	}
}
