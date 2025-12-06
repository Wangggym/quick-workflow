package logger

// HandlerOptions configures a handler
type HandlerOptions struct {
	Level      Level
	AddSource  bool
	JSON       bool // Use JSON format (for file logs)
	Buffered   bool // Enable buffered writes (default: true for file handlers)
	BufferSize int  // Buffer size in bytes (default: 4096)
}
