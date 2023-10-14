#!/bin/bash
script_dir="$(dirname "$0")"

chmod +x $script_dir/*.sh

if ! "$script_dir/check-git-status.sh"; then
    exit 1
fi

if ! "$script_dir/check-network.sh"; then
    exit 1
fi

exit 0
