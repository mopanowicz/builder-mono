# AGENTS.md - S2I Builder & Runtime Images

## Project Overview
This monorepo provides Source-to-Image (S2I) **builder** and **runtime** images for **Java and Node.js**. The `s2i/` directory contains the image definitions. The `s2i-tests/` directory contains minimal apps to verify the builder/runtime flow.

Goal: Build secure, small, reproducible builder and runtime images and test them against real apps.

## Repository Structure
.
├── s2i/
│   ├── java-builder/      # S2I builder image for Java (Maven/Gradle -> artifact)
│   ├── java-runtime/      # S2I runtime image for Java (JRE only to run artifact)
│   ├── nodejs-builder/    # S2I builder image for Node.js (npm/pnpm build)
│   └── nodejs-runtime/    # S2I runtime image for Node.js (node + prod deps)
├── s2i-tests/
│   ├── java-maven/        # Tests java builder/runtime
│   ├── java-gradle/
│   ├── nodejs-npm/
│   └── nodejs-yarn/
└── AGENTS.md
- Builder images are responsible for: `s2i/scripts/assemble` (install deps, build)
- Runtime images are responsible for: `s2i/scripts/run` (start app)
- Builders should produce artifacts that runtimes can consume. Do not mix concerns.

## Commands - Use these exact commands

### Build Images Locally
We use Docker/Podman. Use podman if available, fallback to docker.

```bash
# Build a specific builder/runtime
podman build -t localhost/java-builder:dev -f s2i/java-builder/Dockerfile s2i/java-builder/
podman build -t localhost/java-runtime:dev -f s2i/java-runtime/Dockerfile s2i/java-runtime/
podman build -t localhost/nodejs-builder:dev -f s2i/nodejs-builder/Dockerfile s2i/nodejs-builder/
podman build -t localhost/nodejs-runtime:dev -f s2i/nodejs-runtime/Dockerfile s2i/nodejs-runtime/

# Build all (if Makefile exists)
make build
### Test Images with Examples (The main workflow)
This is how you must verify any change.
# S2I build using local images - Java example
s2i build s2i-tests/java-maven localhost/java-builder:dev localhost/java-test-app --as-dockerfile Dockerfile.gen
podman build -t localhost/java-test-app -f Dockerfile.gen .

# Or if you use builder -> runtime chain:
s2i build s2i-tests/java-maven localhost/java-builder:dev localhost/java-intermediate
s2i build localhost/java-intermediate localhost/java-runtime:dev localhost/java-final --as-dockerfile
# Then run
podman run --rm -p 8080:8080 localhost/java-final

# Node.js example
s2i build s2i-tests/nodejs-npm localhost/nodejs-builder:dev localhost/node-test-app
podman run --rm -p 3000:3000 localhost/node-test-app
### Test All Examples
./hack/test-examples.sh
# or
make test-examples
## Conventions for Agents

### 1. Image Rules
- *NEVER* install build tools (maven, npm, gcc) in runtime images. Runtimes must be minimal + secure.
- *NEVER* run as root in final runtime. Use USER 1001 (S2I default).
- Keep layers small. Combine RUN statements.
- All images must have labels: `io.k8s.description`, `io.openshift.s2i.scripts-url`, `maintainer`.
- S2I scripts must be in `s2i/<image>/s2i/bin/` or `s2i/<image>/scripts/` and be executable.

### 2. When Editing s2i/
- After changing `s2i/_-builder`, you MUST test with corresponding example from `s2i-tests/`.
- After changing `s2i/_-runtime`, you MUST test the full builder -> runtime flow.
- Check CVE scanning: `trivy image localhost/java-runtime:dev`

### 3. When Editing s2i-tests/
- Examples must stay minimal and copy-pasteable. They are not full apps.
- Each example must be buildable without external secrets.
- Don't add monorepo dependencies to examples. They must simulate a real user project.
- Keep examples isolated: to get one example only, users use sparse-checkout:
  `git clone --filter=blob:none --sparse <url> && git sparse-checkout set s2i-tests/java-maven`

### 4. Workflow for a Task
1. Understand which image is affected (builder vs runtime, java vs node)
2. Build the image locally
3. Test with the matching example from `s2i-tests/`
4. Verify `podman run` starts and responds on expected port (8080 for java, 3000 for node)
5. Check logs for errors

## What NOT to do
- Do not modify all 4 s2i images at once. Change one at a time and test.
- Do not add language-specific logic to the root Makefile/Dockerfile - keep it inside s2i/<lang>-*/
- Do not commit built artifacts, `node_modules/`, `target/`, `dist/` inside s2i-tests/
- Do not use `npm install` - use `npm ci` / `pnpm install --frozen-lockfile` in builders

## PR Checklist
- [ ] Image builds
- [ ] Matching example in `s2i-tests/` builds and runs
- [ ] Runtime image is non-root
- [ ] No secrets or large binaries committed

Want me to generate one for a specific subfolder too, like `s2i/java-builder/AGENTS.md` with the exact assemble/run contract?