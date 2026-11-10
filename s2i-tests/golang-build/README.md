# golang-build

Minimal Go HTTP/HTTPS server used to verify the Go S2I builder and runtime
images. It is built as a **static binary with plain `go build`** (no Makefile).

## Endpoints

| Method | Path       | Response                                     |
|--------|------------|----------------------------------------------|
| `GET`  | `/healthz` | `200` `{"status":"UP"}` (`application/json`) |
| `GET`  | `/`        | `200` `Hello from golang-build`              |

`/healthz` follows the Kubernetes health endpoint convention. When `KEYSTORE_P12`
is configured, both HTTP and HTTPS serve the same routes.

## Configuration

| Variable            | Default | Description                                                  |
|---------------------|---------|--------------------------------------------------------------|
| `PORT`              | `8080`  | HTTP listening port                                          |
| `TLS_PORT`          | `8443`  | HTTPS listening port (used when `KEYSTORE_P12` is set)       |
| `KEYSTORE_P12`      | —       | Path to a PKCS#12 (`.p12` / `.pfx`) keystore file            |
| `KEYSTORE_PASSWORD` | `""`    | Password for the PKCS#12 keystore                            |

The server shuts down gracefully on `SIGINT` / `SIGTERM`.

## Build & Run

Requires Go 1.24 or newer.

```bash
# Run unit tests
go test ./...

# Build a static, stripped binary
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o app .

# Run HTTP only
./app &
curl http://localhost:8080/healthz

# Run HTTP + HTTPS with a PKCS#12 keystore
KEYSTORE_P12=/path/to/keystore.p12 KEYSTORE_PASSWORD=changeit ./app &
curl http://localhost:8080/healthz
curl -k https://localhost:8443/healthz
```

## S2I Test

`test-config.mk` sets the health-check path used by the repository's
end-to-end test to `/healthz`.
