# This node image is only for usage with kind v0.12.0
FROM kindest/node:v1.22.7 as kube
ARG TARGETPLATFORM
RUN case ${TARGETPLATFORM} in \
         "linux/amd64")  export ARCH=amd64  ;; \
         "linux/arm64")  export ARCH=arm64  ;; \
         "linux/arm/v7") export ARCH=arm  ;; \
         "linux/386")    export ARCH=i386   ;; \
    esac && \
    curl -L https://github.com/edgefarm/edgefarm/releases/download/cni-0.8.0/cni-plugins-linux-$ARCH-v0.8.0.tgz -o /tmp/cni-plugins.tgz && \
    cd /opt/cni/bin/ && \
    tar xfz /tmp/cni-plugins.tgz && \
    rm /tmp/cni-plugins.tgz

# use kind node v1.24.0 as it contains systemd as default cgroup driver
FROM kindest/node:v1.24.0
COPY --from=kube /opt/cni/ /opt/cni
COPY --from=kube /usr/bin/kubelet /usr/bin/kubelet
COPY --from=kube /usr/bin/kubeadm /usr/bin/kubeadm
COPY --from=kube /usr/bin/kubectl /usr/bin/kubectl
