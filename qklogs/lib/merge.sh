#!/bin/bash

# Merge split files
merge_logs() {
    local output_dir="$1"
    cd "$output_dir" || { echo "❌ Cannot enter output directory" >&2; return 1; }

    # Check for split files
    if [ ! -f "log.zip" ]; then
        echo "ℹ️  No split files found"
        return 0
    fi

    echo "ℹ️  Found split files, merging..."
    
    # Calculate split file count
    local split_count=$(ls log.z[0-9][0-9] 2>/dev/null | wc -l)
    echo "ℹ️  Found $split_count split files"

    # Merge using cat
    echo "ℹ️  Merging files..."
    if cat log.zip log.z[0-9][0-9] > merged.zip; then
        echo "✅ Merge completed"
        
        # Check total size of all split files
        local total_size=0
        for f in log.zip log.z[0-9][0-9]; do
            local size=$(stat -f%z "$f")
            total_size=$((total_size + size))
        done
        
        local merged_size=$(stat -f%z "merged.zip")
        
        if [ "$merged_size" -eq "$total_size" ]; then
            echo "✅ Merge successful: file size matches ($(du -h "merged.zip" | cut -f1))"
        else
            echo "⚠️  File size mismatch, but file was created" >&2
            echo "ℹ️  Expected size: $total_size bytes"
            echo "ℹ️  Actual size: $merged_size bytes"
        fi
        return 0
    fi

    echo "❌ Merge failed"
    return 1
} 