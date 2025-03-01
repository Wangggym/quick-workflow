#!/bin/bash

# Get attachment URLs from JIRA issue
get_attachment_urls() {
    local issue_key="$1"
    local filter_pattern="$2"  # Optional filter pattern, defaults to "log."
    
    # Check JIRA auth
    if [ -z "$JIRA_API_TOKEN" ]; then
        echo "❌ Error: JIRA_API_TOKEN environment variable not set" >&2
        return 1
    fi

    echo "ℹ️  Getting attachments for $issue_key..." >&2

    # Save original output to temp file
    local temp_file=$(mktemp)
    jira issue view "$issue_key" --plain > "$temp_file"

    # Extract attachments section to temp file - stop at the next section header
    local attachments_file=$(mktemp)
    awk '/^[[:space:]]*## \*\*Attachments\*\*/ {p=1; next} /^[[:space:]]*## / {if(p) exit} p' "$temp_file" > "$attachments_file"

    # Process attachments with awk
    local attachments=$(awk -v filter="$filter_pattern" '
        BEGIN {
            current_file = ""
            current_url = ""
            in_url = 0
            processed_urls[0] = 0  # Initialize processed URLs array
        }
        
        # Match file number lines
        /^[[:space:]]*[0-9]+\.[[:space:]]/ {
            if (current_url != "" && (filter == "" || current_file ~ "^" filter)) {
                url_key = current_file "§" current_url
                if (!(url_key in processed_urls)) {
                    processed_urls[url_key] = 1
                    print url_key
                }
            }
            
            # Get new filename
            current_file = $0
            sub(/^[[:space:]]*[0-9]+\.[[:space:]]*/, "", current_file)
            gsub(/[[:space:]]*$/, "", current_file)
            current_url = ""
            in_url = 0
            next
        }
        
        # Collect URL fragments
        {
            if ($0 ~ /^[[:space:]]*http/) {
                if (current_url != "") {
                    # If current URL exists, start a new one
                    current_url = ""
                }
                in_url = 1
            }
            
            if (in_url) {
                line = $0
                gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
                current_url = current_url line
            }
        }
        
        END {
            if (current_url != "" && (filter == "" || current_file ~ "^" filter)) {
                url_key = current_file "§" current_url
                if (!(url_key in processed_urls)) {
                    processed_urls[url_key] = 1
                    print url_key
                }
            }
        }
    ' "$attachments_file")

    # Clean up temp files
    rm "$temp_file" "$attachments_file"

    if [ -z "$attachments" ]; then
        echo "❌ No attachments found" >&2
        return 1
    fi

    echo "$attachments"
    return 0
} 