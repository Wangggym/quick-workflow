#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/base.sh
source $script_dir/multiselect.sh
source $script_dir/util.sh

get_jirs_status_list() {
    email=$(jira me)
    local json_data=$(curl -s "$JIRA_SERVICE_ADDRESS/rest/api/2/project/$1/statuses" \
        --user "$email:$JIRA_API_TOKEN")
    local options_string=$(echo "$json_data" | jq '.[0].statuses[].untranslatedName')
    echo $options_string
}

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
    local jira_ticket=$2
    local jira_project=$(get_jira_project_name "$jira_ticket")

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

read_status_pr_created() {
    read_status created-pr $1
}

read_status_pr_merged() {
    read_status merged-pr $1
}

write_status_dialog_func() {
    local jira_ticket=$1

    jira_project=$(get_jira_project_name "$jira_ticket")

    option_string=$(get_jirs_status_list "$jira_project")
    eval "result_options=($option_string)"

    echo "Select one of the following states to change when PR is ready or In progress:"
    multiselect "false" result_pr_create result_options defualts

    echo "Select one of the following states to change when PR is merged or Done:"
    multiselect "false" result_pr_merged result_options defualts

    idx=0
    for option in "${result_options[@]}"; do
        if [ "${result_pr_create[idx]}" = "true" ]; then
            result_pr_create="${result_options[idx]}"
        fi
        if [ "${result_pr_merged[idx]}" = "true" ]; then
            result_pr_merged="${result_options[idx]}"
        fi
        ((idx++))
    done

    write_status "$jira_project" "$result_pr_create" "$result_pr_merged"
}
