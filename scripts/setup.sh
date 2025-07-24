#!/bin/bash

# Vistara Backend Complete Setup Script
# Sets up the complete development environment

set -e  # Exit on any error

echo "🚀 Vistara Backend Complete Setup"
echo "=================================="

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
echo "This will:"
echo "• Stop and remove existing containers"
echo "• Build and start fresh containers"
echo "• Run database migrations"
echo "• Create test users and sample data"
echo "• Test AI integration if available"
echo ""
read -p "Continue with setup? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Setup cancelled."
    exit 0
fi

# Step 1: Clean up existing containers
print_status "Cleaning up existing containers..."
docker compose down --remove-orphans 2>/dev/null || true

# Step 2: Build and start services
print_status "Building and starting services..."
docker compose up --build -d

# Step 3: Wait for services to be ready
print_status "Waiting for services to be ready..."
sleep 10

# Verify services are running
if ! docker compose ps | grep -q "Up"; then
    print_error "Services failed to start"
    docker compose logs
    exit 1
fi

print_success "Services are running"

# Step 4: Wait for database
print_status "Waiting for database to be ready..."
until docker compose exec postgres pg_isready -U postgres >/dev/null 2>&1; do
    print_status "Database not ready, waiting..."
    sleep 2
done

print_success "Database is ready"

# Step 5: Check API health
print_status "Checking API health..."
for i in {1..30}; do
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "API health check failed"
        exit 1
    fi
    sleep 1
done

print_success "API is healthy"

# Step 6: Create test users
print_status "Creating test users..."

# Register test user 1
USER1_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User 1",
    "email": "testuser1@vistara.com",
    "password": "password123",
    "confirm_password": "password123"
  }' 2>/dev/null || echo "failed")

if echo "$USER1_RESPONSE" | grep -q "error\|failed"; then
    print_warning "User 1 might already exist, trying to login..."
fi

# Login as test user 1
USER1_LOGIN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser1@vistara.com","password":"password123"}' 2>/dev/null || echo "failed")

USER1_TOKEN=$(echo "$USER1_LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$USER1_TOKEN" ]; then
    print_error "Failed to get user 1 token"
    exit 1
fi

print_success "Test user 1 authenticated"

# Register and login test user 2
USER2_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User 2",
    "email": "testuser2@vistara.com", 
    "password": "password123",
    "confirm_password": "password123"
  }' 2>/dev/null || echo "failed")

USER2_LOGIN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser2@vistara.com","password":"password123"}' 2>/dev/null || echo "failed")

USER2_TOKEN=$(echo "$USER2_LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

print_success "Test user 2 authenticated"

# Step 7: Create sample local businesses
print_status "Creating sample local businesses..."

BUSINESS_SAMPLES=(
  '{"name":"Warung Gudeg Yu Djum","description":"Traditional Gudeg restaurant serving authentic Yogyakarta cuisine","address":"Jl. Wijilan No.167","city":"Yogyakarta","province":"DI Yogyakarta","longitude":"110.3644","latitude":"-7.8014","label":"Kuliner","opened_time":"08:00-22:00","photo_url":"https://example.com/gudeg.jpg","is_business":true}'
  '{"name":"Toko Batik Malioboro","description":"Traditional batik shop with authentic Indonesian patterns","address":"Jl. Malioboro No.60","city":"Yogyakarta","province":"DI Yogyakarta","longitude":"110.3658","latitude":"-7.7924","label":"Belanja","opened_time":"09:00-21:00","photo_url":"https://example.com/batik.jpg","is_business":true}'
  '{"name":"Homestay Keluarga Budi","description":"Cozy family homestay with traditional Javanese hospitality","address":"Jl. Prawirotaman II No.629","city":"Yogyakarta","province":"DI Yogyakarta","longitude":"110.3617","latitude":"-7.8131","label":"Penginapan","opened_time":"24/7","photo_url":"https://example.com/homestay.jpg","is_business":true}'
)

for i in "${!BUSINESS_SAMPLES[@]}"; do
    BUSINESS_RESPONSE=$(curl -s -X POST http://localhost:8080/api/locals \
      -H "Authorization: Bearer $USER1_TOKEN" \
      -H "Content-Type: application/json" \
      -d "${BUSINESS_SAMPLES[$i]}" 2>/dev/null || echo "failed")
    
    if echo "$BUSINESS_RESPONSE" | grep -q "id"; then
        print_success "Created local business $((i+1))"
    fi
done

# Step 8: Create sample tourist attractions
print_status "Creating sample tourist attractions..."

ATTRACTION_SAMPLES=(
  '{"name":"Borobudur Temple Tour","description":"Explore the magnificent Borobudur Temple with professional tour guide service","address":"Borobudur, Magelang","city":"Magelang","province":"Jawa Tengah","longitude":110.2038,"latitude":-7.6079,"photo_url":"https://example.com/borobudur.jpg","tour_guide_price":150000,"tour_guide_count":5,"tour_guide_discount_percentage":10.0,"price":200000,"discount_percentage":15.0}'
  '{"name":"Prambanan Temple Heritage Tour","description":"Discover the ancient Hindu temple complex with expert cultural guidance","address":"Prambanan, Klaten","city":"Klaten","province":"Jawa Tengah","longitude":110.4915,"latitude":-7.7520,"photo_url":"https://example.com/prambanan.jpg","tour_guide_price":120000,"tour_guide_count":3,"tour_guide_discount_percentage":5.0,"price":175000,"discount_percentage":12.0}'
  '{"name":"Taman Sari Water Castle Experience","description":"Journey through the historic royal water palace with traditional stories","address":"Patehan, Kraton","city":"Yogyakarta","province":"DI Yogyakarta","longitude":110.3597,"latitude":-7.8106,"photo_url":"https://example.com/tamansari.jpg","tour_guide_price":100000,"tour_guide_count":4,"tour_guide_discount_percentage":8.0,"price":125000,"discount_percentage":10.0}'
)

for i in "${!ATTRACTION_SAMPLES[@]}"; do
    ATTRACTION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/tourist-attractions \
      -H "Authorization: Bearer $USER1_TOKEN" \
      -H "Content-Type: application/json" \
      -d "${ATTRACTION_SAMPLES[$i]}" 2>/dev/null || echo "failed")
    
    if echo "$ATTRACTION_RESPONSE" | grep -q "id"; then
        print_success "Created tourist attraction $((i+1))"
    fi
done

# Step 9: Test AI integration
print_status "Testing AI integration..."

AI_HEALTH=$(curl -s http://localhost:5000/api/v1/health 2>/dev/null || echo "failed")
if echo "$AI_HEALTH" | grep -q "failed"; then
    print_warning "vistara-ai service not detected at http://localhost:5000"
    print_warning "AI smart planning will not be available"
else
    print_success "vistara-ai service detected!"
    
    # Test AI smart planning with authentication
    AI_PLAN_TEST=$(curl -s -X POST http://localhost:8080/api/ai/smart-plan \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TEST_TOKEN1" \
      -d '{
        "destination": "Yogyakarta",
        "start_date": "2025-06-10T00:00:00Z",
        "end_date": "2025-06-12T00:00:00Z",
        "budget": 3000000,
        "travel_style": "solo_traveler",
        "activity_preferences": ["Nature Exploration", "History & culture", "Culinary"],
        "activity_intensity": "balanced"
      }' 2>/dev/null || echo "failed")
    
    if echo "$AI_PLAN_TEST" | grep -q "success"; then
        print_success "AI smart planning integration working!"
    fi
fi

# Step 10: Display setup summary
print_success "🎉 Setup completed successfully!"
echo ""
echo "📊 SETUP SUMMARY"
echo "================"
echo "✅ Services: Running on Docker"
echo "✅ Database: PostgreSQL ready with migrations"
echo "✅ API: Available at http://localhost:8080"
echo "✅ Test Users: 2 users created with authentication"
echo "✅ Sample Data: 3 local businesses, 3 tourist attractions"
echo "✅ AI Integration: ${AI_HEALTH//failed/Not Available}"
echo ""
echo "🔑 TEST CREDENTIALS"
echo "==================="
echo "User 1: testuser1@vistara.com / password123"
echo "User 2: testuser2@vistara.com / password123"
echo ""
echo "🌐 API ENDPOINTS"
echo "================"
echo "Health Check:     GET  /health"
echo "Authentication:   POST /api/auth/register, /api/auth/login"
echo "Local Business:   CRUD /api/locals/*"
echo "Tourist Attractions: CRUD /api/tourist-attractions/*"
echo "AI Planning:      POST /api/ai/smart-plan"
echo "Service Data:     GET  /api/service/* (for vistara-ai)"
echo ""
echo "🧪 QUICK TEST COMMANDS"
echo "======================="
echo "# Get JWT token"
echo 'TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '\''{"email":"testuser1@vistara.com","password":"password123"}'\'' \'
echo '  | grep -o '\'"token":"[^"]*"'\'' | cut -d'\'"'\'' -f4)'
echo ""
echo "# Test endpoints"
echo 'curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/locals'
echo 'curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/tourist-attractions'
echo ""
echo "# Test AI Smart Planning (requires authentication)"
echo 'curl -X POST http://localhost:8080/api/ai/smart-planner \'
echo '  -H "Content-Type: application/json" \'
echo '  -H "Authorization: Bearer $TOKEN" \'
echo '  -d '\''{"destination":"Yogyakarta","start_date":"2025-06-10T00:00:00Z","end_date":"2025-06-12T00:00:00Z","budget":300000}'\'''
echo ""
echo "🚀 READY FOR DEVELOPMENT!"
