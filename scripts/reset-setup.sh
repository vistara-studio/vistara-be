#!/bin/bash

# Vistara Backend Reset Script  
# Completely resets the development environment including AI integration

set -e  # Exit on any error

echo "🔄 Vistara Backend Environment Reset"
echo "===================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Confirmation prompt
echo ""
echo "⚠️  WARNING: This will completely reset your environment!"
echo ""
echo "This will:"
echo "• Stop and remove all containers (Backend, Database, Nginx)"
echo "• Delete all database data and volumes"
echo "• Remove Docker images (optional)"
echo "• Clean build artifacts and caches"
echo "• Remove temporary files and SSL certificates"
echo "• Clear AI service integration logs"
echo ""
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Reset cancelled."
    exit 0
fi

# Step 1: Stop all containers and processes
print_status "Stopping all containers and processes..."
docker compose down --remove-orphans 2>/dev/null || true

# Kill any remaining processes on port 8080 (if running outside Docker)
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    print_status "Stopping process on port 8080..."
    lsof -Pi :8080 -sTCP:LISTEN -t | xargs kill -TERM 2>/dev/null || true
fi

print_success "All processes stopped"

# Step 2: Remove containers and volumes
print_status "Removing containers and volumes..."
docker compose down --volumes --remove-orphans 2>/dev/null || true

# Step 3: Optional image removal
echo ""
read -p "Remove Docker images as well? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Removing Docker images..."
    docker compose down --rmi all --volumes --remove-orphans 2>/dev/null || true
    print_success "Docker images removed"
else
    print_status "Keeping Docker images for faster rebuilds"
fi

# Step 4: Clean up orphaned Docker resources
print_status "Cleaning up orphaned Docker resources..."
docker container prune -f 2>/dev/null || true
docker volume prune -f 2>/dev/null || true
docker network prune -f 2>/dev/null || true

# Step 5: Clean Go build artifacts and caches
print_status "Cleaning Go build artifacts and caches..."
rm -f vistara-backend vistara-be 2>/dev/null || true
go clean -cache -modcache -testcache 2>/dev/null || true

# Step 6: Clean temporary files and logs
print_status "Cleaning temporary files and logs..."
rm -rf tmp/ temp/ *.log .docker/ 2>/dev/null || true
rm -rf .ai-service-logs/ 2>/dev/null || true

# Step 7: Clean SSL certificates (keep examples)
print_status "Cleaning SSL certificates..."
find nginx/ssl/ -name "*.crt" -not -name "*.example" -delete 2>/dev/null || true
find nginx/ssl/ -name "*.key" -not -name "*.example" -delete 2>/dev/null || true

# Step 8: Clean any remaining application state
print_status "Cleaning application state..."
rm -rf .data/ data/ postgres-data/ 2>/dev/null || true

# Step 9: Final system cleanup
print_status "Performing final cleanup..."
docker system prune -f 2>/dev/null || true

print_success "🎉 Environment reset completed!"
echo ""
echo "📊 RESET SUMMARY"
echo "================"
echo "✅ Containers: Stopped and removed (Backend + Database + Nginx)"
echo "✅ Volumes: Removed (all data cleared)"
echo "✅ Networks: Cleaned up"
echo "✅ Build artifacts: Removed"
echo "✅ Temporary files: Cleaned up"
echo "✅ SSL certificates: Reset (examples preserved)"
echo "✅ AI service logs: Cleared"
echo "✅ Application state: Reset"
echo ""
echo "🚀 NEXT STEPS"
echo "============="
echo "To start fresh:"
echo "• Run 'make setup' for complete setup with test data and AI integration"
echo "• Run 'make dev-setup' for basic development setup"
echo "• Run 'docker compose up -d' for manual startup"
echo "• Ensure vistara-ai service is running at localhost:5000 for AI features"
echo ""
print_success "Ready for fresh setup! 🔄"
