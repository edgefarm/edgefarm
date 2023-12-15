# Glossary

This glossary is intended to be a comprehensive, standardized list of EdgeFarm. It includes technical terms that are specific to EdgeFarm and Kubernetes, as well as more general terms that provide useful context.
Please take a look at The [Kubernetes Glossary](https://kubernetes.io/docs/reference/glossary/?fundamental=true) as EdgeFarm is based on Kubernetes.

**Container** - a container is a standard unit of software that packages up code and all its dependencies so the application runs quickly and reliably from one computing environment to another. (1)
{ .annotate }

1.  A container represents a [container](https://kubernetes.io/docs/reference/glossary/?all=true#term-container) in general. A container is a standard unit of software that packages up code and all its dependencies so the application runs quickly and reliably from one computing environment to another. A container is a standard unit of software that packages up code and all its dependencies so the application runs quickly and reliably from one computing environment to another. A Docker container image is a lightweight, standalone, executable package of software that includes everything needed to run an application: code, runtime, system tools, system libraries and settings. Container images become containers at runtime and in the case of Docker containers - images become containers when they run on Docker Engine. Available for both Linux and Windows based apps, containerized software will always run the same, regardless of the infrastructure. Containers isolate software from its environment and ensure that it works uniformly despite differences for instance between development and staging.

**Edge** computing processes data closer to its source, improving real-time performance compared to centralized cloud systems.  (1)
{ .annotate }
   
1.  Cloud and Edge refere to an IT environment that abstracts IT resources in a network, combines them into pools and distributes them. An edge is a computing location at the edge of a network. The associated hardware and software is located at these physical locations. Cloud computing involves the execution of [workloads](https://kubernetes.io/docs/reference/glossary/?all=true#term-workload) in clouds. Edge computing is the execution of workloads on edge devices.

**Edge Node** - a node with special edge case requirements that is located at the edge. (1)
{ .annotate }

1.  A node represents a [Kubernetes node](https://kubernetes.io/docs/reference/glossary/?all=true#term-node) in general. An Edge Node is a remote device located somewhere completely different e.g. Raspberry Pi connected to the [Kubernetes cluster](https://kubernetes.io/docs/reference/glossary/?all=true#term-cluster) running your [workload](https://kubernetes.io/docs/reference/glossary/?all=true#term-workload). Edge Nodes are managed by `edgefarm.core`. Edge Nodes have advanced features enabled aht are needed to run your workload on the Edge in a reliable way. Edge nodes can be accessed via SSH. Edge autonomy is enabled.
  
**Edge autonomy** - edge nodes can operate fully autonomously, even when the connection to the cloud is lost

**edgefarm.core** - core component of EdgeFarm (1)
{ .annotate }

1.  `edgefarm.core` is responsible for node related things like enabling Edge Nodes, node registration and node autonomy.

**edgefarm.applications** - workload definition for EdgeFarm (1)
{ .annotate }

1.  `edgefarm.applications` is responsible for workload related things like rolling out your workload to your Edge Nodes. It is an abstraction layer on top of Kubernetes workload definitions. It allows you to define your workload in a very minimalist format that can be extended to your needs. Using `edgefarm.applications` you can roll out your custom OCI images and configure advanced settings like volumes, envs, commands, args. You decide which Edge Nodes shall run your workload by using labels.


**edgefarm.network** - communication between workload running on Edge Nodes and/or in the cloud (1)
{ .annotate }

1.  `edgefarm.network` is responsible for communication between your workloads running on Edge Nodes and/or in the cloud. It allows you to define streams that act as buffer for your data. Each Edge Device that is part of a Network runs such a stream. Streams can also be used in the cloud aggregating the streams of the Edge Nodes. These streams act as buffers even when the device is offline and needs to operate fully autonomously. Create a Network and let your applications communicate no matter if running in Cloud, Edge or even exported to a third party system. `edgefarm.network` uses [NATS](https://nats.io/) as a message broker. NATS is a lightweight, high-performance cloud native messaging system. It is part of the [CNCF](https://www.cncf.io/).

**edgefarm.monitor** - monitoring of EdgeFarm related components like Edge Nodes, workloads, networks (1)
{ .annotate }

1.  `edgefarm.monitor` is responsible for monitoring of EdgeFarm related components like Edge Nodes, workloads, networks. It is based on [Grafana Mimir](https://grafana.com/oss/mimir/) and [Grafana](https://grafana.com/). 

**edgefarm.portal** - web interface for EdgeFarm (1)
{ .annotate }

1.  Based on Spotifys [Backstage](https://backstage.io/) edgefarm.Portal is the web interface for EdgeFarm. It allows you to manage your EdgeFarm installation. It is the central place to manage your Edge Nodes, your workloads and your networks. 