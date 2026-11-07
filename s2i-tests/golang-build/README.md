# golang-build

Minimal, dependency-free Go HTTP server used to verify the Go S2I builder and
runtime images. It is built as a **static binary with plain `go build`** (no
Makefile).

## Endpoints

| Method | Path       | Response                                 |
|--------|------------|------------------------------------------|
| `GET`  | `/healthz` | `200` `{"status":"UP"}` (`application/json`) |
| `GET`  | `/`        | `200` `Hello from golang-build`          |

`/healthz` follows the Kubernetes health endpoint convention.

## Configuration

| Variable | Default | Description         |
|----------|---------|---------------------|
| `PORT`   | `8080`  | HTTP listening port |

The server shuts down gracefully on `SIGINT` / `SIGTERM`.

## Build & Run

Requires Go 1.24 or newer.

```bash
# Run unit tests
go test ./...

# Build a static, stripped binary
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o app .

# Run and check health
./app &
curl http://localhost:8080/healthz
```

## S2I Test

`test-config.mk` sets the health-check path used by the repository's
end-to-end test to `/healthz`.
