#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

read -p 'PR id: ' pr_id

gh browse $pr_id

echo 'Do you want to continue? (y/n)'

read choice

if [ "$choice" = "n" ]; then
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
