FROM rockylinux:9-minimal

LABEL io.k8s.description="Common S2I base builder image" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root

RUN microdnf update -y && \
    microdnf install -y \
        ca-certificates \
        curl \
        git \
        gzip \
        tar \
        unzip \
        xz && \
    microdnf clean all && \
    mkdir -p /opt/app-root/src /tmp/src /usr/libexec/s2i && \
    chown -R 1001:0 /opt/app-root /tmp/src /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root /tmp/src /usr/libexec/s2i

WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
