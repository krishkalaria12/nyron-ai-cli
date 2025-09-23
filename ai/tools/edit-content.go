package tools

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

type EditParams struct {
	FilePath        string
	EditMode        string
	SearchText      string
	ReplacementText string
	LineNumber      int
}

type EditResult struct {
	Success  bool
	Message  string
	Path     string
	EditMode string
}

var EditToolParams = jsonschema.Definition{
	Type: jsonschema.Object,
	Properties: map[string]jsonschema.Definition{
		"FilePath": {
			Type:        jsonschema.String,
			Description: "Path to the file to edit",
		},
		"EditMode": {
			Type:        jsonschema.String,
			Description: "Edit mode: 'replace', 'insert_at_line', 'append_line', 'prepend_line', 'replace_line'",
		},
		"SearchText": {
			Type:        jsonschema.String,
			Description: "Text to search for (required for 'replace' mode)",
		},
		"ReplacementText": {
			Type:        jsonschema.String,
			Description: "Text to replace with or insert",
		},
		"LineNumber": {
			Type:        jsonschema.Integer,
			Description: "Line number (required for 'insert_at_line' and 'replace_line' modes)",
		},
	},
	Required: []string{
		"FilePath",
		"EditMode",
	},
}

var EditOpenrouterFn = openrouter.FunctionDefinition{
	Name:        "edit_content",
	Description: "Edit file content with various editing modes",
	Parameters:  EditToolParams,
}

var EditTool = openrouter.Tool{
	Type:     openrouter.ToolTypeFunction,
	Function: &EditOpenrouterFn,
}

func EditFileContent(params EditParams) (EditResult, ToolError) {
	// Read the current file content
	currentContentBytes, err := os.ReadFile(params.FilePath)
	if err != nil {
		return EditResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error reading file: %s", err.Error()),
			Err:     err,
		}
	}

	currentContent := string(currentContentBytes)
	var newContent string

	switch params.EditMode {
	case "replace":
		if params.SearchText == "" {
			return EditResult{}, ToolError{
				Success: false,
				Message: "search_text is required for replace mode",
				Err:     fmt.Errorf("search_text is required for replace mode"),
			}
		}

		// Use regex for global replacement (equivalent to 'g' flag in JS)
		re := regexp.MustCompile(regexp.QuoteMeta(params.SearchText))
		newContent = re.ReplaceAllString(currentContent, params.ReplacementText)

	case "insert_at_line":
		if params.LineNumber == 0 {
			return EditResult{}, ToolError{
				Success: false,
				Message: "line_number is required for insert_at_line mode",
				Err:     fmt.Errorf("line_number is required for insert_at_line mode"),
			}
		}

		lines := strings.Split(currentContent, "\n")

		// Check if line number is valid
		if params.LineNumber > len(lines)+1 {
			return EditResult{}, ToolError{
				Success: false,
				Message: fmt.Sprintf("Line %d is out of range. File has %d lines.", params.LineNumber, len(lines)),
				Err:     fmt.Errorf("line number out of range"),
			}
		}

		// Insert at specified line (1-based indexing)
		insertIndex := params.LineNumber - 1
		if insertIndex < 0 {
			insertIndex = 0
		}

		// Create new slice with inserted line
		newLines := make([]string, 0, len(lines)+1)
		newLines = append(newLines, lines[:insertIndex]...)
		newLines = append(newLines, params.ReplacementText)
		newLines = append(newLines, lines[insertIndex:]...)

		newContent = strings.Join(newLines, "\n")

	case "append_line":
		newContent = currentContent + "\n" + params.ReplacementText

	case "prepend_line":
		newContent = params.ReplacementText + "\n" + currentContent

	case "replace_line":
		if params.LineNumber == 0 {
			return EditResult{}, ToolError{
				Success: false,
				Message: "line_number is required for replace_line mode",
				Err:     fmt.Errorf("line_number is required for replace_line mode"),
			}
		}

		lines := strings.Split(currentContent, "\n")

		if params.LineNumber > len(lines) {
			return EditResult{}, ToolError{
				Success: false,
				Message: fmt.Sprintf("Line %d does not exist. File has %d lines.", params.LineNumber, len(lines)),
				Err:     fmt.Errorf("line number out of range"),
			}
		}

		// Replace the specified line (1-based indexing)
		lines[params.LineNumber-1] = params.ReplacementText
		newContent = strings.Join(lines, "\n")

	default:
		return EditResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Unknown edit mode: %s", params.EditMode),
			Err:     fmt.Errorf("unknown edit mode"),
		}
	}

	// Write the modified content back to the file
	err = os.WriteFile(params.FilePath, []byte(newContent), 0644)
	if err != nil {
		return EditResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error writing file: %s", err.Error()),
			Err:     err,
		}
	}

	return EditResult{
		Success:  true,
		Message:  fmt.Sprintf("File edited successfully using %s mode", params.EditMode),
		Path:     params.FilePath,
		EditMode: params.EditMode,
	}, ToolError{}
}
