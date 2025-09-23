package tools

type ToolError struct {
	Success bool
	Message string
	Err     error
}

type ToolResponse struct {
	Result interface{}
	Error  ToolError
}
