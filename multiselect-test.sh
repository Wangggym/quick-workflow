#!/bin/bash
script_dir="$(dirname "$0")"

source $script_dir/multiselect.sh

types_of_changes=("Option 1" "Option 2" "Option 3")
preselection=("false" "true" "false")

multiselect "false" result types_of_changes preselection

idx=0
for option in "${types_of_changes[@]}"; do
    echo -e "$option\t=> ${result[idx]}"
    ((idx++))
done
