#!/bin/bash

# Check if both arguments are provided
if [ $# -lt 2 ]; then
    echo "Usage: $0 <log_file> <request_id> [jira_id]"
    exit 1
fi

# Get the last argument as request_id if no jira_id provided, otherwise second to last
if [ $# -eq 3 ]; then
    request_id="$2"
    jira_id="$3"
else
    request_id="${@: -1}"
    jira_id=""
fi

# Remove the last argument(s), leaving only file paths
if [ $# -eq 3 ]; then
    set -- "$1"
else
    set -- "${@:1:$(($#-1))}"
fi

# Get the domain (log file name) and path from the first argument
if [ -n "$jira_id" ]; then
    domain="${JIRA_SERVICE_ADDRESS}/browse/$jira_id"
else
    domain=$(basename "$1")
fi

# Extract name from the URL in the log line
name=$(awk -v rid="$request_id" '
    $0 ~ "#" rid {
        for (i=1; i<=NF; i++) {
            if ($i ~ /^https?:.*/) {
                # First remove domain and any prefix paths
                sub(/^https?:\/\/[^\/]+\/[^\/]+\/[^\/]+\//, "", $i)
                # Then get only the last two path segments if they exist
                if (match($i, /[^\/]+\/[^\/]+$/)) {
                    print substr($i, RSTART, RLENGTH)
                    exit
                } else if (match($i, /[^\/]+$/)) {
                    # If only one segment exists, use that
                    print substr($i, RSTART, RLENGTH)
                    exit
                }
            }
        }
    }
' "$@")

# If name is empty, use a default value
if [ -z "$name" ]; then
    name="unknown"
else
    name="#${request_id} ${name}"
fi

# Detect the operating system and choose the appropriate clipboard command
if [[ "$OSTYPE" == "darwin"* ]]; then
    CLIP_CMD="pbcopy"
elif command -v xclip >/dev/null 2>&1; then
    CLIP_CMD="xclip -selection clipboard"
elif command -v xsel >/dev/null 2>&1; then
    CLIP_CMD="xsel --clipboard --input"
else
    CLIP_CMD="tee"  # If no clipboard command available, just output to stdout
fi

# Create a temporary file to store the response
temp_file=$(mktemp)

# Use awk to find and extract the response content
awk -v rid="$request_id" '
    /response:/ { 
        # Look for exact match of request ID
        if (prev ~ "#" rid) {
            # Print everything after "response: " until the next empty line
            gsub(/^.*response: /, "")
            p = 1
            print
            next
        }
    }
    p && NF {
        # Continue printing until we hit an empty line
        print
        next    
    }
    p && !NF { 
        # Stop printing when we hit an empty line
        exit 
    }
    { 
        prev = $0 
    }
' "$@" | tee >(eval "$CLIP_CMD") > "$temp_file"

# Check if temp_file is empty
if [ ! -s "$temp_file" ]; then
    echo "Error: No response content found for request ID: $request_id"
    rm "$temp_file"
    exit 1
fi

# Create JSON payload file
json_file=$(mktemp)

# Generate timestamp
timestamp=$(date "+%m/%d/%Y, %I:%M:%S %p")

# Create a temporary file for the Node.js script
node_script_file=$(mktemp)

# Create the Node.js script that will process the data
cat > "$node_script_file" << 'EOF'
const fs = require('fs');
const inputFile = process.argv[2];
const data = fs.readFileSync(inputFile, 'utf8');
const domain = process.argv[3];
const name = process.argv[4];
const timestamp = process.argv[5];

console.log(JSON.stringify({
    encodedKey: '',
    data: data,
    combineLine: '',
    separator: '',
    domain: domain,
    name: name,
    timestamp: timestamp
}, null, 2));
EOF

# Execute Node.js script with the data file
node "$node_script_file" "$temp_file" "$domain" "$name" "$timestamp" > "$json_file"

# Clean up the temporary files
rm "$node_script_file"
rm "$temp_file"

# Debug output
# echo "Sending request with payload:"
# cat "$json_file"

# Send to Streamock and capture the response
response=$(curl -s -w "\n%{http_code}" 'http://localhost:3001/api/submit' \
  -H 'Content-Type: application/json' \
  -H 'Connection: keep-alive' \
  -d "@$json_file")

# Clean up the temporary JSON file
rm "$json_file"

# Extract status code and response body
http_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | sed '$d')

# Check if the request was successful
if [ -n "$http_code" ] && [ "$http_code" -eq 200 ]; then
    echo "Result has been copied to clipboard and successfully sent to Streamock!"
else
    echo "Result has been copied to clipboard, but Error sending to Streamock. Status code: ${http_code:-'unknown'}"
    echo "Response: ${response_body:-'no response'}"
    exit 1
fi
