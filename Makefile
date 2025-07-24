# Vistara Backend Makefile
# Comprehensive development and deployment commands

include .env

.PHONY: help setup reset-setup build run test clean docker-build docker-run docker-clean logs dev-setup production

# Default target
help: ## Show this help message
	@echo "Vistara Backend - Available Commands:"
	@echo "===================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development Setup
setup: ## 🚀 Complete setup with test data (ONE COMMAND SETUP)
	@echo "🚀 Starting Vistara Backend complete setup..."
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

reset-setup: ## 🔄 Reset and clean environment completely
	@echo "🔄 Resetting Vistara Backend environment..."
	@chmod +x scripts/reset-setup.sh
	@./scripts/reset-setup.sh

dev-setup: ## 🛠️ Setup for development (without test data)
	@echo "🛠️ Setting up development environment..."
	@docker compose up --build -d
	@echo "✅ Development environment ready!"

# Build Commands
build: ## 🔨 Build the Go application
	@echo "🔨 Building Vistara Backend..."
	@go mod tidy
	@go build -o vistara-backend cmd/api/main.go
	@echo "✅ Build completed!"

test: ## 🧪 Run tests
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ Tests completed!"

# Docker Commands
docker-build: ## 🐳 Build Docker images
	@echo "🐳 Building Docker images..."
	@docker compose build
	@echo "✅ Docker images built!"

docker-run: ## 🐳 Run with Docker Compose
	@echo "🐳 Starting Docker containers..."
	@docker compose up -d
	@echo "✅ Containers started!"

docker-clean: ## 🧹 Clean Docker containers and images
	@echo "🧹 Cleaning Docker resources..."
	@docker compose down --remove-orphans
	@docker system prune -f
	@echo "✅ Docker cleanup completed!"

# Utility Commands
logs: ## 📋 Show application logs
	@echo "📋 Showing application logs..."
	@docker compose logs -f app

logs-all: ## 📋 Show all service logs
	@echo "📋 Showing all service logs..."
	@docker compose logs -f

db-logs: ## 📋 Show database logs
	@echo "📋 Showing database logs..."
	@docker compose logs -f postgres

status: ## 📊 Show service status
	@echo "📊 Service Status:"
	@docker compose ps

health: ## ❤️ Check API health
	@echo "❤️ Checking API health..."
	@curl -s http://localhost:8080/health || echo "❌ API not responding"

# Database Commands - Legacy support
migrate-up: ## 🗃️ Run database migrations (legacy)
	@docker compose run --rm migrate -path /db/migrations -database "postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose up

migrate-down: ## 🔄 Rollback database migrations (legacy)
	@echo "y" | docker compose run --rm -T migrate -path /db/migrations -database "postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose down

migrate-drop: ## 🗑️ Drop all database tables (legacy)
	@echo "y" | docker compose run --rm -T migrate -path /db/migrations -database "postgresql://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose drop

# Legacy Docker Commands
compose-up: ## 🐳 Start services (legacy)
	@docker compose up --detach --build

compose-down: ## 🛑 Stop services (legacy)
	@docker compose down

# Development Commands
run: ## 🏃 Run application locally
	@echo "🏃 Running Vistara Backend locally..."
	@go run cmd/api/main.go

dev: ## 🔄 Run in development mode with hot reload
	@echo "🔄 Starting development mode..."
	@docker compose up --build

# Cleanup Commands
clean: ## 🧹 Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -f vistara-backend
	@go clean
	@echo "✅ Cleanup completed!"

stop: ## 🛑 Stop all services
	@echo "🛑 Stopping all services..."
	@docker compose down
	@echo "✅ Services stopped!"

restart: ## 🔄 Restart all services
	@echo "🔄 Restarting all services..."
	@docker compose restart
	@echo "✅ Services restarted!"

# API Testing Commands
test-auth: ## 🔐 Test authentication endpoints
	@echo "🔐 Testing authentication..."
	@curl -X POST http://localhost:8080/api/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"testuser1@vistara.com","password":"password123"}' | jq . || echo "Response received"

# AI Integration Testing Commands
test-ai-integration: ## 🤖 Test AI integration endpoints
	@echo "🤖 Testing AI integration..."
	@curl -X POST http://localhost:8080/api/ai/smart-plan \
		-H "Content-Type: application/json" \
		-d '{"destination":"Bali","start_date":"2025-08-01T00:00:00Z","end_date":"2025-08-05T00:00:00Z","budget":5000000,"travel_style":"romantic_couple","activity_preferences":["beach","culture","culinary"],"activity_intensity":"balanced"}' | jq . || echo "Response received"

test-service-endpoints: ## 🔗 Test service-to-service endpoints (for vistara-ai)
	@echo "🔗 Testing service endpoints for vistara-ai..."
	@echo "📍 Testing local businesses endpoint:"
	@curl -X GET http://localhost:8080/api/service/locals \
		-H "X-Service: vistara-ai" | jq . || echo "Response received"
	@echo ""
	@echo "📍 Testing tourist attractions endpoint:"
	@curl -X GET http://localhost:8080/api/service/tourist-attractions \
		-H "X-Service: vistara-ai" | jq . || echo "Response received"

test-notification: ## 🔔 Test AI notification endpoint
	@echo "🔔 Testing AI notification endpoint..."
	@curl -X POST http://localhost:8080/api/service/ai/notify \
		-H "X-Service: vistara-ai" \
		-H "Content-Type: application/json" \
		-d '{"event":"plan_generated","user_id":"test-user","data":{"destination":"Bali"},"timestamp":"2025-07-24T10:00:00Z"}' | jq . || echo "Response received"

# Information Commands
info: ## ℹ️ Show system information
	@echo "ℹ️ Vistara Backend Information:"
	@echo "=============================="
	@echo "📂 Project: Vistara Backend API"
	@echo "🐹 Language: Go 1.23"
	@echo "🌐 Framework: Fiber v2"
	@echo "🗄️ Database: PostgreSQL 17"
	@echo "🐳 Container: Docker + Docker Compose"
	@echo "🚀 Port: 8080"
	@echo "📍 Base URL: http://localhost:8080"
	@echo ""
	@echo "📋 Available Endpoints:"
	@echo "• Authentication: /api/auth/*"
	@echo "• Local Business: /api/locals/*"
	@echo "• Tourist Attractions: /api/tourist-attractions/*"
	@echo "• AI Integration: /api/ai/*"
	@echo "• Service Endpoints: /api/service/*"
	@echo "• Health Check: /health"

# Quick shortcuts
up: docker-run ## 🔼 Quick start (alias for docker-run)
down: stop ## 🔽 Quick stop (alias for stop)
reset: reset-setup ## 🔄 Reset environment (alias for reset-setup)
ps: status ## 📋 Show status (alias for status)