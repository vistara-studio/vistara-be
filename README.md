# 🏛️ Vistara Backend

Tourism platform backend API for Indonesian local business and tourist attractions management.

## ✨ Features

- 🔐 **Authentication** - JWT-based user registration & login
- 🏢 **Local Business CRUD** - Manage local businesses (restaurants, shops, etc.)
- 🏛️ **Tourist Attractions CRUD** - Manage tourist destinations with booking
- 🤖 **AI Integration** - Smart travel planning via vistara-ai service
- 💳 **Payment Integration** - Midtrans payment gateway
- 📱 **RESTful API** - Complete REST endpoints with proper error handling

## 🛠️ Tech Stack

- **Go 1.23** + **Fiber v2** (Web Framework)
- **PostgreSQL** + **GORM** (Database)
- **JWT** (Authentication)
- **Docker** (Containerization)
- **Midtrans** (Payment)

## 🚀 Quick Start

### One Command Setup
```bash
make setup
```

This will:
- ✅ Start PostgreSQL & API containers
- ✅ Run database migrations  
- ✅ Create test users & sample data
- ✅ Ready for testing at `http://localhost:8080`

### Manual Setup
```bash
# 1. Clone repository
git clone https://github.com/vistara-studio/vistara-be.git
cd vistara-be

# 2. Copy environment file
cp .env.example .env

# 3. Start services
docker compose up -d

# 4. Check health
curl http://localhost:8080/health
```

## 📋 Available Commands

```bash
make setup      # 🚀 Complete setup with test data
make up         # 🔼 Start services  
make down       # 🔽 Stop services
make logs       # 📋 View logs
make health     # ❤️ Check API health
make help       # ℹ️ Show all commands
```

## 🌐 API Endpoints

### 🔐 Authentication
```http
POST /api/auth/register  # User registration
POST /api/auth/login     # User login
```

### 🏢 Local Business
```http
GET    /api/locals           # Get all businesses
POST   /api/locals           # Create business
GET    /api/locals/{id}      # Get specific business
PUT    /api/locals/{id}      # Update business
DELETE /api/locals/{id}      # Delete business
```

### 🏛️ Tourist Attractions  
```http
GET    /api/tourist-attractions           # Get all attractions
POST   /api/tourist-attractions           # Create attraction
GET    /api/tourist-attractions/{id}      # Get specific attraction
PUT    /api/tourist-attractions/{id}      # Update attraction
DELETE /api/tourist-attractions/{id}      # Delete attraction
POST   /api/tourist-attractions/{id}/book # Create booking
```

### ❤️ System
```http
GET /health  # Health check
```

## 🧪 Testing

Get JWT token:
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser1@vistara.com","password":"password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
```

Test endpoints:
```bash
# Get all local businesses
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/locals

# Get all tourist attractions
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/tourist-attractions
```

## 📁 Project Structure

```
├── cmd/api/              # Application entry point
├── internal/domain/      # Business domains
│   ├── user/            # User & authentication  
│   ├── session/         # JWT session management
│   └── local/           # Local business & attractions
├── internal/infra/       # Infrastructure layer
├── db/migrations/        # Database migrations
└── scripts/              # Automation scripts
```

## 🔧 Environment Variables

Key variables in `.env`:
```bash
APP_PORT=8080
POSTGRES_DB=vistara_db
POSTGRES_USERNAME=vistara_user  
POSTGRES_PASSWORD=vistara_password
JWT_SECRET=your-jwt-secret
MIDTRANS_SERVER_KEY=your-midtrans-key
VISTARA_AI_URL=http://localhost:5000  # AI service integration
```

## 🤖 AI Integration

Vistara-BE integrates with vistara-ai service for smart travel planning:

### AI Endpoints
```bash
# Generate smart travel plan
curl -X POST http://localhost:8080/api/ai/smart-plan \
  -H "Content-Type: application/json" \
  -d '{
    "destination": "Bali",
    "start_date": "2025-08-01T00:00:00Z",
    "end_date": "2025-08-05T00:00:00Z",
    "budget": 5000000,
    "travel_style": "romantic_couple",
    "activity_preferences": ["beach", "culture", "culinary"]
  }'
```

### Service-to-Service Endpoints (for vistara-ai)
```bash
# Get local businesses (AI service access)
curl -H "X-Service: vistara-ai" \
     http://localhost:8080/api/service/locals

# Get tourist attractions (AI service access)  
curl -H "X-Service: vistara-ai" \
     http://localhost:8080/api/service/tourist-attractions
```

### Testing AI Integration
```bash
# Test all AI integration features
make test-ai-integration
make test-service-endpoints
make test-notification
```
## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/tourism-local-business-management-system`
3. Commit changes: `git commit -m 'Add tourism management features'`
4. Push to branch: `git push origin feature/tourism-local-business-management-system`
5. Create Pull Request

## 📄 License

Copyright (c) 2025 Muhammad Rafly Ash Shiddiqi
