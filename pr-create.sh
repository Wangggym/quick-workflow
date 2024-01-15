#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

source $script_dir/base.sh
source $script_dir/pr-body.sh
source $script_dir/multiselect.sh
source $script_dir/pr-jira.sh
source $script_dir/jira-status.sh

jira_ticket=$1
if [ -z "$jira_ticket" ]; then
    read -p 'Jira ticket (It is optional when there is no ticket): ' jira_ticket
fi
if [ -n "${jira_ticket}" ]; then
    status=$(read_status_pr_created $jira_ticket)
    if [ -z "$status" ]; then
        write_status_dialog_func "$jira_ticket"
    fi
    status=$(read_status_pr_created $jira_ticket)
fi

if [ -n "${jira_ticket}" ]; then
    issue_json=$(aiwflow issue-desc $jira_ticket)
    issue_desc=$(echo "$issue_json" | jq -r '.issue_desc')
    need_translate=$(echo "$issue_json" | jq -r '.need_translate')
    translated_desc=$(echo "$issue_json" | jq -r '.translated_desc')

    if [ "$issue_desc" != "null" ]; then
        echo -e $y 'PR title: '$issue_desc
        if [ "$need_translate" = "true" ]; then
            if [ "$translated_desc" != "null" ]; then
                issue_desc=$translated_desc
                echo -e $y 'Chat GPT translated: '$translated_desc
            else

                echo -e $n 'Chat GPT translated failed, but it dose not affect the use.'
            fi
        fi
    else
        issue_desc=""
    fi
fi

while [ -z "$issue_desc" ]; do
    read -p 'PR title and git branch name (require) : ' issue_desc
done

read -p 'Short description (optional): ' short_description

github_short_description=${short_description:-"Not yet"}

echo 'Types of changes:'
multiselect "true" result types_of_changes preselection

commit_title=${jira_ticket}': '${issue_desc}
pr_body=$(getPRbody result "${github_short_description}" $jira_ticket)
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
git add --all
git commit -m "${commit_title}"
git push -u origin $branch_name

pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H $branch_name)

if [ -n "${jira_ticket}" ]; then
    jira_create "$jira_ticket" "$pr_url" "$status" "$short_description"
fi

echo $pr_url | pbcopy
echo -e $y Successfully copied $pr_url to clipboard

# Clearly show users the copied information, sleep 1 second
sleep 1

open $pr_url
