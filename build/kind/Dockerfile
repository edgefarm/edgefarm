# This node image is only for usage with kind v0.12.0
ARG VERSION
FROM kindest/node:v1.22.7 as kube
FROM kindest/node:v1.24.0
RUN curl -L https://github.com/edgefarm/edgefarm/releases/download/cni-0.8.0/cni-plugins-linux-amd64-v0.8.0.tgz -o /tmp/cni-plugins.tgz && \
    cd /opt/cni/bin/ && \
    tar xfz /tmp/cni-plugins.tgz && \
    rm /tmp/cni-plugins.tgz
COPY --from=kube /usr/bin/kubelet /usr/bin/kubelet
COPY --from=kube /usr/bin/kubeadm /usr/bin/kubeadm
COPY --from=kube /usr/bin/kubectl /usr/bin/kubectl
