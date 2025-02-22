# Use linker flags to provide version/build settings to the target
CONFIGURATION_PACKAGE := github.com/gphotosuploader/gphotos-uploader-cli/version

# This VERSION could be set calling `make VERSION=0.2.0`
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X ${CONFIGURATION_PACKAGE}.versionString=$(VERSION)"

TEST_VERSION="0.0.0-test.preview"
TEST_LDFLAGS=-ldflags "-X ${CONFIGURATION_PACKAGE}.versionString=$(TEST_VERSION)"

# go source files, ignore vendor directory
PKGS = $(shell go list ./... | grep -v /vendor)
SRC := main.go
BINARY := gphotos-uploader-cli

# Temporary files to be used, you can changed it calling `make TMP_DIR=/tmp`
TMP_DIR ?= .tmp
COVERAGE_FILE := $(TMP_DIR)/coverage.txt
COVERAGE_HTML_FILE := $(TMP_DIR)/coverage.html
GOLANGCI := $(TMP_DIR)/golangci-lint
GOLANGCI_VERSION := 1.64.5

# set how to open files based on OS and ARCH.
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)
ifeq ($(UNAME_ARCH),x86_64)
ifeq ($(UNAME_OS),Darwin)
OPEN_CMD := open
endif
ifeq ($(UNAME_OS),Linux)
OPEN_CMD := xdg-open
endif
endif

.DEFAULT_GOAL := help
.PHONY: test
test: ## Run all the tests
	@echo "--> Running tests..."
	@mkdir -p $(dir $(COVERAGE_FILE))
	@go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) -race -failfast -timeout=30s ${TEST_LDFLAGS} $(PKGS)

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	@echo "--> Creating HTML coverage report at $(COVERAGE_HTML_FILE)..."
	@mkdir -p $(dir $(COVERAGE_FILE)) $(dir $(COVERAGE_HTML_FILE))
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML_FILE)
ifndef COVEROPEN
	@echo "--> Open HTML coverage report: $(OPEN_CMD) $(COVERAGE_HTML_FILE)"
else
	$(OPEN_CMD) $(COVERAGE_HTML_FILE)
endif

build: ## Build the app
	@echo "--> Building binary artifact ($(BINARY) $(VERSION))..."
	@go build ${LDFLAGS} -o $(BINARY) $(SRC)

.PHONY: clean
clean: ## Clean all built artifacts
	@echo "--> Cleaning all built artifacts..."
	@rm -f $(GOLANGCI) $(COVERAGE_FILE) $(COVERAGE_HTML_FILE)
	@rm -rf dist
	@go clean
	@go mod tidy -v

$(GOLANGCI):
	@echo "--> Installing golangci v$(GOLANGCI_VERSION)..."
	@mkdir -p $(dir $(GOLANGCI))
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh  | sh -s -- -b $(dir $(GOLANGCI)) v$(GOLANGCI_VERSION)

.PHONY: lint
lint: $(GOLANGCI) ## Run linter
	@echo "--> Running linter golangci v$(GOLANGCI_VERSION)..."
	@$(GOLANGCI) run

.PHONY: ci
ci: lint test cover build ## Run all the tests and code checks

.PHONY: version
version: ## Show current version
	@echo "$(VERSION)"

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
