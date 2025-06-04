#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/base.sh



function getPRbody {
    formatted_output=""
    array_length=${#types_of_changes[@]}

    eval "local -a result=(\"\${$1[@]}\")"

    for ((i = 0; i < array_length; i++)); do
        if [ "${result[i]}" == "true" ]; then
            formatted_output+="- [x] ${types_of_changes[i]}\n"
        else
            formatted_output+="- [ ] ${types_of_changes[i]}\n"
        fi
    done

    pr_body="\n# PR Ready\n\n## Types of changes\n\n$formatted_output\n"

    if [ -n "$2" ]; then
        pr_body+="\n#### Short description:\n\n$2\n"
    fi

    if [ -n "$3" ]; then
        pr_body+="\n#### Jira Link:\n\n${JIRA_SERVICE_ADDRESS}/browse/$3\n"
    fi

    echo -e "$pr_body"
}
