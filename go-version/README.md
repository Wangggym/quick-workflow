# Quick Workflow (Go Version)

> A modern, blazing-fast CLI tool for streamlined GitHub and Jira workflows

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Wangggym/quick-workflow?style=flat&logo=github)](https://github.com/Wangggym/quick-workflow/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/Wangggym/quick-workflow/build.yml?branch=master&style=flat&logo=github-actions)](https://github.com/Wangggym/quick-workflow/actions)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat)](https://github.com/Wangggym/quick-workflow)
[![iCloud Sync](https://img.shields.io/badge/iCloud-Sync%20Enabled-blue?style=flat&logo=icloud)](ICLOUD_MIGRATION.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

## 🚀 What's New in Go Version

This is a complete rewrite of the original Shell-based quick-workflow tool in Go, bringing:

- **📦 Single Binary** - No dependencies, just download and run
- **⚡ Faster** - Native performance, instant startup
- **🔒 Type Safe** - Catch errors at compile time
- **🎨 Better UX** - Interactive prompts, colored output
- **🧪 Testable** - Comprehensive test coverage
- **🌍 Cross-platform** - Works on macOS, Linux, and Windows
- **☁️ iCloud Sync** - Automatic config sync across Mac devices (macOS only)

## ✨ Features

- **PR Creation** - Create PRs with automatic branch creation, commit, and push
- **PR Merging** - Merge PRs and clean up branches automatically
- **Quick Update** - Commit and push with PR title as commit message
- **Jira Integration** - Automatically update Jira status and add PR links
- **Interactive CLI** - Beautiful prompts and progress indicators
- **Configuration Management** - Simple setup with `qkflow init`
- **iCloud Sync** - Seamlessly sync configs across all your Mac devices ☁️

## 📦 Installation

### Option 1: Download Binary (Recommended)

```bash
# macOS (Apple Silicon)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Linux
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/
```

### Option 2: Install with Go

```bash
go install github.com/Wangggym/quick-workflow/cmd/qkflow@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/Wangggym/quick-workflow.git
cd quick-workflow/go-version
make build
sudo cp bin/qkflow /usr/local/bin/
```

## ⚙️ Setup

### Prerequisites

- Git installed and configured
- GitHub CLI (`gh`) installed and authenticated: `gh auth login`
- Jira API token: [Get one here](https://id.atlassian.com/manage-profile/security/api-tokens)

### Initial Configuration

Run the interactive setup:

```bash
qkflow init
```

This will prompt you for:
- Email address
- GitHub token (auto-detected from `gh` CLI)
- Jira service address (e.g., https://your-domain.atlassian.net)
- Jira API token
- Optional: Branch prefix
- Optional: OpenAI API key for AI features

**Configuration Storage:**

✨ **NEW**: On macOS, all configs are automatically saved to iCloud Drive in a single directory and synced across all your devices!

- **macOS with iCloud Drive**: Synced across devices ☁️
  - 📂 All configs in: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
- **Local Storage** (fallback): 
  - 📂 All configs in: `~/.qkflow/`

Both locations contain:
- `config.yaml` - Main configuration
- `jira-status.json` - Jira status mappings

Run `qkflow config` to see your actual storage location.

📖 See [iCloud Migration Guide](ICLOUD_MIGRATION.md) for more details.

## 🎯 Usage

### Create a Pull Request

```bash
# With Jira ticket
qkflow pr create PROJ-123

# Without Jira ticket (will prompt)
qkflow pr create

# Interactive mode (no arguments)
qkflow pr create
```

**What it does:**
1. ✅ Fetches Jira issue details (if ticket provided)
2. ✅ Prompts for PR title and description
3. ✅ Lets you select change types (feat, fix, docs, etc.)
4. ✅ Creates a new git branch
5. ✅ Commits your staged changes
6. ✅ Pushes to remote
7. ✅ Creates GitHub PR
8. ✅ Adds PR link to Jira
9. ✅ Updates Jira status (optional)
10. ✅ Copies PR URL to clipboard

### Merge a Pull Request

```bash
# Merge PR by number
qkflow pr merge 123

# Interactive mode
qkflow pr merge
```

**What it does:**
1. ✅ Fetches PR details
2. ✅ Confirms merge with you
3. ✅ Merges the PR on GitHub
4. ✅ Deletes remote branch (optional)
5. ✅ Switches to main branch
6. ✅ Deletes local branch (optional)
7. ✅ Updates Jira status (optional)
8. ✅ Adds merge comment to Jira

### Quick Update (qkupdate)

```bash
# Quick commit and push with PR title as commit message
qkflow update
```

**What it does:**
1. ✅ Gets the current PR title from GitHub
2. ✅ Stages all changes (git add --all)
3. ✅ Commits with PR title as commit message
4. ✅ Pushes to origin
5. ✅ Falls back to "update" if no PR found

This is perfect for quick updates to an existing PR!

### Other Commands

```bash
# Show current configuration
qkflow config

# Show version
qkflow version

# Get help
qkflow help
qkflow pr --help
```

## 🏗️ Project Structure

```
go-version/
├── cmd/
│   └── qkflow/
│       ├── main.go              # Entry point
│       └── commands/
│           ├── root.go          # Root command
│           ├── init.go          # Setup wizard
│           ├── pr.go            # PR commands
│           ├── pr_create.go     # Create PR
│           ├── pr_merge.go      # Merge PR
│           ├── update.go        # Quick update
│           └── jira.go          # Jira commands
├── internal/
│   ├── github/
│   │   └── client.go           # GitHub API client
│   ├── jira/
│   │   ├── client.go           # Jira API client
│   │   └── status_cache.go     # Status cache
│   ├── git/
│   │   └── operations.go       # Git operations
│   ├── ai/
│   │   └── client.go           # AI client
│   └── ui/
│       └── prompt.go           # User interface
├── pkg/
│   └── config/
│       └── config.go           # Configuration management
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make (optional but recommended)

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to GOPATH/bin
make install

# Run tests
make test

# Run linters
make lint

# Format code
make fmt
```

### Quick Development

```bash
# Run without building
go run ./cmd/qkflow pr create

# Or use Makefile
make run-pr-create
make run-pr-merge
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test ./internal/github/...
go test ./internal/jira/...
```

## 📝 Configuration Files

### Storage Location

qkflow intelligently stores configurations based on your system:

**macOS with iCloud Drive** (Recommended):
- Configs are stored in iCloud Drive and automatically synced across your Mac devices
- All configs in: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
  - `config.yaml` - Main configuration
  - `jira-status.json` - Jira status mappings

**Local Storage** (Fallback):
- All configs in: `~/.qkflow/`
  - `config.yaml` - Main configuration
  - `jira-status.json` - Jira status mappings

Run `qkflow config` to see your actual storage location.

### Configuration Format

```yaml
email: your.email@example.com
jira_api_token: your_jira_token
jira_service_address: https://your-domain.atlassian.net
github_token: ghp_your_github_token
branch_prefix: feature  # optional
openai_key: sk-your_openai_key  # optional
```

## 🔒 Security

- Tokens are stored securely in your config directory (local or iCloud)
- File permissions are set to `0600` (user read/write only)
- iCloud storage is encrypted and secure
- Never commit the config file or share it
- Use environment variables for CI/CD:
  ```bash
  export QK_GITHUB_TOKEN=xxx
  export QK_JIRA_API_TOKEN=xxx
  ```

**Note**: If using iCloud Drive, your configurations will sync across your Mac devices automatically, providing a seamless experience.

## 🚧 Migration from Shell Version

See [MIGRATION.md](MIGRATION.md) for detailed migration guide.

**Quick comparison:**

| Feature | Shell Version | Go Version |
|---------|--------------|------------|
| Installation | Clone + dependencies | Single binary |
| Configuration | `.zshrc` / `.bashrc` | `qk init` |
| Dependencies | `gh`, `jira`, `jq`, etc. | None (self-contained) |
| Speed | ~1-2s startup | <100ms startup |
| Platform | macOS/Linux | macOS/Linux/Windows |
| Updates | `git pull` | Download new binary |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Original Shell version: [quick-workflow](https://github.com/Wangggym/quick-workflow)
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts
- [go-github](https://github.com/google/go-github) - GitHub API client
- [go-jira](https://github.com/andygrunwald/go-jira) - Jira API client

## 📞 Support

- 🐛 [Report a bug](https://github.com/Wangggym/quick-workflow/issues/new?labels=bug)
- 💡 [Request a feature](https://github.com/Wangggym/quick-workflow/issues/new?labels=enhancement)
- 📖 [Documentation](https://github.com/Wangggym/quick-workflow/wiki)

---

Made with ❤️ by [Wangggym](https://github.com/Wangggym)

