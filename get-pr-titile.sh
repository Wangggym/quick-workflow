#!/bin/bash

script_dir="$(dirname "$0")"

source "$script_dir/base.sh"

pr_title=$(gh pr view --json title --jq .title)

if [ -n "$pr_title" ]; then
    echo -e "${y}PR Title: $pr_title"
else
    echo -e "${r}Failed to retrieve PR title."
fi
