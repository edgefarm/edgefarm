# edgefarm.applications Examples

This section contains examples of how to use the edgefarm.applications API.

## Application parameters

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

