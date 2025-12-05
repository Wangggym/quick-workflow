package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/fatih/color"
)

// UIHandler is a custom handler for interactive CLI output
// It provides colored output with emoji for better UX
type UIHandler struct {
	writer  io.Writer
	opts    *HandlerOptions
	mu      sync.RWMutex
	enabled map[slog.Level]bool
}

// ------------------- Public Methods -------------------
// 外部访问

// NewUIHandler creates a new UIHandler
func NewUIHandler(writer io.Writer, opts *HandlerOptions) *UIHandler {
	if opts == nil {
		opts = &HandlerOptions{Level: LevelInfo}
	}

	h := &UIHandler{
		writer:  writer,
		opts:    opts,
		enabled: make(map[slog.Level]bool),
	}
	h.setEnabledLevels(opts.Level.ToSlogLevel())
	return h
}

// SetLevel changes the log level
func (h *UIHandler) SetLevel(level Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.opts.Level = level
	h.setEnabledLevels(level.ToSlogLevel())
}

// Enabled reports whether the handler handles records at the given level
func (h *UIHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.enabled[level]
}

// Handle handles a log record
func (h *UIHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.Enabled(ctx, record.Level) {
		return nil
	}

	emoji, colorFunc := h.getEmojiAndColor(record)

	// Format message
	msg := record.Message

	// Add attributes if any (excluding type attribute for UI)
	var attrs []string
	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "type" { // Skip type attribute for UI
			attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
		}
		return true
	})

	if len(attrs) > 0 {
		msg = fmt.Sprintf("%s (%s)", msg, fmt.Sprint(attrs))
	}

	// Write formatted message
	formatted := colorFunc("%s %s", emoji, msg)
	fmt.Fprintf(h.writer, "%s\n", formatted)

	return nil
}

// WithAttrs returns a new handler with the given attributes
// Note: UIHandler doesn't store attributes at the handler level.
// Attributes are already included in the log record and will be displayed
// in the formatted output. Use Logger.With() to add attributes to log messages.
func (h *UIHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with the given group
// Note: UIHandler doesn't support groups at the handler level.
// Use Logger.With() to add structured context to log messages.
func (h *UIHandler) WithGroup(name string) slog.Handler {
	return h
}

// ------------------- Private Methods -------------------
// 内部方法

// isSuccessMessage checks if the record represents a success message
func (h *UIHandler) isSuccessMessage(record slog.Record) bool {
	isSuccess := false
	record.Attrs(func(a slog.Attr) bool {
		if a.Key == "type" && a.Value.String() == "success" {
			isSuccess = true
			return false // Stop iteration
		}
		return true
	})
	return isSuccess
}

// getEmojiAndColor returns the appropriate emoji and color function for the log record
func (h *UIHandler) getEmojiAndColor(record slog.Record) (string, func(format string, a ...interface{}) string) {
	switch record.Level {
	case slog.LevelDebug:
		return "🔍", color.New(color.FgCyan).SprintfFunc()
	case slog.LevelInfo:
		if h.isSuccessMessage(record) {
			return "✅", color.New(color.FgGreen).SprintfFunc()
		}
		return "ℹ️", color.New(color.FgBlue).SprintfFunc()
	case slog.LevelWarn:
		return "⚠️", color.New(color.FgYellow).SprintfFunc()
	case slog.LevelError:
		return "❌", color.New(color.FgRed).SprintfFunc()
	default:
		return "ℹ️", color.New(color.FgBlue).SprintfFunc()
	}
}

// setEnabledLevels sets which log levels are enabled
func (h *UIHandler) setEnabledLevels(minLevel slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for level := slog.LevelDebug; level <= slog.LevelError; level++ {
		h.enabled[level] = level >= minLevel
	}
}
