package config

import (
	"bytes"
	_ "embed"
	"text/template"
)

// Embed the markdown files at compile time
//
//go:embed gemini.md
var geminiSysPrompt string

//go:embed openai.md
var openaiSysPrompt string

func FinalPrompt(userPrompt string, provider string) string {
	// Template for user prompt
	userPromptTemplate := `User Prompt: {{.UserPrompt}}`

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

	// Combine system prompt with user prompt based on provider
	switch provider {
	case "gemini":
		return geminiSysPrompt + "\n\n" + formattedUserPrompt
	case "openai":
		return openaiSysPrompt + "\n\n" + formattedUserPrompt
	default:
		return formattedUserPrompt // fallback to just user prompt
	}
}
