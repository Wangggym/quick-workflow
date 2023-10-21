#!/bin/bash

github="github.com"

if curl -IsSf $github -o /dev/null && curl -IsSf $JIRA_SERVICE_ADDRESS -o /dev/null; then
    echo ✓ The network is available
    exit 0
fi

echo "Network error, please make sure network is able to use"

exit 1
