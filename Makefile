CONTAINER_ENGINE ?= $(shell command -v podman 2>/dev/null || echo docker)

.PHONY: build test-examples

build:
	$(CONTAINER_ENGINE) build -t localhost/java-builder:dev -f s2i/java-builder/Dockerfile s2i/java-builder/
	$(CONTAINER_ENGINE) build -t localhost/java-runtime:dev -f s2i/java-runtime/Dockerfile s2i/java-runtime/
	$(CONTAINER_ENGINE) build -t localhost/nodejs-builder:dev -f s2i/nodejs-builder/Dockerfile s2i/nodejs-builder/
	$(CONTAINER_ENGINE) build -t localhost/nodejs-runtime:dev -f s2i/nodejs-runtime/Dockerfile s2i/nodejs-runtime/

test-examples:
	./hack/test-examples.sh

