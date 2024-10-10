#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/history.sh
source $script_dir/jira-status.sh

# New function to extract Jira ticket ID from PR title
extract_jira_ticket_id() {
    local pr_title="$1"
    local jira_ticket_id=$(echo "$pr_title" | grep -oE '^[A-Z]+-[0-9]+')
    echo "$jira_ticket_id"
}

jira_create() {
    local jira_ticket=$1
    local pr_url=$2
    local status=$3
    local short_description=$4

    jira issue assign $jira_ticket $(jira me)
    jira issue move $jira_ticket "${status}"
    echo $pr_url | jira issue comment add $jira_ticket
    
    if [ -n "${short_description}" ]; then
        echo $short_description | jira issue comment add $jira_ticket
    fi

    # write history
    write_history $pr_url $jira_ticket
}

jira_merge() {
    local pr_id=$1

    # read history
    jira_ticket=$(read_history $pr_id)

    if [ -z "${jira_ticket}" ]; then
        # If no Jira ticket found in history, try to extract it from PR title
        pr_title=$(gh pr view $pr_id --json title --jq .title)
        jira_ticket=$(extract_jira_ticket_id "$pr_title")
        
        if [ -z "${jira_ticket}" ]; then
            echo "No Jira ticket found in history or PR title"
            return 0
        fi
    fi

    status=$(read_status_pr_merged $jira_ticket)

    jira issue move $jira_ticket "${status}"
}

