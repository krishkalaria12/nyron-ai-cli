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
	ProviderGemini = Provider{
		ID:   "gemini",
		Name: "Google Gemini",
	}
	ProviderOpenAI = Provider{
		ID:   "openai",
		Name: "OpenAI",
	}
	ProviderOpenRouter = Provider{
		ID:   "openrouter",
		Name: "OpenRouter",
	}
)

// Available models by provider
var (
	GeminiModels = []Model{
		{
			ID:          "gemini-2.5-flash",
			Name:        "Gemini 2.5 Flash",
			Description: "Fast and efficient model for quick responses",
		},
		{
			ID:          "gemini-2.5-pro",
			Name:        "Gemini 2.5 Pro",
			Description: "Advanced model with superior reasoning capabilities",
		},
	}

	OpenAIModels = []Model{
		{
			ID:          "gpt-5-2025-08-07",
			Name:        "GPT-5",
			Description: "Latest GPT-5 model with enhanced capabilities",
		},
		{
			ID:          "gpt-5-mini-2025-08-07",
			Name:        "GPT-5 Mini",
			Description: "Lightweight version of GPT-5 for faster responses",
		},
		{
			ID:          "gpt-4o-2024-08-06",
			Name:        "GPT-4o",
			Description: "Optimized GPT-4 model for improved performance",
		},
	}

	OpenRouterModels = []Model{
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
		ProviderGemini,
		ProviderOpenAI,
		ProviderOpenRouter,
	}
}

// GetModelsByProvider returns models for a specific provider
func GetModelsByProvider(providerID string) []Model {
	switch providerID {
	case ProviderGemini.ID:
		return GeminiModels
	case ProviderOpenAI.ID:
		return OpenAIModels
	case ProviderOpenRouter.ID:
		return OpenRouterModels
	default:
		return []Model{}
	}
}
