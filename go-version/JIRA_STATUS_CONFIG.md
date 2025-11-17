# Jira 状态配置指南

## 📋 配置存储位置

Jira 每个项目的状态配置会根据你的系统智能存储：

**macOS with iCloud Drive** (推荐)：
```
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```
配置会自动在你的所有 Mac 设备间同步 ☁️

**本地存储** (回退方案)：
```
~/.qkflow/jira-status.json
```

运行 `qkflow config` 查看实际的存储位置。

配置文件结构：
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

## 🎯 配置说明

### 1. 基本配置（必需）

在 `~/.config/quick-workflow/config.yaml` 中配置 Jira 基本信息：

```yaml
email: your.email@example.com
jira_api_token: your_jira_api_token
jira_service_address: https://your-domain.atlassian.net
github_token: ghp_your_github_token
```

运行 `qkflow init` 进行初始化配置。

**提示**：如果你使用 macOS 并启用了 iCloud Drive，所有配置会自动同步到你的其他 Mac 设备！

### 2. 项目状态映射（按项目配置）

每个 Jira 项目需要配置两个状态：

- **PR Created Status**（PR 创建时的状态）：当创建 PR 时，Jira issue 会更新到这个状态
  - 通常是：`In Progress`、`进行中`、`开发中` 等
  
- **PR Merged Status**（PR 合并时的状态）：当 PR 合并后，Jira issue 会更新到这个状态
  - 通常是：`Done`、`已完成`、`Resolved` 等

## 🛠️ 如何配置

### 方式 1：首次使用时自动配置（推荐）

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

### 方式 2：手动设置/更新项目配置

```bash
# 为指定项目设置状态映射
qkflow jira setup PROJ

# 系统会：
# 1. 连接到 Jira 获取该项目的所有可用状态
# 2. 显示交互式选择界面
# 3. 保存你的选择
```

### 方式 3：查看已配置的项目

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

### 方式 4：删除项目配置

```bash
# 删除指定项目的状态映射
qkflow jira delete PROJ

# 会要求确认后删除
```

### 方式 5：手动编辑配置文件

你也可以直接编辑配置文件：

```bash
# 编辑配置
vim ~/.qkflow/jira-status.json

# 或
code ~/.qkflow/jira-status.json
```

## 🔄 工作流程

### 创建 PR 时（`qkflow pr create`）

1. 检查项目是否已有状态映射
2. 如果没有，自动触发配置流程
3. 创建 PR
4. 将 Jira issue 更新为 `PR Created Status`（如：In Progress）
5. 在 Jira issue 中添加 PR 链接

### 合并 PR 时（`qkflow pr merge`）

1. 读取项目的状态映射
2. 合并 PR
3. 将 Jira issue 更新为 `PR Merged Status`（如：Done）
4. 在 Jira issue 中添加合并备注

## 📝 示例

### 完整配置示例

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

## ⚙️ 技术实现

### 状态获取

系统通过 Jira REST API 自动获取项目的所有可用状态：

```
GET /rest/api/2/project/{projectKey}/statuses
```

这确保你只能选择该项目实际支持的状态，避免配置错误。

### 缓存机制

- 配置保存后会一直生效，除非手动更新或删除
- 每次操作时会自动读取对应项目的配置
- 如果配置被删除，下次使用时会重新提示配置

## 🔍 故障排除

### 问题 1：找不到状态配置

```bash
# 检查配置文件是否存在
ls -la ~/.qkflow/jira-status.json

# 如果不存在，重新配置
qkflow jira setup YOUR_PROJECT_KEY
```

### 问题 2：状态名称不匹配

如果 Jira 中的状态名称发生变化：

```bash
# 重新配置该项目
qkflow jira setup YOUR_PROJECT_KEY

# 或手动编辑配置文件
vim ~/.qkflow/jira-status.json
```

### 问题 3：无法获取项目状态

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

## 🎨 最佳实践

1. **首次使用时配置**：第一次为某个项目创建 PR 时就会提示配置，建议此时完成配置
2. **统一命名**：如果团队有多个项目，尽量使用统一的状态名称
3. **定期检查**：使用 `qkflow jira list` 定期检查配置是否正确
4. **备份配置**：可以备份 `~/.qkflow/jira-status.json` 文件

## 📚 相关命令

```bash
# Jira 相关命令
qkflow jira list           # 列出所有状态映射
qkflow jira setup [key]    # 设置/更新项目状态映射
qkflow jira delete [key]   # 删除项目状态映射

# 配置相关命令
qkflow init               # 初始化配置
qkflow config             # 查看当前配置

# 使用配置的命令
qkflow pr create [ticket] # 创建 PR（会使用状态配置）
qkflow pr merge [number]  # 合并 PR（会使用状态配置）
```

## 🔗 相关文件

### 配置文件位置

**macOS with iCloud Drive**:
- 配置目录: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
  - 基本配置: `config.yaml`
  - 状态映射: `jira-status.json`

**本地存储**:
- 配置目录: `~/.qkflow/`
  - 基本配置: `config.yaml`
  - 状态映射: `jira-status.json`

### 源码
- `internal/utils/paths.go` - 路径管理和 iCloud 检测
- `internal/jira/status_cache.go` - 状态缓存管理
- `cmd/qkflow/commands/jira.go` - Jira 命令
- `cmd/qkflow/commands/pr_create.go` - PR 创建逻辑
- `cmd/qkflow/commands/pr_merge.go` - PR 合并逻辑

