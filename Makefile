.PHONY: test test-unit test-integration test-verbose test-coverage coverage-report test-watch help lint test-routes test-route-registration test-http-methods

# Run all tests
test:
	@echo "Running all tests..."
	@go test -v ./...

# Run only unit tests
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -run Integration ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v -race -count=1 ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Display coverage report
coverage-report: test-coverage
	@echo "Opening coverage report..."
	@if command -v xdg-open > /dev/null; then xdg-open coverage.html; fi

# Run specific test package
test-package:
	@echo "Usage: make test-package PKG=./internal/Bill"
	@if [ -z "$(PKG)" ]; then \
		echo "Please specify PKG parameter"; \
		exit 1; \
	fi
	@go test -v $(PKG)

# Run tests for a specific service
test-service:
	@echo "Usage: make test-service SERVICE=core"
	@if [ -z "$(SERVICE)" ]; then \
		echo "Available services: core, crawler, flowreservation, operator, rsps, client, scaffold-gen, wadl-extract"; \
		exit 1; \
	fi
	@go test -v ./cmd/$(SERVICE)/... ./internal/...

# Route Registration Tests (Phase 3)
test-route-registration:
	@echo "Running route registration tests..."
	@go test -v ./cmd/... -run "RouteRegistration" -parallel 8

# HTTP Method Tests (Phase 4)
test-http-methods:
	@echo "Running HTTP method tests..."
	@go test -v ./cmd/... -run "HTTPMethod" -parallel 8

# All Route Tests (Phases 3-4)
test-routes: test-route-registration test-http-methods
	@echo "✅ All route tests passed!"

# Generate mocks using mockgen
generate-mocks:
	@echo "Generating mocks..."
	@find ./internal -name "*.go" -type f | while read file; do \
		grep -l "^type.*interface {" "$$file" > /dev/null && echo "Found interface in $$file"; \
	done
	@echo "To generate a mock, run: mockgen -source=<file.go> -destination=<mock_file.go> -package=<package>"

# Run tests and show short summary
test-short:
	@go test -short ./... -v

# Run tests with race detector (detects data races)
test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

# Run lint checks (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running linters..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Display help
help:
	@echo "Microservices Testing Commands:"
	@echo ""
	@echo "  make test                    - Run all tests"
	@echo "  make test-unit               - Run only unit tests"
	@echo "  make test-integration        - Run only integration tests"
	@echo "  make test-verbose            - Run all tests with verbose output"
	@echo "  make test-coverage           - Run tests and generate coverage report"
	@echo "  make coverage-report         - Open coverage HTML report"
	@echo "  make test-short              - Run short tests only"
	@echo "  make test-race               - Run tests with race detector"
	@echo "  make test-package PKG=...    - Run tests for specific package"
	@echo "  make test-service SERVICE=...  - Run tests for specific service"
	@echo ""
	@echo "Route Testing Commands (NEW):"
	@echo "  make test-routes             - Run all route registration & HTTP method tests"
	@echo "  make test-route-registration - Run route registration verification tests"
	@echo "  make test-http-methods       - Run HTTP method support tests"
	@echo ""
	@echo "  make generate-mocks          - Generate test mocks (mockgen)"
	@echo "  make lint                    - Run linters"
	@echo "  make help                    - Show this help message"
