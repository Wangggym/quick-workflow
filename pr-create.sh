#!/bin/bash
script_dir="$(dirname "$0")"

if ! "$script_dir/check.sh"; then
    exit 1
fi

if ! "$script_dir/check-pre-commit.sh"; then
    exit 1
fi

source $script_dir/base.sh
source $script_dir/pr-body.sh
source $script_dir/multiselect.sh
source $script_dir/pr-jira.sh
source $script_dir/jira-status.sh
source $script_dir/generate-branch-name/generate-branch-name.sh
source $script_dir/check-code-fix.sh

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

if [ -n "${jira_ticket}" ] && [ -n "${OPENAI_KEY}" ]; then
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
                echo -e $n 'Chat GPT translation failed. Please make sure you have installed aiwflow correctly (https://github.com/Wangggym/quick-workflow/blob/master/2.0.md).'
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

read -p "Excution make fix? (y/n): " answer

case $answer in
    [Yy][Ee][Ss]|[Yy])
        check_code_fix
        ;;
    *)
        ;;
esac

echo 'Types of changes:'
multiselect "true" result types_of_changes preselection

commit_title=${jira_ticket}': '${issue_desc}
pr_body=$(getPRbody result "${github_short_description}" $jira_ticket)
branch_name=${jira_ticket}--$(echo "$issue_desc" | sed 's/[^a-zA-Z0-9]/-/g')

if [ -z "${jira_ticket}" ]; then
    commit_title=$issue_desc
    branch_name=$(echo "$issue_desc" | sed 's/[^a-zA-Z0-9]/-/g')
fi

if [[ -n "${BRAIN_AI_KEY}" && -z "${OPENAI_KEY}" ]]; then
    start_time=$(date +%s.%N)
    generate_id=$(echo -n $commit_title | md5sum | awk '{print $1}')
    echo -e "Start fetch branch name from AI with ID: $generate_id"
    generate_branch_name "$commit_title" "$BRAIN_AI_KEY" "$generate_id"

    result=$(cat /tmp/branch_name_$generate_id.txt)
    echo -e $y Fetch branch name from AI success $result
    rm -f /tmp/branch_name_$generate_id.txt
    if [ -n "$result" ]; then
        branch_name=$result
    fi

    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    echo "Fetch branch name cost $duration seconds"
fi

if [ -n "${GH_BRANCH_PREFIX}" ]; then
    branch_name=${GH_BRANCH_PREFIX}/${branch_name}
fi

# echo $branch_name
# echo $commit_title
# echo $pr_body
git checkout -b $branch_name

# We have checked the commit at `check-pre-commit`, so we do not need to verify it again.
git commit -m "${commit_title}" --no-verify
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
