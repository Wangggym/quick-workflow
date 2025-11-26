# PR Approve Command Guide

Quick guide for using the new `qkflow pr approve` command.

## 🚀 Quick Start

### Basic Approval

```bash
# Approve a specific PR by number (uses default 👍 comment)
qkflow pr approve 123

# Approve a PR by URL (works from anywhere!)
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# Also works with /files, /commits, /checks URLs
qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

# Auto-detect PR from current branch
qkflow pr approve
```

### With Comment

By default, all approvals use 👍 as the comment. Customize with `-c` flag:

```bash
# Default approval (with 👍)
qkflow pr approve 123

# Custom comment
qkflow pr approve 123 -c "LGTM! 🎉"

# Approve by URL with custom comment
qkflow pr approve https://github.com/owner/repo/pull/456 -c "Great work!"

# Long comment with flag
qkflow pr approve 123 --comment "Great work! All tests passed. Approved for merge."
```

### Auto-Merge

```bash
# Approve and merge in one step
qkflow pr approve 123 --merge

# Approve by URL and merge
qkflow pr approve https://github.com/owner/repo/pull/789 --merge

# Short flag
qkflow pr approve 123 -m

# With comment and merge
qkflow pr approve 123 -c "LGTM!" -m

# URL with comment and merge
qkflow pr approve https://github.com/owner/repo/pull/789 -c "Ship it! 🚀" -m
```

## 🌐 URL Support (NEW!)

Now you can approve PRs from **any repository** without being in the git directory!

### Why Use URLs?

1. **Cross-Repository**: Approve PRs from different repos
2. **No Directory Change**: Work from anywhere
3. **Browser to CLI**: Copy URL from GitHub directly
4. **Batch Operations**: Script approvals across multiple repos

### URL Examples

```bash
# Approve a PR from a different repo
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# Your colleague shares a PR link, approve it instantly
qkflow pr approve https://github.com/company/frontend/pull/456 -c "Looks good!"

# Merge someone else's PR by URL
qkflow pr merge https://github.com/team/backend/pull/789

# Approve and merge with URL
qkflow pr approve https://github.com/org/project/pull/123 -c "LGTM! 🎉" -m
```

### Supported URL Formats

All these formats work:

```bash
# HTTPS (most common)
https://github.com/owner/repo/pull/123

# With /files suffix (Files tab)
https://github.com/brain/planning-api/pull/2001/files

# With /commits suffix (Commits tab)
https://github.com/owner/repo/pull/123/commits

# With /checks suffix (Checks tab)
https://github.com/owner/repo/pull/123/checks

# HTTP
http://github.com/owner/repo/pull/123

# Without protocol
github.com/owner/repo/pull/123

# With query params (parsed correctly)
https://github.com/owner/repo/pull/123?comments=all
https://github.com/owner/repo/pull/123/files?file-filters%5B%5D=.js

# With fragments (parsed correctly)
https://github.com/owner/repo/pull/123#discussion_r123456
```

### When to Use URL vs Number

| Scenario | Use | Example |
|----------|-----|---------|
| Your repo, in git dir | Number | `qkflow pr approve 123` |
| Your repo, current branch | None | `qkflow pr approve` |
| Different repo | URL | `qkflow pr approve https://...` |
| Shared link from browser | URL | Copy & paste URL |
| Scripting multiple repos | URL | Loop through URLs |

## 📋 Common Workflows

### Workflow 1: Quick Code Review

You're reviewing a colleague's PR:

```bash
# 1. Check out their branch (optional)
git fetch origin
git checkout feature-branch

# 2. Review the code...

# 3. Approve
qkflow pr approve
# Finds PR automatically from branch
# Adds optional comment
# Done!
```

### Workflow 2: Approve and Merge

You have approval rights and want to merge immediately:

```bash
# One command to approve and merge
qkflow pr approve 123 -c "Approved! Merging now." --merge

# What happens:
# ✅ Approves PR #123
# ✅ Adds comment
# ✅ Checks if mergeable
# ✅ Confirms with you
# ✅ Merges PR
# ✅ Deletes remote branch
# ✅ Switches to main
# ✅ Deletes local branch
```

### Workflow 3: Batch Approvals

Multiple PRs to review:

```bash
# List all open PRs first
gh pr list

# Approve them one by one
qkflow pr approve 101 -c "Approved"
qkflow pr approve 102 -c "Approved"
qkflow pr approve 103 -c "Approved"
```

### Workflow 4: Interactive Mode

Don't know the PR number? Let the tool help:

```bash
# Run without arguments
qkflow pr approve

# What happens:
# 1. Tries to find PR from current branch
# 2. If not found, shows list of all open PRs
# 3. You select one
# 4. Asks for optional comment
# 5. Approves!
```

## 🔍 Use Cases

### Use Case 1: Team Lead Approval

As a team lead, you need to approve PRs daily:

```bash
# Morning routine: approve all ready PRs
for pr in 121 122 123; do
  qkflow pr approve $pr -c "Reviewed and approved"
done
```

### Use Case 2: CI/CD Integration

Add to your CI pipeline:

```bash
#!/bin/bash
# Auto-approve dependabot PRs after tests pass
if [[ "$PR_AUTHOR" == "dependabot" ]] && [[ "$TESTS_PASSED" == "true" ]]; then
  qkflow pr approve $PR_NUMBER -c "Auto-approved: Tests passed" -m
fi
```

### Use Case 3: Hotfix Workflow

Fast-track urgent fixes:

```bash
# Create hotfix
git checkout -b hotfix/critical-bug
# ... fix the bug ...
git add .
qkflow pr create

# Get it approved and merged ASAP
qkflow pr approve 999 -c "Critical hotfix - merging immediately" -m
```

## 🎯 Pro Tips

### Tip 1: Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
alias approve='qkflow pr approve'
alias merge='qkflow pr approve --merge'

# Usage:
approve 123 -c "LGTM"
merge 123
```

### Tip 2: Comment Templates

Save common comments:

```bash
# In your shell config
export APPROVE_LGTM="Looks good to me! 👍"
export APPROVE_MINOR="Approved with minor comments. Please address in follow-up."
export APPROVE_EXCELLENT="Excellent work! 🎉"

# Usage:
qkflow pr approve 123 -c "$APPROVE_LGTM"
```

### Tip 3: Check Before Merge

Before using `--merge`, verify the PR:

```bash
# View PR details
gh pr view 123

# Check CI status
gh pr checks 123

# Approve and merge if all good
qkflow pr approve 123 -m
```

### Tip 4: Branch Protection

If branch protection is enabled:

```bash
# Just approve - let GitHub merge rules handle the rest
qkflow pr approve 123 -c "Approved"

# Don't use --merge if you need multiple approvals
```

## ⚠️ Error Handling

### PR Not Found

```bash
$ qkflow pr approve 999
❌ Failed to get PR: Pull request not found
```

**Solution:** Check PR number with `gh pr list`

### PR Already Closed

```bash
$ qkflow pr approve 123
❌ PR is not open (state: closed)
```

**Solution:** PR is already merged or closed

### Not Mergeable

```bash
$ qkflow pr approve 123 -m
✅ PR approved!
⚠️ Cannot merge PR: PR has conflicts and cannot be merged
```

**Solution:** Resolve conflicts first, then merge

### No PR for Branch

```bash
$ qkflow pr approve
⚠️ No PR found for branch: feature-xyz
Do you want to select a PR from the list? (Y/n)
```

**Solution:** Either select from list or specify PR number

## 🆚 Comparison

### GitHub CLI (`gh`)

```bash
# Approve
gh pr review 123 --approve --body "LGTM"

# Then merge separately
gh pr merge 123

# Then cleanup
git checkout main
git pull
git branch -D feature-branch
```

### qkflow (New!)

```bash
# All in one!
qkflow pr approve 123 -c "LGTM" -m
```

**Benefits:**
- ✅ Fewer commands
- ✅ Auto-cleanup
- ✅ Interactive prompts
- ✅ Branch auto-detection
- ✅ Merge validation

## 🔗 Related Commands

- `qkflow pr create` - Create a new PR
- `qkflow pr merge` - Merge without approving first
- `gh pr view 123` - View PR details
- `gh pr checks 123` - Check CI status

## 📚 Full Reference

### Command Syntax

```
qkflow pr approve [pr-number] [flags]
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--comment` | `-c` | Add a comment with the approval |
| `--merge` | `-m` | Automatically merge after approval |
| `--help` | `-h` | Show help information |

### Exit Codes

- `0`: Success
- `1`: Error (PR not found, not mergeable, etc.)

### Environment

Requires:
- GitHub token configured (`qkflow config`)
- Git repository with remote
- Write access to repository

---

**Need Help?**

```bash
qkflow pr approve --help
qkflow help
```

**Found a Bug?**

Open an issue on GitHub with details!

