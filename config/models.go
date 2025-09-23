package config

type SelectedModel struct {
	// The model id as used by the provider API.
	// Required.
	Model string `json:"model" jsonschema:"required,description=The model ID as used by the provider API,example=gpt-4o"`
	// The model provider, same as the key/id used in the providers config.
	// Required.
	Provider string `json:"provider" jsonschema:"required,description=The model provider ID that matches a key in the providers config,example=openai"`
}

type Provider struct {
	ID   string
	Name string
}

type Model struct {
	ID          string
	Name        string
	Description string
}

// Available providers
var (
	ProviderOpenRouter = Provider{
		ID:   "openrouter",
		Name: "OpenRouter",
	}
)

// Available models by provider
var (
	OpenRouterModels = []Model{
		{
			ID:          "openai/gpt-5",
			Name:        "GPT 5",
			Description: "GPT-5 is OpenAI’s most advanced model, offering major improvements in reasoning, code quality, and user experience.",
		},
		{
			ID:          "openai/gpt-5-mini",
			Name:        "GPT 5 Mini",
			Description: "GPT-5 Mini is a compact version of GPT-5, designed to handle lighter-weight reasoning tasks.",
		},
		{
			ID:          "openai/gpt-4.1",
			Name:        "GPT 4.1",
			Description: "GPT-4.1 is a flagship large language model optimized for advanced instruction following, real-world software engineering, and long-context reasoning.",
		},
		{
			ID:          "google/gemini-2.5-pro",
			Name:        "Gemini 2.5 Pro",
			Description: "Gemini 2.5 Pro is Google’s state-of-the-art AI model designed for advanced reasoning, coding, mathematics, and scientific tasks.",
		},
		{
			ID:          "google/gemini-2.5-flash",
			Name:        "Gemini 2.5 Flash",
			Description: "Gemini 2.5 Flash is Google’s state-of-the-art AI model designed for advanced reasoning, coding, mathematics, and scientific tasks.",
		},
		{
			ID:          "x-ai/grok-4-fast:free",
			Name:        "Grok-4 Fast",
			Description: "xAI's Grok-4 model optimized for speed",
		},
		{
			ID:          "deepseek/deepseek-chat-v3.1:free",
			Name:        "Deepseek V3",
			Description: "Deepseek v3 is a large hybrid model",
		},
		{
			ID:          "z-ai/glm-4.5-air:free",
			Name:        "GLM 4.5 Air",
			Description: "GLM-4.5-Air is the lightweight variant of GLM 4.5",
		},
		{
			ID:          "moonshotai/kimi-k2:free",
			Name:        "Kimi K2",
			Description: "Kimi K2 Instruct is a large-scale Mixture-of-Experts (MoE) language model",
		},
	}
)

// GetAllProviders returns all available providers
func GetAllProviders() []Provider {
	return []Provider{
		ProviderOpenRouter,
	}
}

// GetModelsByProvider returns models for a specific provider
func GetModelsByProvider(providerID string) []Model {
	switch providerID {
	case ProviderOpenRouter.ID:
		return OpenRouterModels
	default:
		return []Model{}
	}
}
