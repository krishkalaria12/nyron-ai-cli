package tools

import (
	"fmt"
	"os"
)

type WriteParams struct {
	file_path string
	content   string
	mode      string
}

type WriteSuccess struct {
	Success bool
	Message string
	Path    string
}

func WriteContent(params WriteParams) (WriteSuccess, ToolError) {
	switch params.mode {
	case "overwrite":
		err := os.WriteFile(params.file_path, []byte(params.content), 0644)
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
			Message: fmt.Sprintf("Content '%s' filepath %s", params.mode, params.file_path),
			Path:    params.file_path,
		}
		return toolSuc, ToolError{}

	case "append":
		file, err := os.OpenFile(params.file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error writing content to the file: %s", err.Error()),
				Err:     err,
			}
			return WriteSuccess{}, toolErr
		}
		defer file.Close()

		_, err = file.WriteString(params.content)
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
			Message: fmt.Sprintf("Content '%s' filepath %s", params.mode, params.file_path),
			Path:    params.file_path,
		}
		return toolSuc, ToolError{}

	default:
		return WriteSuccess{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Invalid mode: %s", params.mode),
			Err:     fmt.Errorf("invalid mode: %s", params.mode),
		}
	}
}
