BINARY := gphotos-uploader-cli
.DEFAULT_GOAL: release

# This VERSION could be set calling `make VERSION=0.2.0`
VERSION ?= 0.1.2

# This BUILD is automatically calculated and used inside the command
BUILD := $(shell git rev-parse --short HEAD)

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
PKGS = $(shell go list ./... | grep -v /vendor)
SRC := cmd/gphotos-uploader-cli/main.go

PLATFORMS := linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-v$(VERSION)-$(os)-amd64 $(SRC)

.PHONY: release
release: linux darwin

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: test
test:
	@go test -v -race $(PKGS)

.PHONY: clean
clean:
	@rm -rf release

$(GOMETALINTER):
	@echo "--> Installing gometalinter"
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	@gometalinter ./... --vendor
