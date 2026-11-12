# AGENTS.md - S2I Builder & Runtime Images

Source-to-Image (S2I) **builder** and **runtime** images for **Java, Node.js, and Go**, with minimal verification apps in `s2i-tests/`.

This repository is intended as a reusable image-building foundation: builders install dependencies and compile or package the app, while runtime images contain only the runnable artifact and the minimum runtime environment needed to execute it.

## Repository Layout
```text
.
├── s2i/
│   ├── base-builder/         # Common builder base image (OS updates, common packages, S2I dirs)
│   ├── java-builder/         # Java build image; compiles Maven/Gradle apps
│   ├── java-runtime/         # Java runtime image; JRE only
│   ├── nodejs-builder/       # Node.js build image; installs/builds app
│   ├── nodejs-runtime/       # Node.js runtime; no package managers
│   ├── golang-builder/       # Go build image; compiles a static binary (go build / make)
│   └── golang-runtime/       # Go runtime image; binary only, no Go toolchain
├── s2i-tests/
│   ├── spring-boot-{maven,gradle}/
│   ├── quarkus-{maven,gradle}/
│   ├── nodejs-{npm,yarn,pnpm}/
│   └── golang-{build,make}/
├── Makefile
├── AGENTS.md
├── README.md
└── .gitignore
```

Each S2I test app can configure its HTTP health-check path in
`s2i-tests/<app>/test-config.mk` by defining `S2I_TEST_HEALTH_PATH_<app>`.
If no value is configured, the test uses `/health`.

## Commands
Use `podman` if available, fallback to `docker`.

```bash
# Build all images (or base/per language: make build-base / make build-java / make build-nodejs / make build-golang)
make build

# Run all end-to-end tests (or per language: make test-java / make test-nodejs / make test-golang)
make test

# Test a single example (runs builder assemble, copies the artifact to runtime, and verifies UID, tools, and health)
make test-spring-boot-maven
make test-spring-boot-gradle
make test-quarkus-maven
make test-quarkus-gradle
make test-nodejs-npm
make test-nodejs-yarn
make test-nodejs-pnpm
make test-golang-build
make test-golang-make

# Clean generated test files and test images
make clean

# Scan runtime for CVEs
trivy image localhost/java-runtime:dev
```

## Rules & Conventions

### Project goals
- Keep the builder and runtime responsibilities strictly separated.
- Optimize for minimal, secure runtime images.
- Validate both build correctness and runtime behavior with repeatable automation.
- Produce examples that are small, isolated, and easy to reason about.

### Images (`s2i/`)
- **Base image**: `base-builder` and all runtime images use `rockylinux:9-minimal`. Language builder images (`java-builder`, `nodejs-builder`, `golang-builder`) use `localhost/base-builder:dev` as their base image, which provides OS updates, common utilities (`ca-certificates`, `curl`, `git`, `gzip`, `tar`, `unzip`, `xz`), `HOME=/opt/app-root`, `/opt/app-root/src`, `/tmp/src`, `/usr/libexec/s2i` (owned by `1001:0` with `g=u` permissions), `WORKDIR /tmp/src`, `CMD ["/usr/libexec/s2i/usage"]`, and common labels (`io.openshift.s2i.scripts-url`, `maintainer`).
- **Separation of concerns**: Builders compile/install dependencies (`s2i/bin/assemble`) from `/tmp/src` and place runnable artifacts in `/opt/app-root/src` (for example, `/opt/app-root/src/app.jar`); runtimes only start the app (`s2i/bin/run`). Never install build tools (`mvn`, `gradle`, `javac`, `npm`, `go`, `gcc`, `g++`, `make`) in runtime images.
- **Security & permissions**: Always run as `USER 1001` (never root) in language builder and runtime images, with `HOME=/opt/app-root`, `/opt/app-root/src`, `/tmp/src` (builders), and `/usr/libexec/s2i` owned by `1001:0` (`g=u`). Language builders inherit directory setup from `base-builder`, copy their scripts via `COPY --chown=1001:0 ./s2i/bin/ /usr/libexec/s2i`, and run `chmod +x /usr/libexec/s2i/*`. Combine `RUN` statements to keep layers small.
- **Required metadata**: Include labels `io.k8s.description`, `io.openshift.s2i.scripts-url`, and `maintainer` (language builders inherit `io.openshift.s2i.scripts-url` and `maintainer` from `base-builder` and override `io.k8s.description`).
- **Scripts**: Place executable S2I scripts in `s2i/<layer>/s2i/bin/` for language builder and runtime images (`base-builder` only creates `/usr/libexec/s2i`).
- **Java**: Use only Rocky's OpenJDK 25 packages (`java-25-openjdk-devel` in the builder and `java-25-openjdk-headless` in the runtime). Install Apache Maven (`3.10.0`) from its official binary distribution under `/opt/maven` and Gradle (`9.8.2`) under `/opt/gradle` in the builder. Prefer `./mvnw` / `./gradlew` wrappers when present, falling back to those installed tools.
- **Node.js**: Install Node.js (`24.22.2`) from the official Linux x64 binary distribution under `/usr/local` in both builder and runtime images (installing and enabling `corepack` in the builder for `yarn` and `pnpm`, and stripping package managers from the runtime). Package manager is detected from the lockfile (`package-lock.json` / `npm-shrinkwrap.json`, `yarn.lock`, or `pnpm-lock.yaml`); exactly one is required. Use `npm ci`, `yarn install --frozen-lockfile` / `--immutable`, or `pnpm install --frozen-lockfile` (never `npm install`). Yarn is supported in both classic and Berry workflows; Yarn/pnpm are provided via corepack (`packageManager` field pins the version). The runtime starts the app with plain `node` (`APP_MAIN` -> `package.json` `main` -> `server.js` / `index.js` / `app.js` / `dist/server.js` / `dist/index.js`).
- **Go**: `go.mod` is required; modules are downloaded and verified (`go mod verify`) when `go.sum` exists. If a `Makefile` is present the builder runs `make build` (`MAKE_TARGET`) and copies `bin/app` (`MAKE_OUTPUT`); otherwise it runs `go build -trimpath -ldflags="-s -w"` on `GO_MAIN_PACKAGE` (default `.`). Builds are static by default (`CGO_ENABLED=0`) and use `GOTOOLCHAIN=local` (no toolchain downloads). `RUN_TESTS=true` runs `make test` / `go test ./...` before building. The artifact is `/opt/app-root/src/app`; the runtime `exec`s it directly (`APP_BINARY` overrides the path).
- **Scope**: Modify and test one S2I image at a time.

### Test Apps (`s2i-tests/`)
- `spring-boot-maven` and `spring-boot-gradle` are minimal Spring Boot applications exposing `/actuator/health` (`{"status":"UP"}`) on port `8080`.
- `quarkus-maven` and `quarkus-gradle` are minimal Quarkus applications exposing `/q/health` (status `UP`) on port `8080`. Local wrapper scripts in the Java test apps are disabled (`mvnw.off`, `gradlew.off`) so tests exercise the builder's installed Maven and Gradle distributions.
- `nodejs-npm`, `nodejs-yarn`, and `nodejs-pnpm` are dependency-free Node.js http servers exposing `/health` (`{"status":"UP"}`) on port `3000`. Commit their lockfiles.
- `golang-build` and `golang-make` are minimal Go `net/http` servers exposing `/healthz` (`{"status":"UP"}`, Kubernetes naming convention) over HTTP on port `8080` (`PORT`) and optionally HTTPS on port `8443` (`TLS_PORT`) when `KEYSTORE_P12` and `KEYSTORE_PASSWORD` (PKCS#12 keystore) are set. `golang-build` is compiled with plain `go build`; `golang-make` provides a `Makefile` with `build` and `test` targets.
- Keep examples minimal, isolated (sparse-checkout friendly), and buildable without external secrets.
- Never commit `node_modules/`, `target/`, `dist/`, `build/`, `.gradle/`, `upload/`, `Dockerfile.gen`, or compiled Go binaries (`bin/`, `app`).

## Verification Workflow
1. Build the modified image locally (`make build-base`, `make build-java`, `make build-nodejs`, `make build-golang`, or `make build`).
2. Run the corresponding `Makefile` test target (`make test-java`, `make test-nodejs`, `make test-golang`, or `make test-<app>`) which verifies:
   - Java validation covers Maven and Gradle examples.
   - Node.js validation covers npm, Yarn, and pnpm examples.
   - Go validation covers `go build` and `Makefile` examples.
   - The builder's S2I `assemble` script runs, then the app artifact is copied into the runtime image.
   - Final container runs as non-root (`UID 1001`).
   - Build tools are absent from the runtime image (Java: `mvn`, `gradle`, `javac`; Node.js: `npm`, `npx`, `yarn`, `pnpm`, `corepack`, `gcc`, `g++`, `make`; Go: `go`, `gofmt`, `gcc`, `g++`, `cc`, `make`, `git`).
   - Application responds with status `UP` on its configured health endpoint (Spring Boot: `http://localhost:8080/actuator/health`, Quarkus: `http://localhost:8080/q/health`, Node.js: `http://localhost:3000/health`, Go: `http://localhost:8080/healthz`).
   - The test fails if the container cannot start (for example, when the host port is already in use).

## Contributor Notes
- Prefer small, focused changes over broad rewrites.
- Keep examples minimal and self-contained.
- Validate behavior using the repo's `Makefile` targets before finalizing a change.
- Prefer reproducible, lockfile-driven dependency resolution over ad hoc package manager commands.
- Maintain a clear separation between build artifacts and runtime artifacts in all custom image work.
