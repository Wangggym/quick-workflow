package logger

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

// MultiHandler handles logs by forwarding to multiple handlers
type MultiHandler struct {
	handlers []slog.Handler
	mu       sync.RWMutex
	enabled  map[slog.Level]bool
}

// ------------------- Public Methods -------------------
// 外部访问

// NewMultiHandler creates a new MultiHandler
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	h := &MultiHandler{
		handlers: handlers,
		enabled:  make(map[slog.Level]bool),
	}
	// Enable all levels by default (each handler will filter)
	for level := slog.LevelDebug; level <= slog.LevelError; level++ {
		h.enabled[level] = true
	}
	return h
}

// SetLevel changes the log level for all handlers that support it
func (h *MultiHandler) SetLevel(level Level) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, handler := range h.handlers {
		if setter, ok := handler.(LevelSetter); ok {
			setter.SetLevel(level)
		}
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.enabled[level]
}

// Handle handles a log record by forwarding to all handlers
func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// WithAttrs returns a new handler with the given attributes
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

// WithGroup returns a new handler with the given group
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}
