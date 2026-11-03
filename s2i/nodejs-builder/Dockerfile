FROM node:20-slim

LABEL io.k8s.description="S2I builder image for Node.js (npm/yarn/pnpm)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root \
    COREPACK_ENABLE_DOWNLOAD_PROMPT=0 \
    NPM_CONFIG_UPDATE_NOTIFIER=false

COPY ./s2i/bin/ /usr/libexec/s2i

# python3/make/g++ are needed by node-gyp to compile native addons (builder only, never in runtime)
RUN apt-get update && \
    apt-get install -y --no-install-recommends python3 make g++ ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    rm -f /usr/local/bin/yarn /usr/local/bin/yarnpkg && \
    corepack enable && \
    mkdir -p /opt/app-root/src /tmp/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /tmp/src /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root /tmp/src

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
