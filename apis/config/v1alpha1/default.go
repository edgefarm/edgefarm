/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

// SetDefaultsCluster sets uninitialized fieWlds to their default value.
func SetDefaultsCluster(obj *Cluster) {
	obj.Kind = "Cluster"
	obj.APIVersion = "config.edgefarm.io/v1alpha1"
	obj.Spec.Type = "local"
}

func SetDefaultsGeneral(obj *General) {
	obj.KubeConfigPath = "~/.edgefarm-local-up/kubeconfig"
	obj.StatePath = "~/.edgefarm-local-up/local.json"
}

func SetDefaultsLocal(obj *Local) {
	obj.Name = "edgefarm"
	obj.ApiServerPort = 6443
	obj.NatsPort = 4222
	obj.HttpPort = 80
	obj.HttpsPort = 443
	obj.VirtualEdgeNodes = 2
}

func SetDefaultsHetzner(obj *Hetzner) {
	obj.Name = "edgefarm"
	obj.HCloudToken = "<your hcloud token>"
	obj.KubeConfigPath = "~/.edgefarm-local-up/hetzner"
	obj.ControlPlane = HetznerMachines{
		Count:       3,
		MachineType: "cx21",
	}
	obj.Workers = HetznerMachines{
		Count:       2,
		MachineType: "cx31",
	}
	obj.HetznerCloudRegion = "nbg1"
	obj.HetznerCloudSSHKeyName = "<your ssh key name>"
}

func SetDefaultNetbird(obj *Netbird) {
	obj.SetupKey = "<your netbird.io setup key>"
}
