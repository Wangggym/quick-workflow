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

# Create directory
mkdir -p "$OUTPUT_DIR"

# Get script directory and import required functions
LIB_DIR="$SCRIPT_DIR/lib"
source "$LIB_DIR/get-urls.sh"
source "$LIB_DIR/download.sh"
source "$LIB_DIR/merge.sh"

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

# Step 3: Merge files
echo "ℹ️  Step 3: Merging files"
if ! merge_logs "$OUTPUT_DIR"; then
    echo "❌ Merge failed" >&2
    exit 1
fi

echo "✅ All done! Files are in: $OUTPUT_DIR"
echo "ℹ️  File list:"
ls -lh "$OUTPUT_DIR"

# Open the merged file
if [ -f "$OUTPUT_DIR/merged.zip" ]; then
    echo "ℹ️  Opening merged file..."
    open "$OUTPUT_DIR/merged.zip"
fi 