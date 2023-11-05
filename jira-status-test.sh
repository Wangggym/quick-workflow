#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/jira-status.sh

write_status CSV "In Progress" "Done"

read_status merged-pr CSV
read_status created-pr CSV

read_statsu_pr_merged CSV
read_statsu_pr_created CSV
