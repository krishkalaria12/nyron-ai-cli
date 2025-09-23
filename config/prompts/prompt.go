package config

import (
	"bytes"
	_ "embed"
	"text/template"
)

// Embed the markdown files at compile time

//go:embed openrouter.md
var openrouterSysPrompt string

type PromptPair struct {
	SystemPrompt string
	UserPrompt   string
}

func GetPrompts(userPrompt string, provider string) PromptPair {
	// Template for user prompt
	userPromptTemplate := `{{.UserPrompt}}`

	// Parse the user prompt template
	tmpl, err := template.New("userPrompt").Parse(userPromptTemplate)
	if err != nil {
		panic(err) // Handle error appropriately in production
	}

	// Execute template with user data
	var userPromptBuffer bytes.Buffer
	data := struct {
		UserPrompt string
	}{
		UserPrompt: userPrompt,
	}

	err = tmpl.Execute(&userPromptBuffer, data)
	if err != nil {
		panic(err)
	}

	formattedUserPrompt := userPromptBuffer.String()

	// Return separate system and user prompts based on provider
	switch provider {
	case "openrouter":
		return PromptPair{
			SystemPrompt: openrouterSysPrompt,
			UserPrompt:   formattedUserPrompt,
		}
	default:
		return PromptPair{
			SystemPrompt: "",
			UserPrompt:   formattedUserPrompt,
		}
	}
}
