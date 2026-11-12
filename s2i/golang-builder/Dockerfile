FROM localhost/base-builder:dev

LABEL io.k8s.description="S2I builder image for Go (go build / Makefile)"

# GOTOOLCHAIN=local: never download a different toolchain at build time (reproducible, offline-friendly)
ENV GOPATH=/opt/app-root/go \
    GOCACHE=/opt/app-root/.cache/go-build \
    GOTOOLCHAIN=local \
    GOFLAGS=-mod=readonly

COPY --chown=1001:0 ./s2i/bin/ /usr/libexec/s2i

RUN microdnf install -y \
        gcc \
        gcc-c++ \
        golang \
        make && \
    microdnf clean all && \
    chmod +x /usr/libexec/s2i/*

USER 1001
