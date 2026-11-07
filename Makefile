.DEFAULT_GOAL := help

.PHONY: all
all: test-unit vet build

.PHONY: help
help: ## Display available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-18s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build all packages
	go build -v ./...

.PHONY: run
run: ## Run the service from the module root
	go run .

.PHONY: test
test: test-unit ## Run the test suite with race detection

.PHONY: test-unit
test-unit: ## Run unit tests with race detection
	go test -race -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests and generate an HTML coverage report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint or fall back to go vet
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		printf '%s\n' 'golangci-lint not found; running go vet instead'; \
		go vet ./...; \
	fi

.PHONY: tidy
tidy: ## Tidy and verify Go module dependencies
	go mod tidy
	go mod verify

.PHONY: clean
clean: ## Remove generated coverage files
	rm -f coverage.out coverage.html profile.out