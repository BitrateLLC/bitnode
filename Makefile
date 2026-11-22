VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

.PHONY: build clean install build-linux build-linux-arm64 build-all

# Build for current platform
build:
	go build -ldflags "$(LDFLAGS)" -tags "with_quic with_utls with_wireguard with_clash_api" -o bitnode ./cmd/bitnode
	ln -sfn bitnode bitctl

# Build for Linux amd64
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -tags "with_quic with_utls with_wireguard with_acme with_clash_api" -o bitnode-linux-amd64 ./cmd/bitnode


# Build for Linux arm64
build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -tags "with_quic with_utls with_wireguard with_acme with_clash_api" -o bitnode-linux-arm64 ./cmd/bitnode


# Build all platforms
build-all: build-linux build-linux-arm64


# Clean build artifacts
clean:
	rm -f bitnode bitnode-linux-* bitctl


# Install to system
install: build
	sudo install -m 755 bitnode /usr/local/bin/bitnode
	sudo ln -sfn /usr/local/bin/bitnode /usr/local/bin/bitctl
	sudo mkdir -p /etc/bitnode
	@if [ ! -f /etc/bitnode/config.yml ]; then \
		sudo cp config.yml.example /etc/bitnode/config.yml; \
		echo "Config copied to /etc/bitnode/config.yml - please edit it"; \
	fi
