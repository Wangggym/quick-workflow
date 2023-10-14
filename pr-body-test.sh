#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/pr-body.sh

jira_ticket='jira_ticket'
short_description='short_description'
result=("true" "true" "true")

pr_body=$(getPRbody $jira_ticket $short_description result)

echo $pr_body
