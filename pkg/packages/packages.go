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
	"fmt"
	"log"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/edgefarm/edgefarm/pkg/args"
	mycontext "github.com/edgefarm/edgefarm/pkg/context"
)

var (
	ClusterBootstrapFlannel = []Packages{
		{
			Helm: []*Helm{
				// flannel-edge used for physical edge devices
				{
					Repo: &repo.Entry{
						Name: "flannel",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "flannel-edge",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/flannel",
								Namespace:   "kube-flannel",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.16.0",
								Timeout:     time.Second * 90,
								ValuesYaml: `nameOverride: flannel-edge
flannel:
  installCNIPlugin: false
  image:
    repository: siredmar/flannel
    tag: v0.23.4-siredmar
  extraVolumes:
    - name: ip
      hostPath:
        path: /usr/local/etc/wt0.ip
        type: File
  extraVolumeMounts:
    - name: ip
      mountPath: /usr/local/etc/wt0.ip
  command: ["bash"]
  args: ["-c" ,"/opt/bin/flanneld --ip-masq --kube-subnet-mgr --iface=edge0 --public-ip=$(cat /usr/local/etc/wt0.ip) --persistent-mac"]
  tolerations:
  - key: edgefarm.io
    effect: NoSchedule
  - effect: NoSchedule
    operator: Exists
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
            - key: openyurt.io/is-edge-worker
              operator: In
              values:
              - "true"
            - key: node.edgefarm.io/type
              operator: DoesNotExist`,
							},
							// flannel-cloud for kind nodes including virtual edge nodes
							{
								ReleaseName: "flannel-cloud",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/flannel",
								Namespace:   "kube-flannel",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.16.0",
								Timeout:     time.Second * 90,
								ValuesYaml: `nameOverride: flannel-cloud
flannel:
  installCNIPlugin: false
  installCNIConfig: true
  image:
    repository: siredmar/flannel
    tag: v0.23.4-siredmar
  command:
  - "bash"
  - "-c"
  - "/opt/bin/flanneld --ip-masq --kube-subnet-mgr --iface=wt0 --iface=eth0 & p=$(ls /sys/class/net); while true; do c=$(ls /sys/class/net); if [ \"$p\" != \"$c\" ]; then echo \"Network changed!\"; sleep 5; pkill -f flanneld; /opt/bin/flanneld --ip-masq --kube-subnet-mgr --iface=wt0 --iface=eth0 & p=$c; fi; sleep 5; done"
  args: []
  tolerations:
  - key: edgefarm.io
    effect: NoSchedule
  - effect: NoSchedule
    operator: Exists
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
            - key: openyurt.io/is-edge-worker
              operator: Exists
            - key: node.edgefarm.io/type
              operator: In
              values:
              - "virtual"
          - matchExpressions:
            - key: node-role.kubernetes.io/control-plane
              operator: Exists`,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	ClusterBootstrapKruise = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "openkruise",
						URL:  "https://openkruise.github.io/charts",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "kruise",
								ChartName:   "openkruise/kruise",
								Namespace:   "kruise-system",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.1.0",
								Timeout:     time.Second * 90,
								ValuesYaml: `featureGates: "PodWebhook=false,KruiseDaemon=false,DaemonWatchingPod=false"
installation:
  namespace: kruise-system
  createNamespace: false
manager:
  replicas: 1
  nodeSelector:
    kubernetes.io/hostname: edgefarm-worker`,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	ClusterBootstrapYurtManager = []Packages{
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
							},
						},
					},
				},
			},
		},
	}

	NodeServantApplier = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "node-servant-applier",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "node-servant-applier",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/node-servant-applier",
								Namespace:   "kube-system",
								Version:     "1.16.0",
								UpgradeCRDs: true,
								Wait:        true,
								Timeout:     time.Second * 90,
							},
						},
						ValuesFunc: func() string {
							if !mycontext.Exists("node-servant-applier") {
								log.Fatalf("context for node-servenat-applier does not exist!")
							}
							ctx := mycontext.Context("node-servant-applier")
							workingMode := ""
							nodeServantImage := ""
							yurthubImage := ""
							enableDummyIf := ""
							if val, ok := ctx.Get("working_mode"); ok {
								workingMode = val.(string)
							} else {
								log.Fatalf("workingMode does not exist!")
							}
							if val, ok := ctx.Get("node_servant_image"); ok {
								nodeServantImage = val.(string)
							} else {
								log.Fatalf("")
							}

							if val, ok := ctx.Get("yurthub_image"); ok {
								yurthubImage = val.(string)
							} else {
								log.Fatalf("")
							}
							if val, ok := ctx.Get("enable_dummy_if"); ok {
								enableDummyIf = val.(string)
							} else {
								log.Fatalf("")
							}

							valuesStr := `parameters:
  workingMode: %s
  nodeServantImage: %s
  yurthubImage: %s
  enableDummyIf: %s

image:
  registry: ghcr.io/edgefarm/edgefarm
  repository: node-servant-applier
  tag: v3

tolerations:
  - effect: NoSchedule
    key: edgefarm.io`
							return fmt.Sprintf(valuesStr, workingMode, nodeServantImage, yurthubImage, enableDummyIf)
						},
					},
				},
			},
		},
	}

	ClusterBootstrapYurtHub = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "yurt-hub",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "yurt-hub",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/yurthub",
								Namespace:   "kube-system",
								Version:     "1.14.0",
								UpgradeCRDs: true,
								Wait:        true,
								Timeout:     time.Second * 90,
								ValuesYaml: `kuberneteServerAddr:
  manual:
    enabled: true
    host: edgefarm-control-plane
    port: 6443
  lookup:
    enabled: false
image:
  registry: ghcr.io/openyurtio/openyurt
  repository: yurthub
  tag: v1.4.0`,
							},
						},
					},
				},
			},
		},
	}

	ClusterBootstrapCoreDNS = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "coredns",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "coredns",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/coredns",
								Namespace:   "kube-system",
								Version:     "1.16.0",
								UpgradeCRDs: true,
								Wait:        true,
								Timeout:     time.Second * 90,
								ValuesYaml: `features:
  log:
    enabled: true
dnsPolicy: None
dnsConfig:
  nameservers:
  - 8.8.8.8
  - 8.8.4.4`,
							},
						},
						CreateNamespace: false,
					},
				},
			},
		},
	}

	ClusterBootstrapKubeProxy = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "kube-proxy-default",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "kube-proxy-default",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/kube-proxy",
								Namespace:   "kube-system",
								Version:     "1.16.0",
								UpgradeCRDs: true,
								Wait:        true,
								Timeout:     time.Second * 90,
								ValuesYaml: `nameOverride: kube-proxy-default
kuberneteServerAddr:
  manual:
    enabled: true
    host: edgefarm-control-plane
    port: 6443
  lookup:
    enabled: false
features:
  openyurt:
    enabled: false
tolerations:
- effect: NoSchedule
  operator: Exists
- effect: NoExecute
  operator: Exists
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: openyurt.io/is-edge-worker
              operator: DoesNotExist`,
							},
						},
						CreateNamespace: false,
					},
				},
				{
					Repo: &repo.Entry{
						Name: "kube-proxy-openyurt",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "kube-proxy-openyurt",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/kube-proxy",
								Namespace:   "kube-system",
								Version:     "1.16.0",
								UpgradeCRDs: true,
								Wait:        true,
								Timeout:     time.Second * 90,
								ValuesYaml: `nameOverride: kube-proxy-openyurt
kuberneteServerAddr:
  manual:
    enabled: true
    host: edgefarm-control-plane
    port: 6443
  lookup:
    enabled: false    
features:
  openyurt:
    enabled: true
tolerations:
- effect: NoSchedule
  operator: Exists
- effect: NoExecute
  operator: Exists
- key: edgefarm.io
  effect: NoSchedule
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: openyurt.io/is-edge-worker
              operator: Exists`,
							},
						},
						CreateNamespace: false,
					},
				},
			},
		},
	}

	ClusterBootstrapVPN = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "netbird",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "netbird",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/netbird-client",
								Namespace:   "vpn",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.17.0",
								Timeout:     time.Second * 90,
							},
						},
						ValuesFunc: func() string {
							return fmt.Sprintf(`config:
  managementURL: https://api.wiretrustee.com:443
  auth:
    secret: netbird-auth
    secretKey: NB_SETUP_KEY
    create: true
    value: %s
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
          - key: openyurt.io/is-edge-worker
            operator: Exists
          - key: node.edgefarm.io/type
            operator: In
            values:
              - "virtual"
        - matchExpressions:
          - key: node-role.kubernetes.io/control-plane
            operator: Exists`, args.NetbirdToken)
						},
						CreateNamespace: true,
						// Only install this helm chart if NetbirdToken is set via args
						Condition: func() bool {
							return args.NetbirdToken != ""
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
)
