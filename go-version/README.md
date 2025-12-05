# Quick Workflow (Go 版本)

> 一个现代化、极速的 CLI 工具，用于简化 GitHub 和 Jira 工作流

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Wangggym/quick-workflow?style=flat&logo=github)](https://github.com/Wangggym/quick-workflow/releases)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat&logo=github-actions)](https://github.com/Wangggym/quick-workflow/actions)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat)](https://github.com/Wangggym/quick-workflow)

## 🚀 快速开始

**新用户？** 查看 [📖 快速开始指南](docs/QUICKSTART.md) - 5 分钟快速上手！

**快速预览：**

```bash
# 1. 安装（macOS Apple Silicon）
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# 2. 初始化配置
qkflow init

# 3. 开始使用！
qkflow pr create PROJ-123
```

> 💡 **提示**：上面的命令只是快速预览。完整的安装步骤、配置说明和详细示例都在 [快速开始指南](docs/QUICKSTART.md) 中。

---

## ✨ 核心功能

- **PR 管理** - 创建、审批、合并 PR，支持 URL 操作
- **PR 编辑器** - 基于 Web 的编辑器，支持添加图片/视频描述
- **Jira 集成** - 自动更新 Jira 状态并添加 PR 链接
- **Jira 阅读器** - 读取和导出 Jira 问题（针对 Cursor AI 优化）
- **监控守护进程** - 自动监控 PR 并在合并时更新 Jira
- **iCloud 同步** - 在所有 Mac 设备间无缝同步配置（仅限 macOS）
- **自动更新** - 自动检查并安装更新（24 小时间隔）

📖 完整功能列表和使用说明见 [使用指南](docs/README.md)。

---

## 📦 安装

### 方式 1: 下载预编译二进制（推荐）

#### macOS
```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Intel
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/
```

> **⚠️ macOS 安全提示**：如果看到安全警告，运行 `xattr -d com.apple.quarantine qkflow` 移除隔离属性。

#### Linux
```bash
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/
```

### 方式 2: 使用 Go 安装

```bash
go install github.com/Wangggym/quick-workflow/cmd/qkflow@latest
```

### 方式 3: 从源码构建

```bash
git clone https://github.com/Wangggym/quick-workflow.git
cd quick-workflow/go-version
make build
sudo cp bin/qkflow /usr/local/bin/
```

---

## ⚙️ 配置

### 前置要求

- 已安装并配置 Git
- 已安装并认证 GitHub CLI (`gh`)：`gh auth login`
- Jira API 令牌：[在此获取](https://id.atlassian.com/manage-profile/security/api-tokens)

### 初始化配置

运行交互式设置：

```bash
qkflow init
```

这将提示你输入邮箱、GitHub 令牌、Jira 配置等。

**配置存储：**

在 macOS 上，配置会自动保存到 iCloud Drive，并在所有设备间同步：
- 📂 `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`

其他系统使用本地存储：
- 📂 `~/.qkflow/`

运行 `qkflow config` 查看实际存储位置。

📖 更多详情请参阅 [迁移指南](docs/MIGRATION.md)。

---

## 🎯 常用命令

### PR 操作

```bash
# 创建 PR
qkflow pr create PROJ-123

# 审批 PR（支持 URL）
qkflow pr approve 123
qkflow pr approve https://github.com/owner/repo/pull/123 -c "LGTM!" -m

# 合并 PR
qkflow pr merge 123

# 快速更新（使用 PR 标题作为提交信息）
qkflow update
```

### Jira 操作

```bash
# 读取 Jira Issue（Cursor AI 优化）
qkflow jira read NA-9245

# 查看 Issue
qkflow jira show NA-9245

# 导出 Issue（包含图片）
qkflow jira export NA-9245 --with-images

# 配置 Jira 状态映射
qkflow jira setup PROJECT-KEY
```

### 其他命令

```bash
qkflow config      # 显示配置
qkflow version     # 显示版本
qkflow update-cli  # 更新到最新版本
qkflow --help      # 获取帮助
```

📖 完整命令说明见 [PR 使用指南](docs/guidelines/usage/PR_GUIDELINES.md) 和 [Jira 使用指南](docs/guidelines/usage/JIRA_GUIDELINES.md)。

---

## 🎓 工作流示例

### 完整的功能开发流程

```bash
# 1. 创建功能分支
git checkout -b feature/add-login

# 2. 开发功能...
# (编写代码)

# 3. 快速提交和推送
qkflow update

# 4. 创建 PR
qkflow pr create PROJ-123

# 5. Code Review...
# (等待审核通过)

# 6. 合并 PR（自动更新 Jira）
qkflow pr merge
```

📖 更多工作流示例见 [快速开始指南](docs/QUICKSTART.md)。

---

## 💡 小贴士

1. **第一次使用**: 运行 `qkflow init` 配置
2. **查看配置**: 运行 `qkflow config` 查看存储位置
3. **快速更新**: 使用 `qkflow update` 代替繁琐的 git 命令
4. **Jira 集成**: 配置后 PR 操作自动更新 Jira 状态
5. **iCloud 同步**: macOS 用户配置自动同步到所有设备

---

## 🔧 故障排除

### 常见问题

**命令未找到**
```bash
which qkflow  # 检查是否在 PATH 中
export PATH="/usr/local/bin:$PATH"  # 如需要，添加到 ~/.zshrc
```

**GitHub 认证失败**
```bash
gh auth status  # 检查认证状态
gh auth login   # 如未认证，先登录
qkflow init     # 重新运行初始化
```

**Jira 连接失败**
```bash
# 验证 Jira 凭证
curl -u "your.email@example.com:your_jira_token" \
  https://your-domain.atlassian.net/rest/api/2/myself

# 如果失败，获取新的 API token 并重新运行 qkflow init
```

📖 **详细故障排除**：更多常见问题和解决方案见 [快速开始指南 - 常见问题](docs/QUICKSTART.md#-常见问题)。

---

## 🚧 从 Shell 版本迁移

详细迁移指南请参阅 [MIGRATION.md](docs/MIGRATION.md)。

**快速对比：**

| 功能 | Shell 版本 | Go 版本 |
|------|-----------|---------|
| 安装 | 克隆 + 依赖 | 单一二进制文件 |
| 配置 | `.zshrc` 环境变量 | `qkflow init` |
| 启动时间 | ~1-2 秒 | <100ms |
| 平台 | macOS/Linux | macOS/Linux/Windows |

---

## 🛠️ 开发

### 前置要求

- Go 1.21 或更高版本
- Make（可选但推荐）

### 构建

```bash
make build      # 为当前平台构建
make build-all  # 为所有平台构建
make test       # 运行测试
make lint       # 运行代码检查
```

📖 详细开发指南见 [开发规范](docs/guidelines/development/DEVELOPMENT_GUIDELINES.md)。

---

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

📖 详细贡献指南请参阅 [CONTRIBUTING.md](docs/guidelines/development/CONTRIBUTING.md)。

---

## 📚 文档

- 📖 [文档索引](docs/README.md) - 所有文档的索引
- 🚀 [快速开始](docs/QUICKSTART.md) - 5 分钟快速上手
- 📝 [PR 使用指南](docs/guidelines/usage/PR_GUIDELINES.md) - PR 功能完整说明
- 🎫 [Jira 使用指南](docs/guidelines/usage/JIRA_GUIDELINES.md) - Jira 功能完整说明
- 🔄 [迁移指南](docs/MIGRATION.md) - 从 Shell 版本迁移
- 🏗️ [架构文档](docs/architecture/ARCHITECTURE.md) - 项目架构说明

---

## 📄 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

---

## 📞 支持

- 🐛 [报告 Bug](https://github.com/Wangggym/quick-workflow/issues/new?labels=bug)
- 💡 [请求功能](https://github.com/Wangggym/quick-workflow/issues/new?labels=enhancement)
- 📖 [文档](docs/README.md)

---

由 [Wangggym](https://github.com/Wangggym) 用 ❤️ 制作
