# Exchange Crypto Service API

A API for cryptocurrency exchange operations, built in Go, following Clean Architecture principles with integrated observability.

## Tech Stack
- **Go** (>=1.24)
- **PostgreSQL** (data persistence)
- **Liquibase** (database migration)
- **Docker & Docker Compose** (containerization and orchestration)
- **OpenTelemetry** (tracing and observability)
- **JWT Authentication** (secure token-based authentication)
- **Rate Limiting** (DDoS protection and traffic control)
- **Uber/FX** (dependency injection framework)
- **Mockery** (mock generation for testing)
- **golangci-lint** (code linting and static analysis)
- **Makefile** (command automation)

## Quick Start

### Using Makefile
```bash
# Run application locally (starts containers + API)
make run-local

# Run complete application via Docker
make run-docker

# Fresh start (rebuild containers and run migrations)
make fs

# Build containers
make build

# Run database migrations
make migrate

# Run tests
make test

# Run code linting
make lint-code

# Stop all services
make down

# Show help with available commands
make help
```

## 🐳 Docker

This application is fully containerized using Docker Compose. You can run the entire stack (API + dependencies) with a single command.

### Quick Docker Commands
```bash
# Start everything via Docker
make run-docker

# View logs and manage containers
make logs
make status
make restart
```

**📋 For detailed Docker documentation, see [Docker Guide](./docs/DOCKER.md)**

## Project Structure
- `cmd/api/main.go`: application entry point
- `internal/app`: business modules (account, exchange, transaction, user)
- `internal/config`: configurations
- `internal/database`: database connection
- `liquibase/`: database scripts and migrations
- `docker-compose.yml`: services orchestration
- `Makefile`: automated commands

## Configuration
- Environment variables can be defined in `env.yaml`

### Environment Variables

| Variable | Description | Default | Values |
|----------|-------------|---------|--------|
| `FX_VERBOSE` | Enable verbose Uber/FX dependency injection logs | `false` | `true`, `false` |
| `DB_URL` | Database host | `postgres` | hostname |
| `DB_PORT` | Database port | `5432` | port number |
| `DB_APP_USER` | Database user | `owner` | username |
| `DB_APP_PASSWORD` | Database password | `owner123` | password |
| `DB_NAME` | Database name | `exchange_crypto_local` | database name |
| `DB_SSL_MODE` | PostgreSQL SSL mode | `disable` | `disable`, `require` |
| `DB_SCHEMA` | Database schema | `exchange_crypto` | schema name |

**Example: Enable FX verbose logs for debugging**
```bash
FX_VERBOSE=true go run cmd/api/main.go
# or
FX_VERBOSE=true make run-local
```

## Database
- PostgreSQL is used as the main database
- Automatic migrations via Liquibase (`liquibase/changelog/migrations/*.sql`)

## Documentation

📚 **Detailed Documentation:**
- **[Dependency Injection (Uber/FX)](./docs/DEPENDENCY_INJECTION.md)** - DI architecture and module structure
- **[Authentication (JWT)](./docs/AUTHENTICATION.md)** - JWT token generation, validation, and usage
- **[API Documentation](./docs/API.md)** - Endpoints, examples, and Postman collection
- **[Development Guide](./docs/DEVELOPMENT.md)** - Testing, linting, and mock generation
- **[Docker Guide](./docs/DOCKER.md)** - Complete Docker setup, troubleshooting, and deployment
- **[Observability](./docs/OBSERVABILITY.md)** - Monitoring, tracing, and rate limiting

## Health Check
```bash
curl http://localhost:8080/health
```

## Getting Help
- Check the Makefile: `make help`
- Review the documentation files in `/docs`
- Import the Postman collection for API examples
