FROM cruizba/ubuntu-dind:19.03.11

ENV KUBEEDGE_VERSION=v1.8.2
ENV KUBEEDGE_ARCH=amd64

RUN curl -L https://github.com/kubeedge/kubeedge/releases/download/${KUBEEDGE_VERSION}/kubeedge-${KUBEEDGE_VERSION}-linux-${KUBEEDGE_ARCH}.tar.gz -o kubeedge-${KUBEEDGE_VERSION}-linux-${KUBEEDGE_ARCH}.tar.gz && \
    tar xfz kubeedge-${KUBEEDGE_VERSION}-linux-${KUBEEDGE_ARCH}.tar.gz && \
    cp kubeedge-${KUBEEDGE_VERSION}-linux-${KUBEEDGE_ARCH}/edge/edgecore /usr/local/bin/edgecore && \
    rm -rf kubeedge-${KUBEEDGE_VERSION}-linux-${KUBEEDGE_ARCH} kubeedge-${KUBEEDGE_VERSION}-linux-${KUBEEDGE_ARCH}.tar.gz

COPY dev/kubeedge-node/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN mkdir -p /etc/kubeedge/config
RUN mkdir -p /kubeedge

COPY dev/kubeedge-node/edgecore.yaml.TEMPLATE /kubeedge/
COPY dev/kubeedge-node/generate.sh /kubeedge/

RUN mkdir -p /etc/kubeedge/ca && \
    mkdir -p /etc/kubeedge/certs

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
