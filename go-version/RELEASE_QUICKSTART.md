# Release 快速入门

## 🚀 快速发布

### 第一次 Release

```bash
# 1. 确保代码已推送
git push origin main

# 2. 运行 release 检查（会运行测试和构建）
make release VERSION=v1.0.0

# 3. 如果检查通过，创建并推送 tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. 等待 2-3 分钟，GitHub Actions 会自动构建和发布
```

### 日常 Release

```bash
# 补丁版本 (bug 修复)
make release VERSION=v1.0.1
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1

# 次版本 (新功能)
make release VERSION=v1.1.0
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 主版本 (不兼容更改)
make release VERSION=v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0
```

## 📋 Release Checklist

```bash
# ✅ 所有更改已提交
git status

# ✅ 运行 release 检查
make release VERSION=vX.Y.Z

# ✅ 创建 tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# ✅ 推送 tag
git push origin vX.Y.Z

# ✅ 查看 GitHub Actions
# https://github.com/Wangggym/quick-workflow/actions

# ✅ 验证 Release
# https://github.com/Wangggym/quick-workflow/releases
```

## 🎯 版本号规则

- `v1.0.0` → `v1.0.1` - Bug 修复
- `v1.0.0` → `v1.1.0` - 新功能
- `v1.0.0` → `v2.0.0` - 破坏性更改

## 🛠️ 常用命令

```bash
# 初始化依赖
make gen

# 运行测试
make test

# 构建本地版本
make build

# 安装到系统
make install

# 清理构建产物
make clean

# 查看帮助
make help
```

## 📦 发布后的产物

每次 Release 会生成以下文件：
- `qkflow-darwin-amd64` - macOS Intel
- `qkflow-darwin-arm64` - macOS Apple Silicon
- `qkflow-linux-amd64` - Linux
- `qkflow-windows-amd64.exe` - Windows

## 🔍 查看构建状态

- Actions: https://github.com/Wangggym/quick-workflow/actions
- Releases: https://github.com/Wangggym/quick-workflow/releases

## 💡 最佳实践

1. **每次 Release 前运行 `make release VERSION=vX.Y.Z`**
2. **使用语义化版本号**
3. **在 tag message 中添加简短的更新说明**
4. **发布后测试下载的二进制文件**

## 🚨 回滚操作

```bash
# 删除远程 tag
git push --delete origin v1.0.0

# 删除本地 tag
git tag -d v1.0.0

# 在 GitHub 上删除 Release
# 访问 Releases 页面手动删除
```

---

详细文档请查看 [RELEASE.md](RELEASE.md)

