# AGENTS.md

## Build / Lint / Test Commands

From the `Makefile` at the project root:

- **Build everything** (server + client + tamper tool): `make build`
- **Build server only**: `make build-server`
- **Build client only**: `make build-client`
- **Build tamper tool only**: `make build-tamper`
- **Cross-compile (Linux + Windows binaries)**: `make cross-compile`
- **Run all tests with race detector**: `make test-all`
- **Run integration tests only**: `make test-integration`
- **Lint** (uses `golangci-lint` if installed, otherwise no-op): `make lint`
- **Vet**: `make vet`
- **Format**: `make fmt`
- **Seed the database**: `make seed`
- **Run server / client locally**: `make run-server` / `make run-client`
- **Clean build artifacts**: `make clean`

Direct Go equivalents:
- Lint: `golangci-lint run ./...`
- Vet: `go vet ./...`
- Test: `go test -v -race ./...`
