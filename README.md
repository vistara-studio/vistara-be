# Vistara Backend

<div align="center">
  <img src="assets/vistara-mockup.png" alt="Vistara AI Mockup" width="800"/>
</div>

<br/>

Backend API for Vistara - A culture and tourism-based digital platform integrating education, ticketing, navigation, and local language translation to preserve and promote Indonesian heritage.

## 🚀 Quick Start

```bash
# Complete setup with test data
make setup

# Development mode
make dev

# Test all endpoints
make test-all
```

## 📋 Main Features

### 🔐 Authentication System
JWT-based authentication with user management and premium features.

**Endpoints:**
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### 🏪 Local Business Management
Comprehensive local business and tourist attraction management.

**Endpoints:**
- `GET /api/locals` - List local businesses
- `GET /api/tourist-attractions` - List tourist attractions
- `POST /api/locals` - Create local business (premium)

### 🤖 AI Integration
Seamless integration with vistara-ai service for intelligent features.

**User Endpoints:** `POST /api/v1/user/*`
```json
// Smart Planner
{
  "destination": "Yogyakarta",
  "start_date": "2025-06-10T00:00:00Z",
  "end_date": "2025-06-12T00:00:00Z",
  "budget": 300000,
  "activity_preferences": ["Nature Exploration", "History & culture"],
  "travel_style": "solo_traveler",
  "activity_intensity": "balanced"
}

// Nusalingo Translation
{
  "from_language": "English",
  "to_language": "Banjar",
  "text": "Hello, how are you today?"
}

// Historical Story
{
  "location": "Borobudur Temple"
}
```

**Service Endpoints:** `POST /api/v1/service/*` (for vistara-ai communication)

### 📊 Data Management
Complete CRUD operations for tourism data with PostgreSQL.

## ⚙️ Configuration

### Environment Variables
```bash
# Server
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=vistara_user
DB_PASSWORD=vistara_password
DB_NAME=vistara_db

# AI Service Integration
VISTARA_AI_URL=http://localhost:5000
VISTARA_AI_KEY=vistara-be-service-key

# JWT
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRY=24h

# Payment (Midtrans)
MIDTRANS_SERVER_KEY=your_midtrans_server_key
MIDTRANS_CLIENT_KEY=your_midtrans_client_key
MIDTRANS_ENVIRONMENT=sandbox
```

### Performance Features
- 🚀 **Fiber Framework** - High-performance HTTP framework
- 🗄️ **PostgreSQL** - Robust relational database
- 🔐 **JWT Authentication** - Secure token-based auth
- 🐳 **Docker Support** - Containerized deployment
- 🔄 **Auto Migration** - Database schema management

## 🛠️ Development

```bash
# Setup development environment
make dev-setup

# Build application
make build

# Run locally
make run

# Watch logs
make logs

# Database operations
make db-logs
make db-shell

# Docker commands
make docker-build
make docker-run
make docker-clean
```

## 🧪 Testing

```bash
# Test individual services
make test-auth         # Authentication
make test-ai          # AI integration
make test-local       # Local business
make test-service     # Service endpoints

# Health check
make health

# Run all tests
make test-all
```

## 📁 Project Structure

```
vistara-be/
├── cmd/api/                 # Application entry point
├── internal/
│   ├── bootstrap/          # App initialization
│   ├── domain/             # Business logic
│   │   ├── ai/            # AI integration
│   │   ├── local/         # Local business
│   │   ├── session/       # Authentication
│   │   └── user/          # User management
│   ├── infra/             # Infrastructure
│   │   ├── ai/            # AI service client
│   │   ├── config/        # Configuration
│   │   ├── db/            # Database connection
│   │   └── http/          # HTTP server
│   └── middleware/        # HTTP middlewares
├── db/migrations/          # Database migrations
├── pkg/                   # Shared utilities
├── scripts/               # Setup & maintenance scripts
└── nginx/                 # Reverse proxy config
```

## 🔧 Maintenance

```bash
# Reset project completely
make reset-setup

# Clean build artifacts
make clean

# Update dependencies
make go-mod

# Restart services
make restart

# Open container shell
make shell
```

## 🐳 Docker Support

```bash
# Build & run with Docker Compose
docker-compose up --build

# Or with make commands
make docker-build
make docker-run

# Production deployment
make deploy
```

## 🌐 API Documentation

All user endpoints require JWT authorization:
```
Authorization: Bearer <your-jwt-token>
```

Service endpoints require service authentication:
```
X-Service: vistara-ai
X-API-Key: <service-api-key>
```

Response format follows consistent JSON structure with proper error handling.

## 🚀 Deployment

### Development
```bash
make setup          # Complete setup
make dev            # Start development
```

### Production VM
```bash
make deploy         # Auto deployment
make setup-ssl      # HTTPS setup guide
```

## 🔗 Integration

Works seamlessly with:
- **vistara-ai** - AI service for smart features
- **Midtrans** - Payment gateway integration
- **PostgreSQL** - Primary database
- **Nginx** - Reverse proxy and SSL

---

## 📄 License

Copyright © 2025 [Muhammad Rafly Ash Shiddiqi](https://github.com/einrafh)

Licensed under MIT License - see [LICENSE](LICENSE) file for complete details.

---

*Made with ❤️ for Indonesian tourism*
