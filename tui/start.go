package tui

import (
	"fmt"
	"os"

	"github.com/krishkalaria12/nyron-ai-cli/tui/components"
)


func StartTUI() {
	prompt := components.RunInputModel()

	if prompt == "" {
		fmt.Println("No prompt provided")
		os.Exit(1)
	}

	// Use streaming response model
	components.RunStreamingResponseModel(prompt)
}
