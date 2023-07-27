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

EdgeFarm is an open-source cloud native development platform tailor-made for edge- and hybrid applications. You get the freedom to move those precious application assets back and forth between the edge and the cloud, no strings attached. By hybrid applications, we mean those cool apps that get the best of both worlds, with some parts running on edge devices and others working it in the cloud. It's all about flexibility and making the most of what you've got!

EdgeFarm is all about Kubernetes - it's heavily built on Kubernetes and takes it to the next level by integrating a bunch of awesome open-source projects. With EdgeFarm, we carefully pick and extend these projects to create a platform that gives you that cozy, familiar vibe of native cloud development. It's like getting the best of both worlds, with Kubernetes as the solid foundation and cool extensions to take your edge and hybrid applications to the next level!

## Features

EdgeFarm's got your back with all of that, and it's done the cloud native way:

- First up, we've got **dynamic and secure registration of edge nodes** and smooth life cycle management of edge node firmware `(edgefarm.core)`.
- Then, for your **edge- or hybrid applications**, we've got top-notch life cycle management covered `(edgefarm.applications)`.
- Don't worry about network hiccups—our **reliable communication** ensures data retention during network loss and gives secure access to third-party systems `(edgefarm.network)`.
- And last but not least, we keep an eye on the whole stack with our **monitoring** capabilities `(edgefarm.monitor)`.

It's all about being cloud native and making sure you have everything you need to handle those edge and hybrid challenges like a boss!

## Why EdgeFarm?

It would be absolutely fantastic! Imagine being able to write edge software with the same ease and convenience as cloud software for your Kubernetes-based cloud backend. The possibilities are endless:

1. **Freedom of Placement**: Move application assets between the edge and the cloud with ease, adapting to dynamic requirements.
3. **Optimized Edge Management**: EdgeFarm simplifies dynamic and secure edge node registration and life cycle management, streamlining operations.
4. **Reliable Communication & Data Retention**: Ensure seamless communication and data retention during network loss, enhancing the user experience.
5. **Secure Access for Third-Party Systems**: Guarantee secure access for authorized third-party systems, bolstering your edge infrastructure's security.
6. **Comprehensive Monitoring**: Gain valuable insights by monitoring the entire stack, optimizing performance, and resolving issues.
7. **Unified CI/CD Pipeline**: Effortlessly roll out your edge software alongside cloud deployments using the same CI/CD system, promoting consistency.
2. **Scalability & Flexibility**: Benefit from Kubernetes' scalability and flexibility, allowing your edge software to handle changing demands.

Overall, being able to write edge software just like cloud software for your Kubernetes-based cloud backend empowers you with unmatched freedom, efficiency, and adaptability. It's a game-changer that opens up exciting possibilities for your edge infrastructure!

# Quick Start

Here's the deal: The `local-up` tool will hook you up with an awesome local EdgeFarm cluster running on [kind](https://kind.sigs.k8s.io/). The whole thing will be popping up in a docker environment, and all those EdgeFarm goodies will be deployed right there. You'll even get two virtual edge nodes, going by the names of `edgefarm-worker2` and `edgefarm-worker3`, and they'll act just like the real deal edge nodes. How cool is that?!

Before you move on, make sure you've got the following stuff installed:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [docker](https://docs.docker.com/get-docker/)
- [kind](https://kind.sigs.k8s.io/)

Hop over to [releases](https://github.com/edgefarm/edgefarm/releases) and grab the freshest release of the EdgeFarm `local-up` binary for your operating system.

Time to get the cluster rolling! Execute these commands:

```console
# This step is gonna take a minute, because we're creating a full-blown Kubernetes cluster. So, hang tight and let it do its thing!
$ ./local-up-amd64 cluster create

# Once it's all set, fire up `kubectl` and give the cluster a quick check to see if it's ready to roll.
$ kubectl get nodes
NAME                     STATUS   ROLES                  AGE   VERSION
edgefarm-control-plane   Ready    control-plane,master   7m    v1.22.7
edgefarm-worker          Ready    <none>                 6m    v1.22.7
edgefarm-worker2         Ready    <none>                 6m    v1.22.7
edgefarm-worker3         Ready    <none>                 7m    v1.22.7
```

Let's get that example application deployed!

```console
$ kubectl apply -f examples/basic/producer/deploy
application.core.oam.dev/example-producer created
network.streams.network.edgefarm.io/example-network created
$ kubectl apply -f examples/basic/consumer/deploy
deployment.apps/example-consumer created
service/example-consumer created
```

What next? We're gonna deploy the example producer on an edge node, and the consumer will kick it on a cloud node. It's gonna be a perfect match!
You're the boss now! Label the right node pool to run that edge application. Show 'em who's in charge!

```console
$ kubectl label nodepools.apps.openyurt.io edgefarm-worker3 example=producer
nodepool.apps.openyurt.io/edgefarm-worker3 labeled
```

Hang tight, it's almost showtime! Just give it a moment and wait for those pods to get all set and ready. Almost there!

```console
$ kubectl get pods     
NAME                                                              READY   STATUS    RESTARTS   AGE   IP            NODE               NOMINATED NODE   READINESS GATES
example-consumer-d69db86c8-n25vb                                  1/1     Running   0          10m   10.244.3.35   edgefarm-worker    <none>           <none>
example-network-default-edge-to-cloud-edgefarm-worker3-5fqfsxl5   1/1     Running   0          12m   10.244.1.5    edgefarm-worker3   <none>           <none>
producer-edgefarm-worker3-s9pbw-5d6f874f65-qfqmf                  2/2     Running   0          16m   10.244.1.6    edgefarm-worker3   <none>           <none>
```

Let's take a peek and see what streams resources were created. Time to investigate! It could take a hot minute for those streams to be ready, so don't rush! Just sit back and be patient, it'll be worth the wait!

```console
$ kubectl get streams.nats.crossplane.io -o wide                       
NAME                          EXTERNAL-NAME      READY   SYNCED   DOMAIN                                                   AGE   ADDRESS                     ACCOUNT PUB KEY                                            MESSAGES   BYTES    CONSUMERS
example-network-25gn7-6bhcs   aggregate-stream   True    True     main                                                     10m   nats://nats.nats.svc:4222   ACDB55OTMWZM6LP4R3I3E5WRLJWWVHCWEBLN5ECYOQCN3BTH5NPDMLD4   321        2.0 MB   1
example-network-25gn7-qxc2v   sensor-stream      True    True     example-network-default-edge-to-cloud-edgefarm-worker3   10m   nats://nats.nats.svc:4222   ACDB55OTMWZM6LP4R3I3E5WRLJWWVHCWEBLN5ECYOQCN3BTH5NPDMLD4   321        1.9 MB   0
```

Alright, let's get things rolling! Port forward the `example-consumer` service to your local machine and fire up your browser at http://localhost:5006/serve. Enjoy the view!

```console
$ kubectl port-forward service/example-consumer 5006:5006
```

No sweat at all! Getting the EdgeFarm cluster up and running was a walk in the park. Your local machine is now rocking with a fully functional EdgeFarm cluster, piece of cake!

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

Code contributions are very much **welcome** 🔥🚀

1. Fork the Project
2. Create your Branch (`git checkout -b AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature")
4. Push to the Branch (`git push origin AmazingFeature`)
5. Open a Pull Request targetting the `beta` branch.

# 🫶 Acknowledgements

TODO
