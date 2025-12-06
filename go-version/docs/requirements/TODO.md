# 功能开发 TODO

> 本文档列出了从 `workflow.go` 和 `workflow.rs` 项目中可以添加到 `go-version` 项目的功能需求。

---

## 📋 目录

- [命令对比](#命令对比)
- [高优先级功能](#高优先级功能)
- [中优先级功能](#中优先级功能)
- [低优先级功能](#低优先级功能)
- [功能详细说明](#功能详细说明)

---

## 命令对比

### go-version 当前命令

| 命令 | 子命令 | 说明 |
|------|--------|------|
| `init` | - | 初始化配置 |
| `version` | - | 显示版本 |
| `pr` | `create`, `merge`, `approve` | PR 操作 |
| `config` | - | 显示配置 |
| `jira` | `show`, `export`, `read`, `clean`, `list`, `setup`, `delete` | Jira 操作 |
| `update` | - | 快速更新代码 |
| `update-cli` | - | 更新 CLI |
| `watch` | `check`, `start`, `stop`, `restart`, `status`, `install`, `uninstall`, `log`, `history`, `config`, `daemon` | 监控守护进程 |

### workflow.go / workflow.rs 命令（完整列表）

| 命令 | 子命令 | 说明 |
|------|--------|------|
| `check` | - | 环境检查 |
| `proxy` | `on`, `off`, `check` | 代理管理 |
| `config` | `setup`, `show`, `log-level`, `completion` | 配置管理 |
| `github` | `list`, `current`, `add`, `remove`, `switch`, `update` | GitHub 账号管理 |
| `branch` | `clean`, `ignore` | 分支管理 |
| `pr` | `create`, `merge`, `status`, `close`, `comment`, `approve`, `list`, `update`, `summarize`, `sync`, `pick`, `rebase` | PR 操作 |
| `log` | `download`, `find`, `search` | 日志操作 |
| `jira` | `info`, `attachments`, `clean` | Jira 操作 |
| `llm` | `show`, `setup` | LLM 配置管理 |
| `lifecycle` | `install`, `uninstall`, `update`, `version` | 生命周期管理 |

**说明**：`workflow.rs` 还额外支持 `llm language` 命令（设置摘要语言）。

---

## 高优先级功能

### 1. 分支管理 (branch) ⭐⭐⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- `branch clean` - 清理已合并的分支
- `branch ignore` - 管理分支忽略列表

**理由**：提高开发效率，自动清理无用分支

**命令示例**：
```bash
qkflow branch clean                    # 清理已合并的分支
qkflow branch clean --dry-run          # 预览将要删除的分支
qkflow branch ignore add <BRANCH_NAME> # 添加分支到忽略列表
qkflow branch ignore remove <BRANCH_NAME> # 从忽略列表移除分支
qkflow branch ignore list              # 列出当前仓库的忽略分支
```

**实现位置**：
- `workflow.go/internal/commands/branch/`
- `workflow.rs/src/commands/branch/`

**状态**：❌ 未实现

---

### 2. PR 操作增强 ⭐⭐⭐

**来源**：`workflow.go` / `workflow.rs`

#### 2.1 PR 状态查询 (status)

**功能**：显示 PR 的详细信息

**命令示例**：
```bash
qkflow pr status [PR_ID_OR_BRANCH]  # 显示 PR 状态信息
```

**状态**：❌ 未实现

#### 2.2 PR 列表 (list)

**功能**：列出仓库中的所有 PR

**命令示例**：
```bash
qkflow pr list                    # 列出所有 PR
qkflow pr list --state open       # 按状态过滤
qkflow pr list --limit 10        # 限制结果数量
```

**状态**：❌ 未实现

#### 2.3 PR 同步 (sync)

**功能**：将源分支同步到当前分支

**命令示例**：
```bash
qkflow pr sync <SOURCE_BRANCH>              # 将指定分支同步到当前分支（merge）
qkflow pr sync <SOURCE_BRANCH> --rebase     # 使用 rebase 同步
qkflow pr sync <SOURCE_BRANCH> --squash      # 使用 squash 合并
qkflow pr sync <SOURCE_BRANCH> --ff-only     # 只允许 fast-forward 合并
qkflow pr sync <SOURCE_BRANCH> --no-push     # 不推送到远程
```

**状态**：❌ 未实现

#### 2.4 PR 关闭 (close)

**功能**：关闭指定的 PR

**命令示例**：
```bash
qkflow pr close [PR_ID]  # 关闭 PR
```

**状态**：❌ 未实现

#### 2.5 PR 评论 (comment)

**功能**：在 PR 上添加评论

**命令示例**：
```bash
qkflow pr comment [PR_ID] --message "..."  # 评论 PR
```

**状态**：❌ 未实现

---

### 3. 日志操作 (log) ⭐⭐⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- `log download` - 下载日志文件
- `log find` - 查找请求 ID
- `log search` - 搜索关键词

**理由**：实用工具，提高调试效率

**命令示例**：
```bash
qkflow log download [PROJ-123]               # 下载日志文件
qkflow log find [PROJ-123] [REQUEST_ID]     # 查找请求 ID
qkflow log search [PROJ-123] [SEARCH_TERM]  # 搜索关键词
```

**实现位置**：
- `workflow.go/internal/commands/log/`
- `workflow.rs/src/commands/log/`

**状态**：❌ 未实现

---

### 4. 环境检查 (check) ⭐⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- 检查 Git 仓库状态
- 检查到 GitHub 的网络连接
- 提供环境健康检查报告

**命令示例**：
```bash
qkflow check  # 运行环境检查
```

**实现位置**：
- `workflow.go/internal/commands/check/check.go`
- `workflow.rs/src/commands/check/check.rs`

**状态**：❌ 未实现

---

## 中优先级功能

### 5. 代理管理 (proxy) ⭐⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- 从系统设置读取代理配置
- 自动设置环境变量
- 支持临时和持久化配置

**命令示例**：
```bash
qkflow proxy on          # 开启代理
qkflow proxy off         # 关闭代理
qkflow proxy check       # 检查代理状态
```

**实现位置**：
- `workflow.go/internal/commands/proxy/`
- `workflow.rs/src/commands/proxy/`

**状态**：❌ 未实现

---

### 6. 配置管理增强 ⭐⭐

**来源**：`workflow.go` / `workflow.rs`

#### 6.1 Shell Completion 管理

**功能**：
- 自动检测 shell 类型
- 生成 completion 脚本
- 自动配置 shell 配置文件

**命令示例**：
```bash
qkflow config completion generate      # 生成 completion 脚本
qkflow config completion check         # 检查 completion 状态
qkflow config completion remove       # 移除 completion 配置
```

**状态**：❌ 未实现

#### 6.2 日志级别管理

**功能**：动态设置日志级别

**命令示例**：
```bash
qkflow config log-level set            # 设置日志级别
qkflow config log-level check          # 检查日志级别
```

**状态**：❌ 未实现

---

### 7. LLM 配置管理 (llm) ⭐⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- 独立的 LLM 配置管理
- 支持多种 LLM 提供者

**命令示例**：
```bash
qkflow llm show   # 显示当前 LLM 配置
qkflow llm setup   # 设置 LLM 配置
```

**Rust 版本额外功能**：
```bash
qkflow llm language  # 设置摘要语言（Rust 版本独有）
```

**实现位置**：
- `workflow.go/internal/commands/llm/`
- `workflow.rs/src/commands/llm/`

**状态**：⚠️ 部分存在（在 `init` 和 `config` 中，无独立命令）

---

### 8. PR 操作增强（续） ⭐⭐

#### 8.1 PR Rebase (rebase)

**功能**：将当前分支 rebase 到目标分支

**命令示例**：
```bash
qkflow pr rebase <TARGET_BRANCH>             # Rebase 当前分支到目标分支
qkflow pr rebase <TARGET_BRANCH> --no-push   # 只 rebase 到本地，不推送
qkflow pr rebase <TARGET_BRANCH> --dry-run    # 预览模式
```

**状态**：❌ 未实现

#### 8.2 PR Pick (pick)

**功能**：从源分支 cherry-pick 提交到目标分支并创建新 PR

**命令示例**：
```bash
qkflow pr pick <FROM_BRANCH> <TO_BRANCH>              # Pick 提交并创建新 PR
qkflow pr pick <FROM_BRANCH> <TO_BRANCH> --dry-run   # 预览模式
```

**状态**：❌ 未实现

#### 8.3 PR 总结 (summarize)

**功能**：使用 LLM 总结 PR 内容

**命令示例**：
```bash
qkflow pr summarize [PR_ID]                 # 使用 LLM 总结 PR
qkflow pr summarize --language zh            # 指定总结语言（Rust 版本支持）
```

**状态**：❌ 未实现

---

### 9. Jira 操作增强 ⭐⭐

**来源**：`workflow.go` / `workflow.rs`

#### 9.1 Jira 附件下载 (attachments)

**功能**：下载 Jira ticket 的所有附件

**命令示例**：
```bash
qkflow jira attachments [PROJ-123]  # 下载所有附件
```

**状态**：❌ 未实现

#### 9.2 Jira 清理增强 (clean)

**功能**：清理 Jira 相关的日志和附件文件

**命令示例**：
```bash
qkflow jira clean                    # 交互式清理
qkflow jira clean PROJ-123          # 清理指定 JIRA ID 的日志目录
qkflow jira clean --all             # 清理整个日志基础目录
qkflow jira clean --dry-run PROJ-123 # 预览清理操作
qkflow jira clean --list PROJ-123    # 只列出将要删除的内容
```

**状态**：⚠️ 部分存在（有 `jira clean` 命令，但无 `--dry-run` 和 `--list` 选项）

---

## 低优先级功能

### 10. GitHub 账号管理 (github) ⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- 多账号管理（当前版本支持单账号，架构支持多账号扩展）
- 账号切换功能

**命令示例**：
```bash
qkflow github list       # 列出所有 GitHub 账号
qkflow github current    # 显示当前激活的账号
qkflow github add        # 添加新的 GitHub 账号
qkflow github remove     # 删除 GitHub 账号
qkflow github switch     # 切换 GitHub 账号
qkflow github update     # 更新 GitHub 账号信息
```

**状态**：❌ 未实现

**优先级说明**：如果不需要多账号支持，可暂缓实现

---

### 11. 生命周期管理增强 ⭐

**来源**：`workflow.go` / `workflow.rs`

**功能**：
- 安装管理（二进制文件和 shell completion 脚本）
- 卸载管理（删除二进制文件、completion 脚本、配置文件）

**命令示例**：
```bash
qkflow install              # 安装 Workflow CLI
qkflow install --binaries   # 只安装二进制文件
qkflow install --completions # 只安装 shell completion 脚本
qkflow uninstall            # 卸载 Workflow CLI
```

**状态**：⚠️ 部分存在
- ✅ 有 `update-cli` 命令（功能类似 `update`）
- ✅ 有 `version` 命令
- ❌ 无 `install` 命令（通过 Makefile 或脚本安装）
- ❌ 无 `uninstall` 命令

---

## 功能详细说明

### 实现建议

1. **复用现有代码**：
   - 可以直接参考 `workflow.go` 或 `workflow.rs` 中的实现
   - 需要适配 `go-version` 的配置格式和架构

2. **保持兼容性**：
   - 保持现有命令的兼容性
   - 新功能作为增强添加

3. **渐进式迁移**：
   - 先实现高优先级功能
   - 逐步完善其他功能

4. **测试覆盖**：
   - 为新功能添加测试
   - 确保不影响现有功能

### 开发顺序

1. **第一阶段**：高优先级功能
   - 分支管理 (branch)
   - PR 操作增强（status, list, sync, close, comment）
   - 日志操作 (log)
   - 环境检查 (check)

2. **第二阶段**：中优先级功能
   - 代理管理 (proxy)
   - 配置管理增强（completion, log-level）
   - LLM 配置管理 (llm)
   - PR 操作增强（rebase, pick, summarize）
   - Jira 操作增强

3. **第三阶段**：低优先级功能
   - GitHub 账号管理
   - 生命周期管理增强

### 技术考虑

1. **配置格式**：
   - `go-version` 使用 YAML 格式
   - `workflow.go` / `workflow.rs` 使用 TOML 格式
   - 需要适配配置格式

2. **架构差异**：
   - `go-version` 使用 `internal/config/` 管理配置
   - `workflow.go` / `workflow.rs` 使用 `internal/base/settings/` 管理配置
   - 需要适配架构差异

3. **依赖管理**：
   - 评估是否需要新的依赖
   - 确保依赖兼容性

---

## 参考文档

- [OPTIMIZATION.md](./OPTIMIZATION.md) - 代码优化需求文档（包含架构差异说明）
- [workflow.rs/README.md](../../workflow.rs/README.md) - workflow.rs 项目文档
- [workflow.go/README.md](../../workflow.go/README.md) - workflow.go 项目文档

---

**文档创建时间**：2024年
