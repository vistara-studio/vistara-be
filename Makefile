# Vistara Backend Makefile
# Comprehensive development and deployment automation

include .env

.PHONY: help setup reset-setup dev-setup build run test clean docker-build docker-run docker-clean logs status health info

# Default target
help: ## 📋 Show available commands
	@echo "🏛️ Vistara Backend - Available Commands"
	@echo "========================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?##/ {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# === SETUP COMMANDS ===
setup: ## 🚀 Complete setup with test data (recommended)
	@echo "🚀 Starting complete Vistara Backend setup..."
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

reset-setup: ## 🔄 Reset and clean environment completely  
	@echo "🔄 Resetting Vistara Backend environment..."
	@chmod +x scripts/reset-setup.sh
	@./scripts/reset-setup.sh

dev-setup: ## 🛠️ Basic development setup (no test data)
	@echo "🛠️ Setting up development environment..."
	@docker compose up --build -d
	@echo "✅ Development environment ready at http://localhost:8080"

# === BUILD COMMANDS ===
build: ## 🔨 Build the Go application
	@echo "🔨 Building Vistara Backend..."
	@go mod tidy
	@go build -o vistara-backend cmd/api/main.go
	@echo "✅ Build completed!"

test: ## 🧪 Run all tests
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ Tests completed!"

# === DOCKER COMMANDS ===
docker-build: ## 🐳 Build Docker images
	@echo "🐳 Building Docker images..."
	@docker compose build
	@echo "✅ Docker images built!"

docker-run: ## 🐳 Start services with Docker Compose
	@echo "🐳 Starting Docker containers..."
	@docker compose up -d
	@echo "✅ Containers started!"

docker-clean: ## 🧹 Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	@docker compose down --remove-orphans
	@docker system prune -f
	@echo "✅ Docker cleanup completed!"

# === RUNTIME COMMANDS ===
run: ## 🏃 Run application locally (development)
	@echo "🏃 Running Vistara Backend locally..."
	@go run cmd/api/main.go

dev: ## � Start development mode with hot reload
	@echo "� Starting development mode..."
	@docker compose up --build

start: docker-run ## ▶️ Start services (alias for docker-run)
stop: ## ⏹️ Stop all services
	@echo "⏹️ Stopping services..."
	@docker compose down
	@echo "✅ Services stopped!"

restart: ## 🔄 Restart all services
	@echo "� Restarting services..."
	@docker compose restart
	@echo "✅ Services restarted!"

# === MONITORING COMMANDS ===
logs: ## � Show application logs
	@echo "📋 Showing application logs..."
	@docker compose logs -f app

logs-all: ## � Show all service logs
	@echo "📋 Showing all service logs..."
	@docker compose logs -f

db-logs: ## 📋 Show database logs
	@echo "📋 Showing database logs..."
	@docker compose logs -f postgres

status: ## � Show service status
	@echo "� Service Status:"
	@docker compose ps

health: ## ❤️ Check API health
	@echo "❤️ Checking API health..."
	@curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health || echo "❌ API not responding"

# === TESTING COMMANDS ===
test-auth: ## 🔐 Test authentication endpoints
	@echo "🔐 Testing authentication..."
	@curl -X POST http://localhost:8080/api/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"testuser1@vistara.com","password":"password123"}' | jq . 2>/dev/null || echo "Response received"

test-ai: ## 🤖 Test AI integration endpoints
	@echo "🤖 Testing AI integration..."
	@curl -X POST http://localhost:8080/api/ai/smart-plan \
		-H "Content-Type: application/json" \
		-d '{"destination":"Bali","start_date":"2025-08-01T00:00:00Z","end_date":"2025-08-05T00:00:00Z","budget":5000000,"travel_style":"romantic_couple","activity_preferences":["beach","culture","culinary"],"activity_intensity":"balanced"}' | jq . 2>/dev/null || echo "Response received"

test-service: ## 🔗 Test service endpoints (for vistara-ai)
	@echo "🔗 Testing service endpoints for vistara-ai..."
	@echo "📍 Local businesses:"
	@curl -X GET http://localhost:8080/api/service/locals \
		-H "X-Service: vistara-ai" | jq . 2>/dev/null || echo "Response received"
	@echo ""
	@echo "📍 Tourist attractions:"
	@curl -X GET http://localhost:8080/api/service/tourist-attractions \
		-H "X-Service: vistara-ai" | jq . 2>/dev/null || echo "Response received"

test-notification: ## 🔔 Test AI notification endpoint
	@echo "🔔 Testing AI notification endpoint..."
	@curl -X POST http://localhost:8080/api/service/ai/notify \
		-H "X-Service: vistara-ai" \
		-H "Content-Type: application/json" \
		-d '{"event":"plan_generated","user_id":"test-user","data":{"destination":"Bali"},"timestamp":"2025-07-25T10:00:00Z"}' | jq . 2>/dev/null || echo "Response received"

test-all: test-auth test-ai test-service test-notification ## 🧪 Run all endpoint tests

# === UTILITY COMMANDS ===
clean: ## 🧹 Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -f vistara-backend
	@go clean
	@echo "✅ Cleanup completed!"

info: ## ℹ️ Show system information
	@echo "ℹ️ Vistara Backend Information"
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

# === QUICK SHORTCUTS ===
up: docker-run     ## 🔼 Quick start (alias for docker-run)
down: stop         ## 🔽 Quick stop (alias for stop)
reset: reset-setup ## 🔄 Reset environment (alias for reset-setup)
ps: status         ## 📋 Show status (alias for status)
log: logs          ## 📋 Show logs (alias for logs)