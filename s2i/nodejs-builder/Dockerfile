FROM localhost/base-builder:dev

LABEL io.k8s.description="S2I builder image for Node.js (npm/yarn/pnpm)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV NODE_VERSION=24.22.2 \
    COREPACK_ENABLE_DOWNLOAD_PROMPT=0 \
    NPM_CONFIG_UPDATE_NOTIFIER=false \
    PATH="/usr/local/bin:${PATH}"

COPY ./s2i/bin/ /usr/libexec/s2i

# python3/make/g++ are needed by node-gyp to compile native addons (builder only, never in runtime)
# Node.js 24+ no longer bundles corepack; install and enable it to provide yarn/pnpm shims (honoring package.json "packageManager")
RUN microdnf install -y \
        gcc \
        gcc-c++ \
        make \
        python3 && \
    microdnf clean all && \
    curl -fsSL "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz" -o /tmp/node.tar.xz && \
    tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1 && \
    rm -f /tmp/node.tar.xz && \
    npm install -g --no-audit --no-fund corepack && \
    corepack enable && \
    npm cache clean --force && \
    rm -rf /root/.npm /opt/app-root/.npm && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /usr/libexec/s2i

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
