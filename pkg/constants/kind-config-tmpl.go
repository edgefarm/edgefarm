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
	KindConfigTemplate = `apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
name: {{.cluster_name}}
networking:
  disableDefaultCNI: {{.disable_default_cni}}
nodes:
  - role: control-plane
    kubeadmConfigPatches:
    - |
      kind: ClusterConfiguration
      apiServer:
        extraArgs:
          kubelet-preferred-address-types: ExternalIP,Hostname,InternalDNS,ExternalDNS
      controllerManager:
        extraArgs:
          cluster-signing-duration: 87600h0m0s
    image: {{.kind_node_image}}
    extraPortMappings:
    - containerPort: 6443
      hostPort: {{.host_api_server_port}}
    - containerPort: {{.host_vpn_port}}
      hostPort: {{.host_vpn_port}}
      listenAddress: "0.0.0.0"
    kubeadmConfigPatches:
    - |
      kind: ClusterConfiguration
      kubernetesVersion: v1.22.7
    - |
      kind: KubeletConfiguration
      cgroupDriver: systemd`

	KindWorkerRoleTemplate = `  - role: worker
    image: {{.kind_node_image}}
    extraPortMappings:
    - containerPort: 4222
      hostPort: {{.host_nats_port}}
      listenAddress: "0.0.0.0"
    - containerPort: 80
      hostPort: {{.host_http_port}}
      listenAddress: "0.0.0.0"
    - containerPort: 443
      hostPort: {{.host_https_port}}
      listenAddress: "0.0.0.0"
    labels:
      ingress-ready: "true"
    kubeadmConfigPatches:
    - |
      kind: ClusterConfiguration
      kubernetesVersion: v1.22.7
    - |
      kind: KubeletConfiguration
      cgroupDriver: systemd`

	KindEdgeRole = `  - role: worker
    image: {{.kind_node_image}}`
)
