package provider

import (
	"context"
	"fmt"

	"github.com/krishkalaria12/nyron-ai-cli/ai/tools"
	"github.com/krishkalaria12/nyron-ai-cli/config"
	openrouter "github.com/revrost/go-openrouter"
)

// OpenRouterAPI generates a complete response using OpenRouter API
func OpenRouterAPI(systemPrompt, userPrompt string, model string, toolchan chan<- ToolCallingResponse) AIResponseMessage {
	client := openrouter.NewClient(
		config.Config("OPENROUTER_API_KEY"),
	)

	msgInput := []openrouter.ChatCompletionMessage{
		{
			Role: openrouter.ChatMessageRoleSystem,
			Content: openrouter.Content{
				Text: systemPrompt,
			},
		},
		{
			Role: openrouter.ChatMessageRoleUser,
			Content: openrouter.Content{
				Text: userPrompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model:    model,
			Messages: msgInput,
			Tools:    tools.GetAllTools(),
		},
	)

	var res AIResponseMessage

	if err != nil {
		res = AIResponseMessage{
			Thinking: "",
			Content:  "",
			Err:      fmt.Errorf("ChatCompletion error: %v\n", err),
		}
		return res
	}

	msg := resp.Choices[0].Message
	for len(msg.ToolCalls) > 0 {
		msgInput = append(msgInput, msg)

		tool_call_id := msg.ToolCalls[0].ID
		fn_name := msg.ToolCalls[0].Function.Name
		fn_arguements := msg.ToolCalls[0].Function.Arguments

		tool_response := tools.ExecuteTool(fn_name, fn_arguements)
		msgInput = append(msgInput, openrouter.ChatCompletionMessage{
			Role: openrouter.ChatMessageRoleTool,
			Content: openrouter.Content{
				Text: tool_response,
			},
			ToolCallID: tool_call_id,
		})

		toolchan <- ToolCallingResponse{
			Step:    fn_name,
			Content: fn_arguements,
		}

		// demonstrate the tool call in here
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openrouter.ChatCompletionRequest{
				Model:    model,
				Messages: msgInput,
				Tools:    tools.GetAllTools(),
			},
		)

		if err != nil || len(resp.Choices) != 1 {
			res = AIResponseMessage{
				Thinking: "",
				Content:  "",
				Err:      fmt.Errorf("Tool completion error: err:%v len(choices):%v\n", err, len(resp.Choices)),
			}

			return res
		}

		msg = resp.Choices[0].Message
	}

	close(toolchan)
	thinking := ""

	if msg.Reasoning != nil {
		thinking = *msg.Reasoning
	}

	res = AIResponseMessage{
		Thinking: thinking,
		Content:  msg.Content.Text,
		Err:      nil,
	}

	return res
}
