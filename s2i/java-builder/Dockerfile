FROM eclipse-temurin:21-jdk-jammy

LABEL io.k8s.description="S2I builder image for Java (Maven/Gradle)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV HOME=/opt/app-root \
    GRADLE_VERSION=8.10.2 \
    PATH="/opt/gradle/bin:${PATH}"

COPY ./s2i/bin/ /usr/libexec/s2i

RUN apt-get update && \
    apt-get install -y --no-install-recommends maven unzip curl ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    curl -fsSL "https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip" -o /tmp/gradle.zip && \
    unzip -q /tmp/gradle.zip -d /opt && \
    mv "/opt/gradle-${GRADLE_VERSION}" /opt/gradle && \
    rm -f /tmp/gradle.zip && \
    mkdir -p /opt/app-root/src /tmp/src && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /opt/app-root /tmp/src /usr/libexec/s2i && \
    chmod -R g=u /opt/app-root /tmp/src

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
