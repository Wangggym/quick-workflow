#!/bin/bash

# Download function
download_file() {
    local url="$1"
    local filename="$2"
    local output_path="$OUTPUT_DIR/$filename"
    
    echo "ℹ️  Downloading: $filename"
    # Add more curl options for redirects and timeouts
    curl -# -L --retry 3 --retry-delay 2 --connect-timeout 10 -o "$output_path" "$url"
    local ret=$?
    
    # Check if download was successful
    if [ $ret -ne 0 ]; then
        echo "❌ curl download failed with code: $ret" >&2
        return 1
    fi
    
    if [ ! -f "$output_path" ] || [ ! -s "$output_path" ]; then
        echo "❌ File $filename download failed or is empty" >&2
        return 1
    fi
    
    echo "✅ Downloaded: $filename ($(du -h "$output_path" | cut -f1))"
    return 0
}

# Main download function
download_attachments() {
    local attachments="$1"
    local output_dir="$2"

    # Set output directory
    export OUTPUT_DIR="$output_dir"
    
    # Export download function for subprocesses
    export -f download_file

    # Download all log files
    echo "ℹ️  Starting downloads..."
    local failed=0
    while IFS='§' read -r filename url; do
        if ! download_file "$url" "$filename"; then
            failed=1
            echo "⚠️  Failed to download: $filename" >&2
            continue
        fi
    done <<< "$attachments"

    # Check for download failures
    if [ $failed -eq 1 ]; then
        echo "❌ Some downloads failed" >&2
        return 1
    fi

    # Check output directory
    if [ ! -d "$OUTPUT_DIR" ] || [ -z "$(ls -A "$OUTPUT_DIR")" ]; then
        echo "❌ Output directory is missing or empty" >&2
        return 1
    fi

    # List files
    echo "✅ Downloads completed. Files:"
    ls -lh "$OUTPUT_DIR"
    return 0
} 