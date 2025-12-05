# Release 操作指南

> 本文档包含 qkflow 的完整发布流程，从快速入门到详细操作说明。

---

## 📋 目录

- [快速入门](#-快速入门)
- [详细发布流程](#-详细发布流程)
- [版本管理](#-版本管理)
- [Release Notes 模板](#-release-notes-模板)
- [手动 Release](#-手动-release-如果自动失败)
- [检查 CI/CD 状态](#-检查-cicd-状态)
- [回滚 Release](#-回滚-release)
- [Release Checklist](#-release-checklist)
- [快速命令参考](#-快速命令参考)
- [最佳实践](#-最佳实践)

---

## 🚀 快速检查清单

发布前请确认：

- [ ] 所有更改已提交并推送：`git status` 和 `git push origin main`
- [ ] 运行 release 检查：`make release VERSION=vX.Y.Z`
- [ ] 创建并推送 tag：`git tag -a vX.Y.Z -m "Release vX.Y.Z" && git push origin vX.Y.Z`
- [ ] 查看 GitHub Actions：[Actions](https://github.com/Wangggym/quick-workflow/actions)
- [ ] 验证 Release：[Releases](https://github.com/Wangggym/quick-workflow/releases)

**版本号规则**：
- `v1.0.0` → `v1.0.1` - Bug 修复（PATCH）
- `v1.0.0` → `v1.1.0` - 新功能（MINOR）
- `v1.0.0` → `v2.0.0` - 破坏性更改（MAJOR）

---

## 📦 详细发布流程

### 1. 准备发布

确保所有更改已提交并推送到 main 分支：

```bash
# 确保在 main 分支
git checkout main
git pull origin main

# 查看当前状态
git status

# 如果有未提交的更改，先提交
git add .
git commit -m "chore: prepare for release v1.0.0"
git push origin main
```

### 2. 创建并推送 Tag

```bash
# 创建 tag (遵循语义化版本)
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送 tag 到远程仓库 (这会触发 GitHub Actions)
git push origin v1.0.0
```

**版本号规范** (语义化版本 SemVer):
- `v1.0.0` - 主版本.次版本.修订号
- `v1.0.0-beta.1` - 预发布版本
- `v1.0.0-rc.1` - 候选发布版本

### 3. 自动构建和发布

推送 tag 后，GitHub Actions 会自动：

1. ✅ 运行测试
2. ✅ 代码检查 (linting)
3. ✅ 构建多平台二进制文件：
   - `qkflow-darwin-amd64` (macOS Intel)
   - `qkflow-darwin-arm64` (macOS Apple Silicon)
   - `qkflow-linux-amd64` (Linux)
   - `qkflow-windows-amd64.exe` (Windows)
4. ✅ 创建 GitHub Release
5. ✅ 上传所有二进制文件到 Release

### 4. 验证 Release

1. 访问 GitHub Releases 页面：
   ```
   https://github.com/Wangggym/quick-workflow/releases
   ```

2. 确认以下内容：
   - ✅ Release 标题和描述正确
   - ✅ 所有平台的二进制文件都已上传
   - ✅ 文件大小合理 (通常 10-20MB)
   - ✅ 下载链接可用

3. 本地测试下载的二进制：
   ```bash
   # 下载并测试
   curl -L https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-darwin-arm64 -o qkflow
   chmod +x qkflow
   ./qkflow version
   ```

---

## 🔄 版本管理

### 版本号递增规则

- **MAJOR (主版本)**: 不兼容的 API 修改
  - 例如: `v1.0.0` → `v2.0.0`
  - 示例: 重大重构、删除旧功能

- **MINOR (次版本)**: 向后兼容的功能新增
  - 例如: `v1.0.0` → `v1.1.0`
  - 示例: 添加新命令、新功能

- **PATCH (修订号)**: 向后兼容的问题修正
  - 例如: `v1.0.0` → `v1.0.1`
  - 示例: Bug 修复、文档更新

### 预发布版本

测试新功能时使用预发布版本：

```bash
# Beta 版本
git tag -a v1.1.0-beta.1 -m "Beta release for testing"
git push origin v1.1.0-beta.1

# Release Candidate
git tag -a v1.1.0-rc.1 -m "Release candidate"
git push origin v1.1.0-rc.1
```

---

## 📝 Release Notes 模板

创建 tag 时，建议使用详细的 release notes：

```bash
git tag -a v1.0.0 -m "Release v1.0.0

## 🚀 New Features
- Add iCloud Drive sync support for configs
- Implement qkflow update command
- Add jira status mapping management

## 🐛 Bug Fixes
- Fix config path resolution on macOS
- Improve error handling in PR creation

## 📚 Documentation
- Add iCloud migration guide
- Update README with new features

## 🔧 Changes
- Rename binary from qkg to qkflow
- Unify config directory structure
"
```

或者在 GitHub Release 页面手动编辑。

---

## 🛠️ 手动 Release (如果自动失败)

如果 GitHub Actions 失败，可以手动构建和发布：

```bash
# 1. 构建所有平台
make build-all

# 2. 验证构建产物
ls -lh bin/

# 3. 在 GitHub 网页创建 Release
# - 访问: https://github.com/Wangggym/quick-workflow/releases/new
# - 选择 tag: v1.0.0
# - 填写标题和描述
# - 上传 bin/ 目录下的所有文件
# - 点击 "Publish release"
```

---

## 🔍 检查 CI/CD 状态

查看 GitHub Actions 执行状态：

1. 访问 Actions 页面：
   ```
   https://github.com/Wangggym/quick-workflow/actions
   ```

2. 查看构建日志：
   - 点击最近的 workflow run
   - 查看每个 job 的详细日志
   - 如有错误，根据日志修复

---

## 🚨 回滚 Release

如果发现严重问题需要回滚：

```bash
# 1. 在 GitHub 上删除 Release (标记为 draft 或删除)

# 2. 删除远程 tag
git push --delete origin v1.0.0

# 3. 删除本地 tag
git tag -d v1.0.0

# 4. 修复问题后重新发布
git tag -a v1.0.1 -m "Fix critical bug in v1.0.0"
git push origin v1.0.1
```

---

## 📊 Release Checklist

发布前检查清单：

```markdown
## Pre-release Checklist
- [ ] 所有测试通过 (`make test`)
- [ ] 代码检查通过 (`make check`)
- [ ] 文档已更新
- [ ] docs/CHANGELOG.md 已更新
- [ ] 版本号符合规范
- [ ] 已在本地测试构建 (`make build-all`)
- [ ] README 中的版本号已更新

## Release Checklist
- [ ] 创建 git tag
- [ ] 推送 tag 到远程
- [ ] GitHub Actions 构建成功
- [ ] 所有二进制文件已上传
- [ ] Release notes 已填写
- [ ] 下载链接可用

## Post-release Checklist
- [ ] 在社交媒体宣布（如需要）
- [ ] 更新文档网站（如需要）
- [ ] 通知团队成员
- [ ] 关闭相关 issues
```

---

## 🎯 快速命令参考

```bash
# 查看所有 tags
git tag -l

# 查看 tag 详情
git show v1.0.0

# 删除本地 tag
git tag -d v1.0.0

# 删除远程 tag
git push --delete origin v1.0.0

# 推送所有 tags
git push origin --tags

# 获取最新 tag
git describe --tags --abbrev=0

# 基于某个 commit 创建 tag
git tag -a v1.0.0 <commit-hash> -m "Release v1.0.0"
```

---

## 💡 最佳实践

1. **每次 Release 前运行 `make release VERSION=vX.Y.Z`**
2. **定期发布**: 建议每 2-4 周发布一个新版本
3. **语义化版本**: 严格遵循 SemVer 规范
4. **详细的 Release Notes**: 帮助用户了解变更
5. **测试充分**: 发布前在本地充分测试
6. **备份**: 保留旧版本的二进制文件
7. **通知用户**: 通过 GitHub、文档等方式通知用户
8. **发布后测试下载的二进制文件**

---

## 📞 遇到问题？

- 查看 [GitHub Actions 文档](https://docs.github.com/en/actions)
- 查看构建日志排查问题
- 联系维护者获取帮助

---

**记住**: 每次推送 tag 都会触发自动构建和发布，请谨慎操作！

**最后更新**：2025-12-05
