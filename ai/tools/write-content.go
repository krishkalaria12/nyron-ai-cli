package tools

import (
	"fmt"
	"os"

	"github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

type WriteParams struct {
	FilePath string
	Content  string
	Mode     string
}

type WriteSuccess struct {
	Success bool
	Message string
	Path    string
}

var WriteToolParams = jsonschema.Definition{
	Type: jsonschema.Object,
	Properties: map[string]jsonschema.Definition{
		"FilePath": {
			Type:        jsonschema.String,
			Description: "Path to the file to write content to",
		},
		"Content": {
			Type:        jsonschema.String,
			Description: "Content to write to the file",
		},
		"Mode": {
			Type:        jsonschema.String,
			Description: "Write mode: 'overwrite' or 'append'",
		},
	},
	Required: []string{
		"FilePath",
		"Content",
		"Mode",
	},
}

var WriteOpenrouterFn = openrouter.FunctionDefinition{
	Name:        "write_content",
	Description: "Write content to a file with overwrite or append mode",
	Parameters:  WriteToolParams,
}

var WriteContentTool = openrouter.Tool{
	Type:     openrouter.ToolTypeFunction,
	Function: &WriteOpenrouterFn,
}

func WriteContent(params WriteParams) (WriteSuccess, ToolError) {
	switch params.Mode {
	case "overwrite":
		err := os.WriteFile(params.FilePath, []byte(params.Content), 0644)
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error overwriting file content: %s", err.Error()),
				Err:     err,
			}
			return WriteSuccess{}, toolErr
		}

		toolSuc := WriteSuccess{
			Success: true,
			Message: fmt.Sprintf("Content '%s' filepath %s", params.Mode, params.FilePath),
			Path:    params.FilePath,
		}
		return toolSuc, ToolError{}

	case "append":
		file, err := os.OpenFile(params.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error writing content to the file: %s", err.Error()),
				Err:     err,
			}
			return WriteSuccess{}, toolErr
		}
		defer file.Close()

		_, err = file.WriteString(params.Content)
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error writing content to the file: %s", err.Error()),
				Err:     err,
			}
			return WriteSuccess{}, toolErr
		}

		toolSuc := WriteSuccess{
			Success: true,
			Message: fmt.Sprintf("Content '%s' filepath %s", params.Mode, params.FilePath),
			Path:    params.FilePath,
		}
		return toolSuc, ToolError{}

	default:
		return WriteSuccess{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Invalid mode: %s", params.Mode),
			Err:     fmt.Errorf("invalid mode: %s", params.Mode),
		}
	}
}
