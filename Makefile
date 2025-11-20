.PHONY: help build run test lint migrate migrate-down migrate-up migrate-status clean docker-build docker-up docker-down docker-logs install-deps swagger fmt vet

# Variables
APP_NAME := feedbacklab
HEALTHCHECK_NAME := healthcheck
BIN_DIR := bin
CMD_DIR := cmd
MIGRATIONS_DIR := migrations
SWAGGER_DIR := docs
GO_VERSION := 1.21
DOCKER_COMPOSE := docker-compose -f docker-compose.local.yml

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

##@ General

help: ## Display this help message
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

##@ Development

install-deps: ## Install Go dependencies
	@echo "$(GREEN)Installing dependencies...$(NC)"
	go mod download
	go mod tidy

run: build ## Build and run the application
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	./$(BIN_DIR)/$(APP_NAME)

run-dev: ## Run the application in development mode (with hot reload if you have air/realize)
	@echo "$(GREEN)Running $(APP_NAME) in development mode...$(NC)"
	go run $(CMD_DIR)/$(APP_NAME)/main.go

##@ Building

build: ## Build the application binaries
	@echo "$(GREEN)Building $(APP_NAME) and $(HEALTHCHECK_NAME)...$(NC)"
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)/$(APP_NAME)
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(HEALTHCHECK_NAME) ./$(CMD_DIR)/healthcheck
	@echo "$(GREEN)Build complete!$(NC)"

build-linux: ## Build for Linux (amd64)
	@echo "$(GREEN)Building for Linux...$(NC)"
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux ./$(CMD_DIR)/$(APP_NAME)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(HEALTHCHECK_NAME)-linux ./$(CMD_DIR)/healthcheck

build-windows: ## Build for Windows (amd64)
	@echo "$(GREEN)Building for Windows...$(NC)"
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME).exe ./$(CMD_DIR)/$(APP_NAME)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(HEALTHCHECK_NAME).exe ./$(CMD_DIR)/healthcheck

build-darwin: ## Build for macOS (amd64)
	@echo "$(GREEN)Building for macOS...$(NC)"
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-darwin ./$(CMD_DIR)/$(APP_NAME)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(HEALTHCHECK_NAME)-darwin ./$(CMD_DIR)/healthcheck

build-all: build-linux build-windows build-darwin ## Build for all platforms

##@ Testing

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...

test-cover: test ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"


##@ Code Quality

lint: ## Run golangci-lint
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(GREEN)Running linter with auto-fix...$(NC)"
	golangci-lint run ./... --fix

fmt: ## Format code with gofmt
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	go vet ./...

check: fmt vet lint ## Run all code quality checks

##@ Database

migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations up...$(NC)"
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "$(RED)Error: DATABASE_URL is not set$(NC)"; \
		exit 1; \
	fi
	go run $(CMD_DIR)/feedbacklab/main.go --migrate-only || \
	(go run ./pkg/db/migrate.go || echo "Run: go run -tags migrate ./cmd/feedbacklab")

migrate: migrate-up ## Alias for migrate-up

migrate-create: ## Create a new migration file
	@echo "$(GREEN)Creating new migration...$(NC)"
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=create_users_table$(NC)"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

##@ Swagger

swagger: ## Generate Swagger documentation
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@which swag > /dev/null || (echo "$(RED)swag is not installed. Install it with: go install github.com/swaggo/swag/cmd/swag@latest$(NC)" && exit 1)
	swag init -g $(CMD_DIR)/$(APP_NAME)/main.go -o $(SWAGGER_DIR)
	@echo "$(GREEN)Swagger documentation generated in $(SWAGGER_DIR)/$(NC)"

##@ Docker

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(APP_NAME):latest .

docker-up: ## Start Docker Compose services
	@echo "$(GREEN)Starting Docker Compose services...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Services started. Use 'make docker-logs' to view logs$(NC)"

docker-down: ## Stop Docker Compose services
	@echo "$(YELLOW)Stopping Docker Compose services...$(NC)"
	$(DOCKER_COMPOSE) down

docker-logs: ## View Docker Compose logs
	$(DOCKER_COMPOSE) logs -f

docker-restart: docker-down docker-up ## Restart Docker Compose services

docker-clean: docker-down ## Stop and remove Docker containers and volumes
	@echo "$(YELLOW)Removing Docker volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

docker-ps: ## Show running Docker containers
	$(DOCKER_COMPOSE) ps


##@ Default

.DEFAULT_GOAL := help

