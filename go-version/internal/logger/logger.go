package logger

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

// Logger is a unified logger interface that wraps slog.Logger
// It provides a consistent API for both UI and file logging
type Logger struct {
	*slog.Logger
	mu       sync.RWMutex
	handlers []io.Closer // For closing file handlers
}

// Level represents log levels
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// ------------------- Public Methods -------------------
// 外部访问

// String returns the string representation of the level
func (l Level) String() string {
	return string(l)
}

// ToSlogLevel converts our Level to slog.Level
func (l Level) ToSlogLevel() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// ParseLevel parses a string to Level
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "warn", "WARNING", "WARN":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// New creates a new Logger with the given handler
func New(handler slog.Handler) *Logger {
	return &Logger{
		Logger:   slog.New(handler),
		handlers: make([]io.Closer, 0),
	}
}

// LevelSetter is an interface for handlers that support dynamic level changes
type LevelSetter interface {
	SetLevel(level Level)
}

// SetLevel dynamically changes the log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Use interface instead of type assertions for better extensibility
	if setter, ok := l.Handler().(LevelSetter); ok {
		setter.SetLevel(level)
	}
}

// Close closes all file handlers
// Returns all errors if any occurred, using errors.Join to combine them
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var errs []error
	for _, handler := range l.handlers {
		if err := handler.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Info logs an info message with optional formatting
// Usage:
//
//	log.Info("simple message")
//	log.Info("formatted: %s", value)
func (l *Logger) Info(format string, args ...interface{}) {
	l.Logger.Info(formatMessage(format, args...))
}

// Warning logs a warning message with optional formatting
// Usage:
//
//	log.Warning("simple message")
//	log.Warning("formatted: %s", value)
func (l *Logger) Warning(format string, args ...interface{}) {
	l.Logger.Warn(formatMessage(format, args...))
}

// Error logs an error message with optional formatting
// Usage:
//
//	log.Error("simple message")
//	log.Error("formatted: %s", value)
func (l *Logger) Error(format string, args ...interface{}) {
	l.Logger.Error(formatMessage(format, args...))
}

// Success logs a success message with optional formatting (maps to Info level with success context)
// Usage:
//
//	log.Success("simple message")
//	log.Success("formatted: %s", value)
func (l *Logger) Success(format string, args ...interface{}) {
	l.Logger.Info(formatMessage(format, args...), "type", "success")
}

// Debug logs a debug message with optional formatting
// Usage:
//
//	log.Debug("simple message")
//	log.Debug("formatted: %s", value)
func (l *Logger) Debug(format string, args ...interface{}) {
	l.Logger.Debug(formatMessage(format, args...))
}

// With returns a new Logger with the given attributes
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		Logger:   l.Logger.With(args...),
		handlers: l.handlers,
	}
}

// ------------------- Private Methods -------------------
// 内部方法

// formatMessage formats a message with optional arguments
func formatMessage(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}
