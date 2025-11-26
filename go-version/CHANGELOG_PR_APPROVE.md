# Changelog - PR Approve Feature

## [New Feature] PR Approve Command with URL Support

Added `qkflow pr approve` command to streamline PR approval workflow, with full support for GitHub PR URLs!

### What's New

#### 🎯 PR Approve Command

Quickly approve pull requests with optional auto-merge capability.

**Basic Usage:**
```bash
# Approve PR by number (default comment: 👍)
qkflow pr approve 123

# 🆕 Approve PR by URL (works from anywhere!)
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# 🆕 Works with /files, /commits, /checks URLs too!
qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

# Auto-detect from current branch
qkflow pr approve

# Custom comment
qkflow pr approve 123 --comment "LGTM! 🎉"
qkflow pr approve 123 -c "Looks good!"

# 🆕 Approve by URL with comment
qkflow pr approve https://github.com/owner/repo/pull/456 -c "Great work!"

# Approve and auto-merge (with default 👍)
qkflow pr approve 123 --merge
qkflow pr approve 123 -m

# 🆕 Approve by URL and merge
qkflow pr approve https://github.com/owner/repo/pull/789 -c "Ship it! 🚀" -m
```

### Features

1. **🆕 URL Support**
   - Approve PRs from **any repository** using GitHub URLs
   - No need to be in the git directory
   - Copy URL from browser and approve instantly
   - Works with: `https://github.com/owner/repo/pull/NUMBER`
   - **🆕 Also works with tab-specific URLs:**
     - `/files` - Files tab
     - `/commits` - Commits tab
     - `/checks` - Checks tab
   - Also supports PR numbers for local repos

2. **🆕 Default Comment**
   - All approvals use 👍 as default comment
   - Customize with `-c` flag if needed
   - No more empty approvals!

3. **Auto-Detection**
   - Automatically finds PR from current branch if no argument provided
   - Interactive selection if multiple PRs exist

4. **Smart Approval**
   - Approves PR on GitHub
   - Optional comment customization
   - Validates PR state before approval

5. **Auto-Merge (Optional)**
   - `--merge` flag enables automatic merging after approval
   - Checks if PR is mergeable before proceeding
   - Confirms with user before merging
   - Cleans up branches after successful merge

6. **Interactive Mode**
   - Asks for confirmation before merge
   - Lists all open PRs for selection

### Command Flags

- `-c, --comment string`: Add a comment with the approval
- `-m, --merge`: Automatically merge the PR after approval
- `-h, --help`: Display help information

### Workflow Examples

**Quick Approve:**
```bash
# In a branch with PR
qkflow pr approve
# → Finds PR automatically → Approves → Done!
```

**🆕 Approve by URL (Cross-Repo):**
```bash
# Someone shares a PR link in Slack/Email
qkflow pr approve https://github.com/brain/planning-api/pull/2001 -c "LGTM!"
# → Works from anywhere! → Approves → Done!
```

**Approve with Comment:**
```bash
qkflow pr approve 123 -c "Great work! Ship it! 🚀"
# → Approves with comment → Done!
```

**Approve and Merge:**
```bash
qkflow pr approve 123 --merge
# → Approves → Checks mergeability → Confirms → Merges → Cleans up branches
```

**🆕 Cross-Repository Workflow:**
```bash
# Approve multiple PRs from different repos
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"
qkflow pr approve https://github.com/team/backend/pull/200 -c "Approved"
qkflow pr approve https://github.com/team/mobile/pull/300 -c "Approved"
# → No directory changes needed!
```

**Code Review Workflow:**
```bash
# Review someone's PR
qkflow pr approve 456 -c "Reviewed and approved. Minor suggestions in comments."

# Or review and merge immediately
qkflow pr approve 456 -c "LGTM!" --merge
```

### Technical Details

#### New Functions in `internal/github/client.go`

- `ApprovePullRequest(owner, repo, number, body)`: Approve a PR with optional comment
- `IsPRMergeable(owner, repo, number)`: Check if a PR can be merged
- `ParsePRFromURL(url)`: 🆕 Parse owner, repo, and PR number from GitHub URL
- `IsPRURL(s)`: 🆕 Check if a string is a GitHub PR URL

#### Updated Functions

- `pr_approve.go`: Added URL parsing support at the start of command execution
- `pr_merge.go`: Added URL parsing support for consistency

#### New Files

- `cmd/qkflow/commands/pr_approve.go`: PR approve command implementation
- `internal/git/operations.go`: Added `CheckoutBranch()` function

#### Supported URL Formats

```bash
# Full HTTPS URL
https://github.com/owner/repo/pull/123

# With tab suffixes (NEW!)
https://github.com/brain/planning-api/pull/2001/files
https://github.com/owner/repo/pull/123/commits
https://github.com/owner/repo/pull/123/checks

# HTTP URL
http://github.com/owner/repo/pull/123

# Without protocol
github.com/owner/repo/pull/123

# With query parameters (parsed correctly)
https://github.com/owner/repo/pull/123?comments=all
https://github.com/owner/repo/pull/123/files?file-filters%5B%5D=.js

# With URL fragments (parsed correctly)
https://github.com/owner/repo/pull/123#discussion_r123456
```

**Pro Tip:** Copy URL from any PR tab and it just works! 🎉

### Integration

The approve command integrates seamlessly with existing PR workflow:

```bash
# Full PR workflow
qkflow pr create          # Create PR
qkflow pr approve 123 -m  # Review and merge
# Done! 🎉
```

### Error Handling

- ✅ Validates PR exists and is open
- ✅ Checks GitHub authentication
- ✅ Verifies PR is mergeable before merge
- ✅ Handles branch cleanup failures gracefully
- ✅ Provides clear error messages

### Comparison with GitHub CLI

**Before (with gh):**
```bash
# For local repo PR
gh pr review 123 --approve --body "LGTM"
gh pr merge 123
git checkout main
git branch -D feature-branch

# For different repo (need to specify owner/repo)
gh pr review 123 --approve --body "LGTM" --repo owner/repo
gh pr merge 123 --repo owner/repo
```

**After (with qkflow):**
```bash
# For local repo PR
qkflow pr approve 123 -c "LGTM" -m
# Everything done automatically! 🚀

# 🆕 For different repo (just copy the URL!)
qkflow pr approve https://github.com/owner/repo/pull/123 -c "LGTM" -m
# Works from anywhere! No --repo flag needed! 🎉
```

**Key Advantages:**
- ✅ URL support - no need to specify owner/repo separately
- ✅ Copy-paste friendly from browser
- ✅ Auto cleanup after merge
- ✅ Single command for approve + merge
- ✅ Works cross-repository seamlessly

### Future Enhancements

Potential improvements for future versions:

- [ ] Request changes workflow
- [ ] Bulk approve multiple PRs
- [ ] Conditional merge (wait for CI)
- [ ] Custom merge strategies (squash, rebase)
- [ ] Slack/notification integration

---

**Related Commands:**
- `qkflow pr create` - Create a new PR
- `qkflow pr merge` - Merge an existing PR
- `qkflow jira read` - Read Jira issues

**Documentation:**
- See [README.md](README.md) for general usage
- See [CONTRIBUTING.md](CONTRIBUTING.md) for development guide

