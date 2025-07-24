#!/bin/bash

# Vistara Backend Setup Script
# This script sets up the complete development environment and populates test data

set -e  # Exit on any error

echo "🚀 Starting Vistara Backend Setup..."

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

# Step 1: Clean up existing containers
print_status "Cleaning up existing containers..."
docker compose down --remove-orphans || true

# Step 2: Build and start services
print_status "Building and starting services..."
docker compose up --build -d

# Step 3: Wait for services to be ready
print_status "Waiting for services to be ready..."
sleep 10

# Check if services are running
if ! docker compose ps | grep -q "Up"; then
    print_error "Services failed to start properly"
    docker compose logs
    exit 1
fi

print_success "Services are running successfully"

# Step 4: Wait for database to be ready
print_status "Waiting for database to be ready..."
until docker compose exec postgres pg_isready -U postgres; do
    print_status "Waiting for PostgreSQL..."
    sleep 2
done

print_success "Database is ready"

# Step 5: Check API health
print_status "Checking API health..."
for i in {1..30}; do
    if curl -s http://localhost:8080/health > /dev/null; then
        print_success "API is healthy"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "API health check failed"
        exit 1
    fi
    sleep 1
done

# Step 6: Create test users and get JWT tokens
print_status "Creating test users..."

# Register test user 1
USER1_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User 1",
    "email": "testuser1@vistara.com",
    "password": "password123",
    "confirm_password": "password123"
  }' || echo '{"error": "failed"}')

if echo "$USER1_RESPONSE" | grep -q "error"; then
    print_warning "User 1 might already exist, trying to login..."
fi

# Login as test user 1
USER1_LOGIN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser1@vistara.com",
    "password": "password123"
  }')

USER1_TOKEN=$(echo "$USER1_LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$USER1_TOKEN" ]; then
    print_error "Failed to get user 1 token"
    echo "Response: $USER1_LOGIN"
    exit 1
fi

print_success "Test user 1 created and authenticated"

# Register test user 2
USER2_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User 2",
    "email": "testuser2@vistara.com",
    "password": "password123",
    "confirm_password": "password123"
  }' || echo '{"error": "failed"}')

# Login as test user 2
USER2_LOGIN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser2@vistara.com",
    "password": "password123"
  }')

USER2_TOKEN=$(echo "$USER2_LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

print_success "Test user 2 created and authenticated"

# Step 7: Create test local businesses
print_status "Creating test local businesses..."

BUSINESS1=$(curl -s -X POST http://localhost:8080/api/locals \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Warung Sate Maranggi Purwakarta",
    "description": "Warung sate khas Purwakarta dengan bumbu rahasia turun temurun yang sudah terkenal sejak puluhan tahun",
    "address": "Jl. Raya Purwakarta No. 123",
    "city": "Purwakarta",
    "province": "Jawa Barat",
    "longitude": "107.4341",
    "latitude": "-6.5569",
    "label": "Kuliner",
    "opened_time": "08:00-22:00",
    "photo_url": "https://example.com/sate-maranggi.jpg",
    "is_business": true
  }')

BUSINESS1_ID=$(echo "$BUSINESS1" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

BUSINESS2=$(curl -s -X POST http://localhost:8080/api/locals \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gudeg Yu Djum Yogyakarta",
    "description": "Gudeg legendaris Yogyakarta dengan cita rasa autentik dan pelayanan ramah keluarga",
    "address": "Jl. Malioboro No. 456",
    "city": "Yogyakarta",
    "province": "DI Yogyakarta",
    "longitude": "110.3650",
    "latitude": "-7.7956",
    "label": "Kuliner",
    "opened_time": "06:00-23:00",
    "photo_url": "https://example.com/gudeg.jpg",
    "is_business": true
  }')

BUSINESS2_ID=$(echo "$BUSINESS2" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

BUSINESS3=$(curl -s -X POST http://localhost:8080/api/locals \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Toko Batik Lasem Heritage",
    "description": "Toko batik tradisional dengan koleksi batik Lasem asli dan berkualitas tinggi",
    "address": "Jl. Veteran No. 789",
    "city": "Rembang",
    "province": "Jawa Tengah",
    "longitude": "111.3433",
    "latitude": "-6.7089",
    "label": "UMKM",
    "opened_time": "09:00-17:00",
    "photo_url": "https://example.com/batik-lasem.jpg",
    "is_business": true
  }')

BUSINESS3_ID=$(echo "$BUSINESS3" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

print_success "Created 3 test local businesses"

# Step 8: Create test tourist attractions
print_status "Creating test tourist attractions..."

ATTRACTION1=$(curl -s -X POST http://localhost:8080/api/tourist-attractions \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Borobudur Temple Tour",
    "description": "Explore the magnificent Borobudur Temple, a UNESCO World Heritage Site with professional tour guide service",
    "address": "Borobudur, Magelang",
    "city": "Magelang",
    "province": "Jawa Tengah",
    "longitude": 110.2038,
    "latitude": -7.6079,
    "photo_url": "https://example.com/borobudur.jpg",
    "tour_guide_price": 150000,
    "tour_guide_count": 5,
    "tour_guide_discount_percentage": 10.0,
    "price": 200000,
    "discount_percentage": 15.0
  }')

ATTRACTION1_ID=$(echo "$ATTRACTION1" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

ATTRACTION2=$(curl -s -X POST http://localhost:8080/api/tourist-attractions \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bromo Sunrise Tour",
    "description": "Witness the spectacular sunrise over Mount Bromo with experienced local guides and comfortable transportation",
    "address": "Bromo Tengger Semeru National Park",
    "city": "Probolinggo",
    "province": "Jawa Timur", 
    "longitude": 112.9533,
    "latitude": -7.9425,
    "photo_url": "https://example.com/bromo.jpg",
    "tour_guide_price": 300000,
    "tour_guide_count": 3,
    "tour_guide_discount_percentage": 5.0,
    "price": 450000,
    "discount_percentage": 20.0
  }')

ATTRACTION2_ID=$(echo "$ATTRACTION2" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

ATTRACTION3=$(curl -s -X POST http://localhost:8080/api/tourist-attractions \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kawah Putih Ciwidey Tour",
    "description": "Visit the stunning white crater lake in Ciwidey with complete tour package and photography service",
    "address": "Kawah Putih, Ciwidey",
    "city": "Bandung",
    "province": "Jawa Barat",
    "longitude": 107.4019,
    "latitude": -7.1661,
    "photo_url": "https://example.com/kawah-putih.jpg",
    "tour_guide_price": 120000,
    "tour_guide_count": 4,
    "tour_guide_discount_percentage": 12.0,
    "price": 180000,
    "discount_percentage": 8.0
  }')

ATTRACTION3_ID=$(echo "$ATTRACTION3" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

print_success "Created 3 test tourist attractions"

# Step 9: Display setup summary
print_success "🎉 Setup completed successfully!"
echo ""
echo "📊 SETUP SUMMARY:"
echo "=================="
echo "✅ Services: Running on Docker"
echo "✅ Database: PostgreSQL ready"
echo "✅ API: Healthy at http://localhost:8080"
echo "✅ Test Users: 2 users created"
echo "✅ Local Businesses: 3 businesses created"
echo "✅ Tourist Attractions: 3 attractions created"
echo ""
echo "🔑 TEST CREDENTIALS:"
echo "==================="
echo "User 1: testuser1@vistara.com / password123"
echo "User 2: testuser2@vistara.com / password123"
echo ""
echo "🚀 ENDPOINT TESTING:"
echo "===================="
echo "Base URL: http://localhost:8080"
echo ""
echo "📁 Available Endpoints:"
echo "• POST /api/auth/register - Register new user"
echo "• POST /api/auth/login - Login user"
echo "• GET /health - Health check"
echo ""
echo "🏢 Local Business CRUD (Auth Required):"
echo "• GET /api/locals - Get all businesses"
echo "• GET /api/locals/{id} - Get specific business"
echo "• POST /api/locals - Create business"
echo "• PUT /api/locals/{id} - Update business"
echo "• DELETE /api/locals/{id} - Delete business"
echo ""
echo "🏛️ Tourist Attraction CRUD (Auth Required):"
echo "• GET /api/tourist-attractions - Get all attractions"
echo "• GET /api/tourist-attractions/{id} - Get specific attraction"
echo "• POST /api/tourist-attractions - Create attraction"
echo "• PUT /api/tourist-attractions/{id} - Update attraction"
echo "• DELETE /api/tourist-attractions/{id} - Delete attraction"
echo "• GET /api/tourist-attractions/{id}/availability - Get booking availability"
echo "• POST /api/tourist-attractions/{id}/book - Create booking"
echo ""
echo "💡 EXAMPLE TEST COMMANDS:"
echo "========================="
echo "# Get JWT Token:"
echo 'TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '\''{"email":"testuser1@vistara.com","password":"password123"}'\'' \'
echo '  | grep -o '\'"token":"[^"]*"'\'' | cut -d'\'"'\'' -f4)'
echo ""
echo "# Test GET all businesses:"
echo 'curl -X GET http://localhost:8080/api/locals \'
echo '  -H "Authorization: Bearer $TOKEN"'
echo ""
echo "# Test GET all attractions:"
echo 'curl -X GET http://localhost:8080/api/tourist-attractions \'
echo '  -H "Authorization: Bearer $TOKEN"'
echo ""
if [ -n "$BUSINESS1_ID" ]; then
    echo "📋 CREATED RESOURCE IDs:"
    echo "======================="
    echo "Business 1 ID: $BUSINESS1_ID"
    echo "Business 2 ID: $BUSINESS2_ID"
    echo "Business 3 ID: $BUSINESS3_ID"
    echo "Attraction 1 ID: $ATTRACTION1_ID"
    echo "Attraction 2 ID: $ATTRACTION2_ID"
    echo "Attraction 3 ID: $ATTRACTION3_ID"
fi
echo ""
print_success "Ready for testing! 🚀"
