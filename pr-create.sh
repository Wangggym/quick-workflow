#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

read -p 'Jira ticket: ' jira_ticket

read -p 'issue name: ' issue_name

# read -p 'Types of changes: ' types_of_changes

read -p 'Short description: ' short_description

branch_name=${jira_ticket}--$(echo "$issue_name" | sed 's/ /-/g')
commit_title=${jira_ticket}': '${issue_name}

create_new_branch=$(git checkout -b "${branch_name}")
create_new_commit=$(git add . && git commit -m "${commit_title}" && git push)

pr_body="
# PR Ready

## Types of changes
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Refactoring (non-breaking change which does not change functionality)

#### Short description:

$short_description

#### Trello Issue Link:

${JIRA_SERVICE_ADDRESS}/${jira_ticket}

#### Dependency
"

pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H "${create_new_branch}")


if [[ $jira_ticket == *BSF* ]]; then
    status="In Review"
else
    status="UNDER Review"
fi
echo $(jira issue assign $jira_ticket $(jira me))

echo $(jira issue move $jira_ticket "${status}")

echo $pr_url | jira issue comment add $jira_ticket

# write history
pr_id=$(echo "$pr_url" | grep -oE '[0-9]+$')

echo "${jira_ticket},${pr_id}" >>"${script_dir}/work-history.txt"

echo $pr_url
open $pr_url


# TODO add Types of changes
