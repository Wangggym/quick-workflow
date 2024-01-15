#!/bin/bash
script_dir="$(dirname "$0")"
source $script_dir/base.sh

jira_ticket=$1

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

echo $issue_desc

while [ -z "$issue_desc" ]; do
    read -p 'PR title and git branch name (require) : ' issue_desc
done
