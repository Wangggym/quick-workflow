#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

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

text=$(grep $pr_id "${script_dir}/work-history.txt")
jira_ticket=$(echo "$text" | awk -F ',' '{print $1}')

if [[ $jira_ticket == *BSF* ]]; then
    status="FE fixed"
else
    status="IN QA"
fi

gh pr merge $pr_id --squash --delete-branch

jira issue move $jira_ticket "${status}"
