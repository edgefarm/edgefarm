# edgefarm.applications Examples

This section contains examples of how to use the edgefarm.applications API.

## Application properties

All options are member of [Application.spec.components[index].properties](../application-spec/#applicationspeccomponentsindexproperties).

### Scheduling

#### nodepoolSelector

The nodePoolSelector is a standard [Kubernetes LabelSelector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#resources-that-support-set-based-requirements).

```yaml
nodePoolSelector:
  matchLabels:
    app: "myapp"
  matchExpressions:
    - key: foo
      operator: In
      values:
        - bar
```

#### tolerations

Tolerations are default Kubernetes tolerations. Follow the [Kubernetes documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) to learn more about tolerations.

```yaml
tolerations:
  - key: "key"
    operator: "Equal"
    value: "value"
    effect: "NoSchedule"
  - operator: "Exists"
    effect: "NoExecute"
```

### Image

#### image

Follow the [Kubernetes documentation](https://kubernetes.io/docs/concepts/containers/images/#image-names) to learn more about image names.

```yaml
image: nginx:latest
```

#### imagePullPolicy

Follow the [Kubernetes documentation](https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy) to learn more about imagePullPolicy.

```yaml
imagePullPolicy: Always
```

#### imagePullSecrets

Follow the [Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) on how to use imagePullSecrets.

```yaml
imagePullSecrets: 
  - name: mysecret
```

### Entrypoint

#### command

Overrides the container image's ENTRYPOINT. Not executed within a shell. 

```yaml
command: 
  - sh
  - "-c"
  - "sleep infinity" 
```

#### args

Arguments to the container image's ENTRYPOINT. The container image's CMD is used if this is not provided. 

```yaml
args:
  - "-c"
  - "sleep infinity"
```

### Environment variables

Define environment variables for the container by value or referencing a ConfigMap or a Secret.

```yaml
envs:
  - name: MY_ENV
    value: "my value"
  - name: MY_ENV_FROM_SECRET
    valueFrom:
      secretKeyRef:
        name: mysecret
        key: mykey
  - name: MY_ENV_FROM_CONFIGMAP
    valueFrom:
      configMapKeyRef:
        name: myconfigmap
        key: mykey
```

### Security context

#### securityContext

The securityContext is a subset of the [Kubernetes SecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/).

```yaml
securityContext:
  allowPrivilegedEscalation: false
  capabilities:
    add:
      - "NET_ADMIN"
    drop:
      - "NET_RAW"
  privileged: true  
  readOnlyRootFilesystem: true
  runAsGroup: 1000
  runAsNonRoot: true
  runAsUser: 1000
```

### Ports

#### ports

```yaml
ports:
  - name: http
    containerPort: 80
    protocol: TCP
    hostPort: 8080
  - containerPort: 9090
    hostPort: 12345
```

### Resources

The notation for cpu and memory resources is the same as in Kubernetes. Follow the [Kubernetes docs](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#resource-units-in-kubernetes) for more information.

#### cpu

Used for requests.cpu and limits.cpu if bot are not specified.

```yaml
cpu: 100m
```

#### memory

Used for requests.memory and limits.memory if bot are not specified.

```yaml
memory: 100Mi
```
 
#### requests

```yaml
requests:
  cpu: 100m
  memory: 100Mi
```

#### limits

```yaml
limits:
  cpu: 100m
  memory: 100Mi
```

### Traits

#### edgefarm-network

By adding the `edgefarm-network` trait, you can connect your application to a `edgefarm.network`.

```yaml
traits:
  - type: edgefarm-network
    properties:
      network:
        name: mynetwork
        subnetwork: big
        user: myuser
```

#### edgefarm-storage

You can define multiple storages for a component. Supported storages are `hostPath`, `emptyDir`, `configMap` and `secret`. Define multiple if you need to mount multiple storages.

```yaml
traits:
  - type: edgefarm-storage
    properties:
      hostPath: 
        - name: myhostpath
          type: DirectoryOrCreate #(1)!
          mountPath: /in-container/mypath
          path: /on-host/mypath
      emptyDir:
        - name: myemptydir
          mountPath: /in-container-/emptydir
      configMap: #(2)!
        - name: myconfigmap
          mountPath: /in-container/myconfigmap
          items:
            - key: foo
              path: foo
      secret: #(3)!
        - name: mysecret
          mountPath: /in-container/mysecret
```

1. If you specify the optional `type` and the path or file to mount does not match the type, the application will fail to start.
2. Let's mount the item `foo` from the ConfigMap `myconfigmap` to `/in-container/myconfigmap/foo`.
3. Let's mount all items from the Secret `mysecret` to `/in-container/mysecret` inside the container.

You can even create ConfigMaps and Secrets with pre-populated data.

```yaml
traits:
  - type: edgefarm-storage
    properties:
      configMap:
        - name: myconfigmap
          mountPath: /in-container/myconfigmap
          mountOnly: false #(1)!
          data: #(2)!
            foo: bar
      secret:
        - name: myconfigmap
          mountPath: /in-container/myconfigmap
          mountOnly: false 
          data: #(3)!
            base64encoded: YmFyCg== #(4)!
          stringData: #(5)!
            foo: bar

```

1. Let's express that the ConfigMap shall be created by setting `mountOnly` to `false`
2. Let's pre-populate the ConfigMap with the key `foo` and the value `bar`.
3. Secrets's `data` field is a map of keys and base64 encoded values.
4. `YmFyCg==` is base64 encoded for `bar`.
5. Secrets's `stringData` field is a map of keys and values as strings.

## Full example appliations
### Basic stress test

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: stress
spec:
  components:
    - name: stress
      type: edgefarm-applications
      properties:
        image: alexeiled/stress-ng
        nodepoolSelector:
          matchLabels:
            app/stress: "" #(1)!
        name: stress
        command: #(2)!
          - "/stress-ng"
          - "--cpu"
          - "4"
          - "--io"
          - "2"
          - "--vm-bytes"
          - "1G"
          - "timeout"
          - "600s"
```

1. This label is used to select the nodepool to deploy the application to. All nodepools that have the label `app/stress=` will be selected. Keep in mind that the values for label selectors can be unset.
2. To run the stress test we'll override the command of the container.

### Application with network and storage

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: myapp
spec:
  components:
    - name: myapp
      type: edgefarm-applications
      properties:
        image: natsio/nats-box:latest
        nodepoolSelector:
          matchLabels:
            mynetwork-big: "" #(1)!
        name: myapp
        command: #(2)!
          - sh
          - "-c"
          - "sleep infinity" 
      traits:
        - type: edgefarm-network #(3)!
          properties:
            network: #(4)!
              name: mynetwork 
              subnetwork: big 
              user: bigonly
        - type: edgefarm-storage #(5)!
          properties:
            configMap: #(6)!
              - name: mycm
                data:
                  foo: bar
                mountPath: /mypath
            emptyDir: #(7)!
              - name: test1 
                mountPath: /test/mount/emptydir
```

1. This application will be deployed to all nodes that have the label `mynetwork-big=`. Keep in mind that the values for label selectors can be unset.
2. Override the command of the container. 
3. A trait `edgefarm-network` is added allowing the application to connect to a network.
4. By configuring the `name`, `subnetwork` and `user` the application will be able to connect to the network. Note, that the network must exist in order to connect to it. In fact, without the network, the application won't be able to start.
5. A trait `edgefarm-storage` is added allowing the application to mount volumes.
6. We'll set-up a ConfigMap called `mycm` and mount it to `/mypath` inside the container. There is also some data pre-populated in the ConfigMap that can be used by the application with the key `foo` and the value `bar`.
7. We'll set-up an emptyDir called `test1` and mount it to `/test/mount/emptydir` inside the container. The emptyDir will be created on the node where the application is deployed. The emptyDir will be deleted when the application is deleted.

