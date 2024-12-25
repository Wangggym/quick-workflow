#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 检查依赖
check_dependencies() {
    # 检查 Python3
    if ! command -v python3 &> /dev/null; then
        echo "错误：未安装 Python3"
        exit 1
    fi

    # 检查 requests 包
    python3 -c "import requests" 2>/dev/null || {
        echo "正在安装 requests 包..."
        pip3 install requests
    }
}

# 执行 Python 脚本
generate_branch_name() {
    if [ ! -f "$SCRIPT_DIR/generate-branch-name.py" ]; then
        echo "错误：找不到 generate-branch-name.py 文件"
    fi

    # 将输入文本传递给 Python 脚本
    python3 $SCRIPT_DIR/generate-branch-name.py "$1" "$2" "$3"

    # 检查函数返回值
    if [ $? -ne 0 ]; then
        echo -e $n "Generate branch name failed: $result"
    fi
}