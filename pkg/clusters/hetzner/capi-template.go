package hetzner

const (
	capiTemplate = `apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
kind: KubeadmConfigTemplate
metadata:
  name: {{.CLUSTER_NAME}}-md-0
spec:
  template:
    spec:
      files:
      - content: |
          net.ipv4.conf.lxc*.rp_filter = 0
        owner: root:root
        path: /etc/sysctl.d/99-cilium.conf
        permissions: "0744"
      - content: |
          overlay
          br_netfilter
        owner: root:root
        path: /etc/modules-load.d/crio.conf
        permissions: "0744"
      - content: |
          net.bridge.bridge-nf-call-iptables  = 1
          net.bridge.bridge-nf-call-ip6tables = 1
          net.ipv4.ip_forward                 = 1
        owner: root:root
        path: /etc/sysctl.d/99-kubernetes-cri.conf
        permissions: "0744"
      - content: |
          vm.overcommit_memory=1
          kernel.panic=10
          kernel.panic_on_oops=1
        owner: root:root
        path: /etc/sysctl.d/99-kubelet.conf
        permissions: "0744"
      - content: |
          nameserver 1.1.1.1
          nameserver 1.0.0.1
          nameserver 2606:4700:4700::1111
        owner: root:root
        path: /etc/kubernetes/resolv.conf
        permissions: "0744"
      joinConfiguration:
        nodeRegistration:
          kubeletExtraArgs:
            anonymous-auth: "false"
            authentication-token-webhook: "true"
            authorization-mode: Webhook
            cloud-provider: external
            event-qps: "5"
            kubeconfig: /etc/kubernetes/kubelet.conf
            max-pods: "220"
            read-only-port: "0"
            resolv-conf: /etc/kubernetes/resolv.conf
            rotate-server-certificates: "true"
            tls-cipher-suites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256
      preKubeadmCommands:
      - set -x
      - grep VERSION= /etc/os-release; uname -a
      - export CONTAINERD=1.7.13
      - export KUBERNETES_VERSION=$(echo {{.KUBERNETES_VERSION}} | sed 's/^v//')
      - export TRIMMED_KUBERNETES_VERSION=$(echo {{.KUBERNETES_VERSION}} | sed 's/^v//'
        | awk -F . '{print $1 "." $2}')
      - ARCH=amd64
      - if [ "$(uname -m)" = "aarch64" ]; then ARCH=arm64; fi
      - localectl set-locale LANG=en_US.UTF-8
      - localectl set-locale LANGUAGE=en_US.UTF-8
      - apt-get update -y
      - apt-get -y install at jq unzip wget socat mtr logrotate apt-transport-https
      - sed -i '/swap/d' /etc/fstab
      - swapoff -a
      - modprobe overlay && modprobe br_netfilter && sysctl --system
      - wget https://github.com/containerd/containerd/releases/download/v$CONTAINERD/cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz
      - wget https://github.com/containerd/containerd/releases/download/v$CONTAINERD/cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz.sha256sum
      - sha256sum --check cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz.sha256sum
      - tar --no-overwrite-dir -C / -xzf cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz
      - rm -f cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz.sha256sum
      - chmod -R 644 /etc/cni && chown -R root:root /etc/cni
      - mkdir -p /etc/containerd
      - containerd config default > /etc/containerd/config.toml
      - sed -i  "s/SystemdCgroup = false/SystemdCgroup = true/" /etc/containerd/config.toml
      - systemctl daemon-reload && systemctl enable containerd && systemctl start
        containerd
      - curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key
        add -
      - echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
      - apt-get update
      - apt-get install -y kubelet=$KUBERNETES_VERSION-00 kubeadm=$KUBERNETES_VERSION-00
        kubectl=$KUBERNETES_VERSION-00  bash-completion && apt-mark hold kubelet kubectl
        kubeadm && systemctl enable kubelet
      - kubeadm config images pull --kubernetes-version $KUBERNETES_VERSION
      - echo 'source <(kubectl completion bash)' >>~/.bashrc
      - echo 'export KUBECONFIG=/etc/kubernetes/admin.conf' >>~/.bashrc
      - apt-get -y autoremove && apt-get -y clean all
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: {{.CLUSTER_NAME}}
  labels:
    cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - 10.244.0.0/16
  controlPlaneRef:
    apiVersion: controlplane.cluster.x-k8s.io/v1beta1
    kind: KubeadmControlPlane
    name: {{.CLUSTER_NAME}}-control-plane
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: HetznerCluster
    name: {{.CLUSTER_NAME}}
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  labels:
    nodepool: {{.CLUSTER_NAME}}-md-0
  name: {{.CLUSTER_NAME}}-md-0
spec:
  clusterName: {{.CLUSTER_NAME}}
  replicas: {{.WORKER_MACHINE_COUNT}}
  selector:
    matchLabels: null
  template:
    metadata:
      labels:
        nodepool: {{.CLUSTER_NAME}}-md-0
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfigTemplate
          name: {{.CLUSTER_NAME}}-md-0
      clusterName: {{.CLUSTER_NAME}}
      failureDomain: {{.HCLOUD_REGION}}
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: HCloudMachineTemplate
        name: {{.CLUSTER_NAME}}-md-0
      version: {{.KUBERNETES_VERSION}}
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineHealthCheck
metadata:
  name: {{.CLUSTER_NAME}}-control-plane-unhealthy-5m
spec:
  clusterName: {{.CLUSTER_NAME}}
  maxUnhealthy: 100%
  nodeStartupTimeout: 15m
  remediationTemplate:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: HCloudRemediationTemplate
    name: control-plane-remediation-request
  selector:
    matchLabels:
      cluster.x-k8s.io/control-plane: ""
  unhealthyConditions:
  - status: Unknown
    timeout: 180s
    type: Ready
  - status: "False"
    timeout: 180s
    type: Ready
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineHealthCheck
metadata:
  name: {{.CLUSTER_NAME}}-md-0-unhealthy-5m
spec:
  clusterName: {{.CLUSTER_NAME}}
  maxUnhealthy: 100%
  nodeStartupTimeout: 10m
  remediationTemplate:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: HCloudRemediationTemplate
    name: worker-remediation-request
  selector:
    matchLabels:
      nodepool: {{.CLUSTER_NAME}}-md-0
  unhealthyConditions:
  - status: Unknown
    timeout: 180s
    type: Ready
  - status: "False"
    timeout: 180s
    type: Ready
---
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: KubeadmControlPlane
metadata:
  name: {{.CLUSTER_NAME}}-control-plane
spec:
  kubeadmConfigSpec:
    clusterConfiguration:
      apiServer:
        extraArgs:
          authorization-mode: Node,RBAC
          client-ca-file: /etc/kubernetes/pki/ca.crt
          cloud-provider: external
          default-not-ready-toleration-seconds: "45"
          default-unreachable-toleration-seconds: "45"
          enable-aggregator-routing: "true"
          enable-bootstrap-token-auth: "true"
          encryption-provider-config: /etc/kubernetes/encryption-provider.yaml
          etcd-cafile: /etc/kubernetes/pki/etcd/ca.crt
          etcd-certfile: /etc/kubernetes/pki/etcd/server.crt
          etcd-keyfile: /etc/kubernetes/pki/etcd/server.key
          kubelet-client-certificate: /etc/kubernetes/pki/apiserver-kubelet-client.crt
          kubelet-client-key: /etc/kubernetes/pki/apiserver-kubelet-client.key
          kubelet-preferred-address-types: ExternalIP,Hostname,InternalDNS,ExternalDNS
          profiling: "false"
          proxy-client-cert-file: /etc/kubernetes/pki/front-proxy-client.crt
          proxy-client-key-file: /etc/kubernetes/pki/front-proxy-client.key
          requestheader-allowed-names: front-proxy-client
          requestheader-client-ca-file: /etc/kubernetes/pki/front-proxy-ca.crt
          requestheader-extra-headers-prefix: X-Remote-Extra-
          requestheader-group-headers: X-Remote-Group
          requestheader-username-headers: X-Remote-User
          service-account-key-file: /etc/kubernetes/pki/sa.pub
          service-account-lookup: "true"
          tls-cert-file: /etc/kubernetes/pki/apiserver.crt
          tls-cipher-suites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256
          tls-private-key-file: /etc/kubernetes/pki/apiserver.key
        extraVolumes:
        - hostPath: /etc/kubernetes/encryption-provider.yaml
          mountPath: /etc/kubernetes/encryption-provider.yaml
          name: encryption-provider
      controllerManager:
        extraArgs:
          allocate-node-cidrs: "true"
          authentication-kubeconfig: /etc/kubernetes/controller-manager.conf
          authorization-kubeconfig: /etc/kubernetes/controller-manager.conf
          bind-address: 0.0.0.0
          cloud-provider: external
          cluster-signing-cert-file: /etc/kubernetes/pki/ca.crt
          cluster-signing-duration: 87600h0m0s
          cluster-signing-key-file: /etc/kubernetes/pki/ca.key
          kubeconfig: /etc/kubernetes/controller-manager.conf
          profiling: "false"
          requestheader-client-ca-file: /etc/kubernetes/pki/front-proxy-ca.crt
          root-ca-file: /etc/kubernetes/pki/ca.crt
          secure-port: "10257"
          service-account-private-key-file: /etc/kubernetes/pki/sa.key
          terminated-pod-gc-threshold: "10"
          use-service-account-credentials: "true"
      etcd:
        local:
          dataDir: /var/lib/etcd
          extraArgs:
            auto-tls: "false"
            cert-file: /etc/kubernetes/pki/etcd/server.crt
            cipher-suites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256
            client-cert-auth: "true"
            key-file: /etc/kubernetes/pki/etcd/server.key
            peer-auto-tls: "false"
            peer-client-cert-auth: "true"
            trusted-ca-file: /etc/kubernetes/pki/etcd/ca.crt
      scheduler:
        extraArgs:
          bind-address: 0.0.0.0
          kubeconfig: /etc/kubernetes/scheduler.conf
          profiling: "false"
          secure-port: "10259"
    files:
    - content: |
        apiVersion: apiserver.config.k8s.io/v1
        kind: EncryptionConfiguration
        resources:
          - resources:
            - secrets
            providers:
            - aescbc:
                keys:
                - name: key1
                  secret: 8d7iAcg3/NwN9aijhtEXj5kL2NOHIgokGFjbIBfL6X0=
            - identity: {}
      owner: root:root
      path: /etc/kubernetes/encryption-provider.yaml
      permissions: "0600"
    - content: |
        net.ipv4.conf.lxc*.rp_filter = 0
      owner: root:root
      path: /etc/sysctl.d/99-cilium.conf
      permissions: "0744"
    - content: |
        overlay
        br_netfilter
      owner: root:root
      path: /etc/modules-load.d/crio.conf
      permissions: "0744"
    - content: |
        net.bridge.bridge-nf-call-iptables  = 1
        net.bridge.bridge-nf-call-ip6tables = 1
        net.ipv4.ip_forward                 = 1
      owner: root:root
      path: /etc/sysctl.d/99-kubernetes-cri.conf
      permissions: "0744"
    - content: |
        vm.overcommit_memory=1
        kernel.panic=10
        kernel.panic_on_oops=1
      owner: root:root
      path: /etc/sysctl.d/99-kubelet.conf
      permissions: "0744"
    - content: |
        nameserver 1.1.1.1
        nameserver 1.0.0.1
        nameserver 2606:4700:4700::1111
      owner: root:root
      path: /etc/kubernetes/resolv.conf
      permissions: "0744"
    initConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          anonymous-auth: "false"
          authentication-token-webhook: "true"
          authorization-mode: Webhook
          cloud-provider: external
          event-qps: "5"
          kubeconfig: /etc/kubernetes/kubelet.conf
          max-pods: "120"
          read-only-port: "0"
          resolv-conf: /etc/kubernetes/resolv.conf
          rotate-server-certificates: "true"
          tls-cipher-suites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256
    joinConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          anonymous-auth: "false"
          authentication-token-webhook: "true"
          authorization-mode: Webhook
          cloud-provider: external
          event-qps: "5"
          kubeconfig: /etc/kubernetes/kubelet.conf
          max-pods: "120"
          read-only-port: "0"
          resolv-conf: /etc/kubernetes/resolv.conf
          rotate-server-certificates: "true"
          tls-cipher-suites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256
    preKubeadmCommands:
    - set -x
    - export CONTAINERD=1.7.13
    - export KUBERNETES_VERSION=$(echo {{.KUBERNETES_VERSION}} | sed 's/^v//')
    - export TRIMMED_KUBERNETES_VERSION=$(echo {{.KUBERNETES_VERSION}} | sed 's/^v//'
      | awk -F . '{print $1 "." $2}')
    - ARCH=amd64
    - if [ "$(uname -m)" = "aarch64" ]; then ARCH=arm64; fi
    - localectl set-locale LANG=en_US.UTF-8
    - localectl set-locale LANGUAGE=en_US.UTF-8
    - apt-get update -y
    - apt-get -y install at jq unzip wget socat mtr logrotate apt-transport-https
    - sed -i '/swap/d' /etc/fstab
    - swapoff -a
    - modprobe overlay && modprobe br_netfilter && sysctl --system
    - wget https://github.com/containerd/containerd/releases/download/v$CONTAINERD/cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz
    - wget https://github.com/containerd/containerd/releases/download/v$CONTAINERD/cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz.sha256sum
    - sha256sum --check cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz.sha256sum
    - tar --no-overwrite-dir -C / -xzf cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz
    - rm -f cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz cri-containerd-cni-$CONTAINERD-linux-$ARCH.tar.gz.sha256sum
    - chmod -R 644 /etc/cni && chown -R root:root /etc/cni
    - mkdir -p /etc/containerd
    - containerd config default > /etc/containerd/config.toml
    - sed -i  "s/SystemdCgroup = false/SystemdCgroup = true/" /etc/containerd/config.toml
    - systemctl daemon-reload && systemctl enable containerd && systemctl start containerd
    - curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    - echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
    - apt-get update
    - apt-get install -y kubelet=$KUBERNETES_VERSION-00 kubeadm=$KUBERNETES_VERSION-00
      kubectl=$KUBERNETES_VERSION-00  bash-completion && apt-mark hold kubelet kubectl
      kubeadm && systemctl enable kubelet
    - kubeadm config images pull --kubernetes-version $KUBERNETES_VERSION
    - echo 'source <(kubectl completion bash)' >>~/.bashrc
    - echo 'export KUBECONFIG=/etc/kubernetes/admin.conf' >>~/.bashrc
    - apt-get -y autoremove && apt-get -y clean all
  machineTemplate:
    infrastructureRef:
      apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
      kind: HCloudMachineTemplate
      name: {{.CLUSTER_NAME}}-control-plane
  replicas: {{.CONTROL_PLANE_MACHINE_COUNT}}
  version: {{.KUBERNETES_VERSION}}
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: HCloudMachineTemplate
metadata:
  name: {{.CLUSTER_NAME}}-control-plane
spec:
  template:
    spec:
      imageName: ubuntu-22.04
      placementGroupName: control-plane
      type: {{.HCLOUD_CONTROL_PLANE_MACHINE_TYPE}}
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: HCloudMachineTemplate
metadata:
  name: {{.CLUSTER_NAME}}-md-0
spec:
  template:
    spec:
      imageName: ubuntu-22.04
      placementGroupName: md-0
      type: {{.HCLOUD_WORKER_MACHINE_TYPE}}
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: HCloudRemediationTemplate
metadata:
  name: control-plane-remediation-request
spec:
  template:
    spec:
      strategy:
        retryLimit: 1
        timeout: 180s
        type: Reboot
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: HCloudRemediationTemplate
metadata:
  name: worker-remediation-request
spec:
  template:
    spec:
      strategy:
        retryLimit: 1
        timeout: 180s
        type: Reboot
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: HetznerCluster
metadata:
  name: {{.CLUSTER_NAME}}
spec:
  controlPlaneEndpoint:
    host: ""
    port: 443
  controlPlaneLoadBalancer:
    region: {{.HCLOUD_REGION}}
  controlPlaneRegions:
  - {{.HCLOUD_REGION}}
  hcloudNetwork:
    enabled: false
  hcloudPlacementGroups:
  - name: control-plane
    type: spread
  - name: md-0
    type: spread
  hetznerSecretRef:
    key:
      hcloudToken: hcloud
      hetznerRobotPassword: robot-password
      hetznerRobotUser: robot-user
    name: hetzner
  sshKeys:
    hcloud:
    - name: {{.HCLOUD_SSH_KEY}}`

	hetznerCCMCSI = `
manifest:
apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: hcloud-ccm
  namespace: default
spec:
  clusterSelector:
    matchLabels:
      cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  repoURL: https://charts.syself.com
  chartName: ccm-hcloud
  version: 1.0.11
  releaseName: ccm
  namespace: kube-system
  valuesTemplate: |
                secret:
                  name: hetzner
                  tokenKeyName: hcloud
---
apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: hcloud-csi
  namespace: default
spec:
  clusterSelector:
    matchLabels:
      cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  repoURL: https://charts.hetzner.cloud/
  chartName: hcloud-csi
  version: 2.5.1
  releaseName: csi
  namespace: kube-system
  valuesTemplate: |
                controller:
                  hcloudToken:
                    existingSecret:
                      name: hetzner
                      key: hcloud
                node:
                  affinity:
                    nodeAffinity:
                      requiredDuringSchedulingIgnoredDuringExecution:
                        nodeSelectorTerms:
                        - matchExpressions:
                          - key: "node-role.kubernetes.io/worker"
                            operator: Exists
                        - matchExpressions:
                          - key: "node-role.kubernetes.io/control-plane"
                            operator: Exists

`

	hetznerSecret = `apiVersion: v1
kind: Secret
type: Opaque 
metadata:
  name: hetzner
  namespace: default
  labels: 
    clusterctl.cluster.x-k8s.io/move: ""
data:
  hcloud: {{.HCLOUD_TOKEN}}
  robot-user: {{.HETZNER_ROBOT_USER}}
  robot-password: {{.HETZNER_ROBOT_PASSWORD}}`

	hetznerSSHSecret = `apiVersion: v1
kind: Secret
type: Opaque 
metadata:
  name: robot-ssh
  namespace: default
  labels: 
    clusterctl.cluster.x-k8s.io/move: ""
stringData:
  sshkey-name: {{.HCLOUD_SSH_KEY}}
data:
  ssh-privatekey: {{.HETZNER_SSH_PRIVATE_KEY}}
  ssh-publickey: {{.HETZNER_SSH_PUBLIC_KEY}}`
)
