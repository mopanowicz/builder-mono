FROM eclipse-temurin:21-jre-jammy

LABEL io.k8s.description="S2I runtime image for Java (JRE only)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root

COPY ./s2i/bin/ /usr/libexec/s2i

RUN mkdir -p /opt/app-root/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root

USER 1001
WORKDIR /opt/app-root/src
EXPOSE 8080
CMD ["/usr/libexec/s2i/run"]
