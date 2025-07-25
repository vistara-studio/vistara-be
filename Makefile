# Vistara Backend Makefile
# Comprehensive development and deployment automation with AI integration

.PHONY: help setup reset-setup dev-setup build run test clean docker-build docker-run docker-clean logs status health info ai-test

# Default target
help: ## 📋 Show available commands
	@echo "🏛️ Vistara Backend - Available Commands"
	@echo "========================================"
	@echo "🤖 Includes AI Integration (NusaLingo, Historical Stories, Smart Planning)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?##/ {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# === SETUP COMMANDS ===
setup: ## 🚀 Complete setup with test data and AI integration
	@echo "🚀 Starting complete Vistara Backend setup with AI integration..."
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

reset-setup: ## 🔄 Reset and clean environment completely
	@echo "🔄 Resetting environment..."
	@chmod +x scripts/reset-setup.sh
	@./scripts/reset-setup.sh

dev-setup: ## 🛠️ Basic development setup (no test data)
	@echo "🛠️ Starting basic development setup..."
	@if [ ! -f .env ]; then \
		echo "📄 Creating .env from template..."; \
		cp .env.example .env; \
		echo "📝 Please edit .env with your configuration"; \
	fi
	@docker compose up --build -d
	@echo "✅ Development setup complete!"

# === BASIC COMMANDS ===
dev: ## 🚀 Start development server
	@echo "🚀 Starting development server..."
	@docker compose up -d
	@echo "✅ Server running at http://localhost:8080 & https://localhost:8080"

prod: ## 🔐 Start production server (requires SSL setup)
	@echo "🔐 Starting production server..."
	@docker compose up -d
	@echo "✅ Production server running!"

restart: ## 🔄 Restart all services
	@echo "🔄 Restarting services..."
	@docker compose restart
	@echo "✅ Services restarted!"

stop: ## ⏹️ Stop all services
	@echo "⏹️ Stopping all services..."
	@docker compose down
	@echo "✅ All services stopped!"

# === BUILD COMMANDS ===
build: ## 🏗️ Build Docker images
	@echo "🏗️ Building Docker images..."
	@docker compose build --no-cache
	@echo "✅ Build complete!"

rebuild: ## 🔄 Rebuild and restart services
	@echo "🔄 Rebuilding and restarting services..."
	@docker compose down
	@docker compose up --build -d
	@echo "✅ Rebuild complete!"

# === DEVELOPMENT COMMANDS ===
logs: ## 📋 Show all service logs
	@echo "📋 Showing service logs..."
	@docker compose logs -f

app-logs: ## 📋 Show only application logs
	@echo "📋 Showing application logs..."
	@docker compose logs -f app

db-logs: ## 📋 Show database logs
	@echo "📋 Showing database logs..."
	@docker compose logs -f postgres

nginx-logs: ## 📋 Show nginx logs
	@echo "📋 Showing nginx logs..."
	@docker compose logs -f nginx

status: ## 📊 Show service status
	@echo "📊 Service Status:"
	@echo "=================="
	@docker compose ps

health: ## 🏥 Check API health
	@echo "🏥 Checking API health..."
	@curl -s http://localhost:8080/health || echo "❌ API not responding"

# === AI INTEGRATION COMMANDS ===
ai-test: ## 🤖 Test AI integration endpoints
	@echo "🤖 Testing AI integration..."
	@echo "Checking vistara-ai service..."
	@curl -s http://localhost:5000/api/v1/health > /dev/null && echo "✅ vistara-ai service is running" || echo "❌ vistara-ai service not available"
	@echo ""
	@echo "To test AI endpoints, first get a JWT token:"
	@echo "TOKEN=\$$(curl -s -X POST http://localhost:8080/api/auth/login -H 'Content-Type: application/json' -d '{\"email\":\"testuser1@vistara.com\",\"password\":\"password123\"}' | grep -o '\"token\":\"[^\"]*\"' | cut -d'\"' -f4)"
	@echo ""
	@echo "Then test AI endpoints:"
	@echo "curl -X POST http://localhost:8080/api/v1/user/nusalingo -H 'Authorization: Bearer \$$TOKEN' -H 'Content-Type: application/json' -d '{\"from_language\":\"English\",\"to_language\":\"Banjar\",\"text\":\"Hello world\"}'"

# === DATABASE COMMANDS ===
db-connect: ## 🗄️ Connect to database
	@echo "🗄️ Connecting to database..."
	@docker compose exec postgres psql -U vistara_user -d vistara_db

db-migrate: ## 📊 Run database migrations
	@echo "📊 Running database migrations..."
	@docker compose exec app go run cmd/migrate/main.go

db-reset: ## 🔄 Reset database (drop and recreate)
	@echo "🔄 Resetting database..."
	@docker compose exec postgres dropdb -U vistara_user vistara_db --if-exists
	@docker compose exec postgres createdb -U vistara_user vistara_db
	@echo "✅ Database reset complete!"

# === SSL COMMANDS ===
setup-ssl: ## 🔐 Manual SSL setup guide for production
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
	@echo "   6. sudo chown \$$USER:\$$USER ./nginx/ssl/einrafh.com.*"
	@echo "   7. make prod  # Restart with SSL"

setup-dev-ssl: ## 🔧 Setup self-signed SSL for development
	@echo "🔧 Setting up development SSL certificates..."
	@mkdir -p nginx/ssl
	@openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout nginx/ssl/einrafh.com.key \
		-out nginx/ssl/einrafh.com.crt \
		-subj "/C=ID/ST=South Kalimantan/L=Banjarmasin/O=Vistara/CN=einrafh.com"
	@echo "✅ Development SSL certificates created!"

# === TESTING COMMANDS ===
test: ## 🧪 Run tests
	@echo "🧪 Running tests..."
	@docker compose exec app go test ./...

test-coverage: ## 📊 Run tests with coverage
	@echo "📊 Running tests with coverage..."
	@docker compose exec app go test -coverprofile=coverage.out ./...
	@docker compose exec app go tool cover -html=coverage.out -o coverage.html

# === CLEANUP COMMANDS ===
clean: ## 🧹 Clean Docker artifacts
	@echo "🧹 Cleaning Docker artifacts..."
	@docker compose down
	@docker system prune -f
	@echo "✅ Cleanup complete!"

clean-all: ## 🧹 Clean everything (images, volumes, networks)
	@echo "🧹 Cleaning everything..."
	@docker compose down --rmi all --volumes --remove-orphans
	@docker system prune -af --volumes
	@echo "✅ Deep cleanup complete!"

# === DEPLOYMENT COMMANDS ===
deploy: ## 🚀 Complete deployment for production VM
	@echo "🚀 Deploying to production VM..."
	@echo "================================"
	@echo "📋 Checking prerequisites..."
	@which docker > /dev/null || (echo "❌ Docker not found. Install: curl -fsSL https://get.docker.com | sh" && exit 1)
	@docker compose version > /dev/null || (echo "❌ Docker Compose not found" && exit 1)
	@echo "✅ Docker ready"
	@echo ""
	@echo "📂 Setting up environment..."
	@cp .env.example .env || echo "⚠️ Using existing .env"
	@echo "✅ Environment configured"
	@echo ""
	@echo "🏗️ Building and starting services..."
	@docker compose up --build -d
	@echo ""
	@echo "🎉 Deployment complete!"
	@echo "📍 Backend: http://[your-vm-ip]:8080"
	@echo "📍 Frontend proxy: http://[your-vm-ip]:80"
	@echo ""
	@echo "🔐 For HTTPS setup, run: make setup-ssl"

# === UTILITY COMMANDS ===
nginx-reload: ## 🔄 Reload nginx configuration
	@echo "🔄 Reloading nginx configuration..."
	@docker compose exec nginx nginx -s reload
	@echo "✅ Nginx configuration reloaded!"

shell: ## 🐚 Open shell in app container
	@echo "🐚 Opening shell in app container..."
	@docker compose exec app sh

go-mod: ## 📦 Update Go modules
	@echo "📦 Updating Go modules..."
	@docker compose exec app go mod tidy
	@echo "✅ Go modules updated!"

# === TESTING COMMANDS ===
test-auth: ## 🔐 Test authentication endpoints
	@echo "🔐 Testing authentication..."
	@echo "📝 User registration:"
	@curl -X POST http://localhost:8080/api/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@vistara.com","password":"password123","name":"Test User"}' | jq . 2>/dev/null || echo "Response received"
	@echo ""
	@echo "🔑 User login:"
	@curl -X POST http://localhost:8080/api/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"testuser1@vistara.com","password":"password123"}' | jq . 2>/dev/null || echo "Response received"

test-ai: ## 🤖 Test AI integration endpoints
	@echo "🤖 Testing AI integration..."
	@echo "🗺️ Testing Smart Planner:"
	@curl -X POST http://localhost:8080/api/v1/user/smart-planner \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer test-token" \
		-d '{"destination":"Yogyakarta","start_date":"2025-08-01T00:00:00Z","end_date":"2025-08-05T00:00:00Z","budget":5000000,"travel_style":"romantic_couple","activity_preferences":["beach","culture","culinary"],"activity_intensity":"balanced"}' | jq . 2>/dev/null || echo "Response received"
	@echo ""
	@echo "🗣️ Testing Nusalingo Translation:"
	@curl -X POST http://localhost:8080/api/v1/user/nusalingo \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer test-token" \
		-d '{"from_language":"English","to_language":"Banjar","text":"Hello, how are you today?"}' | jq . 2>/dev/null || echo "Response received"
	@echo ""
	@echo "📚 Testing Historical Story:"
	@curl -X POST http://localhost:8080/api/v1/user/historical-story \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer test-token" \
		-d '{"location":"Borobudur Temple"}' | jq . 2>/dev/null || echo "Response received"

test-local: ## 🏪 Test local business endpoints
	@echo "🏪 Testing local business endpoints..."
	@echo "📋 Local businesses list:"
	@curl -X GET http://localhost:8080/api/locals \
		-H "Authorization: Bearer test-token" | jq . 2>/dev/null || echo "Response received"
	@echo ""
	@echo "🗺️ Tourist attractions:"
	@curl -X GET http://localhost:8080/api/tourist-attractions \
		-H "Authorization: Bearer test-token" | jq . 2>/dev/null || echo "Response received"

test-service: ## 🔗 Test service endpoints (for vistara-ai)
	@echo "🔗 Testing service endpoints for vistara-ai..."
	@echo "📍 Local businesses (service):"
	@curl -X GET http://localhost:8080/api/service/locals \
		-H "X-Service: vistara-ai" \
		-H "X-API-Key: vistara-ai-service-key" | jq . 2>/dev/null || echo "Response received"
	@echo ""
	@echo "📍 Tourist attractions (service):"
	@curl -X GET http://localhost:8080/api/service/tourist-attractions \
		-H "X-Service: vistara-ai" \
		-H "X-API-Key: vistara-ai-service-key" | jq . 2>/dev/null || echo "Response received"

test-notification: ## 🔔 Test AI notification endpoint
	@echo "🔔 Testing AI notification endpoint..."
	@curl -X POST http://localhost:8080/api/service/ai/notify \
		-H "X-Service: vistara-ai" \
		-H "X-API-Key: vistara-ai-service-key" \
		-H "Content-Type: application/json" \
		-d '{"event":"plan_generated","user_id":"test-user","data":{"destination":"Bali"},"timestamp":"2025-07-25T10:00:00Z"}' | jq . 2>/dev/null || echo "Response received"

test-all: test-auth test-ai test-local test-service test-notification ## 🧪 Run all endpoint tests

health: ## ❤️ Check API health
	@echo "❤️ Checking API health..."
	@curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health || echo "❌ API not responding"

info: ## ℹ️ Show system information
	@echo "ℹ️ Vistara Backend System Information"
	@echo "===================================="
	@echo "🐳 Docker:"
	@docker --version
	@echo "🏗️ Docker Compose:"
	@docker compose version
	@echo "📊 Container Status:"
	@docker compose ps
	@echo "🌐 Network Status:"
	@docker network ls | grep vistara || echo "No vistara networks found"
	@echo "💾 Volume Status:"
	@docker volume ls | grep vistara || echo "No vistara volumes found"
	@echo ""
	@echo "🤖 AI Integration:"
	@curl -s http://localhost:5000/api/v1/health > /dev/null && echo "✅ vistara-ai service: Running" || echo "❌ vistara-ai service: Not available"
	@echo ""
	@echo "🔗 Available Endpoints:"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:80"
	@echo "HTTPS: https://localhost:443"

# === QUICK START ===
quick-start: ## ⚡ Quick start for new developers
	@echo "⚡ Vistara Backend Quick Start"
	@echo "=============================="
	@echo "🚀 Setting up everything for you..."
	@make setup
	@echo ""
	@echo "🎉 Quick start complete!"
	@echo "✅ Backend running at: http://localhost:8080"
	@echo "✅ Test user: testuser1@vistara.com / password123"
	@echo "🤖 AI features available if vistara-ai service is running"

# Include environment variables if .env exists
ifneq (,$(wildcard .env))
    include .env
endif
