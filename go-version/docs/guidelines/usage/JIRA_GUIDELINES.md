# Jira 使用指南

> 本文档包含 qkflow 中所有 Jira 相关功能的使用说明，包括 Issue 读取工具和状态配置管理。

---

## 📋 目录

- [Issue 读取工具](#-issue-读取工具)
- [状态配置管理](#-状态配置管理)
- [相关命令](#-相关命令)

---

## 📖 Issue 读取工具

qkflow 提供了强大的 Jira Issue 读取和导出功能，特别针对 Cursor AI 进行了优化。

### 🎯 快速开始

#### 在 Cursor 中使用（推荐）

最简单的方式：

```bash
# 在 Cursor 终端中运行：
qkflow jira read NA-9245

# 然后在 Cursor 聊天中直接说：
"总结刚才读取的 Jira ticket 内容"
```

Cursor 会自动读取导出的文件并提供全面的总结！

### 📚 可用命令

#### 1. `show` - 快速终端查看

直接在终端显示 Issue 信息。

```bash
# 基本视图（仅元数据）
qkflow jira show NA-9245

# 完整视图（包含描述和评论）
qkflow jira show NA-9245 --full
```

**适用场景**：
- 需要快速查看 Issue
- 只需要文本内容
- 需要最快的响应速度

#### 2. `export` - 完整导出（包含文件）

将 Issue 导出到本地文件，可选包含图片。

```bash
# 仅导出文本
qkflow jira export NA-9245

# 导出包含所有图片和附件
qkflow jira export NA-9245 --with-images

# 导出到自定义目录
qkflow jira export NA-9245 -o ~/jira-exports/ --with-images
```

**输出结构**：
```
/tmp/qkflow/jira/NA-9245/
├── README.md           # 如何在 Cursor 中使用
├── content.md          # 主要内容（Markdown）
└── attachments/        # 下载的文件（如果使用 --with-images）
    ├── screenshot.png
    └── diagram.jpg
```

**适用场景**：
- 需要图片/附件
- 需要保留本地副本
- 需要格式化的 Markdown

#### 3. `read` - 智能模式 ⭐️ **推荐**

自动决定最佳展示方式。

```bash
# 自动模式（智能决策）
qkflow jira read NA-9245
```

**工作原理**：
- ✅ **有图片？** → 导出到文件（包含图片）
- ✅ **仅文本？** → 直接在终端显示
- ✅ 自动优化 Cursor 使用体验

**适用场景**：
- 与 Cursor AI 配合使用（最佳体验）
- 希望工具自动决定最佳格式
- 不确定 Issue 是否包含图片

#### 4. `clean` - 清理导出文件

删除导出的文件以释放磁盘空间。

```bash
# 清理指定 Issue
qkflow jira clean NA-9245

# 清理所有导出
qkflow jira clean --all

# 预览将要删除的内容（dry run）
qkflow jira clean --all --dry-run

# 强制删除（无需确认）
qkflow jira clean --all --force
```

### 🎨 Cursor 使用示例

#### 示例 1：简单文本分析

```
你在 Cursor 中： "通过 qkflow 读取 NA-9245 并总结"

Cursor 执行： qkflow jira read NA-9245
Cursor 响应： "这个 ticket (NA-9245) 是关于..."
```

#### 示例 2：包含图片

```
你在 Cursor 中： "用 qkflow 读取 NA-9245 的所有内容包括图片，分析架构设计"

Cursor 执行： qkflow jira export NA-9245 --with-images
Cursor 读取： content.md + attachments/ 中的所有图片
Cursor 响应： "根据架构图，这个系统包含..."
```

#### 示例 3：手动控制

```bash
# 步骤 1：导出（你运行）
qkflow jira export NA-9245 --with-images

# 步骤 2：告诉 Cursor 读取什么
"读取 /tmp/qkflow/jira/NA-9245/content.md 并分析架构图"

# 步骤 3：完成后清理
qkflow jira clean NA-9245
```

### 🔧 配置要求

确保已配置 Jira 凭据：

```bash
qkflow init
```

必需设置：
- `jira_service_address`: Jira 实例 URL（如 https://brain-ai.atlassian.net）
- `jira_api_token`: Jira API Token
- `email`: Jira 邮箱

### 💡 使用技巧

#### 针对 Cursor 用户

1. **默认使用 `read` 命令** - 专为 AI 消费优化
2. **提示要具体** - 告诉 Cursor 你想知道什么
3. **定期清理** - 使用 `clean --all` 释放空间

#### 命令对比

| 命令 | 速度 | 图片 | 输出 | 最佳场景 |
|------|------|------|------|----------|
| `show` | ⚡️ 最快 | ❌ | 终端 | 快速查看 |
| `show --full` | ⚡️ 快 | ❌ | 终端 | 完整文本 |
| `export` | 🐌 较慢 | ❌ | 文件 | 文本归档 |
| `export --with-images` | 🐌 最慢 | ✅ | 文件 | 完整归档 |
| `read` ⭐️ | ⚡️ 智能 | ✅ 智能 | 智能 | **Cursor AI** |

#### Cursor 提示模板

```bash
# 通用总结
"通过 qkflow 读取 <ISSUE-KEY> 并总结内容"

# 特定分析
"用 qkflow 读取 <ISSUE-KEY>，分析技术方案"

# 带上下文
"读取 <ISSUE-KEY>，对比我们当前的实现方式"

# 包含图片
"qkflow 读取 <ISSUE-KEY> 包括所有图片，分析架构设计"
```

### 📊 输出格式

#### 终端输出（show 命令）

```
╔═══════════════════════════════════════════════════════╗
║ 🎫 NA-9245: 实现用户认证功能                          ║
╚═══════════════════════════════════════════════════════╝

📋 类型:        Story
📊 状态:        In Progress
🏷️  优先级:      High
👤 负责人:      John Doe

🔗 在 Jira 中查看: https://brain-ai.atlassian.net/browse/NA-9245
```

#### Markdown 输出（export 命令）

```markdown
---
issue_key: NA-9245
title: 实现用户认证功能
type: Story
status: In Progress
priority: High
---

# NA-9245: 实现用户认证功能

## 📊 元数据
...

## 📝 描述
...

## 📎 附件 (3)
1. **screenshot.png** (245 KB)
   ![screenshot.png](./attachments/screenshot.png)
...
```

#### Cursor 优化输出（read 命令）

`read` 命令提供 Cursor 可识别的特殊输出标记：

```
✅ 已导出到: /tmp/qkflow/jira/NA-9245/

主文件: /tmp/qkflow/jira/NA-9245/content.md
图片: /tmp/qkflow/jira/NA-9245/attachments/ (3 个文件)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💡 CURSOR: 请读取以下文件：
1. /tmp/qkflow/jira/NA-9245/content.md
2. /tmp/qkflow/jira/NA-9245/attachments/ 中的所有图片
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 🚀 高级用法

#### 批量处理

```bash
# 导出多个 Issue
for issue in NA-9245 NA-9246 NA-9247; do
  qkflow jira export $issue --with-images
done

# 告诉 Cursor 分析所有
"读取 /tmp/qkflow/jira/ 中的所有 Jira 导出并总结"
```

#### 自定义工作流

```bash
# 为团队创建脚本
#!/bin/bash
ISSUE_KEY=$1
qkflow jira read "$ISSUE_KEY"
echo "准备就绪，Cursor 可以分析了！"
```

#### 与其他工具集成

```bash
# 导出并在 VS Code/Cursor 中打开
qkflow jira export NA-9245 --with-images
code /tmp/qkflow/jira/NA-9245/content.md
```

### 🐛 故障排除

#### "Failed to create Jira client"
- 检查配置：`cat ~/.config/quick-workflow/config.yaml`
- 验证 API Token 是否有效
- 确保 `jira_service_address` 正确

#### "Failed to get issue"
- 验证 Issue Key 是否正确（如 NA-9245）
- 检查是否有查看该 Issue 的权限
- 先在浏览器中访问该 Issue 试试

#### "Failed to download attachment"
- 检查网络连接
- 验证 API Token 是否有附件下载权限
- 某些文件可能受限

#### Cursor 无法读取文件
- 确保导出命令成功完成
- 检查输出中的文件路径
- 尝试手动附加文件：`@/tmp/qkflow/jira/NA-9245/content.md`

### 📝 注意事项

- 导出文件默认存储在 `/tmp/qkflow/jira/`，是临时文件
- 定期使用 `clean` 命令释放空间
- 图片仅在 `--with-images` 标志时下载
- `read` 命令专为 Cursor AI 集成设计

---

## ⚙️ 状态配置管理

Jira 状态配置允许你为每个项目设置 PR 创建和合并时的状态映射。

### 📋 配置存储位置

Jira 每个项目的状态配置会根据你的系统智能存储：

**macOS with iCloud Drive**（推荐）：
```
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```
配置会自动在你的所有 Mac 设备间同步 ☁️

**本地存储**（回退方案）：
```
~/.qkflow/jira-status.json
```

运行 `qkflow config` 查看实际的存储位置。

### 📝 配置文件结构

```json
{
  "mappings": {
    "PROJ-123": {
      "project_key": "PROJ",
      "pr_created_status": "In Progress",
      "pr_merged_status": "Done"
    },
    "TEAM-456": {
      "project_key": "TEAM",
      "pr_created_status": "进行中",
      "pr_merged_status": "已完成"
    }
  }
}
```

### 🎯 配置说明

#### 1. 基本配置（必需）

在配置文件中设置 Jira 基本信息：

```yaml
email: your.email@example.com
jira_api_token: your_jira_api_token
jira_service_address: https://your-domain.atlassian.net
github_token: ghp_your_github_token
```

运行 `qkflow init` 进行初始化配置。

**提示**：如果你使用 macOS 并启用了 iCloud Drive，所有配置会自动同步到你的其他 Mac 设备！

#### 2. 项目状态映射（按项目配置）

每个 Jira 项目需要配置两个状态：

- **PR Created Status**（PR 创建时的状态）：当创建 PR 时，Jira issue 会更新到这个状态
  - 通常是：`In Progress`、`进行中`、`开发中` 等

- **PR Merged Status**（PR 合并时的状态）：当 PR 合并后，Jira issue 会更新到这个状态
  - 通常是：`Done`、`已完成`、`Resolved` 等

### 🛠️ 如何配置

#### 方式 1：首次使用时自动配置（推荐）

当你第一次为某个项目创建 PR 时，系统会自动提示你配置状态映射：

```bash
# 创建 PR
qkflow pr create PROJ-123

# 如果是首次使用该项目，会自动弹出交互式配置：
# 1. 从 Jira 获取该项目所有可用的状态
# 2. 让你选择 "PR Created" 状态（如：In Progress）
# 3. 让你选择 "PR Merged" 状态（如：Done）
# 4. 自动保存配置到 ~/.qkflow/jira-status.json
```

#### 方式 2：手动设置/更新项目配置

```bash
# 为指定项目设置状态映射
qkflow jira setup PROJ

# 系统会：
# 1. 连接到 Jira 获取该项目的所有可用状态
# 2. 显示交互式选择界面
# 3. 保存你的选择
```

#### 方式 3：查看已配置的项目

```bash
# 列出所有已配置的项目状态映射
qkflow jira list

# 输出示例：
# 📋 Jira Status Mappings:
#
# Project: PROJ
#   PR Created → In Progress
#   PR Merged  → Done
#
# Project: TEAM
#   PR Created → 进行中
#   PR Merged  → 已完成
```

#### 方式 4：删除项目配置

```bash
# 删除指定项目的状态映射
qkflow jira delete PROJ

# 会要求确认后删除
```

#### 方式 5：手动编辑配置文件

你也可以直接编辑配置文件：

```bash
# 编辑配置
vim ~/.qkflow/jira-status.json

# 或
code ~/.qkflow/jira-status.json
```

### 🔄 工作流程

#### 创建 PR 时（`qkflow pr create`）

1. 检查项目是否已有状态映射
2. 如果没有，自动触发配置流程
3. 创建 PR
4. 将 Jira issue 更新为 `PR Created Status`（如：In Progress）
5. 在 Jira issue 中添加 PR 链接

#### 合并 PR 时（`qkflow pr merge`）

1. 读取项目的状态映射
2. 合并 PR
3. 将 Jira issue 更新为 `PR Merged Status`（如：Done）
4. 在 Jira issue 中添加合并备注

### 📝 完整示例

```bash
# 1. 初始化基本配置
qkflow init

# 2. 查看当前配置
qkflow config

# 3. 为项目 PROJ 设置状态映射
qkflow jira setup PROJ
# 选择 PR Created: In Progress
# 选择 PR Merged: Done

# 4. 查看所有状态映射
qkflow jira list

# 5. 创建 PR（会自动使用配置的状态）
qkflow pr create PROJ-123

# 6. 合并 PR（会自动使用配置的状态）
qkflow pr merge 456
```

### 多项目配置示例

如果你在多个 Jira 项目工作：

```bash
# 为项目 A 配置
qkflow jira setup PROJA
# 选择: In Progress / Done

# 为项目 B 配置（可能用中文状态）
qkflow jira setup PROJB
# 选择: 进行中 / 已完成

# 为项目 C 配置（可能用自定义状态）
qkflow jira setup PROJC
# 选择: Development / Resolved

# 查看所有配置
qkflow jira list
```

### ⚙️ 技术实现

#### 状态获取

系统通过 Jira REST API 自动获取项目的所有可用状态：

```
GET /rest/api/2/project/{projectKey}/statuses
```

这确保你只能选择该项目实际支持的状态，避免配置错误。

#### 缓存机制

- 配置保存后会一直生效，除非手动更新或删除
- 每次操作时会自动读取对应项目的配置
- 如果配置被删除，下次使用时会重新提示配置

### 🔍 故障排除

#### 问题 1：找不到状态配置

```bash
# 检查配置文件是否存在
ls -la ~/.qkflow/jira-status.json

# 如果不存在，重新配置
qkflow jira setup YOUR_PROJECT_KEY
```

#### 问题 2：状态名称不匹配

如果 Jira 中的状态名称发生变化：

```bash
# 重新配置该项目
qkflow jira setup YOUR_PROJECT_KEY

# 或手动编辑配置文件
vim ~/.qkflow/jira-status.json
```

#### 问题 3：无法获取项目状态

确保：
1. Jira API Token 有效
2. 有该项目的访问权限
3. Jira Service Address 正确

```bash
# 检查基本配置
qkflow config

# 重新初始化配置
qkflow init
```

### 🎨 最佳实践

1. **首次使用时配置**：第一次为某个项目创建 PR 时就会提示配置，建议此时完成配置
2. **统一命名**：如果团队有多个项目，尽量使用统一的状态名称
3. **定期检查**：使用 `qkflow jira list` 定期检查配置是否正确
4. **备份配置**：可以备份 `~/.qkflow/jira-status.json` 文件

---

## 📚 相关命令

### Jira 相关命令

```bash
# Issue 读取
qkflow jira show <ISSUE-KEY>        # 终端查看
qkflow jira export <ISSUE-KEY>      # 导出文件
qkflow jira read <ISSUE-KEY>        # 智能模式（推荐）
qkflow jira clean <ISSUE-KEY>       # 清理导出文件

# 状态配置
qkflow jira list                    # 列出所有状态映射
qkflow jira setup [key]             # 设置/更新项目状态映射
qkflow jira delete [key]            # 删除项目状态映射
```

### 配置相关命令

```bash
qkflow init                         # 初始化配置
qkflow config                       # 查看当前配置
```

### 使用配置的命令

```bash
qkflow pr create [ticket]           # 创建 PR（会使用状态配置）
qkflow pr merge [number]            # 合并 PR（会使用状态配置）
```

---

## 🔗 相关文档

- [PR 使用指南](./PR_GUIDELINES.md)
- [入门指南](../../README.md#-快速开始)
- [自动更新指南](./AUTO_UPDATE_GUIDELINES.md)

---

**最后更新**：2025-12-05
