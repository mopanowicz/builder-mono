FROM node:20-slim

LABEL io.k8s.description="S2I runtime image for Node.js (node + prod deps)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root \
    NODE_ENV=production

COPY ./s2i/bin/ /usr/libexec/s2i

# Runtime only needs `node`: remove package managers shipped with the base image
RUN rm -rf /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack \
           /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/corepack \
           /usr/local/bin/yarn /usr/local/bin/yarnpkg /opt/yarn-* && \
    mkdir -p /opt/app-root/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root

USER 1001
WORKDIR /opt/app-root/src
EXPOSE 3000
CMD ["/usr/libexec/s2i/run"]
