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

package packages

import (
	"context"
	"fmt"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/s0rg/retry"
	"helm.sh/helm/v3/pkg/repo"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"tideland.dev/go/wait"

	"github.com/edgefarm/edgefarm/pkg/args"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/route"
)

type Spec struct {
	Chart           []*helmclient.ChartSpec
	CreateNamespace bool
	ValuesFunc      func() string
}

type Helm struct {
	Repo *repo.Entry
	Spec *Spec
}

type Packages struct {
	Helm []*Helm
}

var (
	ClusterBootstrapFlannel = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "flannel",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "flannel",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/flannel",
								Namespace:   "kube-flannel",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.4.0",
								Timeout:     time.Second * 90,
								ValuesYaml: `flannel:
  command:
  - "bash"
  - "-c"
  - "/opt/bin/flanneld --ip-masq --kube-subnet-mgr --iface=wt0 --iface=eth0 & p=$(ls /sys/class/net); while true; do c=$(ls /sys/class/net); if [ \"$p\" != \"$c\" ]; then echo \"Network changed!\"; sleep 5; pkill -f flanneld; /opt/bin/flanneld --ip-masq --kube-subnet-mgr --iface=wt0 --iface=eth0 & p=$c; fi; sleep 5; done"
  args: []`,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	ClusterBootstrapStage2 = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "headscale",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "headscale",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/headscale",
								Namespace:   "vpn",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.4.0",
								Timeout:     time.Second * 90,
							},
						},
						CreateNamespace: true,
						ValuesFunc: func() string {
							return fmt.Sprintf("config:\n  port: %d\n  hostPort:\n    enabled: true", args.Ports.HostVPNPort)
						},
					},
				},
				{
					Repo: &repo.Entry{
						Name: "tailscale",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "tailscale",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/tailscale",
								Namespace:   "vpn",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.4.0",
								Timeout:     time.Second * 90,
							},
						},
						CreateNamespace: true,
						ValuesFunc: func() string {
							r, err := route.GetRoute(args.Interface)
							if err != nil {
								panic(err)
							}
							return fmt.Sprintf("config:\n  loginServer: http://%s:%d", r.IP, args.Ports.HostVPNPort)
						},
					},
				},
			},
		},
	}

	ClusterDependencies = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "ingress-nginx",
						URL:  "https://kubernetes.github.io/ingress-nginx",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "ingress-nginx",
								ChartName:   "ingress-nginx/ingress-nginx",
								Namespace:   "ingress-nginx",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "4.7.1",
								Timeout:     time.Second * 300,
								ValuesYaml: `controller:
  extraArgs:
    publish-status-address: "localhost"
  publishService:
    enabled: false
  watchIngressWithoutClass: true
  terminationGracePeriodSeconds: 0
  nodeSelector:
    ingress-ready: "true"
  service:
    internal:
      enabled: false
    type: "NodePort"
    nodePorts:
      http: 32080
      https: 32443
  hostPort:
    enabled: true`,
							},
						},
						CreateNamespace: true,
					},
				},
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
								ValuesYaml:  `installCRDs: true`,
							},
						},
						CreateNamespace: true,
					},
				},
				{
					Repo: &repo.Entry{
						Name: "crossplane-stable",
						URL:  "https://charts.crossplane.io/stable",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "crossplane",
								ChartName:   "crossplane-stable/crossplane",
								Namespace:   "crossplane-system",
								UpgradeCRDs: true,
								Version:     "1.12.2",
								Wait:        true,
								Timeout:     time.Second * 300,
								ValuesYaml: `args:
  - --enable-composition-functions
  - --debug
resourcesCrossplane:
  limits:
    cpu: 100m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 256Mi
resourcesRBACManager:
  limits:
    cpu: 100m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 256Mi
xfn:
  enabled: true
  args:
    - --debug
  resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi`,
							},
						},
						CreateNamespace: true,
					},
				},
				{
					Repo: &repo.Entry{
						Name: "vault",
						URL:  "https://kubernetes-charts.banzaicloud.com",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "vault-operator",
								ChartName:   "vault/vault-operator",
								Namespace:   "vault",
								UpgradeCRDs: true,
								Version:     "1.19.0",
								Wait:        true,
								Timeout:     time.Second * 300,
							},
							{
								ReleaseName: "vault-secrets-webhook",
								ChartName:   "vault/vault-secrets-webhook",
								Namespace:   "vault",
								UpgradeCRDs: true,
								Version:     "1.19.0",
							},
						},
						CreateNamespace: true,
					},
				},
				{
					Repo: &repo.Entry{
						Name: "metacontroller",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "metacontroller",
								ChartName:   "oci://ghcr.io/metacontroller/metacontroller-helm",
								Namespace:   "metacontroller",
								UpgradeCRDs: true,
								Version:     "v4.10.4",
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}
	Base = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "vault",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "vault",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/vault",
								Namespace:   "vault",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.4.0",
								Timeout:     time.Second * 300,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	Network = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "edgefarm-network",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "edgefarm-network",
								ChartName:   "oci://ghcr.io/edgefarm/edgefarm.network/edgefarm-network",
								Namespace:   "edgefarm-network",
								UpgradeCRDs: true,
								Version:     "1.0.0-beta.37",
								Wait:        true,
								Timeout:     time.Second * 600,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	Applications = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "edgefarm-applications",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "applications",
								ChartName:   "oci://ghcr.io/edgefarm/edgefarm.applications/edgefarm-applications",
								Namespace:   "edgefarm-applications",
								UpgradeCRDs: true,
								Version:     "1.0.0-beta.27",
								Wait:        true,
								Timeout:     time.Second * 300,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	Monitor = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "edgefarm-monitor",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "edgefarm-monitor",
								ChartName:   "oci://ghcr.io/edgefarm/edgefarm.monitor/edgefarm-monitor",
								Namespace:   "edgefarm-monitor",
								UpgradeCRDs: true,
								Version:     "1.0.0-beta.14 ",
								Wait:        true,
								Timeout:     time.Second * 300,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	YurtManager = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "openyurt",
						URL:  "https://openyurtio.github.io/openyurt-helm",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "yurt-manager",
								ChartName:   "openyurt/yurt-manager",
								Namespace:   "kube-system",
								Version:     "1.3.4",
								UpgradeCRDs: true,
								// 								ValuesYaml: `image:
								//   registry: ghcr.io/openyurtio
								//   repository: openyurt/yurt-manager
								//   tag: v1.3.4`,
							},
						},
					},
				},
			},
		},
	}
)

func Template(h *Helm, index int) ([]byte, error) {
	if h == nil {
		return nil, fmt.Errorf("helm is nil")
	}
	if h.Spec == nil {
		return nil, fmt.Errorf("helm spec is nil")
	}
	if h.Spec.Chart == nil {
		return nil, fmt.Errorf("helm spec chart is nil")
	}
	if len(h.Spec.Chart) <= index {
		return nil, fmt.Errorf("helm spec chart index out of range")
	}

	client, err := helmclient.New(&helmclient.Options{
		Namespace: h.Spec.Chart[index].Namespace,
		Debug:     true,
		Linting:   false,
		DebugLog:  klog.Infof,
	})
	if h.Repo.URL != "" {
		if err := client.AddOrUpdateChartRepo(*h.Repo); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return client.TemplateChart(h.Spec.Chart[index])
}

func (h *Helm) Install() error {
	for _, spec := range h.Spec.Chart {
		client, err := helmclient.New(&helmclient.Options{
			Namespace: spec.Namespace,
			Debug:     true,
			Linting:   false,
			DebugLog:  klog.Infof,
		})
		if h.Repo != nil {
			if h.Repo.URL != "" {
				if err := client.AddOrUpdateChartRepo(*h.Repo); err != nil {
					return err
				}
			}
		}
		if err != nil {
			return err
		}

		if h.Spec.ValuesFunc != nil {
			spec.ValuesYaml = h.Spec.ValuesFunc()
		}

		if h.Spec.CreateNamespace {
			klog.Infof("chart: %s: creating namespace %s", spec.ChartName, spec.Namespace)
			_, err := k8s.CreateNamespaceIfNotExist(spec.Namespace)
			if err != nil {
				return err
			}
		}

		// Install helm chart with repeat mechanism if failed to install
		try := retry.New(
			retry.Count(5),
			retry.Sleep(time.Second*2),
			retry.Verbose(true),
		)
		if err := try.Single(fmt.Sprintf("InstallOrUpgradeChart for chart %s", spec.ChartName),
			func() error {
				release, err := client.InstallOrUpgradeChart(context.Background(), spec)
				if err != nil {
					return err
				}
				if release == nil {
					return fmt.Errorf("release failed")
				}
				klog.Infof("chart: %s: installed release %s", spec.ChartName, release.Name)
				return nil
			}); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packages) Install() error {
	if p.Helm != nil {
		for _, helm := range p.Helm {
			if err := helm.Install(); err != nil {
				return err
			}
		}
	}
	return nil
}

func InstallBase() error {
	for _, pkg := range Base {
		if err := pkg.Install(); err != nil {
			return err
		}
	}
	return nil
}

func InstallDependencies() error {
	for _, pkg := range ClusterDependencies {
		if err := pkg.Install(); err != nil {
			return err
		}
	}
	return nil
}

func Install(packages []Packages) error {
	for _, pkg := range packages {
		if err := pkg.Install(); err != nil {
			return err
		}
	}
	return nil
}

func WaitForBootstrapConditions(stepTimeout time.Duration) error {
	ticker := wait.MakeExpiringIntervalTicker(time.Second, stepTimeout)

	// Checks for flannel pods to be ready on all nodes
	flannelCondition := func() (bool, error) {
		pods, err := k8s.GetPods("kube-system", "app=flannel")
		if err != nil {
			return false, err
		}
		for _, pod := range pods {
			if pod.Status.Phase != v1.PodRunning {
				return false, nil
			}
		}
		return true, nil
	}
	wait.Poll(context.Background(), ticker, flannelCondition)

	// // Checks for core-dns pods to be ready on all nodes
	// corednsCondition := func() (bool, error) {
	// 	pods, err := k8s.GetPods("kube-system", "k8s-app=kube-dns")
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	for _, pod := range pods {
	// 		if pod.Status.Phase != v1.PodRunning {
	// 			return false, nil
	// 		}
	// 	}
	// 	return true, nil
	// }
	// wait.Poll(context.Background(), ticker, corednsCondition)

	// Checks for ready state of all nodes
	nodesCondition := func() (bool, error) {
		nodes, err := k8s.GetNodes(metav1.LabelSelector{})
		if err != nil {
			return false, err
		}
		for _, node := range nodes {
			for _, condition := range node.Status.Conditions {
				if condition.Type == v1.NodeReady && condition.Status != v1.ConditionTrue {
					return false, nil
				}
			}
		}
		return true, nil
	}
	wait.Poll(context.Background(), ticker, nodesCondition)

	return nil
}

func InstallAndWaitBootstrapStage1() error {
	err := k8s.ReplaceCoreDNS()
	if err != nil {
		return err
	}

	for _, pkg := range ClusterBootstrapFlannel {
		if err := pkg.Install(); err != nil {
			return err
		}
	}

	err = WaitForBootstrapConditions(time.Minute * 5)
	if err != nil {
		return err
	}
	return nil
}
