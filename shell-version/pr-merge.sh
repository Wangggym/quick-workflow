#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

source $script_dir/base.sh
source $script_dir/pr-jira.sh

pr_id=$1

if [ -z "$pr_id" ]; then
    pr_url=$(gh pr status --json url -q '.currentBranch.url')
    pr_id=$(echo "$pr_url" | grep -oE '[0-9]+$')
    if [ -n $"$pr_id" ]; then
        echo -e $y Find current branch PR: $pr_url
    fi
fi

while [ -z "$pr_id" ]; do
    read -p 'PR id(require): ' pr_id
done

# 直接使用普通合并
if ! gh pr merge $pr_id --merge --delete-branch; then
    echo -e "${r}Normal merge failed${n}"
    exit 1
fi

jira_merge $pr_id
