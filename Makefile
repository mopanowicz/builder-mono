CONTAINER_ENGINE ?= $(shell command -v podman 2>/dev/null || echo docker)
JAVA_TEST_APPS   := java-maven java-gradle

.PHONY: build build-java test-examples test-java $(addprefix test-,$(JAVA_TEST_APPS))

build-java:
	$(CONTAINER_ENGINE) build -t localhost/java-builder:dev -f s2i/java-builder/Dockerfile s2i/java-builder/
	$(CONTAINER_ENGINE) build -t localhost/java-runtime:dev -f s2i/java-runtime/Dockerfile s2i/java-runtime/

build: build-java
	$(CONTAINER_ENGINE) build -t localhost/nodejs-builder:dev -f s2i/nodejs-builder/Dockerfile s2i/nodejs-builder/
	$(CONTAINER_ENGINE) build -t localhost/nodejs-runtime:dev -f s2i/nodejs-runtime/Dockerfile s2i/nodejs-runtime/

$(addprefix test-,$(JAVA_TEST_APPS)): test-%: build-java
	@echo "===> [1/4] Building S2I image for $* (java-builder -> java-runtime)"
	@trap 'rm -rf Dockerfile.gen upload' EXIT; \
	s2i build s2i-tests/$* localhost/java-builder:dev localhost/$*-test:dev \
		--as-dockerfile Dockerfile.gen && \
	printf '\nFROM localhost/java-runtime:dev\nCOPY --from=0 --chown=1001:0 /opt/app-root/src/app.jar /opt/app-root/src/app.jar\nUSER 1001\nCMD ["/usr/libexec/s2i/run"]\n' >> Dockerfile.gen && \
	$(CONTAINER_ENGINE) build -t localhost/$*-test:dev -f Dockerfile.gen .
	@echo "===> [2/4] Verifying non-root execution (UID 1001)"
	@test "$$($(CONTAINER_ENGINE) run --rm --entrypoint id localhost/$*-test:dev -u)" = "1001"
	@echo "===> [3/4] Verifying no build tools in runtime image"
	@! $(CONTAINER_ENGINE) run --rm --entrypoint sh localhost/$*-test:dev -c "command -v mvn || command -v gradle || command -v javac" >/dev/null 2>&1
	@echo "===> [4/4] Verifying HTTP response on port 8080"
	@cid=$$($(CONTAINER_ENGINE) run -d -p 8080:8080 localhost/$*-test:dev); \
	trap "$(CONTAINER_ENGINE) rm -f $$cid >/dev/null 2>&1 || true" EXIT; \
	if ! curl -fsS --retry 15 --retry-delay 1 --retry-all-errors http://localhost:8080/actuator/health | grep -q '"status":"UP"'; then \
		echo "ERROR: Health check failed. Container logs:"; \
		$(CONTAINER_ENGINE) logs $$cid; \
		exit 1; \
	fi
	@echo "===> PASS: $*"

test-java: $(addprefix test-,$(JAVA_TEST_APPS))

test-all: test-java
