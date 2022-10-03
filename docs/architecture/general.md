# Design Decissions

## Strict separation into independent components

Strict separation between management of the edge nodes (edgefarm.devices), management of the applications (edgefarm.appplications) and the delivery of the application data (edgefarm.network).
These components can be used independently of each other and serve as a supplement to existing systems.

## Linux as edge node OS

Linux has become the standard for more complex device software in recent years.
software. Furthermore, Linux has the largest embedded hardware support,
driver support and the most complete open source software stack.

## Containers as a core technology

In pure cloud/data center scenarios, container technologies (such as docker) have been established and have displaced running software directly on the OS stack.

The usage of software and services installed directly in the OS is kept to a minimum (drivers, basic services and container runtime) and everything else is run as a container.
