# PR URL Support - Feature Summary

## 🎉 What's New

You can now use **GitHub PR URLs** directly with `pr approve` and `pr merge` commands! No need to be in the repository directory or remember PR numbers.

## ✨ Quick Examples

### Before (PR number only)
```bash
# Had to be in the repo directory
cd ~/projects/my-repo
qkflow pr approve 123
```

### After (URL support!)
```bash
# Works from anywhere! 🚀
qkflow pr approve https://github.com/brain/planning-api/pull/2001
qkflow pr approve https://github.com/brain/planning-api/pull/2001 -c "LGTM!" -m
```

## 📖 Usage

### Basic Syntax

Both commands now accept:
1. **PR Number** (requires being in git repo)
2. **Full GitHub URL** (works from anywhere!)
3. **No argument** (auto-detect from current branch)

### Commands Updated

#### `qkflow pr approve`

```bash
# By PR number (in repo)
qkflow pr approve 123

# By URL (anywhere)
qkflow pr approve https://github.com/owner/repo/pull/456

# With options
qkflow pr approve https://github.com/owner/repo/pull/789 -c "LGTM!" -m
```

#### `qkflow pr merge`

```bash
# By PR number (in repo)
qkflow pr merge 123

# By URL (anywhere)
qkflow pr merge https://github.com/owner/repo/pull/456
```

## 🌟 Use Cases

### 1. Cross-Repository Reviews

Review PRs from multiple repositories without changing directories:

```bash
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"
qkflow pr approve https://github.com/team/backend/pull/200 -c "Approved"
qkflow pr approve https://github.com/team/mobile/pull/300 -c "Approved"
```

### 2. Browser to CLI Workflow

1. Open PR in GitHub web
2. Copy URL from address bar
3. Paste into terminal:

```bash
qkflow pr approve https://github.com/company/project/pull/1234 -c "Looks good!"
```

### 3. Slack/Email Integration

Teammate shares a PR link in Slack? Approve it instantly:

```bash
# Copy link from Slack
qkflow pr approve <paste-url-here> -c "Reviewed and approved"
```

### 4. Scripting Across Repos

Automate approvals across multiple repositories:

```bash
#!/bin/bash
PR_URLS=(
  "https://github.com/org/repo1/pull/10"
  "https://github.com/org/repo2/pull/20"
  "https://github.com/org/repo3/pull/30"
)

for url in "${PR_URLS[@]}"; do
  qkflow pr approve "$url" -c "Auto-approved by bot"
done
```

## 🔧 Technical Details

### Supported URL Formats

All these formats work:

```
✅ https://github.com/owner/repo/pull/123
✅ https://github.com/owner/repo/pull/123/files      (Files tab)
✅ https://github.com/owner/repo/pull/123/commits    (Commits tab)
✅ https://github.com/owner/repo/pull/123/checks     (Checks tab)
✅ http://github.com/owner/repo/pull/123
✅ github.com/owner/repo/pull/123
✅ https://github.com/owner/repo/pull/123?comments=all
✅ https://github.com/owner/repo/pull/123#discussion_r123456
✅ https://github.com/owner/repo/pull/123/files?file-filters%5B%5D=.js
```

**Pro Tip:** Just copy the URL from any PR tab (Overview, Files, Commits, Checks) and it will work! 🎉

### URL Parsing

The tool automatically:
- Detects if argument is a URL or number
- Extracts owner, repo, and PR number from URL
- Handles query parameters and fragments
- Validates URL format

### Error Handling

Clear error messages for common issues:

```bash
$ qkflow pr approve invalid-url
❌ Invalid PR number or URL: invalid-url
ℹ️  Expected: PR number (e.g., '123') or GitHub URL (e.g., 'https://github.com/owner/repo/pull/123')
```

## 🧪 Testing

Comprehensive test coverage for URL parsing:

```bash
# Run tests
cd go-version
go test -v ./internal/github/

# All tests pass! ✅
# - TestParsePRFromURL (8 test cases)
# - TestIsPRURL (8 test cases)
# - TestParseRepositoryFromURL (5 test cases)
```

## 📚 Documentation

Updated documentation:
- ✅ `README.md` - Main documentation with examples
- ✅ `PR_APPROVE_GUIDE.md` - Detailed usage guide
- ✅ `CHANGELOG_PR_APPROVE.md` - Feature changelog
- ✅ `PR_URL_SUPPORT.md` - This file
- ✅ Command help text (`--help`)

## 🎯 Benefits

1. **🌐 Cross-Repository**: Approve PRs from any repo
2. **⚡ Faster**: No need to navigate to repo directory
3. **📋 Copy-Paste Friendly**: Direct from browser
4. **🤖 Scriptable**: Easy batch operations
5. **🔗 Shareable**: Use links from Slack/Email
6. **✨ Backwards Compatible**: Still works with PR numbers

## 🚀 Get Started

1. **Update qkflow** (if already installed):
   ```bash
   qkflow update-cli
   ```

2. **Try it out**:
   ```bash
   # Find a PR on GitHub, copy the URL
   qkflow pr approve https://github.com/.../pull/123 -c "Testing URL support!"
   ```

3. **See help**:
   ```bash
   qkflow pr approve --help
   qkflow pr merge --help
   ```

## 💡 Tips

### Alias for Quick Access

Add to your `.bashrc` or `.zshrc`:

```bash
# Quick approve
alias gha='qkflow pr approve'

# Usage:
gha https://github.com/owner/repo/pull/123 -c "LGTM!"
```

### Use with `pbpaste` (macOS)

```bash
# Copy URL in browser, then:
qkflow pr approve "$(pbpaste)" -c "Approved!"
```

### Integration with GitHub CLI

Combine with `gh` CLI:

```bash
# List PRs with gh, approve with qkflow
gh pr list
qkflow pr approve https://github.com/owner/repo/pull/123 -m
```

## 📝 Examples in Action

### Example 1: Quick Review (Default Comment)

```bash
$ qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

ℹ️  Detected GitHub PR URL, parsing...
✅ Parsed: brain/planning-api PR #2001
ℹ️  Fetching PR #2001...
ℹ️  PR: feat: Add user authentication
ℹ️  Branch: feature/auth -> main
ℹ️  State: open
ℹ️  Using default comment: 👍 (use -c flag to customize)
ℹ️  Approving PR #2001...
✅ PR approved with comment: 👍

ℹ️  PR approved. Use 'qkg pr merge' to merge it later, or run with --merge flag to auto-merge.
```

### Example 1b: Quick Review (Custom Comment)

```bash
$ qkflow pr approve https://github.com/brain/planning-api/pull/2001 -c "LGTM!"

ℹ️  Detected GitHub PR URL, parsing...
✅ Parsed: brain/planning-api PR #2001
ℹ️  Fetching PR #2001...
ℹ️  PR: feat: Add user authentication
ℹ️  Branch: feature/auth -> main
ℹ️  State: open
ℹ️  Approving PR #2001...
✅ PR approved with comment: LGTM!

ℹ️  PR approved. Use 'qkg pr merge' to merge it later, or run with --merge flag to auto-merge.
```

### Example 2: Approve and Merge

```bash
$ qkflow pr approve https://github.com/team/backend/pull/456 -c "Ship it! 🚀" -m

ℹ️  Detected GitHub PR URL, parsing...
✅ Parsed: team/backend PR #456
ℹ️  Fetching PR #456...
ℹ️  PR: fix: Database connection timeout
ℹ️  Branch: fix/db-timeout -> main
ℹ️  State: open
ℹ️  Approving PR #456...
✅ PR approved with comment: Ship it! 🚀

ℹ️  Checking if PR is mergeable...
❓ Proceed with merging the PR? (Y/n) y
ℹ️  Merging PR #456...
✅ 🎉 PR merged successfully!
ℹ️  Deleting remote branch fix/db-timeout...
✅ Remote branch deleted

✅ All done! 🎉
```

## 🤝 Backwards Compatibility

All existing workflows still work:

```bash
# PR number (requires being in repo)
qkflow pr approve 123

# Auto-detect from current branch
qkflow pr approve

# Interactive selection
qkflow pr approve
# → Shows list of PRs to choose from
```

---

**Enjoy the new URL support! 🎉**

For questions or feedback, please open an issue on GitHub.

