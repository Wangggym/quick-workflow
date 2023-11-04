#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

source $script_dir/base.sh
source $script_dir/pr-body.sh
source $script_dir/multiselect.sh
source $script_dir/pr-jira.sh

jira_ticket=$1
if [ -z "$jira_ticket" ]; then
    read -p 'Jira ticket(It is optional when there is no ticket): ' jira_ticket
fi

read -p 'Issue desc(require): ' issue_desc
while [ -z "$issue_desc" ]; do
    read -p 'Issue desc(require): ' issue_desc
done

read -p 'Short description(optional): ' short_description

short_description=${short_description:-"Not yet"}

echo 'Types of changes:'
multiselect "false" result types_of_changes preselection

commit_title=${jira_ticket}': '${issue_desc}
pr_body=$(getPRbody result "${short_description}" $jira_ticket)
branch_name=${jira_ticket}--$(echo "$issue_desc" | sed 's/[^a-zA-Z0-9]/-/g')

if [ -z "${jira_ticket}" ]; then
    commit_title=$issue_desc
    branch_name=$(echo "$issue_desc" | sed 's/[^a-zA-Z0-9]/-/g')
fi

if [ -n "${GH_BRANCH_PREFIX}" ]; then
    branch_name=${GH_BRANCH_PREFIX}/${branch_name}
fi

# echo $branch_name
# echo $commit_title
# echo $pr_body

git checkout -b $branch_name
git add . && git commit -m "${commit_title}" && git push -u origin $branch_name
pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H $branch_name)

if [ -n "${jira_ticket}" ]; then
    jira_create "$jira_ticket" "$pr_url"
fi

echo $pr_url | pbcopy
echo -e $y Successfully copied $pr_url to clipboard

open $pr_url
