.PHONY: help build up up-build down restart logs logs-api logs-db clean db-reset db-backup db-restore shell-api shell-db admin status health \
	fmt lint test tidy run all

# Variables
COMPOSE_FILE = docker-compose.yaml
PROJECT_NAME = my-api

# --- Help ---
help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Docker Compose ---
build: ## Build Docker images
	docker-compose -f $(COMPOSE_FILE) build --no-cache

build-optimized: ## Build optimized production Docker image
	./scripts/docker-build.sh

build-prod: ## Build production Docker image with tag
	./scripts/docker-build.sh prod

up: ## Start all services
	docker-compose -f $(COMPOSE_FILE) up -d

up-build: ## Build and start all services
	docker-compose -f $(COMPOSE_FILE) up -d --build

down: ## Stop all services
	docker-compose -f $(COMPOSE_FILE) down

restart: ## Restart all services
	docker-compose -f $(COMPOSE_FILE) restart

logs: ## Show logs for all services
	docker-compose -f $(COMPOSE_FILE) logs -f

logs-api: ## Show logs for the API service
	docker-compose -f $(COMPOSE_FILE) logs -f api

logs-db: ## Show logs for PostgreSQL
	docker-compose -f $(COMPOSE_FILE) logs -f postgres

clean: ## Remove unused containers, images, and volumes
	docker-compose -f $(COMPOSE_FILE) down -v --remove-orphans
	docker system prune -f
	docker volume prune -f

db-reset: ## Reset database (WARNING: deletes all data)
	docker-compose -f $(COMPOSE_FILE) down -v
	docker volume rm $(PROJECT_NAME)_postgres_data || true
	docker-compose -f $(COMPOSE_FILE) up -d postgres

db-backup: ## Create a database backup
	mkdir -p backups
	docker-compose -f $(COMPOSE_FILE) exec postgres pg_dump -U api_user my_api_db > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

db-restore: ## Restore a backup (usage: make db-restore FILE=backup.sql)
	docker-compose -f $(COMPOSE_FILE) exec -T postgres psql -U api_user -d my_api_db < $(FILE)

shell-api: ## Access the API container shell
	docker-compose -f $(COMPOSE_FILE) exec api /bin/sh

shell-db: ## Access the PostgreSQL shell
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U api_user -d my_api_db

admin: ## Start with pgAdmin included
	docker-compose -f $(COMPOSE_FILE) --profile admin up -d

status: ## Show status of services
	docker-compose -f $(COMPOSE_FILE) ps

health: ## Check health of services
	@echo "=== Service status ==="
	@docker-compose -f $(COMPOSE_FILE) ps
	@echo "\n=== Health checks ==="
	@docker-compose -f $(COMPOSE_FILE) exec api wget -q --spider http://localhost:8080/health && echo "✅ API: OK" || echo "❌ API: FAIL"
	@docker-compose -f $(COMPOSE_FILE) exec postgres pg_isready -U api_user -d my_api_db && echo "✅ PostgreSQL: OK" || echo "❌ PostgreSQL: FAIL"

# --- Go Development ---
fmt: ## Format code and organize imports
	goimports -w .
	go fmt ./...

lint: ## Linting with golangci-lint (requires prior installation)
	golangci-lint run ./...

test: ## Run tests
	go test ./...

tidy: ## Clean and update dependencies
	go mod tidy

run: ## Run the main application
	go run ./cmd/api/main.go

all: fmt lint test ## Run all: format, lint, and test

# --- Production Deployment ---
prod-up: ## Start production environment
	docker-compose -f docker-compose.prod.yml up -d

prod-down: ## Stop production environment
	docker-compose -f docker-compose.prod.yml down

prod-logs: ## Show production logs
	docker-compose -f docker-compose.prod.yml logs -f

prod-status: ## Check production status
	docker-compose -f docker-compose.prod.yml ps

prod-health: ## Check production health
	@echo "=== Production Health Check ==="
	@curl -s http://localhost:8080/health/ | jq . || echo "Health check failed"

image-size: ## Show Docker image sizes
	@echo "=== Docker Image Sizes ==="
	@docker images gin-template --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

image-scan: ## Scan Docker image for vulnerabilities (requires trivy)
	@trivy image gin-template:latest 2>/dev/null || echo "Install trivy for vulnerability scanning: brew install trivy"
