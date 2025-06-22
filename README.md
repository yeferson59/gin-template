# 🚀 Professional Gin API Template

## Description

**Gin API Template** is an enterprise-grade template for building robust, scalable REST APIs in Go using the Gin framework. It incorporates industry best practices including JWT authentication, structured logging, rate limiting, comprehensive validation, health checks, graceful shutdown, and production-ready security measures.

## ✨ Key Features

- 🔐 **JWT Authentication** with secure validation
- 📊 **Structured Logging** with Logrus (JSON/Text formats)
- 🛡️ **Rate Limiting** (IP-based with different rules for auth endpoints)
- ✅ **Input Validation** with comprehensive password security
- 🏥 **Health Checks** (Kubernetes-ready liveness/readiness probes)
- 🔒 **Security Headers** (OWASP recommendations)
- 📝 **Consistent API Responses** with standardized error handling
- 🔄 **Graceful Shutdown** with configurable timeouts
- 🐳 **Docker & Kubernetes Ready** with optimized containers
- 🧪 **Comprehensive Testing** with in-memory database
- 📋 **Complete Documentation** with API examples
- ⚙️ **Environment-based Configuration** (dev/prod/test)

---

## 📁 Project Structure

```
gin-template/
├── pkg/                    # Reusable packages
│   ├── response/          # Standardized API responses
│   └── logger/            # Structured logging
├── internal/               # Private application code
│   ├── auth/              # JWT authentication utilities
│   ├── config/            # Configuration management
│   ├── database/          # Database initialization and utilities
│   ├── handlers/          # HTTP controllers and business logic
│   ├── middlewares/       # Custom middlewares (auth, rate limiting, etc.)
│   ├── models/            # Data models (GORM)
│   ├── routes/            # Route definitions and registration
│   └── validators/        # Input validation logic
├── cmd/api/               # Application entrypoint
│   └── main.go           # Main application file
├── configs/               # Environment-specific configurations
│   ├── development.yaml  # Development settings
│   └── production.yaml   # Production settings
├── docs/                  # Documentation
│   └── api.md            # API documentation
├── scripts/               # Utility scripts
│   └── setup-db.sh      # Database setup script
├── tests/                 # Integration and load tests
├── data/                  # Database files (SQLite)
├── .env.example          # Environment variables template
├── .golangci.yml         # Linter configuration
├── Makefile              # Development and deployment commands
├── Dockerfile            # Optimized Docker container
├── docker-compose.yaml   # Multi-service development environment
├── go.mod / go.sum       # Go dependencies
└── README.md             # This documentation
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+ installed
- Docker and Docker Compose (optional, for containerized development)
- Make (for using Makefile commands)

### Quick Start

#### 1. Clone and Setup
```bash
git clone <repo-url>
cd gin-template

# Copy environment template
cp .env.example .env
```

#### 2. Configure Environment
Edit `.env` file with your settings:
```bash
# Required: Change this for production!
JWT_SECRET=your-super-secret-jwt-key

# Database (SQLite for development)
DB_DRIVER=sqlite
DB_DSN=./data/app.db

# Server
PORT=8080
APP_ENV=development

# Logging
LOG_LEVEL=info
LOG_FORMAT=text
```

#### 3. Install Dependencies and Tools
```bash
# Install Go dependencies
go mod tidy

# Install development tools (optional but recommended)
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 4. Setup Database
```bash
# Create data directory for SQLite
mkdir -p data

# Or use the setup script
./scripts/setup-db.sh
```

#### 5. Run the Application

**Option A: Direct Go execution**
```bash
make run
# or
go run ./cmd/api/main.go
```

**Option B: Docker Compose (Recommended for development)**
```bash
# Start all services (API + PostgreSQL + pgAdmin)
make up-build

# View logs
make logs

# Stop services
make down
```

#### 6. Test the API
```bash
# Health check
curl http://localhost:8080/health/

# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com", "password": "SecurePass123!"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "SecurePass123!"}'
```

---

## Useful Commands (Makefile)

- `make fmt`   — Format and organize imports.
- `make lint`  — Run linters and best practice checks.
- `make test`  — Run all tests.
- `make run`   — Run the application.
- `make all`   — Format, lint, and test in one command.

### Docker Compose

- `make build`        — Build Docker images.
- `make up`           — Start all services.
- `make up-build`     — Build and start all services.
- `make down`         — Stop all services.
- `make restart`      — Restart all services.
- `make logs`         — Show logs for all services.
- `make logs-api`     — Show logs for the API service.
- `make logs-db`      — Show logs for the PostgreSQL service.
- `make clean`        — Remove containers, images, and volumes.
- `make db-reset`     — Reset the database (WARNING: deletes all data).
- `make db-backup`    — Create a database backup.
- `make db-restore`   — Restore a database backup.
- `make shell-api`    — Access the API container shell.
- `make shell-db`     — Access the PostgreSQL shell.
- `make admin`        — Start with pgAdmin included.
- `make status`       — Show service status.
- `make health`       — Check health of services.

---

## 🌐 API Endpoints

### Public Endpoints
- `GET /health/` — Complete health check with service status
- `GET /health/live` — Kubernetes liveness probe
- `GET /health/ready` — Kubernetes readiness probe
- `POST /api/auth/register` — User registration (enhanced validation)
- `POST /api/auth/login` — User authentication (returns JWT + user info)

### Protected Endpoints (Require JWT)
- `GET /api/protected/` — Example protected resource
- `GET /api/users/me` — Current user profile

### Legacy Endpoints (Backward Compatibility)
- `POST /api/register` — User registration
- `POST /api/login` — User authentication

### Security Features
- **Rate Limiting**: 10 req/sec for general API, 5 req/min for auth endpoints
- **Input Validation**: Comprehensive password requirements and email validation
- **Security Headers**: OWASP-recommended headers automatically applied
- **Request Tracking**: Unique request IDs for debugging and monitoring

---

## Best Practices and Included Tools

- **Modular and scalable structure**
- **Environment variables** with `.env` and `.env.example`
- **Linters and formatters** (`golangci-lint`, `goimports`)
- **Makefile** for automating development tasks
- **Code comments and documentation**
- **Ready for CI/CD integration**
- **Support for SQLite, PostgreSQL, and MySQL (GORM)**
- **Docker and Docker Compose** for local development and deployment

---

## Testing

Add your tests in `*_test.go` files within the corresponding packages.
Run all tests with:

```sh
make test
```

---

## Continuous Integration (CI/CD)

You can easily add a GitHub Actions workflow, for example:

```yaml
# .github/workflows/ci.yml
name: Go CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Install linters
        run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
      - name: Lint
        run: make lint
      - name: Test
        run: make test
```

---

## Contributing

Contributions are welcome!
Please open an issue or pull request following the project's best practices.

---

## License

MIT

---
