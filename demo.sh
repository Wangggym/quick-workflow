#!/bin/bash

script_dir="$(dirname "$0")"

# echo "BSF-123,9283" >> "${script_dir}/work-history.txt"

# echo -e "${script_dir}/work-history.txt" | grep "9283"

text=$(grep "9283" "${script_dir}/work-history.txt")

echo $text


var1=$(echo "$text" | awk -F ',' '{print $1}')
var2=$(echo "$text" | awk -F ',' '{print $2}')

echo "var1: $var1"
echo "var2: $var2"