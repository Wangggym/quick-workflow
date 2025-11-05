#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/history.sh

# should be written history
write_history 9323 ABC-134

# should be readed history
jira_ticket=$(read_history 9323)
if [ -n "$jira_ticket" ]; then
    echo Readed $jira_ticket
fi

# shold not be readed history
jira_ticket=$(read_history 1234)
if [ -z "$jira_ticket" ]; then
    echo Cannot read
fi

jira_ticket=$(read_history 2)
if [ -z "$jira_ticket" ]; then
    echo Cannot read
fi

jira_ticket=$(read_history 7)
if [ -z "$jira_ticket" ]; then
    echo Cannot read
fi
