/*
Copyright © 2023 EdgeFarm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package constants

const (
	OpenYurtKindConfig = `apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
name: {{.cluster_name}}
networking:
  disableDefaultCNI: {{.disable_default_cni}}
nodes:
  - role: control-plane
    image: {{.kind_node_image}}
    extraPortMappings:
    - containerPort: 6443
      hostPort: {{.host_api_server_port}}
    - containerPort: 6555
      hostPort: 6555
      listenAddress: "0.0.0.0"`

	KindWorkerRole = `  - role: worker
    image: {{.kind_node_image}}
    extraPortMappings:
    - containerPort: 4222
      hostPort: {{.host_nats_port}}
      listenAddress: "0.0.0.0"
    - containerPort: 7422
      hostPort: {{.host_nats_leafnode_port}}
      listenAddress: "0.0.0.0"
    - containerPort: 80
      hostPort: {{.host_http_port}}
      listenAddress: "0.0.0.0"
    - containerPort: 443
      hostPort: {{.host_https_port}}
      listenAddress: "0.0.0.0"
    labels:
      ingress-ready: "true"`

	KindEdgeRole = `  - role: worker
    image: {{.kind_node_image}}`
)
