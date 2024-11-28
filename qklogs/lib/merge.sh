#!/bin/bash

# Merge split files
merge_logs() {
    local output_dir="$1"
    cd "$output_dir" || { echo "❌ Cannot enter output directory" >&2; return 1; }

    # Check for split files
    if [ ! -f "log.z01" ] || [ ! -f "log.zip" ]; then
        echo "ℹ️  No split files found"
        return 0
    fi

    echo "ℹ️  Found split files, merging..."
    
    # Calculate split file count
    local split_count=$(ls log.z* 2>/dev/null | wc -l)
    echo "ℹ️  Found $split_count split files"

    # Merge using cat
    echo "ℹ️  Merging files..."
    if cat log.zip log.z* > merged.zip; then
        echo "✅ Merge completed"
        
        # Check file size
        local merged_size=$(stat -f%z "merged.zip")
        local expected_size=$(($(stat -f%z "log.zip") + $(stat -f%z "log.z01")))
        
        if [ "$merged_size" -eq "$expected_size" ]; then
            echo "✅ Merge successful: file size matches ($(du -h "merged.zip" | cut -f1))"
        else
            echo "⚠️  File size mismatch, but file was created" >&2
            echo "ℹ️  Expected size: $expected_size bytes"
            echo "ℹ️  Actual size: $merged_size bytes"
        fi
        return 0
    fi

    echo "❌ Merge failed"
    return 1
} 