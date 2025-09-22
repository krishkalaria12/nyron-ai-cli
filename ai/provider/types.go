package provider

// StreamMessage represents a chunk of response or an error
type StreamMessage struct {
	Content string
	Error   error
	Done    bool
}

type AIResponseMessage struct {
	Thinking string
	Content  string
	Err      error
}
