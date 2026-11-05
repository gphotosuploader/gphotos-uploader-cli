# Use linker flags to provide version/build settings to the target
CONFIGURATION_PACKAGE := github.com/gphotosuploader/gphotos-uploader-cli/version

# This VERSION could be set calling `make VERSION=0.2.0`
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X $(CONFIGURATION_PACKAGE).versionString=$(VERSION)"

TEST_VERSION := 0.0.0-test.preview
TEST_LDFLAGS := -ldflags "-X $(CONFIGURATION_PACKAGE).versionString=$(TEST_VERSION)"
TEST_FLAGS := -race -failfast -timeout=2m $(TEST_LDFLAGS)

PKGS := ./...
SRC := main.go
BINARY := gphotos-uploader-cli

# Temporary files to be used, you can changed it calling `make TMP_DIR=/tmp`
TMP_DIR ?= .tmp
COVERAGE_FILE := $(TMP_DIR)/coverage.txt
COVERAGE_HTML_FILE := $(TMP_DIR)/coverage.html

# Keep in sync with GOLANGCI_LINT_VERSION in .github/workflows/lint.yml.
GOLANGCI_VERSION := 2.15.0
# The version is part of the path, so bumping GOLANGCI_VERSION triggers a reinstall.
GOLANGCI := $(TMP_DIR)/golangci-lint-v$(GOLANGCI_VERSION)/golangci-lint

# How to open the HTML coverage report (`make cover OPEN=1`).
UNAME_OS := $(shell uname -s)
ifeq ($(UNAME_OS),Darwin)
OPEN_CMD := open
else
OPEN_CMD := xdg-open
endif

.DEFAULT_GOAL := help

.PHONY: test
test: ## Run all the tests
	@echo "--> Running tests..."
	@go test $(TEST_FLAGS) $(PKGS)

.PHONY: cover
cover: ## Run all the tests with coverage and create an HTML report (OPEN=1 to open it)
	@echo "--> Running tests with coverage..."
	@mkdir -p $(TMP_DIR)
	@go test $(TEST_FLAGS) -covermode=atomic -coverprofile=$(COVERAGE_FILE) $(PKGS)
	@go tool cover -func=$(COVERAGE_FILE) | tail -n 1
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML_FILE)
	@echo "--> HTML coverage report created at $(COVERAGE_HTML_FILE)"
ifdef OPEN
	@$(OPEN_CMD) $(COVERAGE_HTML_FILE)
endif

.PHONY: build
build: ## Build the app
	@echo "--> Building binary artifact ($(BINARY) $(VERSION))..."
	@go build $(LDFLAGS) -o $(BINARY) $(SRC)

.PHONY: clean
clean: ## Clean all built artifacts
	@echo "--> Cleaning all built artifacts..."
	@rm -rf $(BINARY) $(TMP_DIR) dist
	@go clean

$(GOLANGCI):
	@echo "--> Installing golangci-lint v$(GOLANGCI_VERSION)..."
	@mkdir -p $(dir $(GOLANGCI))
	@# Use the install script from the release tag, not master, so it matches the release assets.
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v$(GOLANGCI_VERSION)/install.sh | sh -s -- -b $(dir $(GOLANGCI)) v$(GOLANGCI_VERSION)

.PHONY: lint
lint: $(GOLANGCI) ## Run linter
	@echo "--> Running linter golangci-lint v$(GOLANGCI_VERSION)..."
	@$(GOLANGCI) run

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	@go mod tidy -v

.PHONY: ci
ci: lint cover build ## Run all the tests and code checks

.PHONY: version
version: ## Show current version
	@echo "$(VERSION)"

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
