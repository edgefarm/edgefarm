# Join a physical edge node running vanilla Ubuntu 22.04 to the local EdgeFarm cluster

## Prerequisites

You need a running EdgeFarm cluster. See [Create a local EdgeFarm cluster for testing](create-local-cluster.md) for instructions.

This How-To shows how to add a Raspberry Pi 4 with a vanilla Ubuntu 22.04 as an Edge Node. 
If you have a different hardware or OS your mileage may vary but the general steps should be the same.

The device should be running

- Ubuntu 22.04 LTS
- systemd
- Docker
- udev
- cgroupv2

## Pre-register the Edge Node using `local-up`

Pre-register the node using `local-up` to create the corresponding resources for the edge node. 

By default, the nodename is the hostname of the edge node. If you didn't define a different TTL using the `--ttl` argument for your token, you have 24 hours to join the node to the cluster. After that, the token expires and you won't be able to join using this token.

See this example to pre-register a node with the name `eagle`:

```bash
$ local-up node join --name eagle
I1206 12:42:34.430305  958575 nodeJoin.go:121] Adding empty node resource for eagle
I1206 12:42:34.430418  958575 nodeJoin.go:129] Adding nodepool for node eagle
Here is some information you need to join a physical edge node to this cluster.
Ensure that the /etc/hosts file contain the following entry: # (1)!
192.168.1.46 edgefarm-control-plane 

Use this token m7i0gd.vduhnle5x4dzs3lo to join the cluster. You have 1 day to join the cluster before this token expires. 
# (2)!

If you experience any problems, please consult the documentation at 
https://edgefarm.github.io/edgefarm/ or file an issue at https://github.com/edgefarm/edgefarm/issues/new?template=question.md
```

1.  '192.168.1.46' is the IP address of the local cluster. It mus be connected to the physical node. This might be different depending on your network setup.
2.  You need this example token 'm7i0gd.vduhnle5x4dzs3lo' later to join the node. Every node gets its own token.

## Edge Node preparations

Install Ubuntu Server 22.04 on your target machine. It doesn't matter if you use a Virtual Machine, a Raspberry Pi or any other physical machine.

Make sure that you have SSH access to the Pi and that you can log in as root.

Install required packages
```bash
sudo apt update
sudo apt-get install curl -y
```
  
## Run the Edge Node prepare script

Download the `ubuntu-22.04-edge-node-config.tar.gz` from the [EdgeFarm releases page](https://github.com/edgefarm/edgefarm/releases) and extract it on your edge node.

```bash
$ tar xvfz ubuntu-22.04-edge-node-config.tar.gz
```

Run the `prepare.sh` script as root. This step takes care of the following tasks:

- install Docker, socat and conntrack
- install netbird
- disable swap
- setup some udev rules and scripts

```bash
./prepare.sh
```

Once done, connect netbird using a netbird setup-key.
```bash
netbird up -k <your token>
```

You should see the interface `wt0` using `ip a` now.

```bash
$ ip addr show dev wt0
13: wt0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1280 qdisc noqueue state UNKNOWN group default qlen 1000
    link/none # (1)!
    inet 100.127.123.145/16 brd 100.127.255.255 scope global wt0 
       valid_lft forever preferred_lft forever
```

1.  The IP address is 100.127.123.145

## Join the Edge Node to the EdgeFarm cluster

Run the `install.sh` script as root to join the Edge Node to the EdgeFarm cluster.
The bootstrap-token to enter is the one you got from the `local-up` command above (`m7i0gd.vduhnle5x4dzs3lo`). 

```console
./install.sh --address <IP:port> --token <bootstrap-token> --node-ip $(cat /usr/local/etc/wt0.ip) --join
```

??? "Raspberry Pi 4 notes"

    For the Raspberry Pi you need to install the package `linux-modules-extra-raspi` and enable some boot options to make it work with Kubernetes.

    ```bash
    apt update
    apt install linux-modules-extra-raspi
    sed -i '$s/$/ cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory/' /boot/firmware/cmdline.txt
    ```

## Verify the Edge Node

Verify that the Edge Node is ready to use. Let's say the hostname of the node is `eagle`.

```bash
$ kubectl get nodes | grep eagle
NAME     STATUS   ROLES    AGE   VERSION
eagle    Ready    <none>   1M    v1.22.17
```