# Create a local EdgeFarm cluster for testing

Before you decide wether EdgeFarm is the right solution for you, you might want to try it out. This guide will help you to set up a EdgeFarm cluster for testing that runs on your local machine.
In the end you will have a EdgeFarm cluster running in a dockerized environment. You can then follow the [Getting Started](../getting-started.md) guide to deploy your first application.

## Prerequisites

=== "Linux"

    [How to install Docker on Linux](https://docs.docker.com/desktop/install/linux-install/)

    Choose the right version of kind and kubectl for your architecture:

    === "amd64"

        [kind-linux-amd64 v0.20.0](https://github.com/kubernetes-sigs/kind/releases/download/v0.20.0/kind-linux-amd64) - click to download

        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/linux/amd64/kubectl) - click to download
    === "arm64" 

        [kind-linux-arm64 v0.20.0](https://github.com/kubernetes-sigs/kind/releases/download/v0.20.0/kind-linux-arm64) - click to download
    
        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/linux/arm64/kubectl) - click to download

=== "MacOS"
    
    [How to install Docker on macOS](https://docs.docker.com/desktop/install/mac-install/)
    
    Choose the right version of kind and kubectl for your architecture:

    === "amd64"

        [kind-darwin-amd64 v0.20.0](https://github.com/kubernetes-sigs/kind/releases/download/v0.20.0/kind-darwin-amd64) - click to download

        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/darwin/amd64/kubectl) - click to download
    === "arm64" 

        [kind-darwin-arm64 v0.20.0](https://github.com/kubernetes-sigs/kind/releases/download/v0.20.0/kind-darwin-arm64) - click to download
    
        [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/darwin/arm64/kubectl) - click to download

=== "Windows running WSL2"

    [How to install Docker in Windows](https://docs.docker.com/desktop/install/windows-install/) and [WSL2 backend](https://docs.docker.com/desktop/wsl/)

    [kind-linux-amd64 v0.20.0](https://github.com/kubernetes-sigs/kind/releases/download/v0.20.0/kind-linux-amd64) - click to download

    [kubectl v1.22.17](https://dl.k8s.io/v1.22.17/bin/linux/amd64/kubectl) - click to download

## Create a local cluster

Download the `local-up` tool from the [EdgeFarm releases page](https://github.com/edgefarm/edgefarm/releases) or run this command:

```console
curl -sfL https://raw.githubusercontent.com/edgefarm/edgefarm/main/install.sh | sh -s -- -b ~/bin
```

Once you've got everything set up, go ahead and run the local-up tool. This could take a while, so grab a coffee while you wait.

```console
local-up cluster create
```

If everything went well, you should see something like this:

```{: .console .no-copy}
The local EdgeFarm cluster is ready to use! Have fun exploring EdgeFarm.
To access the cluster use 'kubectl', e.g.
  % kubectl get nodes
```