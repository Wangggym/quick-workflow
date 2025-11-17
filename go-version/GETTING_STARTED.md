# 🚀 Getting Started with qkflow

## 📥 安装

### 方式 1: 下载预编译二进制 (推荐)

访问 [Releases 页面](https://github.com/Wangggym/quick-workflow/releases) 下载适合你系统的版本：

```bash
# macOS Apple Silicon (M1/M2/M3)
curl -L https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# macOS Intel
curl -L https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Linux
curl -L https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-linux-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-windows-amd64.exe -OutFile qkflow.exe
# 将 qkflow.exe 移动到 PATH 中的目录
```

### 方式 2: 从源码构建

```bash
git clone https://github.com/Wangggym/quick-workflow.git
cd quick-workflow/go-version
make gen      # 初始化依赖
make build    # 构建
make install  # 安装到 GOPATH/bin
```

## ⚙️ 初始化配置

首次使用需要配置 GitHub 和 Jira 信息：

```bash
qkflow init
```

按提示输入：
- GitHub Token (从 https://github.com/settings/tokens 获取)
- GitHub Owner (你的用户名或组织名)
- GitHub Repo (仓库名)
- Jira URL (如 https://your-company.atlassian.net)
- Jira Email
- Jira API Token (从 Jira 账户设置获取)

### 📱 iCloud 同步 (macOS)

配置会自动保存到 iCloud Drive（如果可用），实现多设备同步！

查看配置位置：
```bash
qkflow config
```

## 🎯 核心功能

### 1. 创建 PR

```bash
# 在当前分支创建 PR
qkflow pr create

# 交互式选择 AI 生成标题和描述
```

### 2. 合并 PR

```bash
# 合并当前分支的 PR
qkflow pr merge

# 自动更新 Jira 状态
```

### 3. 快速更新 (qkupdate)

```bash
# 自动使用 PR 标题作为 commit message
# 等同于: git add --all && git commit -m "PR Title" && git push
qkflow update
```

### 4. Jira 状态管理

```bash
# 查看已配置的项目
qkflow jira list

# 为项目配置状态映射
qkflow jira setup PROJECT-KEY

# 删除项目配置
qkflow jira delete PROJECT-KEY
```

## 📚 常用命令

```bash
# 查看版本
qkflow version

# 查看配置
qkflow config

# 查看帮助
qkflow --help
qkflow pr --help
qkflow jira --help
```

## 🔧 开发者命令

```bash
cd go-version

# 初始化依赖
make gen

# 格式化代码
make fix

# 运行测试
make test

# 构建
make build

# 安装到系统
make install

# 清理构建产物
make clean

# 查看所有命令
make help
```

## 🎓 工作流示例

### 完整的功能开发流程

```bash
# 1. 创建并切换到功能分支
git checkout -b feature/add-login

# 2. 开发功能...
# (编写代码)

# 3. 快速提交和推送
qkflow update

# 4. 创建 PR
qkflow pr create

# 5. Code Review...
# (等待审核通过)

# 6. 合并 PR（自动更新 Jira）
qkflow pr merge
```

### Bug 修复流程

```bash
# 1. 创建 bugfix 分支
git checkout -b bugfix/fix-login-error

# 2. 修复 bug...
# (修改代码)

# 3. 快速更新
qkflow update

# 4. 创建 PR
qkflow pr create

# 5. 合并
qkflow pr merge
```

## 🔐 环境变量（可选）

可以使用环境变量代替配置文件：

```bash
export GITHUB_TOKEN=your_token
export GITHUB_OWNER=your_username
export GITHUB_REPO=your_repo
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your@email.com
export JIRA_TOKEN=your_jira_token
export OPENAI_API_KEY=your_openai_key  # 可选
export DEEPSEEK_API_KEY=your_deepseek_key  # 可选
```

## 📖 更多文档

- [README.md](README.md) - 完整功能介绍
- [RELEASE.md](RELEASE.md) - Release 详细指南
- [RELEASE_QUICKSTART.md](RELEASE_QUICKSTART.md) - Release 快速入门
- [ICLOUD_MIGRATION.md](ICLOUD_MIGRATION.md) - iCloud 同步指南
- [JIRA_STATUS_CONFIG.md](JIRA_STATUS_CONFIG.md) - Jira 配置详解
- [CONTRIBUTING.md](CONTRIBUTING.md) - 贡献指南

## 💡 小贴士

1. **第一次使用**: 运行 `qkflow init` 配置
2. **查看配置**: 运行 `qkflow config` 查看存储位置
3. **快速更新**: 使用 `qkflow update` 代替繁琐的 git 命令
4. **Jira 集成**: 配置后 PR 操作自动更新 Jira 状态
5. **iCloud 同步**: macOS 用户配置自动同步到所有设备

## 🆘 获取帮助

- 遇到问题？查看 [Issues](https://github.com/Wangggym/quick-workflow/issues)
- 想要新功能？提交 [Feature Request](https://github.com/Wangggym/quick-workflow/issues/new)
- 贡献代码？查看 [CONTRIBUTING.md](CONTRIBUTING.md)

## 🎉 开始使用

```bash
# 安装
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# 初始化
qkflow init

# 开始使用！
qkflow pr create
```

Happy coding! 🚀

