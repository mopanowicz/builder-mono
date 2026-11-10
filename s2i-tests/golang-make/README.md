# golang-make

Minimal Go HTTP/HTTPS server used to verify the Go S2I builder and runtime
images. It is built as a **static binary through its `Makefile`**.

## Endpoints

| Method | Path       | Response                                     |
|--------|------------|----------------------------------------------|
| `GET`  | `/healthz` | `200` `{"status":"UP"}` (`application/json`) |
| `GET`  | `/`        | `200` `Hello from golang-make`               |

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

## Make Targets

Requires Go 1.24 or newer and `make`.

| Target  | Description                                                   |
|---------|---------------------------------------------------------------|
| `all`   | Run `test`, then `build` (default)                            |
| `build` | Build a static, stripped binary (`CGO_ENABLED=0`) to `bin/app` |
| `test`  | Run unit tests (`go test -v ./...`)                           |
| `clean` | Remove the `bin/` directory                                   |

Overridable variables:

| Variable  | Default | Description               |
|-----------|---------|---------------------------|
| `APP`     | `app`   | Output binary name        |
| `BIN_DIR` | `bin`   | Output directory          |
| `GO`      | `go`    | Go toolchain command      |

## Build & Run

```bash
make test build

# Custom output location (e.g. dist/server)
make build BIN_DIR=dist APP=server

# Run HTTP only
./bin/app &
curl http://localhost:8080/healthz

# Run HTTP + HTTPS with a PKCS#12 keystore
KEYSTORE_P12=/path/to/keystore.p12 KEYSTORE_PASSWORD=changeit ./bin/app &
curl http://localhost:8080/healthz
curl -k https://localhost:8443/healthz
```

## S2I Test

`test-config.mk` sets the health-check path used by the repository's
end-to-end test to `/healthz`.
