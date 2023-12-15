# edgefarm.network Overview

`edgefarm.network` uses [NATS.io](https://nats.io) as the technology behind the scenes. It provides a simple way of describing a communication network scheme and deploying it. We take care of security, scalability and reliability. To make `edgefarm.network` possible, we utilize several Open Source projects:

- [Hashicorp Vault](https://www.vaultproject.io/) with our [custom natssecrets plugin](https://github.com/edgefarm/vault-plugin-secrets-nats) to manage NATS.io credentials
- [NATS.io](https://nats.io) as the messaging system
- [Crossplane](https://www.crossplane.io/) for providing our custom networking resources
- [OpenYurt](https://openyurt.io) for providing the workload resources needed for Edge Nodes
- [Metacontroller](https://github.com/metacontroller/metacontroller) for writing custom Kubernets controllers we need
  
## Network

Writing a manifest for `edgefarm.network` is easy. Networks consist of `users`, `subNetworks`, `streams` and `consumers`. All is defined in one manifest file. Using the network in a Application is done by using the [edgefarm-network](../applications/network-trait-spec.md) trait.

See the [network-spec](../network-spec.md) page for more details.

## Examples

See the [examples](./examples) page on how to learn more about `edgefarm.network` and how to use it.

