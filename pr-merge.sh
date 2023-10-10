#!/bin/bash

# merge pr and delete branch
read -p 'pr id: ' pr_id
pr_merge=$(gh pr merge $pr_id --squash --delete-branch)
echo $pr_merge

# find jira ticket and change status to done
script_dir="$(dirname "$0")"

text=$(grep $pr_id "${script_dir}/work-history.txt")

echo $text

jira_ticket=$(echo "$text" | awk -F ',' '{print $1}')

result=$(jira issue move $jira_ticket "FE fixed")

echo Merged pr and changed status successfully!
