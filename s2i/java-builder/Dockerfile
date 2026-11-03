FROM eclipse-temurin:21-jdk-jammy

LABEL io.k8s.description="S2I builder image for Java (Maven/Gradle)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

COPY ./s2i/bin/ /usr/libexec/s2i
RUN chmod +x /usr/libexec/s2i/*

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]

