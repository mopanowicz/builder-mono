# builder-mono

A reusable Source-to-Image (S2I) image foundation for Java, Node.js, and Go applications.

This repository provides a small, opinionated set of builder and runtime images built on Rocky Linux 9 minimal, designed to:
- separate build-time dependencies from runtime dependencies,
- produce minimal runtime containers,
- enforce non-root execution,
- and validate application health with automated checks.

The project is intentionally minimal and reusable: it targets common patterns for Java (Maven/Gradle), Node.js (npm/yarn/pnpm), and Go (`go build`/Makefile), while keeping the final runtime images free from build tooling.

## Why this project exists

Modern containerized application workflows often need to:
- build dependencies in a dedicated builder stage,
- copy only the compiled or packaged result into a lean runtime image,
- run as a non-root user for better security,
- and verify the app can start cleanly and serve health checks.

This repository provides a foundation for exactly that workflow, with reusable S2I scripts and validation apps that exercise the pattern end-to-end.

## Design principles

This project follows the conventions described in [AGENTS.md](./AGENTS.md):

- Builder images are responsible for installation and build steps.
- Runtime images are responsible only for execution.
- All images use the `rockylinux:9-minimal` base. Java images use Rocky OpenJDK 25; the builder installs Maven (`3.10.0`) under `/opt/maven` and Gradle (`9.8.2`) under `/opt/gradle`. Node.js images install Node.js (`24.22.2`) from the official Linux x64 binary distribution under `/usr/local` (with `corepack` installed and enabled in the builder to provide `yarn` and `pnpm`), while Go builder images use Rocky's `golang` package.
- Build tools such as Maven, Gradle, npm, yarn, pnpm, the Go toolchain, gcc, g++, and make are not shipped in runtime images.
- Applications run as non-root user `1001`.
- Java projects prefer local wrappers (`./mvnw`, `./gradlew`) when present, falling back to the builder's installed Maven/Gradle (in `s2i-tests/`, wrapper scripts are disabled as `mvnw.off` / `gradlew.off` so tests verify the system-installed tools).
- Node.js projects require exactly one package manager lockfile (`package-lock.json` / `npm-shrinkwrap.json`, `yarn.lock`, or `pnpm-lock.yaml`).
- Go projects are compiled into a static binary, using `make build` when a `Makefile` is present and `go build` otherwise. The Go test applications support HTTP on port `8080` and optional HTTPS on port `8443` via PKCS#12 keystores (`KEYSTORE_P12` and `KEYSTORE_PASSWORD`).
- The runtime image starts an app without build tooling and with a lean production environment.

## Repository layout

```
.
├── s2i/
│   ├── java-builder/         # Java build image; compiles Maven/Gradle apps
│   ├── java-runtime/         # Java runtime image; JRE only
│   ├── nodejs-builder/       # Node.js build image; installs/builds app
│   ├── nodejs-runtime/       # Node.js runtime image; no package managers
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

Generated local build output such as `target/`, `build/`, `.gradle/`, and `node_modules/` is intentionally excluded from version control and should not be committed.

## Command cheat sheet

GNU Make does not include a built-in option to list targets. To print target names from this Makefile, run:

```bash
make -qp | awk -F: '/^[a-zA-Z0-9][^$#[:space:]]*:/ { print $1 }' | sort -u
```

Common targets:

| Target | Description |
| --- | --- |
| `make build` | Build all builder and runtime images. |
| `make build-java` / `make build-nodejs` / `make build-golang` | Build images for one language. |
| `make test` | Run all S2I end-to-end tests. |
| `make test-java` / `make test-nodejs` / `make test-golang` | Test all examples for one language. |
| `make test-<app>` | Test one example, such as `make test-nodejs-pnpm`. |
| `make clean` | Remove generated test files and test images. |
