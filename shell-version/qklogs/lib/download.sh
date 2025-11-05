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
    local max_parallel=${3:-5}  # 默认最大并行数为5

    # Set output directory
    export OUTPUT_DIR="$output_dir"
    
    # Export download function for subprocesses
    export -f download_file

    # Download all log files
    echo "ℹ️  Starting parallel downloads..."
    local failed=0
    local running=0
    local pids=()

    while IFS='§' read -r filename url; do
        # 如果当前运行的进程数达到最大值，等待任意一个完成
        while [ $running -ge $max_parallel ]; do
            for pid in ${pids[@]}; do
                if ! kill -0 $pid 2>/dev/null; then
                    wait $pid || failed=1
                    running=$((running - 1))
                    break
                fi
            done
            sleep 0.1
        done

        # 启动新的下载进程
        download_file "$url" "$filename" &
        pids+=($!)
        running=$((running + 1))
    done <<< "$attachments"

    # 等待所有下载完成
    for pid in ${pids[@]}; do
        wait $pid || failed=1
    done

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