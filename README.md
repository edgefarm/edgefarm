[contributors-shield]: https://img.shields.io/github/contributors/edgefarm/edgefarm.svg?style=for-the-badge
[contributors-url]: https://github.com/edgefarm/edgefarm/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/edgefarm/edgefarm.svg?style=for-the-badge
[forks-url]: https://github.com/edgefarm/edgefarm/network/members
[stars-shield]: https://img.shields.io/github/stars/edgefarm/edgefarm.svg?style=for-the-badge
[stars-url]: https://github.com/edgefarm/edgefarm/stargazers
[issues-shield]: https://img.shields.io/github/issues/edgefarm/edgefarm.svg?style=for-the-badge
[issues-url]: https://github.com/edgefarm/edgefarm/issues
[license-shield]: https://img.shields.io/github/license/edgefarm/edgefarm?logo=mit&style=for-the-badge
[license-url]: https://opensource.org/licenses/AGPL-3.0

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![AGPL 3.0 License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/edgefarm/edgefarm">
    <img src="https://github.com/edgefarm/edgefarm/raw/beta/.images/EdgefarmLogoWithText.png" alt="Logo" height="112">
  </a>

  <h2 align="center">edgefarm</h2>

  <p align="center">
    Seamless edge computing
  </p>
  <hr />
</p>

# About The Project

EdgeFarm is an open source cloud native development platform for edge- and hybrid applications where application assets can be freely moved between edge and cloud. Hybrid applications refer to applications that are partially deployed on edge devices and partially in the cloud.

Edgefarm is heavily based on Kubernetes. EdgeFarm extends Kubernetes with a bunch of great open source projects. EdgeFarm selectively combines and extends these to provide a platform that offers the same comfort of native cloud development.

## Features

- dynamic and secure registration of edge nodes and life cycle management of edge node firmware (edgefarm.core)
- life cycle management of edge- or hybrid applications (edgefarm.applications)
- reliable communication with data retention in the event of network loss and providing secure access for third party systems (edgefarm.network)
- monitoring the whole stack (edgefarm.monitor)

... all done in a cloud native way.

## Why EdgeFarm?

How great would it be if you could write edge software just like cloud software for your Kubernetes based cloud backend? You'd be free to try out a new piece of software nearly effortless, you'd have access to a huge pool of open source software, you could use your existing CD/CD system to roll out your edge software, and so on.

But edge computing differs from cloud computing in one fundamental way. While compute power in the cloud can be scaled automatically at any time, edge devices are tied to specific locations and replacements or upgrades must be done manually on site. This means that network failures or outages cannot simply be bridged by redundancies and taken over by other compute resources.

This results in the requirement that egde devices must be able to run autonomously over a longer period of time and that the acquired data must be buffered until the connection is restored.

All software used on the edge devices must be able to handle unreliable network connections and synchronize with the backend system when the connection is restored.

If this was solved and my Edge device behaved like another Kubernetes node handling unreliable connections, it would make my day-to-day life as a programmer much more pleasant.

And that is why EdgeFarm is being developed.

# Quick Start

What to expect? The `local-up` tool will create a local EdgeFarm cluster running with [kind](https://kind.sigs.k8s.io/). The cluster will be created in a docker container and the EdgeFarm components will be deployed to the cluster. There will also be two edge nodes that are created purely virtual called `edgefarm-worker2` and `edgefarm-worker3` that behave just like real edge nodes.

Before continuing make sure you have the following installed:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [docker](https://docs.docker.com/get-docker/)

Go to [releases](https://github.com/edgefarm/edgefarm/releases) and download the latest release of the EdgeFarm `local-up` binary for your operating system.

To create the cluster run the following commands:

```console
# this step will take a while as a complete kubernetes cluster is created
$ ./local-up-amd64 cluster create

# use kubectl to check if the cluster is ready
$ kubectl get nodes
NAME                     STATUS   ROLES                  AGE   VERSION
edgefarm-control-plane   Ready    control-plane,master   7m    v1.22.7
edgefarm-worker          Ready    <none>                 6m    v1.22.7
edgefarm-worker2         Ready    <none>                 6m    v1.22.7
edgefarm-worker3         Ready    <none>                 7m    v1.22.7

```



## ⚙️ Configuration

## 🎯 Installation

TODO

## 🧪 Testing

TODO

# 💡 Usage

TODO

# 📖 Examples

TODO

# 🐞 Debugging

TODO

# 📜 History

TODO

# 🤝🏽 Contributing

Code contributions are very much **welcome**.

1. Fork the Project
2. Create your Branch (`git checkout -b AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature")
4. Push to the Branch (`git push origin AmazingFeature`)
5. Open a Pull Request targetting the `beta` branch.

# 🫶 Acknowledgements

TODO
