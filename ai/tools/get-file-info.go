package tools

import (
	"fmt"
	"os"
	"time"
)

type GetFileInfoParams struct {
	Path string
}

type FileInfo struct {
	Type         string
	Size         int64
	Created      time.Time
	Modified     time.Time
	Accessed     time.Time
	Permissions  string
	IsReadable   bool
	IsWritable   bool
	IsExecutable bool
}

type GetFileInfoResult struct {
	Success bool
	Message string
	Path    string
	Info    FileInfo
}

func GetFileInfo(params GetFileInfoParams) (GetFileInfoResult, ToolError) {
	info, err := os.Stat(params.Path)
	if err != nil {
		return GetFileInfoResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error getting file info: %s", err.Error()),
			Err:     err,
		}
	}

	fileType := "file"
	if info.IsDir() {
		fileType = "folder"
	}

	mode := info.Mode()
	permissions := fmt.Sprintf("0%o", mode&0777)

	fileInfo := FileInfo{
		Type:         fileType,
		Size:         info.Size(),
		Modified:     info.ModTime(),
		Permissions:  permissions,
		IsReadable:   mode&0444 != 0,
		IsWritable:   mode&0222 != 0,
		IsExecutable: mode&0111 != 0,
	}

	// Note: Go's os.FileInfo doesn't provide creation and access times directly
	// These would need platform-specific implementations
	fileInfo.Created = info.ModTime()  // Fallback to modification time
	fileInfo.Accessed = info.ModTime() // Fallback to modification time

	return GetFileInfoResult{
		Success: true,
		Message: fmt.Sprintf("Retrieved info for %s", params.Path),
		Path:    params.Path,
		Info:    fileInfo,
	}, ToolError{}
}