#!/bin/bash

# Quick Deployment Script for Vistara Backend
# Run this script on your VPS/Cloud server for automated deployment

set -e

echo "🚀 Vistara Backend - Quick Deployment"
echo "====================================="

# Check if we're in the right directory
if [ ! -f "docker-compose.yaml" ]; then
    echo "❌ docker-compose.yaml not found. Please run this script from the vistara-be directory."
    exit 1
fi

# Step 1: Check system requirements
echo "🔧 Step 1: Checking system requirements..."

# Check if Docker is installed and running
if command -v docker &> /dev/null; then
    if docker info &> /dev/null; then
        echo "✅ Docker is installed and running"
    else
        echo "⚠️  Docker is installed but not running. Please start Docker daemon."
        echo "   sudo systemctl start docker"
        exit 1
    fi
else
    echo "❌ Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo "✅ Docker installed. Please logout and login again, then re-run this script."
    exit 0
fi

# Check if Docker Compose is available
if docker compose version &> /dev/null; then
    echo "✅ Docker Compose is available"
else
    echo "❌ Docker Compose is not available. Installing..."
    sudo apt update
    sudo apt install docker-compose-plugin -y
    echo "✅ Docker Compose installed"
fi

# Step 2: Setup environment
echo "⚙️  Step 2: Setting up environment..."
if [ ! -f ".env" ]; then
    echo "📄 Copying environment template..."
    cp .env.production .env
    echo "⚠️  Please edit .env file with your actual credentials before continuing."
    echo "   nano .env"
    echo ""
    read -p "Press Enter after updating .env file..." -r
fi

# Step 3: Test with development setup
echo "🧪 Step 3: Testing with development setup (HTTP)..."
make nginx-dev

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 30

# Check if services are running
if curl -f http://localhost/health > /dev/null 2>&1; then
    echo "✅ Services are running successfully!"
else
    echo "❌ Services failed to start. Checking logs..."
    make logs-all
    exit 1
fi

# Step 4: Setup SSL
echo "🔐 Step 4: Setting up SSL certificates..."
echo "⚠️  Make sure your domain einrafh.com points to this VM before continuing."
echo "   Check DNS: dig einrafh.com"
echo ""
read -p "Continue with SSL setup? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    make setup-ssl
    
    # Step 5: Start production
    echo "🎉 Step 5: Starting production setup..."
    make nginx-prod
    
    echo ""
    echo "🎊 Deployment completed successfully!"
    echo "✅ Your API is now available at:"
    echo "   - https://einrafh.com/health"
    echo "   - https://einrafh.com/api/"
    echo "   - https://einrafh.com/ai/"
else
    echo "ℹ️  SSL setup skipped. You can run 'make setup-ssl' later."
    echo "✅ Development setup completed!"
    echo "   - http://localhost/health"
    echo "   - http://einrafh.com/health (after DNS propagation)"
fi

echo ""
echo "📋 Useful commands:"
echo "   make logs-all      # View all logs"
echo "   make nginx-logs    # View nginx logs"
echo "   make nginx-reload  # Reload nginx config"
echo "   make setup-ssl     # Setup SSL certificates"
