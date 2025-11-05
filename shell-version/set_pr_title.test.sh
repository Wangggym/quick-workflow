#!/bin/bash
script_dir="$(dirname "$0")"
source $script_dir/base.sh

jira_ticket=$1

if [ -n "${jira_ticket}" ]; then
    issue_desc=$(aiwflow pr-create $jira_ticket)
fi

echo  -e $y 'PR title: ' $issue_desc
