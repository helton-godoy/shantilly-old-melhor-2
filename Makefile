# Makefile for the shantilly project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLINT=./lint.sh

# Binaries
BINARY_NAME=shantilly
BINARY_DIR=bin

# Cross-compilation targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Targets
.PHONY: all build test vet lint ci clean build-all $(PLATFORMS)

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/shantilly

test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

lint:
	@echo "Running linter..."
	@$(GOLINT)

ci:
	@echo "Running full CI suite (test + vet + lint)..."
	@$(MAKE) test
	@$(MAKE) vet
	@$(MAKE) lint

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
