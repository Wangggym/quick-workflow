#!/bin/bash

# Show usage if no arguments provided
if [ $# -lt 2 ]; then
    echo "Usage: qk <JIRA-ID> [-d|-f|-s]"
    echo "  -d: Download logs (qklogs)"
    echo "  -f: Find request by ID (qkfind)"
    echo "  -s: Search in logs (qksearch)"
    exit 1
fi

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
QKLOGS_PATH="$SCRIPT_DIR/qklogs/qklogs.sh"
QKFIND_PATH="$SCRIPT_DIR/qkfind.sh"
QKSEARCH_PATH="$SCRIPT_DIR/qksearch.sh"

JIRA_ID="$1"
ACTION="$2"
LOGS_DIR="$HOME/Downloads/logs_${JIRA_ID}"
LOG_FILE="$LOGS_DIR/merged/flutter-api.log"

# Function to check if logs exist
check_logs() {
    if [ ! -f "$LOG_FILE" ]; then
        echo "❌ Log file not found at: $LOG_FILE"
        echo "💡 Try downloading logs first with: qk $JIRA_ID -d"
        exit 1
    fi
}

case "$ACTION" in
    "-d")
        # Call qklogs with full path
        "$QKLOGS_PATH" "$JIRA_ID"
        ;;
    
    "-f")
        check_logs
        # Use third argument if provided, otherwise prompt
        if [ -n "$3" ]; then
            REQUEST_ID="$3"
        else
            echo -n "Enter request ID to find: "
            read REQUEST_ID
        fi
        "$QKFIND_PATH" "$LOG_FILE" "$REQUEST_ID"
        ;;
    
    "-s")
        check_logs
        # Use third argument if provided, otherwise prompt
        if [ -n "$3" ]; then
            SEARCH_TERM="$3"
        else
            echo -n "Enter search term: "
            read SEARCH_TERM
        fi
        "$QKSEARCH_PATH" "$LOG_FILE" "$SEARCH_TERM"
        ;;
    
    *)
        echo "❌ Invalid action: $ACTION"
        echo "Usage: qk <JIRA-ID> [-d|-f|-s]"
        exit 1
        ;;
esac 