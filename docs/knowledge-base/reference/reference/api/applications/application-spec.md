---
hide:
- toc
---
# API Reference

## Application
<sup><sup>[↩ Parent](#application)</sup></sup>

An `Application` is a definition of a set of containers that can run on a Kubernetes Node. This resource is created by clients and scheduled onto hosts.

| Field                        | Type                 | Description                                              | Required |
| ---------------------------- | -------------------- | -------------------------------------------------------- | -------- |
| **apiVersion**               | core.oam.dev/v1beta1 | Version of the API                                       | true     |
| **kind**                     | Application          | Type of the resource                                     | true     |
| [**metadata**](#objectmeta)  | object               | Standard object's metadata                               | true     |
| [**spec**](#applicationspec) | object               | Specification of the desired behavior of the application | true     |

## Application.spec
<sup><sup>[↩ Parent](#application)</sup></sup>

| Field                                             | Type     | Description                                                                                                                                                                                                                                 | Required |
| ------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| [**components**](#applicationspeccomponentsindex) | []object | List of components belonging to the application. Having multiple components in one application means that there are multiple containers managed by the same application. All components are deployed together and share the same lifecycle. | true     |

## Application.spec.components[index]
<sup><sup>[↩ Parent](#applicationspec)</sup></sup>


| Field                                                       | Type   | Description                                                                                                                                                                                                | Required |
| ----------------------------------------------------------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **name**                                                    | string | The name of the component. This name must be unique between all components in an application.                                                                                                              | true     |
| **type**                                                    | string | The type of the component. This is used to allows the the templating engine behind `EdgeFarm.applications` to generate the correct manifest. Currently the supported component is: *edgefarm-applications* | true     |
| [**properties**](#applicationspeccomponentsindexproperties) | object | Properties of the component. This is used to configure the component. The properties are specific to the component type.                                                                                   | true     |

## Application.spec.components[index].properties
<sup><sup>[↩ Parent](#applicationspeccomponentsindex)</sup></sup>


| Field                                                                             | Type     | Description                                                                                                                                                                                                                                                                                                                                                                                                                                    | Required |
| --------------------------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **name**                                                                          | string   | The name of the container.                                                                                                                                                                                                                                                                                                                                                                                                                     | true     |
| [**nodePoolSelector**](#applicationspeccomponentsindexpropertiesnodepoolselector) | object   | Label selector for nodepools. Every Edge Node has a corresponding nodepool. The nodepool is used to select the Edge Nodes that shall run the component. The nodePoolSelector specifies the nodepools that shall run the component. A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects. | true     |  | [tolerations](#applicationspeccomponentsindexpropertiestolerationsindex) | []Toleration | The pod this Toleration is attached to tolerates any taint that matches the triple `<key,value,effect>` using the matching operator. |  |
| **image**                                                                         | string   | OCI container image name.                                                                                                                                                                                                                                                                                                                                                                                                                      | true     |
| **imagePullPolicy**                                                               | string   | Image pull policy. One of `Always`, `Never`, `IfNotPresent`. Defaults to `Always` if :latest tag is specified, or `IfNotPresent` otherwise.                                                                                                                                                                                                                                                                                                    | true     |
| **imagePullSecrets**                                                              | []string | Specify image pull secrets.                                                                                                                                                                                                                                                                                                                                                                                                                    |          |
| **command**                                                                       | []string | Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided.                                                                                                                                                                                                                                                                                                                               | false    |
| **args**                                                                          | []string | Arguments to the entrypoint. The container image's CMD is used if this is not provided.                                                                                                                                                                                                                                                                                                                                                        | false    |
| [**envs**](#applicationspeccomponentsindexpropertiesenvsindex)                    | []object | List of environment variables to set in the container.                                                                                                                                                                                                                                                                                                                                                                                         | false    |
| [**tolerations**](#applicationspeccomponentsindexpropertiestolerationsindex)      | []object | The pod this Toleration is attached to tolerates any taint that matches the triple `<key,value,effect>` using the matching operator.                                                                                                                                                                                                                                                                                                           | false    |
| [**ports**](#applicationspeccomponentsindexpropertiesportsindex)                  | []object | List of ports to expose from the container.                                                                                                                                                                                                                                                                                                                                                                                                    | false    |
| **cpu**                                                                           | string   | Default values for CPU resources for Requests or Limits is unset.                                                                                                                                                                                                                                                                                                                                                                              | false    |
| **memory**                                                                        | string   | Default values for Requests and Limits on Memory resources for a container.                                                                                                                                                                                                                                                                                                                                                                    |
| [**requests**](#applicationspeccomponentsindexpropertiesrequests)                 | object   | Resources requested by the container.                                                                                                                                                                                                                                                                                                                                                                                                          | false    |
| [**limits**](#applicationspeccomponentsindexpropertieslimits)                     | object   | Resources allowed for the container.                                                                                                                                                                                                                                                                                                                                                                                                           | false    |
| [**securityContext**](#applicationspeccomponentsindexpropertiessecuritycontext)   | object   | The security context to apply.                                                                                                                                                                                                                                                                                                                                                                                                                 | false    |
| [**traits**](#applicationspeccomponentsindexpropertiestraitsindex)                | []object | Traits of the component. Traits of the component. This is used to configure the component. The traits are specific to the component type.                                                                                                                                                                                                                                                                                                      |          |

## Application.spec.components[index].properties.securityContext
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field                                                                                    | Type   | Description                                                                                                                                                                                                                                                                                                                                                                                                               | Required |
| ---------------------------------------------------------------------------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **allowPrivilegedEscalation**                                                            | bool   | AllowPrivilegedEscalation determines if a pod can request to allow privilege escalation. If unspecified, defaults to true.                                                                                                                                                                                                                                                                                                | false    |
| [**capabilities**](#applicationspeccomponentsindexpropertiessecuritycontextcapabilities) | object | The capabilities to add/drop when running containers. Defaults to the default set of capabilities granted by the container runtime.                                                                                                                                                                                                                                                                                       | false    |
| **privileged**                                                                           | bool   | Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false.                                                                                                                                                                                                                                                                                   | false    |
| **readOnlyRootFilesystem**                                                               | bool   | Whether this container has a read-only root filesystem. Default is false.                                                                                                                                                                                                                                                                                                                                                 | false    |
| **runAsGroup**                                                                           | int64  | The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.                                                                                                                                                             | false    |
| **runAsNonRoot**                                                                         | bool   | Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. | false    |
| **runAsUser**                                                                            | int64  | The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.                                                                                                                               | false    |

## Application.spec.components[index].properties.securityContext.capabilities
<sup><sup>[↩ Parent](#applicationspeccomponentsindexpropertiessecuritycontext)</sup></sup>

| Field    | Type     | Description           | Required |
| -------- | -------- | --------------------- | -------- |
| **add**  | []string | Added capabilities.   | false    |
| **drop** | []string | Removed capabilities. | false    |

## Application.spec.components[index].properties.nodepoolSelector
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field                                                                                                   | Type              | Description                                                                                                                                                                                                                                                     | Required |
| ------------------------------------------------------------------------------------------------------- | ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **matchLabels**                                                                                         | map[string]string | matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed. | false    |
| [matchExpressions](#applicationspeccomponentsindexpropertiesnodepoolselectorindexmatchexpressionsindex) | []MatchExpression | matchExpressions is a list of label selector requirements. The requirements are ANDed.                                                                                                                                                                          | false    |

## Application.spec.components[index].properties.nodepoolSelector[index].matchExpressions[index]
<sup><sup>[↩ Parent](#applicationspeccomponentsindexpropertiesnodepoolselector)</sup></sup>

| Field        | Type     | Description                                                                                                                                                                                                                                         | Required |
| ------------ | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **key**      | string   | The label key that the selector applies to.                                                                                                                                                                                                         | true     |
| **operator** | string   | Represents a key's relationship to a set of values. Valid operators are `In`, `NotIn`, `Exists`, `DoesNotExist`. `In` and `NotIn` operators can be used with non-empty values. `Exists` and `DoesNotExist` operators can be used with empty values. | true     |
| **values**   | []string | An array of string values.                                                                                                                                                                                                                          | false    |



## Application.spec.components[index].properties.tolerations[index]
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field                 | Type   | Description                                                                                                                                                                                                                                                                                                                 | Required |
| --------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **key**               | string | The taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.                                                                                                                                      |          |
| **operator**          | string | Operator represents a key's relationship to the value. Valid operators are `Exists` and `Equal`. Defaults to `Equal`. `Exists` is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.                                                                                         |          |
| **value**             | string | Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.                                                                                                                                                                                  |          |
| **effect**            | string | Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are `NoSchedule`, `PreferNoSchedule` and `NoExecute`.                                                                                                                                                       |          |
| **tolerationSeconds** | int64  | TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system. |          |

## Application.spec.components[index].properties.ports[index]
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field             | Type   | Description                                                                                          | Required |
| ----------------- | ------ | ---------------------------------------------------------------------------------------------------- | -------- |
| **name**          | string | The name of the port mapping.                                                                        | false    |
| **containerPort** | int32  | Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.   | true     |
| **hostPort**      | int32  | Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. | false    |
| **protocol**      | string | Protocol for port. Must be UDP or TCP. Defaults to "TCP".                                            | false    |


## Application.spec.components[index].properties.envs[index]
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field                                                                        | Type   | Description                                                                        | Required |
| ---------------------------------------------------------------------------- | ------ | ---------------------------------------------------------------------------------- | -------- |
| **name**                                                                     | string | Name of the environment variable.                                                  | true     |
| **value**                                                                    | string | The value of the environment variable.                                             | false    |
| [**valueFrom**](#applicationspeccomponentsindexpropertiesenvsindexvaluefrom) | object | Source for the environment variable's value. Cannot be used if value is not empty. | false    |

## Application.spec.components[index].properties.envs[index].valueFrom
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field                                                                                             | Type                 | Description                   | Required |
| ------------------------------------------------------------------------------------------------- | -------------------- | ----------------------------- | -------- |
| [**configMapKeyref**](#applicationspeccomponentsindexpropertiesenvsindexvaluefromconfigmapkeyref) | ConfigMapKeySelector | Selects a key of a ConfigMap. | false    |
| [**secretKeyref**](#applicationspeccomponentsindexpropertiesenvsindexvaluefromsecretkeyref)       | SecretKeySelector    | Selects a key of a Secret.    | false    |

## Application.spec.components[index].properties.envs[index].valueFrom.configMapKeyref
<sup><sup>[↩ Parent](#applicationspeccomponentsindexpropertiesenvsindexvaluefrom)</sup></sup>

| Field    | Type   | Description                                                          | Required |
| -------- | ------ | -------------------------------------------------------------------- | -------- |
| **name** | string | The name of the config map in the pod's namespace to select from.    | true     |
| **key**  | string | The key of the config map to select from. Must be a valid secret key | true     |

## Application.spec.components[index].properties.envs[index].valueFrom.secretKeyref
<sup><sup>[↩ Parent](#applicationspeccomponentsindexpropertiesenvsindexvaluefrom)</sup></sup>

| Field    | Type   | Description                                                       | Required |
| -------- | ------ | ----------------------------------------------------------------- | -------- |
| **name** | string | The name of the secret in the pod's namespace to select from.     | true     |
| **key**  | string | The key of the secret to select from. Must be a valid secret key. | true     |

## Application.spec.components[index].properties.requests
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>


| Field      | Type   | Description                                    | Required |
| ---------- | ------ | ---------------------------------------------- | -------- |
| **memory** | string | Memory resource requests. Defaults to "256Mi". | false    |
| **cpu**    | string | CPU resource requests. Defaults to "250m".     | false    |

## Application.spec.components[index].properties.requests.limits
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

| Field      | Type   | Description                                  | Required |
| ---------- | ------ | -------------------------------------------- | -------- |
| **memory** | string | Memory resource limits. Defaults to "256Mi". | false    |
| **cpu**    | string | CPU resource limits. Defaults to "250m".     | false    |

## Application.spec.components[index].properties.traits[index]
<sup><sup>[↩ Parent](#applicationspeccomponentsindexproperties)</sup></sup>

Currently supported traits are

* [edgefarm-network](../network-trait-spec)
* [edgefarm-storage](../storage-trait-spec)

| Field          | Type                   | Description                                                                                      | Required |
| -------------- | ---------------------- | ------------------------------------------------------------------------------------------------ | -------- |
| **name**       | string                 | The name of the trait.                                                                           | true     |
| **properties** | map[string]interface{} | Properties of the trait. Used to configure the trait. Properties are specific to the trait type. | false    |

