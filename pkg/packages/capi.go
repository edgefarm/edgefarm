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

package packages

import (
	"strings"
)

var (
	ClusterAPIOperator = []Packages{
		{
			Manifest: []*Manifest{
				{
					Name: "cluster-api-operator",
					URI:  "https://github.com/kubernetes-sigs/cluster-api-operator/releases/download/v0.8.1/operator-components.yaml",
				},
				{
					Name: "bootstrap-components",
					URI:  "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.6.1/bootstrap-components.yaml",
					PreHook: func(manifest string) (string, error) {
						new := strings.Replace(manifest, "${CAPI_DIAGNOSTICS_ADDRESS:=:8443}", ":8443", -1)
						new = strings.Replace(new, "${CAPI_INSECURE_DIAGNOSTICS:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_MACHINE_POOL:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_KUBEADM_BOOTSTRAP_FORMAT_IGNITION:=false}", "false", -1)
						new = strings.Replace(new, "${KUBEADM_BOOTSTRAP_TOKEN_TTL:=15m}", "15m", -1)
						return new, nil
					},
				},
				{
					Name: "control-plane-components",
					URI:  "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.6.1/control-plane-components.yaml",
					PreHook: func(manifest string) (string, error) {
						new := strings.Replace(manifest, "${CAPI_DIAGNOSTICS_ADDRESS:=:8443}", ":8443", -1)
						new = strings.Replace(new, "${CAPI_INSECURE_DIAGNOSTICS:=false}", "false", -1)
						new = strings.Replace(new, "${CLUSTER_TOPOLOGY:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_KUBEADM_BOOTSTRAP_FORMAT_IGNITION:=false}", "false", -1)
						return new, nil
					},
				},
				{
					Name: "core-components",
					URI:  "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.6.1/core-components.yaml",
					PreHook: func(manifest string) (string, error) {
						new := strings.Replace(manifest, "${CAPI_DIAGNOSTICS_ADDRESS:=:8443}", ":8443", -1)
						new = strings.Replace(new, "${CAPI_INSECURE_DIAGNOSTICS:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_MACHINE_POOL:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_CLUSTER_RESOURCE_SET:=false}", "false", -1)
						new = strings.Replace(new, "${CLUSTER_TOPOLOGY:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_RUNTIME_SDK:=false}", "false", -1)
						new = strings.Replace(new, "${EXP_MACHINE_SET_PREFLIGHT_CHECKS:=false}", "false", -1)
						return new, nil
					},
				},
				{
					Name: "cluster-api-addon-provider-helm",
					URI:  "https://github.com/kubernetes-sigs/cluster-api-addon-provider-helm/releases/download/v0.1.1-alpha.1/addon-components.yaml",
					PreHook: func(manifest string) (string, error) {
						new := strings.Replace(manifest, "${CAAPH_DIAGNOSTICS_ADDRESS:=:8443}", ":8443", -1)
						new = strings.Replace(new, "${CAAPH_INSECURE_DIAGNOSTICS:=false}", "false", -1)
						return new, nil
					},
				},
			},
		},
	}

	ClusterAPIOperatorHetzner = []Packages{
		{
			Manifest: []*Manifest{
				{
					Name: "capi-provider-hetzner",
					URI:  "https://github.com/syself/cluster-api-provider-hetzner/releases/download/v1.0.0-beta.30/infrastructure-components.yaml",
				},
			},
		},
	}
)
