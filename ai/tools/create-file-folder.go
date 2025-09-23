package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

type CreateParams struct {
	BasePath     string
	TypeOfCreate string
	Name         string
}

type CreateSuccess struct {
	Success bool
	Message string
	Path    string
}

func CreateFileOrFolder(params CreateParams) (CreateSuccess, ToolError) {
	// Join the base path with the name to get full path
	fullPath := filepath.Join(params.BasePath, params.Name)

	switch params.TypeOfCreate {
	case "folder":
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error creating folder: %s", err.Error()),
				Err:     err,
			}
			return CreateSuccess{}, toolErr
		}

		toolSuc := CreateSuccess{
			Success: true,
			Message: fmt.Sprintf("Folder '%s' created at %s", params.Name, fullPath),
			Path:    fullPath,
		}
		return toolSuc, ToolError{}

	case "file":
		err := os.MkdirAll(params.BasePath, 0755)
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error creating directory: %s", err.Error()),
				Err:     err,
			}
			return CreateSuccess{}, toolErr
		}

		file, err := os.Create(fullPath)
		defer file.Close()
		if err != nil {
			toolErr := ToolError{
				Success: false,
				Message: fmt.Sprintf("Error creating file: %s", err.Error()),
				Err:     err,
			}
			return CreateSuccess{}, toolErr
		}

		toolSuc := CreateSuccess{
			Success: true,
			Message: fmt.Sprintf("File '%s' created at %s", params.Name, fullPath),
			Path:    fullPath,
		}

		return toolSuc, ToolError{}

	default:
		toolErr := ToolError{
			Success: false,
			Message: fmt.Sprintf("Invalid type: %s. Must be 'file' or 'folder'", params.TypeOfCreate),
			Err:     fmt.Errorf("invalid type: %s", params.TypeOfCreate),
		}
		return CreateSuccess{}, toolErr
	}
}
