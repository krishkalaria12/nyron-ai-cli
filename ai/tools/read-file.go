package tools

import (
	"fmt"
	"os"
	"strings"
)

type ReadFileParams struct {
	FilePath string
	Encoding string
	MaxLines *int
}

type ReadFileResult struct {
	Success    bool
	Message    string
	Path       string
	Content    string
	Truncated  bool
	TotalLines int
}

func ReadFile(params ReadFileParams) (ReadFileResult, ToolError) {
	content, err := os.ReadFile(params.FilePath)
	if err != nil {
		return ReadFileResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error reading file: %s", err.Error()),
			Err:     err,
		}
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	totalLines := len(lines)

	displayContent := contentStr
	truncated := false

	if params.MaxLines != nil && *params.MaxLines > 0 {
		if totalLines > *params.MaxLines {
			displayContent = strings.Join(lines[:*params.MaxLines], "\n")
			truncated = true
		}
	}

	message := fmt.Sprintf("Read file %s", params.FilePath)
	if truncated {
		message = fmt.Sprintf("Read file %s (showing first %d lines)", params.FilePath, *params.MaxLines)
	}

	return ReadFileResult{
		Success:    true,
		Message:    message,
		Path:       params.FilePath,
		Content:    displayContent,
		Truncated:  truncated,
		TotalLines: totalLines,
	}, ToolError{}
}