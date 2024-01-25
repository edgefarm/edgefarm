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
	NodeManifest = `apiVersion: v1
kind: Node
metadata:
  labels:
    kubernetes.io/os: "{{.name}}"
    kubernetes.io/hostname: "{{.name}}"
    openyurt.io/is-edge-worker: "true"
    apps.openyurt.io/desired-nodepool: "{{.name}}"
    node.edgefarm.io/machine: "physical"
    node.edgefarm.io/type: "edge"
  annotations:
    apps.openyurt.io/binding: "true"
  name: "{{.name}}"
spec:
  taints:
  - effect: NoSchedule
    key: edgefarm.io`

	NodepoolManifest = `apiVersion: apps.openyurt.io/v1beta1
kind: NodePool
metadata:
  name: "{{.name}}"
  labels:
    "{{.name}}": ""
    monitor.edgefarm.io/metrics: default
spec:
  type: Edge
  selector:
    matchLabels:
      apps.openyurt.io/nodepool: "{{.name}}"`
)