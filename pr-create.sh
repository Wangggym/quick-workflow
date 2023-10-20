#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

source $script_dir/base.sh
source $script_dir/pr-body.sh
source $script_dir/multiselect.sh

read -p 'Jira ticket: ' jira_ticket
read -p 'Issue name: ' issue_name
read -p 'Short description: ' short_description

if [[ $jira_ticket == *BSF* ]]; then
    status="In Review"
    preselection=("true" "false" "false")
else
    status="UNDER Review"
    preselection=("false" "true" "false")
fi

echo 'Types of changes:'
multiselect "false" result types_of_changes preselection

branch_name=${jira_ticket}--$(echo "$issue_name" | sed 's/[^a-zA-Z0-9]/-/g')
commit_title=${jira_ticket}': '${issue_name}
pr_body=$(getPRbody $jira_ticket $short_description result)

# echo $branch_name
# echo $commit_title
# echo $pr_body

git checkout -b $branch_name
git add . && git commit -m "${commit_title}" && git push
pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H $branch_name)

jira issue assign $jira_ticket $(jira me)
jira issue move $jira_ticket "${status}"
echo $pr_url | jira issue comment add $jira_ticket

# write history
pr_id=$(echo "$pr_url" | grep -oE '[0-9]+$')
echo "${jira_ticket},${pr_id}" >>"${script_dir}/work-history.txt"


echo $pr_url | pbcopy
echo ✓ Successfully copied $pr_url to clipboard 

open $pr_url
