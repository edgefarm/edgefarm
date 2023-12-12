# EdgeFarm.applications Overview

`EdgeFarm.applications` is a simple convenience wrapper around standard Kubernets and [Openyurt](https://openyurt.io) APIs.

It provides a simple way to describe an application and deploy it to the EdgeFarm cluster. The application delivery model behind it is [Open Application Model](https://oam.dev/), or OAM for short and [KubeVela](https://kubevela.io/).

## Application

Writing a manifest for `EdgeFarm.applications` is easy. Applications are described in `components` of a given type. So called `taints` are used to configure/add/remove specific settings to the components. 

### Components

Currently supported component types are:

- [edgefarm-applications](../application-spec) - a component that allows you to run your custom OCI images on Edge Nodes

### Traits

Currently supported traits are:

- [edgefarm-network](../network-trait-spec) - a trait that allows you to connect your application to a network
- [edgefarm-storage](../storage-trait-spec) - a trait that allows you to mount a volume to your application





