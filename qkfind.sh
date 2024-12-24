#!/bin/bash

# Check if both arguments are provided
if [ $# -lt 2 ]; then
    echo "Usage: $0 <log_file> <request_id>"
    exit 1
fi

# Get the last argument as request_id
request_id="${@: -1}"
# Remove the last argument, leaving only file paths
set -- "${@:1:$(($#-1))}"

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

# Use awk to find and extract the response content and copy to clipboard
awk -v rid="$request_id" '
    /response:/ { 
        if (prev ~ rid) {
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
' "$@" | eval "$CLIP_CMD"

echo "Result has been copied to clipboard!"
