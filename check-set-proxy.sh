#!/bin/bash

script_dir="$(dirname "$0")"

source $script_dir/base.sh


proxy_info=$(scutil --proxy)

http_proxy_enable=$(echo "$proxy_info" | awk -F' : ' '/HTTPEnable/ {print $2}')
http_proxy_address=$(echo "$proxy_info" | awk -F' : ' '/HTTPProxy/ {print $2}')
http_proxy_port=$(echo "$proxy_info" | awk -F' : ' '/HTTPPort/ {print $2}')
https_proxy_enable=$(echo "$proxy_info" | awk -F' : ' '/HTTPSEnable/ {print $2}')
https_proxy_address=$(echo "$proxy_info" | awk -F' : ' '/HTTPSProxy/ {print $2}')
https_proxy_port=$(echo "$proxy_info" | awk -F' : ' '/HTTPSPort/ {print $2}')
socks_proxy_enable=$(echo "$proxy_info" | awk -F' : ' '/SOCKSEnable/ {print $2}')
socks_proxy_address=$(echo "$proxy_info" | awk -F' : ' '/SOCKSProxy/ {print $2}')
socks_proxy_port=$(echo "$proxy_info" | awk -F' : ' '/SOCKSPort/ {print $2}')

match_http_proxy="http://$http_proxy_address:$http_proxy_port"
match_https_proxy="http://$https_proxy_address:$https_proxy_port"
match_all_proxy="socks5://$socks_proxy_address:$socks_proxy_port"

copy_proxy="export http_proxy=http://$http_proxy_address:$http_proxy_port https_proxy=http://$https_proxy_address:$https_proxy_port all_proxy=socks5://$socks_proxy_address:$socks_proxy_port"

if [ "$http_proxy" == "$match_http_proxy" ] && [ "$https_proxy" == "$match_https_proxy" ] && [ "$all_proxy" == "$match_all_proxy" ]; then
    echo -e $y The proxy server is available
    exit 0
else
    echo $copy_proxy | pbcopy
    echo -e $n The proxy service is not set on the command line, which may cause the operation to fail.
    echo Copied proxy command to clipboard, you can use it or set it by yourself
    exit 1
fi
