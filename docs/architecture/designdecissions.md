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

## Kubernetes as base system

Kubernetes already brings many useful features that edgefarm relies on and that can be used within edgefarm.

Functions such as declarative setup, secure authorization, application livecycle management, storage management, redudancy, etc. are available out of the box and only need to be adapted by EdgeFarm.

## Strict separation between device firmware and application

Classically, on edge devices are distributed together with the firmware of the device. But this has some disadvantages:

* Updates take a long time
* The amount of data to be transferred is difficult to reduce
* Application development is difficult to separate from firmware development.

EdgeFarm addresses this disadvantages by strictly separating the device firmware from the applications and allowing different groups of users to deploy and operate applications independently from the device software.

The functionality contained in the device firmware is kept to a minimum.

## Compatible with different platforms and architectures

EdgeFarm is developed to support as much edge hardware as possible and to avoid dependencies to the underlying linux os as much as possible.

EdgeFarm is developed from the start to support arm64 and x86 architectures.

## Declarative setup of all components

Declarative setup in the form of manifest files have become established at the latest since the dominance of kubernetes.

Here, the target state is defined in the form of .yaml manifests and made available to the system. The system takes care that the defined target state is established. The system configuration can be administered with it e.g. comfortably over git repositories, versioned, gereviewed and in the case of error also again turned back. This approach is called GitOps.

EdgeFarm adapts this approach for all system configurations brought in by the user, like firmware,system configurations, device configurations, applications, networks etc.

## Automation-first for all components

EdgeFarm provides a basis for edge or hybrid applications.

However, the added value of such an application is usually created outside of EdgeFarm by other application systems that use the derived data to realize e.g. a desktop application for end users.

Accordingly, it is important that EdgeFarm can be easily and fully integrated into external development and delivery processes by making all functions available via APIs and offering them secured to these systems.

## Flexible hosting

The backend is developed without dependencies on cloud providers or other cloud services. The only dependencies are kubernetes and open source software components.

This allows EdgeFarm can be deployed in different scenarios, from public cloud to hybrid cloud to private cloud in your own data center.

## Open Source

Both our own software and all third-party software used in the form of
libraries/frameworks are open source.

All EdgeFarm components are licensed under open source license (either AGPL V3 or MIT).

## Independent of industry/market

EdgeFarm solves a generic problem, which exists independently of the industry/market.
Accordingly, EdgeFarm does not contain any industry specific code parts and can be used for a wide range of use cases.
