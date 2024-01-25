# Join a physical edge node running vanilla Ubuntu 22.04 to the local EdgeFarm cluster

## Prerequisites

You need a running and VPN enabled EdgeFarm cluster. See [Create a local EdgeFarm cluster for testing](create-local-cluster.md) for instructions.

This How-To shows how to add a Raspberry Pi 4 with a vanilla Ubuntu 22.04 as an Edge Node. 
If you have a different hardware or OS your mileage may vary but the general steps should be adaptable.

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
Here is some information you need to join a physical edge node to this cluster.

VPN:
Unless you already connected the physical node to netbird.io VPN, you need to connect it to the VPN first.

Use can use this setup-key 78A12F38-7E48-4068-97D5-8172E3017C58  to connect to netbird.io VPN. # (1)!f

Kubernetes:
Ensure that the /etc/hosts file on your physical edge node contains the following entry:
100.127.213.101 edgefarm-control-plane # (2)!

Use this token ny32as.wrzxg9fzn3r5tjsj to join the cluster. You have 1 day to join the cluster before this token expires. # (3)!

If you experience any problems, please consult the documentation at 
https://edgefarm.github.io/edgefarm/ or file an issue at https://github.com/edgefarm/edgefarm/issues/new?template=question.md

```

1.  Keep the netbird setup-key for connecting the node to the VPN.
2.  Add this entry in /etc/hosts. '100.127.213.101' is the VPN IP address of the control-plane node of your local cluster.
3.  Keep this token for later to join the node. Every node gets its own token.

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

## /etc/hosts entry

Modify the `/etc/hosts` file on your edge node the way the `local-up node join`'s output told you. 

```bash
<IP ADDRESS> edgefarm-control-plane # (1)!
```

1.  In the example above, the IP address is `100.127.213.101`. So your entry must be `100.127.213.101 edgefarm-control-plane`.

## Join the Edge Node to the EdgeFarm cluster

Run the `install.sh` script as root to join the Edge Node to the EdgeFarm cluster.
The bootstrap-token to enter is the one you got from the `local-up` command above (`ny32as.wrzxg9fzn3r5tjsj`). 

```console
./install.sh --address edgefarm-control-plane:6443 --token <bootstrap-token> --node-ip $(cat /usr/local/etc/wt0.ip) --join
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
$ KUBECONFIG=~/.edgefarm-local-up/kubeconfig kubectl get nodes | grep eagle
NAME     STATUS   ROLES    AGE   VERSION
eagle    Ready    <none>   1M    v1.22.17
```
