# Join a physical edge node running vanilla Ubuntu 22.04 to an EdgeFarm cluster

First, select the type of EdgeFarm cluster you have:

=== "Local EdgeFarm cluster"

    You need a running and VPN enabled EdgeFarm cluster. See [Create a local EdgeFarm cluster for testing](create-local-cluster.md) for instructions.


=== "Cloud Cluster (e.g. Hetzner)"
     
    You need a running and Cloud driven EdgeFarm cluster, e.g. running on Hetzner Cloud. See [Create a EdgeFarm cluster that runs on Hetzner Cloud](create-hetzner-cluster.md) for instructions.

!!! note  "Raspberry Pi 4 notes"

    For the Raspberry Pi you need to install the package `linux-modules-extra-raspi` and enable some boot options to make it work with Kubernetes.

    ```bash
    apt update
    apt install linux-modules-extra-raspi
    sed -i '$s/$/ cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory/' /boot/firmware/cmdline.txt
    ```

## Prerequisites
    
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


=== "Local EdgeFarm cluster"

    ```bash
    $ local-up node join --name eagle
    Here is some information you need to join a edge node to this cluster.
    
    VPN:
    If you haven't already linked the node to the netbird.io VPN, you must establish the connection to the VPN beforehand.
    
    Use can use this setup-key 375EE66D-E647-4B75-A9FC-########### to connect to netbird.io VPN. 
    # (1)!
    
    Kubernetes:
    Ensure that the /etc/hosts file on your physical edge node contains the following entry:
    123.123.123.123 edgefarm-control-plane 
    # (2)!
    
    Use this token 123456.0wm9r81m6dlozk30 to join the cluster. You have 1 day to join the cluster before this token expires. 
    # (3)!
    
    If you experience any problems, please consult the documentation at 
    https://edgefarm.github.io/edgefarm/ or file an issue at https://github.com/edgefarm/edgefarm/issues/new?template=question.md
    ```

    1.  Keep the netbird setup-key for connecting the edge node to the VPN.
    2.  Add this entry in /etc/hosts. '100.127.213.101' is the VPN IP address of the control-plane node of your local cluster.
    3.  Keep this token for later to join the node.
   
=== "Cloud Cluster (e.g. Hetzner)"

    ```bash
    $ local-up node join --name eagle --config ~/hetzner-config.yaml 
    Here is some information you need to join a edge node to this cluster.

    VPN:
    If you haven't already linked the node to the netbird.io VPN, you must establish the connection to the VPN beforehand.

    Use can use this setup-key B0BB22F7-4890-41D7-908A-########### to connect to netbird.io VPN. 
    # (1)!

    Kubernetes:
    Use this token 123456.5emtbddc32xjlplg to join the cluster reachable here: 123.123.123.123:443 
    # (2)!
    You have 1 day to join the cluster before this token expires.
    
    If you experience any problems, please consult the documentation at 
    https://edgefarm.github.io/edgefarm/ or file an issue at https://github.com/edgefarm/edgefarm/issues/new?template=question.md
    ```

    1.  Keep the netbird setup-key for connecting the edge node to the VPN.
    2.  Keep this token and cluster ip:port for later to join the node.
   

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

## Join the Edge Node 

=== "Local EdgeFarm cluster"

    Modify the `/etc/hosts` file on your edge node the way the `local-up node join`'s output told you. 

    ```bash
    <IP ADDRESS> edgefarm-control-plane # (1)!
    ```

    1.  In the example above, the IP address is `100.127.213.101`. So your entry must be `100.127.213.101 edgefarm-control-plane`.

    Join the Edge Node to the EdgeFarm cluster by running the `install.sh` script as root.
    Use the bootstrap token from the `local-up node join` output.

    ```console
    ./install.sh --address edgefarm-control-plane:6443 --token <bootstrap-token> --node-ip $(cat /usr/local/etc/wt0.ip) --join --convert --node-type kubeadm
    ```

=== "Cloud Cluster (e.g. Hetzner)"

    Join the Edge Node to the EdgeFarm cluster by running the `install.sh` script as root.
    Use the token and ip:port from the `local-up node join` output.

    ```console
    ./install.sh --address <ip:port> --token <bootstrap-token> --node-ip $(cat /usr/local/etc/wt0.ip) --join --convert --node-type kubeadm 
    ```

## Verify the Edge Node

Verify that the Edge Node is ready to use. Let's say the hostname of the node is `eagle`.

Ensure that the envirnment variable `KUBECONFIG` is set to the kubeconfig file of the EdgeFarm cluster.

```bash
$ KUBECONFIG=~/.edgefarm-local-up/kubeconfig kubectl get nodes | grep eagle
NAME     STATUS   ROLES    AGE   VERSION
eagle    Ready    <none>   1M    v1.22.17
```
