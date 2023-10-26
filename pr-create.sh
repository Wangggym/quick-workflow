#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

source $script_dir/base.sh
source $script_dir/pr-body.sh
source $script_dir/multiselect.sh

jira_ticket=$1
while [ -z "$jira_ticket" ]; do
    read -p 'Jira ticket(require): ' jira_ticket
done

read -p 'Issue desc(require): ' issue_desc
while [ -z "$issue_desc" ]; do
    read -p 'Issue desc(require): ' issue_desc
done

read -p 'Short description(optional): ' short_description

short_description=${short_description:-"Not yet"}

if [[ $jira_ticket == *BSF* ]]; then
    status="In Review"
    preselection=("true" "false" "false")
else
    status="UNDER Review"
    preselection=("false" "true" "false")
fi

echo 'Types of changes:'
multiselect "false" result types_of_changes preselection

commit_title=${jira_ticket}': '${issue_desc}
pr_body=$(getPRbody $jira_ticket $short_description result)
branch_name=${jira_ticket}--$(echo "$issue_desc" | sed 's/[^a-zA-Z0-9]/-/g')

if [ -n "${GH_BRANCH_PREFIX}" ]; then
    branch_name=${GH_BRANCH_PREFIX}/${branch_name}
fi

# echo $branch_name
# echo $commit_title
# echo $pr_body

git checkout -b $branch_name
git add . && git commit -m "${commit_title}" && git push -u origin $branch_name
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
