#!/bin/bash

# Check if both arguments are provided
if [ $# -lt 2 ]; then
    echo "Usage: $0 <log_file> <search_term>"
    exit 1
fi

# Get the last argument as search_term
search_term="${@: -1}"

# Remove the last argument, leaving only file paths
set -- "${@:1:$(($#-1))}"

# Use awk to search for the term and extract relevant API URLs and IDs
awk -v term="$search_term" '
    function print_result() {
        if (block_url != "" && block_id != "" && found_in_current_block) {
            if (!(block_id in printed)) {
                printf "URL: %s, ID: %s\n", block_url, block_id
                printed[block_id] = 1
            }
        }
    }
    
    /^💡/ {
        # New log entry starts
        found_in_current_block = 0
        if ($0 ~ /#[0-9]+/) {
            match($0, /#([0-9]+)/)
            block_id = substr($0, RSTART+1, RLENGTH-1)
            
            # Try to extract URL - first check for HTTP method
            if ($0 ~ /(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)/) {
                match($0, /(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)/)
                method_pos = RSTART + RLENGTH
                rest = substr($0, method_pos)
                if (match(rest, /[[:space:]]+https?:\/\/[^[:space:]]+/)) {
                    block_url = substr(rest, RSTART+1, RLENGTH-1)
                    gsub(/["\047\s,}]$/, "", block_url)
                }
            } 
            # If no HTTP method found, try to find URL in response line
            else if ($0 ~ /[0-9]+ https?:\/\//) {
                match($0, /[0-9]+ (https?:\/\/[^[:space:]]+)/)
                if (RSTART) {
                    block_url = substr($0, RSTART + RLENGTH - match($0, /https:\/\/[^[:space:]]+/), match($0, /https:\/\/[^[:space:]]+/))
                    gsub(/["\047\s,}]$/, "", block_url)
                }
            }
        }
    }
    
    tolower($0) ~ tolower(term) {
        found_in_current_block = 1
        print_result()
    }
    
    /^$/ {
        # Empty line marks end of a block
        block_url = ""
        block_id = ""
        found_in_current_block = 0
    }
' "$@"
