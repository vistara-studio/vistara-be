# 🏛️ Vistara Backend

Indonesian tourism platform backend with AI-powered travel planning.

## ✨ Key Features

- 🔐 **JWT Authentication** - Secure user registration & login
- 🏢 **Local Business** - Manage restaurants, shops, accommodations  
- 🏛️ **Tourist Attractions** - Destination management with bookings
- 🤖 **AI Integration** - Smart travel planning via vistara-ai service
- 💳 **Payment Gateway** - Midtrans integration
- 🔗 **Microservice API** - Service-to-service communication

## 🛠️ Tech Stack

**Backend:** Go 1.23 + Fiber v2 + PostgreSQL + Docker  
**AI Integration:** HTTP client for vistara-ai service  
**Payment:** Midtrans gateway  
**Storage:** Supabase integration

## 🚀 Quick Start

```bash
# One command setup (recommended)
make setup

# Alternative: basic setup
make dev-setup

# Check status
make status
make health
```

**Setup includes:** Docker containers, database migrations, test users, sample data, AI integration testing.

## 📋 Essential Commands

```bash
make setup          # 🚀 Complete setup with test data
make reset-setup    # 🔄 Full environment reset  
make start          # ▶️ Start services
make stop           # ⏹️ Stop services
make logs           # 📋 View logs
make test-all       # 🧪 Test all endpoints
make help           # 📋 Show all commands (30+)
```

## 🌐 API Endpoints

### Core Endpoints
```http
# Authentication
POST /api/auth/register
POST /api/auth/login

# Business Management (requires auth)
CRUD /api/locals/*
CRUD /api/tourist-attractions/*

# AI Integration
POST /api/ai/smart-planner

# Service-to-Service (for vistara-ai)
GET /api/service/locals
GET /api/service/tourist-attractions
POST /api/service/ai/notify
```

## 🧪 Testing

### Quick Test
```bash
# Get auth token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser1@vistara.com","password":"password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Test endpoints
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/locals
```

### Test Commands
```bash
make test-auth      # Authentication endpoints
make test-ai        # AI integration  
make test-service   # Service endpoints
make test-all       # All endpoints
```

## 🤖 AI Integration Example

**Note:** AI endpoints require JWT authentication. First get a token:

```bash
# Get JWT Token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser1@vistara.com","password":"password123"}' \
  | jq -r '.payload.token')

# Use AI Smart Planning (requires authentication)
curl -X POST http://localhost:8080/api/ai/smart-planner \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "destination": "Yogyakarta",
    "start_date": "2025-06-10T00:00:00Z", 
    "end_date": "2025-06-12T00:00:00Z",
    "budget": 3000000,
    "travel_style": "solo_traveler",
    "activity_preferences": ["Nature Exploration", "History & culture", "Culinary"],
    "activity_intensity": "balanced"
  }'
```

**Note:** This endpoint proxies to `http://localhost:5000/api/v1/smart-planner` (vistara-ai service)

## 🔧 Configuration

Key environment variables:
```bash
APP_PORT=8080
POSTGRES_DB=vistara_db
JWT_SECRET=your-jwt-secret
MIDTRANS_SERVER_KEY=your-key
VISTARA_AI_URL=http://localhost:5000
```

### AI Service Requirements
For AI integration to work, ensure `vistara-ai` service is running at:
- **Health Check:** `http://localhost:5000/api/v1/health`
- **Smart Planner:** `http://localhost:5000/api/v1/smart-planner`

## 📁 Structure

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── domain/        # Business logic (user, local, ai)
│   ├── infra/         # Infrastructure (db, http, ai client)
│   ├── middleware/    # HTTP middleware
│   └── bootstrap/     # App initialization
├── scripts/           # Setup automation
└── db/migrations/     # Database migrations
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch: `git checkout -b feature/your-feature`
3. Commit: `git commit -m 'Add feature'`
4. Push: `git push origin feature/your-feature`
5. Create Pull Request

---
**Ready for development!** Run `make setup` to get started. 🚀
