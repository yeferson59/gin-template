# my-api

## Description

**my-api** is a professional template for building robust APIs in Go using the Gin framework. It includes JWT authentication, modular architecture, database connection with GORM, environment variable management, middlewares, and best development practices.

---

## Project Structure

```
my-api/
├── cmd/
│   └── api/
│       └── main.go         # Main entrypoint
├── internal/
│   ├── auth/               # JWT and authentication utilities
│   ├── config/             # Configuration and environment variable management
│   ├── database/           # Database initialization and connection
│   ├── handlers/           # HTTP controllers
│   ├── middlewares/        # Custom middlewares
│   ├── models/             # Data models (GORM)
│   ├── modules/            # Additional modules (scalable)
│   └── routes/             # Main route registration
├── data/                   # Database files (e.g., SQLite)
├── .env                    # Environment variables (do not commit to git)
├── .env.example            # Environment variable template
├── .golangci.yml           # Linter configuration
├── Makefile                # Development commands
├── go.mod / go.sum         # Go dependencies
├── Dockerfile              # Docker build file
├── .dockerignore           # Docker context ignore file
├── .gitignore              # Git ignore file
└── README.md               # This file
```

---

## Getting Started

### 1. Clone the repository and enter the directory

```sh
git clone <repo-url>
cd my-api
```

### 2. Copy and edit the environment file

```sh
cp .env.example .env
# Edit .env according to your environment (port, database, JWT_SECRET, etc.)
```

### 3. Install recommended tools

```sh
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

Make sure `~/go/bin` is in your `$PATH`.

### 4. Install dependencies

```sh
go mod tidy
```

### 5. Run the application

```sh
make run
# or directly
go run ./cmd/api/main.go
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

## Main Endpoints

- `POST /api/register` — User registration (requires username, email, password)
- `POST /api/login`    — User login (returns JWT)
- `GET /api/protected` — Protected endpoint, requires JWT in Authorization header

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
