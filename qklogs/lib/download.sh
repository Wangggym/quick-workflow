#!/bin/bash

# 下载函数
download_file() {
    local url="$1"
    local filename="$2"
    local output_path="$DOWNLOAD_DIR/$filename"
    
    echo "下载: $filename"
    curl -# -L -o "$output_path" "$url"
    
    # 检查下载是否成功
    if [ ! -f "$output_path" ] || [ ! -s "$output_path" ]; then
        echo "错误: 文件 $filename 下载失败或为空" >&2
        return 1
    fi
    return 0
}

# 主下载函数
download_attachments() {
    local attachments="$1"
    local download_dir="$2"

    # 设置下载目录
    export DOWNLOAD_DIR="$download_dir"
    
    # 导出下载函数供子进程使用
    export -f download_file

    # 下载所有log文件
    echo "开始下载日志文件..."
    echo "$attachments" | while IFS='§' read -r filename url; do
        if ! download_file "$url" "$filename"; then
            return 1
        fi
    done

    # 检查下载是否成功
    if [ ! -d "$DOWNLOAD_DIR" ] || [ -z "$(ls -A "$DOWNLOAD_DIR")" ]; then
        echo "错误: 下载失败或目录为空" >&2
        return 1
    fi

    # 检查文件
    echo "下载的文件列表:"
    ls -l "$DOWNLOAD_DIR"
    return 0
} 