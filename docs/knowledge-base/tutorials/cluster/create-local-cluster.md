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


Once you've got everything set up, go ahead and run the local-up tool. This could take a while, so grab a coffee while you wait. 

```console
local-up cluster create 
```

!!! note "Using cluster config file"
    You can use a custom cluster configuration file by providing the `--config` flag and change the default values. Create a config file by running 
    `local-up cluster create --generate-config --type local` and adjust the values to your needs.
    You need to pass the config file to the `local-up` command with the `--config` flag.

??? "kubeconfig path notes"

    The default location for the kubeconfig file is `~/.edgefarm-local-up/kubeconfig`. This is by intention not to interfer with any existing clusters you might have. This means, you have to set the `KUBECONFIG` environment variable to use the local cluster with `kubectl`:

    ```console
    export KUBECONFIG=~/.edgefarm-local-up/kubeconfig
    ```

    or use the `--kubeconfig` flag with `kubectl`:

    ```console
    kubectl --kubeconfig ~/.edgefarm-local-up/kubeconfig get nodes
    ```

    or the even better choice: **use a tool to manage multiple kube contexts for you, e.g. [kubecm](https://github.com/sunny0826/kubecm) or  [kubie](https://github.com/sbstp/kubie)**
    
## Deploy EdgeFarm components

Now you can deploy the EdgeFarm components to the local cluster. 

```console
$ local-up deploy
```

This will deploy the EdgeFarm components and its dependencies to the local cluster. This might take a few minutes.

If everything went well, you should see something like this:

```{: .console .no-copy}
The EdgeFarm cluster is ready to use! Have fun exploring EdgeFarm.
To access the cluster use 'kubectl', e.g.
 $ kubectl get nodes
```

!!! note "deploy notes"

    See `deploy` options using the `--help` flag to skip certain components or to deploy a specific component only.

If everything went well, you should see something like this:

```{: .console .no-copy}
The local EdgeFarm cluster is ready to use! Have fun exploring EdgeFarm.
To access the cluster use 'kubectl', e.g.
  $ kubectl get nodes
```

## Enable VPN

If you want to join physical edge nodes to the cluster, you need to setup a free account at [netbird.io](https://netbird.io) and create a personal access token first. Follow the [netbird docs](https://docs.netbird.io/how-to/access-netbird-public-api#creating-an-access-token). 
If you don't want to join physical edge nodes, you can skip this step.

If you have a personal access token, you can enable the VPN with the following command:

```console
local-up vpn enable --netbird-token <your-access-token>
```

## Delete the local cluster

If you want to delete the local cluster, you can run the following command:

```console
local-up cluster delete
```

This will delete the local cluster and all its components.

# Considerations

## Costs

Running a local cluster on your machine is totaly free. However, if you want to join physical edge nodes to the cluster, you need to setup a free account at [netbird.io](https://netbird.io).
