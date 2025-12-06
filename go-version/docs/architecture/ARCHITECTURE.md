# 项目架构文档

> 本文档包含 qkflow 项目的完整架构说明，包括项目概述、技术栈、项目结构和设计原则。

---

## 📋 项目概述

### 📊 项目摘要

**状态**：✅ 已完成并可使用
**语言**：Go 1.21+
**类型**：CLI 工具
**目的**：简化 GitHub 和 Jira 工作流

### 🎯 核心功能

#### 核心功能
- ✅ 自动分支管理的 PR 创建
- ✅ 带清理功能的 PR 合并
- ✅ Jira 集成（状态更新、评论、链接）
- ✅ GitHub API 集成（PR CRUD 操作）
- ✅ Git 操作（分支、提交、推送、合并）
- ✅ 交互式 CLI 和美观的提示
- ✅ 通过 `qk init` 进行配置管理

#### 用户体验
- ✅ 单二进制分发（无依赖）
- ✅ 跨平台支持（macOS、Linux、Windows）
- ✅ 彩色输出和进度指示器
- ✅ 清晰的错误消息
- ✅ 交互式用户输入提示
- ✅ 剪贴板集成（macOS）

#### 开发者体验
- ✅ 模块化架构
- ✅ 类型安全的代码
- ✅ 全面的错误处理
- ✅ 易于测试和扩展
- ✅ 文档完善的代码
- ✅ 支持 GitHub Actions 的 CI/CD

### 📈 与 Shell 版本对比

| 指标 | Shell 版本 | Go 版本 | 改进 |
|------|-----------|---------|------|
| **二进制大小** | N/A | ~15MB | 自包含 |
| **启动时间** | ~1.5s | <100ms | 15 倍更快 |
| **依赖项** | 4+ 工具 | 0 | 减少 100% |
| **安装** | 多步骤 | 一个命令 | 更简单 |
| **错误处理** | 基础 | 全面 | 更好 |
| **类型安全** | 无 | 完整 | 更安全 |
| **测试** | 有限 | 全面 | 更可靠 |
| **平台** | macOS/Linux | macOS/Linux/Windows | 更多平台 |
| **维护** | 手动 | 自动化 | 更容易更新 |

---

## 🔧 技术栈

### 核心依赖

- **cobra**：CLI 框架，用于命令结构
- **viper**：配置管理
- **survey**：交互式提示和用户输入
- **go-github**：官方 GitHub API 客户端
- **go-jira**：Jira API 客户端
- **oauth2**：OAuth2 认证
- **fatih/color**：终端颜色

### 构建工具

- **go 1.21+**：语言和工具链
- **make**：构建自动化
- **golangci-lint**：代码检查
- **GitHub Actions**：CI/CD

---

## 📁 项目结构

```
go-version/
│
├── 📝 配置文件
│   ├── go.mod                          # Go 模块定义
│   ├── go.sum                          # 依赖校验和
│   ├── Makefile                        # 构建自动化
│   ├── .gitignore                      # Git 忽略规则
│   └── .golangci.yml                   # Linter 配置
│
├── 🎯 主应用程序
│   └── cmd/
│       └── qkflow/
│           ├── main.go                 # 应用程序入口点
│           └── commands/
│               ├── root.go             # 根命令和应用设置
│               ├── init.go             # 设置向导（qkflow init）
│               ├── pr.go               # PR 命令组
│               ├── pr_create.go        # PR 创建逻辑
│               └── pr_merge.go         # PR 合并逻辑
│
├── 🔒 内部包
│   └── internal/
│       ├── github/
│       │   └── client.go               # GitHub API 客户端
│       │       ├── NewClient()
│       │       ├── CreatePullRequest()
│       │       ├── GetPullRequest()
│       │       ├── MergePullRequest()
│       │       └── ParseRepositoryFromURL()
│       │
│       ├── jira/
│       │   └── client.go               # Jira API 客户端
│       │       ├── NewClient()
│       │       ├── GetIssue()
│       │       ├── UpdateStatus()
│       │       ├── AddComment()
│       │       ├── AddPRLink()
│       │       └── GetProjectStatuses()
│       │
│       ├── git/
│       │   └── operations.go           # Git 命令包装器
│       │       ├── CheckStatus()
│       │       ├── GetCurrentBranch()
│       │       ├── CreateBranch()
│       │       ├── Commit()
│       │       ├── Push()
│       │       ├── DeleteBranch()
│       │       ├── DeleteRemoteBranch()
│       │       ├── GetRemoteURL()
│       │       └── SanitizeBranchName()
│       │
│       ├── config/
│       │   └── config.go               # 配置管理
│       │       ├── Load()
│       │       ├── Get()
│       │       ├── Save()
│       │       ├── Validate()
│       │       └── IsConfigured()
│       │
│       └── ui/
│           └── prompt.go               # 用户界面辅助函数
│               ├── Success()
│               ├── Error()
│               ├── Warning()
│               ├── Info()
│               ├── PromptInput()
│               ├── PromptPassword()
│               ├── PromptConfirm()
│               ├── PromptSelect()
│               └── PromptMultiSelect()
│
├── 📦 公共包
│   └── (无 - 所有包都是内部的)
│
├── 🛠️ 脚本
│   └── scripts/
│       ├── install.sh                  # 安装脚本
│       ├── test.sh                     # 测试运行脚本
│       └── release.sh                  # 发布自动化
│
├── 🤖 CI/CD
│   └── .github/
│       └── workflows/
│           └── build.yml               # GitHub Actions 工作流
│
├── 📚 文档
│   ├── README.md                       # 主文档
│   ├── README.md                        # 主文档（包含快速开始）
│   └── docs/                           # 详细文档
│       ├── MIGRATION.md                # 迁移指南（Shell → Go、iCloud、功能更新）
│       ├── CHANGELOG.md                # 变更日志
│       ├── guidelines/                 # 使用和开发指南
│       └── architecture/                # 架构文档
│
├── 📄 法律
│   └── LICENSE                         # MIT 许可证
│
└── 🔨 构建输出（gitignored）
    └── bin/                            # 编译的二进制文件
        ├── qkflow                      # 当前平台
        ├── qkflow-darwin-amd64         # macOS Intel
        ├── qkflow-darwin-arm64         # macOS Apple Silicon
        ├── qkflow-linux-amd64          # Linux
        └── qkflow-windows-amd64.exe   # Windows
```

### 📊 文件统计

| 类别 | 文件数 | 代码行数（估算） |
|------|--------|-----------------|
| Go 源代码 | 10+ | ~2,000+ |
| 文档 | 10+ | ~3,000+ |
| 脚本 | 3 | ~300 |
| 配置 | 5 | ~200 |
| **总计** | **28+** | **~5,500+** |

---

## 🏗️ 架构设计

### 设计原则

1. **模块化**：将关注点分离到不同的包中
2. **可测试性**：易于模拟和测试
3. **用户优先**：优先考虑用户体验
4. **类型安全**：利用 Go 的类型系统
5. **错误处理**：清晰、可操作的错误
6. **性能**：快速启动和执行

### 包结构说明

#### `cmd/qkflow/commands`
- CLI 命令定义
- 用户交互逻辑
- 命令编排

#### `internal/github`
- GitHub API 客户端包装器
- PR 操作（创建、获取、合并）
- 仓库解析

#### `internal/jira`
- Jira API 客户端包装器
- Issue 操作（获取、更新）
- 状态管理

#### `internal/git`
- Git 命令执行
- 分支管理
- 提交和推送操作

#### `internal/ui`
- 用户提示和输入
- 彩色输出
- 进度指示器

#### `internal/config`
- 配置加载和保存
- 环境变量支持
- 验证

### 🔗 包依赖关系

```
cmd/qkflow/commands
  ├─→ internal/github
  ├─→ internal/jira
  ├─→ internal/git
  ├─→ internal/ui
  └─→ internal/config

internal/github
  └─→ internal/config

internal/jira
  └─→ internal/config

internal/config
  └─→ internal/utils

internal/git
  └─→ (无内部依赖)

internal/ui
  └─→ (无内部依赖)
```

### 📖 关键文件说明

#### 入口点
- **`cmd/qkflow/main.go`**：应用程序入口点，调用命令执行

#### 命令
- **`commands/root.go`**：根命令设置、版本、配置显示
- **`commands/init.go`**：首次配置的交互式设置向导
- **`commands/pr.go`**：PR 命令组（create/merge 的父命令）
- **`commands/pr_create.go`**：完整的 PR 创建工作流
- **`commands/pr_merge.go`**：完整的 PR 合并工作流

#### 核心库
- **`internal/github/client.go`**：带类型接口的 GitHub API 包装器
- **`internal/jira/client.go`**：带状态管理的 Jira API 包装器
- **`internal/git/operations.go`**：Git 命令执行和分支管理
- **`internal/ui/prompt.go`**：用户交互和彩色输出

#### 基础设施
- **`internal/config/config.go`**：配置加载、保存、验证
- **`Makefile`**：构建命令（build、test、lint、install）
- **`.github/workflows/build.yml`**：多平台构建的 CI/CD 管道

---

## 🎨 设计模式

1. **工厂模式**：客户端创建（`NewClient()`）
2. **命令模式**：CLI 命令结构
3. **仓库模式**：API 客户端抽象数据访问
4. **外观模式**：简化复杂操作的接口
5. **策略模式**：不同的 PR 类型和工作流

---

## 🔍 代码组织原则

### 1. 关注点分离
- `cmd/` - CLI 接口和用户交互
- `internal/` - 业务逻辑、API 客户端和配置

### 2. 依赖方向
- 命令依赖内部包
- 内部包可以依赖其他内部包
- 无循环依赖

### 3. 可见性
- `internal/` - 模块私有（外部项目无法导入）
- `cmd/` - 应用程序入口点

### 4. 测试
- 每个包都有自己的测试
- 模拟外部依赖
- 表驱动测试模式

---

## 📦 外部依赖

```
核心框架:
├── github.com/spf13/cobra           # CLI 框架
├── github.com/spf13/viper            # 配置
└── github.com/AlecAivazis/survey/v2  # 交互式提示

API 客户端:
├── github.com/google/go-github/v57  # GitHub API
├── github.com/andygrunwald/go-jira  # Jira API
└── golang.org/x/oauth2              # OAuth2 认证

工具:
└── github.com/fatih/color           # 终端颜色
```

---

## 🚀 构建产物

运行 `make build-all` 后：

```
bin/
├── qkflow-darwin-amd64      # macOS Intel (12-15MB)
├── qkflow-darwin-arm64      # macOS M1/M2 (12-15MB)
├── qkflow-linux-amd64       # Linux x86_64 (12-15MB)
└── qkflow-windows-amd64.exe # Windows 64-bit (12-15MB)
```

---

## 📈 项目指标

- **总代码行数**：~5,500+
- **Go 文件**：10+
- **包**：6+
- **命令**：4+
- **函数**：~80+
- **结构体**：~15+
- **接口**：~5+

---

## 🧪 测试策略

### 单元测试
- 测试单个函数
- 模拟外部依赖
- 使用表驱动测试

### 集成测试
- 测试 API 客户端（使用模拟）
- 测试命令执行
- 测试配置管理

### 手动测试
- 测试完整工作流
- 测试错误场景
- 在不同平台上测试

---

## 📊 构建和发布流程

### 开发构建
```bash
make build          # 为当前平台构建
make test           # 运行测试
make lint           # 运行 linter
```

### 多平台构建
```bash
make build-all      # 为 macOS、Linux、Windows 构建
```

### 发布流程
```bash
./scripts/release.sh v1.0.0
# 创建 tag，触发 CI/CD
# GitHub Actions 构建并上传二进制文件
```

---

## 🎯 导航指南

### 对于用户
1. 开始 → `README.md`（概述）
2. 设置 → `README.md` 的[快速开始](#-快速开始)部分（5 分钟）
3. 迁移 → `docs/MIGRATION.md`（如果从 Shell 版本迁移）

### 对于开发者
1. 架构 → 本文档
2. 贡献 → `docs/guidelines/development/CONTRIBUTING.md`
3. 代码 → 从 `cmd/qkflow/main.go` 开始

### 对于构建
1. 依赖 → `go.mod`
2. 构建 → `Makefile`
3. CI/CD → `.github/workflows/build.yml`
4. 发布 → `scripts/release.sh`

---

## 🔮 未来增强

### 高优先级
- [ ] GitLab 支持
- [ ] Bitbucket 支持
- [ ] Draft PR 支持
- [ ] PR 模板
- [ ] 自定义工作流

### 中优先级
- [ ] 更好的 Windows 集成
- [ ] Shell 补全脚本
- [ ] PR 审查自动化
- [ ] 批量操作
- [ ] Webhooks 集成

### 低优先级
- [ ] GUI 版本
- [ ] VS Code 扩展
- [ ] 指标和分析
- [ ] 团队仪表板

---

## 📞 支持和社区

### 获取帮助
- 📖 首先阅读文档
- 🐛 通过 GitHub Issues 报告 Bug
- 💡 通过 GitHub Issues 请求功能
- 💬 在 GitHub Discussions 中提问

### 贡献
- Fork、分支、编码、测试、PR
- 遵循编码标准
- 为新功能编写测试
- 更新文档

---

## 📄 许可证

MIT 许可证 - 详见 LICENSE 文件

---

## 🙏 致谢

### 原始项目
- Shell 版本由 [Wangggym](https://github.com/Wangggym) 创建

### Go 版本
- 架构和实现：AI 辅助开发
- 测试和改进：社区贡献者

### 开源库
- cobra、viper、survey（CLI 框架）
- go-github、go-jira（API 客户端）
- 以及更多优秀的 Go 包

---

## 📈 项目状态

**当前版本**：1.0.0（初始发布）
**状态**：✅ 可用于生产
**稳定性**：稳定
**维护**：活跃

---

## 🎉 总结

qkflow 的 Go 版本代表了原始 Shell 工具的全面现代化。它在以下方面带来了显著改进：

- **可用性**：更简单的安装和设置
- **性能**：更快的启动和执行
- **可靠性**：类型安全、经过充分测试的代码
- **可维护性**：清晰的架构、良好的文档
- **可扩展性**：易于添加新功能

项目已准备好用于生产，并欢迎社区贡献！

---

**最后更新**：2025-12-05
**维护者**：Wangggym
**仓库**：https://github.com/Wangggym/quick-workflow
