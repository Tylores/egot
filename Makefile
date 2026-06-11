.PHONY: test test-unit test-integration test-verbose test-coverage coverage-report test-watch help lint test-routes test-route-registration test-http-methods
.PHONY: build-all build-services build-tools
.PHONY: build-BRS build-Bill build-DCAP build-DR build-EDevice build-File build-MUP build-Messaging build-Notify build-PPY build-SDevice build-TariffProfile build-TimeOfUse build-UPT
.PHONY: build-crawler build-client build-scaffold-gen build-wadl-extract
.PHONY: start stop status restart nginx-config
.PHONY: ssl-refresh ssl-clients
.PHONY: cleanup-services cleanup-services-dry-run cleanup-services-list

BIN_DIR := ./bin

# ─── Build Targets ────────────────────────────────────────────────────────────

build-FlowReservation:
	@echo "Building FlowReservation..."
	@go build -o $(BIN_DIR)/FlowReservation ./cmd/FlowReservation/

build-DER:
	@echo "Building DER..."
	@go build -o $(BIN_DIR)/DER ./cmd/DER/

build-MUP:
	@echo "Building MUP..."
	@go build -o $(BIN_DIR)/MUP ./cmd/MUP/

build-UPT:
	@echo "Building UPT..."
	@go build -o $(BIN_DIR)/UPT ./cmd/UPT/

build-emulator-der:
	@echo "Building emulator-der..."
	@go build -o $(BIN_DIR)/emulator-der ./cmd/emulator-der/

build-operator:
	@echo "Building operator..."
	@go build -o $(BIN_DIR)/operator ./cmd/operator/

build-rsps:
	@echo "Building rsps..."
	@go build -o $(BIN_DIR)/rsps ./cmd/rsps/

build-BRS:
	@echo "Building BRS..."
	@go build -o $(BIN_DIR)/BRS ./cmd/BRS/

build-Bill:
	@echo "Building Bill..."
	@go build -o $(BIN_DIR)/Bill ./cmd/Bill/

build-DCAP:
	@echo "Building DCAP..."
	@go build -o $(BIN_DIR)/DCAP ./cmd/DCAP/

build-DR:
	@echo "Building DR..."
	@go build -o $(BIN_DIR)/DR ./cmd/DR/

build-EDevice:
	@echo "Building EDevice..."
	@go build -o $(BIN_DIR)/EDevice ./cmd/EDevice/

build-File:
	@echo "Building File..."
	@go build -o $(BIN_DIR)/File ./cmd/File/

build-MUP:
	@echo "Building MUP..."
	@go build -o $(BIN_DIR)/MUP ./cmd/MUP/

build-Messaging:
	@echo "Building Messaging..."
	@go build -o $(BIN_DIR)/Messaging ./cmd/Messaging/

build-Notify:
	@echo "Building Notify..."
	@go build -o $(BIN_DIR)/Notify ./cmd/Notify/

build-PPY:
	@echo "Building PPY..."
	@go build -o $(BIN_DIR)/PPY ./cmd/PPY/

build-SDevice:
	@echo "Building SDevice..."
	@go build -o $(BIN_DIR)/SDevice ./cmd/SDevice/

build-TariffProfile:
	@echo "Building TariffProfile..."
	@go build -o $(BIN_DIR)/TariffProfile ./cmd/TariffProfile/

build-TimeOfUse:
	@echo "Building TimeOfUse..."
	@go build -o $(BIN_DIR)/TimeOfUse ./cmd/TimeOfUse/

build-UPT:
	@echo "Building UPT..."
	@go build -o $(BIN_DIR)/UPT ./cmd/UPT/

build-crawler:
	@echo "Building crawler..."
	@go build -o $(BIN_DIR)/crawler ./cmd/crawler/

build-client:
	@echo "Building client..."
	@go build -o $(BIN_DIR)/client ./cmd/client/

build-scaffold-gen:
	@echo "Building scaffold-gen..."
	@go build -o $(BIN_DIR)/scaffold-gen ./cmd/scaffold-gen/

build-wadl-extract:
	@echo "Building wadl-extract..."
	@go build -o $(BIN_DIR)/wadl-extract ./cmd/wadl-extract/

build-nginx-config-gen:
	@echo "Building nginx-config-gen..."
	@go build -o $(BIN_DIR)/nginx-config-gen ./cmd/nginx-config-gen/

# Build all microservice servers (excludes tools/clients)
build-services: build-FlowReservation build-operator build-rsps \
	build-BRS build-Bill build-DCAP build-DR build-EDevice \
	build-File build-MUP build-Messaging build-Notify build-PPY \
	build-SDevice build-TariffProfile build-TimeOfUse build-UPT build-DER
	@echo "✅ All services built in $(BIN_DIR)/"

# Build tools and clients only
build-tools: build-client build-scaffold-gen build-wadl-extract build-nginx-config-gen build-emulator-der
	@echo "✅ All tools built in $(BIN_DIR)/"

# Cleanup stale microservices
cleanup-services-dry-run:
	@echo "Scanning for stale microservices (dry-run)..."
	@go run ./tools/cleanup-services --dry-run=true --interactive=false

cleanup-services:
	@go run ./tools/cleanup-services --dry-run=false --interactive=true

cleanup-services-list:
	@go run ./tools/cleanup-services --list

# Build everything
build-all: build-services build-tools
	@echo "✅ All binaries built in $(BIN_DIR)/"

# ─── Service Management ───────────────────────────────────────────────────────

PIDS_DIR := $(BIN_DIR)/pids
LOGS_DIR := $(BIN_DIR)/logs

SERVICES := FlowReservation operator rsps \
	BRS Bill DCAP DR EDevice File MUP Messaging Notify PPY \
	SDevice TariffProfile TimeOfUse UPT DER

# Start all microservices in the background
start: build-services nginx-config
	@echo "Starting all services..."
	@mkdir -p $(PIDS_DIR) $(LOGS_DIR)
	@for svc in $(SERVICES); do \
		if [ -f $(PIDS_DIR)/$$svc.pid ] && kill -0 $$(cat $(PIDS_DIR)/$$svc.pid) 2>/dev/null; then \
			echo "  $$svc already running (PID $$(cat $(PIDS_DIR)/$$svc.pid))"; \
		else \
			nohup $(BIN_DIR)/$$svc > $(LOGS_DIR)/$$svc.log 2>&1 & echo $$! > $(PIDS_DIR)/$$svc.pid; \
			echo "  ▶ $$svc started (PID $$(cat $(PIDS_DIR)/$$svc.pid))"; \
		fi; \
	done
	@echo "Logs: $(LOGS_DIR)/"
	@echo "Starting nginx API gateway..."
	@bash ./scripts/start-nginx.sh ./nginx/nginx.yaml

# Stop all microservices
stop:
	@echo "Stopping nginx API gateway..."
	@bash ./scripts/stop-nginx.sh
	@echo "Stopping all services..."
	@for svc in $(SERVICES); do \
		if [ -f $(PIDS_DIR)/$$svc.pid ]; then \
			PID=$$(cat $(PIDS_DIR)/$$svc.pid); \
			if kill -0 $$PID 2>/dev/null; then \
				kill $$PID && echo "  ■ $$svc stopped (PID $$PID)"; \
			else \
				echo "  $$svc was not running"; \
			fi; \
			rm -f $(PIDS_DIR)/$$svc.pid; \
		else \
			echo "  $$svc not started (no PID file)"; \
		fi; \
	done

# Show running/stopped status for each service
status:
	@echo "Service status:"
	@for svc in $(SERVICES); do \
		if [ -f $(PIDS_DIR)/$$svc.pid ] && kill -0 $$(cat $(PIDS_DIR)/$$svc.pid) 2>/dev/null; then \
			echo "  ✅ $$svc running (PID $$(cat $(PIDS_DIR)/$$svc.pid))"; \
		else \
			echo "  ⬜ $$svc stopped"; \
		fi; \
	done

# Restart all services without rebuilding (useful after ssl-refresh)
restart: stop
	@echo "Restarting all services..."
	@mkdir -p $(PIDS_DIR) $(LOGS_DIR)
	@for svc in $(SERVICES); do \
		if [ -f $(PIDS_DIR)/$$svc.pid ] && kill -0 $$(cat $(PIDS_DIR)/$$svc.pid) 2>/dev/null; then \
			echo "  $$svc already running (PID $$(cat $(PIDS_DIR)/$$svc.pid))"; \
		else \
			nohup $(BIN_DIR)/$$svc > $(LOGS_DIR)/$$svc.log 2>&1 & echo $$! > $(PIDS_DIR)/$$svc.pid; \
			echo "  ▶ $$svc started (PID $$(cat $(PIDS_DIR)/$$svc.pid))"; \
		fi; \
	done
	@echo "Logs: $(LOGS_DIR)/"
	@echo "Starting nginx API gateway..."
	@bash ./scripts/start-nginx.sh ./nginx/nginx.yaml

# ─── Nginx Configuration ──────────────────────────────────────────────────────

# Generate nginx.conf from routes.go and nginx.yaml configuration
nginx-config: build-tools
	@echo "Generating nginx configuration..."
	@mkdir -p ./logs
	@$(BIN_DIR)/nginx-config-gen -yaml ./nginx/nginx.yaml -output ./nginx/nginx.conf
	@echo "✅ Nginx configuration generated"

# Regenerate CA, server, and default client certificates
ssl-refresh:
	@bash scripts/gen-ssl.sh refresh
	@if ls $(PIDS_DIR)/*.pid 2>/dev/null | grep -q .; then \
		$(MAKE) --no-print-directory restart; \
	fi

# Generate N numbered client certificates (default N=1)
# Usage: make ssl-clients N=5
ssl-clients:
	@bash scripts/gen-ssl.sh clients $(or $(N),1)

# ─── Tests ────────────────────────────────────────────────────────────────────

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
	@echo "Usage: make test-service SERVICE=DCAP"
	@if [ -z "$(SERVICE)" ]; then \
		echo "Available services: DCAP, crawler, flowreservation, operator, rsps, client, scaffold-gen, wadl-extract"; \
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
	@echo "SSL Certificate Commands:"
	@echo ""
	@echo "  make ssl-refresh             - Regenerate CA, server, and default client certs (restarts running services)"
	@echo "  make ssl-clients N=5         - Generate N numbered client certs (client-0001...)"
	@echo ""
	@echo "Build Commands:"
	@echo ""
	@echo "  make build-all               - Build all services and tools to ./bin/"
	@echo "  make build-services          - Build all microservice servers to ./bin/"
	@echo "  make build-tools             - Build crawler, client, scaffold-gen, wadl-extract"
	@echo "  make build-<name>            - Build a single service (e.g. make build-DCAP)"
	@echo ""
	@echo "Service Management:"
	@echo ""
	@echo "  make start                   - Build and start all microservices in background"
	@echo "  make stop                    - Stop all running microservices"
	@echo "  make restart                 - Restart all services without rebuilding"
	@echo "  make status                  - Show running/stopped status for each service"
	@echo ""
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
	@echo "Route Testing Commands:"
	@echo "  make test-routes             - Run all route registration & HTTP method tests"
	@echo "  make test-route-registration - Run route registration verification tests"
	@echo "  make test-http-methods       - Run HTTP method support tests"
	@echo ""
	@echo "Cleanup Commands:"
	@echo "  make cleanup-services-list   - List all services and identify stale ones"
	@echo "  make cleanup-services-dry-run - Preview what would be deleted (no changes)"
	@echo "  make cleanup-services        - Remove stale services (with confirmation)"
	@echo ""
	@echo "  make generate-mocks          - Generate test mocks (mockgen)"
	@echo "  make lint                    - Run linters"
	@echo "  make help                    - Show this help message"
