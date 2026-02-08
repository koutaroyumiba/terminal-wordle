# Project name (binary name)
APP_NAME := wordle

# Go settings
GO := go
CMD_DIR := ./cmd/$(APP_NAME)
BIN_DIR := ./dist
BIN := $(BIN_DIR)/$(APP_NAME)

# Default target
all: build

## Build the binary
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(CMD_DIR)

## Run the app (no build artifact)
run:
	$(GO) run $(CMD_DIR)

## Run tests
test:
	$(GO) test ./...

## Format code
fmt:
	$(GO) fmt ./...

## Vet code
vet:
	$(GO) vet ./...

## Clean build artifacts
clean:
	rm -rf $(BIN_DIR)

## Install binary into GOPATH/bin
install:
	$(GO) install $(CMD_DIR)

## Help
help:
	@echo "Available targets:"
	@echo "  make build    - Build binary into ./bin/"
	@echo "  make run      - Run the app"
	@echo "  make test     - Run all tests"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make install  - Install binary to GOPATH/bin"

.PHONY: all build run test fmt vet clean install help
