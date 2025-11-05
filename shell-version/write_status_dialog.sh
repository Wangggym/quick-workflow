#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/base.sh
source $script_dir/jira-status.sh

jira_project=$1

if [ -z "$jira_project" ]; then
    echo -e $n Need input jira project name or jira ticket
    exit 0
fi

write_status_dialog_func "$jira_project"
