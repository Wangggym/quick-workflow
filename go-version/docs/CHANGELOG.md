# 变更日志

> 本文档记录 qkflow 的所有版本更新和功能变更。

---

## 📋 目录

- [v1.2.0](#v120---2024-11-18) - PR 审批和 Jira Reader 功能
- [v1.1.0](#v110---2024-11-01) - PR 编辑器功能
- [v1.0.0](#v100---2024-10-15) - 初始版本

---

## [v1.2.0] - 2024-11-18

### 🎉 新增

#### PR 审批功能

- **PR 审批命令**：`qkflow pr approve` 快速审批 Pull Request
  - 支持 PR 编号和 GitHub URL
  - 支持 `/files`、`/commits`、`/checks` 等 URL 路径
  - 默认使用 👍 作为审批评论
  - 支持自定义评论（`-c` 标志）
  - 支持审批后自动合并（`-m` 标志）
  - 自动检测当前分支的 PR

**使用示例**：
```bash
# 基本审批
qkflow pr approve 123

# 通过 URL 审批（跨仓库）
qkflow pr approve https://github.com/owner/repo/pull/456 -c "LGTM!" -m

# 自动检测当前分支
qkflow pr approve
```

#### Jira Issue 读取和导出功能

- **Jira 阅读器**：针对 Cursor AI 优化的 Issue 读取工具
  - `qkflow jira show` - 快速终端查看
  - `qkflow jira export` - 导出到文件（支持图片）
  - `qkflow jira read` - 智能模式（推荐用于 Cursor AI）
  - `qkflow jira clean` - 清理导出文件

**使用示例**：
```bash
# 智能读取（Cursor AI 优化）
qkflow jira read NA-9245

# 导出包含图片
qkflow jira export NA-9245 --with-images
```

### ✨ 改进

- **URL 支持**：`pr approve` 和 `pr merge` 命令现在支持 GitHub PR URL
  - 可以从任何地方使用，无需在 git 目录中
  - 支持跨仓库操作
  - 支持从浏览器直接复制粘贴 URL

- **交互体验**：改进 PR 创建时的编辑器选择界面
  - 从简单的 Yes/No 改为更直观的选择界面
  - 默认选项更清晰

### 🐛 修复

- 修复配置加载问题
- 改进错误提示信息

---

## [v1.1.0] - 2024-11-01

### 🎉 新增

#### PR 编辑器功能

- **基于 Web 的 PR 描述编辑器**
  - GitHub 风格界面，带深色主题
  - Markdown 编辑器，带实时预览
  - 支持拖放图片和视频
  - 支持从剪贴板粘贴图片
  - 自动上传到 GitHub 和 Jira
  - 支持的格式：PNG、JPG、GIF、WebP、SVG、MP4、MOV、WebM、AVI

**使用示例**：
```bash
# 创建 PR 时选择添加详细描述
qkflow pr create PROJ-123
# 选择 "是，继续" 打开 Web 编辑器
```

### ✨ 改进

- 优化 PR 创建流程
- 改进文件上传错误处理

---

## [v1.0.0] - 2024-10-15

### 🎉 初始版本

qkflow Go 版本的首次发布，包含以下核心功能：

#### 核心功能

- **PR 管理**
  - `qkflow pr create` - 创建 Pull Request
  - `qkflow pr merge` - 合并 Pull Request
  - `qkflow update` - 快速更新（使用 PR 标题作为提交信息）

- **Jira 集成**
  - 自动更新 Jira 状态
  - 添加 PR 链接到 Jira
  - Jira 状态配置管理

- **配置管理**
  - `qkflow init` - 交互式配置向导
  - `qkflow config` - 查看配置
  - iCloud Drive 同步（macOS）

- **监控守护进程**
  - `qkflow watch install` - 安装监控守护进程
  - 自动监控 PR 并在合并时更新 Jira

### ✨ 特性

- 📦 单一二进制文件，无需依赖
- ⚡ 快速启动（<100ms）
- 🌍 跨平台支持（macOS、Linux、Windows）
- 🎨 交互式 CLI，美观的提示
- ☁️ iCloud 同步配置（macOS）

---

## 📚 相关文档

- [PR 使用指南](./guidelines/usage/PR_GUIDELINES.md) - PR 功能完整使用指南
- [Jira 使用指南](./guidelines/usage/JIRA_GUIDELINES.md) - Jira 功能完整使用指南
- [快速开始指南](./QUICKSTART.md) - 5 分钟快速上手
- [迁移指南](./MIGRATION.md) - 从 Shell 版本迁移

---

## 🔗 版本导航

- [最新版本](#v120---2024-11-18)
- [v1.1.0](#v110---2024-11-01)
- [v1.0.0](#v100---2024-10-15)

---

**最后更新**：2025-12-05
