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

There are some predefined handy commands that make the setup easier.

`devspace run init`: Initialization, incl. k3d cluster setup.

`devspace run purge`: Remove all created resources, incl. clusters.

`devspace run activate`: Set the kubernetes context pointing to the cluster.

`devspace run update`: Update all dependencies.

To initializing and setup a fresh environment simply call:

```
devspace run init
devspace deploy
```
