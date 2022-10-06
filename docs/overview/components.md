# EdgeFarm Components

![](../assets/components-overview.png)


EdgeFarm consists of several independently deployable components, each extending the EdgeFarm system with specific functionalities.

## Kubernetes

Kubernetes is not an edgefarm component, but the basis for all components.

Kubernetes already brings many useful features that edgefarm relies on and that can be used within edgefarm.

Functions such as declarative setup, secure authorization, application livecycle management, storage management, redudancy, etc. are available out of the box and only need to be adapted by EdgeFarm.

## edgefarm.core

edgefarm.core extends the Kubernetes cluster with the ability to add edge nodes. This makes it possible to roll out Kubernetes workload to edge devices, the basis for all other components.

In particular, the following functions are available to the user:

* secure token-bassed node authentication
* node revocation
* scheduling of any workload on edge nodes
* retrieve status information from workload

edgefarm.core is made possible by the great open source projects [KubeEdge](https://kubeedge.io/en/) and [HashiCorp Vault](https://www.hashicorp.com/products/vault). Many thanks for this!

## edgefarm.devices

Kubernetes itself can manage and update workload on corresponding nodes. The node os itself must be updated via other ways.

edgefarm.device extends the kubernetes cluster with the functionality to update the os of the edge nodes, e.g. to update the linux kernel or to install driver modules.

In particular, the following functions are available to the user:

* Image signing
* A/B updates with rollback
* scheduled rollouts
* phased rollouts
* gitops ready

## edgefarm.applications

Kubernetes takes an infrastructure-centric approach to defining workload. Workload is defined by deployments, statefullsets, pods, ingresses, etc., which is very powerful. However, this approach also entails a very steep learning curve.

edgefarm.applications encapsulates this approach with the help of [KubeVela](https://kubevela.io) and follows an application-oriented approach.

edgefarm.applications is complemented by other great open source projects to implement a fully comprehensive application livecycle management.

In particular, the following functions are available to the user:

* Application-centric approach
* Application rollout with rollback
* sheduled rollouts
* phased rollouts
* gitops ready

## edgefarm.network

Dealing with unreliable networks is not easy and is usually solved again and again at the application level.

With edgefarm.network a solution is available that encapsulates and solves the problem away from the application developer. The application developer is only offered an API through which he sends his data. edgefarm.network then takes care that the data is transferred reliably.

edgefarm.network uses the open source project [nats](https://docs.nats.io) under the hood for this, a swiss army knife for messaging.

In particular, the following functions are available to the user:

* Issolated messaging networks
* User defined buffer sizes and data redention
* Dapr-based convenience layer for easy access
* Secure access from thirt party systems

## edgefarm.monitor

At the latest when I have installed the devices and distributed the applications, I want to know how my system is doing.

edgefarm.monitor provides an out-of-the-box solution that can monitor all assets of my system, including my applications and inform me about error conditions.

edgefarm.monitor is based on the open source [Grafana](https://grafana.com) stack. Thanks for the great work!

In particular, the following functions are available to the user:

* Monitoring of node health
* Monitoring of application health
* Monitorinf of network capacity
* Preconfigured Grafana Dashboard
* Get Alerts on errors
