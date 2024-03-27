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
	"fmt"
	"log"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/klog/v2"

	mycontext "github.com/edgefarm/edgefarm/pkg/context"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
)

var (
	Flannel = []Packages{
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
								Version:     "1.21.0",
								Timeout:     time.Second * 90,
								ValuesYaml: `nameOverride: flannel-edge
flannel:
  installCNIPlugin: true
  installCNIConfig: true
  image:
    repository: docker.io/flannel/flannel
    tag: v0.24.2
  image_cni:
    repository: ghcr.io/edgefarm/edgefarm/cni-plugins
    tag: v1
    command:
      - /bin/sh
      - -c
    args:
      - mkdir -p /opt/cni/bin && cp /cni/* /opt/cni/bin && chmod +x /opt/cni/bin/*
  extraVolumes:
    - name: ip
      hostPath:
        path: /usr/local/etc/wt0.ip
        type: File
  extraVolumeMounts:
    - name: ip
      mountPath: /usr/local/etc/wt0.ip
  command: ["bash"]
  args: ["-c" ,"/opt/bin/flanneld --ip-masq --kube-subnet-mgr --iface=edge0 --public-ip=$(cat /usr/local/etc/wt0.ip)"]
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
            - key: node.edgefarm.io/type
              operator: In
              values:
              - "edge"
            - key: node.edgefarm.io/machine
              operator: In
              values:
              - "physical"`,
							},

							// flannel-cloud for kind nodes including virtual edge nodes
							{
								ReleaseName: "flannel-kind",
								ChartName:   "oci://ghcr.io/edgefarm/helm-charts/flannel",
								Namespace:   "kube-flannel",
								UpgradeCRDs: true,
								Wait:        true,
								Version:     "1.21.0",
								Timeout:     time.Second * 90,
								ValuesYaml: `nameOverride: flannel-kind
flannel:
  installCNIPlugin: true
  installCNIConfig: true
  image:
    repository: docker.io/flannel/flannel
    tag: v0.24.2
  image_cni:
    repository: ghcr.io/edgefarm/edgefarm/cni-plugins
    tag: v1
    command:
      - /bin/sh
      - -c
    args:
      - mkdir -p /opt/cni/bin && cp /cni/* /opt/cni/bin && chmod +x /opt/cni/bin/*
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
            - key: node.edgefarm.io/machine
              operator: NotIn
              values:
              - "physical"
          - matchExpressions:
            - key: node.edgefarm.io/type
              operator: In
              values:
                - "cloud"`,
							},
						},
						CreateNamespace: true,
					},
				},
			},
		},
	}

	Kyverno = []Packages{
		{
			Helm: []*Helm{
				{
					Repo: &repo.Entry{
						Name: "kyverno",
						URL:  "https://kyverno.github.io/kyverno",
					},
					Spec: &Spec{
						Chart: []*helmclient.ChartSpec{
							{
								ReleaseName: "kyverno",
								ChartName:   "kyverno/kyverno",
								Namespace:   "kyverno",
								Version:     "v2.5.5",
								UpgradeCRDs: true,
								ValuesYaml: `generatecontrollerExtraResources:
  - nodes
config:
  resourceFilters:
  - '[Event,*,*]'
  - '[*,kube-system,*]'
  - '[*,kube-public,*]'
  - '[*,kube-node-lease,*]'
  - '[APIService,*,*]'
  - '[TokenReview,*,*]'
  - '[SubjectAccessReview,*,*]'
  - '[SelfSubjectAccessReview,*,*]'
  - '[Binding,*,*]'
  - '[ReplicaSet,*,*]'
  - '[ReportChangeRequest,*,*]'
  - '[ClusterReportChangeRequest,*,*]'
  - '[ClusterRole,*,{{ template "kyverno.fullname" . }}:*]'
  - '[ClusterRoleBinding,*,{{ template "kyverno.fullname" . }}:*]'
  - '[ServiceAccount,{{ include "kyverno.namespace" . }},{{ template "kyverno.serviceAccountName" . }}]'
  - '[ConfigMap,{{ include "kyverno.namespace" . }},{{ template "kyverno.configMapName" . }}]'
  - '[ConfigMap,{{ include "kyverno.namespace" . }},{{ template "kyverno.metricsConfigMapName" . }}]'
  - '[Deployment,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}]'
  - '[Job,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}-hook-pre-delete]'
  - '[NetworkPolicy,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}]'
  - '[PodDisruptionBudget,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}]'
  - '[Role,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}:*]'
  - '[RoleBinding,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}:*]'
  - '[Secret,{{ include "kyverno.namespace" . }},{{ template "kyverno.serviceName" . }}.{{ template "kyverno.namespace" . }}.svc.*]'
  - '[Service,{{ include "kyverno.namespace" . }},{{ template "kyverno.serviceName" . }}]'
  - '[Service,{{ include "kyverno.namespace" . }},{{ template "kyverno.serviceName" . }}-metrics]'
  - '[ServiceMonitor,{{ if .Values.serviceMonitor.namespace }}{{ .Values.serviceMonitor.namespace }}{{ else }}{{ template "kyverno.namespace" . }}{{ end }},{{ template "kyverno.serviceName" . }}-service-monitor]'
  - '[Pod,{{ include "kyverno.namespace" . }},{{ template "kyverno.fullname" . }}-test]'`,
							},
						},
						CreateNamespace: true,
					},
				},
			},
			Manifest: []*Manifest{
				{
					Name: "kyverno policy edge-node-annotation",
					Condition: func() bool {
						exists, err := k8s.CrdExists(shared.KubeConfigRestConfig, "clusterpolicies.kyverno.io")
						if err != nil {
							klog.Error(err)
							return false
						}
						if !exists {
							return false
						}
						est, err := k8s.CrdEstablished(shared.KubeConfigRestConfig, "clusterpolicies.kyverno.io")
						if err != nil {
							return false
						}
						return est
					},
					WaitForCondition: true,
					Manifest: `apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: edge-node-annotation
  annotations:
    policies.kyverno.io/title: Add OpenYurt Node binding annotation
spec:
  rules:
  - name: edge-node-annotation
    match:
      any:
      - resources:
          kinds:
          - Node
    preconditions:
      all:
      - key: "{{request.operation || 'BACKGROUND'}}"
        operator: AnyIn
        value:
          - CREATE
          - UPDATE
    mutate:
      targets:
      - apiVersion: v1
        kind: Node
        name: "{{ request.object.metadata.name }}"
      patchStrategicMerge:
        metadata:
          annotations:
            +(apps.openyurt.io/binding): "true"`,
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
								Version:     "1.4.1",
								UpgradeCRDs: true,
								ValuesYaml: `log:
  level: 4
replicaCount: 1
nameOverride: ""
image:
  registry: openyurt
  repository: yurt-manager
  tag: v1.4.1
ports:
  metrics: 10271
  healthProbe: 10272
  webhook: 10273
controllers: "csr-approver-controller,daemon-pod-updater-controller,delegate-lease-controller,node-life-cycle-controller,nodepool-controller,platform-admin-controller,pod-binding-controller,service-topology-endpoints-controller,service-topology-endpointslice-controller,yurt-app-daemon-controller,yurt-app-set-controller,yurt-coordinator-cert-controller,yurt-static-set-controller"
disableIndependentWebhooks: ""
leaderElectResourceName: "cloud-yurt-manager"
resources:
  limits:
    cpu: 2000m
    memory: 1024Mi
  requests:
    cpu: 100m
    memory: 256Mi
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: node-role.kubernetes.io/control-plane
              operator: Exists`,
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
								Version:     "1.23.0",
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
  tag: v8
tolerations:
  - effect: NoSchedule
    key: edgefarm.io
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: node.edgefarm.io/to-be-converted
          operator: In
          values:
          - "true"`
							return fmt.Sprintf(valuesStr, workingMode, nodeServantImage, yurthubImage, enableDummyIf)
						},
					},
				},
			},
		},
	}

	YurtHubLocal = []Packages{
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

	YurtHubCloud = []Packages{
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
    enabled: false
  lookup:
    enabled: true
    secretRef:
      name: hetzner
      namespace: kube-system
      keys:
        host: apiserver-host
        port: apiserver-port
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

	CoreDNS = []Packages{
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

	KubeProxy = []Packages{
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
            - key: node.edgefarm.io/converted
              operator: NotIn
              values:
              - "true"
        - matchExpressions:
            - key: node-role.kubernetes.io/master
              operator: Exists
        - matchExpressions:
            - key: node-role.kubernetes.io/control-plane
              operator: Exists`,
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
            - key: node.edgefarm.io/converted
              operator: In
              values:
                - "true"`,
							},
						},
						CreateNamespace: false,
					},
				},
			},
		},
	}

	VPN = []Packages{
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
							return fmt.Sprintf(`image:
  repository: docker.io/netbirdio/netbird
  pullPolicy: IfNotPresent
  tag: "0.24.4"
config: 
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
          - key: node.edgefarm.io/machine
            operator: NotIn
            values:
            - "physical"
        - matchExpressions:
          - key:   node.edgefarm.io/type
            operator: In
            values:
            - "cloud"`, args.NetbirdSetupKey)
						},
						CreateNamespace: true,
						// Only install/uninstall this helm chart if NetbirdSetupKey is set
						Condition: func() bool {
							return args.NetbirdSetupKey != ""
						},
					},
				},
			},
		},
	}

	Ingress = []Packages{
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
			},
		},
	}
	Crossplane = []Packages{
		{
			Helm: []*Helm{
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
								Version:     "1.12.3",
								Wait:        true,
								Timeout:     time.Second * 300,
								ValuesYaml: `customAnnotations:
  "container.apparmor.security.beta.kubernetes.io/crossplane-xfn": "unconfined"
args:
  - --enable-composition-functions
  - --enable-environment-configs
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
			},
		},
	}
	VaultOperator = []Packages{
		{
			Helm: []*Helm{
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
			},
		},
	}
	Metacontroller = []Packages{
		{
			Helm: []*Helm{
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
	Vault = []Packages{
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
								Version:     "1.24.0",
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
								Version:     "1.0.0-beta.41",
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
								Version:     "1.0.0-beta.29",
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
								Version:     "1.0.0-beta.15",
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
