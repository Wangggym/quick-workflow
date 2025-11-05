#!/bin/bash

# 自动加载 .env 文件（如果存在）
script_dir="$(cd "$(dirname "$0")" && pwd)"
if [ -f "$script_dir/.env" ]; then
  set -o allexport
  source "$script_dir/.env"
  set +o allexport
fi 