#!/bin/bash

script_dir="$(dirname "$0")"

source $script_dir/base.sh

github="github.com"

if curl -IsSf $github -o /dev/null; then
    echo -e $y The github network is available
    exit 0
fi

echo -e $n Github Network error, please make sure network is able to use

exit 1
