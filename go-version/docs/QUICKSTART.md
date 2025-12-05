# 快速开始指南

5 分钟快速上手 qkflow！

> 💡 **提示**：这是快速入门指南。需要完整文档？请查看 [README.md](../README.md)。

## 📦 安装（30 秒）

选择你的平台，复制粘贴即可：

### macOS
```bash
# Apple Silicon (M1/M2/M3) - 推荐
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow && \
chmod +x qkflow && \
sudo mv qkflow /usr/local/bin/

# Intel Mac
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow && \
chmod +x qkflow && \
sudo mv qkflow /usr/local/bin/
```

### Linux
```bash
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow && \
chmod +x qkflow && \
sudo mv qkflow /usr/local/bin/
```

> 📖 **其他安装方式**：使用 Go 安装或从源码构建？查看 [README.md](../README.md#-安装)。

> ⚠️ **macOS 安全提示**：如果看到安全警告，运行 `xattr -d com.apple.quarantine qkflow` 移除隔离属性。

## ⚙️ 配置（2 分钟）

### 前置条件

在运行 `qkflow init` 之前，确保：

1. **GitHub CLI 已安装并认证**
   ```bash
   brew install gh  # macOS
   gh auth login
   ```

2. **获取 Jira API 令牌**（如果使用 Jira）
   - 访问：https://id.atlassian.com/manage-profile/security/api-tokens
   - 创建新令牌并复制

### 运行设置向导

```bash
qkflow init
```

按照提示输入：
- **邮箱**：你的工作邮箱
- **GitHub Token**：自动从 `gh` CLI 检测（无需手动输入）
- **Jira 地址**：`https://your-domain.atlassian.net`（可选）
- **Jira Token**：粘贴刚才获取的令牌（可选）
- **分支前缀**：可选（例如：`feature` 或你的用户名）

> 📖 **配置详情**：关于配置存储位置（iCloud 同步等），查看 [README.md](../README.md#-配置)。

## 🎯 创建第一个 PR（2 分钟）

### 步骤 1：进行更改

```bash
cd your-project
git checkout -b feature/test

# 进行一些更改
echo "# Test" >> README.md
git add README.md
```

### 步骤 2：创建 PR

```bash
qkflow pr create PROJ-123
```

按照提示操作：
1. **标题**：接受建议或输入自定义标题
2. **描述**：可选的简短描述
3. **变更类型**：选择适用的类型（feat、fix 等）
4. **Jira 状态**：选择新状态（可选）

**完成！** 你的 PR 已创建，Jira 已更新！🎉

## 🔄 合并 PR（1 分钟）

```bash
qkflow pr merge 123
```

按照提示操作：
1. **确认合并**：查看 PR 详情
2. **删除分支**：选择是否清理
3. **更新 Jira**：设置最终状态

**完成！** PR 已合并并清理！🎉

## 💡 专业技巧

### 不使用 Jira

```bash
# 跳过 Jira ticket（提示时按 Enter）
qkflow pr create
```

### 提示中的键盘快捷键

- **方向键**：导航选项
- **空格**：选择/取消选择（多选）
- **Enter**：确认选择
- **Ctrl+C**：取消操作

### 快速命令

```bash
# 显示配置
qkflow config

# 显示版本
qkflow version

# 获取帮助
qkflow --help
qkflow pr --help
```

## 🎨 工作流示例

### 典型开发流程

```bash
# 1. 创建功能分支
git checkout -b feature/awesome-feature

# 2. 进行更改并暂存
# ... 编写代码 ...
git add .

# 3. 创建 PR（自动提交、推送、创建 PR、更新 Jira）
qkflow pr create PROJ-456
# 按照提示输入标题、描述、变更类型等

# 4. 等待代码审查...

# 5. 合并 PR（自动清理分支、更新 Jira）
qkflow pr merge 789
```

### 快速更新现有 PR

```bash
# 修改代码后，快速提交和推送
qkflow update  # 使用 PR 标题作为提交信息
```

> 📖 **更多工作流示例**：查看 [README.md](../README.md#-工作流示例) 和 [PR 使用指南](guidelines/usage/PR_GUIDELINES.md)。

## 🐛 常见问题

### "Command not found: qkflow"

```bash
# 检查二进制文件是否存在
ls -l /usr/local/bin/qkflow

# 检查 PATH
echo $PATH | grep -q "/usr/local/bin" && echo "OK" || echo "Add to PATH"

# 如需要，添加到 PATH（添加到 ~/.zshrc）
export PATH="/usr/local/bin:$PATH"
```

### "Failed to create GitHub client"

```bash
# 确保 gh 已认证
gh auth status

# 如果未认证
gh auth login

# 重新运行 qkflow init
qkflow init
```

### "Failed to get Jira issue"

```bash
# 验证 Jira 凭据
curl -u "your.email@example.com:your_jira_token" \
  https://your-domain.atlassian.net/rest/api/2/myself

# 如果失败，获取新的 API token 并重新运行 qkflow init
```

## 📚 下一步

现在你已经掌握了基础操作！想要了解更多？

- 📖 **完整功能文档**：[README.md](../README.md) - 所有命令和功能的详细说明
- 📝 **PR 使用指南**：[PR_GUIDELINES.md](guidelines/usage/PR_GUIDELINES.md) - PR 功能的完整指南
- 🎫 **Jira 使用指南**：[JIRA_GUIDELINES.md](guidelines/usage/JIRA_GUIDELINES.md) - Jira 集成详细说明
- 🔄 **迁移指南**：[MIGRATION.md](MIGRATION.md) - 从 Shell 版本迁移
- 🐛 **遇到问题？**：[GitHub Issues](https://github.com/Wangggym/quick-workflow/issues)

## 🎉 准备就绪！

恭喜！你现在已经设置好 qkflow。享受流畅的工作流吧！

**常用命令备忘：**
```bash
qkflow pr create      # 创建 PR
qkflow pr merge       # 合并 PR
qkflow config         # 显示配置
qkflow --help         # 获取帮助
```

---

**祝编码愉快！🚀**
