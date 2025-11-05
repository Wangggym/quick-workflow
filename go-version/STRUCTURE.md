# Project Structure

```
go-version/
│
├── 📝 Configuration Files
│   ├── go.mod                          # Go module definition
│   ├── go.sum                          # Dependency checksums
│   ├── Makefile                        # Build automation
│   ├── .gitignore                      # Git ignore rules
│   └── .golangci.yml                   # Linter configuration
│
├── 🎯 Main Application
│   └── cmd/
│       └── qk/
│           ├── main.go                 # Application entry point
│           └── commands/
│               ├── root.go             # Root command & app setup
│               ├── init.go             # Setup wizard (qk init)
│               ├── pr.go               # PR command group
│               ├── pr_create.go        # PR creation logic
│               └── pr_merge.go         # PR merging logic
│
├── 🔒 Internal Packages
│   └── internal/
│       ├── github/
│       │   └── client.go               # GitHub API client
│       │       ├── NewClient()
│       │       ├── CreatePullRequest()
│       │       ├── GetPullRequest()
│       │       ├── MergePullRequest()
│       │       └── ParseRepositoryFromURL()
│       │
│       ├── jira/
│       │   └── client.go               # Jira API client
│       │       ├── NewClient()
│       │       ├── GetIssue()
│       │       ├── UpdateStatus()
│       │       ├── AddComment()
│       │       ├── AddPRLink()
│       │       └── GetProjectStatuses()
│       │
│       ├── git/
│       │   └── operations.go           # Git command wrappers
│       │       ├── CheckStatus()
│       │       ├── GetCurrentBranch()
│       │       ├── CreateBranch()
│       │       ├── Commit()
│       │       ├── Push()
│       │       ├── DeleteBranch()
│       │       ├── DeleteRemoteBranch()
│       │       ├── GetRemoteURL()
│       │       └── SanitizeBranchName()
│       │
│       └── ui/
│           └── prompt.go               # User interface helpers
│               ├── Success()
│               ├── Error()
│               ├── Warning()
│               ├── Info()
│               ├── PromptInput()
│               ├── PromptPassword()
│               ├── PromptConfirm()
│               ├── PromptSelect()
│               └── PromptMultiSelect()
│
├── 📦 Public Packages
│   └── pkg/
│       └── config/
│           └── config.go               # Configuration management
│               ├── Load()
│               ├── Get()
│               ├── Save()
│               ├── Validate()
│               └── IsConfigured()
│
├── 🛠️ Scripts
│   └── scripts/
│       ├── install.sh                  # Installation script
│       ├── test.sh                     # Test runner script
│       └── release.sh                  # Release automation
│
├── 🤖 CI/CD
│   └── .github/
│       └── workflows/
│           └── build.yml               # GitHub Actions workflow
│
├── 📚 Documentation
│   ├── README.md                       # Main documentation
│   ├── MIGRATION.md                    # Migration guide (Shell → Go)
│   ├── QUICKSTART.md                   # 5-minute quick start
│   ├── CONTRIBUTING.md                 # Contribution guidelines
│   ├── PROJECT_OVERVIEW.md             # Technical overview
│   └── STRUCTURE.md                    # This file
│
├── 📄 Legal
│   └── LICENSE                         # MIT License
│
└── 🔨 Build Output (gitignored)
    └── bin/                            # Compiled binaries
        ├── qk                          # Current platform
        ├── qk-darwin-amd64             # macOS Intel
        ├── qk-darwin-arm64             # macOS Apple Silicon
        ├── qk-linux-amd64              # Linux
        └── qk-windows-amd64.exe        # Windows
```

## 📊 File Statistics

| Category | Files | Lines of Code (est.) |
|----------|-------|---------------------|
| Go Source | 10 | ~2,000 |
| Documentation | 6 | ~2,500 |
| Scripts | 3 | ~300 |
| Config | 5 | ~200 |
| **Total** | **24** | **~5,000** |

## 🔗 Package Dependencies

```
cmd/qk/commands
  ├─→ internal/github
  ├─→ internal/jira
  ├─→ internal/git
  ├─→ internal/ui
  └─→ pkg/config

internal/github
  └─→ pkg/config

internal/jira
  └─→ pkg/config

internal/git
  └─→ (no internal deps)

internal/ui
  └─→ (no internal deps)

pkg/config
  └─→ (no internal deps)
```

## 📖 Key Files Explained

### Entry Point
- **`cmd/qk/main.go`**: Application entry point, calls command execution

### Commands
- **`commands/root.go`**: Root command setup, version, config display
- **`commands/init.go`**: Interactive setup wizard for first-time config
- **`commands/pr.go`**: PR command group (parent of create/merge)
- **`commands/pr_create.go`**: Complete PR creation workflow
- **`commands/pr_merge.go`**: Complete PR merging workflow

### Core Libraries
- **`internal/github/client.go`**: GitHub API wrapper with typed interfaces
- **`internal/jira/client.go`**: Jira API wrapper with status management
- **`internal/git/operations.go`**: Git command execution and branch management
- **`internal/ui/prompt.go`**: User interaction and colored output

### Infrastructure
- **`pkg/config/config.go`**: Configuration loading, saving, validation
- **`Makefile`**: Build commands (build, test, lint, install)
- **`.github/workflows/build.yml`**: CI/CD pipeline for multi-platform builds

### Documentation
- **`README.md`**: User-facing documentation (installation, usage)
- **`MIGRATION.md`**: Detailed migration guide from Shell version
- **`QUICKSTART.md`**: 5-minute getting started guide
- **`CONTRIBUTING.md`**: Guidelines for contributors
- **`PROJECT_OVERVIEW.md`**: Technical architecture and design

## 🎯 Navigation Guide

### For Users
1. Start → `README.md` (overview)
2. Setup → `QUICKSTART.md` (5 min)
3. Migration → `MIGRATION.md` (if from Shell)

### For Developers
1. Architecture → `PROJECT_OVERVIEW.md`
2. Structure → This file (`STRUCTURE.md`)
3. Contributing → `CONTRIBUTING.md`
4. Code → Start from `cmd/qk/main.go`

### For Building
1. Dependencies → `go.mod`
2. Build → `Makefile`
3. CI/CD → `.github/workflows/build.yml`
4. Release → `scripts/release.sh`

## 🔍 Code Organization Principles

### 1. **Separation of Concerns**
- `cmd/` - CLI interface and user interaction
- `internal/` - Business logic and API clients
- `pkg/` - Reusable utilities

### 2. **Dependency Direction**
- Commands depend on internal packages
- Internal packages depend on pkg
- No circular dependencies

### 3. **Visibility**
- `internal/` - Private to this module
- `pkg/` - Can be imported by other modules
- `cmd/` - Application entry points

### 4. **Testing**
- Each package has its own tests
- Mock external dependencies
- Table-driven test patterns

## 📦 External Dependencies

```
Core Framework:
├── github.com/spf13/cobra           # CLI framework
├── github.com/spf13/viper           # Configuration
└── github.com/AlecAivazis/survey/v2 # Interactive prompts

API Clients:
├── github.com/google/go-github/v57  # GitHub API
├── github.com/andygrunwald/go-jira  # Jira API
└── golang.org/x/oauth2              # OAuth2 auth

Utilities:
└── github.com/fatih/color           # Terminal colors
```

## 🎨 Design Patterns Used

1. **Factory Pattern**: Client creation (`NewClient()`)
2. **Command Pattern**: CLI commands structure
3. **Repository Pattern**: API clients abstract data access
4. **Facade Pattern**: Simplified interfaces for complex operations
5. **Strategy Pattern**: Different PR types and workflows

## 🚀 Build Artifacts

After running `make build-all`:

```
bin/
├── qk-darwin-amd64      # macOS Intel (12-15MB)
├── qk-darwin-arm64      # macOS M1/M2 (12-15MB)
├── qk-linux-amd64       # Linux x86_64 (12-15MB)
└── qk-windows-amd64.exe # Windows 64-bit (12-15MB)
```

## 📈 Metrics

- **Total Lines**: ~5,000
- **Go Files**: 10
- **Packages**: 6
- **Commands**: 4
- **Functions**: ~80
- **Structs**: ~15
- **Interfaces**: ~5

---

Last Updated: 2025-11-04  
Version: 1.0.0

