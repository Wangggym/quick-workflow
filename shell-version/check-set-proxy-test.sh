#!/bin/bash
script_dir="$(dirname "$0")"

echo $http_proxy
echo $https_proxy
echo $all_proxy

$script_dir/check-set-proxy.sh