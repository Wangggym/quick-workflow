#!/bin/bash

# Enable error checking
set -e

# Import base functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        echo "❌ Script failed with exit code: $exit_code" >&2
    fi
    exit $exit_code
}

trap cleanup EXIT

# Check parameters
if [ $# -eq 0 ]; then
    echo "❌ Usage: $0 <JIRA-ISSUE-KEY>" >&2
    exit 1
fi

ISSUE_KEY="$1"
OUTPUT_DIR="$HOME/Downloads/logs_${ISSUE_KEY}"

# Delete files if exists
if [ -d "$OUTPUT_DIR" ]; then
    rm -rf "$OUTPUT_DIR"
fi
# Create directory
mkdir -p "$OUTPUT_DIR"

# Get script directory and import required functions
LIB_DIR="$SCRIPT_DIR/lib"
source "$LIB_DIR/get-urls.sh"
source "$LIB_DIR/download.sh"
source "$LIB_DIR/merge.sh"
source "$LIB_DIR/handle-single-zip.sh"

# Step 1: Get attachment URLs
echo "ℹ️  Step 1: Getting attachment URLs"
ATTACHMENTS=$(get_attachment_urls "$ISSUE_KEY")
RET=$?
if [ $RET -ne 0 ]; then
    echo "❌ Failed to get attachment URLs (code: $RET)" >&2
    exit $RET
fi

if [ -z "$ATTACHMENTS" ]; then
    echo "❌ No attachments found" >&2
    exit 1
fi

echo "✅ Found attachments:"
echo "$ATTACHMENTS"

# Step 2: Download files
echo "ℹ️  Step 2: Downloading attachments"
if ! download_attachments "$ATTACHMENTS" "$OUTPUT_DIR"; then
    echo "❌ Download failed" >&2
    exit 1
fi

# # 检查是否只有一个 zip 文件
if handle_single_zip_file "$OUTPUT_DIR" "$LOG_OUTPUT_FOLDER_NAME"; then
    exit 0
fi

# Step 3: Merge files (只有多个文件时才执行)
echo "ℹ️  Step 3: Merging files"
if ! merge_logs "$OUTPUT_DIR"; then
    echo "❌ Merge failed" >&2
    exit 1
fi

if [[ "${LOG_DELETE_WHEN_OPERATION_COMPLETED}" == 1 ]]; then
    echo "ℹ️  Deleting downloaded files..."
    # 删除原始日志文件
    rm "$OUTPUT_DIR"/log.z*
fi

# Open the merged file
if [ -f "$OUTPUT_DIR/merged.zip" ]; then
    echo "ℹ️  Opening merged file..."

    if [[ -n "${LOG_OUTPUT_FOLDER_NAME}" ]]; then
        unzip "$OUTPUT_DIR/merged.zip" -d "$OUTPUT_DIR/${LOG_OUTPUT_FOLDER_NAME}"
    else
        unzip "$OUTPUT_DIR/merged.zip" -d "$OUTPUT_DIR/merged"
    fi

    if [[ "${LOG_DELETE_WHEN_OPERATION_COMPLETED}" == 1 ]]; then
        echo "ℹ️  Deleting merged.zip ..."
        rm "$OUTPUT_DIR/merged.zip"
    fi

    echo "✅ All done! Files are in: $OUTPUT_DIR"
    echo "ℹ️  File list:"
    ls -lh "$OUTPUT_DIR"

    open "$OUTPUT_DIR"
else
    echo "ℹ️  Opening file directory..."
    open "$OUTPUT_DIR"
fi

