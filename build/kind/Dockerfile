ARG VERSION
FROM kindest/node:${VERSION}
RUN curl -L https://github.com/containernetworking/plugins/releases/download/v1.3.0/cni-plugins-linux-amd64-v1.3.0.tgz -o /tmp/cni-plugins-linux-amd64-v1.3.0.tgz && \
    cd /opt/cni/bin/ && \
    tar xfz /tmp/cni-plugins-linux-amd64-v1.3.0.tgz && \
    rm /tmp/cni-plugins-linux-amd64-v1.3.0.tgz
