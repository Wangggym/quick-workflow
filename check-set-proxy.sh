#!/bin/bash

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

if [ "$http_proxy_enable" == "1" ] && [ "$https_proxy_enable" == "1" ] && [ "$socks_proxy_enable" == "1" ]; then
    export http_proxy="http://$http_proxy_address:$http_proxy_port" https_proxy="http://$https_proxy_address:$https_proxy_port" all_proxy="socks5://$socks_proxy_address:$socks_proxy_port"

    echo ✓ export http_proxy="http://$http_proxy_address:$http_proxy_port" https_proxy="http://$https_proxy_address:$https_proxy_port" all_proxy="socks5://$socks_proxy_address:$socks_proxy_port"
else
    unset http_proxy
    unset https_proxy
    unset all_proxy
    echo 'The proxy server is not available, please check'
    exit 1
fi

exit 0
