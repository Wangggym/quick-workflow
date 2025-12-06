# qkflow 文档索引

## 📚 文档概览

本文档目录包含 qkflow CLI 工具的完整使用指南、架构文档和开发文档。

---

## 🚀 快速开始

### [README.md](../README.md#-快速开始)
**入门指南**

- 安装方法（预编译二进制、源码构建）
- 初始化配置
- 基本使用示例
- 常见问题解答

### [QUICKSTART.md](./QUICKSTART.md)
**快速开始指南**

- 快速上手指南
- 常用命令速查
- 典型工作流示例

---

## 📖 使用指南

> 所有使用指南位于 [`guidelines/`](./guidelines/) 目录下

### [README.md](./guidelines/README.md)
**指南目录说明**

- 目录作用和功能说明
- 用户指南和开发指南分类
- 文档列表和使用场景

### 用户使用指南 (`usage/`)

> 位于 [`guidelines/usage/`](./guidelines/usage/) 目录下

### [PR_GUIDELINES.md](./guidelines/usage/PR_GUIDELINES.md)
**PR 功能完整指南**

- PR 审批功能（支持 PR 编号、URL、自动检测）
- PR URL 支持（通过 URL 操作 PR）
- PR 编辑器功能（添加详细描述和截图）
- 自定义审批评论和批量审批

### [JIRA_GUIDELINES.md](./guidelines/usage/JIRA_GUIDELINES.md)
**Jira 使用指南**

- Jira Issue 读取和导出功能
- Jira 状态配置管理
- 与 Cursor AI 集成使用
- 批量处理功能

### [AUTO_UPDATE_GUIDELINES.md](./guidelines/usage/AUTO_UPDATE_GUIDELINES.md)
**自动更新功能指南**

- 自动更新机制说明
- 配置自动更新
- 手动更新方法
- 更新频率和安全性

---

## 🏗️ 架构文档

> 所有架构文档位于 [`architecture/`](./architecture/) 目录下

### [README.md](./architecture/README.md)
**架构目录说明**

- 目录作用和功能说明
- 文档内容概述
- 使用场景指南

### [ARCHITECTURE.md](./architecture/ARCHITECTURE.md)
**项目架构文档**

- 项目概述和特性
- 技术栈和依赖
- 项目结构详解
- 模块组织方式
- 设计原则和模式
- 代码组织规范

---

## 🔄 迁移文档

> 迁移文档位于 [`docs/`](./) 目录下

### [MIGRATION.md](./MIGRATION.md)
**迁移指南**

- Shell 到 Go 版本迁移
- iCloud 配置迁移
- 功能更新说明

---

## 📋 需求文档

> 需求文档位于 [`requirements/`](./requirements/) 目录下

### [README.md](./requirements/README.md)
**需求文档目录说明**

- 功能开发需求文档
- 代码优化需求文档
- 使用指南

### [TODO.md](./requirements/TODO.md)
**功能开发需求文档**

- 从 workflow.go 和 workflow.rs 可以添加的功能
- 按优先级分类的功能列表
- 实现建议和开发顺序

### [OPTIMIZATION.md](./requirements/OPTIMIZATION.md)
**代码优化需求文档**

- 架构优化需求
- 性能优化需求
- 用户体验优化需求
- 代码质量优化需求

---

## 📝 变更日志

> 变更日志位于 [`docs/`](./) 目录下

### [CHANGELOG.md](./CHANGELOG.md)
**功能变更日志**

- 所有功能的版本更新记录
- PR 审批功能变更
- Jira Reader 功能变更
- 新功能添加
- Bug 修复和改进

---

## 🚢 发布文档

> 发布文档位于 [`guidelines/development/`](./guidelines/development/) 目录下

### [RELEASE_GUIDELINES.md](./guidelines/development/RELEASE_GUIDELINES.md)
**发布操作指南**

- 快速入门
- 发布流程说明
- 版本号规范
- GitHub Actions 自动构建
- 发布检查清单
- 详细操作说明

---

## 📖 快速导航

### 新用户
- 第一次使用？ → [README.md](../README.md#-快速开始)
- 想快速上手？ → [QUICKSTART.md](./QUICKSTART.md)

### 功能使用
- 想使用 PR 功能？ → [PR_GUIDELINES.md](./guidelines/usage/PR_GUIDELINES.md)
- 想使用 Jira 功能？ → [JIRA_GUIDELINES.md](./guidelines/usage/JIRA_GUIDELINES.md)
- 想了解自动更新？ → [AUTO_UPDATE_GUIDELINES.md](./guidelines/usage/AUTO_UPDATE_GUIDELINES.md)

### 架构和开发
- 想了解项目架构？ → [ARCHITECTURE.md](./architecture/ARCHITECTURE.md)
- 想了解开发规范？ → [DEVELOPMENT_GUIDELINES.md](./guidelines/development/DEVELOPMENT_GUIDELINES.md)
- 想了解文档规范？ → [DOCUMENT_GUIDELINES.md](./guidelines/development/DOCUMENT_GUIDELINES.md)
- 想贡献代码？ → [CONTRIBUTING.md](./guidelines/development/CONTRIBUTING.md)

### 迁移和更新
- 需要迁移？ → [MIGRATION.md](./MIGRATION.md)
- 查看变更日志？ → [CHANGELOG.md](./CHANGELOG.md)

### 需求和优化
- 想了解功能需求？ → [TODO.md](./requirements/TODO.md)
- 想了解优化需求？ → [OPTIMIZATION.md](./requirements/OPTIMIZATION.md)

### 发布
- 想发布新版本？ → [RELEASE_GUIDELINES.md](./guidelines/development/RELEASE_GUIDELINES.md)

---

## 📝 文档结构说明

文档按照以下分类组织：

1. **guidelines/** - 所有指南文档
   - **usage/** - 用户使用指南和功能文档
   - **development/** - 开发规范和文档编写指南（包括贡献指南）
2. **architecture/** - 项目架构和结构文档
3. **requirements/** - 需求文档（`TODO.md`、`OPTIMIZATION.md`）
4. **docs/** - 迁移指南和变更日志（`MIGRATION.md`、`CHANGELOG.md`）

---

## 🔄 文档更新记录

- **2025-12-05** - 文档优化和重构，统一文档风格，精简内容
- **2025-11-18** - 添加 PR 审批和 Jira Reader 功能文档
- **2025-10-15** - 创建文档索引和目录结构

---

## 📋 开发文档

> 开发相关文档位于 [`guidelines/development/`](./guidelines/development/) 目录下

> 详见 [指南目录说明](./guidelines/README.md)

### [DOCUMENT_GUIDELINES.md](./guidelines/development/DOCUMENT_GUIDELINES.md)
**文档编写指南**

- 文档类型和模板说明
- 使用指南模板
- 架构文档模板
- 迁移文档模板
- 变更日志模板
- 发布文档模板
- 文档格式规范

### [DEVELOPMENT_GUIDELINES.md](./guidelines/development/DEVELOPMENT_GUIDELINES.md)
**开发规范文档**

- 代码风格规范（格式化、Lint、命名约定）
- 错误处理规范
- 文档规范
- 命名规范
- 模块组织规范
- Git 工作流
- 提交规范
- 测试规范
- 代码审查
- 依赖管理
- 开发工具

### [CONTRIBUTING.md](./guidelines/development/CONTRIBUTING.md)
**贡献指南**

- 如何开始贡献
- 开发环境设置
- 代码风格和提交规范
- 测试指南
- PR 审查流程
- 错误报告和功能请求
- 贡献领域

---

## 📚 相关资源

- [主 README](../README.md) - 项目主文档
- [GitHub Repository](https://github.com/Wangggym/quick-workflow) - 源代码仓库
- [Issues](https://github.com/Wangggym/quick-workflow/issues) - 问题反馈
