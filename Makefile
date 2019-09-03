BINARY := gphotos-uploader-cli
.DEFAULT_GOAL := help

# This VERSION could be set calling `make VERSION=0.2.0`
VERSION ?= $(shell git describe --tags --abbrev=0)

# This BUILD is automatically calculated and used inside the command
BUILD := $(shell git rev-parse --short HEAD)

# Use linker flags to provide version/build settings to the target
VERSION_IMPORT_PATH := github.com/gphotosuploader/gphotos-uploader-cli/cmd
RELEASE_VERSION_FLAGS=-X=${VERSION_IMPORT_PATH}.version=$(VERSION) -X=${VERSION_IMPORT_PATH}.build=$(BUILD)
LDFLAGS=-ldflags "$(RELEASE_VERSION_FLAGS)"

# go source files, ignore vendor directory
PKGS = $(shell go list ./... | grep -v /vendor)
SRC := main.go
COVERAGE_FILE := coverage.txt

# Get first path on multiple GOPATH environments
GOPATH := $(shell echo ${GOPATH} | cut -d: -f1)

.PHONY: test
test: ## Run all the tests
	@echo "--> Running tests..."
	@go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) -race -failfast -timeout=30s $(PKGS)

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	@echo "--> Openning coverage report..."
	@go tool cover -html=$(COVERAGE_FILE)

.PHONY: coveralls
coveralls: test ## Run all the tests and send it to Coveralls (only CI)
	@echo "--> Sending coverage report to Coveralls..."
	@go get github.com/mattn/goveralls
	@goveralls -coverprofile $(COVERAGE_FILE) -service drone.io

build: ## Build the app
	@echo "--> Building binary artifact ($(BINARY) $(VERSION) (build: $(BUILD)))..."
	@go build ${LDFLAGS} -o $(BINARY) $(SRC)

.PHONY: clean
clean: ## Clean all built artifacts
	@echo "--> Cleaning all built artifacts..."
	@rm -f $(BINARY) $(COVERAGE_FILE)
	@rm -rf dist

BIN_DIR := $(GOPATH)/bin

GOLANGCI := $(BIN_DIR)/golangci-lint
GOLANGCI_VERSION := 1.12.3

GORELEASER := $(BIN_DIR)/goreleaser

$(GOLANGCI):
	@echo "--> Installing golangci v$(GOLANGCI_VERSION)..."
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_VERSION)

$(GORELEASER):
	@echo "--> Installing goreleaser..."
	@curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh -s -- -b $(BIN_DIR)

.PHONY: lint
lint: $(GOLANGCI) ## Run linter
	@echo "--> Running linter golangci v$(GOLANGCI_VERSION)..."
	@$(GOLANGCI) run

.PHONY: ci
ci: build test lint ## Run all the tests and code checks

.PHONY: release
release: $(GORELEASER) ## Release a new version using goreleaser (only CI)
	@echo "--> Releasing $(BINARY) $(VERSION) (build: $(BUILD))..."
	@RELEASE_VERSION_TAG="$(RELEASE_VERSION_FLAGS)" $(GORELEASER) release

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version: ## Show current version
	@echo "$(VERSION) (build: $(BUILD))"
