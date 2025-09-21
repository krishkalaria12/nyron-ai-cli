# Nyron AI CLI

A beautiful terminal-based AI chat interface built in Go that supports multiple AI providers including Gemini, OpenAI, and OpenRouter.

## Features ✨

- **Multi-Provider Support**: Switch between Gemini, OpenAI, and OpenRouter AI models
- **Beautiful TUI**: Modern terminal interface built with Bubble Tea
- **Model Selection**: Dynamic model selection dialog with `Ctrl+P`
- **Markdown Rendering**: Rich markdown support for AI responses
- **Responsive Design**: Adapts to terminal size changes
- **Real-time Chat**: Smooth conversational experience with loading indicators

## Supported AI Providers

- **Google Gemini** (gemini-2.5-flash)
- **OpenAI** (GPT-5 Mini)
- **OpenRouter** (Various models)

## Installation

### Prerequisites

- Go 1.25.1 or later
- API keys for your preferred AI providers

### Setup

1. Clone the repository:
```bash
git clone https://github.com/krishkalaria12/nyron-ai-cli.git
cd nyron-ai-cli
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file with your API keys:
```env
GEMINI_API_KEY=your_gemini_api_key_here
OPENAI_API_KEY=your_openai_api_key_here
OPENROUTER_API_KEY=your_openrouter_api_key_here
```

4. Run the application:
```bash
go run .
```

## Usage

### Basic Commands

- **Type your message** and press `Enter` to send
- **Shift+Enter** or **Ctrl+J** to add new lines without sending
- **Tab** to switch focus between chat history and input
- **Ctrl+P** to open model selection dialog
- **↑/↓** or **k/j** to scroll through chat history (when focused on viewport)
- **Page Up/Down** or **Ctrl+U/Ctrl+D** for page navigation
- **Ctrl+C** to quit

### Model Selection

Press `Ctrl+P` to open the model selection dialog where you can choose between:
- Gemini models
- OpenAI models
- OpenRouter models

## Project Structure

```
├── ai/                     # AI client implementations
│   ├── client.go          # API clients for different providers
│   └── markdown-renderer.go # Markdown rendering utilities
├── config/                # Configuration management
│   ├── config.go          # Environment configuration
│   ├── models.go          # Model definitions
│   └── prompts/           # System prompts
├── tui/                   # Terminal UI components
│   ├── components/        # Reusable UI components
│   │   ├── chat/          # Main chat interface
│   │   ├── dialogs/       # Modal dialogs
│   │   └── editor/        # Input editor
│   └── runner.go          # TUI runner
├── util/                  # Utility functions
└── main.go               # Application entry point
```

## Configuration

The application uses environment variables for configuration:

- `GEMINI_API_KEY`: Your Google Gemini API key
- `OPENAI_API_KEY`: Your OpenAI API key
- `OPENROUTER_API_KEY`: Your OpenRouter API key

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## API Key Setup

### Google Gemini
1. Visit [Google AI Studio](https://aistudio.google.com/)
2. Create an API key
3. Add it to your `.env` file

### OpenAI
1. Visit [OpenAI API](https://platform.openai.com/api-keys)
2. Create an API key
3. Add it to your `.env` file

### OpenRouter
1. Visit [OpenRouter](https://openrouter.ai/)
2. Create an account and get an API key
3. Add it to your `.env` file

---

Built with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and Go