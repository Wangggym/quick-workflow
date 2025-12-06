package logger

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
	"time"
)

// FileHandler is a custom handler for file logging
// It outputs structured logs (JSON or text) with timestamps
// Supports buffered writes for better performance
type FileHandler struct {
	writer     io.Writer
	bufWriter  *bufio.Writer
	underlying io.Closer // The underlying file or writer that needs to be closed
	opts       *HandlerOptions
	mu         sync.RWMutex
	enabled    map[slog.Level]bool
}

// ------------------- Public Methods -------------------
// 外部访问

// NewFileHandler creates a new FileHandler
// If opts.Buffered is true (default), it wraps the writer with bufio.Writer for better performance
func NewFileHandler(writer io.Writer, opts *HandlerOptions) *FileHandler {
	if opts == nil {
		opts = &HandlerOptions{Level: LevelInfo, JSON: true, Buffered: true, BufferSize: 4096}
	}

	// Set defaults for buffering
	if opts.BufferSize <= 0 {
		opts.BufferSize = 4096 // Default 4KB buffer
	}
	// Default to buffered for file handlers (unless explicitly disabled)
	// This is a reasonable default for performance

	h := &FileHandler{
		opts:    opts,
		enabled: make(map[slog.Level]bool),
	}

	// Set up buffered writer if enabled
	if opts.Buffered {
		h.bufWriter = bufio.NewWriterSize(writer, opts.BufferSize)
		h.writer = h.bufWriter
	} else {
		h.writer = writer
	}

	// Store underlying closer if writer implements io.Closer
	if closer, ok := writer.(io.Closer); ok {
		h.underlying = closer
	}

	h.setEnabledLevels(opts.Level.ToSlogLevel())
	return h
}

// SetLevel changes the log level
func (h *FileHandler) SetLevel(level Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.opts.Level = level
	h.setEnabledLevels(level.ToSlogLevel())
}

// Enabled reports whether the handler handles records at the given level
func (h *FileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.enabled[level]
}

// Handle handles a log record
func (h *FileHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.Enabled(ctx, record.Level) {
		return nil
	}

	if h.opts.JSON {
		return h.handleJSON(record)
	}
	return h.handleText(record)
}

// WithAttrs returns a new handler with the given attributes
// Note: FileHandler doesn't store attributes at the handler level.
// Attributes are already included in the log record and will be written
// to the log file (JSON or text format). Use Logger.With() to add attributes to log messages.
func (h *FileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with the given group
// Note: FileHandler doesn't support groups at the handler level.
// Use Logger.With() to add structured context to log messages.
func (h *FileHandler) WithGroup(name string) slog.Handler {
	return h
}

// Close flushes the buffer (if buffered) and closes the underlying writer
// This implements io.Closer interface
func (h *FileHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var errs []error

	// Flush buffer if buffered
	if h.bufWriter != nil {
		if err := h.bufWriter.Flush(); err != nil {
			errs = append(errs, fmt.Errorf("failed to flush buffer: %w", err))
		}
	}

	// Close underlying writer if it implements io.Closer
	if h.underlying != nil {
		if err := h.underlying.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close underlying writer: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Flush flushes the buffer if buffered
// This is useful when you want to ensure logs are written immediately
func (h *FileHandler) Flush() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.bufWriter != nil {
		return h.bufWriter.Flush()
	}
	return nil
}

// ------------------- Private Methods -------------------
// 内部方法

// setEnabledLevels sets which log levels are enabled
func (h *FileHandler) setEnabledLevels(minLevel slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for level := slog.LevelDebug; level <= slog.LevelError; level++ {
		h.enabled[level] = level >= minLevel
	}
}

// handleJSON handles a log record in JSON format
func (h *FileHandler) handleJSON(record slog.Record) error {
	logEntry := map[string]interface{}{
		"time":    record.Time.Format(time.RFC3339),
		"level":   record.Level.String(),
		"message": record.Message,
	}

	// Add attributes
	record.Attrs(func(a slog.Attr) bool {
		logEntry[a.Key] = a.Value.Any()
		return true
	})

	// Add source if requested
	if h.opts.AddSource && record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := fs.Next()
		logEntry["source"] = fmt.Sprintf("%s:%d", frame.File, frame.Line)
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(h.writer, "%s\n", jsonData)
	return err
}

// handleText handles a log record in text format
func (h *FileHandler) handleText(record slog.Record) error {
	timestamp := record.Time.Format("2006-01-02 15:04:05")
	level := record.Level.String()

	// Format: [timestamp LEVEL] message
	msg := fmt.Sprintf("[%s %s] %s", timestamp, level, record.Message)

	// Add attributes with better formatting
	var attrs []string
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
		return true
	})

	if len(attrs) > 0 {
		// Format attributes as key=value pairs separated by spaces
		msg = fmt.Sprintf("%s %s", msg, fmt.Sprint(attrs))
	}

	_, err := fmt.Fprintf(h.writer, "%s\n", msg)
	return err
}
