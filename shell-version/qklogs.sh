#!/bin/bash

# 开启调试模式
set -x

# 检查参数
if [ $# -eq 0 ]; then
    echo "使用方法: $0 <JIRA-ISSUE-KEY>"
    exit 1
fi

ISSUE_KEY="$1"
# 修改目录结构，所有文件都放在这个目录下
BASE_DIR="$HOME/Downloads/logs_${ISSUE_KEY}"
DOWNLOAD_DIR="$BASE_DIR/downloads"

# 创建目录
mkdir -p "$DOWNLOAD_DIR"

# 获取JIRA认证信息
JIRA_TOKEN="${JIRA_API_TOKEN}"
if [ -z "$JIRA_TOKEN" ]; then
    echo "错误: 未设置JIRA_API_TOKEN环境变量"
    exit 1
fi

# 使用jira-cli获取issue信息并提取附件URL
echo "正在获取 $ISSUE_KEY 的附件信息..."

# 保存原始输出到临时文件
TEMP_FILE=$(mktemp)
jira issue view "$ISSUE_KEY" --plain > "$TEMP_FILE"

# 提取附件信息 - 使用更严格的格式匹配，并确保正确的文件名格式
ATTACHMENTS=$(awk '
    /^[[:space:]]*[0-9]+\. log\.(zip|z[0-9]+)[[:space:]]*$/ {
        filename = $2
        getline
        if ($0 ~ /^[[:space:]]*http/) {
            url = $0
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", url)
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", filename)
            # 移除序号，只保留文件名
            sub(/^[0-9]+\. /, "", filename)
            print filename "|" url
            next
        }
    }
' "$TEMP_FILE" | sort | uniq)

# 清理临时文件
rm "$TEMP_FILE"

# 检查是否有附件
if [ -z "$ATTACHMENTS" ]; then
    echo "未找到附件"
    exit 1
fi

echo "找到以下附件:"
echo "$ATTACHMENTS"

# 调用下载脚本
SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
source "$SCRIPT_DIR/download-logs.sh"

# 执行下载
if ! download_attachments "$ATTACHMENTS" "$DOWNLOAD_DIR"; then
    echo "下载过程中发生错误"
    exit 1
fi

# 检查是否需要合并
if [ -f "$DOWNLOAD_DIR/log.z01" ]; then
    echo "检测到分片文件，需要合并..."
    # 调用merge-logs.sh合并日志，合并后的文件也放在同一目录下
    echo "合并日志文件..."
    "$SCRIPT_DIR/merge-logs.sh" "$DOWNLOAD_DIR"
else
    echo "只有单个文件，不需要合并..."
fi

echo "完成! 所有文件位于: $BASE_DIR"
echo "文件列表:"
ls -l "$BASE_DIR"
echo ""
echo "如果需要解压，请进入目录 $BASE_DIR 后操作"