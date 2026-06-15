.PHONY: all build build-server build-client test test-all test-integration clean seed run-server run-client lint fmt cross-compile docker-build

BINARY_SERVER=bin/server
BINARY_CLIENT_LINUX=bin/thelastvideostore-linux
BINARY_CLIENT_WINDOWS=bin/thelastvideostore.exe
DB_FILE=thelastvideostore.db

all: test build

build: build-server build-client

build-server:
	@echo "Building server..."
	@mkdir -p bin
	go build -o $(BINARY_SERVER) ./cmd/server/

build-client:
	@echo "Building client..."
	@mkdir -p bin
	go build -o $(BINARY_CLIENT_LINUX) ./cmd/client/

cross-compile: build-server
	@echo "Cross-compiling client..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_CLIENT_LINUX) ./cmd/client/
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_CLIENT_WINDOWS) ./cmd/client/
	@echo "Binaries:"
	@ls -lh bin/

test:
	go test -v -race ./internal/...

test-all:
	go test -v -race ./internal/... ./tests/...

test-integration:
	go test -v -count=1 -timeout 120s ./tests/

seed:
	go run ./data/seed.go $(DB_FILE)

run-server: build-server
	./$(BINARY_SERVER)

run-client: build-client
	./$(BINARY_CLIENT_LINUX)

clean:
	rm -rf bin/
	rm -f $(DB_FILE)

fmt:
	go fmt ./...

lint:
	@which golangci-lint > /dev/null && golangci-lint run ./... || echo "golangci-lint not installed, skipping"

docker-build:
	docker build -f Dockerfile.server -t thelastvideostore:latest .

docker-run:
	docker run -p 8080:8080 -e TLVS_JWT_SECRET=change-me -e TLVS_AES_KEY=0123456789abcdef0123456789abcdef thelastvideostore:latest
