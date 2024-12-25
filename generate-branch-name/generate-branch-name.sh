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
    if [ ! -f "generate-branch-name.py" ]; then
        echo "错误：找不到 generate-branch-name.py 文件"
        exit 1
    fi

    # 将输入文本传递给 Python 脚本
    result=$(python3 $SCRIPT_DIR/generate-branch-name.py $1 $2)

    # 检查函数返回值
    if [ $? -ne 0 ]; then
        echo "Generate branch name failed: $result"
        return 1
    fi

    echo "$result"
    # echo "Generate branch name success: $result"
    # echo "$result" > /tmp/branch_name_$3.txt

    exit 0
}

# # 主函数
# __main__() {
#     check_dependencies
#     generate_branch_name "测试一些API是否成功" "$BRAIN_AI_KEY"  # 传递第一个命令行参数

#     # 读取保存的结果
#     result=$(cat /tmp/stream_result.txt)

#     # 这里可以添加对结果的进一步处理
#     echo "处理后的结果："
#     echo "$result"

#     # 清理临时文件
#     rm -f /tmp/stream_result.txt
# }

# # 执行主函数，传递所有命令行参数
# __main__