#!/bin/bash

# Name of the dummy interface
dummy_interface="edge0"
mac_address_file="/usr/local/etc/edge0.mac"
ip_address="192.168.168.1"
# Create or read MAC address from the file
if [ -f "$mac_address_file" ]; then
  mac_address=$(cat "$mac_address_file")
else
  mac_address=$(echo $FQDN|md5sum|sed 's/^\(..\)\(..\)\(..\)\(..\)\(..\).*$/02:\1:\2:\3:\4:\5/')
  echo "$mac_address" > "$mac_address_file"
fi

# Create the dummy interface and assign the IP address and subnet
ip link add name edge0 type dummy
ip link set dev edge0 address "$mac_address"
ip addr add "$ip_address/24" dev edge0
ip link set dev edge0 up

echo 1 > /proc/sys/net/ipv4/ip_forward
exit 0
