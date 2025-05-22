#!/bin/bash
# 用法: ./qkpy.sh <python_file_path>
# 例子: ./qkpy.sh src/case/planning/travel_extractor.py

if [ $# -lt 1 ]; then
  echo "Usage: $0 <python_file_path>"
  exit 1
fi

PYTHONPATH=$(pwd) python "$1"