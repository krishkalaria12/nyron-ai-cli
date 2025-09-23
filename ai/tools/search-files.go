package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

type SearchFilesParams struct {
	SearchPath string
	Pattern    string
	Recursive  bool
	Type       string
}

type SearchResult struct {
	Name string
	Type string
	Path string
	Size *int64
}

type SearchFilesResult struct {
	Success    bool
	Message    string
	Pattern    string
	SearchPath string
	Results    []SearchResult
}

var SearchFilesToolParams = jsonschema.Definition{
	Type: jsonschema.Object,
	Properties: map[string]jsonschema.Definition{
		"SearchPath": {
			Type:        jsonschema.String,
			Description: "Path to search in (default: current directory)",
		},
		"Pattern": {
			Type:        jsonschema.String,
			Description: "File name pattern to search for (supports wildcards)",
		},
		"Recursive": {
			Type:        jsonschema.Boolean,
			Description: "Whether to search recursively in subdirectories",
		},
		"Type": {
			Type:        jsonschema.String,
			Description: "Filter by type: 'files', 'folders', or empty for all",
		},
	},
	Required: []string{
		"Pattern",
	},
}

var SearchFilesOpenrouterFn = openrouter.FunctionDefinition{
	Name:        "search_files",
	Description: "Search for files and folders by pattern with optional filtering",
	Parameters:  SearchFilesToolParams,
}

var SearchFilesTool = openrouter.Tool{
	Type:     openrouter.ToolTypeFunction,
	Function: &SearchFilesOpenrouterFn,
}

func SearchFiles(params SearchFilesParams) (SearchFilesResult, ToolError) {
	searchPath := params.SearchPath
	if searchPath == "" {
		searchPath = "."
	}

	var results []SearchResult

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking even if there's an error with one file
		}

		// If not recursive, only check the immediate directory
		if !params.Recursive && filepath.Dir(path) != searchPath {
			if info.IsDir() && path != searchPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip the root directory itself
		if path == searchPath {
			return nil
		}

		// Check if the file matches the pattern
		matched, err := filepath.Match(params.Pattern, info.Name())
		if err != nil {
			return nil // Continue if pattern matching fails
		}

		if matched {
			itemType := "file"
			if info.IsDir() {
				itemType = "folder"
			}

			// Filter by type
			if params.Type == "files" && itemType != "file" {
				return nil
			}
			if params.Type == "folders" && itemType != "folder" {
				return nil
			}

			var size *int64
			if itemType == "file" {
				fileSize := info.Size()
				size = &fileSize
			}

			results = append(results, SearchResult{
				Name: info.Name(),
				Type: itemType,
				Path: path,
				Size: size,
			})
		}

		return nil
	})

	if err != nil {
		return SearchFilesResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error searching files: %s", err.Error()),
			Err:     err,
		}
	}

	return SearchFilesResult{
		Success:    true,
		Message:    fmt.Sprintf("Found %d matches for pattern \"%s\"", len(results), params.Pattern),
		Pattern:    params.Pattern,
		SearchPath: searchPath,
		Results:    results,
	}, ToolError{}
}