#!/bin/bash

script_dir="$(dirname "$0")"

get_jira_project_name() {
    local jira_ticket=$1
    local jira_project="${jira_ticket%%-*}"

    if [ "$jira_project" = "$jira_ticket" ]; then
        echo ''
    else
        echo "$jira_project"
    fi
}
