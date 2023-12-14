# Basic example application

# The components

This example shows how a `hybrid-application` could be defined and deployed.
The term `hybrid-application` means that parts of an application is deployed to cloud nodes and parts of it is deployed on edge nodes. All parts, wherever they are deployed, form the overall application.

This example consists of two parts: 

* A `producer` that run on one or more edge nodes. It is a very simple application that produces simulated sensor data.
* A `consumer` that run on cloud nodes and consume the sensor data. The consumer is a basic web application that can be viewed in a browser.

Both parts are connected to the `edgefarm.network`. The producer drops the generated data in a stream that is buffered locally on the Edge Node. Note, that the producer can run on many edge nodes, each collecting data individually. As long as the Edge Devices are connected, the consuming part of the network aggregates the data from all Edge Nodes and puts them into another stream running in the cloud. The `consumer` application reads this stream and provides the data via the web browser.

The following picture shows the overall architecture of the example:

![Architecture](https://github.com/edgefarm/edgefarm/blob/main/examples/basic/.images/architecture.png?raw=true)

# The development

You can use whatever programming language you like to develop your application. This example uses golang for the producer and python for the consumer.

The source code of both `consumer` and `producer` are located [here]([examples/basic](https://github.com/edgefarm/edgefarm/tree/main/examples/basic)). Have a look if you are interesested in the details.
Both components are published as OCI images to ghcr.io. 

The `consumer` and `producer` OCI images are published as `ghcr.io/edgefarm/edgefarm/example-basic-consumer:latest` and `ghcr.io/edgefarm/edgefarm/example-basic-producer:latest` as multi-arch images for amd64 and arm64.

# The deployment

Here you'll learn how the manifests of the `producer` application and the `network` are defined. The `consumer` is a standard Kubernetes deployment and service resource.

## Producer explained

The manifest files are located [here](https://github.com/edgefarm/edgefarm/tree/main/examples/basic/producer/deploy).

The producer is deployed to the edge nodes. The following snippet shows the edgefarm.applications manifest:

```yaml
apiVersion: core.oam.dev/v1beta1 #(1)!
kind: Application
metadata:
  name: example-producer #(2)!
spec:
  components: #(3)!
    - name: producer #(4)!
      type: edgefarm-applications #(5)!
      properties: 
        image: ghcr.io/edgefarm/edgefarm/example-basic-producer:latest #(6)!
        nodepoolSelector: #(7)!
          matchLabels:
            example: "producer" #(8)!
        name: producer #(9)!
        cpu: 0.25 #(10)!
        memory: 256Mi #(11)!
      traits: #(12)!
        - type: edgefarm-network  #(13)!
          properties:
            network:
              name: example-network #(14)!
              subnetwork: edge-to-cloud #(15)!
              user: publish #(16)!
```

1. Every application is of `kind: Application` and `apiversion: core.oam.dev/v1beta1`.
2. Give your application a meaningful name. This name is used to identify the application resource in the cluster.
3. The `components` section defines the components that are part of the application. In this case, there is only one component called `producer`.
4. The `name` of the component. This name must be unique between all components of the application.
5. The `type` of the component. Using `edgefarm-applications` means that the component is deployed to edge nodes.
6. The `image` of the component. This is the OCI image that is deployed to the edge nodes.
7. Every Edge Node is located in it's unique nodepool called the same as the node. The `nodepoolSelector` defines which edge nodes shall run the component.
8. This example uses a label called `example=producer` to identify the edge nodes that shall run the component. The label is set on the edge node using `kubectl label nodepools.apps.openyurt.io <your node> example=producer`.
9. The `name` of the container to run.
10. The `cpu` resources the container is allowed to consume.
11. The `memory` the container is allowed to consume.
12. `traits` are additional configuration parameters that can be added to the component. 
13. There is one trait of type `edgefarm-network` added to the component. This trait is used to connect the component to the `edgefarm.network`.
14. The `name` of the network to connect to.
15. The `subnetwork` of the network to connect to.
16. The `user` of the network to connect with.

The application contains one component called `producer` that runs our OCI image mentioned before. The component is deployed to all nodes, that the corresponding nodepool has the label `example=producer`. There are also some limits defined how much CPU resources and memory the container is allowed to consume. If the application shall be enabled to communicate with a network, a trait of type `edgefarm-network` must be added. In this case, the component is connected to the network `example-network` and the user `publish` is allowed to publish data to the network. We can define multiple sub-networks in the network definition. In this case, the component is connected to the sub-network `edge-to-cloud`.
We referenced a `edgefarm.network`. So let's define it. Without the network resource the application would not be able to start. The following snippet shows the network definition:

```yaml
apiVersion: streams.network.edgefarm.io/v1alpha1 #(1)!
kind: Network
metadata:
  name: example-network #(2)!
spec:
  compositeDeletePolicy: Foreground #(3)!
  parameters: #(4)!
    users: #(5)!
      - name: publish #(6)!
        limits: #(7)!
          payload: -1 #(8)!
          data: -1 #(9)!
          subscriptions: -1 #(10)!
        permissions: #(11)!
          pub: #(12)!
            allow: #(13)!
              - "*.sensor" #(14)!
              - "$JS.API.CONSUMER.>"
              - "$JS.ACK.>"
            deny: [] #(15)!
          sub: #(16)!
            allow: #(17)!
              - "*.sensor" 
            deny: [] #(18)!
    subNetworks: #(19)!
      - name: edge-to-cloud #(20)!
        limits: #(21)!
          fileStorage: 1G #(22)!
          inMemoryStorage: 100M #(23)!
        nodepoolSelector: #(24)!
          matchLabels:
            example: "producer" #(25)!

    streams: #(26)!
      - name: sensor-stream #(27)!
        type: Standard #(28)!
        subNetworkRef: edge-to-cloud #(29)!
        config:
          subjects: #(30)!
            - "sensor" #(31)!
          discard: Old #(32)!
          retention: Limits #(33)!
          storage: File #(34)!
          maxConsumers: -1 #(35)!
          maxMsgSize: -1 #(36)!
          maxMsgs: -1 #(37)!
          maxMsgsPerSubject: -1 #(38)!
          maxBytes: 10000000 #(39)!

      - name: aggregate-stream #(40)!
        type: Aggregate #(41)!
        config:
          discard: Old #(32)!
          retention: Limits #(43)!
          storage: File #(44)!
          maxConsumers: -1 
          maxMsgSize: -1
          maxMsgs: -1
          maxMsgsPerSubject: -1
          maxBytes: 500000000 #(45)!
        references: #(46)!
          - sensor-stream #(47)!
```

1. Every network is of `kind: Network` and `apiversion: streams.network.edgefarm.io/v1alpha1`.
2. Give your network a meaningful name. This name is used to identify the network resource in the cluster.
3. The `compositeDeletePolicy` defines how the network is deleted. This is mandatory to set to `Foreground` to prevent the network from being deleted before all components are deleted.
4. The `parameters` section defines the parameters of the network.
5. The `users` section defines the users that are allowed to publish or subscribe to the network.
6. The `name` of the user.
7. Users can be limited in their actions. The `limits` section defines the limits of the user.
8. The `payload` limit defines how much data a user is allowed to publish to the network. `-1` means unlimited.
9. TBD
10. The `subscriptions` limit defines how many subscriptions a user is allowed to create. `-1` means unlimited.
11. The `permissions` section defines the permissions of the user regarding publishing and subscribing.
12. The `pub` section defines the permissions for publishing.
13. `allow` defines which subjects the user is allowed to publish to.
14. `*.sensor` means that the user is allowed to publish to any subject that ends with `.sensor`. See (Subjet-Based messaging)[https://docs.nats.io/nats-concepts/subjects] for more information on how subjects work.
15. `deny` defines which subjects the user is not allowed to publish to.
16. The `sub` section defines the permissions for subscribing.
17. `allow` defines which subjects the user is allowed to subscribe to.
18. `deny` defines which subjects the user is not allowed to subscribe to.
19. The `subNetworks` section defines the sub-networks that are part of the network. Imagine you have two types of Edge Nodes. One is highly potent while the other has a smaller footprint. The potent one can reserve way more file storage than the smaller one. For both types you can define differnt sub-networks with different characteristics.
20. `name` is of the sub-network. This get referenced in the `edgefarm.application` manifest.
21. 

The network section is split up into several sub-sections in the spec. 
There are `users` that are allowed to publish or subscribe to specific subjects.
There are `subNetworks` that specifies which parts of the network shall run on which nodes. In this case, the sub-network `edge-to-cloud` is deployed to all nodes that have the label `example=producer` and have some file-storage and memory limits set. 
There are `streams` that basically act as buckets. Each bucket has a configurable size. It is configured how long data is kept in the bucket and how much data can be stored in the bucket. It can be defined what to do if the bucket is full. It can drop old messages or block incoming messages. The streams can be defined where to run. They can be either be located in the cloud (no subNetworkRef) or on a edge node (subNetworkRef is set). 

In this example there are two stream definitions that basically act like this:
`sensor-stream`: Create a bucket with the given size and characteristics on each edge node that matches the subNetworks selector. Each edge node gets its own, unique instance of the stream located on the device. The bucket is named `sensor-stream`. The bucket is configured to accept messages with the subject `sensor`. The size is 10000000 bytes. If the bucket is full, drop old messages.

`aggregate-stream`: Create a bucket with the given size and characteristics in the cloud. The bucket is named `aggregate-stream`. The size is 500000000 bytes. If the bucket is full, drop old messages. The bucket is referenced to the `sensor-stream` meaning that it aggregates all data from all sensor stream instances.

There is also a user called `publish` that is allowed to publish n specific subjects - also "*.sensor". The `*` acts as a wildcard. This is needed, because the producer prefixes it's messages with its unique name. 

The suNetwork `cloud-to-edge` defines that each matching edge node is equipped with a component that is part of that specific edgefarm.network. In the end, there is a pod running on each edge node that connects to the network. 

## Consumer explained

The manifest files are located in `./consumer/deploy`

The consumer is deployed to the cloud nodes. The following snippet shows a standard Kubernetes deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/instance: example-consumer
  name: example-consumer
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/instance: example-consumer
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/instance: example-consumer
    spec:
      containers:
        - env:
            - name: NATS_SERVER
              value: nats://nats.nats.svc:4222
            - name: NATS_EXPORT_SUBJECT
              value: "*.sensor"
            - name: NATS_STREAM_NAME
              value: aggregate-stream
            - name: NATS_CREDS
              value: /creds/network.creds
          image: ghcr.io/edgefarm/edgefarm/example-basic-consumer:latest
          imagePullPolicy: IfNotPresent
          name: consumer
          resources:
            limits:
              cpu: 500m
              memory: 256Mi
            requests:
              cpu: 250m
              memory: 128Mi
          volumeMounts:
            - mountPath: /creds/network.creds
              name: creds
              readOnly: true
              subPath: creds
      restartPolicy: Always
      volumes:
        - name: creds
          secret:
            defaultMode: 420
            secretName: example-network-publish
```

This defines a deployment manifest that runs our OCI images referenced earlier. It uses the network credentials as volumes to create the connection to the network. The consumer is configured to consume all messages with the subject `*.sensor` from the stream `aggregate-stream`.

To be able to access the `consumer` application we need a standard Kubernetes service resoure. 

```yaml
apiVersion: v1
kind: Service
metadata:
  name: example-consumer
spec:
  ports:
    - port: 5006
      targetPort: 5006
  selector:
    app.kubernetes.io/instance: example-consumer
```

This service exposes the port 5006 of the consumer application to the cluster.

# In Action

Deploy both, producer and consumer, to the cluster:

```console
$ kubectl apply -f ./producer/deploy
application.core.oam.dev/example-producer created
network.streams.network.edgefarm.io/example-network created
$ kubectl apply -f ./consumer/deploy
deployment.apps/example-consumer created
service/example-consumer created
```

Decide what edge nodes shall run the producer component:

```console
$ kubectl label nodepools.apps.openyurt.io edgefarm-worker3 example=producer
nodepool.apps.openyurt.io/edgefarm-worker3 labeled
```

Wait until everything is up and running:

```console
$ kubectl get pods     
NAME                                                              READY   STATUS    RESTARTS   AGE   IP            NODE               NOMINATED NODE   READINESS GATES
example-consumer-d69db86c8-n25vb                                  1/1     Running   0          10m   10.244.3.35   edgefarm-worker    <none>           <none>
example-network-default-edge-to-cloud-edgefarm-worker3-5fqfsxl5   1/1     Running   0          12m   10.244.1.5    edgefarm-worker3   <none>           <none>
producer-edgefarm-worker3-s9pbw-5d6f874f65-qfqmf                  2/2     Running   0          16m   10.244.1.6    edgefarm-worker3   <none>           <none>
```

Check the streams that are created and watch the messages flowing in:

```console
$ kubectl get streams.nats.crossplane.io -o wide                       
NAME                          EXTERNAL-NAME      READY   SYNCED   DOMAIN                                                   AGE   ADDRESS                     ACCOUNT PUB KEY                                            MESSAGES   BYTES    CONSUMERS
example-network-25gn7-6bhcs   aggregate-stream   True    True     main                                                     10m   nats://nats.nats.svc:4222   ACDB55OTMWZM6LP4R3I3E5WRLJWWVHCWEBLN5ECYOQCN3BTH5NPDMLD4   321        2.0 MB   1
example-network-25gn7-qxc2v   sensor-stream      True    True     example-network-default-edge-to-cloud-edgefarm-worker3   10m   nats://nats.nats.svc:4222   ACDB55OTMWZM6LP4R3I3E5WRLJWWVHCWEBLN5ECYOQCN3BTH5NPDMLD4   321        1.9 MB   0
```

Forward the service of the `consumer` to your local machine and open a browser at http://localhost:5006/serve.

```console
$ kubectl port-forward service/example-consumer 5006:5006
```
