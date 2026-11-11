FROM rockylinux:9-minimal

LABEL io.k8s.description="S2I builder image for Node.js (npm/yarn/pnpm)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root \
    NODE_VERSION=24.22.2 \
    NPM_CONFIG_UPDATE_NOTIFIER=false \
    PATH="/usr/local/bin:${PATH}"

COPY ./s2i/bin/ /usr/libexec/s2i

# python3/make/g++ are needed by node-gyp to compile native addons (builder only, never in runtime)
RUN microdnf update -y && \
    microdnf install -y \
        ca-certificates \
        curl \
        gcc \
        gcc-c++ \
        git \
        make \
        python3 \
        tar \
        xz && \
    microdnf clean all && \
    curl -fsSL "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz" -o /tmp/node.tar.xz && \
    tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1 && \
    rm -f /tmp/node.tar.xz && \
    mkdir -p /opt/app-root/src /tmp/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /tmp/src /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root /tmp/src

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
