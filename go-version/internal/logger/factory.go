package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wangggym/quick-workflow/internal/utils"
)

// LoggerType 定义 logger 类型
type LoggerType string

const (
	// LoggerTypeUI 仅 UI 输出（stdout，带颜色和 emoji）
	LoggerTypeUI LoggerType = "ui"
	// LoggerTypeFile 仅文件输出（JSON 格式）
	LoggerTypeFile LoggerType = "file"
	// LoggerTypeMulti 文件 + UI 输出
	LoggerTypeMulti LoggerType = "multi"
)

// LoggerOptions 配置选项
// 统一的配置结构，用于所有 logger 创建场景
type LoggerOptions struct {
	Type     LoggerType `yaml:"type" mapstructure:"type"`               // Logger 类型
	Level    Level      `yaml:"level" mapstructure:"level"`             // 日志级别（可选，未设置时从环境变量 QKFLOW_LOG_LEVEL 或默认值 LevelInfo 读取）
	FileName string     `yaml:"file_name" mapstructure:"file_name"`     // 文件名（LoggerTypeFile 和 LoggerTypeMulti 需要，自动组合到 configDir）
	JSON     bool       `yaml:"json_format" mapstructure:"json_format"` // 是否使用 JSON 格式（默认 true）
}

// ------------------- Public Methods -------------------
// 外部访问

// DefaultLoggerOptions 返回默认配置
func DefaultLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		Type:  LoggerTypeUI,
		Level: LevelInfo,
		JSON:  true,
	}
}

// LoadOptions loads logger options from environment variables
func LoadOptions() *LoggerOptions {
	opts := DefaultLoggerOptions()

	// Load from environment variable
	if levelStr := os.Getenv("QKFLOW_LOG_LEVEL"); levelStr != "" {
		opts.Level = ParseLevel(levelStr)
	}

	// Load file name from environment (backward compatibility)
	// QKFLOW_LOG_FILE can be a full path or just a filename
	if filePath := os.Getenv("QKFLOW_LOG_FILE"); filePath != "" {
		opts.Type = LoggerTypeMulti
		// If it's a full path, extract just the filename
		// Otherwise use it as filename
		if filepath.IsAbs(filePath) {
			opts.FileName = filepath.Base(filePath)
		} else {
			opts.FileName = filePath
		}
	}

	return opts
}

// NewLogger 统一的工厂方法，根据类型和配置创建 logger
// 如果 opts.Level 未设置（零值），会自动从环境变量 QKFLOW_LOG_LEVEL 或默认值 LevelInfo 读取
// 优先级：显式传入 > 环境变量 > 默认值
func NewLogger(opts *LoggerOptions) (*Logger, error) {
	if opts == nil {
		opts = DefaultLoggerOptions()
	}

	// 如果 level 未设置（零值），从环境变量或默认值读取
	if opts.Level == "" {
		opts.Level = getDefaultLogLevel()
	}

	switch opts.Type {
	case LoggerTypeUI:
		return newUILogger(opts.Level), nil

	case LoggerTypeFile:
		filePath, err := resolveLogFilePath(opts)
		if err != nil {
			return nil, err
		}
		return newFileLogger(filePath, opts.Level, opts.JSON)

	case LoggerTypeMulti:
		filePath, err := resolveLogFilePath(opts)
		if err != nil {
			return nil, err
		}
		return newMultiLogger(filePath, opts.Level, opts.JSON)

	default:
		return nil, fmt.Errorf("unknown logger type: %s", opts.Type)
	}
}

// ------------------- Private Methods -------------------
// 内部方法

// getDefaultLogLevel 统一获取默认日志级别
// 优先级：环境变量 > 默认值
// 未来可以扩展支持从配置文件读取
func getDefaultLogLevel() Level {
	// 1. 优先从环境变量读取
	if levelStr := os.Getenv("QKFLOW_LOG_LEVEL"); levelStr != "" {
		return ParseLevel(levelStr)
	}

	// 2. 未来可以支持从配置文件读取
	// if cfg := config.Get(); cfg != nil && cfg.LogLevel != "" {
	//     return ParseLevel(cfg.LogLevel)
	// }

	// 3. 返回默认值
	return LevelInfo
}

// newUILogger 创建 UI logger
func newUILogger(level Level) *Logger {
	handler := NewUIHandler(os.Stdout, &HandlerOptions{
		Level:     level,
		AddSource: false,
	})
	logger := New(handler)
	logger.SetLevel(level)
	return logger
}

// ensureLogFile ensures the log file directory exists and opens the file
func ensureLogFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return file, nil
}

// createFileHandler creates a FileHandler with default options
func createFileHandler(file *os.File, level Level, json bool) *FileHandler {
	return NewFileHandler(file, &HandlerOptions{
		Level:      level,
		AddSource:  false,
		JSON:       json,
		Buffered:   true, // Enable buffering by default for better performance
		BufferSize: 4096,
	})
}

// newFileLogger 创建文件 logger
func newFileLogger(filePath string, level Level, json bool) (*Logger, error) {
	file, err := ensureLogFile(filePath)
	if err != nil {
		return nil, err
	}

	handler := createFileHandler(file, level, json)
	logger := New(handler)
	logger.mu.Lock()
	// FileHandler now implements io.Closer and will flush buffer and close file
	logger.handlers = append(logger.handlers, handler)
	logger.mu.Unlock()
	logger.SetLevel(level)

	return logger, nil
}

// newMultiLogger 创建多输出 logger（文件 + UI）
func newMultiLogger(filePath string, level Level, json bool) (*Logger, error) {
	file, err := ensureLogFile(filePath)
	if err != nil {
		return nil, err
	}

	fileHandler := createFileHandler(file, level, json)

	// 创建 UI handler
	uiHandler := NewUIHandler(os.Stdout, &HandlerOptions{
		Level:     level,
		AddSource: false,
	})

	// 创建 multi-handler
	multiHandler := NewMultiHandler(uiHandler, fileHandler)
	logger := New(multiHandler)
	logger.mu.Lock()
	// FileHandler now implements io.Closer and will flush buffer and close file
	logger.handlers = append(logger.handlers, fileHandler)
	logger.mu.Unlock()
	logger.SetLevel(level)

	return logger, nil
}

// resolveLogFilePath 解析日志文件路径
// 优先级：FileName > 环境变量 QKFLOW_LOG_FILE > 默认值 app.log
func resolveLogFilePath(opts *LoggerOptions) (string, error) {
	// 优先级1: FileName（代码指定）
	if opts.FileName != "" {
		configDir, err := utils.GetConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get config directory: %w", err)
		}
		return filepath.Join(configDir, opts.FileName), nil
	}

	// 优先级2: 环境变量（向后兼容）
	if filePath := os.Getenv("QKFLOW_LOG_FILE"); filePath != "" {
		// 如果是绝对路径，直接使用
		if filepath.IsAbs(filePath) {
			return filePath, nil
		}
		// 否则组合到 configDir
		configDir, err := utils.GetConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get config directory: %w", err)
		}
		return filepath.Join(configDir, filePath), nil
	}

	// 优先级3: 默认值
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(configDir, "app.log"), nil
}
