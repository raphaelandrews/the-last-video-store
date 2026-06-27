.PHONY: all build build-server build-client build-tamper test test-all test-integration clean seed run-server run-client lint fmt vet cross-compile docker-build tamper-list tamper-demo tamper-restore

BINARY_SERVER=bin/server
BINARY_CLIENT_LINUX=bin/thelastvideostore-linux
BINARY_CLIENT_WINDOWS=bin/thelastvideostore.exe
BINARY_TAMPER_LINUX=bin/tamper-linux
BINARY_TAMPER_WINDOWS=bin/tamper.exe
DB_FILE=thelastvideostore.db

all: build

build: build-server build-client build-tamper

build-server:
	@echo "Building server..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -o $(BINARY_SERVER) ./cmd/server/

build-client:
	@echo "Building client..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -o $(BINARY_CLIENT_LINUX) ./cmd/client/

build-tamper:
	@echo "Building tamper tool..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -o $(BINARY_TAMPER_LINUX) ./cmd/tamper/

cross-compile:
	@echo "Cross-compiling all binaries (Linux + Windows)..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o $(BINARY_SERVER)         ./cmd/server/
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o $(BINARY_CLIENT_LINUX)   ./cmd/client/
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o $(BINARY_TAMPER_LINUX)   ./cmd/tamper/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/server.exe          ./cmd/server/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_CLIENT_WINDOWS) ./cmd/client/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_TAMPER_WINDOWS) ./cmd/tamper/
	@echo "Binaries:"
	@ls -lh bin/

test: test-integration

test-all:
	go test -v -race ./tests/...

test-integration:
	go test -v -count=1 -timeout 120s ./tests/

seed:
	go run ./data/seed.go $(DB_FILE)

# ─── Audit chain demo ──────────────────────────────────────────────────────

tamper-list: build-tamper
	./$(BINARY_TAMPER_LINUX) list

tamper-demo: build-tamper
	./$(BINARY_TAMPER_LINUX) demo

tamper-restore: build-tamper
	./$(BINARY_TAMPER_LINUX) restore $(ID)

# ───────────────────────────────────────────────────────────────────────────

run-server: build-server
	./$(BINARY_SERVER)

run-client: build-client
	./$(BINARY_CLIENT_LINUX)

clean:
	rm -rf bin/
	rm -f $(DB_FILE)

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@which golangci-lint > /dev/null && golangci-lint run ./... || echo "golangci-lint not installed, skipping"

docker-build:
	docker build -f Dockerfile.server -t thelastvideostore:latest .

docker-run:
	docker run -p 8080:8080 -e TLVS_JWT_SECRET=change-me -e TLVS_AES_KEY=0123456789abcdef0123456789abcdef thelastvideostore:latest
