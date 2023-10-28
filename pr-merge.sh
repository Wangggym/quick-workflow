#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

source $script_dir/pr-jira.sh

pr_id=$1
while [ -z "$pr_id" ]; do
    read -p 'PR id(require): ' pr_id
done

gh browse $pr_id

echo 'Do you want to continue merging PR? (y/n)'

read choice

if [ ! "$choice" = "y" ]; then
    exit 0
fi

gh pr merge $pr_id --squash --delete-branch

jira_merge $pr_id
