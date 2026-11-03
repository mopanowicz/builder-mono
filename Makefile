CONTAINER_ENGINE ?= $(shell command -v podman 2>/dev/null || echo docker)
JAVA_TEST_APPS   := java-maven java-gradle
NODEJS_TEST_APPS := nodejs-npm nodejs-pnpm

JAVA_TESTS   := $(addprefix test-,$(JAVA_TEST_APPS))
NODEJS_TESTS := $(addprefix test-,$(NODEJS_TEST_APPS))

.PHONY: build build-java build-nodejs test test-java test-nodejs clean $(JAVA_TESTS) $(NODEJS_TESTS)
.NOTPARALLEL:

build: build-java build-nodejs

build-java:
	$(CONTAINER_ENGINE) build -t localhost/java-builder:dev -f s2i/java-builder/Dockerfile s2i/java-builder/
	$(CONTAINER_ENGINE) build -t localhost/java-runtime:dev -f s2i/java-runtime/Dockerfile s2i/java-runtime/

build-nodejs:
	$(CONTAINER_ENGINE) build -t localhost/nodejs-builder:dev -f s2i/nodejs-builder/Dockerfile s2i/nodejs-builder/
	$(CONTAINER_ENGINE) build -t localhost/nodejs-runtime:dev -f s2i/nodejs-runtime/Dockerfile s2i/nodejs-runtime/

# S2I end-to-end test: builder -> runtime build, UID 1001, no build tools, HTTP healthcheck
# $(1)=app  $(2)=lang (java|nodejs)  $(3)=artifact path  $(4)=port  $(5)=health path  $(6)=forbidden tools
define S2I_TEST
	@echo "===> [1/4] Building S2I image for $(1) ($(2)-builder -> $(2)-runtime)"
	@trap 'rm -rf Dockerfile.gen upload' EXIT; \
	s2i build s2i-tests/$(1) localhost/$(2)-builder:dev localhost/$(1)-test:dev \
		--as-dockerfile Dockerfile.gen && \
	printf '\nFROM localhost/$(2)-runtime:dev\nCOPY --from=0 --chown=1001:0 $(3) $(3)\nUSER 1001\nCMD ["/usr/libexec/s2i/run"]\n' >> Dockerfile.gen && \
	$(CONTAINER_ENGINE) build -t localhost/$(1)-test:dev -f Dockerfile.gen .
	@echo "===> [2/4] Verifying non-root execution (UID 1001)"
	@test "$$($(CONTAINER_ENGINE) run --rm --entrypoint id localhost/$(1)-test:dev -u)" = "1001"
	@echo "===> [3/4] Verifying no build tools in runtime image ($(6))"
	@if $(CONTAINER_ENGINE) run --rm --entrypoint sh localhost/$(1)-test:dev -c 'for t in $(6); do command -v "$$t" && exit 0; done; exit 1'; then \
		echo "ERROR: Build tool found in runtime image"; \
		exit 1; \
	fi
	@echo "===> [4/4] Verifying HTTP response on http://localhost:$(4)$(5)"
	@cid=$$($(CONTAINER_ENGINE) run -d -p $(4):$(4) localhost/$(1)-test:dev); \
	trap "$(CONTAINER_ENGINE) rm -f $$cid >/dev/null 2>&1 || true" EXIT; \
	if ! curl -fsS --retry 15 --retry-delay 1 --retry-all-errors http://localhost:$(4)$(5) | grep -q '"status":"UP"'; then \
		echo "ERROR: Health check failed. Container logs:"; \
		$(CONTAINER_ENGINE) logs $$cid; \
		exit 1; \
	fi
	@echo "===> PASS: $(1)"
endef

$(JAVA_TESTS): test-%: build-java
	$(call S2I_TEST,$*,java,/opt/app-root/src/app.jar,8080,/actuator/health,mvn gradle javac)

$(NODEJS_TESTS): test-%: build-nodejs
	$(call S2I_TEST,$*,nodejs,/opt/app-root/src,3000,/health,npm npx yarn pnpm corepack gcc g++ make)

test-java: $(JAVA_TESTS)

test-nodejs: $(NODEJS_TESTS)

test: test-java test-nodejs

clean:
	rm -rf Dockerfile.gen upload
	-$(CONTAINER_ENGINE) rmi $(patsubst %,localhost/%-test:dev,$(JAVA_TEST_APPS) $(NODEJS_TEST_APPS)) 2>/dev/null || true
