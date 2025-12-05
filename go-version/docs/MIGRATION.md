# 迁移指南

> 本文档包含从 Shell 版本迁移到 Go 版本以及 iCloud 配置迁移的完整指南。

---

## 📋 目录

- [Shell 到 Go 版本迁移](#-shell-到-go-版本迁移)
- [iCloud 配置迁移](#-icloud-配置迁移)

---

## 🔄 Shell 到 Go 版本迁移

### 🎯 为什么迁移？

| 方面 | Shell 版本 | Go 版本 | 改进 |
|------|-----------|---------|------|
| **安装** | Clone repo + 安装 4+ 依赖 | 下载 1 个二进制文件 | ✅ 简化 90% |
| **配置** | 手动在 `.zshrc` 中设置环境变量 | 交互式 `qkflow init` | ✅ 更容易 |
| **启动时间** | ~1-2 秒 | <100ms | ✅ 快 10-20 倍 |
| **跨平台** | 仅 macOS/Linux | macOS/Linux/Windows | ✅ 通用 |
| **更新** | `git pull` + 重新安装 | 下载新二进制文件 | ✅ 更简单 |

### 📋 前置条件

迁移前，请确保：
- ✅ 当前 Shell 版本正常工作
- ✅ 可以访问 Jira 和 GitHub 凭证
- ✅ 记录当前配置（特别是环境变量）

### 🔄 迁移步骤

#### 步骤 1: 备份当前配置

```bash
# 保存当前环境变量
cat ~/.zshrc | grep -E "(EMAIL|JIRA|GH_|OPENAI|DEEPSEEK)" > ~/qk-backup.txt
```

#### 步骤 2: 安装 Go 版本

```bash
# macOS (Apple Silicon)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# 验证安装
qkflow version
```

#### 步骤 3: 运行初始设置

```bash
qkflow init
```

这将提示你输入：
- **Email**: 使用 `EMAIL` 环境变量的值
- **GitHub Token**: 将从 `gh auth token` 自动检测
- **Jira Service Address**: 使用 `JIRA_SERVICE_ADDRESS` 的值
- **Jira API Token**: 使用 `JIRA_API_TOKEN` 的值
- **Branch Prefix**（可选）: 使用 `GH_BRANCH_PREFIX` 的值

#### 步骤 4: 测试新版本

```bash
# 测试 PR 创建
cd your-project
git checkout -b test-qkflow-migration
echo "test" > test.txt
git add test.txt
qkflow pr create

# 如果成功，你会看到：
# ✅ Branch created
# ✅ Changes committed
# ✅ Pushed to remote
# ✅ Pull request created: https://github.com/...
```

#### 步骤 5: 更新 Shell 别名（可选）

如果你在 Shell 版本中有自定义别名，请更新它们：

```bash
# 旧别名（Shell 版本）
alias prc='~/quick-workflow/pr-create.sh'
alias prm='~/quick-workflow/pr-merge.sh'

# 新别名（Go 版本）
alias prc='qkflow pr create'
alias prm='qkflow pr merge'
```

#### 步骤 6: 清理旧安装（可选）

一旦你验证了 Go 版本正常工作：

```bash
# 从 .zshrc 中删除旧环境变量
# 编辑 ~/.zshrc 并删除：
# export JIRA_API_TOKEN=...
# export JIRA_SERVICE_ADDRESS=...
# export GH_BRANCH_PREFIX=...

# 重新加载 shell
source ~/.zshrc

# 归档旧安装
mv ~/quick-workflow ~/quick-workflow-shell-backup
```

### 🔍 命令映射

| Shell 版本 | Go 版本 | 说明 |
|-----------|---------|------|
| `pr-create.sh` | `qkflow pr create` | 相同功能，更快 |
| `pr-merge.sh` | `qkflow pr merge` | 相同功能，更快 |
| N/A | `qkflow init` | 新功能：设置向导 |
| N/A | `qkflow config` | 新功能：显示配置 |

### 🐛 故障排除

#### 问题：找不到命令 "qkflow"

```bash
# 检查二进制文件是否在 PATH 中
which qkflow

# 如果未找到，确保 /usr/local/bin 在 PATH 中
export PATH="/usr/local/bin:$PATH"
```

#### 问题："Config not found"

```bash
# 运行设置向导
qkflow init
```

#### 问题："Failed to create GitHub client"

```bash
# 确保 gh CLI 已认证
gh auth status

# 如果未认证，登录
gh auth login

# 重新运行 qkflow init
qkflow init
```

#### 问题："Failed to get Jira issue"

1. 验证 Jira API token 是否正确
2. 检查 Jira service address 格式：`https://your-domain.atlassian.net`
3. 测试 Jira 凭证：
   ```bash
   curl -u your.email@example.com:your_jira_token \
     https://your-domain.atlassian.net/rest/api/2/myself
   ```

### 🎉 迁移检查清单

- [ ] 已备份当前配置
- [ ] 已安装 Go 版本
- [ ] 已运行 `qkflow init` 并配置
- [ ] 已测试 PR 创建
- [ ] 已验证 Jira 集成工作正常
- [ ] 已更新 shell 别名（如果有）
- [ ] 已删除旧环境变量
- [ ] 已归档旧 Shell 安装

---

## ☁️ iCloud 配置迁移

### 🌟 新特性

在 macOS 上，`qkflow` 会优先将配置存储到 iCloud Drive，实现多设备自动同步！

### 📍 存储位置

#### macOS with iCloud Drive（推荐）

```
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/
├── config.yaml          # 主配置文件
└── jira-status.json     # Jira 状态映射
```

配置会自动在你的所有 Mac 设备间同步 ☁️

#### 本地存储（回退方案）

```
~/.qkflow/
├── config.yaml          # 主配置文件
└── jira-status.json     # Jira 状态映射
```

运行 `qkflow config` 查看实际存储位置。

### 🔄 迁移到 iCloud

如果你之前使用的是本地配置，可以手动迁移到 iCloud：

```bash
# 1. 确保 iCloud Drive 已启用
# 打开"系统设置" → "Apple ID" → "iCloud" → 确保"iCloud Drive"已开启

# 2. 迁移配置文件
if [ -f ~/.qkflow/config.yaml ]; then
  cp ~/.qkflow/config.yaml \
     ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/config.yaml
fi

if [ -f ~/.qkflow/jira-status.json ]; then
  cp ~/.qkflow/jira-status.json \
     ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/jira-status.json
fi

# 3. 验证迁移
qkflow config
```

### 🎯 多设备使用

在新的 Mac 设备上：

1. **安装 qkflow**
   ```bash
   curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
   chmod +x qkflow
   sudo mv qkflow /usr/local/bin/
   ```

2. **等待 iCloud 同步**
   - 打开 Finder → iCloud Drive
   - 确保 `.qkflow` 文件夹已同步完成

3. **验证配置**
   ```bash
   qkflow config
   ```

配置会自动从 iCloud Drive 读取，无需重新配置！

### 🔒 安全说明

- iCloud Drive 存储是加密的
- 配置文件权限设置为 `0600`（仅用户可读写）
- Token 和密钥安全地存储在你的 iCloud 账户中
- 只有登录同一 Apple ID 的设备才能访问

### ⚠️ 注意事项

#### iCloud 同步延迟

iCloud Drive 同步可能需要几秒到几分钟，取决于网络连接速度和文件大小。

#### 离线工作

如果没有网络连接：
- 仍然可以读取和修改配置
- 更改会在网络恢复后自动同步

#### 回退到本地存储

如果你不想使用 iCloud Drive：

```bash
# 1. 禁用 iCloud Drive（系统设置）
# 2. qkflow 会自动回退到本地存储

# 或者手动移动配置回本地
cp ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/config.yaml \
   ~/.qkflow/config.yaml
```

### 🐛 故障排除

#### 问题 1: "No configuration found" 错误

```bash
# 检查 iCloud Drive 是否可用
ls -la ~/Library/Mobile\ Documents/com~apple~CloudDocs/

# 如果目录不存在，启用 iCloud Drive
# 系统设置 → Apple ID → iCloud → iCloud Drive

# 手动创建配置
qkflow init
```

#### 问题 2: 配置不同步

```bash
# 1. 检查 iCloud 同步状态
# 打开 Finder → iCloud Drive → 检查文件是否有云图标

# 2. 强制同步
# 右键点击文件 → "从 iCloud 下载"

# 3. 检查网络连接
ping icloud.com
```

---

## 📚 相关文档

- [快速开始指南](../README.md#-快速开始) - 快速开始
- [Jira 使用指南](./guidelines/usage/JIRA_GUIDELINES.md) - Jira 功能完整使用指南
- [PR 使用指南](./guidelines/usage/PR_GUIDELINES.md) - PR 功能完整使用指南

---

**最后更新**：2025-12-05
