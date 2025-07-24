#!/bin/bash

# Vistara Backend Reset Setup Script
# This script completely resets the development environment

set -e  # Exit on any error

echo "🔄 Starting Vistara Backend Reset..."

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

# Step 1: Stop all running containers
print_status "Stopping all running containers..."
docker compose down --remove-orphans || true

# Step 2: Remove containers and volumes
print_status "Removing containers and volumes..."
docker compose down --volumes --remove-orphans || true

# Step 3: Remove images (optional - ask user)
read -p "Do you want to remove Docker images as well? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Removing Docker images..."
    docker compose down --rmi all --volumes --remove-orphans || true
    print_success "Docker images removed"
else
    print_status "Keeping Docker images"
fi

# Step 4: Clean up any orphaned containers
print_status "Cleaning up orphaned containers..."
docker container prune -f || true

# Step 5: Clean up any orphaned volumes
print_status "Cleaning up orphaned volumes..."
docker volume prune -f || true

# Step 6: Clean up any orphaned networks
print_status "Cleaning up orphaned networks..."
docker network prune -f || true

# Step 7: Remove any build artifacts
print_status "Cleaning up build artifacts..."
rm -f vistara-backend vistara-be || true
go clean -cache -modcache -testcache || true

# Step 8: Clean up temporary files
print_status "Cleaning up temporary files..."
rm -rf tmp/ temp/ *.log || true

print_success "🎉 Reset completed successfully!"
echo ""
echo "📊 RESET SUMMARY:"
echo "================="
echo "✅ Containers: Stopped and removed"
echo "✅ Volumes: Removed (database data cleared)"
echo "✅ Networks: Cleaned up"
echo "✅ Build artifacts: Removed"
echo "✅ Temporary files: Cleaned up"
echo ""
echo "🚀 NEXT STEPS:"
echo "=============="
echo "Run 'make setup' to start fresh setup"
echo "Or run 'docker compose up -d' for basic startup"
echo ""
print_success "Environment reset complete! 🔄"
