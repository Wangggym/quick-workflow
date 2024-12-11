#!/bin/bash

# 下载函数
download_file() {
    local url="$1"
    local filename="$2"
    local output_path="$DOWNLOAD_DIR/$filename"
    
    echo "下载: $filename"
    curl -# -L "$url" -o "$output_path"
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
    echo "$attachments" | while read -r line; do
        filename=$(echo "$line" | cut -d'|' -f1)
        url=$(echo "$line" | cut -d'|' -f2)
        download_file "$url" "$filename"
    done

    # 检查下载是否成功
    if [ ! -d "$DOWNLOAD_DIR" ] || [ -z "$(ls -A "$DOWNLOAD_DIR")" ]; then
        echo "错误: 下载失败或目录为空"
        return 1
    fi

    # 检查文件
    echo "下载的文件列表:"
    ls -l "$DOWNLOAD_DIR"
    return 0
} 