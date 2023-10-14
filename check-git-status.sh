#!/bin/bash

working_tree_clean="nothing to commit, working tree clean"

git_status=$(git status)

if [[ $git_status == *"${working_tree_clean}"* ]]; then
    echo $git_status
    exit 1
fi

exit 0
