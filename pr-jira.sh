#!/bin/bash
script_dir="$(dirname "$0")"

jira_create() {
    local jira_ticket=$1
    local pr_url=$2
    local status=$3

    jira issue assign $jira_ticket $(jira me)
    jira issue move $jira_ticket "${status}"
    echo $pr_url | jira issue comment add $jira_ticket

    # write history
    pr_id=$(echo "$pr_url" | grep -oE '[0-9]+$')
    echo "${jira_ticket},${pr_id}" >>"${script_dir}/work-history.txt"
}

jira_merge() {
    local pr_id=$1

    # read history
    text=$(grep $pr_id "${script_dir}/work-history.txt")
    jira_ticket=$(echo "$text" | awk -F ',' '{print $1}')

    echo $jira_ticket

    if [ -z "${jira_ticket}" ]; then
        return 0
    fi

    if [[ $jira_ticket == *BSF* ]]; then
        status="FE fixed"
    else
        status="IN QA"
    fi

    jira issue move $jira_ticket "${status}"
}
