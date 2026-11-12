FROM localhost/base-builder:dev

LABEL io.k8s.description="S2I builder image for Go (go build / Makefile)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

# GOTOOLCHAIN=local: never download a different toolchain at build time (reproducible, offline-friendly)
ENV GOPATH=/opt/app-root/go \
    GOCACHE=/opt/app-root/.cache/go-build \
    GOTOOLCHAIN=local \
    GOFLAGS=-mod=readonly

COPY ./s2i/bin/ /usr/libexec/s2i

RUN microdnf install -y \
        gcc \
        gcc-c++ \
        golang \
        make && \
    microdnf clean all && \
    mkdir -p /opt/app-root/go /opt/app-root/.cache && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
