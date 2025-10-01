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
- **Mockery** (mock generation for testing)
- **golangci-lint** (code linting and static analysis)
- **Makefile** (command automation)

## Quick Start

### Using Makefile
```bash
# Run application locally (starts containers + API)
make run-local

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

## Database
- PostgreSQL is used as the main database
- Automatic migrations via Liquibase (`liquibase/changelog/migrations/*.sql`)

## Documentation

📚 **Detailed Documentation:**
- **[Authentication (JWT)](./docs/AUTHENTICATION.md)** - JWT token generation, validation, and usage
- **[API Documentation](./docs/API.md)** - Endpoints, examples, and Postman collection
- **[Development Guide](./docs/DEVELOPMENT.md)** - Testing, linting, and mock generation
- **[Observability](./docs/OBSERVABILITY.md)** - Monitoring, tracing, and rate limiting

## Health Check
```bash
curl http://localhost:8080/health
```

## Getting Help
- Check the Makefile: `make help`
- Review the documentation files in `/docs`
- Import the Postman collection for API examples
