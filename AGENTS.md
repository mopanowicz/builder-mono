# AGENTS.md - S2I Builder & Runtime Images

Source-to-Image (S2I) **builder** and **runtime** images for **Java and Node.js**, with minimal verification apps in `s2i-tests/`.

This repository is intended as a reusable image-building foundation: builders install dependencies and compile or package the app, while runtime images contain only the runnable artifact and the minimum runtime environment needed to execute it.

## Repository Layout
```text
.
├── s2i/
│   ├── java-builder/         # Java build image; compiles Maven/Gradle apps
│   ├── java-runtime/         # Java runtime image; JRE only
│   ├── nodejs-builder/       # Node.js build image; installs/builds app
│   └── nodejs-runtime/       # Node.js runtime; no package managers
├── s2i-tests/
│   ├── spring-boot-{maven,gradle}/
│   ├── quarkus-{maven,gradle}/
│   └── nodejs-{npm,yarn,pnpm}/
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
# Build all images (or per language: make build-java / make build-nodejs)
make build

# Run all end-to-end tests (or per language: make test-java / make test-nodejs)
make test

# Test a single example (runs builder assemble, copies the artifact to runtime, and verifies UID, tools, and health)
make test-spring-boot-maven
make test-spring-boot-gradle
make test-quarkus-maven
make test-quarkus-gradle
make test-nodejs-npm
make test-nodejs-yarn
make test-nodejs-pnpm

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
- **Separation of concerns**: Builders compile/install dependencies (`s2i/bin/assemble`) from `/tmp/src` and place runnable artifacts in `/opt/app-root/src` (for example, `/opt/app-root/src/app.jar`); runtimes only start the app (`s2i/bin/run`). Never install build tools (`mvn`, `gradle`, `javac`, `npm`, `gcc`, `g++`, `make`) in runtime images.
- **Security & permissions**: Always run as `USER 1001` (never root) with `HOME=/opt/app-root` and `/opt/app-root/src` owned by `1001:0`. Combine `RUN` statements to keep layers small.
- **Required metadata**: Include labels `io.k8s.description`, `io.openshift.s2i.scripts-url`, and `maintainer`.
- **Scripts**: Place executable S2I scripts in `s2i/<image>/s2i/bin/`.
- **Java**: Prefer `./mvnw` / `./gradlew` wrappers when present in the project, falling back to system `mvn` / `gradle`.
- **Node.js**: Package manager is detected from the lockfile (`package-lock.json`, `yarn.lock`, or `pnpm-lock.yaml`); exactly one is required. Use `npm ci`, `yarn install --frozen-lockfile` / `--immutable`, or `pnpm install --frozen-lockfile` (never `npm install`). Yarn is supported in both classic and Berry workflows; Yarn/pnpm are provided via corepack (`packageManager` field pins the version). The runtime starts the app with plain `node` (`APP_MAIN` -> `package.json` `main` -> `index.js`/`server.js`).
- **Scope**: Modify and test one S2I image at a time.

### Test Apps (`s2i-tests/`)
- `spring-boot-maven` and `spring-boot-gradle` are minimal Spring Boot applications exposing `/actuator/health` (`{"status":"UP"}`) on port `8080`.
- `quarkus-maven` and `quarkus-gradle` are minimal Quarkus applications exposing `/q/health` (status `UP`) on port `8080`.
- `nodejs-npm`, `nodejs-yarn`, and `nodejs-pnpm` are dependency-free Node.js http servers exposing `/health` (`{"status":"UP"}`) on port `3000`. Commit their lockfiles.
- Keep examples minimal, isolated (sparse-checkout friendly), and buildable without external secrets.
- Never commit `node_modules/`, `target/`, `dist/`, `build/`, `.gradle/`, `upload/`, or `Dockerfile.gen`.

## Verification Workflow
1. Build the modified image locally (`make build-java`, `make build-nodejs`, or `make build`).
2. Run the corresponding `Makefile` test target (`make test-java-*`, `make test-nodejs-*`) which verifies:
   - Java validation covers Maven and Gradle examples.
   - Node.js validation covers npm, Yarn, and pnpm examples.
   - The builder's S2I `assemble` script runs, then the app artifact is copied into the runtime image.
   - Final container runs as non-root (`UID 1001`).
   - Build tools are absent from the runtime image (Java: `mvn`, `gradle`, `javac`; Node.js: `npm`, `npx`, `yarn`, `pnpm`, `corepack`, `gcc`, `g++`, `make`).
   - Application responds with status `UP` on its configured health endpoint (Spring Boot: `http://localhost:8080/actuator/health`, Quarkus: `http://localhost:8080/q/health`, Node.js: `http://localhost:3000/health`).

## Contributor Notes
- Prefer small, focused changes over broad rewrites.
- Keep examples minimal and self-contained.
- Validate behavior using the repo's `Makefile` targets before finalizing a change.
- Prefer reproducible, lockfile-driven dependency resolution over ad hoc package manager commands.
- Maintain a clear separation between build artifacts and runtime artifacts in all custom image work.
