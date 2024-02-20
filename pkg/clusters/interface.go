/*
Copyright © 2024 EdgeFarm Authors

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

package clusters

type Cluster interface {
	// CreateCluster creates a new cloud cluster
	CreateCluster() error
	// DeleteCluster deletes a cloud cluster
	DeleteCluster() error
	// GetKubeConfig returns the kubeconfig
	GetKubeConfig() (string, error)
	// // GetClusterStatus returns the status of a cloud cluster
	// GetClusterStatus() (string, error)
	// // GetClusterConfig returns the config of a cloud cluster
	// GetClusterConfig() (string, error)
}
