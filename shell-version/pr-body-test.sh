#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/pr-body.sh

jira_ticket=""
short_description='short description'
result=("true" "true" "true")

pr_body=$(getPRbody result "${short_description}" $jira_ticket)

echo $pr_body
