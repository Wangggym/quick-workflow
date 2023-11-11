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

    pr_body="
# PR Ready\n

## Types of changes\n

$formatted_output

#### Short description:\n

$2\n

#### Trello Issue Link:\n

${JIRA_SERVICE_ADDRESS}/browse/$3\n

#### Dependency\n
"
    echo -e "$pr_body"
}
