package tools

import (
	"encoding/json"
)

func ExecuteTool(toolName string, arguements string) string {
	var response ToolResponse
	switch toolName {
	case "create_file_or_folder":
		createPar := CreateParams{}
		if err := json.Unmarshal([]byte(arguements), &createPar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := CreateFileOrFolder(createPar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "edit_content":
		editPar := EditParams{}
		if err := json.Unmarshal([]byte(arguements), &editPar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := EditFileContent(editPar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "write_content":
		writePar := WriteParams{}
		if err := json.Unmarshal([]byte(arguements), &writePar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := WriteContent(writePar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "list_directory":
		listPar := ListDirectoryParams{}
		if err := json.Unmarshal([]byte(arguements), &listPar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := ListDirectory(listPar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "search_files":
		searchPar := SearchFilesParams{}
		if err := json.Unmarshal([]byte(arguements), &searchPar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := SearchFiles(searchPar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "read_file":
		readPar := ReadFileParams{}
		if err := json.Unmarshal([]byte(arguements), &readPar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := ReadFile(readPar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "get_current_directory":
		result, err := GetCurrentDirectory()
		response = ToolResponse{Result: result, Error: err}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)

	case "web_search":
		webSearchPar := WebSearchParams{}
		if err := json.Unmarshal([]byte(arguements), &webSearchPar); err != nil {
			response = ToolResponse{
				Result: nil,
				Error: ToolError{
					Success: false,
					Message: "Unknown tool or invalid parameters",
					Err:     err,
				},
			}
		} else {
			result, err := WebSearch(webSearchPar)
			response = ToolResponse{Result: result, Error: err}
		}

		responseStr, _ := json.Marshal(response)
		return string(responseStr)
	}

	response = ToolResponse{
		Result: nil,
		Error: ToolError{
			Success: false,
			Message: "Unknown tool or invalid parameters",
			Err:     nil,
		},
	}

	responseStr, _ := json.Marshal(response)
	return string(responseStr)
}
