# EdgeFarm.applications specification

## Application

An `Application` is a definition of a set of containers that can run on a Kubernetes Node. This resource is created by clients and scheduled onto hosts.

- **apiVersion**: core.oam.dev/v1beta1
- **kind**: Application
- **metadata** ([ObjectMeta](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta)), required <br> Standard object's metadata. More info: [https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata](https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata)
- **spec** ([ApplicationSpec](#applicationspec)), required <br> Specification of the desired behavior of the application.

## ApplicationSpec

ApplicationSpec is the description of an application.

- **components** ([][ComponentSpec](#componentspec)), required <br> List of components belonging to the application. Having multiple components in one applications means that there are multiple containers managed by the same application. All components are deployed together and share the same lifecycle.

## ComponentSpec

ComponentSpec is the description of a component. 

- **name** (string), required <br> The name of the component. This name must be unique between all components in an application.
- **type** (string), required <br> The type of the component. This is used to allows the the templating engine behind `EdgeFarm.applications` to generate the correct manifest. <br> Currently the supported component is: *edgefarm-applications*
- **properties** ([ComponentProperties](#componentproperties)), required <br> Properties of the component. This is used to configure the component. The properties are specific to the component type.

## ComponentProperties

ComponentProperties is the description of a component's properties.

- **name** (string), required <br> The name of the container.
- **traits** ([][Traits](#traits)) <br> Traits of the component. <br> Traits of the component. This is used to configure the component. The traits are specific to the component type.
  
### Image

- **image** (string), required <br> OCI container image name. More info: [https://kubernetes.io/docs/concepts/containers/images](https://kubernetes.io/docs/concepts/containers/images).
- **imagePullPolicy** (string), required <br> Image pull policy. One of `Always`, `Never`, `IfNotPresent`. Defaults to `Always` if :latest tag is specified, or `IfNotPresent` otherwise. More info: [https://kubernetes.io/docs/concepts/containers/images#updating-images](https://kubernetes.io/docs/concepts/containers/images#updating-images)
- **imagePullSecrets** ([]string), optional <br> Specify image pull secrets. More info: [https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry)


### Scheduling

- **nodePoolSelector** (LabelSelector), required <br> Label selector for nodepools. Every Edge Node has a corresponding nodepool. The nodepool is used to select the Edge Nodes that shall run the component. The nodePoolSelector specifies the nodepools that shall run the component.
A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
More info: [https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/#LabelSelector](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/#LabelSelector)

`LabelSelector` represents a selector that matches labels.

  - **matchLabels** (map[string]string), optional <br> matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
  - **matchExpressions** ([]MatchExpression), optional <br> matchExpressions is a list of label selector requirements. The requirements are ANDed.

    `MatchExpression` represents a source for the value of an matching expression that is formed of `<key,operator,values>`. 

      - **key** (string), required <br> The label key that the selector applies to.
      - **operator** (string), required <br> Represents a key's relationship to a set of values. Valid operators are `In`, `NotIn`, `Exists`, `DoesNotExist`. `In` and `NotIn` operators can be used with non-empty values. `Exists` and `DoesNotExist` operators can be used with empty values.
      - **values** ([]string), optional <br> An array of string values. 

- **tolerations** ([]Toleration), optional <br> The pod this Toleration is attached to tolerates any taint that matches the triple `<key,value,effect>` using the matching operator.<br>
  `Toleration` represents a single toleration:

    - **key** (string), optional <br> The taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.
    - **operator** (string), optional <br> Operator represents a key's relationship to the value. Valid operators are `Exists` and `Equal`. Defaults to `Equal`. `Exists` is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.
    - **value** (string), optional <br> Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.
    - **effect** (string), optional <br> Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are `NoSchedule`, `PreferNoSchedule` and `NoExecute`.
    - **tolerationSeconds** (int64), optional <br> TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system.


### Entrypoint

- **command** ([]string), optional <br> Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided. More info: [https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint).

- **args** ([]string), optional <br> Arguments to the entrypoint. The container image's CMD is used if this is not provided. More info: [https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint) and.

### Ports

- **ports** ([][ContainerPort](#ContainerPort)), optional <br> List of ports to expose from the container.<br>  
  `ContainerPort` represents a network port in a single container.

    - **name** (string), optional <br> The name of the port mapping
    - **containerPort** (int32), required <br> Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.
    - **hostPort** (int32), optional <br> Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536.
    - **protocol** (string), optional <br> Protocol for port. Must be UDP or TCP. Defaults to "TCP".

### Environment Variables

- **envs** ([][EnvVar](#envvar)), optional <br> List of environment variables to set in the container.<br>
    `EnvVar` represents an environment variable present in a Container.

    - **name** (string), required <br> Name of the environment variable.
    - **value** (string), optional <br> The value of the environment variable.
    - **valueFrom** (EnvVarSource), optional <br> Source for the environment variable's value. Cannot be used if value is not empty.<br>
      `EnvVarSource` represents a source for the value of an EnvVar. 
  
        - **configMapKeyref** (ConfigMapKeySelector), optional <br> Selects a key of a ConfigMap.
            - **configMapKeyref.name** (string), required <br> The name of the config map in the pod's namespace to select from
            - **configMapKeyref.key** (string), required <br> The key of the config map to select from. Must be a valid secret key
        - **secretKeyref** (SecretKeySelector), optional <br> Selects a key of a Secret.
            - **secretKeyref.name** (string), required <br> The name of the secret in the pod's namespace to select from
            - **secretKeyref.key** (string), required <br> The key of the secret to select from. Must be a valid secret key


### Resources

- **requests** (Requests), optional <br> Resources that are requested by the container. <br>  `Requests` represents resources that are requested by a container.

    - **requests.memory** (string), optional
      Memory resource limits. Defaults to "256Mi".
    - **requests.cpu** (string), optional
      CPU resource limits. Defaults to "250m".

- **limits** (Limits), optional <br> Resources that are allowed for the container.<br>  `Limits` represents resources that are allowed for a container.

    - **limits.memory** (string), optional
      Memory resource limits. Defaults to "256Mi".
    - **limits.cpu** (string), optional
      CPU resource limits. Defaults to "250m".

- **cpu** (string), optional <br>  Default values for CPU resources for Requests or Limits is unset.
- **memory** (string), optional <br> Default values for Requests and Limits on Memory resources for a container. Optional.

### SecurityContext

- **allowPrivilegedEscalation** (bool), optional <br> AllowPrivilegedEscalation determines if a pod can request to allow privilege escalation. If unspecified, defaults to true.
- **capabilities** (Capabilities), optional <br> The capabilities to add/drop when running containers. Defaults to the default set of capabilities granted by the container runtime.<br>
  `Capabilities` describes the set of capabilities that can be requested to add or drop by a pod or container. 
    - **add** ([]string), optional <br> Added capabilities
    - **drop** ([]string), optional <br> Removed capabilities
  
- **privileged** (bool), optional <br> Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false.
- **readOnlyRootFilesystem** (bool), optional <br> Whether this container has a read-only root filesystem. Default is false.
- **runAsGroup** (int64), optional <br> The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.
- **runAsNonRoot** (bool), optional <br> Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
- **runAsUser** (int64), optional <br> The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.

### Traits

- **traits** ([]Trait), optional <br> List of traits to apply to the component. <br>
  `Trait` represents a trait that may be applied to a component.

    - **name** (string), required <br> The name of the trait.
    - **properties** (map[string]interface{}), optional <br> Properties of the trait. This is used to configure the trait. The properties are specific to the trait type.

    Currently there are two traits supported:

    #### edgefarm-network trait

    - **network** (NetworkTraitSpec), required <br> The name of the network the component shall be connected to.<br>
      `NetworkTraitSpec` represents the specification of a network trait.
        - **name** (string), required <br> The name of the network the component shall be connected to.
        - **username** (string), optional <br> The username used to authenticate to the network.
        - **subnetwork** (string), optional <br> The subnetwork used to connect to the network.
    - **daprProtocol** (string), optional <br> The protocol used to connect to the network. Supported protocols are `grpc` and `http`. Defaults to `grpc`.
    - **daprGrpcPort** (int32), optional <br> The port dapr uses for grpc. Defaults to `3500`.
    - **daprHttpPort** (int32), optional <br> The port dapr uses for http. Defaults to `3501`.
    - **daprAppPort** (int32), optional <br> The port the application uses to communicate with dapr. Defaults to `50001`.

    #### edgefarm-storage trait

    Using the `edgefarm-storage` trait you can mount volumes to your application. By using `configMap` or `secret` it is also possible to generate the configMap/secret resources by defining the data/items in the trait.

    - **configMap** (ConfigMapSpec), optional <br> Declare config map type storage.<br>`ConfigMapSpec` represents the specification of a config map.
        - **name** (string), required <br> The name of the config map.
        - **mountOnly** (bool), optional <br> If set to true, the config map will only be mounted and not exposed as environment variables. Defaults to false.
        - **mountToEnv** (MountToEnv), optional <br> Mount the config map to an environment variable.
            - **envName** (string), required <br> The name of the environment variable.
            - **configMapKey** (string), required <br> The key of the config map to be mounted to the environment variable.
        - **mountToEnvs** ([]MountToEnv), optional <br> Mount the config map to multiple environment variables.
            - **envName** (string), required <br> The name of the environment variable.
            - **configMapKey** (string), required <br> The key of the config map to be mounted to the environment variable.
        - **mountPath** (string), optional <br> The path where the config map will be mounted in the container.
        - **subPath** (string), optional <br> The subpath to mount the config map.
        - **defaultMode** (int32), optional <br> The default mode to use when mounting the config map. Defaults to 420.
        - **readOnly** (bool), optional <br> If set to true, the config map will be mounted as read-only. Defaults to false.
        - **data** (map[string]string), optional <br> The data of the config map.
        - **items** ([]ConfigMapItem), optional <br> The items of the config map. <br>`ConfigMapItem` represents the specification of a config map item.
            - **key** (string), required <br> The key of the item.
            - **path** (string), required <br> The path of the item.
            - **mode** (int32), optional <br> The mode of the item. Defaults to 511.
    
    - **secret** (SecretSpec), optional <br> Declare secret type storage.<br>`SecretSpec` represents the specification of a secret.
        - **name** (string), required <br> The name of the secret.
        - **mountOnly** (bool), optional <br> If set to true, the secret will only be mounted and not exposed as environment variables. Defaults to false.
        - **mountToEnv** (MountToEnv), optional <br> Mount the secret to an environment variable.
            - **envName** (string), required <br> The name of the environment variable.
            - **secretKey** (string), required <br> The key of the secret to be mounted to the environment variable.
        - **mountToEnvs** ([]MountToEnv), optional <br> Mount the secret to multiple environment variables.
            - **envName** (string), required <br> The name of the environment variable.
            - **secretKey** (string), required <br> The key of the secret to be mounted to the environment variable.
        - **mountPath** (string), required <br> The path where the secret will be mounted in the container.
        - **subPath** (string), optional <br> The subpath to mount the secret.
        - **defaultMode** (int32), optional <br> The default mode to use when mounting the secret. Defaults to 420.
        - **readOnly** (bool), optional <br> If set to true, the secret will be mounted as read-only. Defaults to false.
        - **stringData** (map[string]string), optional <br> The string data of the secret.
        - **data** (map[string][]byte), optional <br> The data of the secret.
        - **items** (SecretItem), optional <br> The items of the secret.<br> `SecretItem` represents the specification of a secret item.
            - **key** ([]string), required <br> The key of the item.
            - **path** (string), required <br> The path of the item.
            - **mode** (int32), optional <br> The mode of the item. Defaults to 511.
     
    - **emptyDir** (EmptyDirSpec), optional <br> Declare empty dir type storage.<br> `EmptyDirSpec` represents the specification of an empty dir. <br>More Info: [https://kubernetes.io/docs/concepts/storage/volumes/#emptydir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir)
        - **name** (string), required <br> The name of the empty dir.
        - **mountPath** (string), required <br> The path where the empty dir will be mounted in the container.
        - **subPath** (string), optional <br> The subpath to mount the empty dir.
        - **medium** (string), optional <br> By default the medium volumes are store Defaults to "".
    
    - **hostPath** (HostPathSpec), optional <br> Declare host path type storage.<br>`HostPathSpec` represents the specification of a host path. <br>More info: [https://kubernetes.io/docs/concepts/storage/volumes/#hostpath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath)
        - **name** (string), required <br> The name of the host path.
        - **path** (string), required <br> The path on the host.
        - **mountPath** (string), required <br> The path where the host path will be mounted in the container.
        - **type** (string), optional <br> The type of the host path. Valid values are `Directory`, `DirectoryOrCreate`, `File`, `FileOrCreate`, `Socket`, `CharDevice` and `BlockDevice`. Defaults to `Directory`.
