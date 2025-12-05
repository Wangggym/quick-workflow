# 统一日志系统架构文档

基于 Go 标准库 `log/slog` 的统一日志系统，提供一致的 API 用于 UI 输出和文件日志。

## 目录

1. [架构概览](#架构概览)
2. [核心组件](#核心组件)
3. [设计模式](#设计模式)
4. [创建流程](#创建流程)
5. [实现流程](#实现流程)
6. [快速开始](#快速开始)
7. [使用指南](#使用指南)
8. [最佳实践](#最佳实践)

---

## 架构概览

### 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Logger (统一接口)                        │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  *slog.Logger (标准库)                                  │ │
│  │  handlers []io.Closer (资源管理)                        │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ 使用
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    slog.Handler 接口                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ UIHandler    │  │ FileHandler  │  │ MultiHandler │    │
│  │ (交互式输出)  │  │ (文件输出)    │  │ (多输出)      │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 设计原则

1. **统一接口**：所有日志操作通过 `Logger` 结构体，提供一致的 API
2. **标准库基础**：基于 Go 标准库 `log/slog`，保证兼容性和稳定性
3. **工厂模式**：通过 `NewLogger()` 统一创建入口，根据配置创建不同类型
4. **Handler 模式**：使用 `slog.Handler` 接口实现不同输出方式
5. **资源管理**：自动管理文件句柄等资源，支持优雅关闭

### 特性

- ✅ **统一接口**：同一套 API 用于所有日志输出
- ✅ **多种输出方式**：支持 stdout、stderr、文件、多输出
- ✅ **结构化日志**：支持 JSON 和文本格式
- ✅ **日志级别**：支持 Debug、Info、Warn、Error
- ✅ **UI 友好**：交互式输出带颜色和 emoji
- ✅ **文件日志**：后台服务日志写入文件，带时间戳
- ✅ **向后兼容**：提供便捷方法，便于逐步迁移

---

## 核心组件

### 1. Logger 结构体

**位置**：`logger.go`

```go
type Logger struct {
    *slog.Logger          // 标准库 Logger，提供核心日志功能
    mu       sync.RWMutex // 保护 handlers 的并发访问
    handlers []io.Closer  // 需要关闭的资源（如文件句柄）
}
```

**职责**：
- 包装 `slog.Logger`，提供统一的日志接口
- 管理需要关闭的资源（文件句柄等）
- 提供便捷方法（Info, Success, Warning, Error, Debug）
- 支持动态级别设置

### 2. Handler 类型

#### UIHandler
**位置**：`handler_ui.go`

**特点**：
- 输出到 `io.Writer`（通常是 `os.Stdout`）
- 带颜色和 emoji 的格式化输出
- 适合交互式 CLI 应用

**实现**：
- 实现 `slog.Handler` 接口
- 实现 `LevelSetter` 接口（支持动态级别）

#### FileHandler
**位置**：`handler_file.go`

**特点**：
- 输出到文件
- 支持 JSON 和文本两种格式
- 支持缓冲写入（默认启用，4KB 缓冲区）
- 实现 `io.Closer` 接口（自动刷新缓冲区并关闭文件）

**实现**：
- 实现 `slog.Handler` 接口
- 实现 `LevelSetter` 接口
- 实现 `io.Closer` 接口

#### MultiHandler
**位置**：`handler_multi.go`

**特点**：
- 将日志转发到多个 Handler
- 支持同时输出到 UI 和文件
- 统一管理多个 Handler 的级别

**实现**：
- 实现 `slog.Handler` 接口
- 实现 `LevelSetter` 接口（设置所有子 Handler 的级别）

### 3. 配置系统

#### LoggerOptions
**位置**：`factory.go`

```go
type LoggerOptions struct {
    Type     LoggerType // Logger 类型
    Level    Level      // 日志级别
    FileName string     // 文件名（LoggerTypeFile 和 LoggerTypeMulti 需要，自动组合到 configDir）
    JSON     bool       // 是否使用 JSON 格式
}
```

#### LoggerType 枚举

```go
LoggerTypeUI              // 仅 UI 输出（stdout）
LoggerTypeFile            // 仅文件输出
LoggerTypeMulti           // 文件 + UI 输出
```

---

## 设计模式

### 1. 工厂模式
- **统一入口**：`NewLogger(opts *LoggerOptions)`
- **类型枚举**：通过 `LoggerType` 区分不同场景
- **内部实现**：小写开头的私有方法（newUILogger, newFileLogger 等）

### 2. Handler 模式
- **标准接口**：实现 `slog.Handler` 接口
- **职责分离**：每个 Handler 负责一种输出方式
- **组合模式**：MultiHandler 组合多个 Handler

### 3. 接口设计
- **LevelSetter 接口**：支持动态设置日志级别
- **io.Closer 接口**：支持资源清理

---

## 创建流程

### 流程图

```
用户调用
  │
  ▼
NewLogger(opts *LoggerOptions)
  │
  ├─→ 参数验证和默认值设置
  │
  ├─→ 根据 opts.Type 选择创建路径
  │
  ├─→ LoggerTypeUI
  │     └─→ newUILogger(level)
  │           └─→ NewUIHandler(os.Stdout, opts)
  │           └─→ New(handler)
  │
  ├─→ LoggerTypeFile
  │     └─→ resolveLogFilePath(opts)  // 解析文件路径
  │     └─→ newFileLogger(filePath, level, json)
  │           ├─→ 创建目录（如不存在）
  │           ├─→ os.OpenFile(filePath, ...)
  │           ├─→ NewFileHandler(file, opts)
  │           ├─→ New(handler)
  │           └─→ 注册 handler 到 handlers 列表
  │
  └─→ LoggerTypeMulti
        └─→ newMultiLogger(filePath, level, json)
              ├─→ 解析文件路径（FileName > 环境变量 > 默认值）
              ├─→ 创建文件 handler
              ├─→ 创建 UI handler
              ├─→ NewMultiHandler(uiHandler, fileHandler)
              ├─→ New(multiHandler)
              └─→ 注册 fileHandler 到 handlers 列表
```

### 路径解析逻辑

对于 `LoggerTypeFile` 和 `LoggerTypeMulti`，文件路径的解析优先级：

1. **FileName**（代码指定）：自动组合为 `configDir/FileName`
2. **环境变量** `QKFLOW_LOG_FILE`（向后兼容）：
   - 如果是绝对路径，直接使用
   - 如果是相对路径，组合到 `configDir`
3. **默认值**：`configDir/app.log`

### 便捷方法

除了工厂方法，还提供了便捷的创建函数：

```go
// 快速创建 UI logger
func NewUILogger() *Logger

// 快速创建文件 logger（使用 FileName）
func NewLogger(opts *LoggerOptions) (*Logger, error)

// 快速创建多输出 logger（手动组合 handlers）
func NewMultiLogger(handlers ...slog.Handler) *Logger
```

---

## 实现流程

### 日志记录处理流程

```
用户调用 log.Info("message")
  │
  ▼
Logger.Info(format, args...)
  │
  ├─→ 格式化消息（如果有参数）
  │
  └─→ l.Logger.Info(msg)  // 调用 slog.Logger
        │
        ▼
      slog.Logger 内部处理
        │
        ├─→ 创建 slog.Record
        │     ├─→ 时间戳
        │     ├─→ 级别
        │     ├─→ 消息
        │     └─→ 属性（通过 With() 添加的）
        │
        └─→ 调用 Handler.Handle(ctx, record)
              │
              ├─→ Handler.Enabled(ctx, level) 检查
              │     └─→ 如果 false，直接返回
              │
              └─→ Handler.Handle(ctx, record) 处理
                    │
                    ├─→ UIHandler.Handle()
                    │     ├─→ 根据级别选择 emoji 和颜色
                    │     ├─→ 格式化消息
                    │     └─→ 写入 stdout
                    │
                    ├─→ FileHandler.Handle()
                    │     ├─→ 选择 JSON 或文本格式
                    │     ├─→ 序列化日志记录
                    │     └─→ 写入文件（可能缓冲）
                    │
                    └─→ MultiHandler.Handle()
                          ├─→ 遍历所有子 Handler
                          ├─→ 检查每个 Handler.Enabled()
                          └─→ 调用每个 Handler.Handle()
```

### Handler.Enabled() 实现

所有 Handler 都维护一个 `enabled` map，记录每个级别是否启用：

```go
enabled map[slog.Level]bool
```

**初始化**：
```go
func setEnabledLevels(minLevel slog.Level) {
    for level := slog.LevelDebug; level <= slog.LevelError; level++ {
        h.enabled[level] = level >= minLevel
    }
}
```

**检查**：
```go
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.enabled[level]
}
```

### 缓冲机制（FileHandler）

**启用条件**：`opts.Buffered == true`（默认启用）

**实现**：
```go
if opts.Buffered {
    h.bufWriter = bufio.NewWriterSize(writer, opts.BufferSize)
    h.writer = h.bufWriter  // 写入到缓冲区
} else {
    h.writer = writer  // 直接写入
}
```

**刷新时机**：
1. 缓冲区满时自动刷新
2. 调用 `FileHandler.Flush()` 手动刷新
3. 调用 `FileHandler.Close()` 时自动刷新

---

## 快速开始

### 基本使用

```go
import "github.com/Wangggym/quick-workflow/internal/logger"

// 创建 UI logger（用于交互式 CLI）
log := logger.NewUILogger()
log.Info("这是一条信息")
log.Success("操作成功")
log.Warning("这是一条警告")
log.Error("这是一条错误")

// 创建文件 logger（用于后台服务）
fileLog, err := logger.NewLogger(&logger.LoggerOptions{
    Type:     logger.LoggerTypeFile,
    FileName: "app.log",  // 自动组合到 configDir/app.log
    Level:    logger.LevelInfo,
})
if err != nil {
    // handle error
}
defer fileLog.Close()

fileLog.Info("这条日志会写入文件")
```

### 多输出 Logger

```go
// 同时输出到 stdout 和文件
uiHandler := logger.NewUIHandler(os.Stdout, &logger.HandlerOptions{
    Level: logger.LevelInfo,
})

file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
fileHandler := logger.NewFileHandler(file, &logger.HandlerOptions{
    Level: logger.LevelInfo,
    JSON:  true,
})

multiHandler := logger.NewMultiHandler(uiHandler, fileHandler)
log := logger.New(multiHandler)
```

### 结构化日志

```go
log := logger.NewUILogger()

// 添加结构化字段（推荐方式）
log.With("user", "alice", "action", "login").Info("用户登录")
log.With("pr", 123, "status", "merged").Success("PR 已合并")
```

**注意**：`UIHandler` 和 `FileHandler` 不支持 Handler 层面的 `WithAttrs/WithGroup`。
- 属性已经包含在日志记录中，会在输出时自动显示
- 使用 `Logger.With()` 是添加属性的标准方式
- `MultiHandler` 支持 `WithAttrs/WithGroup`，因为它需要将属性传递给子 Handler

### 动态设置日志级别

```go
log := logger.NewUILogger()
log.SetLevel(logger.LevelDebug) // 启用调试日志
log.Debug("这条日志只在 Debug 级别显示")
```

### 使用配置

```go
// 从环境变量加载配置
opts := logger.LoadOptions()
log, err := logger.NewLogger(opts)
if err != nil {
    // handle error
}
```

环境变量：
- `QKFLOW_LOG_LEVEL`: 日志级别 (debug/info/warn/error)
- `QKFLOW_LOG_FILE`: 日志文件名或路径（如果设置，会自动创建 Multi logger）

### Watch Daemon Logger

```go
// 创建 watch daemon 专用的 logger（仅文件输出）
log, err := logger.NewLogger(&logger.LoggerOptions{
    Type:     logger.LoggerTypeFile,
    FileName: "watch.log",  // 自动组合到 configDir/watch.log
    Level:    logger.LevelInfo,
})
if err != nil {
    // handle error
}
defer log.Close()

// 或者同时输出到 stdout 和文件（用于调试）
log, err := logger.NewLogger(&logger.LoggerOptions{
    Type:     logger.LoggerTypeMulti,
    FileName: "watch.log",
    Level:    logger.LevelInfo,
})
```

---

## 使用指南

### 日志级别

- `LevelDebug`: 调试信息
- `LevelInfo`: 一般信息
- `LevelWarn`: 警告信息
- `LevelError`: 错误信息

### 输出格式

#### UI Handler（交互式输出）

- 带颜色和 emoji
- 格式：`✅ 操作成功`、`❌ 错误信息`、`⚠️ 警告信息`、`ℹ️ 一般信息`

#### File Handler（文件输出）

支持两种格式：

1. **JSON 格式**（默认）：
```json
{"time":"2024-01-01T12:00:00Z","level":"INFO","message":"操作成功","type":"success"}
```

2. **文本格式**：
```
[2024-01-01 12:00:00 INFO] 操作成功 type=success
```

### 日志文件位置

#### watch.log（应用程序日志）
- **位置**: `~/.qkflow/watch.log` 或 `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/watch.log` (macOS iCloud)
- **用途**: Watch daemon 的应用程序日志
- **格式**: JSON 格式（结构化日志）
- **创建方式**: `logger.NewLogger(&logger.LoggerOptions{Type: logger.LoggerTypeFile, FileName: "watch.log"})`

#### watch.stdout.log 和 watch.stderr.log（系统日志）
- **位置**: `~/.qkflow/watch.stdout.log` 和 `~/.qkflow/watch.stderr.log`
- **用途**: launchd 的 stdout/stderr 重定向
- **格式**: 原始输出（非结构化）
- **创建方式**: launchd plist 配置

### 迁移指南

#### 从 UI 模块迁移

**旧代码**：
```go
import "github.com/Wangggym/quick-workflow/internal/ui"

ui.Info("信息")
ui.Success("成功")
ui.Error("错误")
```

**新代码**：
```go
import "github.com/Wangggym/quick-workflow/internal/logger"

log := logger.NewUILogger()
log.Info("信息")
log.Success("成功")
log.Error("错误")
```

#### 从 Watcher Logger 迁移

**旧代码**：
```go
import "github.com/Wangggym/quick-workflow/internal/watcher"

logger, _ := watcher.NewLogger()
logger.Info("信息")
```

**新代码**：
```go
import "github.com/Wangggym/quick-workflow/internal/logger"

log, _ := logger.NewLogger(&logger.LoggerOptions{
    Type:     logger.LoggerTypeFile,
    FileName: "watch.log",
    Level:    logger.LevelInfo,
})
log.Info("信息")
```

---

## 最佳实践

1. **交互式命令**：使用 `LoggerTypeUI` 创建 UI logger
2. **后台服务**：使用 `LoggerTypeFile` 或 `LoggerTypeMulti`，通过 `FileName` 指定文件名
3. **调试模式**：使用 `SetLevel(logger.LevelDebug)` 启用调试日志
4. **结构化日志**：使用 `With()` 方法添加上下文信息
5. **资源管理**：记得 `Close()` 文件 logger

### 使用场景

#### 场景 1: 交互式 CLI 命令
```go
log := logger.NewUILogger()
log.Info("开始处理...")
log.Success("操作完成")
```
- **Handler**: UIHandler
- **输出**: stdout（带颜色和 emoji）
- **格式**: 人类可读

#### 场景 2: 后台 Daemon（生产环境）
```go
log, _ := logger.NewLogger(&logger.LoggerOptions{
    Type:     logger.LoggerTypeFile,
    FileName: "watch.log",
    Level:    logger.LevelInfo,
})
defer log.Close()
```
- **Handler**: FileHandler
- **输出**: `{configDir}/watch.log`
- **格式**: JSON

#### 场景 3: 后台 Daemon（调试模式）
```go
log, _ := logger.NewLogger(&logger.LoggerOptions{
    Type:     logger.LoggerTypeMulti,
    FileName: "watch.log",
    Level:    logger.LevelDebug,
})
defer log.Close()
```
- **Handler**: MultiHandler (UIHandler + FileHandler)
- **输出**: stdout + 文件
- **格式**: UI（stdout）+ JSON（文件）

#### 场景 4: 通过环境变量配置
```go
opts := logger.LoadOptions()
log, _ := logger.NewLogger(opts)
```
- **行为**: 根据 `QKFLOW_LOG_FILE` 环境变量决定
- **有文件路径**: MultiHandler
- **无文件路径**: UIHandler

### 示例

#### 示例 1：CLI 命令

```go
func runCommand() {
    log := logger.NewUILogger()
    log.Info("开始处理...")

    // 执行操作
    if err := doSomething(); err != nil {
        log.Error("操作失败: %v", err)
        return
    }

    log.Success("操作完成")
}
```

#### 示例 2：后台服务

```go
func runDaemon() {
    log, err := logger.NewLogger(&logger.LoggerOptions{
        Type:     logger.LoggerTypeFile,
        FileName: "watch.log",
        Level:    logger.LevelInfo,
    })
    if err != nil {
        panic(err)
    }
    defer log.Close()

    log.Info("Daemon 启动")

    for {
        // 处理逻辑
        log.Info("处理了 %d 个任务", count)
    }
}
```

#### 示例 3：结构化日志

```go
log := logger.NewUILogger()

// 添加上下文
prLog := log.With("pr", 123, "repo", "myrepo")
prLog.Info("开始处理 PR")
prLog.Success("PR 已合并")
```

---

## 关键接口实现

### 1. slog.Handler 接口

所有 Handler 必须实现：

```go
type Handler interface {
    Enabled(ctx context.Context, level Level) bool
    Handle(ctx context.Context, record Record) error
    WithAttrs(attrs []Attr) Handler
    WithGroup(name string) Handler
}
```

**实现说明**：
- `UIHandler` 和 `FileHandler` 的 `WithAttrs/WithGroup` 直接返回自身（不存储属性）
- `MultiHandler` 的 `WithAttrs/WithGroup` 会创建新的子 Handler

### 2. LevelSetter 接口

```go
type LevelSetter interface {
    SetLevel(level Level)
}
```

**实现者**：
- `UIHandler`
- `FileHandler`
- `MultiHandler`（设置所有子 Handler 的级别）

**使用**：
```go
func (l *Logger) SetLevel(level Level) {
    if setter, ok := l.Handler().(LevelSetter); ok {
        setter.SetLevel(level)
    }
}
```

### 3. io.Closer 接口

```go
type Closer interface {
    Close() error
}
```

**实现者**：
- `FileHandler`：刷新缓冲区并关闭文件

**使用**：
- `Logger.Close()` 会关闭所有注册的 handlers

---

## 设计优势

1. **统一接口**：所有日志操作通过 `Logger`，API 一致
2. **灵活扩展**：通过 Handler 模式，易于添加新的输出方式
3. **资源管理**：自动管理文件句柄，支持优雅关闭
4. **性能优化**：文件写入支持缓冲，减少 I/O 操作
5. **向后兼容**：保留旧 API，便于逐步迁移

---

## 扩展点

1. **添加新的 Handler 类型**：实现 `slog.Handler` 接口
2. **添加新的 Logger 类型**：在 `NewLogger()` 中添加新的 case
3. **自定义格式化**：修改 Handler 的 `Handle()` 方法
4. **添加新的输出目标**：创建新的 Handler（如网络日志、数据库日志）
