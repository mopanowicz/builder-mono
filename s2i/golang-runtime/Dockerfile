FROM rockylinux:9-minimal

LABEL io.k8s.description="S2I runtime image for Go (compiled binary only, no toolchain)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root

COPY ./s2i/bin/ /usr/libexec/s2i

# ca-certificates: needed by Go binaries making outbound TLS connections
RUN microdnf update -y && \
    microdnf install -y ca-certificates && \
    microdnf clean all && \
    mkdir -p /opt/app-root/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root

USER 1001
WORKDIR /opt/app-root/src
EXPOSE 8080
CMD ["/usr/libexec/s2i/run"]
