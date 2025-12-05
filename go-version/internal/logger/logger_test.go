package logger_test

import (
	"testing"

	"github.com/Wangggym/quick-workflow/internal/logger"
)

// TestLogger is a test helper
func TestLogger(t *testing.T) {
	log, _ := logger.NewLogger(&logger.LoggerOptions{
		Type: logger.LoggerTypeUI,
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
	})
	log.Info("Test log message")
}
