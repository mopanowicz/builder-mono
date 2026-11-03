# AGENTS.md - S2I Builder & Runtime Images

Source-to-Image (S2I) **builder** and **runtime** images for **Java and Node.js**, with minimal verification apps in `s2i-tests/`.

## Repository Structure
```text
.
├── s2i/
│   ├── java-builder/      # Maven/Gradle build -> artifact
│   ├── java-runtime/      # JRE-only runtime (port 8080)
│   ├── nodejs-builder/    # npm/pnpm build -> artifact
│   └── nodejs-runtime/    # Node.js runtime (port 3000)
├── s2i-tests/
│   ├── java-maven/        # Spring Boot (Maven) test app (/actuator/health)
│   ├── java-gradle/       # Spring Boot (Gradle) test app (/actuator/health)
│   ├── nodejs-npm/
│   └── nodejs-pnpm/
└── Makefile
```

## Commands
Use `podman` if available, fallback to `docker`.

```bash
# Build all images (or build individually: podman build -t localhost/<name>:dev -f s2i/<name>/Dockerfile s2i/<name>/)
make build

# Run all end-to-end tests
make test-examples

# Test a single builder -> runtime flow (e.g., java-maven)
s2i build s2i-tests/java-maven localhost/java-builder:dev localhost/java-app --runtime-image localhost/java-runtime:dev
podman run --rm -p 8080:8080 localhost/java-app

# Scan runtime for CVEs
trivy image localhost/java-runtime:dev
```

## Rules & Conventions

### Images (`s2i/`)
- **Separation of concerns**: Builders compile/install deps (`s2i/bin/assemble`); runtimes only start the app (`s2i/bin/run`). Never install build tools (`maven`, `npm`, `gcc`) in runtime images.
- **Security & size**: Always run as `USER 1001` (never root). Combine `RUN` statements to keep layers small.
- **Required metadata**: Include labels `io.k8s.description`, `io.openshift.s2i.scripts-url`, and `maintainer`.
- **Scripts**: Place executable S2I scripts in `s2i/<image>/s2i/bin/`.
- **Node.js**: Use `npm ci` or `pnpm install --frozen-lockfile` in builders (never `npm install`).
- **Scope**: Modify and test one S2I image at a time. Keep language-specific logic inside `s2i/<lang>-*/`, not the root `Makefile`.

### Test Apps (`s2i-tests/`)
- `java-maven` and `java-gradle` are minimal Spring Boot applications with a healthcheck endpoint at `/actuator/health` (port `8080`).
- Keep examples minimal, isolated (sparse-checkout friendly), and buildable without external secrets.
- Never commit `node_modules/`, `target/`, `dist/`, `build/`, or `Dockerfile.gen`.

## Verification Workflow
1. Build the modified image locally.
2. Build and run the matching example from `s2i-tests/` (`8080/actuator/health` for Java, `3000` for Node.js) and check logs for errors.
3. If modifying a runtime image, test the full builder -> runtime chain and verify non-root execution.
