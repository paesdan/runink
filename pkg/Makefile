# Project Metadata
BINARY_NAME = runi
PKG_PATH = ./...

# Default Target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) main.go

# Install locally
.PHONY: install
install:
	go install ./...

# Run the project
.PHONY: run
run:
	go run main.go

# Lint the code
.PHONY: lint
lint:
	go fmt $(PKG_PATH)
	go vet $(PKG_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/

# Format code
.PHONY: fmt
fmt:
	go fmt $(PKG_PATH)

# Test (if you add *_test.go files later)
.PHONY: test
test:
	go test $(PKG_PATH)

# Update dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Show help
.PHONY: help
help:
	@echo "Runink CLI Developer Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build       - Build the binary"
	@echo "  make run         - Run the CLI locally"
	@echo "  make install     - Install the CLI to your Go bin"
	@echo "  make lint        - Format and vet the code"
	@echo "  make clean       - Remove built binaries"
	@echo "  make test        - Run unit tests"
	@echo "  make tidy        - Clean up go.mod and go.sum"
