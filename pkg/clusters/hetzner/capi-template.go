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
          version = 2
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
            runtime_type = "io.containerd.runc.v2"
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
            SystemdCgroup = true
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun]
            runtime_type = "io.containerd.runc.v2"
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun.options]
            BinaryName = "crun"
            Root = "/usr/local/sbin"
            SystemdCgroup = true
          [plugins."io.containerd.grpc.v1.cri".containerd]
            default_runtime_name = "crun"
          [plugins."io.containerd.runtime.v1.linux"]
            runtime = "crun"
            runtime_root = "/usr/local/sbin"
        owner: root:root
        path: /etc/containerd/config.toml
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
      - contentFrom:
          secret:
            name: netbird
            key: setupKey
        owner: root:root
        path: /etc/netbird/setup-key
        permissions: "0744"
      - contentFrom:
          secret:
            name: netbird
            key: domain
        owner: root:root
        path: /etc/netbird/domain
        permissions: "0744"
      - contentFrom:
          secret:
            name: netbird
            key: admin-url
        owner: root:root
        path: /etc/netbird/admin-url
        permissions: "0744"
      - contentFrom:
          secret:
            name: netbird
            key: management-url
        owner: root:root
        path: /etc/netbird/management-url
        permissions: "0744"  
      joinConfiguration:
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "node.edgefarm.io/type=cloud,ingress-ready=true,node-role.edgefarm.io/worker="
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
      - export CRUN=1.8.5
      - export CONTAINERD=1.7.13
      - export KUBERNETES_VERSION=$(echo {{.KUBERNETES_VERSION}} | sed 's/^v//')
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
      - wget https://github.com/containers/crun/releases/download/$CRUN/crun-$CRUN-linux-$ARCH -O /usr/local/sbin/crun && chmod +x /usr/local/sbin/crun
      - rm -f /etc/cni/net.d/10-containerd-net.conflist
      - chmod -R 644 /etc/cni && chown -R root:root /etc/cni
      - systemctl daemon-reload && systemctl enable containerd && systemctl start containerd
      - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubeadm_1.22.17-00_amd64_7b7456beaf364ecf5c14f4d995bc49985cd23273ebf7610717961e2575057209.deb -o /tmp/kubeadm_1.22.17-00.deb
      - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubectl_1.22.17-00_amd64_b3bcd8e4a64fded2873e873301ef68c6c3787dbc5e68f079a2f9c7c283180709.deb -o /tmp/kubectl_1.22.17-00.deb
      - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubelet_1.22.17-00_amd64_3488568197f82b8b8c267058ea7165968560a67daa5cea981ac6bcff43fe0966.deb -o /tmp/kubelet_1.22.17-00.deb
      - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubernetes-cni_1.2.0-00_amd64_0c2be3775ea591dee9ce45121341dd16b3c752763c6898adc35ce12927c977c1.deb -o /tmp/kubernetes-cni_1.2.0-00.deb
      - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/cri-tools_1.26.0-00_amd64_5ba786e8853986c7f9f51fe850086083e5cf3c3d34f3fc09aaadd63fa0b578df.deb -o /tmp/cri-tools_1.26.0-00.deb
      - apt install -y --fix-broken /tmp/kubeadm_1.22.17-00.deb /tmp/kubectl_1.22.17-00.deb /tmp/kubelet_1.22.17-00.deb /tmp/kubernetes-cni_1.2.0-00.deb /tmp/cri-tools_1.26.0-00.deb
      - kubectl=$KUBERNETES_VERSION-00 bash-completion && apt-mark hold kubelet kubectl kubeadm
      - kubeadm && systemctl enable kubelet
      - kubeadm config images pull --kubernetes-version $KUBERNETES_VERSION
      - echo 'source <(kubectl completion bash)' >>~/.bashrc
      - echo 'export KUBECONFIG=/etc/kubernetes/admin.conf' >>~/.bashrc
      - curl -L https://pkgs.wiretrustee.com/debian/public.key | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/wiretrustee.gpg
      - echo 'deb https://pkgs.wiretrustee.com/debian stable main' | sudo tee /etc/apt/sources.list.d/wiretrustee.list
      - apt-get update
      - sudo apt-get install netbird=0.24.3 -y
      - netbird up --setup-key $(cat /etc/netbird/setup-key) --admin-url $(cat /etc/netbird/admin-url) --management-url $(cat /etc/netbird/management-url)
      - apt-get -y autoremove && apt-get -y clean all && rm -rf /var/lib/apt/lists/*
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
        version = 2
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          runtime_type = "io.containerd.runc.v2"
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
          SystemdCgroup = true
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun]
          runtime_type = "io.containerd.runc.v2"
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun.options]
          BinaryName = "crun"
          Root = "/usr/local/sbin"
          SystemdCgroup = true
        [plugins."io.containerd.grpc.v1.cri".containerd]
          default_runtime_name = "crun"
        [plugins."io.containerd.runtime.v1.linux"]
          runtime = "crun"
          runtime_root = "/usr/local/sbin"
      owner: root:root
      path: /etc/containerd/config.toml
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
    - contentFrom:
        secret:
          name: netbird
          key: setupKey
      owner: root:root
      path: /etc/netbird/setup-key
      permissions: "0744"
    - contentFrom:
        secret:
          name: netbird
          key: domain
      owner: root:root
      path: /etc/netbird/domain
      permissions: "0744"
    - contentFrom:
        secret:
          name: netbird
          key: admin-url
        owner: root:root
      path: /etc/netbird/admin-url
      permissions: "0744"
    - contentFrom:
        secret:
          name: netbird
          key: management-url
        owner: root:root
      path: /etc/netbird/management-url
      permissions: "0744"
    - content: |
        #!/bin/bash
        export DNS_SEARCH_DOMAIN=$(cat /etc/netbird/domain)
        MANIFEST="/etc/kubernetes/manifests/kube-apiserver.yaml"
        while true; do
            if [ -f "$MANIFEST" ]; then
                break
            else
                echo "$MANIFEST not found yet. Waiting..."
                exit 1
            fi
        done
        if yq eval '(.spec.dnsPolicy == "None") and (.spec.dnsConfig.nameservers[0] == "127.0.0.53") and (.spec.dnsConfig.searches[0] == env(DNS_SEARCH_DOMAIN))' $MANIFEST | grep -q "true"; then
            echo "DNS settings are present in the kube-apiserver manifest."
            systemctl stop patch-kube-apiserver.service
            systemctl disable patch-kube-apiserver.service
            exit 0
        else
            echo "DNS settings are not present in the kube-apiserver manifest. Applying..."
            yq eval '.spec.dnsPolicy = "None" | .spec.dnsConfig.nameservers[0] = "127.0.0.53" | .spec.dnsConfig.searches[0] = env(DNS_SEARCH_DOMAIN)' -i $MANIFEST
        fi
        systemctl stop patch-kube-apiserver.service
        systemctl disable patch-kube-apiserver.service
        exit 0
      owner: root:root
      path: /usr/local/bin/patch-kube-apiserver-dns.sh
      permissions: "0755"
    - content: |
        [Unit]
        Description=Patch kube-apiserver static manifest
        After=kubelet.target

        [Service]
        ExecStart=/usr/local/bin/patch-kube-apiserver-dns.sh
        Restart=always
        RestartSec=30

        [Install]
        WantedBy=multi-user.target
      owner: root:root
      path: /etc/systemd/system/patch-kube-apiserver.service
      permissions: "0644"      
    initConfiguration:
      skipPhases:
        - addon/coredns
        - addon/kube-proxy
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "node.edgefarm.io/type=cloud"
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
          node-labels: "node.edgefarm.io/type=cloud"
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
    - export CRUN=1.9.2
    - export CONTAINERD=1.7.6
    - export KUBERNETES_VERSION=$(echo {{.KUBERNETES_VERSION}} | sed 's/^v//')
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
    - wget https://github.com/containers/crun/releases/download/$CRUN/crun-$CRUN-linux-$ARCH -O /usr/local/sbin/crun && chmod +x /usr/local/sbin/crun
    - rm -f /etc/cni/net.d/10-containerd-net.conflist
    - chmod -R 644 /etc/cni && chown -R root:root /etc/cni
    - systemctl daemon-reload && systemctl enable containerd && systemctl start containerd
    - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubeadm_1.22.17-00_amd64_7b7456beaf364ecf5c14f4d995bc49985cd23273ebf7610717961e2575057209.deb -o /tmp/kubeadm_1.22.17-00.deb
    - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubectl_1.22.17-00_amd64_b3bcd8e4a64fded2873e873301ef68c6c3787dbc5e68f079a2f9c7c283180709.deb -o /tmp/kubectl_1.22.17-00.deb
    - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubelet_1.22.17-00_amd64_3488568197f82b8b8c267058ea7165968560a67daa5cea981ac6bcff43fe0966.deb -o /tmp/kubelet_1.22.17-00.deb
    - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/kubernetes-cni_1.2.0-00_amd64_0c2be3775ea591dee9ce45121341dd16b3c752763c6898adc35ce12927c977c1.deb -o /tmp/kubernetes-cni_1.2.0-00.deb
    - curl -L https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb/cri-tools_1.26.0-00_amd64_5ba786e8853986c7f9f51fe850086083e5cf3c3d34f3fc09aaadd63fa0b578df.deb -o /tmp/cri-tools_1.26.0-00.deb
    - apt install -y --fix-broken /tmp/kubeadm_1.22.17-00.deb /tmp/kubectl_1.22.17-00.deb /tmp/kubelet_1.22.17-00.deb /tmp/kubernetes-cni_1.2.0-00.deb /tmp/cri-tools_1.26.0-00.deb
    - kubectl=$KUBERNETES_VERSION-00 bash-completion && apt-mark hold kubelet kubectl kubeadm
    - kubeadm && systemctl enable kubelet
    - kubeadm config images pull --kubernetes-version $KUBERNETES_VERSION
    - echo 'source <(kubectl completion bash)' >>~/.bashrc
    - echo 'export KUBECONFIG=/etc/kubernetes/admin.conf' >>~/.bashrc
    - sudo apt install ca-certificates curl gnupg -y
    - curl -L https://pkgs.wiretrustee.com/debian/public.key | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/wiretrustee.gpg
    - echo 'deb https://pkgs.wiretrustee.com/debian stable main' | sudo tee /etc/apt/sources.list.d/wiretrustee.list
    - apt-get update
    - sudo apt-get install netbird=0.24.3 -y
    - netbird up --setup-key $(cat /etc/netbird/setup-key) --admin-url $(cat /etc/netbird/admin-url) --management-url $(cat /etc/netbird/management-url)
    - curl -L https://github.com/mikefarah/yq/releases/download/v4.35.2/yq_linux_${ARCH} -o /usr/local/bin/yq
    - chmod +x /usr/local/bin/yq
    - systemctl enable patch-kube-apiserver.service
    - systemctl start patch-kube-apiserver.service
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

	hetznerCCMCSI = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
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
              - key: "node.edgefarm.io/type"
                operator: In
                values:
                  - "cloud"`

	hetznerSecret = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: hetzner
  namespace: default
  labels:
    clusterctl.cluster.x-k8s.io/move: ""
stringData:
  hcloud: {{.HCLOUD_TOKEN}}`

	netbirdSecret = `apiVersion: v1
kind: Secret
metadata: 
  name: netbird
  namespace: default
  labels:
    clusterctl.cluster.x-k8s.io/move: ""
type: Opaque
data: 
  admin-url: {{.NETBIRD_ADMIN_URL}}
  domain: {{.NETBIRD_DOMAIN}}
  management-url: {{.NETBIRD_MANAGEMENT_URL}}
  setupKey: {{.NETBIRD_SETUP_KEY}}`

	hetznerSSHSecret = `apiVersion: v1
kind: Secret
metadata:
  name: robot-ssh
  namespace: default
  labels:
    clusterctl.cluster.x-k8s.io/move: ""
type: Opaque
stringData:
  sshkey-name: {{.HCLOUD_SSH_KEY}}`

	flannelCloud = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: flannel-cloud
  namespace: default
spec:
  clusterSelector:
    matchLabels:
      cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  version: 1.21.0
  namespace: kube-flannel
  chartName: flannel
  releaseName: flannel-cloud
  repoURL: oci://ghcr.io/edgefarm/helm-charts
  valuesTemplate: |
    nameOverride: flannel-cloud
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
      args: ["--ip-masq", "--kube-subnet-mgr", "--iface=wt0"]
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
                    - "cloud"`

	flannelEdge = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: flannel-edge
  namespace: default
spec:
  clusterSelector:
    matchLabels:
      cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  version: 1.21.0
  namespace: kube-flannel
  chartName: flannel
  releaseName: flannel-edge
  repoURL: oci://ghcr.io/edgefarm/helm-charts
  valuesTemplate: |
    nameOverride: flannel-edge
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
                  - "physical"`

	kubeProxyOpenYurt = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: kube-proxy-openyurt
  namespace: default
spec:
  clusterSelector:
    matchLabels:
    cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  version: 1.23.1
  releaseName: kube-proxy-openyurt
  namespace: kube-system
  chartName: kube-proxy
  repoURL: oci://ghcr.io/edgefarm/helm-charts
  valuesTemplate: |
    nameOverride: kube-proxy-openyurt
    kuberneteServerAddr:
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
    features:
      openyurt:
        enabled: true
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
                operator: In
                values:
                  - "true"`

	kubeProxyDefault = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: kube-proxy-default
  namespace: default
spec:
  clusterSelector:
    matchLabels:
    cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  version: 1.23.1
  releaseName: kube-proxy-default
  namespace: kube-system
  chartName: kube-proxy
  repoURL: oci://ghcr.io/edgefarm/helm-charts
  valuesTemplate: |
    nameOverride: kube-proxy-default
    kuberneteServerAddr:
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
                  operator: Exists`

	coreDNS = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: coredns
  namespace: default
spec:
  clusterSelector:
    matchLabels:
    cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  version: 1.16.0
  namespace: kube-system
  chartName: coredns
  repoURL: oci://ghcr.io/edgefarm/helm-charts
  valuesTemplate: |
    features:
      log:
        enabled: true
    dnsPolicy: None
    dnsConfig:
      nameservers:
        - 8.8.8.8
        - 8.8.4.4`

	localPathProvisioner = `apiVersion: addons.cluster.x-k8s.io/v1alpha1
kind: HelmChartProxy
metadata:
  name: local-path-provisioner
  namespace: default
spec:
  clusterSelector:
    matchLabels:
    cloudprovider.clusters.infrastructure.edgefarm.io/type: hetzner
  version: 1.25.0
  namespace: local-path-provisioner
  chartName: local-path-provisioner
  repoURL: oci://ghcr.io/edgefarm/helm-charts`
)
