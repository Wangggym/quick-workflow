# Jira Issue Reader for Cursor AI

A powerful tool to read and export Jira issues, optimized for use with Cursor AI.

## 🎯 Quick Start

### For Cursor Users (Recommended)

The simplest way to use this in Cursor:

```bash
# In Cursor terminal, run:
qkflow jira read NA-9245

# Then in Cursor chat, simply say:
"总结刚才读取的 Jira ticket 内容"
```

Cursor will automatically read the exported files and provide you with a comprehensive summary!

## 📚 Available Commands

### 1. `show` - Quick Terminal View

Display issue information directly in the terminal.

```bash
# Basic view (metadata only)
qkflow jira show NA-9245

# Full view (includes description and comments)
qkflow jira show NA-9245 --full
```

**Use when:**
- You need a quick peek at the issue
- You only need text content
- You want the fastest response

### 2. `export` - Complete Export with Files

Export issue to local files with optional images.

```bash
# Export text only
qkflow jira export NA-9245

# Export with all images and attachments
qkflow jira export NA-9245 --with-images

# Export to custom directory
qkflow jira export NA-9245 -o ~/jira-exports/ --with-images
```

**Output structure:**
```
/tmp/qkflow/jira/NA-9245/
├── README.md           # How to use in Cursor
├── content.md          # Main content (Markdown)
└── attachments/        # Downloaded files (if --with-images)
    ├── screenshot.png
    └── diagram.jpg
```

**Use when:**
- You need images/attachments
- You want to keep a local copy
- You need formatted Markdown

### 3. `read` - Intelligent Mode ⭐️ **RECOMMENDED**

Automatically decides the best way to present the issue.

```bash
# Auto mode (smart decision)
qkflow jira read NA-9245
```

**How it works:**
- ✅ **Has images?** → Exports to files with images
- ✅ **Text only?** → Displays directly in terminal
- ✅ Automatically optimized for Cursor

**Use when:**
- Working with Cursor AI (best experience)
- You want the tool to decide the best format
- You're not sure if the issue has images

### 4. `clean` - Clean Up Exports

Remove exported files to free up disk space.

```bash
# Clean specific issue
qkflow jira clean NA-9245

# Clean all exports
qkflow jira clean --all

# Preview what would be deleted (dry run)
qkflow jira clean --all --dry-run

# Force delete without confirmation
qkflow jira clean --all --force
```

## 🎨 Usage Examples in Cursor

### Example 1: Simple Text Analysis

```
You in Cursor: "通过 qkflow 读取 NA-9245 并总结"

Cursor executes: qkflow jira read NA-9245
Cursor responds: "这个 ticket (NA-9245) 是关于..."
```

### Example 2: With Images

```
You in Cursor: "用 qkflow 读取 NA-9245 的所有内容包括图片，分析架构设计"

Cursor executes: qkflow jira export NA-9245 --with-images
Cursor reads: content.md + all images in attachments/
Cursor responds: "根据架构图，这个系统包含..."
```

### Example 3: Manual Control

```bash
# Step 1: Export (you run this)
qkflow jira export NA-9245 --with-images

# Step 2: Tell Cursor what to read
"Read /tmp/qkflow/jira/NA-9245/content.md and analyze the architecture diagram"

# Step 3: Clean up when done
qkflow jira clean NA-9245
```

## 🔧 Configuration

Make sure your Jira credentials are configured:

```bash
qkflow init
```

Required settings:
- `jira_service_address`: Your Jira instance URL (e.g., https://brain-ai.atlassian.net)
- `jira_api_token`: Your Jira API token
- `email`: Your Jira email

## 💡 Tips & Best Practices

### For Cursor Users

1. **Use `read` command by default** - It's optimized for AI consumption
2. **Be specific in prompts** - Tell Cursor what you want to know
3. **Clean up regularly** - Use `clean --all` to free up space

### Command Comparison

| Command | Speed | Images | Output | Best For |
|---------|-------|--------|--------|----------|
| `show` | ⚡️ Fastest | ❌ | Terminal | Quick peek |
| `show --full` | ⚡️ Fast | ❌ | Terminal | Full text |
| `export` | 🐌 Slower | ❌ | Files | Text archive |
| `export --with-images` | 🐌 Slowest | ✅ | Files | Complete archive |
| `read` ⭐️ | ⚡️ Smart | ✅ Smart | Smart | **Cursor AI** |

### Cursor Prompt Templates

```bash
# General summary
"通过 qkflow 读取 <ISSUE-KEY> 并总结内容"

# Specific analysis
"用 qkflow 读取 <ISSUE-KEY>，分析技术方案"

# With context
"读取 <ISSUE-KEY>，对比我们当前的实现方式"

# With images
"qkflow 读取 <ISSUE-KEY> 包括所有图片，分析架构设计"
```

## 📊 Output Formats

### Terminal Output (show command)

```
╔═══════════════════════════════════════════════════════╗
║ 🎫 NA-9245: Implement user authentication            ║
╚═══════════════════════════════════════════════════════╝

📋 Type:        Story
📊 Status:      In Progress
🏷️  Priority:    High
👤 Assignee:    John Doe

🔗 View in Jira: https://brain-ai.atlassian.net/browse/NA-9245
```

### Markdown Output (export command)

```markdown
---
issue_key: NA-9245
title: Implement user authentication
type: Story
status: In Progress
priority: High
---

# NA-9245: Implement user authentication

## 📊 Metadata
...

## 📝 Description
...

## 📎 Attachments (3)
1. **screenshot.png** (245 KB)
   ![screenshot.png](./attachments/screenshot.png)
...
```

### Cursor-Optimized Output (read command)

The `read` command provides special output markers that Cursor recognizes:

```
✅ Exported to: /tmp/qkflow/jira/NA-9245/

Main file: /tmp/qkflow/jira/NA-9245/content.md
Images: /tmp/qkflow/jira/NA-9245/attachments/ (3 files)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💡 CURSOR: Please read the following files:
1. /tmp/qkflow/jira/NA-9245/content.md
2. All images in /tmp/qkflow/jira/NA-9245/attachments/
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 🚀 Advanced Usage

### Batch Processing

```bash
# Export multiple issues
for issue in NA-9245 NA-9246 NA-9247; do
  qkflow jira export $issue --with-images
done

# Tell Cursor to analyze all
"Read all Jira exports in /tmp/qkflow/jira/ and summarize"
```

### Custom Workflows

```bash
# Create a script for your team
#!/bin/bash
ISSUE_KEY=$1
qkflow jira read "$ISSUE_KEY"
echo "Ready for Cursor to analyze!"
```

### Integration with Other Tools

```bash
# Export and open in VS Code/Cursor
qkflow jira export NA-9245 --with-images
code /tmp/qkflow/jira/NA-9245/content.md
```

## 🐛 Troubleshooting

### "Failed to create Jira client"
- Check your config: `cat ~/.config/quick-workflow/config.yaml`
- Verify your API token is valid
- Ensure `jira_service_address` is correct

### "Failed to get issue"
- Verify the issue key is correct (e.g., NA-9245)
- Check you have permission to view the issue
- Try accessing the issue in your browser first

### "Failed to download attachment"
- Check your network connection
- Verify your API token has attachment download permissions
- Some files may be restricted

### Cursor not reading files
- Make sure the export command completed successfully
- Check the file paths in the output
- Try manually attaching the file: `@/tmp/qkflow/jira/NA-9245/content.md`

## 📝 Notes

- Exports are temporary and stored in `/tmp/qkflow/jira/` by default
- Use `clean` command regularly to free up space
- Images are only downloaded with `--with-images` flag
- The `read` command is specifically designed for Cursor AI integration

## 🔗 Related Commands

- `qkflow init` - Configure Jira credentials
- `qkflow pr create` - Create PR (can auto-link to Jira)
- `qkflow jira list` - List Jira status mappings

---

**Made with ❤️ for Cursor AI users**

