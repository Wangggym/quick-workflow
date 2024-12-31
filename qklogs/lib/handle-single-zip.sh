#!/bin/bash

# 处理单个zip文件的函数
handle_single_zip_file() {
    local output_dir="$1"
    local log_output_folder_name="$2"
    
    # 检查是否只有一个 zip 文件
    local zip_files=("$output_dir"/*.zip)
    if [ ${#zip_files[@]} -eq 1 ]; then
        local zip_file_name=$(basename "${zip_files[0]}")
        echo "✅ Single zip file found: $zip_file_name, skipping merge step"
        
        local extract_dir
        if [[ -n "${log_output_folder_name}" ]]; then
            extract_dir="$output_dir/${log_output_folder_name}"
        else
            extract_dir="$output_dir/merged"
        fi
        
        unzip "$output_dir/$zip_file_name" -d "$extract_dir"
        return 0
    fi
    return 1
}