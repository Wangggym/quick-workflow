#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/history.sh

jira_create() {
    local jira_ticket=$1
    local pr_url=$2

    if [[ $jira_ticket == *BSF* ]]; then
        status="In Review"
    fi
    if [[ $jira_ticket == *STUD* ]]; then
        status="UNDER Review"
    fi
    if [[ $jira_ticket == *IRB* ]]; then
        status="IN Progress"
    fi

    jira issue assign $jira_ticket $(jira me)
    jira issue move $jira_ticket "${status}"
    echo $pr_url | jira issue comment add $jira_ticket

    # write history
    write_history $pr_url $jira_ticket
}

jira_merge() {
    local pr_id=$1

    # read history
    jira_ticket=$(read_history $pr_id)

    if [ -z "${jira_ticket}" ]; then
        return 0
    fi

    if [[ $jira_ticket == *BSF* ]]; then
        status="FE fixed"
    fi
    if [[ $jira_ticket == *STUD* ]]; then
        status="IN QA"
    fi
    if [[ $jira_ticket == *IRB* ]]; then
        status="DONE"
    fi

    jira issue move $jira_ticket "${status}"
}
