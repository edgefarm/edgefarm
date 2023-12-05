#!/bin/bash

iptables -D FORWARD -i wt0 -o edge0 -p tcp --dport 1:65535 -j ACCEPT
iptables -D FORWARD -i edge0 -o wt0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -t nat -D PREROUTING -i wt0 -p tcp --dport 1:65535 -j DNAT --to-destination 192.168.168.1
iptables -t nat -D POSTROUTING -o wt0 -j MASQUERADE
