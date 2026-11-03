# AGENTS.md - S2I Builder & Runtime Images

Source-to-Image (S2I) **builder** and **runtime** images for **Java and Node.js**, with minimal verification apps in `s2i-tests/`.

## Repository Structure
```text
.
├── s2i/
│   ├── java-builder/      # Maven/Gradle build -> /opt/app-root/src/app.jar
│   ├── java-runtime/      # JRE-only runtime (port 8080)
│   ├── nodejs-builder/    # npm/yarn/pnpm install + build -> /opt/app-root/src (app + prod deps)
│   └── nodejs-runtime/    # Node.js runtime, no package managers (port 3000)
├── s2i-tests/
│   ├── java-maven/        # Spring Boot (Maven, mvnw) test app (/actuator/health)
│   ├── java-gradle/       # Spring Boot (Gradle, gradlew) test app (/actuator/health)
│   ├── nodejs-npm/        # Plain Node.js http server, npm lockfile (/health)
│   └── nodejs-pnpm/       # Plain Node.js http server, pnpm lockfile (/health)
└── Makefile
```

## Commands
Use `podman` if available, fallback to `docker`.

```bash
# Build all images (or per language: make build-java / make build-nodejs)
make build

# Run all end-to-end tests (or per language: make test-java / make test-nodejs)
make test

# Test a single example (builds S2I -> runtime stage, verifies UID 1001, no build tools, and HTTP healthcheck)
make test-java-maven
make test-java-gradle
make test-nodejs-npm
make test-nodejs-pnpm

# Manual S2I builder -> runtime build with Podman (--runtime-image cannot be combined with --as-dockerfile)
s2i build s2i-tests/java-maven localhost/java-builder:dev localhost/java-maven-test:dev --as-dockerfile Dockerfile.gen
printf '\nFROM localhost/java-runtime:dev\nCOPY --from=0 --chown=1001:0 /opt/app-root/src/app.jar /opt/app-root/src/app.jar\nUSER 1001\nCMD ["/usr/libexec/s2i/run"]\n' >> Dockerfile.gen
podman build -t localhost/java-maven-test:dev -f Dockerfile.gen .
rm -rf Dockerfile.gen upload

# Scan runtime for CVEs
trivy image localhost/java-runtime:dev
```

## Rules & Conventions

### Images (`s2i/`)
- **Separation of concerns**: Builders compile/install deps (`s2i/bin/assemble`) from `/tmp/src` and place runnable artifacts in `/opt/app-root/src` (e.g., `/opt/app-root/src/app.jar`); runtimes only start the app (`s2i/bin/run`). Never install build tools (`mvn`, `gradle`, `javac`, `npm`, `gcc`) in runtime images.
- **Security & permissions**: Always run as `USER 1001` (never root) with `HOME=/opt/app-root` and `/opt/app-root/src` owned by `1001:0`. Combine `RUN` statements to keep layers small.
- **Required metadata**: Include labels `io.k8s.description`, `io.openshift.s2i.scripts-url`, and `maintainer`.
- **Scripts**: Place executable S2I scripts in `s2i/<image>/s2i/bin/`.
- **Java**: Prefer `./mvnw` / `./gradlew` wrappers when present in the project, falling back to system `mvn` / `gradle`.
- **Node.js**: Package manager is detected from the lockfile (`package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`); exactly one is required. Use `npm ci`, `yarn install --frozen-lockfile` / `--immutable`, or `pnpm install --frozen-lockfile` (never `npm install`). Yarn/pnpm are provided via corepack (`packageManager` field pins the version). The runtime starts the app with plain `node` (`APP_MAIN` -> `package.json` `main` -> `index.js`/`server.js`).
- **Scope**: Modify and test one S2I image at a time.

### Test Apps (`s2i-tests/`)
- `java-maven` and `java-gradle` are minimal Spring Boot applications exposing `/actuator/health` (`{"status":"UP"}`) on port `8080`.
- `nodejs-npm` and `nodejs-pnpm` are dependency-free Node.js http servers exposing `/health` (`{"status":"UP"}`) on port `3000`. Commit their lockfiles.
- Keep examples minimal, isolated (sparse-checkout friendly), and buildable without external secrets.
- Never commit `node_modules/`, `target/`, `dist/`, `build/`, `.gradle/`, `upload/`, or `Dockerfile.gen`.

## Verification Workflow
1. Build the modified image locally (`make build-java`, `make build-nodejs`, or `make build`).
2. Run the corresponding `Makefile` test target (`make test-java-*`, `make test-nodejs-*`) which verifies:
   - Multi-stage `builder` -> `runtime` build succeeds.
   - Final container runs as non-root (`UID 1001`).
   - Build tools are absent from the runtime image (Java: `mvn`, `gradle`, `javac`; Node.js: `npm`, `npx`, `yarn`, `pnpm`, `corepack`, `gcc`, `g++`, `make`).
   - Application responds with `"status":"UP"` on its health endpoint (Java: `http://localhost:8080/actuator/health`, Node.js: `http://localhost:3000/health`).
