#!/bin/bash
# 用法: ./qkpy.sh <python_file_path>
# 例子: ./qkpy.sh src/case/planning/travel_extractor.py

# 检测可用的Python解释器
get_python_cmd() {
    if command -v python3 &> /dev/null; then
        echo "python3"
    elif command -v python &> /dev/null; then
        echo "python"
    else
        echo "Error: No Python interpreter found" >&2
        exit 1
    fi
}

if [ $# -lt 1 ]; then
  echo "Usage: $0 <python_file_path>"
  exit 1
fi

PYTHON_CMD=$(get_python_cmd)
PYTHONPATH=$(pwd) $PYTHON_CMD "$1"