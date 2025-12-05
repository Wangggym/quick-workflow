package logger_test

import (
	"github.com/Wangggym/quick-workflow/internal/logger"
)

// ExampleNewLogger_ui demonstrates how to create a UI logger for interactive CLI
func ExampleNewLogger_ui() {
	log, _ := logger.NewLogger(&logger.LoggerOptions{
		Type: logger.LoggerTypeUI,
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
	})

	log.Info("这是一条信息")
	log.Success("操作成功")
	log.Warning("这是一条警告")
	log.Error("这是一条错误")
	log.Debug("这是调试信息（默认不显示）")

	// Output:
	// ℹ️ 这是一条信息
	// ✅ 操作成功
	// ⚠️ 这是一条警告
	// ❌ 这是一条错误
}

// ExampleNewLogger_file demonstrates how to create a file logger
func ExampleNewLogger_file() {
	fileLog, err := logger.NewLogger(&logger.LoggerOptions{
		Type:     logger.LoggerTypeFile,
		FileName: "test.log",
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
		JSON: true,
	})
	if err != nil {
		panic(err)
	}
	defer fileLog.Close()

	fileLog.Info("这条日志会写入文件")
	fileLog.Success("成功日志")
	fileLog.Error("错误日志")
}

// ExampleNewLogger_multi demonstrates multi-output logging
func ExampleNewLogger_multi() {
	log, err := logger.NewLogger(&logger.LoggerOptions{
		Type:     logger.LoggerTypeMulti,
		FileName: "multi.log",
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
		JSON: true,
	})
	if err != nil {
		panic(err)
	}
	defer log.Close()

	log.Info("这条日志会同时输出到 stdout 和文件")
	log.Success("成功日志也会同时输出")
}

// ExampleNewLogger_watcher demonstrates watcher logger usage using factory method
func ExampleNewLogger_watcher() {
	// Create watcher logger (for daemon processes)
	log, err := logger.NewLogger(&logger.LoggerOptions{
		Type:     logger.LoggerTypeFile,
		FileName: "watch.log",
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
	})
	if err != nil {
		panic(err)
	}
	defer log.Close()

	log.Info("Daemon 启动")
	log.Info("处理了 %d 个任务", 10)
	log.Success("任务完成")
}

// ExampleLoadOptions demonstrates using configuration
func ExampleLoadOptions() {
	// Load options from environment (recommended)
	opts := logger.LoadOptions()
	opts.Level = logger.LevelDebug

	// Create logger from options
	log, err := logger.NewLogger(opts)
	if err != nil {
		panic(err)
	}

	log.Debug("调试日志")
	log.Info("信息日志")
}

// ExampleNewLogger_withoutLevel demonstrates that Level can be omitted
// and will be automatically loaded from environment variable or default value
func ExampleNewLogger_withoutLevel() {
	// Create logger without specifying Level
	// Level will be automatically loaded from:
	// 1. Environment variable QKFLOW_LOG_LEVEL (if set)
	// 2. Default value (LevelInfo)
	log, _ := logger.NewLogger(&logger.LoggerOptions{
		Type: logger.LoggerTypeUI,
		// Level is omitted - will use default or environment variable
	})

	log.Info("这条日志使用默认级别或环境变量设置的级别")
}
