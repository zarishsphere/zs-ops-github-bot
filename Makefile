# ZarishSphere GitHub Bot Makefile

.PHONY: help build test clean docker-build docker-push deploy helm-install helm-upgrade

# Default target
help: ## Show this help message
	@echo "ZarishSphere GitHub Bot"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the Go binary
	@echo "Building Go binary..."
	go build -o bin/bot ./cmd/bot

test: ## Run all tests
	@echo "Running tests..."
	go test ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t zarishsphere/zs-ops-github-bot:latest -f deploy/Dockerfile .

docker-push: ## Push Docker image
	@echo "Pushing Docker image..."
	docker push zarishsphere/zs-ops-github-bot:latest

# Deployment targets
deploy: docker-build docker-push ## Build and deploy
	@echo "Deployment complete"

# Helm targets
helm-install: ## Install Helm chart
	@echo "Installing Helm chart..."
	helm install zs-ops-github-bot ./deploy/helm/zs-ops-github-bot

helm-upgrade: ## Upgrade Helm chart
	@echo "Upgrading Helm chart..."
	helm upgrade zs-ops-github-bot ./deploy/helm/zs-ops-github-bot

helm-uninstall: ## Uninstall Helm chart
	@echo "Uninstalling Helm chart..."
	helm uninstall zs-ops-github-bot

# Development targets
dev: ## Run in development mode
	@echo "Starting development server..."
	go run ./cmd/bot

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# CI targets
ci: fmt vet lint test ## Run full CI pipeline