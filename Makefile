.PHONY: help build run test lint migrate migrate-down migrate-up migrate-status clean docker-build docker-up docker-down docker-logs install-deps swagger fmt vet

# Variables
APP_NAME := feedbacklab
HEALTHCHECK_NAME := healthcheck
BIN_DIR := bin
CMD_DIR := cmd
MIGRATIONS_DIR := migrations
SWAGGER_DIR := docs
GO_VERSION := 1.25
DOCKER_COMPOSE := docker-compose -f docker-compose.local.yml
DC_CMD := $(DOCKER_COMPOSE)

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

swagger: ## Generate OpenAPI 3.0 documentation from annotations
	@echo "$(GREEN)Generating Swagger documentation from annotations...$(NC)"
	@which swag > /dev/null || (echo "$(RED)swag is not installed. Installing...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest)
	@export PATH=$$PATH:$$(go env GOPATH)/bin:$$HOME/go/bin && swag init -g $(CMD_DIR)/$(APP_NAME)/main.go -o $(SWAGGER_DIR) --parseDependency --parseInternal
	@echo "$(GREEN)Converting Swagger 2.0 to OpenAPI 3.0...$(NC)"
	@which swagger2openapi > /dev/null || (echo "$(RED)swagger2openapi is not installed. Installing...$(NC)" && npm install -g swagger2openapi)
	@swagger2openapi $(SWAGGER_DIR)/swagger.json -o $(SWAGGER_DIR)/openapi-generated.json 2>/dev/null || echo "$(YELLOW)Warning: Could not convert swagger.json to OpenAPI 3.0$(NC)"
	@echo "$(GREEN)Using manually created OpenAPI 3.0 specification...$(NC)"
	@if [ -f $(SWAGGER_DIR)/openapi.json ]; then \
		echo "$(GREEN)✓ Found openapi.json, using it as the source$(NC)"; \
		cp $(SWAGGER_DIR)/openapi.json $(SWAGGER_DIR)/swagger.json; \
		cp $(SWAGGER_DIR)/openapi.json $(SWAGGER_DIR)/swagger.yaml 2>/dev/null || echo "$(YELLOW)Note: Could not create swagger.yaml$(NC)"; \
	else \
		echo "$(YELLOW)Warning: openapi.json not found, using generated swagger.json$(NC)"; \
	fi
	@echo "$(GREEN)OpenAPI 3.0 documentation ready in $(SWAGGER_DIR)/$(NC)"
	@echo "$(GREEN)Access Swagger UI at: http://localhost:8080/swagger/index.html$(NC)"

##@ Docker

docker-build: ## Build Docker image force
	@echo "$(GREEN)Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build --no-cache

docker-up: ## Start Docker Compose services (rebuilds if necessary)
	@echo "$(GREEN)Starting Docker Compose services...$(NC)"
	@if [ ! -f .env.docker ]; then \
		echo "$(YELLOW).env.docker file not found. Creating from template...$(NC)"; \
		make docker-env; \
	fi
	@echo "$(GREEN)Building and starting containers...$(NC)"
	$(DOCKER_COMPOSE) up -d --build
	@echo "$(GREEN)Services started! Waiting for services to be ready...$(NC)"
	@echo "$(YELLOW)This may take a minute or two...$(NC)"
	@sleep 10
	@echo "$(GREEN)Checking service health...$(NC)"
	@$(DOCKER_COMPOSE) ps
	@echo ""
	@echo "$(GREEN)=== Service URLs ===$(NC)"
	@echo "$(GREEN)App: http://localhost:$$(grep APP_PORT .env.docker | cut -d'=' -f2 || echo 8080)$(NC)"
	@echo "$(GREEN)Swagger: http://localhost:$$(grep APP_PORT .env.docker | cut -d'=' -f2 || echo 8080)/swagger/index.html$(NC)"
	@echo "$(GREEN)Mattermost: http://localhost:$$(grep MM_PORT .env.docker | cut -d'=' -f2 || echo 8065)$(NC)"
	@echo "$(GREEN)Keycloak: http://localhost:$$(grep KEYCLOAK_PORT .env.docker | cut -d'=' -f2 || echo 8082)$(NC)"
	@echo "$(GREEN)MinIO Console: http://localhost:$$(grep MINIO_CONSOLE_PORT .env.docker | cut -d'=' -f2 || echo 9001)$(NC)"
	@echo ""
	@echo "$(YELLOW)Use 'make docker-logs' to view logs$(NC)"

docker-down: ## Stop Docker Compose services
	@echo "$(YELLOW)Stopping Docker Compose services...$(NC)"
	$(DOCKER_COMPOSE) down

docker-logs: ## View Docker Compose logs
	$(DOCKER_COMPOSE) logs -f

docker-restart: docker-down docker-up ## Restart Docker Compose services

docker-clean: ## Stop and remove Docker containers AND volumes (FULL RESET)
	@echo "$(RED)Stopping and removing containers and VOLUMES...$(NC)"
	$(DOCKER_COMPOSE) down -v
	@echo "$(YELLOW)Pruning system...$(NC)"
	docker system prune -f
	@echo "$(GREEN)Clean complete. Next start will re-initialize databases.$(NC)"

docker-ps: ## Show running Docker containers
	$(DOCKER_COMPOSE) ps

docker-logs-app: ## View logs for app container
	$(DOCKER_COMPOSE) logs -f app

docker-logs-db: ## View logs for database container
	$(DOCKER_COMPOSE) logs -f db

docker-shell-app: ## Open shell in app container
	docker exec -it feedbacklab /bin/sh

docker-shell-db: ## Open psql shell in database container
	@POSTGRES_USER=$$(grep POSTGRES_USER .env.docker | cut -d'=' -f2 || echo feedback); \
	docker exec -it feedback_db psql -U $$POSTGRES_USER -d innotech

docker-status: ## Check health status of all services
	@echo "$(GREEN)Checking service status...$(NC)"
	@$(DOCKER_COMPOSE) ps
	@echo ""
	@echo "$(GREEN)Checking service health...$(NC)"
	@for service in db keycloak minio app mattermost; do \
		echo -n "$$service: "; \
		$(DOCKER_COMPOSE) ps $$service | grep -q "healthy\|Up" && echo "$(GREEN)✓$(NC)" || echo "$(RED)✗$(NC)"; \
	done

docker-env: ## Create .env.docker file from template
	@echo "$(GREEN)Creating .env.docker file...$(NC)"
	@echo "# Database Configuration" > .env.docker
	@echo "POSTGRES_USER=feedback" >> .env.docker
	@echo "POSTGRES_PASSWORD=feedback" >> .env.docker
	@echo "POSTGRES_DB=innotech" >> .env.docker
	@echo "DB_PORT=5432" >> .env.docker
	@echo "" >> .env.docker
	@echo "# Mattermost Database Configuration" >> .env.docker
	@echo "MM_DB_USER=feedback" >> .env.docker
	@echo "MM_DB_PASSWORD=feedback" >> .env.docker
	@echo "MM_SITEURL=http://localhost:8065" >> .env.docker
	@echo "MM_PORT=8065" >> .env.docker
	@echo "" >> .env.docker
	@echo "# Keycloak Configuration" >> .env.docker
	@echo "KEYCLOAK_ADMIN=admin" >> .env.docker
	@echo "KEYCLOAK_ADMIN_PASSWORD=admin" >> .env.docker
	@echo "KEYCLOAK_PORT=8082" >> .env.docker
	@echo "" >> .env.docker
	@echo "# Application Configuration" >> .env.docker
	@echo "APP_PORT=8080" >> .env.docker
	@echo "HEALTH_PORT=8081" >> .env.docker
	@echo "DATABASE_URL=postgres://feedback:feedback@db:5432/innotech?sslmode=disable" >> .env.docker
	@echo "MIGRATIONS_DIR=/app/migrations" >> .env.docker
	@echo "" >> .env.docker
	@echo "# Swagger Configuration" >> .env.docker
	@echo "SWAGGER_USERNAME=swagger" >> .env.docker
	@echo "SWAGGER_PASSWORD=swagger" >> .env.docker
	@echo "" >> .env.docker
	@echo "# MinIO Configuration" >> .env.docker
	@echo "MINIO_ENDPOINT=minio:9000" >> .env.docker
	@echo "MINIO_ACCESS_KEY=minioadmin" >> .env.docker
	@echo "MINIO_SECRET_KEY=minioadmin123" >> .env.docker
	@echo "MINIO_BUCKET=feedback-files" >> .env.docker
	@echo "MINIO_SSL=false" >> .env.docker
	@echo "MINIO_PORT=9000" >> .env.docker
	@echo "MINIO_CONSOLE_PORT=9001" >> .env.docker
	@echo "" >> .env.docker
	@echo "# Mattermost Webhook" >> .env.docker
	@echo "MATTERMOST_WEBHOOK=" >> .env.docker
	@echo "$(GREEN).env.docker file created!$(NC)"

docker-init-db: ## Manually init databases inside running container
	@echo "$(GREEN)Creating databases manually...$(NC)"
	@if [ -z "$$(docker ps -q -f name=feedback_db)" ]; then \
		echo "$(RED)Error: Database container is not running. Start it first with 'make docker-up'$(NC)"; \
		exit 1; \
	fi
	@POSTGRES_USER=$$(grep POSTGRES_USER .env.docker | cut -d'=' -f2 || echo feedback); \
	docker exec feedback_db psql -U $$POSTGRES_USER -d postgres -c "CREATE DATABASE keycloak;" 2>/dev/null || echo "Database keycloak already exists"; \
	docker exec feedback_db psql -U $$POSTGRES_USER -d postgres -c "CREATE DATABASE mattermost;" 2>/dev/null || echo "Database mattermost already exists"; \
	docker exec feedback_db psql -U $$POSTGRES_USER -d keycloak -c "GRANT ALL PRIVILEGES ON DATABASE keycloak TO $$POSTGRES_USER;" 2>/dev/null || true; \
	docker exec feedback_db psql -U $$POSTGRES_USER -d mattermost -c "GRANT ALL PRIVILEGES ON DATABASE mattermost TO $$POSTGRES_USER;" 2>/dev/null || true; \
	echo "$(GREEN)Databases initialized. Restarting services...$(NC)"; \
	$(DC_CMD) restart keycloak mattermost || true

##@ Default

.DEFAULT_GOAL := help