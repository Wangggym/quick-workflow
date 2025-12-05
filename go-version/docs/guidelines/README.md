
# Guidelines 目录说明

## 📋 目录作用

`guidelines/` 目录包含 qkflow 项目的所有指南文档，分为用户使用指南和开发指南两大类。

## 📂 目录结构

```
guidelines/
├── usage/                    # 用户使用指南
│   ├── PR_GUIDELINES.md     # PR 功能使用指南
│   ├── JIRA_GUIDELINES.md   # Jira 功能使用指南
│   └── AUTO_UPDATE_GUIDELINES.md  # 自动更新功能指南
│
└── development/              # 开发指南
    ├── DEVELOPMENT_GUIDELINES.md  # 开发规范
    ├── DOCUMENT_GUIDELINES.md     # 文档编写指南
    └── RELEASE_GUIDELINES.md      # 发布操作指南
```

> 💡 **注意**：快速开始指南已移至 [`../QUICKSTART.md`](../QUICKSTART.md)

## 📖 文档分类

### 用户使用指南 (`usage/`)

面向最终用户的功能使用文档，帮助用户快速上手和掌握 qkflow 的各项功能。

#### 文档列表

1. **PR_GUIDELINES.md** - PR 功能完整指南
   - PR 审批功能
   - PR URL 支持
   - PR 编辑器功能
   - 工作流示例和最佳实践

2. **JIRA_GUIDELINES.md** - Jira 功能完整指南
   - Jira Issue 读取和导出
   - Jira 状态配置管理
   - 与 Cursor AI 集成使用

3. **AUTO_UPDATE_GUIDELINES.md** - 自动更新功能指南
   - 自动更新机制说明
   - 配置和手动更新方法

> 💡 **快速开始指南**：已移至 [`../QUICKSTART.md`](../QUICKSTART.md)

### 开发指南 (`development/`)

面向开发者的规范和指南文档，帮助开发者理解项目规范、编写文档和发布版本。

#### 文档列表

1. **DEVELOPMENT_GUIDELINES.md** - 开发规范
   - 代码风格规范（格式化、Lint、命名约定）
   - 错误处理规范
   - Git 工作流和提交规范
   - 测试和代码审查规范

2. **DOCUMENT_GUIDELINES.md** - 文档编写指南
   - 文档类型和模板说明
   - 各类文档的编写模板
   - 文档格式规范
   - 文档维护指南

3. **RELEASE_GUIDELINES.md** - 发布操作指南
   - 快速入门
   - 完整的发布流程
   - 版本管理规则
   - CI/CD 检查
   - 发布检查清单

## 🎯 使用场景

### 新用户
- 查看 [`../QUICKSTART.md`](../QUICKSTART.md) 快速上手
- 查看 `usage/PR_GUIDELINES.md` 了解 PR 功能
- 查看 `usage/JIRA_GUIDELINES.md` 了解 Jira 集成

### 开发者
- 查看 `development/DEVELOPMENT_GUIDELINES.md` 了解开发规范
- 查看 `development/DOCUMENT_GUIDELINES.md` 了解文档编写规范
- 查看 `development/RELEASE_GUIDELINES.md` 了解发布流程

### 维护者
- 参考 `development/DOCUMENT_GUIDELINES.md` 维护文档
- 参考 `development/RELEASE_GUIDELINES.md` 发布新版本
- 参考 `development/DEVELOPMENT_GUIDELINES.md` 进行代码审查

## 📝 文档命名规范

所有指南文档遵循 `XX_GUIDELINES.md` 命名规范：
- `PR_GUIDELINES.md` - PR 相关功能指南
- `JIRA_GUIDELINES.md` - Jira 相关功能指南
- `DEVELOPMENT_GUIDELINES.md` - 开发规范指南
- `DOCUMENT_GUIDELINES.md` - 文档编写指南
- `RELEASE_GUIDELINES.md` - 发布操作指南

## 🔗 相关文档

- [文档索引](../README.md) - 查看所有文档
- [架构文档](../architecture/README.md) - 项目架构说明
- [迁移指南](../MIGRATION.md) - 迁移相关文档

---

**最后更新**：2025-12-05
