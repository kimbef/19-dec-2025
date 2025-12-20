# CI/CD Pipeline Application

A production-ready Go application demonstrating best practices for CI/CD pipeline implementation with comprehensive testing, monitoring, and deployment automation.

[![CI/CD Pipeline](https://github.com/kimbef/19-dec-2025/actions/workflows/ci.yml/badge.svg)](https://github.com/kimbef/19-dec-2025/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kimbef/19-dec-2025)](https://goreportcard.com/report/github.com/kimbef/19-dec-2025)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [CI/CD Pipeline](#cicd-pipeline)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **RESTful API**: Built with Gin web framework
- **Database Integration**: PostgreSQL with automated migrations
- **Comprehensive Testing**: Unit and integration tests with high coverage
- **Structured Logging**: Using Uber's Zap logger
- **Graceful Shutdown**: Proper signal handling and cleanup
- **Health Checks**: Readiness and liveness probes for orchestration
- **Configuration Management**: Environment-based configuration
- **Docker Support**: Multi-stage Dockerfile for optimized images
- **CI/CD Ready**: GitHub Actions workflow with multi-version testing
- **Code Quality**: Integrated linting with golangci-lint
- **Security Scanning**: Automated vulnerability detection

## 🏗 Architecture

```
.
├── cmd/
│   └── api/              # Application entrypoint
├── internal/
│   ├── api/              # HTTP handlers and routes
│   ├── config/           # Configuration management
│   ├── database/         # Database connection and migrations
│   └── models/           # Data models
├── pkg/
│   ├── logger/           # Structured logging
│   └── middleware/       # HTTP middlewares
├── .github/
│   └── workflows/        # CI/CD pipeline definitions
├── deployments/          # Deployment configurations
├── scripts/              # Utility scripts
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Local development environment
└── Makefile              # Build automation
```

## 📦 Prerequisites

- **Go**: 1.20 or higher
- **Docker**: For containerization (optional)
- **PostgreSQL**: 13 or higher (or use Docker Compose)
- **Make**: For build automation

## 🚀 Quick Start

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/kimbef/19-dec-2025.git
   cd 19-dec-2025
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL** (using Docker)
   ```bash
   docker run -d \
     --name postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=cicd_app \
     -p 5432:5432 \
     postgres:15-alpine
   ```

4. **Run the application**
   ```bash
   make run
   ```

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📚 API Documentation

### Health Endpoints

#### Health Check
```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-12-19T19:47:38Z",
  "version": "1.0.0",
  "checks": {
    "database": "healthy"
  }
}
```

#### Readiness Check
```
GET /ready
```

Response:
```json
{
  "ready": true,
  "message": "application is ready"
}
```

### Pipeline Endpoints

#### List Pipelines
```
GET /api/v1/pipelines
```

#### Get Pipeline
```
GET /api/v1/pipelines/:id
```

#### Create Pipeline
```
POST /api/v1/pipelines
Content-Type: application/json

{
  "name": "My Pipeline",
  "description": "Pipeline description"
}
```

#### Update Pipeline
```
PUT /api/v1/pipelines/:id
Content-Type: application/json

{
  "name": "Updated Pipeline",
  "description": "Updated description",
  "status": "active"
}
```

#### Delete Pipeline
```
DELETE /api/v1/pipelines/:id
```

### Build Endpoints

#### List Builds
```
GET /api/v1/pipelines/:pipeline_id/builds
```

#### Create Build
```
POST /api/v1/builds
Content-Type: application/json

{
  "pipeline_id": 1,
  "branch": "main",
  "commit_hash": "abc123def456"
}
```

### Deployment Endpoints

#### Create Deployment
```
POST /api/v1/deployments
Content-Type: application/json

{
  "build_id": 1,
  "environment": "production",
  "deployed_by": "user@example.com"
}
```

## 🛠 Development

### Build

```bash
# Build binary
make build

# Build Docker image
make docker-build
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View coverage report
open coverage.html
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

### Database Migrations

Migrations run automatically on application startup. The following tables are created:

- **pipelines**: CI/CD pipeline definitions
- **builds**: Build/run history
- **deployments**: Deployment records

## 🚢 Deployment

### Docker Deployment

```bash
# Build image
docker build -t cicd-pipeline-app:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  cicd-pipeline-app:latest
```

### Kubernetes Deployment

See `deployments/k8s/` directory for example Kubernetes manifests.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `ENVIRONMENT` | Environment name | `development` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `cicd_app` |
| `DB_SSLMODE` | Database SSL mode | `require` |
| `LOG_LEVEL` | Logging level | `info` |
| `LOG_FORMAT` | Log format (json/text) | `json` |

**Security note:** For non-local environments, use a secure SSL mode such as `require` or `verify-full`. The `disable` mode should only be used for local development with trusted connections (for example, when connecting to a PostgreSQL instance on `localhost`).
## 🔄 CI/CD Pipeline

The GitHub Actions workflow automatically:

1. **Tests** across multiple Go versions (1.20, 1.21, 1.22)
2. **Lints** code with golangci-lint
3. **Builds** the application
4. **Generates** code coverage reports
5. **Scans** for security vulnerabilities
6. **Builds and pushes** Docker images
7. **Deploys** to target environments (configurable)

### Pipeline Stages

- **Test**: Runs tests with PostgreSQL service
- **Lint**: Code quality checks
- **Build**: Compiles application
- **Docker**: Builds and pushes images
- **Security**: Vulnerability scanning
- **Deploy**: Example deployment configuration

### Required Secrets

For Docker Hub deployment:
- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub password or token

## 📖 Make Commands

```bash
make help              # Show all available commands
make build             # Build application binary
make test              # Run tests
make test-coverage     # Run tests with coverage report
make clean             # Clean build artifacts
make run               # Run application locally
make docker-build      # Build Docker image
make docker-compose-up # Start services with docker-compose
make lint              # Run linter
make fmt               # Format code
make vet               # Run go vet
make deps              # Download dependencies
make all               # Run all checks and build
```

## 🤝 Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Gin Web Framework](https://github.com/gin-gonic/gin)
- Logging with [Uber Zap](https://github.com/uber-go/zap)
- PostgreSQL driver [lib/pq](https://github.com/lib/pq)

## 📞 Support

For issues and questions:
- Create an [Issue](https://github.com/kimbef/19-dec-2025/issues)
- Submit a [Pull Request](https://github.com/kimbef/19-dec-2025/pulls)

---

**Built with ❤️ for CI/CD excellence**
