Certainly! Here's the complete conversion of the provided Markdown content into tables:

```markdown
# EdgeFarm.applications specification

## Application

An `Application` is a definition of a set of containers that can run on a Kubernetes Node. This resource is created by clients and scheduled onto hosts.

| Field      | Type                                                                                                         | Description                                                                                                                                                                                                                         |
| ---------- | ------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion | core.oam.dev/v1beta1                                                                                         |                                                                                                                                                                                                                                     |
| kind       | Application                                                                                                  |                                                                                                                                                                                                                                     |
| metadata   | [ObjectMeta](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta) | Standard object's metadata. More info: [https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata](https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata) |
| spec       | [ApplicationSpec](#applicationspec)                                                                          | Specification of the desired behavior of the application.                                                                                                                                                                           |

## ApplicationSpec

ApplicationSpec is the description of an application.

| Field      | Type                              | Description                                                                                                                                                                                                                                 |
| ---------- | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| components | [][ComponentSpec](#componentspec) | List of components belonging to the application. Having multiple components in one application means that there are multiple containers managed by the same application. All components are deployed together and share the same lifecycle. |

## ComponentSpec

ComponentSpec is the description of a component. 

| Field      | Type                                        | Description                                                                                                                                                                                            |
| ---------- | ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| name       | string                                      | The name of the component. This name must be unique between all components in an application.                                                                                                          |
| type       | string                                      | The type of the component. This is used to allow the templating engine behind `EdgeFarm.applications` to generate the correct manifest. Currently, the supported component is: *edgefarm-applications* |
| properties | [ComponentProperties](#componentproperties) | Properties of the component. This is used to configure the component. The properties are specific to the component type.                                                                               |

## ComponentProperties

ComponentProperties is the description of a component's properties.

| Field  | Type                | Description                                                                                                      |
| ------ | ------------------- | ---------------------------------------------------------------------------------------------------------------- |
| name   | string              | The name of the container.                                                                                       |
| traits | [][Traits](#traits) | Traits of the component. This is used to configure the component. The traits are specific to the component type. |

### Image

| Field            | Type     | Description                                                                                                                                                                                                                                                                                           |
| ---------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| image            | string   | OCI container image name. More info: [https://kubernetes.io/docs/concepts/containers/images](https://kubernetes.io/docs/concepts/containers/images).                                                                                                                                                  |
| imagePullPolicy  | string   | Image pull policy. One of `Always`, `Never`, `IfNotPresent`. Defaults to `Always` if :latest tag is specified, or `IfNotPresent` otherwise. More info: [https://kubernetes.io/docs/concepts/containers/images#updating-images](https://kubernetes.io/docs/concepts/containers/images#updating-images) |
| imagePullSecrets | []string | Specify image pull secrets. More info: [https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry)                                                                                   |

### Scheduling

| Field            | Type          | Description                                                                                                                                                                                                                        |
| ---------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| nodePoolSelector | LabelSelector | Label selector for nodepools. Every Edge Node has a corresponding nodepool. The nodepool is used to select the Edge Nodes that shall run the component. The nodePoolSelector specifies the nodepools that shall run the component. |

| Field            | Type              | Description                                                                                                                                                                                                                                                     |
| ---------------- | ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| matchLabels      | map[string]string | matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed. |
| matchExpressions | []MatchExpression | matchExpressions is a list of label selector requirements. The requirements are ANDed.                                                                                                                                                                          |

#### MatchExpression

| Field    | Type     | Description                                                                                                                                                                                                                                         |
| -------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| key      | string   | The label key that the selector applies to.                                                                                                                                                                                                         |
| operator | string   | Represents a key's relationship to a set of values. Valid operators are `In`, `NotIn`, `Exists`, `DoesNotExist`. `In` and `NotIn` operators can be used with non-empty values. `Exists` and `DoesNotExist` operators can be used with empty values. |
| values   | []string | An array of string values.                                                                                                                                                                                                                          |

| Field       | Type         | Description                                                                                                                          |
| ----------- | ------------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| tolerations | []Toleration | The pod this Toleration is attached to tolerates any taint that matches the triple `<key,value,effect>` using the matching operator. |

#### Toleration

| Field             | Type   | Description                                                                                                                                                                                                                                                                                                                 |
| ----------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| key               | string | The taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.                                                                                                                                      |
| operator          | string | Operator represents a key's relationship to the value. Valid operators are `Exists` and `Equal`. Defaults to `Equal`. `Exists` is equivalent to a wildcard for value, so that a pod can tolerate all taints of a particular category.                                                                                       |
| value             | string | Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.                                                                                                                                                                                  |
| effect            | string | Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are `NoSchedule`, `PreferNoSchedule`, and `NoExecute`.                                                                                                                                                      |
| tolerationSeconds | int64  | TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system. |

### Entrypoint

| Field   | Type     | Description                                                                                                                                                                                                                                                                                                         |
| ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| command | []string | Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided. More info: [https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint). |
| args    | []string | Arguments to the entrypoint. The container image's CMD is used if this is not provided. More info: [https://kubernetes.io/docs/reference/kubernetes-api                                                                                                                                                             |

/workload-resources/pod-v1/#cmd](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#cmd). |

### Port

| Field         | Type   | Description                                                        |
| ------------- | ------ | ------------------------------------------------------------------ |
| name          | string | The name of this port. This must match the name of a service port. |
| containerPort | int32  | Number or name of the port to access on the container.             |
| protocol      | string | Protocol for port. Must be UDP, TCP, or SCTP. Defaults to "TCP".   |

...

### Traits

| Field  | Type    | Description                               |
| ------ | ------- | ----------------------------------------- |
| traits | []Trait | List of traits to apply to the component. |

#### edgefarm-network trait

| Field        | Type             | Description                                                                                                 |
| ------------ | ---------------- | ----------------------------------------------------------------------------------------------------------- |
| network      | NetworkTraitSpec | The name of the network the component shall be connected to.                                                |
| daprProtocol | string           | The protocol used to connect to the network. Supported protocols are `grpc` and `http`. Defaults to `grpc`. |
| daprGrpcPort | int32            | The port dapr uses for grpc. Defaults to `3500`.                                                            |
| daprHttpPort | int32            | The port dapr uses for http. Defaults to `3501`.                                                            |
| daprAppPort  | int32            | The port the application uses to communicate with dapr. Defaults to `50001`.                                |

#### edgefarm-storage trait

Using the `edgefarm-storage` trait you can mount volumes to your application. By using `configMap` or `secret` it is also possible to generate the configMap/secret resources by defining the data/items in the trait.

| Field     | Type          | Description                                                                                                                                                               |
| --------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| configMap | ConfigMapSpec | Declare config map type storage.                                                                                                                                          |
| secret    | SecretSpec    | Declare secret type storage.                                                                                                                                              |
| emptyDir  | EmptyDirSpec  | Declare empty dir type storage. More Info: [https://kubernetes.io/docs/concepts/storage/volumes/#emptydir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) |
| hostPath  | HostPathSpec  | Declare host path type storage. More info: [https://kubernetes.io/docs/concepts/storage/volumes/#hostpath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) |

...

```