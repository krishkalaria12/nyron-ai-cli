package tools

import (
	"fmt"
	"os"
)

type GetCurrentDirectoryResult struct {
	Success          bool
	Message          string
	CurrentDirectory string
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