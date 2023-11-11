#!/bin/bash

script_dir="$(dirname "$0")"

source $script_dir/util.sh

# should be get project name
echo $(get_jira_project_name "ABC-123")

# should not be get project name
echo $(get_jira_project_name "ABC/123")
echo $(get_jira_project_name "ABC1123")

