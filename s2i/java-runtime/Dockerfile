FROM rockylinux:9-minimal

LABEL io.k8s.description="S2I runtime image for Java (JRE only)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root \
    JAVA_HOME=/usr/lib/jvm/java-25-openjdk \
    PATH="${JAVA_HOME}/bin:${PATH}"

COPY ./s2i/bin/ /usr/libexec/s2i

RUN microdnf update -y && \
    microdnf install -y \
        ca-certificates \
        java-25-openjdk-headless && \
    microdnf clean all && \
    mkdir -p /opt/app-root/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root

USER 1001
WORKDIR /opt/app-root/src
EXPOSE 8080
CMD ["/usr/libexec/s2i/run"]
