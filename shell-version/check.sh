#!/bin/bash
script_dir="$(dirname "$0")"

chmod +x $script_dir/*.sh

if ! "$script_dir/check-git-status.sh"; then
    exit 1
fi

# 先检查网络连接
if "$script_dir/check-network.sh"; then
    # 网络检查通过，直接退出成功
    exit 0
else
    # 网络检查失败，尝试检查代理设置
    if ! "$script_dir/check-set-proxy.sh"; then
        exit 1
    fi
    # 代理设置检查完成后再次检查网络
    if ! "$script_dir/check-network.sh"; then
        exit 1
    fi
fi

exit 0
