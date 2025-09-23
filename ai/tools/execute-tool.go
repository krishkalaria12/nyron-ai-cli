package tools

func ExecuteTool(toolName string, parameters interface{}) ToolResponse {
	switch toolName {
	case "create_file_or_folder":
		if params, ok := parameters.(CreateParams); ok {
			result, err := CreateFileOrFolder(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "edit_content":
		if params, ok := parameters.(EditParams); ok {
			result, err := EditFileContent(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "write_content":
		if params, ok := parameters.(WriteParams); ok {
			result, err := WriteContent(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "list_directory":
		if params, ok := parameters.(ListDirectoryParams); ok {
			result, err := ListDirectory(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "search_files":
		if params, ok := parameters.(SearchFilesParams); ok {
			result, err := SearchFiles(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "read_file":
		if params, ok := parameters.(ReadFileParams); ok {
			result, err := ReadFile(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "get_file_info":
		if params, ok := parameters.(GetFileInfoParams); ok {
			result, err := GetFileInfo(params)
			return ToolResponse{Result: result, Error: err}
		}
	case "get_current_directory":
		result, err := GetCurrentDirectory()
		return ToolResponse{Result: result, Error: err}
	case "web_search":
		if params, ok := parameters.(WebSearchParams); ok {
			result, err := WebSearch(params)
			return ToolResponse{Result: result, Error: err}
		}
	}
	return ToolResponse{
		Result: nil,
		Error: ToolError{
			Success: false,
			Message: "Unknown tool or invalid parameters",
			Err:     nil,
		},
	}
}