# quick-workflow

> Streamlined GitHub and Jira workflow automation tools

A collection of CLI tools to automate your daily GitHub PR and Jira workflow, available in both **Go** and **Shell** implementations.

## 💡 Why use it

> Are you tired of the hassle of changing the status every time you create a PR, such as 'code review,' 'merged,' and dealing with repetitive tasks?

> Especially when you're busy and have opened multiple browser windows, waiting for pages to load just to complete these mundane tasks can be time-consuming, diverting our focus from more important matters.

**Highlighted Benefits:**

- **Efficiency**: Streamline the process by automating status updates and repetitive tasks.
- **Time-Saving**: Eliminate the need to open multiple browser tabs and wait for pages to load.
- **Enhanced Focus**: Allow users to concentrate on more crucial tasks.
- **Reduced Manual Work**: Minimize the need for manual status updates for each PR.
- **Standardized Workflow**: Ensures a consistent and error-free process.

## 🚀 Features

- **Automated PR Creation**: Create branch, commit, push, and open PR in one command
- **Jira Integration**: Automatically link PRs to Jira tickets and update status
- **Smart Branch Management**: Auto-detect default branch (main/master)
- **Interactive CLI**: User-friendly prompts with support for cancellation (Ctrl+C)
- **One-Command Merge**: Merge PR, delete branch, and update Jira status
- **Browser Integration**: Auto-open PR in browser and copy URL to clipboard

## 📦 Two Implementations

### Go Version (Recommended) 🆕

Modern, fast, and cross-platform implementation written in Go.

```bash
# Install
cd go-version
make install

# Usage
qkg pr create    # Create a PR
qkg pr merge     # Merge a PR
qkg init         # Initialize configuration
```

**Benefits:**
- ✅ Single binary, no dependencies
- ✅ Cross-platform (macOS, Linux, Windows)
- ✅ Fast execution
- ✅ Better error handling
- ✅ Type-safe configuration

👉 [Go Version Documentation](./go-version/README.md)

### Shell Version (Legacy)

Original implementation using shell scripts.

```bash
# Installation (one-time setup)
cd quick-workflow/shell-version
./install.sh

# Usage
qk pr create     # Create a PR
qk pr merge      # Merge a PR
```

👉 [Shell Version Documentation](./shell-version/README.md)

## 🎯 Quick Start

### Option 1: Go Version (Recommended)

1. **Prerequisites**:
   ```bash
   # Install Go
   brew install go
   
   # Install GitHub CLI
   brew install gh
   gh auth login
   
   # Install Jira CLI (optional)
   brew install jira-cli
   jira init
   ```

2. **Install qkg**:
   ```bash
   cd go-version
   make install
   ```

3. **Configure**:
   ```bash
   qkg init
   # Follow the prompts to set up GitHub token, Jira credentials, etc.
   ```

4. **Use**:
   ```bash
   # Stage your changes
   git add .
   
   # Create PR with interactive prompts
   qkg pr create
   
   # Or with Jira ticket
   qkg pr create PROJ-123
   ```

### Option 2: Shell Version

See [shell-version/README.md](./shell-version/README.md) for installation instructions.

## 📖 Core Commands

### PR Create Workflow

```bash
qkg pr create [jira-ticket]
```

What it does:
1. ✅ Prompts for Jira ticket (optional)
2. ✅ Fetches Jira issue details if ticket provided
3. ✅ Prompts for PR title and description
4. ✅ Lets you select change types (feat, fix, docs, etc.)
5. ✅ Creates a new branch
6. ✅ Stages and commits changes
7. ✅ Pushes to remote
8. ✅ Creates GitHub PR
9. ✅ Links PR to Jira ticket
10. ✅ Updates Jira status
11. ✅ Copies PR URL to clipboard
12. ✅ Opens PR in browser

### PR Merge Workflow

```bash
qkg pr merge
```

What it does:
1. ✅ Lists your open PRs
2. ✅ Lets you select which PR to merge
3. ✅ Merges the PR
4. ✅ Deletes the remote branch
5. ✅ Deletes the local branch
6. ✅ Updates Jira status
7. ✅ Switches back to default branch

## 🔧 Configuration

### Go Version

Configuration is stored in `~/.config/qk/config.yaml`:

```yaml
github:
  token: ghp_xxx
jira:
  email: your.email@example.com
  api_token: xxx
  server: https://your-domain.atlassian.net
branch_prefix: "ym"  # Optional: prefix for branch names
```

You can also set environment variables:
- `GITHUB_TOKEN` or `GH_TOKEN`
- `JIRA_API_TOKEN`
- `JIRA_SERVICE_ADDRESS`
- `GH_BRANCH_PREFIX`

### Shell Version

Uses environment variables in `.zshrc` or `.bashrc`. See [shell-version/README.md](./shell-version/README.md).

## 🛠️ Development

### Go Version

```bash
cd go-version

# Initialize dependencies
make gen

# Format and fix code
make fix

# Run linters
make check

# Run tests
make test

# Development mode (quick run)
make dev ARGS="pr create"

# Build
make build

# Install
make install
```

See [go-version/README.md](./go-version/README.md) for more details.

## 📚 Documentation

- [Go Version Guide](./go-version/README.md) - Recommended implementation
- [Shell Version Guide](./shell-version/README.md) - Legacy implementation
- [Migration Guide](./go-version/MIGRATION.md) - Migrating from shell to Go version
- [Project Structure](./go-version/STRUCTURE.md) - Codebase overview
- [Contributing](./go-version/CONTRIBUTING.md) - How to contribute

## 🆚 Comparison

| Feature | Go Version | Shell Version |
|---------|-----------|---------------|
| Installation | Single binary | Multiple dependencies |
| Performance | Fast | Slower |
| Cross-platform | ✅ Yes | 🔶 Unix-like only |
| Error handling | ✅ Robust | 🔶 Basic |
| Maintenance | ✅ Active | 🔶 Legacy |
| Cancel support | ✅ Ctrl+C anywhere | ❌ Limited |
| Type safety | ✅ Yes | ❌ No |

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./go-version/CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Uses [GitHub CLI](https://cli.github.com/) for GitHub API
- Uses [Jira CLI](https://github.com/ankitpokhrel/jira-cli) for Jira integration

---

**Maintainer**: [@Wangggym](https://github.com/Wangggym)

💡 **Tip**: Star this repo if you find it useful!
