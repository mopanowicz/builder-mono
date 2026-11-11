FROM rockylinux:9-minimal

LABEL io.k8s.description="S2I builder image for Go (go build / Makefile)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

# GOTOOLCHAIN=local: never download a different toolchain at build time (reproducible, offline-friendly)
ENV HOME=/opt/app-root \
    GOPATH=/opt/app-root/go \
    GOCACHE=/opt/app-root/.cache/go-build \
    GOTOOLCHAIN=local \
    GOFLAGS=-mod=readonly

COPY ./s2i/bin/ /usr/libexec/s2i

RUN microdnf update -y && \
    microdnf install -y \
        ca-certificates \
        gcc \
        gcc-c++ \
        git \
        golang \
        make && \
    microdnf clean all && \
    mkdir -p /opt/app-root/src /opt/app-root/go /opt/app-root/.cache /tmp/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /tmp/src /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root /tmp/src

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
