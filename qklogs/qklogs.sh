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

# 获取脚本所在目录
SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
LIB_DIR="$SCRIPT_DIR/lib"

# 导入所需函数
source "$LIB_DIR/get-urls.sh"
source "$LIB_DIR/download.sh"
source "$LIB_DIR/merge.sh"

# 获取附件URL
ATTACHMENTS=$(get_attachment_urls "$ISSUE_KEY")
if [ $? -ne 0 ]; then
    exit 1
fi

echo "找到以下附件:"
echo "$ATTACHMENTS"

# 执行下载
if ! download_attachments "$ATTACHMENTS" "$DOWNLOAD_DIR"; then
    echo "下载过程中发生错误"
    exit 1
fi

# 合并文件
if ! merge_logs "$DOWNLOAD_DIR"; then
    echo "合并过程中发生错误"
    exit 1
fi

echo "完成! 所有文件位于: $BASE_DIR"
echo "文件列表:"
ls -l "$BASE_DIR"
echo ""
echo "如果需要解压，请进入目录 $BASE_DIR 后操作" 