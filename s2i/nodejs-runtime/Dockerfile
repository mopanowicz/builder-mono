FROM rockylinux:9-minimal

LABEL io.k8s.description="S2I runtime image for Node.js (node + prod deps)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root \
    NODE_ENV=production \
    PATH="/usr/local/bin:${PATH}"

COPY ./s2i/bin/ /usr/libexec/s2i

# Runtime only needs `node`: remove package managers and build tooling from the installed runtime image
RUN microdnf update -y && \
    microdnf install -y \
        ca-certificates \
        curl \
        tar \
        xz && \
    microdnf clean all && \
    curl -fsSL "https://nodejs.org/dist/v24.0.0/node-v24.0.0-linux-x64.tar.xz" -o /tmp/node-v24.tar.xz && \
    tar -xJf /tmp/node-v24.tar.xz -C /usr/local --strip-components=1 && \
    rm -f /tmp/node-v24.tar.xz && \
    rm -rf /usr/local/lib/node_modules /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack \
           /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/corepack \
           /usr/local/bin/yarn /usr/local/bin/yarnpkg /usr/local/bin/pnpm && \
    mkdir -p /opt/app-root/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root

USER 1001
WORKDIR /opt/app-root/src
EXPOSE 3000
CMD ["/usr/libexec/s2i/run"]
