# Makefile for the shantilly project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

# Tools
GOFUMPT=gofumpt
GOLANGCI_LINT=golangci-lint

# Binaries
BINARY_NAME=shantilly
BINARY_DIR=bin

# Cross-compilation targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Targets
.PHONY: all build test lint fmt fmt-check check clean build-all $(PLATFORMS)

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/shantilly

test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

fmt:
	@echo "Formatting code with gofumpt..."
	@$(GOFUMPT) -w .

fmt-check:
	@echo "Checking code format with gofumpt..."
	@output=$$($(GOFUMPT) -l .); \
		if [ -n "$$output" ]; then \
			echo "The following files need formatting:"; \
			echo "$$output"; \
			exit 1; \
		fi

lint:
	@echo "Running linter..."
	@$(GOLANGCI_LINT) run --config=.golangci.yml ./...

check:
	@echo "Running full local checks (format, lint, tests)..."
	@$(MAKE) fmt-check
	@$(MAKE) lint
	@$(MAKE) test

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BINARY_DIR)

# Cross-compilation targets
build-all: $(PLATFORMS)

$(PLATFORMS):
	@echo "Building for $@..."
	@mkdir -p $(BINARY_DIR)/$@
	@if [ "$(word 1, $(subst /, ,$@))" = "windows" ]; then \
		CGO_ENABLED=0 GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) $(GOBUILD) -ldflags="-s -w" -o $(BINARY_DIR)/$@/$(BINARY_NAME).exe ./cmd/shantilly; \
	else \
		CGO_ENABLED=0 GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) $(GOBUILD) -ldflags="-s -w" -o $(BINARY_DIR)/$@/$(BINARY_NAME) ./cmd/shantilly; \
	fi
