package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ListDirectoryParams struct {
	DirectoryPath string
	ShowHidden    bool
	FilterType    string
}

type DirectoryItem struct {
	Name string
	Type string
	Path string
}

type ListDirectoryResult struct {
	Success   bool
	Message   string
	Directory string
	Items     []DirectoryItem
}

func ListDirectory(params ListDirectoryParams) (ListDirectoryResult, ToolError) {
	directoryPath := params.DirectoryPath
	if directoryPath == "" {
		directoryPath = "."
	}

	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return ListDirectoryResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error listing directory: %s", err.Error()),
			Err:     err,
		}
	}

	var items []DirectoryItem
	for _, entry := range entries {
		// Filter hidden files
		if !params.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		itemType := "file"
		if entry.IsDir() {
			itemType = "folder"
		}

		// Filter by type
		if params.FilterType == "files" && itemType != "file" {
			continue
		}
		if params.FilterType == "folders" && itemType != "folder" {
			continue
		}

		items = append(items, DirectoryItem{
			Name: entry.Name(),
			Type: itemType,
			Path: filepath.Join(directoryPath, entry.Name()),
		})
	}

	return ListDirectoryResult{
		Success:   true,
		Message:   fmt.Sprintf("Listed %d items in %s", len(items), directoryPath),
		Directory: directoryPath,
		Items:     items,
	}, ToolError{}
}