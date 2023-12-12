---
hide:
- toc
---
# API Reference

Trait: `edgefarm-storage`

## Edgefarm-Storage
<sup><sup>[↩ Parent](#edgefarm-storage)</sup></sup>

| Field                                       | Type   | Description                      | Required |
| ------------------------------------------- | ------ | -------------------------------- | -------- |
| [**hostPath**](#edgefarm-storagehostpath)   | object | Declare host path type storage.  | false    |
| [**emptyDir**](#edgefarm-storageemptydir)   | object | Declare empty dir type storage.  | false    |
| [**configMap**](#edgefarm-storageconfigmap) | object | Declare config map type storage. | false    |
| [**secret**](#edgefarm-storagesecret)       | object | Declare secret type storage.     | false    |

## Edgefarm-Storage.emptyDir
<sup><sup>[↩ Parent](#edgefarm-storage)</sup></sup>

| Field         | Type   | Description                                                    | Required |
| ------------- | ------ | -------------------------------------------------------------- | -------- |
| **name**      | string | The name of the empty dir.                                     | true     |
| **mountPath** | string | The path where the empty dir will be mounted in the container. | true     |
| **subPath**   | string | The subpath to mount the empty dir.                            | false    |
| **medium**    | string | By default, the medium volumes are stored. Defaults to "".     | false    |

## Edgefarm-Storage.hostPath
<sup><sup>[↩ Parent](#edgefarm-storage)</sup></sup>

| Field         | Type   | Description                                                                                                                                                               | Required |
| ------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| **name**      | string | The name of the host path.                                                                                                                                                | true     |
| **path**      | string | The path on the host.                                                                                                                                                     | true     |
| **mountPath** | string | The path where the host path will be mounted in the container.                                                                                                            | true     |
| **type**      | string | The type of the host path. Valid values are `Directory`, `DirectoryOrCreate`, `File`, `FileOrCreate`, `Socket`, `CharDevice`, and `BlockDevice`. Defaults to `Directory`. | false    |


## Edgefarm-Storage.configMap
<sup><sup>[↩ Parent](#edgefarm-storage)</sup></sup>

| Field                                                      | Type              | Description                                                                                                      | Required |
| ---------------------------------------------------------- | ----------------- | ---------------------------------------------------------------------------------------------------------------- | -------- |
| **name**                                                   | string            | The name of the config map.                                                                                      | true     |
| **mountOnly**                                              | bool              | If set to true, the config map will only be mounted and not exposed as environment variables. Defaults to false. | false    |
| [**mountToEnv**](#edgefarm-storagesecretmounttoenv)        | object            | Mount the config map to an environment variable.                                                                 | false    |
| [**mountToEnvs**](#edgefarm-storagesecretmounttoenvsindex) | []object          | Mount the config map to multiple environment variables.                                                          | false    |
| **mountPath**                                              | string            | The path where the config map will be mounted in the container.                                                  | false    |
| **subPath**                                                | string            | The subpath to mount the config map.                                                                             | false    |
| **defaultMode**                                            | int32             | The default mode to use when mounting the config map. Defaults to 420.                                           | false    |
| **readOnly**                                               | bool              | If set to true, the config map will be mounted as read-only. Defaults to false.                                  | false    |
| **data**                                                   | map[string]string | The data of the config map.                                                                                      | false    |
| [**items**](#edgefarm-storageconfigmapitems)               | []object          | The items of the config map.                                                                                     | false    |

## Edgefarm-Storage.configMap.items
<sup><sup>[↩ Parent](#edgefarm-storageconfigmap)</sup></sup>

| Field    | Type   | Description                            | Required |
| -------- | ------ | -------------------------------------- | -------- |
| **key**  | string | The key of the item.                   | true     |
| **path** | string | The path of the item.                  | true     |
| **mode** | int32  | The mode of the item. Defaults to 511. | false    |

## Edgefarm-Storage.configMap.mountToEnv
<sup><sup>[↩ Parent](#edgefarm-storageconfigmap)</sup></sup>

| Field            | Type   | Description                           | Required |
| ---------------- | ------ | ------------------------------------- | -------- |
| **envName**      | string | The name of the environment variable. | true     |
| **configMapKey** | string | The key of the config map.            | true     |

## Edgefarm-Storage.configMap.mountToEnvs[index]
<sup><sup>[↩ Parent](#edgefarm-storageconfigmap)</sup></sup>

| Field            | Type   | Description                           | Required |
| ---------------- | ------ | ------------------------------------- | -------- |
| **envName**      | string | The name of the environment variable. | true     |
| **configMapKey** | string | The key of the config map.            | true     |


## Edgefarm-Storage.secret
<sup><sup>[↩ Parent](#edgefarm-storage)</sup></sup>

| Field                                                      | Type              | Description                                                                                                  | Required |
| ---------------------------------------------------------- | ----------------- | ------------------------------------------------------------------------------------------------------------ | -------- |
| **name**                                                   | string            | The name of the secret.                                                                                      | true     |
| **mountOnly**                                              | bool              | If set to true, the secret will only be mounted and not exposed as environment variables. Defaults to false. | false    |
| [**mountToEnv**](#edgefarm-storagesecretmounttoenv)        | object            | Mount the secret to an environment variable.                                                                 | false    |
| [**mountToEnvs**](#edgefarm-storagesecretmounttoenvsindex) | []object          | Mount the secret to multiple environment variables.                                                          | false    |
| **mountPath**                                              | string            | The path where the secret will be mounted in the container.                                                  | true     |
| **subPath**                                                | string            | The subpath to mount the secret.                                                                             | false    |
| **defaultMode**                                            | int32             | The default mode to use when mounting the secret. Defaults to 420.                                           | false    |
| **readOnly**                                               | bool              | If set to true, the secret will be mounted as read-only. Defaults to false.                                  | false    |
| **stringData**                                             | map[string]string | The string data of the secret.                                                                               | false    |
| **data**                                                   | map[string][]byte | The data of the secret.                                                                                      | false    |
| [**items**](#edgefarm-storagesecretitems)                  | object            | The items of the secret.                                                                                     | false    |

## Edgefarm-Storage.secret.items
<sup><sup>[↩ Parent](#edgefarm-storagesecret)</sup></sup>

| Field    | Type     | Description          | Required |
| -------- | -------- | -------------------- | -------- |
| **key**  | []string | The key of the item. | true     |
| **path** | string   | The path             |


## Edgefarm-Storage.secret.mountToEnv
<sup><sup>[↩ Parent](#edgefarm-storagesecret)</sup></sup>

| Field         | Type   | Description                           | Required |
| ------------- | ------ | ------------------------------------- | -------- |
| **envName**   | string | The name of the environment variable. | true     |
| **secretKey** | string | The key of the secret.                | true     |

## Edgefarm-Storage.secret.mountToEnvs[index]
<sup><sup>[↩ Parent](#edgefarm-storagesecret)</sup></sup>

| Field       | Type   | Description                           | Required |
| ----------- | ------ | ------------------------------------- | -------- |
| *envName*   | string | The name of the environment variable. | true     |
| *secretKey* | string | The key of the secret.                | true     |
