# Join a physical edge node to the local EdgeFarm cluster

## Prerequisites

A running EdgeFarm cluster. See [Create a local EdgeFarm cluster for testing](create-local-cluster.md) for instructions.

This How To shows how to add a Raspberry Pi 4 with vanilla Ubuntu 22.04 as an Edge Node. 
If you have a different hardware or OS your mileage may vary but the general steps should be the same.

## Edge Node preparations

Download the Ubuntu Server 22.04 ISO and flash it to an SD card. See [Ubuntu 22.04 LTS for Raspberry Pi](https://ubuntu.com/download/raspberry-pi) for instructions.

Make sure that you have SSH access to the Pi and that you can log in as root.

Install required packages
```console
sudo apt update
sudo apt-get install unzip curl -y
```

### Install Docker

Install docker using their convenience script:

https://docs.docker.com/engine/install/ubuntu/#install-using-the-convenience-script

As Ubuntu uses cgroup v2 by default, make sure docker uses cgroup v2. Create a file `/etc/docker/daemon.json` with the following content:

```json
{
  "exec-opts": ["native.cgroupdriver=cgroupfs"],
}
```

Restart docker after the change:

```console
sudo systemctl restart docker
```

### Install yurtadm

Download `yurtadm` from the release page of the openyurt project: https://github.com/openyurtio/openyurt/releases/tag/v1.3.4

```console
# on the Pi
sudo curl -L https://github.com/openyurtio/openyurt/releases/download/v1.3.4/yurtadm-v1.3.4-linux-arm64.zip -o /tmp/yurtadm.zip
sudo unzip /tmp/yurtadm.zip -d /tmp/
mv /tmp/linux-arm64/yurtadm /usr/local/bin/yurtadm
sudo chmod +x /usr/local/bin/yurtadm
rm /tmp/yurtadm.zip
rm -r /tmp/linux-arm64
```

Verify that `yurtadm` is installed correctly:

```console
$ yurtadm --version
yurtadm version: projectinfo.Info{GitVersion:"v1.3.4", GitCommit:"609469f", BuildDate:"2023-07-07T06:46:48Z", GoVersion:"go1.18.10", Compiler:"gc", Platform:"linux/arm64", AllVersions:[]string{"unknown"}}
```

### Raspberry Pi 4 specific preparations

For the Raspberry Pi you need to enable some boot options to make it work with Kubernetes.

```console
sed -i '$s/$/ cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory/' /boot/firmware/cmdline.txt
```

### Register the node using `local-up`

