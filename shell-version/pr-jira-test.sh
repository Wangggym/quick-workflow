#!/bin/bash

# Source the file containing the function we want to test
source ./pr-jira.sh

# Test function
test_extract_jira_ticket_id() {
    local test_case="$1"
    local expected="$2"
    local result=$(extract_jira_ticket_id "$test_case")
    
    if [ "$result" = "$expected" ]; then
        echo "PASS: '$test_case' -> '$result'"
    else
        echo "FAIL: '$test_case' -> '$result' (Expected: '$expected')"
    fi
}

# Run tests
echo "Running tests for extract_jira_ticket_id function:"
echo "------------------------------------------------"

test_extract_jira_ticket_id "IOSNAT-25723: Arrow is upward when I am not expanding" "IOSNAT-25723"
test_extract_jira_ticket_id "ANDROID-1234: Fix login bug" "ANDROID-1234"
test_extract_jira_ticket_id "WEB-56789: Implement new feature" "WEB-56789"
test_extract_jira_ticket_id "No Jira ticket in this title" ""
test_extract_jira_ticket_id "PROJ-123-456: Multiple numbers" "PROJ-123"
test_extract_jira_ticket_id "proj-123: Lowercase project" ""
test_extract_jira_ticket_id "PROJ123: No hyphen" ""
test_extract_jira_ticket_id "PROJ-ABC: Non-numeric ID" ""

echo "------------------------------------------------"
echo "Tests completed."