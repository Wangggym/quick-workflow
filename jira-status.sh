#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/base.sh

write_status() {
    local jira_project=$1
    local jira_status_created_pr=$2
    local jira_status_merged_pr=$3

    new_jira_status_created_pr="created-pr,${jira_project},${jira_status_created_pr}"
    new_jira_status_merged_pr="merged-pr,${jira_project},${jira_status_merged_pr}"

    echo "$jira_project": >>"${script_dir}/jira-status.txt"
    echo $new_jira_status_created_pr >>"${script_dir}/jira-status.txt"
    echo $new_jira_status_merged_pr >>"${script_dir}/jira-status.txt"
    echo -e "\n" >>"${script_dir}/jira-status.txt"

    echo -e $y Added jira status in jira-status.txt
}

check_jira_status_file() {
    local pr_status=$1
    local jira_project=$2

    status=$(grep -E "^$pr_status,$jira_project" "${script_dir}/jira-status.txt")
    echo $status
}

read_status() {
    status=$(check_jira_status_file $1 $2)

    if [ -z "$status" ]; then
        return 0
    fi

    jira_status=$(echo "$status" | awk -F ',' '{print $3}')

    echo "$jira_status"
}

read_statsu_pr_created() {
    read_status created-pr $1
}

read_statsu_pr_merged() {
    read_status merged-pr $1
}
