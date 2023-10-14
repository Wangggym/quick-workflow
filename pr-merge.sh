#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

read -p 'pr id: ' pr_id

pr_merge=$(gh pr merge $pr_id --squash --delete-branch)

echo $pr_merge

# find jira ticket and change status to IN QA

text=$(grep $pr_id "${script_dir}/work-history.txt")

echo $text

jira_ticket=$(echo "$text" | awk -F ',' '{print $1}')

if [[ $jira_ticket == *BSF* ]]; then
    status="FE fixed"
else
    status="IN QA"
fi

result=$(jira issue move $jira_ticket "${status}")

echo $result

echo "✓ 1. PR merged"
echo "✓ 2. Jira status changed"
