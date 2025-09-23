package tools

import (
	"fmt"
	"os"

	"github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

type GetCurrentDirectoryResult struct {
	Success          bool
	Message          string
	CurrentDirectory string
}

var GetCurrentDirectoryToolParams = jsonschema.Definition{
	Type:       jsonschema.Object,
	Properties: map[string]jsonschema.Definition{},
	Required:   []string{},
}

var GetCurrentDirectoryOpenrouterFn = openrouter.FunctionDefinition{
	Name:        "get_current_directory",
	Description: "Get the current working directory path",
	Parameters:  GetCurrentDirectoryToolParams,
}

var GetCurrentDirectoryTool = openrouter.Tool{
	Type:     openrouter.ToolTypeFunction,
	Function: &GetCurrentDirectoryOpenrouterFn,
}

func GetCurrentDirectory() (GetCurrentDirectoryResult, ToolError) {
	currentDir, err := os.Getwd()
	if err != nil {
		return GetCurrentDirectoryResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error getting current directory: %s", err.Error()),
			Err:     err,
		}
	}

	return GetCurrentDirectoryResult{
		Success:          true,
		Message:          fmt.Sprintf("Current directory: %s", currentDir),
		CurrentDirectory: currentDir,
	}, ToolError{}
}