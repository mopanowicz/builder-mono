FROM localhost/base-builder:dev

LABEL io.k8s.description="S2I builder image for Java (Maven/Gradle)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

ENV GRADLE_VERSION=9.8.2 \
    MAVEN_VERSION=3.10.0 \
    JAVA_HOME=/usr/lib/jvm/java-25-openjdk \
    PATH="/opt/maven/bin:/opt/gradle/bin:${JAVA_HOME}/bin:${PATH}"

COPY ./s2i/bin/ /usr/libexec/s2i

RUN microdnf install -y \
        java-25-openjdk-devel && \
    microdnf clean all && \
    mkdir -p /opt/maven && \
    curl -fsSL "https://dlcdn.apache.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz" -o /tmp/maven.tar.gz && \
    curl -fsSL "https://downloads.apache.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz.sha512" -o /tmp/maven.tar.gz.sha512 && \
    echo "$(cut -d' ' -f1 /tmp/maven.tar.gz.sha512)  /tmp/maven.tar.gz" | sha512sum -c - && \
    tar -xzf /tmp/maven.tar.gz --strip-components=1 -C /opt/maven && \
    rm -f /tmp/maven.tar.gz /tmp/maven.tar.gz.sha512 && \
    curl -fsSL "https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip" -o /tmp/gradle.zip && \
    unzip -q /tmp/gradle.zip -d /opt && \
    mv "/opt/gradle-${GRADLE_VERSION}" /opt/gradle && \
    rm -f /tmp/gradle.zip && \
    chmod +x /usr/libexec/s2i/* && \
    chown -R 1001:0 /usr/libexec/s2i

USER 1001
WORKDIR /tmp/src
CMD ["/usr/libexec/s2i/usage"]
