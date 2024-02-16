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

package v1alpha1

// Cluster contains edgefarm cluster configuration
// +k8s:deepcopy-gen=true
type Cluster struct {
	TypeMeta `yaml:",inline" json:",inline"`
	Hetzner  `yaml:"hetzner,omitempty" json:"hetzner,omitempty"`
}

// +k8s:deepcopy-gen=true
type TypeMeta struct {
	Kind       string `yaml:"kind,omitempty" json:"kind,omitempty"`
	APIVersion string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
}

// +k8s:deepcopy-gen=true
type Hetzner struct {
	//  The HCloudToken is created within a Hetzner Cloud project and needs read/write permissions
	HCloudToken string `yaml:"hcloudToken,omitempty" json:"hcloudToken,omitempty"`
	// The Robot user and password are created here https://robot.hetzner.com/preferences/index -> 'Webservice and app settings'
	HetznerRobotUser     string `yaml:"robotUser,omitempty" json:"robotUser,omitempty"`
	HetznerRobotPassword string `yaml:"robotPassword,omitempty" json:"robotPassword,omitempty"`
}
