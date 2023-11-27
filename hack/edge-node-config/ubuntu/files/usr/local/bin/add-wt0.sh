#!/bin/bash

ip_file="/usr/local/etc/wt0.ip"
wt0_info=$(ip addr show dev wt0 | grep -oE 'inet [0-9.]+' | awk '{print $2}')
echo $wt0_info > $ip_file

iptables -A FORWARD -i wt0 -o edge0 -p tcp --dport 1:65535 -j ACCEPT
iptables -A FORWARD -i edge0 -o wt0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -t nat -A PREROUTING -i wt0 -p tcp --dport 1:65535 -j DNAT --to-destination 192.168.168.1
iptables -t nat -A POSTROUTING -o wt0 -j MASQUERADE
