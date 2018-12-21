BINARY := gphotos-uploader-cli
.DEFAULT_GOAL := help

# This VERSION could be set calling `make VERSION=0.2.0`
VERSION ?= 0.1.2

# This BUILD is automatically calculated and used inside the command
BUILD := $(shell git rev-parse --short HEAD)

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
PKGS = $(shell go list ./... | grep -v /vendor)
SRC := cmd/gphotos-uploader-cli/main.go

# Get first path on multiple GOPATH environments
GOPATH := $(shell echo ${GOPATH} | cut -d: -f1)

PLATFORMS := linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):			## Create binary for an specific platform
	@echo "--> Generating binary for $(os) v$(VERSION) (build: $(BUILD))..."
	@mkdir -p release
	@GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-v$(VERSION)-$(os)-amd64 $(SRC)

.PHONY: release
release: linux darwin	## Create binaries for all supported platforms

.PHONY: test
test: lint					## Run tests
	@echo "--> Running tests..."
	@go test -v -race $(PKGS)

.PHONY: clean
clean:					## Clean all built artifacts
	@echo "--> Cleaning all built artifacts..."
	@rm -rf release

BIN_DIR := $(GOPATH)/bin

GOLANGCI := $(BIN_DIR)/golangci-lint
GOLANGCI_VERSION := 1.12.3

$(GOLANGCI):
	@echo "--> Installing golangci v$(GOLANGCI_VERSION)..."
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(BIN_DIR) v$(GOLANGCI_VERSION)

.PHONY: lint
lint: $(GOLANGCI)	## Run linter
	@echo "--> Running linter golangci v$(GOLANGCI_VERSION)..."
	@$(GOLANGCI) run

.PHONY: help
help:					## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
