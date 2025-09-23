package tools

import "github.com/revrost/go-openrouter"

func GetAllTools() []openrouter.Tool {
	return []openrouter.Tool{
		CreateTool,
		EditTool,
		WriteContentTool,
		ListDirectoryTool,
		SearchFilesTool,
		ReadFileTool,
		GetFileInfoTool,
		GetCurrentDirectoryTool,
		WebSearchTool,
	}
}
