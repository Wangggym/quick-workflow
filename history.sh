#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/base.sh

write_history() {
    local pr_url=$1
    local jira_ticket=$2

    pr_id=$(echo "$pr_url" | grep -oE '[0-9]+$')
    new_history="${jira_ticket},${pr_id}"

    echo $new_history >>"${script_dir}/work-history.txt"

    echo -e $y Added a work history in work-history.txt: $new_history
}

read_history() {
    local pr_id=$1

    text=$(grep -E ",$pr_id$" "${script_dir}/work-history.txt")

    if [ -z "$text" ]; then
        return 0
    fi

    jira_ticket=$(echo "$text" | awk -F ',' '{print $1}')

    echo "$jira_ticket"
}
