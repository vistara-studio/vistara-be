# Vistara Backend Makefile
# Comprehensive development and deployment autrestart: ## 🔄 Restart all services
	@echo "🔄 Restarting services..."
	@docker compose restart
	@echo "✅ Services restarted!"

# === BASIC COMMANDS ===
dev: ## 🚀 Start development server
	@echo "🚀 Starting development server..."
	@docker compose up -d
	@echo "✅ Server running at http://localhost/"

prod: ## 🔐 Start production server (requires SSL setup)
	@echo "🔐 Starting production server..."
	@docker compose up -d
	@echo "✅ Production server running!"

setup-ssl: ## 🔐 Manual SSL setup guide
	@echo "🔐 SSL Certificate Setup Guide for GCP VM + Biznetgio Domain"
	@echo "==========================================================="
	@echo ""
	@echo "📋 Prerequisites:"
	@echo "1. GCP VM running with external IP"
	@echo "2. Domain einrafh.com from Biznetgio"
	@echo "3. Ports 80 and 443 open in GCP firewall"
	@echo ""
	@echo "🌐 DNS Setup (in Biznetgio panel):"
	@echo "   A Record: einrafh.com → [your_gcp_vm_ip]"
	@echo "   A Record: www.einrafh.com → [your_gcp_vm_ip]"
	@echo ""
	@echo "🔐 SSL Certificate (run on your GCP VM):"
	@echo "   1. sudo apt update && sudo apt install certbot"
	@echo "   2. make dev  # Start the server first"
	@echo "   3. sudo certbot certonly --webroot -w ./nginx/certbot -d einrafh.com -d www.einrafh.com"
	@echo "   4. sudo cp /etc/letsencrypt/live/einrafh.com/fullchain.pem ./nginx/ssl/einrafh.com.crt"
	@echo "   5. sudo cp /etc/letsencrypt/live/einrafh.com/privkey.pem ./nginx/ssl/einrafh.com.key"
	@echo "   6. sudo chown $$USER:$$USER ./nginx/ssl/einrafh.com.*"
	@echo "   7. make prod  # Restart with SSL"
	@echo ""

setup-dev-ssl: ## 🔧 Setup self-signed SSL for development (no hosting required)
	@echo "🔧 Setting up development SSL certificates..."
	@./scripts/setup-dev-ssl.sh

prod-env: ## 📄 Copy production environment template
	@echo "📄 Copying production environment template..."
	@cp .env.production .env
	@echo "✅ Environment template copied to .env"
	@echo "📝 Please edit .env file with your actual credentials"

deploy: ## 🚀 Complete deployment (automated)
	@echo "🚀 Starting automated deployment..."
	@./scripts/deploy.sh

nginx-logs: ## 📋 Show nginx logs
	@echo "📋 Showing nginx logs..."
	@docker compose logs -f nginx

nginx-reload: ## 🔄 Reload nginx configuration
	@echo "🔄 Reloading nginx configuration..."
	@docker compose exec nginx nginx -s reload
	@echo "✅ Nginx configuration reloaded!"

# === MONITORING COMMANDS ===n

include .env

.PHONY: help setup reset-setup dev-setup build run test clean docker-build docker-run docker-clean logs status health info

deploy: ## 🚀 Complete deployment for production VM
	@echo "🚀 Deploying to production VM..."
	@echo "================================"
	@echo "📋 Checking prerequisites..."
	@which docker > /dev/null || (echo "❌ Docker not found. Install: curl -fsSL https://get.docker.com | sh" && exit 1)
	@which docker-compose > /dev/null || docker compose version > /dev/null || (echo "❌ Docker Compose not found" && exit 1)
	@echo "✅ Docker ready"
	@echo ""
	@echo "📂 Setting up environment..."
	@cp .env.production .env || echo "⚠️  .env.production not found, using existing .env"
	@echo "✅ Environment configured"
	@echo ""
	@echo "🏗️ Building and starting services..."
	@docker compose up --build -d
	@echo ""
	@echo "🎉 Deployment complete!"
	@echo "📍 Server running at: http://[your-vm-ip]:80"
	@echo ""
	@echo "🔐 For HTTPS setup, run: make setup-ssl"

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