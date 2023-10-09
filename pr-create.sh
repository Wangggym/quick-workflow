#!/bin/bash

read -p 'Jira ticket: ' jira_ticket

read -p 'issue name: ' issue_name

# read -p 'Types of changes: ' types_of_changes

# read -p 'Short description: ' short_description

branch_name=${jira_ticket}--$(echo "$issue_name" | sed 's/ /-/g')
commit_title=${jira_ticket}': '${issue_name}

create_new_branch=$(git checkout -b "${branch_name}")
create_new_commit=$(git add . && git commit -m "${commit_title}" && git push)

echo "${create_new_branch}"
echo "${create_new_commit}"

pr_body='fweoiewf12121212'

# must provide `--title` and `--body` (or `--fill` or `fill-first`) when not running interactively
pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H $branch_name)

echo "${pr_url}"
