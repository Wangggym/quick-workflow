#!/bin/bash

# 合并分片文件
merge_logs() {
    local download_dir="$1"

    if [ ! -f "$download_dir/log.z01" ]; then
        echo "没有找到需要合并的分片文件"
        return 0
    fi

    echo "检测到分片文件，开始合并..."
    cd "$download_dir" || exit 1
    
    # 检查所有需要的文件是否存在
    if [ ! -f "log.zip" ]; then
        echo "错误: 未找到 log.zip 文件"
        return 1
    fi

    # 计算分片文件数量
    local split_count=$(ls log.z* 2>/dev/null | wc -l)
    echo "找到 $split_count 个分片文件"

    # 尝试使用 7z 合并（如果安装了的话）
    if command -v 7z &> /dev/null; then
        echo "使用 7z 合并文件..."
        if 7z x log.zip && 7z x log.z01; then
            echo "7z 合并成功"
            return 0
        fi
    fi

    # 尝试使用 cat 直接合并
    echo "使用 cat 合并文件..."
    if cat log.zip log.z* > merged.zip; then
        echo "文件合并成功"
        
        # 验证合并后的文件
        if unzip -t merged.zip >/dev/null 2>&1; then
            echo "合并文件验证成功"
            return 0
        else
            echo "合并文件验证失败"
            rm -f merged.zip
        fi
    fi

    echo "cat 合并失败，尝试使用 zip 命令修复..."
    # 使用 yes 命令自动回答 zip 的提示
    if yes | zip -F log.zip --out merged.zip; then
        if unzip -t merged.zip >/dev/null 2>&1; then
            echo "zip 修复成功"
            return 0
        fi
    fi

    echo "尝试使用 zip -FF..."
    # 使用 yes 命令自动回答 zip 的提示
    if yes | zip -FF log.zip --out merged.zip; then
        if unzip -t merged.zip >/dev/null 2>&1; then
            echo "zip -FF 修复成功"
            return 0
        fi
    fi

    echo "所有合并尝试都失败"
    return 1
} 