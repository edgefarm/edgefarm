# Create a local EdgeFarm cluster for testing

Before you decide wether EdgeFarm is the right solution for you, you might want to try it out. This guide will help you to set up a EdgeFarm cluster for testing that runs on your local machine.
In the end you will have a EdgeFarm cluster running in a dockerized environment. You can then follow the [Getting Started](../getting-started.md) guide to deploy your first application.

!!! warning "Not for production"
    This is not a production-ready setup. It is meant for you trying out EdgeFarm.

## Prerequisites

=== "Linux"

    [How to install Docker on Linux](https://docs.docker.com/desktop/install/linux-install/)

    Choose the right version of kubectl for your architecture:

    === "amd64"

        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/linux/amd64/kubectl) - click to download
    === "arm64" 

        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/linux/arm64/kubectl) - click to download

=== "MacOS"
    
    [How to install Docker on macOS](https://docs.docker.com/desktop/install/mac-install/)
    
    Choose the right version of kubectl for your architecture:

    === "amd64"

        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/darwin/amd64/kubectl) - click to download
    === "arm64" 

        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/darwin/arm64/kubectl) - click to download

=== "Windows running WSL2"

    [How to install Docker in Windows](https://docs.docker.com/desktop/install/windows-install/) and [WSL2 backend](https://docs.docker.com/desktop/wsl/)

    [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/linux/amd64/kubectl) - click to download

## Create a local cluster

Download the `local-up` tool from the [EdgeFarm releases page](https://github.com/edgefarm/edgefarm/releases) or run this command:

```console
curl -sfL https://raw.githubusercontent.com/edgefarm/edgefarm/main/install.sh | sh -s -- -b ~/bin && chmod +x ~/bin/local-up
```

If you want to join physical edge nodes to the cluster, you need to setup a free account at [netbird.io](https://netbird.io) and create a setup-key first. Follow the [netbird docs](https://docs.netbird.io/how-to/register-machines-using-setup-keys). 
If you don't want to join physical edge nodes, you can skip this step and skip the `--netbird-token` flag.
Once you've got everything set up, go ahead and run the local-up tool. This could take a while, so grab a coffee while you wait. 

```console
local-up cluster create --netbird-token <your-netbird-token>
```

If everything went well, you should see something like this:

```{: .console .no-copy}
The local EdgeFarm cluster is ready to use! Have fun exploring EdgeFarm.
To access the cluster use 'kubectl', e.g.
  $ kubectl get nodes
```
