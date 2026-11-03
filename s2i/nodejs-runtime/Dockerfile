FROM node:20-slim

LABEL io.k8s.description="S2I runtime image for Node.js (node + prod deps)" \
      io.openshift.s2i.scripts-url="image:///usr/libexec/s2i" \
      maintainer="platform-team"

COPY ./s2i/bin/ /usr/libexec/s2i
RUN chmod +x /usr/libexec/s2i/*

USER 1001
WORKDIR /opt/app-root/src
EXPOSE 3000
CMD ["/usr/libexec/s2i/run"]

