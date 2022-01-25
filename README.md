# EdgeFarm

Seamless edge computing.

## Best Practices

### Secret Management

see `examples/secrets/README.md`.

## Local development environment

[devspace](https://devspace.sh/) and [k3d](https://k3d.io/) is used for
local development. In the root folder you can find a
customized `devspace.yaml` which can be used directly.

Dependencies:

- devspace
- k3d
- kubectl
- helm
- openssl

There are some predefined handy commands that make the setup easier.

`devspace run init`: Initialization with k3d cluster setup, kubeedge certs and kubeedge instance.
`devspace run purge`: Remove all created resources, incl. clusters, virtual kubeedge nodes.
`devspace run activate`: Set the kubernetes context pointing to the cluster.
`devspace run update`: Update all dependencies.
`devspace run instantiate-nodes`: Instantiate some virtual kubeedge nodes.
`devspace run purge-nodes`: Destroys the virtual kubeedge nodes.

To initializing and setup a fresh environment call the following commands *(please take a look at the setup notes below)*:

```sh
devspace run init
devspace deploy
```

### Setup notes

**Note 1:** If asked for an IP address, use your [tailscale](https://login.tailscale.com/admin/machines) IP if available.
This is used to ensure that external devices can access the kubeedge instance running locally on the developer machine.
If you don't need access from external devices, you can simply put in the IP addresse that your development machine uses on your local network.
To install tailscale, run the following command:

```sh
devspace run tailscale-install
```

**Note 2:** in order to get edgefarm.network up and running correctly you need to create the edgefarm.network secret first.
Clone [ngs-accounts](https://github.com/edgefarm/ngs-accounts) and run `generate_secrets.sh` to generate the secrets and apply them in to the cluster.

**Note 3:** make sure that you've set your `/etc/hosts` to `user-argocd` for `127.0.0.1` to enanble access to the `argocd` instance running on the kubeedge instance.
To do this, run the following command:

```sh
devspace run etc-hosts
```

**Note 4:** The user for argocd is `admin`, while the password can be read with the following command:

```sh
devspace run argocd-password
```

### Accessing services

Once everything is setup correctly the following services can be reached:

- argocd: [https://user-argocd:8443/](https://user-argocd:8443/)
- keycloack: todo
