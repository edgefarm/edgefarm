---
hide:
- toc
---
# API Reference

Trait: `edgefarm-network`

## Edgefarm-Network
<sup><sup>[↩ Parent](#edgefarm-network)</sup></sup>

| Field                                   | Type   | Description                                                                                                 | Required |
| --------------------------------------- | ------ | ----------------------------------------------------------------------------------------------------------- | -------- |
| [**network**](#edgefarm-networknetwork) | object | The name of the network the component shall be connected to.                                                | true     |
| **daprProtocol**                        | string | The protocol used to connect to the network. Supported protocols are `grpc` and `http`. Defaults to `grpc`. | false    |
| **daprGrpcPort**                        | int32  | The port dapr uses for grpc. Defaults to `3500`.                                                            | false    |
| **daprHttpPort**                        | int32  | The port dapr uses for http. Defaults to `3501`.                                                            | false    |
| **daprAppPort**                         | int32  | The port the application uses to communicate with dapr. Defaults to `50001`.                                | false    |

## Edgefarm-Network.network
<sup><sup>[↩ Parent](#edgefarm-network)</sup></sup>

| Field          | Type   | Description                                                  | Required |
| -------------- | ------ | ------------------------------------------------------------ | -------- |
| **name**       | string | The name of the network the component shall be connected to. | true     |
| **username**   | string | The username used to authenticate to the network.            | true     |
| **subnetwork** | string | The subnetwork used to connect to the network.               | true     |
