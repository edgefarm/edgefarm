FROM docker.io/flannel/flannel-cni-plugin:v1.6.0-flannel1 as flannel-cni-plugin
FROM alpine:3.14 
ARG TARGETPLATFORM
RUN apk add curl
COPY --from=flannel-cni-plugin /flannel /cni/flannel
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then ARCHITECTURE=amd64; elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then ARCHITECTURE=arm64; fi \
    && curl -L https://github.com/edgefarm/cni-plugins/releases/download/v1.4.1-host-local/cni-plugins-linux-${ARCHITECTURE}-v1.4.1-host-local.tgz -o /tmp/cni-plugins.tgz && \
    tar -xvf /tmp/cni-plugins.tgz -C /cni/ && \
    rm /tmp/cni-plugins.tgz
