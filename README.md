# EdgeFarm

Seamless edge computing.

## Best Practices

### Secret Management

See `examples/secrets/README.md`.

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

`devspace run instanciate-nodes`: Instanciate some virtual kubeedge nodes.

`devspace run purge-nodes`: Destroys the virtual kubeedge nodes.

To initializing and setup a fresh environment simply call:

```sh
devspace run init
devspace deploy
```
