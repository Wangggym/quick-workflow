#!/bin/bash

read -p 'Jira ticket: ' jira_ticket

read -p 'issue name: ' issue_name

# read -p 'Types of changes: ' types_of_changes

read -p 'Short description: ' short_description

branch_name=${jira_ticket}--$(echo "$issue_name" | sed 's/ /-/g')
commit_title=${jira_ticket}': '${issue_name}

create_new_branch=$(git checkout -b "${branch_name}")
sleep 1
create_new_commit=$(git add . && git commit -m "${commit_title}" && git push)

echo "${create_new_branch}"
echo "${create_new_commit}"

pr_body="
# PR Ready

## Types of changes
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Refactoring (non-breaking change which does not change functionality)

#### Short description:

$short_description

#### Trello Issue Link:


#### Dependency
"

# must provide `--title` and `--body` (or `--fill` or `fill-first`) when not running interactively
pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H)

sleep 1

echo $pr_url

echo | $(jira issue comment add "${jira_ticket}" "${pr_url}")

# write history
script_dir="$(dirname "$0")"

pr_id=$(echo "$pr_url" | grep -oE '[0-9]+$')

echo "${jira_ticket},${pr_id}" >> "${script_dir}/work-history.txt"

# TODO change jira status
# TODO add Types of changes
# TODO open pr on windows
# TODO add jira assign to me