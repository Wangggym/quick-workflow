#!/bin/bash

# 检查输入目录
if [ $# -eq 0 ]; then
    echo "使用方法: $0 <分片文件目录> [输出文件路径]"
    exit 1
fi

# 进入目录
cd "$1" || exit 1

# 设置输出文件路径，如果没有提供就用当前目录的logs.zip
output_file="${2:-$1/logs.zip}"

# 检查并合并文件
if [ ! -f "log.zip" ]; then
    echo "错误: 找不到log.zip文件"
    exit 1
fi

# 合并文件
cat log.zip log.z[0-9][0-9] > "$output_file"

# 验证大小
total_size=0
for f in log.zip log.z[0-9][0-9]; do
    size=$(stat -f%z "$f")
    total_size=$((total_size + size))
done

if [ $(stat -f%z "$output_file") -eq $total_size ]; then
    echo "文件已成功合并到: $output_file"
else
    echo "错误: 合并文件大小不匹配"
fi