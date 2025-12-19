# Production Go REST API

A production-ready REST API built with Go featuring a modular architecture, comprehensive testing, PostgreSQL integration, structured logging, graceful shutdown, and full CI/CD pipeline.

## 🚀 Features

- **Modular Architecture**: Clean separation with `cmd/`, `internal/`, and `pkg/` directories
- **5 RESTful Endpoints**: Full CRUD operations with pagination
- **PostgreSQL Database**: With connection pooling and health checks
- **Structured Logging**: JSON logging with configurable levels using Zap
- **Graceful Shutdown**: Proper cleanup of resources on termination
- **Docker Support**: Multi-stage Dockerfile for optimized images
- **Health Checks**: Built-in health and readiness probes
- **Comprehensive Tests**: Unit tests with coverage reporting
- **CI/CD Pipeline**: GitHub Actions for testing, linting, building, and deployment
- **Environment Configuration**: 12-factor app principles with .env support

## 📋 Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [CI/CD](#cicd)

## 🏗️ Architecture

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection and migrations
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # HTTP middleware (logging, recovery, CORS)
│   ├── models/                  # Data models
│   ├── repository/              # Data access layer
│   └── server/                  # HTTP server setup
├── pkg/
│   ├── logger/                  # Structured logging
│   └── response/                # Standardized API responses
├── .github/
│   └── workflows/
│       └── ci-cd.yml           # CI/CD pipeline
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml          # Local development environment
├── Makefile                    # Development commands
└── .env.example                # Environment variables template
```

## 📦 Prerequisites

- Go 1.24 or higher
- PostgreSQL 16 or higher (or use Docker Compose)
- Docker and Docker Compose (optional, for containerized deployment)
- Make (optional, for convenience commands)

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/kimbef/19-dec-2025.git
cd 19-dec-2025
```

2. Start the services:
```bash
docker compose up -d
```

3. Verify the API is running:
```bash
curl http://localhost:8080/health
```

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/kimbef/19-dec-2025.git
cd 19-dec-2025
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. Start PostgreSQL (if not using Docker):
```bash
# Using Docker
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=apidb \
  -p 5432:5432 \
  postgres:16-alpine
```

5. Run the application:
```bash
make run
# or
go run ./cmd/api
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

### Endpoints

#### Health Check

```http
GET /health
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "services": {
      "api": "up",
      "database": "up"
    },
    "uptime": "1h23m45s"
  }
}
```

#### Create Item

```http
POST /api/items
Content-Type: application/json

{
  "name": "Laptop",
  "description": "High-performance laptop",
  "quantity": 10,
  "price": 1299.99
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Laptop",
    "description": "High-performance laptop",
    "quantity": 10,
    "price": 1299.99,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get Item

```http
GET /api/items/{id}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Laptop",
    "description": "High-performance laptop",
    "quantity": 10,
    "price": 1299.99,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update Item

```http
PUT /api/items/{id}
Content-Type: application/json

{
  "quantity": 15,
  "price": 1199.99
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Laptop",
    "description": "High-performance laptop",
    "quantity": 15,
    "price": 1199.99,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Delete Item

```http
DELETE /api/items/{id}
```

**Response:**
```
204 No Content
```

#### List Items

```http
GET /api/items?page=1&per_page=10
```

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Laptop",
        "description": "High-performance laptop",
        "quantity": 10,
        "price": 1299.99,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "per_page": 10,
    "total_pages": 1
  }
}
```

### Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
```

Error codes:
- `BAD_REQUEST` (400): Invalid request
- `NOT_FOUND` (404): Resource not found
- `INTERNAL_ERROR` (500): Server error

## ⚙️ Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env` and adjust as needed:

```bash
# Server Configuration
PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_SHUTDOWN_TIMEOUT=30s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=apidb
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Logging Configuration
LOG_LEVEL=info          # debug, info, warn, error
LOG_FORMAT=json         # json, console
```

## 🛠️ Development

### Available Make Commands

```bash
make help           # Display all available commands
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage report
make lint           # Run linter
make fmt            # Format code
make tidy           # Tidy dependencies
make clean          # Clean build artifacts

make docker-build   # Build Docker image
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-logs    # View Docker logs
```

### Code Structure Guidelines

- **Handlers**: HTTP request handlers in `internal/handlers/`
- **Repository**: Database operations in `internal/repository/`
- **Models**: Data structures in `internal/models/`
- **Middleware**: HTTP middleware in `internal/middleware/`
- **Shared packages**: Reusable code in `pkg/`

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates:
- `coverage.out`: Coverage data file
- `coverage.html`: HTML coverage report

### Run Specific Tests

```bash
go test -v ./internal/handlers/
go test -v ./pkg/response/
```

## 🚢 Deployment

### Docker Deployment

1. Build the Docker image:
```bash
make docker-build
```

2. Run with Docker Compose:
```bash
docker compose up -d
```

3. Check logs:
```bash
docker compose logs -f api
```

### Production Deployment

1. Set production environment variables
2. Build the application:
```bash
CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
```

3. Deploy the binary with your preferred method (systemd, Kubernetes, etc.)

### Kubernetes Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-rest-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-rest-api
  template:
    metadata:
      labels:
        app: go-rest-api
    spec:
      containers:
      - name: api
        image: your-registry/go-rest-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

## 🔄 CI/CD

The project includes a comprehensive GitHub Actions workflow that:

1. **Testing**: Runs all tests with PostgreSQL service
2. **Linting**: Uses golangci-lint for code quality
3. **Coverage**: Uploads coverage to Codecov
4. **Building**: Builds the application
5. **Docker**: Builds and pushes Docker images
6. **Security**: Scans for vulnerabilities with Trivy
7. **Deployment**: Deploys to production on main branch

### Setup GitHub Secrets

Configure these secrets in your GitHub repository:

- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub password/token
- Additional deployment secrets as needed

### Workflow Triggers

- **Push**: Runs on pushes to `main` and `develop` branches
- **Pull Request**: Runs on PRs to `main` and `develop` branches

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📧 Support

For issues and questions, please open an issue on GitHub.

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- PostgreSQL driver: [lib/pq](https://github.com/lib/pq)
- Logging: [Zap](https://github.com/uber-go/zap)

