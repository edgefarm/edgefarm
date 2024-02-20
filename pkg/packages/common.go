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
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
)

var (
	CertManager = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "cert-manager",
						URL:  "https://charts.jetstack.io",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "cert-manager",
								ChartName:   "cert-manager/cert-manager",
								Namespace:   "cert-manager",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "v1.12.0",
								Timeout:     time.Second * 300,
								ValuesYaml: `installCRDs: true
image:
  repository: ghcr.io/edgefarm/helm-charts/cert-manager-controller
webhook:
  image:
    repository: ghcr.io/edgefarm/helm-charts/cert-manager-webhook
cainjector:
  image:
    repository: ghcr.io/edgefarm/helm-charts/cert-manager-cainjector`,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}
)
